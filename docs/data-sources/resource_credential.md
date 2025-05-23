---
title: "powerflex_resource_credential data source"
linkTitle: "powerflex_resource_credential"
page_title: "powerflex_resource_credential Data Source - powerflex"
subcategory: "Resource Group Management"
description: |-
  This datasource is used to read the Resource Credential entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.
---

# powerflex_resource_credential (Data Source)

This datasource is used to read the Resource Credential entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.



After the successful execution of above said block, we can see the output by executing `terraform output` command. Also, we can fetch information via the variable: `data.powerflex_resource_credential.datasource_block_name.attribute_name` where datasource_block_name is the name of the data source block and attribute_name is the attribute which user wants to fetch.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `id` (String) default datasource id
- `resource_credential_details` (Attributes List) List of Resource Credentials (see [below for nested schema](#nestedatt--resource_credential_details))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `created_by` (Set of String) List of created_by
- `created_date` (Set of String) List of created_date
- `domain` (Set of String) List of domain
- `id` (Set of String) List of id
- `label` (Set of String) List of label
- `type` (Set of String) List of type
- `updated_by` (Set of String) List of updated_by
- `updated_date` (Set of String) List of updated_date
- `username` (Set of String) List of username


<a id="nestedatt--resource_credential_details"></a>
### Nested Schema for `resource_credential_details`

Read-Only:

- `created_by` (String) Who the Resource Credential was created by
- `created_date` (String) Resource Credential created date
- `domain` (String) Resource Credential domain
- `id` (String) Unique identifier of the resource credential instance.
- `label` (String) Resource Credential label
- `type` (String) Resource Credential type
- `updated_by` (String) Who the Resource Credential was last updated by
- `updated_date` (String) Resource Credential updated date
- `username` (String) Resource Credential username