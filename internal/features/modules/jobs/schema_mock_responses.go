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

// puppetModuleVersionsMockResponse represents the response from https://api.opentofu.org/registry/docs/modules/puppetlabs/deployment/ec/index.json
var puppetModuleVersionsMockResponse = `{
  "addr": {
    "display": "puppetlabs/deployment/ec",
    "namespace": "puppetlabs",
    "name": "deployment",
    "target": "ec"
  },
  "description": "",
  "versions": [
    {
      "id": "v0.0.8",
      "published": "2021-08-05T00:26:01Z"
    },
    {
      "id": "v0.0.7",
      "published": "2021-08-03T19:57:07Z"
    },
    {
      "id": "v0.0.6",
      "published": "2021-08-03T19:47:41Z"
    },
    {
      "id": "v0.0.5",
      "published": "2021-08-03T18:56:01Z"
    },
    {
      "id": "v0.0.4",
      "published": "2021-08-03T18:39:44Z"
    },
    {
      "id": "v0.0.3",
      "published": "2021-08-02T21:54:06Z"
    },
    {
      "id": "v0.0.2",
      "published": "2021-08-02T21:48:35Z"
    },
    {
      "id": "v0.0.1",
      "published": "2021-08-02T21:18:58Z"
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
// It is a shortened response from https://api.opentofu.org/registry/docs/modules/cloudposse/label/null/index.json
var labelNullModuleVersionsMockResponse = `{
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

// labelNullModuleDataMockResponse represents response for
// a module that does NOT suffer from "unreliable" input data,
// as described in https://github.com/hashicorp/vscode-terraform/issues/1582
// This is for comparison with the unreliable input data.
// It is a shortened response from https://api.opentofu.org/registry/docs/modules/cloudposse/label/null/v0.25.0/index.json
var labelNullModuleDataMockResponse = `{
  "id": "v0.25.0",
  "published": "2021-08-25T17:45:16Z",
  "readme": true,
  "edit_link": "https://github.com/cloudposse/terraform-null-label/blob/0.25.0/README.md",
  "variables": {
    "enabled": {
      "type": "bool",
      "default": null,
      "description": "Set to false to prevent the module from creating any resources",
      "sensitive": false,
      "required": false
    },
    "environment": {
      "type": "string",
      "default": null,
      "description": "ID element. Usually used for region e.g. 'uw2', 'us-west-2', OR role 'prod', 'staging', 'dev', 'UAT'",
      "sensitive": false,
      "required": false
    }
  },
  "outputs": {
    "additional_tag_map": {
      "sensitive": false,
      "description": "The merged additional_tag_map"
    },
    "attributes": {
      "sensitive": false,
      "description": "List of attributes"
    }
  },
  "schema_error": "",
  "providers": [],
  "dependencies": [],
  "resources": [],
  "link": "https://github.com/cloudposse/terraform-null-label/tree/0.25.0",
  "vcs_repository": "",
  "licenses": [
    {
      "spdx": "Apache-2.0",
      "confidence": 0.98704666,
      "is_compatible": true,
      "file": "LICENSE",
      "link": "https://github.com/cloudposse/terraform-null-label/blob/0.25.0/LICENSE"
    }
  ],
  "incompatible_license": false,
  "examples": {
    "autoscalinggroup": {
      "readme": false,
      "edit_link": "",
      "variables": {
        "enabled": {
          "type": "bool",
          "default": null,
          "description": "Set to false to prevent the module from creating any resources",
          "sensitive": false,
          "required": false
        },
        "environment": {
          "type": "string",
          "default": null,
          "description": "Environment, e.g. 'uw2', 'us-west-2', OR 'prod', 'staging', 'dev', 'UAT'",
          "sensitive": false,
          "required": false
        },
        "label_order": {
          "type": "list of string",
          "default": null,
          "description": "The naming order of the id output and Name tag.\nDefaults to [\"namespace\", \"environment\", \"stage\", \"name\", \"attributes\"].\nYou can omit any of the 5 elements, but at least one must be present.\n",
          "sensitive": false,
          "required": false
        },
        "name": {
          "type": "string",
          "default": null,
          "description": "Solution name, e.g. 'app' or 'jenkins'",
          "sensitive": false,
          "required": false
        },
        "namespace": {
          "type": "string",
          "default": null,
          "description": "Namespace, which could be your organization name or abbreviation, e.g. 'eg' or 'cp'",
          "sensitive": false,
          "required": false
        }
      },
      "outputs": {
        "id": {
          "sensitive": false,
          "description": ""
        },
        "tags": {
          "sensitive": false,
          "description": ""
        },
        "tags_as_list_of_maps": {
          "sensitive": false,
          "description": ""
        }
      },
      "schema_error": ""
    },
    "complete": {
      "readme": false,
      "edit_link": "",
      "variables": {
        "enabled": {
          "type": "bool",
          "default": null,
          "description": "Set to false to prevent the module from creating any resources",
          "sensitive": false,
          "required": false
        },
        "environment": {
          "type": "string",
          "default": null,
          "description": "ID element. Usually used for region e.g. 'uw2', 'us-west-2', OR role 'prod', 'staging', 'dev', 'UAT'",
          "sensitive": false,
          "required": false
        },
        "namespace": {
          "type": "string",
          "default": null,
          "description": "ID element. Usually an abbreviation of your organization name, e.g. 'eg' or 'cp', to help ensure generated IDs are globally unique",
          "sensitive": false,
          "required": false
        },
        "stage": {
          "type": "string",
          "default": null,
          "description": "ID element. Usually used to indicate role, e.g. 'prod', 'staging', 'source', 'build', 'test', 'deploy', 'release'",
          "sensitive": false,
          "required": false
        },
        "tenant": {
          "type": "string",
          "default": null,
          "description": "ID element _(Rarely used, not included by default)_. A customer identifier, indicating who this instance of a resource is for",
          "sensitive": false,
          "required": false
        }
      },
      "outputs": {
        "chained_descriptor_account_name": {
          "sensitive": false,
          "description": ""
        },
        "chained_descriptor_stack": {
          "sensitive": false,
          "description": ""
        },
        "compatible": {
          "sensitive": false,
          "description": ""
        },
        "descriptor_account_name": {
          "sensitive": false,
          "description": ""
        },
        "descriptor_stack": {
          "sensitive": false,
          "description": ""
        }
      },
      "schema_error": ""
    }
  },
  "submodules": {

  }
}`

var labelNullExpectedNewModuleData = &tfregistry.ModuleData{
	Version: version.Must(version.NewVersion("0.25.0")),
	Inputs: []tfregistry.Input{
		{
			Name:        "enabled",
			Type:        cty.Bool,
			Description: lang.Markdown("Set to false to prevent the module from creating any resources"),
			Required:    false,
		},
		{
			Name:        "environment",
			Type:        cty.String,
			Description: lang.Markdown("ID element. Usually used for region e.g. 'uw2', 'us-west-2', OR role 'prod', 'staging', 'dev', 'UAT'"),
			Required:    false,
		},
	},
	Outputs: []tfregistry.Output{
		{
			Name:        "additional_tag_map",
			Description: lang.Markdown("The merged additional_tag_map"),
		},
		{
			Name:        "attributes",
			Description: lang.Markdown("List of attributes"),
		},
	},
}
