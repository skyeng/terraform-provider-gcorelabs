---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_router Resource - terraform-provider-gcorelabs"
subcategory: ""
description: |-
  Represent router. Router enables you to dynamically exchange routes between networks
---

# gcore_router (Resource)

Represent router. Router enables you to dynamically exchange routes between networks

## Example Usage

```terraform
provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_api = "https://api.cloud.gcorelabs.com"
}

resource "gcore_router" "router" {
  name = "router_example"

  dynamic external_gateway_info {
  iterator = egi
  for_each = var.external_gateway_info
  content {
    type = egi.value.type
    enable_snat = egi.value.enable_snat
    network_id = egi.value.network_id
    }
  }

  dynamic interfaces {
  iterator = ifaces
  for_each = var.interfaces
  content {
    type = ifaces.value.type
    subnet_id = ifaces.value.subnet_id
    }
  }

  dynamic routes {
  iterator = rs
  for_each = var.routes
  content {
    destination = rs.value.destination
    nexthop = rs.value.nexthop
    }
  }

  region_id = 1
  project_id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)

### Optional

- **external_gateway_info** (Block List, Max: 1) (see [below for nested schema](#nestedblock--external_gateway_info))
- **id** (String) The ID of this resource.
- **interfaces** (Block List) (see [below for nested schema](#nestedblock--interfaces))
- **last_updated** (String)
- **project_id** (Number)
- **project_name** (String)
- **region_id** (Number)
- **region_name** (String)
- **routes** (Block List) (see [below for nested schema](#nestedblock--routes))

<a id="nestedblock--external_gateway_info"></a>
### Nested Schema for `external_gateway_info`

Required:

- **network_id** (String) Id of the external network

Optional:

- **enable_snat** (Boolean)
- **type** (String) Must be 'manual'

Read-Only:

- **external_fixed_ips** (List of Object) (see [below for nested schema](#nestedatt--external_gateway_info--external_fixed_ips))

<a id="nestedatt--external_gateway_info--external_fixed_ips"></a>
### Nested Schema for `external_gateway_info.external_fixed_ips`

Read-Only:

- **ip_address** (String)
- **subnet_id** (String)



<a id="nestedblock--interfaces"></a>
### Nested Schema for `interfaces`

Required:

- **subnet_id** (String) Subnet for router interface must have a gateway IP
- **type** (String) must be 'subnet'


<a id="nestedblock--routes"></a>
### Nested Schema for `routes`

Required:

- **destination** (String)
- **nexthop** (String) IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR

## Import

Import is supported using the following syntax:

```shell
# import using <project_id>:<region_id>:<router_id> format
terraform import gcore_router.router1 1:6:447d2959-8ae0-4ca0-8d47-9f050a3637d7
```
