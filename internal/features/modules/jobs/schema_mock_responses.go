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

// puppetModuleVersionsMockResponse represents response from https://api.opentofu.org/registry/docs/modules/puppetlabs/deployment/ec/index.json
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

// puppetModuleDataMockResponse represents response from https://registry.terraform.io/v1/modules/puppetlabs/deployment/ec/0.0.8
var puppetModuleDataMockResponse = `{
  "id": "puppetlabs/deployment/ec/0.0.8",
  "owner": "mattkirby",
  "namespace": "puppetlabs",
  "name": "deployment",
  "version": "0.0.8",
  "provider": "ec",
  "provider_logo_url": "/images/providers/generic.svg?2",
  "description": "",
  "source": "https://github.com/puppetlabs/terraform-ec-deployment",
  "tag": "v0.0.8",
  "published_at": "2021-08-05T00:26:33.501756Z",
  "downloads": 3059237,
  "verified": false,
  "root": {
    "path": "",
    "name": "deployment",
    "readme": "# EC project Terraform module\n\nTerraform module which creates a Elastic Cloud project.\n\n## Usage\n\nDetails coming soon\n",
    "empty": false,
    "inputs": [
      {
        "name": "autoscale",
        "type": "string",
        "description": "Enable autoscaling of elasticsearch",
        "default": "\"true\"",
        "required": false
      },
      {
        "name": "ec_stack_version",
        "type": "string",
        "description": "Version of Elastic Cloud stack to deploy",
        "default": "\"\"",
        "required": false
      },
      {
        "name": "name",
        "type": "string",
        "description": "Name of resources",
        "default": "\"ecproject\"",
        "required": false
      },
      {
        "name": "traffic_filter_sourceip",
        "type": "string",
        "description": "traffic filter source IP",
        "default": "\"\"",
        "required": false
      },
      {
        "name": "ec_region",
        "type": "string",
        "description": "cloud provider region",
        "default": "\"gcp-us-west1\"",
        "required": false
      },
      {
        "name": "deployment_templateid",
        "type": "string",
        "description": "ID of Elastic Cloud deployment type",
        "default": "\"gcp-io-optimized\"",
        "required": false
      }
    ],
    "outputs": [
      {
        "name": "elasticsearch_password",
        "description": "elasticsearch password"
      },
      {
        "name": "deployment_id",
        "description": "Elastic Cloud deployment ID"
      },
      {
        "name": "elasticsearch_version",
        "description": "Stack version deployed"
      },
      {
        "name": "elasticsearch_cloud_id",
        "description": "Elastic Cloud project deployment ID"
      },
      {
        "name": "elasticsearch_https_endpoint",
        "description": "elasticsearch https endpoint"
      },
      {
        "name": "elasticsearch_username",
        "description": "elasticsearch username"
      }
    ],
    "dependencies": [],
    "provider_dependencies": [
      {
        "name": "ec",
        "namespace": "elastic",
        "source": "elastic/ec",
        "version": "0.2.1"
      }
    ],
    "resources": [
      {
        "name": "ecproject",
        "type": "ec_deployment"
      },
      {
        "name": "gcp_vpc_nat",
        "type": "ec_deployment_traffic_filter"
      },
      {
        "name": "ec_tf_association",
        "type": "ec_deployment_traffic_filter_association"
      }
    ]
  },
  "submodules": [
    {
      "path": "modules/ec",
      "inputs": [
        {
          "name": "sub_autoscale",
          "type": "string",
          "description": "Enable autoscaling of elasticsearch",
          "default": "\"true\"",
          "required": false
        },
        {
          "name": "sub_ec_stack_version",
          "type": "string",
          "description": "Version of Elastic Cloud stack to deploy",
          "default": "\"\"",
          "required": false
        }
      ],
      "outputs": [
        {
          "name": "sub_elasticsearch_password",
          "description": "elasticsearch password"
        },
        {
          "name": "sub_deployment_id",
          "description": "Elastic Cloud deployment ID"
        }
      ]
    }
  ],
  "examples": [],
  "providers": [
    "ec"
  ],
  "versions": [
    "0.0.1",
    "0.0.2",
    "0.0.3",
    "0.0.4",
    "0.0.5",
    "0.0.6",
    "0.0.7",
    "0.0.8"
  ]
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
      "id": "v0.25.0-rc.1",
      "published": "2021-08-20T18:46:03Z"
    },
    {
      "id": "v0.24.1",
      "published": "2021-02-04T09:17:22Z"
    },
    {
      "id": "v0.24.0",
      "published": "2021-02-04T08:11:56Z"
    },
    {
      "id": "v0.23.0",
      "published": "2021-02-01T19:44:14Z"
    },
    {
      "id": "v0.22.1",
      "published": "2020-12-22T01:51:49Z"
    },
    {
      "id": "v0.22.0",
      "published": "2020-12-03T20:05:03Z"
    },
    {
      "id": "v0.21.0",
      "published": "2020-11-11T20:27:27Z"
    },
    {
      "id": "v0.20.0",
      "published": "2020-11-10T22:26:22Z"
    },
    {
      "id": "v0.19.2",
      "published": "2020-08-28T03:21:39Z"
    },
    {
      "id": "v0.19.1",
      "published": "2020-08-27T21:36:02Z"
    },
    {
      "id": "v0.19.0",
      "published": "2020-08-27T05:47:47Z"
    },
    {
      "id": "v0.18.0",
      "published": "2020-08-26T22:56:53Z"
    },
    {
      "id": "v0.17.0",
      "published": "2020-08-09T20:05:49Z"
    },
    {
      "id": "v0.16.0",
      "published": "2019-10-28T16:13:58Z"
    },
    {
      "id": "v0.15.0",
      "published": "2019-08-14T17:43:57Z"
    },
    {
      "id": "v0.14.1",
      "published": "2019-06-25T00:07:07Z"
    },
    {
      "id": "v0.14.0",
      "published": "2019-06-24T21:38:36Z"
    },
    {
      "id": "v0.13.0",
      "published": "2019-06-19T19:35:18Z"
    },
    {
      "id": "v0.12.2",
      "published": "2019-06-19T05:21:01Z"
    },
    {
      "id": "v0.12.1",
      "published": "2019-06-18T02:43:09Z"
    },
    {
      "id": "v0.12.0",
      "published": "2019-06-06T01:27:08Z"
    },
    {
      "id": "v0.11.1",
      "published": "2019-05-28T06:31:58Z"
    },
    {
      "id": "v0.11.0",
      "published": "2019-05-28T04:27:44Z"
    },
    {
      "id": "v0.10.0",
      "published": "2019-05-28T01:13:16Z"
    },
    {
      "id": "v0.9.0",
      "published": "2019-05-28T01:04:58Z"
    },
    {
      "id": "v0.8.0",
      "published": "2019-05-27T22:15:55Z"
    },
    {
      "id": "v0.7.0",
      "published": "2019-04-01T19:56:19Z"
    },
    {
      "id": "v0.6.3",
      "published": "2019-03-14T19:09:08Z"
    },
    {
      "id": "v0.6.2",
      "published": "2019-02-19T02:26:04Z"
    },
    {
      "id": "v0.6.1",
      "published": "2019-02-14T18:37:52Z"
    },
    {
      "id": "v0.6.0",
      "published": "2019-02-13T17:53:14Z"
    },
    {
      "id": "v0.5.4",
      "published": "2018-12-18T08:22:04Z"
    },
    {
      "id": "v0.5.3",
      "published": "2018-09-11T18:38:54Z"
    },
    {
      "id": "v0.5.2",
      "published": "2018-09-05T19:46:15Z"
    },
    {
      "id": "v0.5.1",
      "published": "2018-08-31T13:35:32Z"
    },
    {
      "id": "v0.5.0",
      "published": "2018-08-24T18:32:14Z"
    },
    {
      "id": "v0.4.1",
      "published": "2018-07-25T10:05:35Z"
    },
    {
      "id": "v0.4.0",
      "published": "2018-07-24T21:02:54Z"
    },
    {
      "id": "v0.3.8",
      "published": "2018-07-24T19:11:52Z"
    },
    {
      "id": "v0.3.7",
      "published": "2018-07-05T12:37:12Z"
    },
    {
      "id": "v0.3.6",
      "published": "2018-06-28T22:05:09Z"
    },
    {
      "id": "v0.3.5",
      "published": "2018-05-16T04:39:24Z"
    },
    {
      "id": "v0.3.4",
      "published": "2018-05-16T01:41:52Z"
    },
    {
      "id": "v0.3.3",
      "published": "2018-02-27T19:03:56Z"
    },
    {
      "id": "v0.3.2",
      "published": "2018-02-27T01:45:31Z"
    },
    {
      "id": "v0.3.1",
      "published": "2017-11-15T15:40:14Z"
    },
    {
      "id": "v0.3.0",
      "published": "2017-10-30T17:46:05Z"
    },
    {
      "id": "v0.2.2",
      "published": "2017-10-13T23:30:33Z"
    },
    {
      "id": "v0.2.1",
      "published": "2017-09-20T18:08:15Z"
    },
    {
      "id": "v0.2.0",
      "published": "2017-08-24T17:42:16Z"
    },
    {
      "id": "v0.1.0",
      "published": "2017-08-03T08:53:22Z"
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

// labelNullModuleDataOldMockResponse represents response for
// a module that suffers from "unreliable" input data, as described in
// https://github.com/hashicorp/vscode-terraform/issues/1582
// It is a shortened response from https://api.opentofu.org/registry/docs/modules/cloudposse/label/null/v0.25.0/index.json
var labelNullModuleDataOldMockResponse = `{
  "id": "v0.25.0",
  "published": "2021-08-25T17:45:16Z",
  "readme": true,
  "edit_link": "https://github.com/cloudposse/terraform-null-label/blob/0.25.0/README.md",
  "variables": {
    "additional_tag_map": {
      "type": "map of string",
      "default": {

      },
      "description": "Additional key-value\n",
      "sensitive": false,
      "required": false
    },
    "attributes": {
      "type": "list of string",
      "default": [],
      "description": "ID element. Additional attributes (\n",
      "sensitive": false,
      "required": false
    },
    "context": {
      "type": "dynamic",
      "default": {
        "additional_tag_map": {

        },
        "attributes": [],
        "delimiter": null,
        "descriptor_formats": {

        },
        "enabled": true,
        "environment": null,
        "id_length_limit": null,
        "label_key_case": null,
        "label_order": [],
        "label_value_case": null,
        "labels_as_tags": [
          "unset"
        ],
        "name": null,
        "namespace": null,
        "regex_replace_chars": null,
        "stage": null,
        "tags": {

        },
        "tenant": null
      },
      "description": "Single object for setting entire context at once.\nSee description of individual variables for details.\nLeave string and numeric variables as \n",
      "sensitive": false,
      "required": false
    },
    "delimiter": {
      "type": "string",
      "default": null,
      "description": "Delimiter to be used between ID elements.\n",
      "sensitive": false,
      "required": false
    },
    "descriptor_formats": {
      "type": "dynamic",
      "default": {

      },
      "description": "Describe additional descriptors to be output in then\n",
      "sensitive": false,
      "required": false
    },
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
    "id_length_limit": {
      "type": "number",
      "default": null,
      "description": "Limit.\n",
      "sensitive": false,
      "required": false
    },
    "label_key_case": {
      "type": "string",
      "default": null,
      "description": "Controls the letter case of the.\n",
      "sensitive": false,
      "required": false
    },
    "label_order": {
      "type": "list of string",
      "default": null,
      "description": "The order in which the labels (ID elements) appear in t\n",
      "sensitive": false,
      "required": false
    },
    "label_value_case": {
      "type": "string",
      "default": null,
      "description": "Controls the letter case .\n",
      "sensitive": false,
      "required": false
    },
    "labels_as_tags": {
      "type": "set of string",
      "default": [
        "default"
      ],
      "description": "Set o.\n",
      "sensitive": false,
      "required": false
    },
    "name": {
      "type": "string",
      "default": null,
      "description": "ID element. Usually the component or solution name, e.g. 'app' or 'jenkins'.\nThis is the only ID element not also included as.\n",
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
    "regex_replace_chars": {
      "type": "string",
      "default": null,
      "description": "Terraform regular expression (regex) string.\nCharacters matching the regex will be removed from the ID elements.\nIf not set.\n",
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
    "tags": {
      "type": "map of string",
      "default": {

      },
      "description": "Additional tags (e..\n",
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
    "additional_tag_map": {
      "sensitive": false,
      "description": "The merged additional_tag_map"
    },
    "attributes": {
      "sensitive": false,
      "description": "List of attributes"
    },
    "context": {
      "sensitive": false,
      "description": "Merged but otherwise unmodified input to this module, to be used as context input to other modules.\nNote: this version will have null values as defaults, not the values actually used as defaults.\n"
    },
    "delimiter": {
      "sensitive": false,
      "description": "Delimiter"
    },
    "descriptors": {
      "sensitive": false,
      "description": "Map of descriptors as configured"
    },
    "enabled": {
      "sensitive": false,
      "description": "True if module is enabled, false otherwise"
    },
    "environment": {
      "sensitive": false,
      "description": "Normalized environment"
    },
    "id": {
      "sensitive": false,
      "description": "Disambiguated ID string restricted tocharacters in total"
    },
    "id_full": {
      "sensitive": false,
      "description": "ID string not restricted in length"
    },
    "id_length_limit": {
      "sensitive": false,
      "description": "The id_length_limit actually used to create the ID, with 0 meaning unlimited"
    },
    "label_order": {
      "sensitive": false,
      "description": "The naming order actually used to create the ID"
    },
    "name": {
      "sensitive": false,
      "description": "Normalized name"
    },
    "namespace": {
      "sensitive": false,
      "description": "Normalized namespace"
    },
    "normalized_context": {
      "sensitive": false,
      "description": "Normalized context of this module"
    },
    "regex_replace_chars": {
      "sensitive": false,
      "description": "The regex_replace_chars actually used to create the ID"
    },
    "stage": {
      "sensitive": false,
      "description": "Normalized stage"
    },
    "tags": {
      "sensitive": false,
      "description": "Normalized Tag map"
    },
    "tags_as_list_of_maps": {
      "sensitive": false,
      "description": "This is a list with one map for each\n"
    },
    "tenant": {
      "sensitive": false,
      "description": "Normalized tenant"
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
        "additional_tag_map": {
          "type": "map of string",
          "default": {

          },
          "description": "Additional tags for appending to tags_as_list_of_maps. Not added to",
          "sensitive": false,
          "required": false
        },
        "attributes": {
          "type": "list of string",
          "default": [],
          "description": "Additional attributes (e.g. 1)",
          "sensitive": false,
          "required": false
        },
        "context": {
          "type": "object",
          "default": {
            "additional_tag_map": {

            },
            "attributes": [],
            "delimiter": null,
            "enabled": true,
            "environment": null,
            "id_length_limit": null,
            "label_key_case": null,
            "label_order": [],
            "label_value_case": null,
            "name": null,
            "namespace": null,
            "regex_replace_chars": null,
            "stage": null,
            "tags": {

            }
          },
          "description": "Single object for setting entire context at once.\nSee description of individual variables for details.\nLeave string and numeric variables as null to use default value.\nIndividual variable settings (non-null) override settings in context object,\nexcept for attributes, tags, and additional_tag_map, which are merged.\n",
          "sensitive": false,
          "required": false
        },
        "delimiter": {
          "type": "string",
          "default": null,
          "description": "Delimiter to be used betwe\n",
          "sensitive": false,
          "required": false
        },
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
        "id_length_limit": {
          "type": "number",
          "default": null,
          "description": "Limit .\n",
          "sensitive": false,
          "required": false
        },
        "label_key_case": {
          "type": "string",
          "default": null,
          "description": "The letter case of label key\n",
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
        "label_value_case": {
          "type": "string",
          "default": null,
          "description": "The letter case of .\n",
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
        },
        "regex_replace_chars": {
          "type": "string",
          "default": null,
          "description": "Regex to replace chars with empty string in .\n",
          "sensitive": false,
          "required": false
        },
        "stage": {
          "type": "string",
          "default": null,
          "description": "Stage, e.g. 'prod', 'staging', 'dev', OR 'source', 'build', 'test', 'deploy', 'release'",
          "sensitive": false,
          "required": false
        },
        "tags": {
          "type": "map of string",
          "default": {

          },
          "description": "Additional tags (e.g",
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
        "additional_tag_map": {
          "type": "map of string",
          "default": {

          },
          "description": "This is for some rare cases where resources want additional configuration of tags\nand therefore take a list of maps with tag key, value, and additional configuration.\n",
          "sensitive": false,
          "required": false
        },
        "attributes": {
          "type": "list of string",
          "default": [],
          "description": "ID element. Additional attributes (e.g. s a single ID element.\n",
          "sensitive": false,
          "required": false
        },
        "context": {
          "type": "dynamic",
          "default": {
            "additional_tag_map": {

            },
            "attributes": [],
            "delimiter": null,
            "descriptor_formats": {

            },
            "enabled": true,
            "environment": null,
            "id_length_limit": null,
            "label_key_case": null,
            "label_order": [],
            "label_value_case": null,
            "labels_as_tags": [
              "unset"
            ],
            "name": null,
            "namespace": null,
            "regex_replace_chars": null,
            "stage": null,
            "tags": {

            },
            "tenant": null
          },
          "description": "Single object for setting entire context at once.\nSee descriptisettings in context object,\nexcept for attributes, tags, and additional_tag_map, which are merged.\n",
          "sensitive": false,
          "required": false
        },
        "delimiter": {
          "type": "string",
          "default": null,
          "description": "Delimiter to be used between ID elements.\nDefaults(hyphen). Set to\n",
          "sensitive": false,
          "required": false
        },
        "descriptor_formats": {
          "type": "dynamic",
          "default": {

          },
          "description": "Describe additional descriptors to be output in the utput map.\nMap of maps. Keys are names of descriptors. Values are maps of the for\n",
          "sensitive": false,
          "required": false
        },
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
        "id_length_limit": {
          "type": "number",
          "default": null,
          "description": "Limit .\n",
          "sensitive": false,
          "required": false
        },
        "label_key_case": {
          "type": "string",
          "default": null,
          "description": "Controls the letter case of thault value:.\n",
          "sensitive": false,
          "required": false
        },
        "label_order": {
          "type": "list of string",
          "default": null,
          "description": "The order in which the labels (ID elements) appear in the.\nDefaults to [\"namespace\", \"environment\", \"stage\", \"name\", \"attributes\"].\nYou can omit any of the 6 labels (\"tenant\" is the 6th), but at least one must be present.\n",
          "sensitive": false,
          "required": false
        },
        "label_value_case": {
          "type": "string",
          "default": null,
          "description": "Controls t.\n",
          "sensitive": false,
          "required": false
        },
        "labels_as_tags": {
          "type": "set of string",
          "default": [
            "default"
          ],
          "description": "Set of labels (ID elements) to include as tags in .\n",
          "sensitive": false,
          "required": false
        },
        "name": {
          "type": "string",
          "default": null,
          "description": "ID element. Usually the component or solution na.\n",
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
        "regex_replace_chars": {
          "type": "string",
          "default": null,
          "description": "Terraform regular expression (regex) string.\nCharacters matching the regex will be ",
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
        "tags": {
          "type": "map of string",
          "default": {

          },
          "description": "Additional tags (e.g..\n",
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
        "compare_22_25_empty": {
          "sensitive": false,
          "description": ""
        },
        "compare_22_25_full": {
          "sensitive": false,
          "description": ""
        },
        "compare_24_25_empty": {
          "sensitive": false,
          "description": ""
        },
        "compare_24_25_full": {
          "sensitive": false,
          "description": ""
        },
        "compare_25_22_empty": {
          "sensitive": false,
          "description": ""
        },
        "compare_25_22_full": {
          "sensitive": false,
          "description": ""
        },
        "compare_25_24_empty": {
          "sensitive": false,
          "description": ""
        },
        "compare_25_24_full": {
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
        },
        "label1": {
          "sensitive": false,
          "description": ""
        },
        "label1_context": {
          "sensitive": false,
          "description": ""
        },
        "label1_normalized_context": {
          "sensitive": false,
          "description": ""
        },
        "label1_tags": {
          "sensitive": false,
          "description": ""
        },
        "label1t1": {
          "sensitive": false,
          "description": ""
        },
        "label1t1_tags": {
          "sensitive": false,
          "description": ""
        },
        "label1t2": {
          "sensitive": false,
          "description": ""
        },
        "label1t2_tags": {
          "sensitive": false,
          "description": ""
        },
        "label2": {
          "sensitive": false,
          "description": ""
        },
        "label2_context": {
          "sensitive": false,
          "description": ""
        },
        "label2_tags": {
          "sensitive": false,
          "description": ""
        },
        "label2_tags_as_list_of_maps": {
          "sensitive": false,
          "description": ""
        },
        "label3c": {
          "sensitive": false,
          "description": ""
        },
        "label3c_context": {
          "sensitive": false,
          "description": ""
        },
        "label3c_normalized_context": {
          "sensitive": false,
          "description": ""
        },
        "label3c_tags": {
          "sensitive": false,
          "description": ""
        },
        "label3n": {
          "sensitive": false,
          "description": ""
        },
        "label3n_context": {
          "sensitive": false,
          "description": ""
        },
        "label3n_normalized_context": {
          "sensitive": false,
          "description": ""
        },
        "label3n_tags": {
          "sensitive": false,
          "description": ""
        },
        "label4": {
          "sensitive": false,
          "description": ""
        },
        "label4_context": {
          "sensitive": false,
          "description": ""
        },
        "label4_tags": {
          "sensitive": false,
          "description": ""
        },
        "label5": {
          "sensitive": false,
          "description": ""
        },
        "label5_context": {
          "sensitive": false,
          "description": ""
        },
        "label5_tags": {
          "sensitive": false,
          "description": ""
        },
        "label6f": {
          "sensitive": false,
          "description": ""
        },
        "label6f_tags": {
          "sensitive": false,
          "description": ""
        },
        "label6t": {
          "sensitive": false,
          "description": ""
        },
        "label6t_tags": {
          "sensitive": false,
          "description": ""
        },
        "label7": {
          "sensitive": false,
          "description": ""
        },
        "label7_attributes": {
          "sensitive": false,
          "description": ""
        },
        "label7_context": {
          "sensitive": false,
          "description": ""
        },
        "label7_id": {
          "sensitive": false,
          "description": ""
        },
        "label8d_chained_context_labels_as_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8d_context": {
          "sensitive": false,
          "description": ""
        },
        "label8d_context_context": {
          "sensitive": false,
          "description": ""
        },
        "label8d_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8d_context_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8d_id": {
          "sensitive": false,
          "description": ""
        },
        "label8d_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8dcd_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8dcd_id": {
          "sensitive": false,
          "description": ""
        },
        "label8dnd_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8dnd_id": {
          "sensitive": false,
          "description": ""
        },
        "label8l_context": {
          "sensitive": false,
          "description": ""
        },
        "label8l_context_context": {
          "sensitive": false,
          "description": ""
        },
        "label8l_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8l_context_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8l_id": {
          "sensitive": false,
          "description": ""
        },
        "label8l_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8n_context": {
          "sensitive": false,
          "description": ""
        },
        "label8n_context_context": {
          "sensitive": false,
          "description": ""
        },
        "label8n_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8n_context_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8n_id": {
          "sensitive": false,
          "description": ""
        },
        "label8n_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8t_context": {
          "sensitive": false,
          "description": ""
        },
        "label8t_context_context": {
          "sensitive": false,
          "description": ""
        },
        "label8t_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8t_context_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8t_id": {
          "sensitive": false,
          "description": ""
        },
        "label8t_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8u_context": {
          "sensitive": false,
          "description": ""
        },
        "label8u_context_context": {
          "sensitive": false,
          "description": ""
        },
        "label8u_context_id": {
          "sensitive": false,
          "description": ""
        },
        "label8u_context_normalized_context": {
          "sensitive": false,
          "description": ""
        },
        "label8u_context_tags": {
          "sensitive": false,
          "description": ""
        },
        "label8u_id": {
          "sensitive": false,
          "description": ""
        },
        "label8u_tags": {
          "sensitive": false,
          "description": ""
        }
      },
      "schema_error": ""
    }
  },
  "submodules": []
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
