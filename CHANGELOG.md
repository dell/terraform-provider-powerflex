# v1.1.0 (June 21, 2023)
## Release Summary
The release supports resources mentioned in the Features section for Dell PowerFlex.
## Features

### Resources
* `powerflex_device` for managing devices in PowerFlex.
* `powerflex_protection_domain` for managing protection domains in PowerFlex.
* `powerflex_package` for managing packages on the PowerFlex Gateway.

### Data Sources:
* `powerflex_device` for reading devices in PowerFlex.

### Enhancements:
* `powerflex_storage_pool` is enhanced to support additional attributes in PowerFlex.
* `powerflex_sdc` is enhanced to create/delete/import multiple SDCs in PowerFlex.

### Deprecations:
* `name` attribute from SDC resource.

### Notes:
* `sdc_list` attribute is removed from volume and snapshot resource.

# v1.0.1 (May 23, 2023)
## Release Summary
The release supports resources mentioned in the Features section for Dell PowerFlex.
## Features

### Resources
* `powerflex_sdc_volumes_mapping` for managing map/unmap operations between SDC and volumes in PowerFlex.

### Deprecations
* sdc_list attribute in powerflex_volume and powerflex_snapshot resource.
***
<br>

# v1.0.0 (Feb 27, 2023)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerFlex.
## Features

### Data Sources:
* `powerflex_protection_domain` for reading protection domain in PowerFlex.
* `powerflex_sdc` for reading SDC in PowerFlex.
* `powerflex_sds` for reading SDS in PowerFlex.
* `powerflex_snapshot_policy` for reading snapshot policy in PowerFlex.
* `powerflex_storage_pool` for reading storage pool in PowerFlex.
* `powerflex_volume` for reading volume in PowerFlex.

### Resources
* `powerflex_sdc` for managing SDC in PowerFlex.
* `powerflex_sds` for managing SDS in PowerFlex.
* `powerflex_snapshot` for managing Snapshot in PowerFlex.
* `powerflex_storage_pool` for managing Storage Pool in PowerFlex.
* `powerflex_volume` for managing Volume in PowerFlex.

### Others
N/A
## Enhancements
N/A

## Bug Fixes
N/A

