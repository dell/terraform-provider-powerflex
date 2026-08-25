<!--
Copyright (c) 2024-2026 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Migration Guide: PowerFlex 4.x to 5.x (Gen2)

This guide helps you migrate your Terraform configurations from PowerFlex 4.x to PowerFlex 5.x (Gen2 architecture).

## Overview

PowerFlex 5.x introduces a new Gen2 architecture with updated REST APIs. The Terraform provider supports both 4.x and 5.x through a hybrid approach:
- **4.x features** continue using the goscaleio SDK
- **5.x Gen2 features** use the OpenAPI-generated client (clientgen)

## New Resources in 5.x

The following new resources are available for PowerFlex 5.x Gen2 systems:

### Storage Management
- `powerflex_storage_node` - Gen2 replacement for SDS
- `powerflex_device_group` - Device group management
- `powerflex_storage_pool_erasure_coding` - Erasure coding configuration

### Actions
- `powerflex_device_action` - Device operations (clear errors, activate)
- `powerflex_volume_action` - Volume operations (refresh, restore, map/unmap host)

### Datasources
- `powerflex_storage_node` - Query storage nodes
- `powerflex_device_group` - Query device groups

## Migration Path

### 1. Device Resource Changes

**4.x (Deprecated):**
```hcl
resource "powerflex_device" "example" {
  device_path         = "/dev/sdc"
  storage_pool_id     = "pool-123"
  sds_id              = "sds-456"  # Deprecated
  # or
  sds_name            = "sds-name"  # Deprecated
}
```

**5.x (Recommended):**
```hcl
resource "powerflex_device" "example" {
  device_path         = "/dev/sdc"
  storage_pool_id     = "pool-123"
  storage_node_id     = "node-789"  # New Gen2 parameter
  # or
  device_group_id     = "dg-012"    # New Gen2 parameter
}
```

**Migration Steps:**
1. Replace `sds_id`/`sds_name` with `storage_node_id` or `device_group_id`
2. Run `terraform plan` to verify the changes
3. Apply the changes with `terraform apply`

### 2. Storage Pool Erasure Coding

**New Resource for 5.x:**
```hcl
resource "powerflex_storage_pool_erasure_coding" "example" {
  storage_pool_id       = "pool-123"
  erasure_coding_policy = "rs_2_1"  # none, rs_2_1, rs_4_1
  protection_scheme     = "TwoPlusTwo"  # Optional
}
```

### 3. Storage Node (Gen2 SDS)

**New Resource for 5.x:**
```hcl
resource "powerflex_storage_node" "example" {
  protection_domain_id = "pd-123"
  ip_list {
    ip   = "192.168.1.1"
    role = "All"  # SdsOnly, SdcOnly, All
  }
  fault_set_id         = "fs-456"  # Optional
  performance_profile  = "HighPerformance"  # Optional
}
```

### 4. Device Group

**New Resource for 5.x:**
```hcl
resource "powerflex_device_group" "example" {
  name                 = "device-group-1"
  protection_domain_id = "pd-123"
}
```

### 5. Volume Actions

**New Resource for 5.x:**
```hcl
# Refresh volume metadata
resource "powerflex_volume_action" "refresh" {
  volume_id = "vol-abc123"
  action    = "refresh"
}

# Restore from snapshot
resource "powerflex_volume_action" "restore" {
  volume_id          = "vol-abc123"
  action             = "restore"
  source_snapshot_id = "snap-def456"
}

# Map to host
resource "powerflex_volume_action" "map" {
  volume_id   = "vol-abc123"
  action      = "map_to_host"
  host_id     = "host-ghi789"
  access_mode = "READ_WRITE"  # READ_ONLY, READ_WRITE
}

# Unmap from host
resource "powerflex_volume_action" "unmap" {
  volume_id = "vol-abc123"
  action    = "unmap_from_host"
  host_id   = "host-ghi789"
}
```

## Backward Compatibility

### Deprecated Parameters

The following parameters are deprecated for PowerFlex 5.x but still supported:

| Resource | Deprecated Parameter | New Parameter |
|----------|---------------------|--------------|
| `powerflex_device` | `sds_id` | `storage_node_id` |
| `powerflex_device` | `sds_name` | `storage_node_id` |

**Deprecation Warnings:**
- Using deprecated parameters will emit a warning in Terraform output
- The warnings suggest migrating to the new Gen2 parameters
- Existing configurations continue to work without breaking changes

### Gradual Migration

You can migrate gradually:
1. Keep existing 4.x configurations working
2. Add new 5.x resources alongside existing ones
3. Migrate resources one at a time
4. Remove deprecated parameters after migration is complete

## Compatibility Matrix

| Feature | PowerFlex 4.x | PowerFlex 5.x |
|---------|---------------|---------------|
| SDS resource | ✅ | ✅ (via goscaleio) |
| Storage Node resource | ❌ | ✅ (Gen2) |
| Device resource (sds_id) | ✅ | ⚠️ Deprecated |
| Device resource (storage_node_id) | ❌ | ✅ (Gen2) |
| Storage Pool EC | ❌ | ✅ (Gen2) |
| Volume actions | ❌ | ✅ (Gen2) |

## Troubleshooting

### Error: "Gen2 API Client Not Available"

**Cause:** The PowerFlex system does not support Gen2 APIs (version < 5.0).

**Solution:** Use 4.x resources (SDS, device with sds_id) instead of Gen2 resources.

### Error: "sds_id is deprecated"

**Cause:** Using deprecated parameter in 5.x environment.

**Solution:** Replace `sds_id` with `storage_node_id` or `device_group_id`.

### Validation Error: "value must be one of"

**Cause:** Invalid value for a Gen2 parameter (e.g., invalid EC policy).

**Solution:** Check the valid values in the resource documentation.

## Additional Resources

- [Storage Node Resource Documentation](docs/resources/storage_node.md)
- [Device Group Resource Documentation](docs/resources/device_group.md)
- [Volume Action Resource Documentation](docs/resources/volume_action.md)
- [Storage Pool Erasure Coding Documentation](docs/resources/storage_pool_erasure_coding.md)
