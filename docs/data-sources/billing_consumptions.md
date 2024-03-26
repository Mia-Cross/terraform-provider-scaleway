---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "scaleway_billing_consumptions Data Source - terraform-provider-scaleway"
subcategory: ""
description: |-
  
---

# scaleway_billing_consumptions (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `project_id` (String) The project_id you want to attach the resource to

### Read-Only

- `consumptions` (List of Object) (see [below for nested schema](#nestedatt--consumptions))
- `id` (String) The ID of this resource.
- `organization_id` (String) The organization_id you want to attach the resource to
- `updated_at` (String)

<a id="nestedatt--consumptions"></a>
### Nested Schema for `consumptions`

Read-Only:

- `billed_quantity` (String)
- `category_name` (String)
- `product_name` (String)
- `project_id` (String)
- `sku` (String)
- `unit` (String)
- `value` (String)