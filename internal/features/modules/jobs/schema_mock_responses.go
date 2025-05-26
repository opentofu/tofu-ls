// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jobs

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	tfregistry "github.com/opentofu/opentofu-schema/registry"
	"github.com/zclconf/go-cty/cty"
)

// puppetModuleVersionsMockResponse represents response from https://registry.terraform.io/v1/modules/puppetlabs/deployment/ec/versions
var puppetModuleVersionsMockResponse = `{
  "modules": [
    {
      "source": "puppetlabs/deployment/ec",
      "versions": [
        {
          "version": "0.0.5",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.0.6",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.0.8",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.0.2",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.0.1",
          "root": {
            "providers": [],
            "dependencies": []
          },
          "submodules": [
            {
              "path": "modules/ec-deployment",
              "providers": [
                {
                  "name": "ec",
                  "namespace": "",
                  "source": "elastic/ec",
                  "version": "0.2.1"
                }
              ],
              "dependencies": []
            }
          ]
        },
        {
          "version": "0.0.4",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.0.3",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.0.7",
          "root": {
            "providers": [
              {
                "name": "ec",
                "namespace": "",
                "source": "elastic/ec",
                "version": "0.2.1"
              }
            ],
            "dependencies": []
          },
          "submodules": []
        }
      ]
    }
  ]
}`

// puppetModuleDataMockResponse represents response from https://api.opentofu.org/registry/docs/modules/puppetlabs/deployment/ec/v0.0.8/index.json
var puppetModuleDataMockResponse = `{
  "id": "v0.0.8",
  "published": "2021-08-05T00:26:01Z",
  "readme": true,
  "edit_link": "https://github.com/puppetlabs/terraform-ec-deployment/blob/v0.0.8/README.md",
  "variables": {
    "autoscale": {
      "type": "string",
      "default": "true",
      "description": "Enable autoscaling of elasticsearch",
      "sensitive": false,
      "required": false
    },
    "deployment_templateid": {
      "type": "string",
      "default": "gcp-io-optimized",
      "description": "ID of Elastic Cloud deployment type",
      "sensitive": false,
      "required": false
    },
    "ec_region": {
      "type": "string",
      "default": "gcp-us-west1",
      "description": "cloud provider region",
      "sensitive": false,
      "required": false
    },
    "ec_stack_version": {
      "type": "string",
      "default": "",
      "description": "Version of Elastic Cloud stack to deploy",
      "sensitive": false,
      "required": false
    },
    "name": {
      "type": "string",
      "default": "ecproject",
      "description": "Name of resources",
      "sensitive": false,
      "required": false
    },
    "traffic_filter_sourceip": {
      "type": "string",
      "default": "",
      "description": "traffic filter source IP",
      "sensitive": false,
      "required": false
    }
  },
  "outputs": {
    "deployment_id": {
      "sensitive": false,
      "description": "Elastic Cloud deployment ID"
    },
    "elasticsearch_cloud_id": {
      "sensitive": false,
      "description": "Elastic Cloud project deployment ID"
    },
    "elasticsearch_https_endpoint": {
      "sensitive": false,
      "description": "elasticsearch https endpoint"
    },
    "elasticsearch_password": {
      "sensitive": true,
      "description": "elasticsearch password"
    },
    "elasticsearch_username": {
      "sensitive": false,
      "description": "elasticsearch username"
    },
    "elasticsearch_version": {
      "sensitive": false,
      "description": "Stack version deployed"
    }
  },
  "schema_error": "",
  "providers": [],
  "dependencies": [],
  "resources": [
    {
      "address": "ec_deployment.ecproject",
      "type": "ec_deployment",
      "name": "ecproject"
    },
    {
      "address": "ec_deployment_traffic_filter.gcp_vpc_nat",
      "type": "ec_deployment_traffic_filter",
      "name": "gcp_vpc_nat"
    },
    {
      "address": "ec_deployment_traffic_filter_association.ec_tf_association",
      "type": "ec_deployment_traffic_filter_association",
      "name": "ec_tf_association"
    },
    {
      "address": "data.ec_stack.latest",
      "type": "ec_stack",
      "name": "latest"
    }
  ],
  "link": "https://github.com/puppetlabs/terraform-ec-deployment/tree/v0.0.8",
  "vcs_repository": "",
  "licenses": [
    {
      "spdx": "Apache-2.0",
      "confidence": 0.9876943,
      "is_compatible": true,
      "file": "LICENSE",
      "link": "https://github.com/puppetlabs/terraform-ec-deployment/blob/v0.0.8/LICENSE"
    }
  ],
  "incompatible_license": false,
  "examples": {

  },
  "submodules": {

  }
}`

// labelNullModuleVersionsMockResponse represents response for
// versions of module that suffers from "unreliable" input data, as described in
// https://github.com/hashicorp/vscode-terraform/issues/1582
// It is a shortened response from https://registry.terraform.io/v1/modules/cloudposse/label/null/versions
var labelNullModuleVersionsMockResponse = `{
  "modules": [
    {
      "source": "cloudposse/label/null",
      "versions": [
        {
          "version": "0.25.0",
          "root": {
            "providers": [],
            "dependencies": []
          },
          "submodules": []
        },
        {
          "version": "0.26.0",
          "root": {
            "providers": [],
            "dependencies": []
          },
          "submodules": []
        }
      ]
    }
  ]
}`

// labelNullModuleDataOldMockResponse represents response for
// a module that does NOT suffer from "unreliable" input data,
// as described in https://github.com/hashicorp/vscode-terraform/issues/1582
// This is for comparison with the unreliable input data.
// It is a shortened response from https://api.opentofu.org/registry/docs/modules/cloudposse/label/null/index.json
var labelNullModuleDataMockResponse = `{
  "addr": {
    "display": "cloudposse/label/null",
    "namespace": "cloudposse",
    "name": "label",
    "target": "null"
  },
  "description": "",
  "versions": [
    {
      "id": "v0.25.0",
      "published": "2021-08-25T17:45:16Z"
    },
    {
      "id": "v0.24.0",
      "published": "2021-02-04T08:11:56Z"
    }
  ],
  "is_blocked": false,
  "popularity": 0,
  "fork_count": 0,
  "fork_of": {
    "display": "//",
    "namespace": "",
    "name": "",
    "target": ""
  },
  "upstream_popularity": 0,
  "upstream_fork_count": 0
}`

var labelNullExpectedOldModuleData = &tfregistry.ModuleData{
	Version: version.Must(version.NewVersion("0.25.0")),
	Inputs: []tfregistry.Input{
		{
			Name:        "environment",
			Type:        cty.String,
			Description: lang.Markdown(""),
		},
		{
			Name:        "label_order",
			Type:        cty.DynamicPseudoType,
			Description: lang.Markdown(""),
		},
		{
			Name:        "descriptor_formats",
			Type:        cty.DynamicPseudoType,
			Description: lang.Markdown(""),
		},
	},
	Outputs: []tfregistry.Output{
		{
			Name:        "id",
			Description: lang.Markdown(""),
		},
	},
}

var labelNullExpectedNewModuleData = &tfregistry.ModuleData{
	Version: version.Must(version.NewVersion("0.26.0")),
	Inputs: []tfregistry.Input{
		{
			Name:        "environment",
			Type:        cty.String,
			Description: lang.Markdown(""),
			Required:    true,
		},
		{
			Name:        "label_order",
			Type:        cty.DynamicPseudoType,
			Description: lang.Markdown(""),
			Default:     cty.NullVal(cty.DynamicPseudoType),
		},
		{
			Name:        "descriptor_formats",
			Type:        cty.DynamicPseudoType,
			Description: lang.Markdown(""),
		},
	},
	Outputs: []tfregistry.Output{
		{
			Name:        "id",
			Description: lang.Markdown(""),
		},
	},
}
