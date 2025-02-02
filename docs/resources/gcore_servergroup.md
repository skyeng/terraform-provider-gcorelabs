---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_servergroup Resource - terraform-provider-gcorelabs"
subcategory: ""
description: |-
  Represent server group resource
---

# gcore_servergroup (Resource)

Represent server group resource

## Example Usage

```terraform
provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_servergroup" "default" {
  name = "default"
  policy = "affinity"
  region_id = 1
  project_id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Displayed server group name
- **policy** (String) Server group policy. Available value is 'affinity', 'anti-affinity'

### Optional

- **id** (String) The ID of this resource.
- **project_id** (Number)
- **project_name** (String)
- **region_id** (Number)
- **region_name** (String)

### Read-Only

- **instances** (List of Object) Instances in this server group (see [below for nested schema](#nestedatt--instances))

<a id="nestedatt--instances"></a>
### Nested Schema for `instances`

Read-Only:

- **instance_id** (String)
- **instance_name** (String)

## Import

Import is supported using the following syntax:

```shell
# import using <project_id>:<region_id>:<servergroup_id> format
terraform import gcore_servergroup.servergroup1 1:6:447d2959-8ae0-4ca0-8d47-9f050a3637d7
```
