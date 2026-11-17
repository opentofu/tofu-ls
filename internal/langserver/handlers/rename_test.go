// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"fmt"
	"github.com/creachadair/jrpc2"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/opentofu/tofu-ls/internal/langserver"
	"github.com/opentofu/tofu-ls/internal/state"
	"github.com/opentofu/tofu-ls/internal/tofu/exec"
	"github.com/opentofu/tofu-ls/internal/walker"
	"github.com/stretchr/testify/mock"
)

func TestRename_variable(t *testing.T) {
	tmpDir := TempDir(t)

	ss, err := state.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	wc := walker.NewWalkerCollector()

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TofuCalls: &exec.TofuMockCalls{
			PerWorkDir: map[string][]*mock.Call{
				tmpDir.Path(): {
					{
						Method:        "Version",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							version.Must(version.NewVersion("0.12.0")),
							nil,
							nil,
						},
					},
					{
						Method:        "GetExecPath",
						Repeatability: 1,
						ReturnArguments: []interface{}{
							"",
						},
					},
				},
			},
		},
		StateStore:      ss,
		WalkerCollector: wc,
	}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {
	    	"definition": {
	    		"linkSupport": true
	    	}
	    },
	    "rootUri": %q,
	    "processId": 12345
	}`, tmpDir.URI)})
	waitForWalkerPath(t, ss, wc, tmpDir)
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "opentofu",
			"text": `+fmt.Sprintf("%q",
			`variable "old" {
}

locals {
  x = 1
}

output "foo" {
  value = "${var.old}-${local.x}"
}

output "bar" {
  value = var.old
}`)+`,
			"uri": "%s/main.tf"
		}
	}`, tmpDir.URI)})
	waitForAllJobs(t, ss)

	edits := func(rows ...[4]int) string {
		s := ""
		for i, r := range rows {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprintf(`{"range":{"start":{"line":%d,"character":%d},"end":{"line":%d,"character":%d}},"newText":"new"}`, r[0], r[1], r[0], r[2])
		}
		return s
	}
	want := fmt.Sprintf(`{"jsonrpc":"2.0","id":%%d,"result":{"changes":{"%s/main.tf":[%s]}}}`, tmpDir.URI,
		edits([4]int{0, 10, 13, 0}, [4]int{8, 17, 20, 0}, [4]int{12, 14, 17, 0}))

	// from a reference
	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/rename",
		ReqParams: fmt.Sprintf(`{"textDocument":{"uri":"%s/main.tf"},"position":{"line":12,"character":15},"newName":"new"}`, tmpDir.URI),
	}, fmt.Sprintf(want, 3))

	// from the definition
	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/rename",
		ReqParams: fmt.Sprintf(`{"textDocument":{"uri":"%s/main.tf"},"position":{"line":0,"character":11},"newName":"new"}`, tmpDir.URI),
	}, fmt.Sprintf(want, 4))

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/prepareRename",
		ReqParams: fmt.Sprintf(`{"textDocument":{"uri":"%s/main.tf"},"position":{"line":12,"character":15}}`, tmpDir.URI),
	}, `{"jsonrpc":"2.0","id":5,"result":{"range":{"start":{"line":12,"character":14},"end":{"line":12,"character":17}},"placeholder":"old"}}`)

	// not a name
	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method:    "textDocument/prepareRename",
		ReqParams: fmt.Sprintf(`{"textDocument":{"uri":"%s/main.tf"},"position":{"line":11,"character":2}}`, tmpDir.URI),
	}, &jrpc2.Error{Code: -32098, Message: errNotRenamable.Error()})
}
