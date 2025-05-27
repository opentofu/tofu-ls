// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

// moduleVersionsMockResponse represents response from https://registry.terraform.io/v1/modules/puppetlabs/deployment/ec/versions
var moduleVersionsMockResponse = `{
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

// moduleDataMockResponse represents response from https://api.opentofu.org/registry/docs/modules/puppetlabs/deployment/ec/v0.0.8/index.json
var moduleDataMockResponse = `{
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
