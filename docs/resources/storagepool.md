---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerflex_storagepool Resource - powerflex"
subcategory: ""
description: |-
  Manages storage pool resource
---

# powerflex_storagepool (Resource)

Manages storage pool resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `media_type` (String) Media Type
- `name` (String) Name of the Storage pool

### Optional

- `protection_domain_id` (String) ID of the Protection domain
- `protection_domain_name` (String) Name of the Protection domain.
- `use_rfcache` (Boolean) Enable/Disable RFcache on a specific storage pool
- `use_rmcache` (Boolean) Enable/Disable RMcache on a specific storage pool

### Read-Only

- `id` (String) ID of the Storage pool
- `last_updated` (String) Last Updated
- `links` (Attributes List) Specifies the links asscociated with Storagepool (see [below for nested schema](#nestedatt--links))
- `systemid` (String) ID of the system

<a id="nestedatt--links"></a>
### Nested Schema for `links`

Read-Only:

- `href` (String) Specifies the exact path to fetch the details
- `rel` (String) Specifies the relationship with the Storagepool

