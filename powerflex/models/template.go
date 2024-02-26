/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TemplateDataSourceModel maps the struct to Template data source schema
type TemplateDataSourceModel struct {
	TemplateIDs     types.Set       `tfsdk:"template_ids"`
	TemplateNames   types.Set       `tfsdk:"template_names"`
	TemplateDetails []TemplateModel `tfsdk:"template_details"`
	ID              types.String    `tfsdk:"id"`
}

// TemplateModel is the tfsdk model of TemplateDetails
type TemplateModel struct {
	ID                     types.String         `tfsdk:"id"`
	TemplateName           types.String         `tfsdk:"template_name"`
	TemplateDescription    types.String         `tfsdk:"template_description"`
	TemplateType           types.String         `tfsdk:"template_type"`
	TemplateVersion        types.String         `tfsdk:"template_version"`
	OriginalTemplateID     types.String         `tfsdk:"original_template_id"`
	TemplateValid          TemplateValid        `tfsdk:"template_valid"`
	TemplateLocked         types.Bool           `tfsdk:"template_locked"`
	InConfiguration        types.Bool           `tfsdk:"in_configuration"`
	CreatedDate            types.String         `tfsdk:"created_date"`
	CreatedBy              types.String         `tfsdk:"created_by"`
	UpdatedDate            types.String         `tfsdk:"updated_date"`
	LastDeployedDate       types.String         `tfsdk:"last_deployed_date"`
	UpdatedBy              types.String         `tfsdk:"updated_by"`
	ManageFirmware         types.Bool           `tfsdk:"manage_firmware"`
	UseDefaultCatalog      types.Bool           `tfsdk:"use_default_catalog"`
	FirmwareRepository     FirmwareRepository   `tfsdk:"firmware_repository"`
	LicenseRepository      LicenseRepository    `tfsdk:"license_repository"`
	AssignedUsers          []AssignedUsers      `tfsdk:"assigned_users"`
	AllUsersAllowed        types.Bool           `tfsdk:"all_users_allowed"`
	Category               types.String         `tfsdk:"category"`
	Components             []Components         `tfsdk:"components"`
	Configuration          ConfigurationDetails `tfsdk:"configuration"`
	ServerCount            types.Int64          `tfsdk:"server_count"`
	StorageCount           types.Int64          `tfsdk:"storage_count"`
	ClusterCount           types.Int64          `tfsdk:"cluster_count"`
	ServiceCount           types.Int64          `tfsdk:"service_count"`
	SwitchCount            types.Int64          `tfsdk:"switch_count"`
	VMCount                types.Int64          `tfsdk:"vm_count"`
	SdnasCount             types.Int64          `tfsdk:"sdnas_count"`
	BrownfieldTemplateType types.String         `tfsdk:"brownfield_template_type"`
	Networks               []Networks           `tfsdk:"networks"`
	Draft                  types.Bool           `tfsdk:"draft"`
}

// Messages is the tfsdk model of Messages
type Messages struct {
	ID              types.String `tfsdk:"id"`
	MessageCode     types.String `tfsdk:"message_code"`
	MessageBundle   types.String `tfsdk:"message_bundle"`
	Severity        types.String `tfsdk:"severity"`
	Category        types.String `tfsdk:"category"`
	DisplayMessage  types.String `tfsdk:"display_message"`
	ResponseAction  types.String `tfsdk:"response_action"`
	DetailedMessage types.String `tfsdk:"detailed_message"`
	CorrelationID   types.String `tfsdk:"correlation_id"`
	AgentID         types.String `tfsdk:"agent_id"`
	TimeStamp       types.String `tfsdk:"time_stamp"`
	SequenceNumber  types.Int64  `tfsdk:"sequence_number"`
}

// TemplateValid is the tfsdk model of TemplateValid
type TemplateValid struct {
	Valid    types.Bool `tfsdk:"valid"`
	Messages []Messages `tfsdk:"messages"`
}

// SoftwareComponents is the tfsdk model of SoftwareComponents
type SoftwareComponents struct {
	ID                  types.String   `tfsdk:"id"`
	PackageID           types.String   `tfsdk:"package_id"`
	DellVersion         types.String   `tfsdk:"dell_version"`
	VendorVersion       types.String   `tfsdk:"vendor_version"`
	ComponentID         types.String   `tfsdk:"component_id"`
	DeviceID            types.String   `tfsdk:"device_id"`
	SubDeviceID         types.String   `tfsdk:"sub_device_id"`
	VendorID            types.String   `tfsdk:"vendor_id"`
	SubVendorID         types.String   `tfsdk:"sub_vendor_id"`
	CreatedDate         types.String   `tfsdk:"created_date"`
	CreatedBy           types.String   `tfsdk:"created_by"`
	UpdatedDate         types.String   `tfsdk:"updated_date"`
	UpdatedBy           types.String   `tfsdk:"updated_by"`
	Path                types.String   `tfsdk:"path"`
	HashMd5             types.String   `tfsdk:"hash_md_5"`
	Name                types.String   `tfsdk:"name"`
	Category            types.String   `tfsdk:"category"`
	ComponentType       types.String   `tfsdk:"component_type"`
	OperatingSystem     types.String   `tfsdk:"operating_system"`
	SystemIDs           []types.String `tfsdk:"system_ids"`
	Custom              types.Bool     `tfsdk:"custom"`
	NeedsAttention      types.Bool     `tfsdk:"needs_attention"`
	Ignore              types.Bool     `tfsdk:"ignore"`
	OriginalVersion     types.String   `tfsdk:"original_version"`
	OriginalComponentID types.String   `tfsdk:"original_component_id"`
	FirmwareRepoName    types.String   `tfsdk:"firmware_repo_name"`
}

// SoftwareBundles is the tfsdk model of SoftwareBundles
type SoftwareBundles struct {
	ID                 types.String         `tfsdk:"id"`
	Name               types.String         `tfsdk:"name"`
	Version            types.String         `tfsdk:"version"`
	BundleDate         types.String         `tfsdk:"bundle_date"`
	CreatedDate        types.String         `tfsdk:"created_date"`
	CreatedBy          types.String         `tfsdk:"created_by"`
	UpdatedDate        types.String         `tfsdk:"updated_date"`
	UpdatedBy          types.String         `tfsdk:"updated_by"`
	Description        types.String         `tfsdk:"description"`
	UserBundle         types.Bool           `tfsdk:"user_bundle"`
	UserBundlePath     types.String         `tfsdk:"user_bundle_path"`
	UserBundleHashMd5  types.String         `tfsdk:"user_bundle_hash_md_5"`
	DeviceType         types.String         `tfsdk:"device_type"`
	DeviceModel        types.String         `tfsdk:"device_model"`
	Criticality        types.String         `tfsdk:"criticality"`
	FwRepositoryID     types.String         `tfsdk:"fw_repository_id"`
	BundleType         types.String         `tfsdk:"bundle_type"`
	Custom             types.Bool           `tfsdk:"custom"`
	NeedsAttention     types.Bool           `tfsdk:"needs_attention"`
	SoftwareComponents []SoftwareComponents `tfsdk:"software_components"`
}

// DeploymentValid is the tfsdk model of DeploymentValid
type DeploymentValid struct {
	Valid    types.Bool `tfsdk:"valid"`
	Messages []Messages `tfsdk:"messages"`
}

// DeploymentDevice is the tfsdk model of DeploymentDevice
type DeploymentDevice struct {
	RefID            types.String `tfsdk:"ref_id"`
	RefType          types.String `tfsdk:"ref_type"`
	LogDump          types.String `tfsdk:"log_dump"`
	Status           types.String `tfsdk:"status"`
	StatusEndTime    types.String `tfsdk:"status_end_time"`
	StatusStartTime  types.String `tfsdk:"status_start_time"`
	DeviceHealth     types.String `tfsdk:"device_health"`
	HealthMessage    types.String `tfsdk:"health_message"`
	CompliantState   types.String `tfsdk:"compliant_state"`
	BrownfieldStatus types.String `tfsdk:"brownfield_status"`
	DeviceType       types.String `tfsdk:"device_type"`
	DeviceGroupName  types.String `tfsdk:"device_group_name"`
	IPAddress        types.String `tfsdk:"ip_address"`
	CurrentIPAddress types.String `tfsdk:"current_ip_address"`
	ServiceTag       types.String `tfsdk:"service_tag"`
	ComponentID      types.String `tfsdk:"component_id"`
	StatusMessage    types.String `tfsdk:"status_message"`
	Model            types.String `tfsdk:"model"`
	CloudLink        types.Bool   `tfsdk:"cloud_link"`
	DasCache         types.Bool   `tfsdk:"das_cache"`
	DeviceState      types.String `tfsdk:"device_state"`
	PuppetCertName   types.String `tfsdk:"puppet_cert_name"`
	Brownfield       types.Bool   `tfsdk:"brownfield"`
}

// Vms is the tfsdk model of Vms
type Vms struct {
	CertificateName types.String `tfsdk:"certificate_name"`
	VMModel         types.String `tfsdk:"vm_model"`
	VMIpaddress     types.String `tfsdk:"vm_ipaddress"`
	VMManufacturer  types.String `tfsdk:"vm_manufacturer"`
	VMServiceTag    types.String `tfsdk:"vm_service_tag"`
}

// LicenseRepository is the tfsdk model of LicenseRepository
type LicenseRepository struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	DiskLocation types.String `tfsdk:"disk_location"`
	Filename     types.String `tfsdk:"filename"`
	State        types.String `tfsdk:"state"`
	CreatedDate  types.String `tfsdk:"created_date"`
	CreatedBy    types.String `tfsdk:"created_by"`
	UpdatedDate  types.String `tfsdk:"updated_date"`
	UpdatedBy    types.String `tfsdk:"updated_by"`
	LicenseData  types.String `tfsdk:"license_data"`
}

// AssignedUsers is the tfsdk model of AssignedUsers
type AssignedUsers struct {
	UserSeqID      types.Int64    `tfsdk:"user_seq_id"`
	UserName       types.String   `tfsdk:"user_name"`
	Password       types.String   `tfsdk:"password"`
	UpdatePassword types.Bool     `tfsdk:"update_password"`
	DomainName     types.String   `tfsdk:"domain_name"`
	GroupDN        types.String   `tfsdk:"group_dn"`
	GroupName      types.String   `tfsdk:"group_name"`
	FirstName      types.String   `tfsdk:"first_name"`
	LastName       types.String   `tfsdk:"last_name"`
	Email          types.String   `tfsdk:"email"`
	PhoneNumber    types.String   `tfsdk:"phone_number"`
	Enabled        types.Bool     `tfsdk:"enabled"`
	SystemUser     types.Bool     `tfsdk:"system_user"`
	CreatedDate    types.String   `tfsdk:"created_date"`
	CreatedBy      types.String   `tfsdk:"created_by"`
	UpdatedDate    types.String   `tfsdk:"updated_date"`
	UpdatedBy      types.String   `tfsdk:"updated_by"`
	Role           types.String   `tfsdk:"role"`
	UserPreference types.String   `tfsdk:"user_preference"`
	ID             types.String   `tfsdk:"id"`
	Roles          []types.String `tfsdk:"roles"`
}

// JobDetails is the tfsdk model of JobDetails
type JobDetails struct {
	Level       types.String `tfsdk:"level"`
	Message     types.String `tfsdk:"message"`
	Timestamp   types.String `tfsdk:"timestamp"`
	ExecutionID types.String `tfsdk:"execution_id"`
	ComponentID types.String `tfsdk:"component_id"`
}

// DeploymentValidationResponse is the tfsdk model of DeploymentValidationResponse
type DeploymentValidationResponse struct {
	Nodes                  types.Int64    `tfsdk:"nodes"`
	StoragePools           types.Int64    `tfsdk:"storage_pools"`
	DrivesPerStoragePool   types.Int64    `tfsdk:"drives_per_storage_pool"`
	MaxScalability         types.Int64    `tfsdk:"max_scalability"`
	VirtualMachines        types.Int64    `tfsdk:"virtual_machines"`
	NumberOfServiceVolumes types.Int64    `tfsdk:"number_of_service_volumes"`
	CanDeploy              types.Bool     `tfsdk:"can_deploy"`
	WarningMessages        []types.String `tfsdk:"warning_messages"`
	StoragePoolDiskType    []types.String `tfsdk:"storage_pool_disk_type"`
	Hostnames              []types.String `tfsdk:"hostnames"`
	NewNodeDiskTypes       []types.String `tfsdk:"new_node_disk_types"`
	NoOfFaultSets          types.Int64    `tfsdk:"no_of_fault_sets"`
	NodesPerFaultSet       types.Int64    `tfsdk:"nodes_per_fault_set"`
	ProtectionDomain       types.String   `tfsdk:"protection_domain"`
	DiskTypeMismatch       types.Bool     `tfsdk:"disk_type_mismatch"`
}

// Deployments is the tfsdk model of Deployments
type Deployments struct {
	ID                           types.String                 `tfsdk:"id"`
	DeploymentName               types.String                 `tfsdk:"deployment_name"`
	DeploymentDescription        types.String                 `tfsdk:"deployment_description"`
	DeploymentValid              DeploymentValid              `tfsdk:"deployment_valid"`
	Retry                        types.Bool                   `tfsdk:"retry"`
	Teardown                     types.Bool                   `tfsdk:"teardown"`
	TeardownAfterCancel          types.Bool                   `tfsdk:"teardown_after_cancel"`
	RemoveService                types.Bool                   `tfsdk:"remove_service"`
	CreatedDate                  types.String                 `tfsdk:"created_date"`
	CreatedBy                    types.String                 `tfsdk:"created_by"`
	UpdatedDate                  types.String                 `tfsdk:"updated_date"`
	UpdatedBy                    types.String                 `tfsdk:"updated_by"`
	DeploymentScheduledDate      types.String                 `tfsdk:"deployment_scheduled_date"`
	DeploymentStartedDate        types.String                 `tfsdk:"deployment_started_date"`
	DeploymentFinishedDate       types.String                 `tfsdk:"deployment_finished_date"`
	ScheduleDate                 types.String                 `tfsdk:"schedule_date"`
	Status                       types.String                 `tfsdk:"status"`
	Compliant                    types.Bool                   `tfsdk:"compliant"`
	DeploymentDevice             []DeploymentDevice           `tfsdk:"deployment_device"`
	Vms                          []Vms                        `tfsdk:"vms"`
	UpdateServerFirmware         types.Bool                   `tfsdk:"update_server_firmware"`
	UseDefaultCatalog            types.Bool                   `tfsdk:"use_default_catalog"`
	FirmwareRepositoryID         types.String                 `tfsdk:"firmware_repository_id"`
	LicenseRepository            LicenseRepository            `tfsdk:"license_repository"`
	LicenseRepositoryID          types.String                 `tfsdk:"license_repository_id"`
	IndividualTeardown           types.Bool                   `tfsdk:"individual_teardown"`
	DeploymentHealthStatusType   types.String                 `tfsdk:"deployment_health_status_type"`
	AssignedUsers                []AssignedUsers              `tfsdk:"assigned_users"`
	AllUsersAllowed              types.Bool                   `tfsdk:"all_users_allowed"`
	Owner                        types.String                 `tfsdk:"owner"`
	NoOp                         types.Bool                   `tfsdk:"no_op"`
	FirmwareInit                 types.Bool                   `tfsdk:"firmware_init"`
	DisruptiveFirmware           types.Bool                   `tfsdk:"disruptive_firmware"`
	PreconfigureSVM              types.Bool                   `tfsdk:"preconfigure_svm"`
	PreconfigureSVMAndUpdate     types.Bool                   `tfsdk:"preconfigure_svm_and_update"`
	ServicesDeployed             types.String                 `tfsdk:"services_deployed"`
	PrecalculatedDeviceHealth    types.String                 `tfsdk:"precalculated_device_health"`
	LifecycleModeReasons         []types.String               `tfsdk:"lifecycle_mode_reasons"`
	JobDetails                   []JobDetails                 `tfsdk:"job_details"`
	NumberOfDeployments          types.Int64                  `tfsdk:"number_of_deployments"`
	OperationType                types.String                 `tfsdk:"operation_type"`
	OperationStatus              types.String                 `tfsdk:"operation_status"`
	OperationData                types.String                 `tfsdk:"operation_data"`
	DeploymentValidationResponse DeploymentValidationResponse `tfsdk:"deployment_validation_response"`
	CurrentStepCount             types.String                 `tfsdk:"current_step_count"`
	TotalNumOfSteps              types.String                 `tfsdk:"total_num_of_steps"`
	CurrentStepMessage           types.String                 `tfsdk:"current_step_message"`
	CustomImage                  types.String                 `tfsdk:"custom_image"`
	OriginalDeploymentID         types.String                 `tfsdk:"original_deployment_id"`
	CurrentBatchCount            types.String                 `tfsdk:"current_batch_count"`
	TotalBatchCount              types.String                 `tfsdk:"total_batch_count"`
	Brownfield                   types.Bool                   `tfsdk:"brownfield"`
	ScaleUp                      types.Bool                   `tfsdk:"scale_up"`
	LifecycleMode                types.Bool                   `tfsdk:"lifecycle_mode"`
	OverallDeviceHealth          types.String                 `tfsdk:"overall_device_health"`
	Vds                          types.Bool                   `tfsdk:"vds"`
	TemplateValid                types.Bool                   `tfsdk:"template_valid"`
	ConfigurationChange          types.Bool                   `tfsdk:"configuration_change"`
	CanMigratevCLSVMs            types.Bool                   `tfsdk:"can_migratev_clsv_ms"`
}

// FirmwareRepository is the tfsdk model of FirmwareRepository
type FirmwareRepository struct {
	ID                      types.String         `tfsdk:"id"`
	Name                    types.String         `tfsdk:"name"`
	SourceLocation          types.String         `tfsdk:"source_location"`
	SourceType              types.String         `tfsdk:"source_type"`
	DiskLocation            types.String         `tfsdk:"disk_location"`
	Filename                types.String         `tfsdk:"filename"`
	Md5Hash                 types.String         `tfsdk:"md_5_hash"`
	Username                types.String         `tfsdk:"username"`
	Password                types.String         `tfsdk:"password"`
	DownloadStatus          types.String         `tfsdk:"download_status"`
	CreatedDate             types.String         `tfsdk:"created_date"`
	CreatedBy               types.String         `tfsdk:"created_by"`
	UpdatedDate             types.String         `tfsdk:"updated_date"`
	UpdatedBy               types.String         `tfsdk:"updated_by"`
	DefaultCatalog          types.Bool           `tfsdk:"default_catalog"`
	Embedded                types.Bool           `tfsdk:"embedded"`
	State                   types.String         `tfsdk:"state"`
	SoftwareComponents      []SoftwareComponents `tfsdk:"software_components"`
	SoftwareBundles         []SoftwareBundles    `tfsdk:"software_bundles"`
	Deployments             []Deployments        `tfsdk:"deployments"`
	BundleCount             types.Int64          `tfsdk:"bundle_count"`
	ComponentCount          types.Int64          `tfsdk:"component_count"`
	UserBundleCount         types.Int64          `tfsdk:"user_bundle_count"`
	Minimal                 types.Bool           `tfsdk:"minimal"`
	DownloadProgress        types.Int64          `tfsdk:"download_progress"`
	ExtractProgress         types.Int64          `tfsdk:"extract_progress"`
	FileSizeInGigabytes     types.Int64          `tfsdk:"file_size_in_gigabytes"`
	SignedKeySourceLocation types.String         `tfsdk:"signed_key_source_location"`
	Signature               types.String         `tfsdk:"signature"`
	Custom                  types.Bool           `tfsdk:"custom"`
	NeedsAttention          types.Bool           `tfsdk:"needs_attention"`
	JobID                   types.String         `tfsdk:"job_id"`
	Rcmapproved             types.Bool           `tfsdk:"rcmapproved"`
}

// ComponentValid is the tfsdk model of ComponentValid
type ComponentValid struct {
	Valid    types.Bool `tfsdk:"valid"`
	Messages []Messages `tfsdk:"messages"`
}


// DependenciesDetails is the tfsdk model of DependenciesDetails
type DependenciesDetails struct {
	ID               types.String `tfsdk:"id"`
	DependencyTarget types.String `tfsdk:"dependency_target"`
	DependencyValue  types.String `tfsdk:"dependency_value"`
}

// NetworkIPAddressList is the tfsdk model of NetworkIPAddressList
type NetworkIPAddressList struct {
	ID        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip_address"`
}

// Partitions is the tfsdk model of Partitions
type Partitions struct {
	ID                   types.String           `tfsdk:"id"`
	Name                 types.String           `tfsdk:"name"`
	Networks             []types.String         `tfsdk:"networks"`
	NetworkIPAddressList []NetworkIPAddressList `tfsdk:"network_ip_address_list"`
	Minimum              types.Int64            `tfsdk:"minimum"`
	Maximum              types.Int64            `tfsdk:"maximum"`
	LanMacAddress        types.String           `tfsdk:"lan_mac_address"`
	IscsiMacAddress      types.String           `tfsdk:"iscsi_mac_address"`
	IscsiIQN             types.String           `tfsdk:"iscsi_iqn"`
	Wwnn                 types.String           `tfsdk:"wwnn"`
	Wwpn                 types.String           `tfsdk:"wwpn"`
	Fqdd                 types.String           `tfsdk:"fqdd"`
	MirroredPort         types.String           `tfsdk:"mirrored_port"`
	MacAddress           types.String           `tfsdk:"mac_address"`
	PortNo               types.Int64            `tfsdk:"port_no"`
	PartitionNo          types.Int64            `tfsdk:"partition_no"`
	PartitionIndex       types.Int64            `tfsdk:"partition_index"`
}

// Interfaces is the tfsdk model of Interfaces
type Interfaces struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	Partitioned   types.Bool     `tfsdk:"partitioned"`
	Partitions    []Partitions   `tfsdk:"partitions"`
	Enabled       types.Bool     `tfsdk:"enabled"`
	Redundancy    types.Bool     `tfsdk:"redundancy"`
	Nictype       types.String   `tfsdk:"nictype"`
	Fqdd          types.String   `tfsdk:"fqdd"`
	MaxPartitions types.Int64    `tfsdk:"max_partitions"`
	AllNetworks   []types.String `tfsdk:"all_networks"`
}

// InterfacesDetails is the tfsdk model of InterfacesDetails
type InterfacesDetails struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Redundancy    types.Bool   `tfsdk:"redundancy"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Partitioned   types.Bool   `tfsdk:"partitioned"`
	Interfaces    []Interfaces `tfsdk:"interfaces"`
	Nictype       types.String `tfsdk:"nictype"`
	Fabrictype    types.String `tfsdk:"fabrictype"`
	MaxPartitions types.Int64  `tfsdk:"max_partitions"`
	Nports        types.Int64  `tfsdk:"nports"`
	CardIndex     types.Int64  `tfsdk:"card_index"`
	NictypeSource types.String `tfsdk:"nictype_source"`
}

// NetworkConfiguration is the tfsdk model of NetworkConfiguration
type NetworkConfiguration struct {
	ID           types.String        `tfsdk:"id"`
	Interfaces   []InterfacesDetails `tfsdk:"interfaces"`
	SoftwareOnly types.Bool          `tfsdk:"software_only"`
}

// ConfigurationDetails is the tfsdk model of ConfigurationDetails
type ConfigurationDetails struct {
	ID              types.String `tfsdk:"id"`
	Disktype        types.String `tfsdk:"disktype"`
	Comparator      types.String `tfsdk:"comparator"`
	Numberofdisks   types.Int64  `tfsdk:"numberofdisks"`
	Raidlevel       types.String `tfsdk:"raidlevel"`
	VirtualDiskFqdd types.String `tfsdk:"virtual_disk_fqdd"`
	ControllerFqdd  types.String `tfsdk:"controller_fqdd"`
	Categories      []Categories `tfsdk:"categories"`
}

// VirtualDisks is the tfsdk model of VirtualDisks
type VirtualDisks struct {
	PhysicalDisks         []types.String       `tfsdk:"physical_disks"`
	VirtualDiskFqdd       types.String         `tfsdk:"virtual_disk_fqdd"`
	RaidLevel             types.String         `tfsdk:"raid_level"`
	RollUpStatus          types.String         `tfsdk:"roll_up_status"`
	Controller            types.String         `tfsdk:"controller"`
	ControllerProductName types.String         `tfsdk:"controller_product_name"`
	Configuration         ConfigurationDetails `tfsdk:"configuration"`
	MediaType             types.String         `tfsdk:"media_type"`
	EncryptionType        types.String         `tfsdk:"encryption_type"`
}

// ExternalVirtualDisks is the tfsdk model of ExternalVirtualDisks
type ExternalVirtualDisks struct {
	PhysicalDisks         []types.String       `tfsdk:"physical_disks"`
	VirtualDiskFqdd       types.String         `tfsdk:"virtual_disk_fqdd"`
	RaidLevel             types.String         `tfsdk:"raid_level"`
	RollUpStatus          types.String         `tfsdk:"roll_up_status"`
	Controller            types.String         `tfsdk:"controller"`
	ControllerProductName types.String         `tfsdk:"controller_product_name"`
	Configuration         ConfigurationDetails `tfsdk:"configuration"`
	MediaType             types.String         `tfsdk:"media_type"`
	EncryptionType        types.String         `tfsdk:"encryption_type"`
}


// RaidConfiguration is the tfsdk model of RaidConfiguration
type RaidConfiguration struct {
	VirtualDisks         []VirtualDisks         `tfsdk:"virtual_disks"`
	ExternalVirtualDisks []ExternalVirtualDisks `tfsdk:"external_virtual_disks"`
	HddHotSpares         []types.String         `tfsdk:"hdd_hot_spares"`
	SsdHotSpares         []types.String         `tfsdk:"ssd_hot_spares"`
	ExternalHddHotSpares []types.String         `tfsdk:"external_hdd_hot_spares"`
	ExternalSsdHotSpares []types.String         `tfsdk:"external_ssd_hot_spares"`
	SizeToDiskMap        types.Map          `tfsdk:"size_to_disk_map"`
}


// OptionsDetails is the tfsdk model of OptionsDetails
type OptionsDetails struct {
	ID           types.String          `tfsdk:"id"`
	Name         types.String          `tfsdk:"name"`
	Value        types.String          `tfsdk:"value"`
	Dependencies []DependenciesDetails `tfsdk:"dependencies"`
	Attributes   types.Map            `tfsdk:"attributes"`
}

// ScaleIOStoragePoolDisks is the tfsdk model of ScaleIOStoragePoolDisks
type ScaleIOStoragePoolDisks struct {
	ProtectionDomainID   types.String   `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String   `tfsdk:"protection_domain_name"`
	StoragePoolID        types.String   `tfsdk:"storage_pool_id"`
	StoragePoolName      types.String   `tfsdk:"storage_pool_name"`
	DiskType             types.String   `tfsdk:"disk_type"`
	PhysicalDiskFqdds    []types.String `tfsdk:"physical_disk_fqdds"`
	VirtualDiskFqdds     []types.String `tfsdk:"virtual_disk_fqdds"`
	SoftwareOnlyDisks    []types.String `tfsdk:"software_only_disks"`
}

// ScaleIODiskConfiguration is the tfsdk model of ScaleIODiskConfiguration
type ScaleIODiskConfiguration struct {
	ScaleIOStoragePoolDisks []ScaleIOStoragePoolDisks `tfsdk:"scale_io_storage_pool_disks"`
}

// ShortWindow is the tfsdk model of ShortWindow
type ShortWindow struct {
	Threshold       types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec types.Int64 `tfsdk:"window_size_in_sec"`
}

// MediumWindow is the tfsdk model of MediumWindow
type MediumWindow struct {
	Threshold       types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec types.Int64 `tfsdk:"window_size_in_sec"`
}

// LongWindow is the tfsdk model of LongWindow
type LongWindow struct {
	Threshold       types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec types.Int64 `tfsdk:"window_size_in_sec"`
}

// SdsDecoupledCounterParameters is the tfsdk model of SdsDecoupledCounterParameters
type SdsDecoupledCounterParameters struct {
	ShortWindow  ShortWindow  `tfsdk:"short_window"`
	MediumWindow MediumWindow `tfsdk:"medium_window"`
	LongWindow   LongWindow   `tfsdk:"long_window"`
}

// SdsConfigurationFailureCounterParameters is the tfsdk model of SdsConfigurationFailureCounterParameters
type SdsConfigurationFailureCounterParameters struct {
	ShortWindow  ShortWindow  `tfsdk:"short_window"`
	MediumWindow MediumWindow `tfsdk:"medium_window"`
	LongWindow   LongWindow   `tfsdk:"long_window"`
}

// MdmSdsCounterParameters is the tfsdk model of MdmSdsCounterParameters
type MdmSdsCounterParameters struct {
	ShortWindow  ShortWindow  `tfsdk:"short_window"`
	MediumWindow MediumWindow `tfsdk:"medium_window"`
	LongWindow   LongWindow   `tfsdk:"long_window"`
}

// SdsSdsCounterParameters is the tfsdk model of SdsSdsCounterParameters
type SdsSdsCounterParameters struct {
	ShortWindow  ShortWindow  `tfsdk:"short_window"`
	MediumWindow MediumWindow `tfsdk:"medium_window"`
	LongWindow   LongWindow   `tfsdk:"long_window"`
}

// SdsReceiveBufferAllocationFailuresCounterParameters is the tfsdk model of SdsReceiveBufferAllocationFailuresCounterParameters
type SdsReceiveBufferAllocationFailuresCounterParameters struct {
	ShortWindow  ShortWindow  `tfsdk:"short_window"`
	MediumWindow MediumWindow `tfsdk:"medium_window"`
	LongWindow   LongWindow   `tfsdk:"long_window"`
}

// General is the tfsdk model of General
type General struct {
	ID                                                  types.String                                        `tfsdk:"id"`
	Name                                                types.String                                        `tfsdk:"name"`
	SystemID                                            types.String                                        `tfsdk:"system_id"`
	ProtectionDomainState                               types.String                                        `tfsdk:"protection_domain_state"`
	RebuildNetworkThrottlingInKbps                      types.Int64                                         `tfsdk:"rebuild_network_throttling_in_kbps"`
	RebalanceNetworkThrottlingInKbps                    types.Int64                                         `tfsdk:"rebalance_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingInKbps                    types.Int64                                         `tfsdk:"overall_io_network_throttling_in_kbps"`
	SdsDecoupledCounterParameters                       SdsDecoupledCounterParameters                       `tfsdk:"sds_decoupled_counter_parameters"`
	SdsConfigurationFailureCounterParameters            SdsConfigurationFailureCounterParameters            `tfsdk:"sds_configuration_failure_counter_parameters"`
	MdmSdsCounterParameters                             MdmSdsCounterParameters                             `tfsdk:"mdm_sds_counter_parameters"`
	SdsSdsCounterParameters                             SdsSdsCounterParameters                             `tfsdk:"sds_sds_counter_parameters"`
	RfcacheOpertionalMode                               types.String                                        `tfsdk:"rfcache_opertional_mode"`
	RfcachePageSizeKb                                   types.Int64                                         `tfsdk:"rfcache_page_size_kb"`
	RfcacheMaxIoSizeKb                                  types.Int64                                         `tfsdk:"rfcache_max_io_size_kb"`
	SdsReceiveBufferAllocationFailuresCounterParameters SdsReceiveBufferAllocationFailuresCounterParameters `tfsdk:"sds_receive_buffer_allocation_failures_counter_parameters"`
	RebuildNetworkThrottlingEnabled                     types.Bool                                          `tfsdk:"rebuild_network_throttling_enabled"`
	RebalanceNetworkThrottlingEnabled                   types.Bool                                          `tfsdk:"rebalance_network_throttling_enabled"`
	OverallIoNetworkThrottlingEnabled                   types.Bool                                          `tfsdk:"overall_io_network_throttling_enabled"`
	RfcacheEnabled                                      types.Bool                                          `tfsdk:"rfcache_enabled"`
}

// StatisticsDetails is the tfsdk model of StatisticsDetails
type StatisticsDetails struct {
	NumOfDevices                             types.Int64 `tfsdk:"num_of_devices"`
	UnusedCapacityInKb                       types.Int64 `tfsdk:"unused_capacity_in_kb"`
	NumOfVolumes                             types.Int64 `tfsdk:"num_of_volumes"`
	NumOfMappedToAllVolumes                  types.Int64 `tfsdk:"num_of_mapped_to_all_volumes"`
	CapacityAvailableForVolumeAllocationInKb types.Int64 `tfsdk:"capacity_available_for_volume_allocation_in_kb"`
	VolumeAllocationLimitInKb                types.Int64 `tfsdk:"volume_allocation_limit_in_kb"`
	CapacityLimitInKb                        types.Int64 `tfsdk:"capacity_limit_in_kb"`
	NumOfUnmappedVolumes                     types.Int64 `tfsdk:"num_of_unmapped_volumes"`
	SpareCapacityInKb                        types.Int64 `tfsdk:"spare_capacity_in_kb"`
	CapacityInUseInKb                        types.Int64 `tfsdk:"capacity_in_use_in_kb"`
	MaxCapacityInKb                          types.Int64 `tfsdk:"max_capacity_in_kb"`
	NumOfSds                                 types.Int64 `tfsdk:"num_of_sds"`
	NumOfStoragePools                        types.Int64 `tfsdk:"num_of_storage_pools"`
	NumOfFaultSets                           types.Int64 `tfsdk:"num_of_fault_sets"`
	ThinCapacityInUseInKb                    types.Int64 `tfsdk:"thin_capacity_in_use_in_kb"`
	ThickCapacityInUseInKb                   types.Int64 `tfsdk:"thick_capacity_in_use_in_kb"`
}

// DiskList is the tfsdk model of DiskList
type DiskList struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ErrorState             types.String `tfsdk:"error_state"`
	SdsID                  types.String `tfsdk:"sds_id"`
	DeviceState            types.String `tfsdk:"device_state"`
	CapacityLimitInKb      types.Int64  `tfsdk:"capacity_limit_in_kb"`
	MaxCapacityInKb        types.Int64  `tfsdk:"max_capacity_in_kb"`
	StoragePoolID          types.String `tfsdk:"storage_pool_id"`
	DeviceCurrentPathName  types.String `tfsdk:"device_current_path_name"`
	DeviceOriginalPathName types.String `tfsdk:"device_original_path_name"`
	SerialNumber           types.String `tfsdk:"serial_number"`
	VendorName             types.String `tfsdk:"vendor_name"`
	ModelName              types.String `tfsdk:"model_name"`
}

// MappedSdcInfoDetails is the tfsdk model of MappedSdcInfoDetails
type MappedSdcInfoDetails struct {
	SdcIP         types.String `tfsdk:"sdc_ip"`
	SdcID         types.String `tfsdk:"sdc_id"`
	LimitBwInMbps types.Int64  `tfsdk:"limit_bw_in_mbps"`
	LimitIops     types.Int64  `tfsdk:"limit_iops"`
}

// VolumeList is the tfsdk model of VolumeList
type VolumeList struct {
	ID                types.String           `tfsdk:"id"`
	Name              types.String           `tfsdk:"name"`
	VolumeType        types.String           `tfsdk:"volume_type"`
	StoragePoolID     types.String           `tfsdk:"storage_pool_id"`
	DataLayout        types.String           `tfsdk:"data_layout"`
	CompressionMethod types.String           `tfsdk:"compression_method"`
	SizeInKb          types.Int64            `tfsdk:"size_in_kb"`
	MappedSdcInfo     []MappedSdcInfoDetails `tfsdk:"mapped_sdc_info"`
	VolumeClass       types.String           `tfsdk:"volume_class"`
}

// StoragePoolList is the tfsdk model of StoragePoolList
type StoragePoolList struct {
	ID                                               types.String      `tfsdk:"id"`
	Name                                             types.String      `tfsdk:"name"`
	RebuildIoPriorityPolicy                          types.String      `tfsdk:"rebuild_io_priority_policy"`
	RebalanceIoPriorityPolicy                        types.String      `tfsdk:"rebalance_io_priority_policy"`
	RebuildIoPriorityNumOfConcurrentIosPerDevice     types.Int64       `tfsdk:"rebuild_io_priority_num_of_concurrent_ios_per_device"`
	RebalanceIoPriorityNumOfConcurrentIosPerDevice   types.Int64       `tfsdk:"rebalance_io_priority_num_of_concurrent_ios_per_device"`
	RebuildIoPriorityBwLimitPerDeviceInKbps          types.Int64       `tfsdk:"rebuild_io_priority_bw_limit_per_device_in_kbps"`
	RebalanceIoPriorityBwLimitPerDeviceInKbps        types.Int64       `tfsdk:"rebalance_io_priority_bw_limit_per_device_in_kbps"`
	RebuildIoPriorityAppIopsPerDeviceThreshold       types.String      `tfsdk:"rebuild_io_priority_app_iops_per_device_threshold"`
	RebalanceIoPriorityAppIopsPerDeviceThreshold     types.String      `tfsdk:"rebalance_io_priority_app_iops_per_device_threshold"`
	RebuildIoPriorityAppBwPerDeviceThresholdInKbps   types.Int64       `tfsdk:"rebuild_io_priority_app_bw_per_device_threshold_in_kbps"`
	RebalanceIoPriorityAppBwPerDeviceThresholdInKbps types.Int64       `tfsdk:"rebalance_io_priority_app_bw_per_device_threshold_in_kbps"`
	RebuildIoPriorityQuietPeriodInMsec               types.Int64       `tfsdk:"rebuild_io_priority_quiet_period_in_msec"`
	RebalanceIoPriorityQuietPeriodInMsec             types.Int64       `tfsdk:"rebalance_io_priority_quiet_period_in_msec"`
	ZeroPaddingEnabled                               types.Bool        `tfsdk:"zero_padding_enabled"`
	BackgroundScannerMode                            types.String      `tfsdk:"background_scanner_mode"`
	BackgroundScannerBWLimitKBps                     types.Int64       `tfsdk:"background_scanner_bw_limit_k_bps"`
	UseRmcache                                       types.Bool        `tfsdk:"use_rmcache"`
	ProtectionDomainID                               types.String      `tfsdk:"protection_domain_id"`
	SpClass                                          types.String      `tfsdk:"sp_class"`
	UseRfcache                                       types.Bool        `tfsdk:"use_rfcache"`
	SparePercentage                                  types.Int64       `tfsdk:"spare_percentage"`
	RmcacheWriteHandlingMode                         types.String      `tfsdk:"rmcache_write_handling_mode"`
	ChecksumEnabled                                  types.Bool        `tfsdk:"checksum_enabled"`
	RebuildEnabled                                   types.Bool        `tfsdk:"rebuild_enabled"`
	RebalanceEnabled                                 types.Bool        `tfsdk:"rebalance_enabled"`
	NumOfParallelRebuildRebalanceJobsPerDevice       types.Int64       `tfsdk:"num_of_parallel_rebuild_rebalance_jobs_per_device"`
	CapacityAlertHighThreshold                       types.Int64       `tfsdk:"capacity_alert_high_threshold"`
	CapacityAlertCriticalThreshold                   types.Int64       `tfsdk:"capacity_alert_critical_threshold"`
	Statistics                                       StatisticsDetails `tfsdk:"statistics"`
	DataLayout                                       types.String      `tfsdk:"data_layout"`
	ReplicationCapacityMaxRatio                      types.String      `tfsdk:"replication_capacity_max_ratio"`
	MediaType                                        types.String      `tfsdk:"media_type"`
	DiskList                                         []DiskList        `tfsdk:"disk_list"`
	VolumeList                                       []VolumeList      `tfsdk:"volume_list"`
	FglAccpID                                        types.String      `tfsdk:"fgl_accp_id"`
}

// SdsListDetails is the tfsdk model of SdsListDetails
type SdsListDetails struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Port                types.Int64  `tfsdk:"port"`
	ProtectionDomainID  types.String `tfsdk:"protection_domain_id"`
	FaultSetID          types.String `tfsdk:"fault_set_id"`
	SoftwareVersionInfo types.String `tfsdk:"software_version_info"`
	SdsState            types.String `tfsdk:"sds_state"`
	MembershipState     types.String `tfsdk:"membership_state"`
	MdmConnectionState  types.String `tfsdk:"mdm_connection_state"`
	DrlMode             types.String `tfsdk:"drl_mode"`
	MaintenanceState    types.String `tfsdk:"maintenance_state"`
	PerfProfile         types.String `tfsdk:"perf_profile"`
	OnVMWare            types.Bool   `tfsdk:"on_vm_ware"`
	IPList              []IPList     `tfsdk:"ip_list"`
}

// SdrListDetails is the tfsdk model of SdrListDetails
type SdrListDetails struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Port                types.Int64  `tfsdk:"port"`
	ProtectionDomainID  types.String `tfsdk:"protection_domain_id"`
	SoftwareVersionInfo types.String `tfsdk:"software_version_info"`
	SdrState            types.String `tfsdk:"sdr_state"`
	MembershipState     types.String `tfsdk:"membership_state"`
	MdmConnectionState  types.String `tfsdk:"mdm_connection_state"`
	MaintenanceState    types.String `tfsdk:"maintenance_state"`
	PerfProfile         types.String `tfsdk:"perf_profile"`
	IPList              []IPList     `tfsdk:"ip_list"`
}

// AccelerationPool is the tfsdk model of AccelerationPool
type AccelerationPool struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	MediaType          types.String `tfsdk:"media_type"`
	Rfcache            types.Bool   `tfsdk:"rfcache"`
}

// ProtectionDomainSettings is the tfsdk model of ProtectionDomainSettings
type ProtectionDomainSettings struct {
	General          General            `tfsdk:"general"`
	Statistics       StatisticsDetails  `tfsdk:"statistics"`
	StoragePoolList  []StoragePoolList  `tfsdk:"storage_pool_list"`
	SdsList          []SdsListDetails   `tfsdk:"sds_list"`
	SdrList          []SdrListDetails   `tfsdk:"sdr_list"`
	AccelerationPool []AccelerationPool `tfsdk:"acceleration_pool"`
}

// FaultSetSettings is the tfsdk model of FaultSetSettings
type FaultSetSettings struct {
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	Name               types.String `tfsdk:"name"`
	ID                 types.String `tfsdk:"id"`
}

// Datacenter is the tfsdk model of Datacenter
type Datacenter struct {
	VcenterID      types.String `tfsdk:"vcenter_id"`
	DatacenterID   types.String `tfsdk:"datacenter_id"`
	DatacenterName types.String `tfsdk:"datacenter_name"`
}

// PortGroupOptions is the tfsdk model of PortGroupOptions
type PortGroupOptions struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// PortGroups is the tfsdk model of PortGroups
type PortGroups struct {
	ID               types.String       `tfsdk:"id"`
	DisplayName      types.String       `tfsdk:"display_name"`
	Vlan             types.Int64        `tfsdk:"vlan"`
	Name             types.String       `tfsdk:"name"`
	Value            types.String       `tfsdk:"value"`
	PortGroupOptions []PortGroupOptions `tfsdk:"port_group_options"`
}

// VdsSettings is the tfsdk model of VdsSettings
type VdsSettings struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	PortGroups  []PortGroups `tfsdk:"port_groups"`
}

// VdsNetworkMtuSizeConfiguration is the tfsdk model of VdsNetworkMtuSizeConfiguration
type VdsNetworkMtuSizeConfiguration struct {
	ID    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
}

// VdsNetworkMTUSizeConfiguration is the tfsdk model of VdsNetworkMTUSizeConfiguration
type VdsNetworkMTUSizeConfiguration struct {
	ID    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
}

// VdsConfiguration is the tfsdk model of VdsConfiguration
type VdsConfiguration struct {
	Datacenter                     Datacenter                       `tfsdk:"datacenter"`
	PortGroupOption                types.String                     `tfsdk:"port_group_option"`
	PortGroupCreationOption        types.String                     `tfsdk:"port_group_creation_option"`
	VdsSettings                    []VdsSettings                    `tfsdk:"vds_settings"`
	VdsNetworkMtuSizeConfiguration []VdsNetworkMtuSizeConfiguration `tfsdk:"vds_network_mtu_size_configuration"`
}

// NodeSelection is the tfsdk model of NodeSelection
type NodeSelection struct {
	ID            types.String `tfsdk:"id"`
	ServiceTag    types.String `tfsdk:"service_tag"`
	MgmtIPAddress types.String `tfsdk:"mgmt_ip_address"`
}

// ParametersDetails is the tfsdk model of ParametersDetails
type ParametersDetails struct {
	GUID                     types.String               `tfsdk:"guid"`
	ID                       types.String               `tfsdk:"id"`
	Type                     types.String               `tfsdk:"type"`
	DisplayName              types.String               `tfsdk:"display_name"`
	Value                    types.String               `tfsdk:"value"`
	ToolTip                  types.String               `tfsdk:"tool_tip"`
	Required                 types.Bool                 `tfsdk:"required"`
	RequiredAtDeployment     types.Bool                 `tfsdk:"required_at_deployment"`
	HideFromTemplate         types.Bool                 `tfsdk:"hide_from_template"`
	Dependencies             []DependenciesDetails      `tfsdk:"dependencies"`
	Group                    types.String               `tfsdk:"group"`
	ReadOnly                 types.Bool                 `tfsdk:"read_only"`
	Generated                types.Bool                 `tfsdk:"generated"`
	InfoIcon                 types.Bool                 `tfsdk:"info_icon"`
	Step                     types.Int64                `tfsdk:"step"`
	MaxLength                types.Int64                `tfsdk:"max_length"`
	Min                      types.Int64                `tfsdk:"min"`
	Max                      types.Int64                `tfsdk:"max"`
	NetworkIPAddressList     []NetworkIPAddressList     `tfsdk:"network_ip_address_list"`
	NetworkConfiguration     NetworkConfiguration       `tfsdk:"network_configuration"`
	RaidConfiguration        RaidConfiguration          `tfsdk:"raid_configuration"`
	Options                  []OptionsDetails           `tfsdk:"options"`
	OptionsSortable          types.Bool                 `tfsdk:"options_sortable"`
	PreservedForDeployment   types.Bool                 `tfsdk:"preserved_for_deployment"`
	ScaleIODiskConfiguration ScaleIODiskConfiguration   `tfsdk:"scale_io_disk_configuration"`
	ProtectionDomainSettings []ProtectionDomainSettings `tfsdk:"protection_domain_settings"`
	FaultSetSettings         []FaultSetSettings         `tfsdk:"fault_set_settings"`
	Attributes               types.Map                 `tfsdk:"attributes"`
	VdsConfiguration         VdsConfiguration           `tfsdk:"vds_configuration"`
	NodeSelection            NodeSelection              `tfsdk:"node_selection"`
}

// Resources is the tfsdk model of Resources
type Resources struct {
	GUID        types.String `tfsdk:"guid"`
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}

// Components is the tfsdk model of Components
type Components struct {
	ID                  types.String      `tfsdk:"id"`
	ComponentID         types.String      `tfsdk:"component_id"`
	Identifier          types.String      `tfsdk:"identifier"`
	ComponentValid      ComponentValid    `tfsdk:"component_valid"`
	Name                types.String      `tfsdk:"name"`
	HelpText            types.String      `tfsdk:"help_text"`
	ClonedFromID        types.String      `tfsdk:"cloned_from_id"`
	Teardown            types.Bool        `tfsdk:"teardown"`
	Type                types.String      `tfsdk:"type"`
	SubType             types.String      `tfsdk:"sub_type"`
	RelatedComponents   types.Map `tfsdk:"related_components"`
	Resources           []Resources       `tfsdk:"resources"`
	Brownfield          types.Bool        `tfsdk:"brownfield"`
	PuppetCertName      types.String      `tfsdk:"puppet_cert_name"`
	OsPuppetCertName    types.String      `tfsdk:"os_puppet_cert_name"`
	ManagementIPAddress types.String      `tfsdk:"management_ip_address"`
	SerialNumber        types.String      `tfsdk:"serial_number"`
	AsmGUID             types.String      `tfsdk:"asm_guid"`
	Cloned              types.Bool        `tfsdk:"cloned"`
	ConfigFile          types.String      `tfsdk:"config_file"`
	ManageFirmware      types.Bool        `tfsdk:"manage_firmware"`
	Instances           types.Int64       `tfsdk:"instances"`
	RefID               types.String      `tfsdk:"ref_id"`
	ClonedFromAsmGUID   types.String      `tfsdk:"cloned_from_asm_guid"`
	Changed             types.Bool        `tfsdk:"changed"`
	IP                  types.String      `tfsdk:"ip"`
}

// IPRange is the tfsdk model of IPRange
type IPRange struct {
	ID         types.String `tfsdk:"id"`
	StartingIP types.String `tfsdk:"starting_ip"`
	EndingIP   types.String `tfsdk:"ending_ip"`
	Role       types.String `tfsdk:"role"`
}

// StaticRoute is the tfsdk model of StaticRoute
type StaticRoute struct {
	StaticRouteSourceNetworkID      types.String `tfsdk:"static_route_source_network_id"`
	StaticRouteDestinationNetworkID types.String `tfsdk:"static_route_destination_network_id"`
	StaticRouteGateway              types.String `tfsdk:"static_route_gateway"`
	SubnetMask                      types.String `tfsdk:"subnet_mask"`
	DestinationIPAddress            types.String `tfsdk:"destination_ip_address"`
}

// StaticNetworkConfiguration is the tfsdk model of StaticNetworkConfiguration
type StaticNetworkConfiguration struct {
	Gateway      types.String  `tfsdk:"gateway"`
	Subnet       types.String  `tfsdk:"subnet"`
	PrimaryDNS   types.String  `tfsdk:"primary_dns"`
	SecondaryDNS types.String  `tfsdk:"secondary_dns"`
	DNSSuffix    types.String  `tfsdk:"dns_suffix"`
	IPRange      []IPRange     `tfsdk:"ip_range"`
	IPAddress    types.String  `tfsdk:"ip_address"`
	StaticRoute  []StaticRoute `tfsdk:"static_route"`
}

// Networks is the tfsdk model of Networks
type Networks struct {
	ID                         types.String               `tfsdk:"id"`
	Name                       types.String               `tfsdk:"name"`
	Description                types.String               `tfsdk:"description"`
	VlanID                     types.Int64                `tfsdk:"vlan_id"`
	StaticNetworkConfiguration StaticNetworkConfiguration `tfsdk:"static_network_configuration"`
	DestinationIPAddress       types.String               `tfsdk:"destination_ip_address"`
	Static                     types.Bool                 `tfsdk:"static"`
	Type                       types.String               `tfsdk:"type"`
}

// Options is the tfsdk model of Options
type Options struct {
	ID           types.String          `tfsdk:"id"`
	Name         types.String          `tfsdk:"name"`
	Dependencies []DependenciesDetails `tfsdk:"dependencies"`
	Attributes   types.Map            `tfsdk:"attributes"`
}

// Parameters is the tfsdk model of Parameters
type Parameters struct {
	ID               types.String          `tfsdk:"id"`
	Value            types.String          `tfsdk:"value"`
	DisplayName      types.String          `tfsdk:"display_name"`
	Type             types.String          `tfsdk:"type"`
	ToolTip          types.String          `tfsdk:"tool_tip"`
	Required         types.Bool            `tfsdk:"required"`
	HideFromTemplate types.Bool            `tfsdk:"hide_from_template"`
	DeviceType       types.String          `tfsdk:"device_type"`
	Dependencies     []DependenciesDetails `tfsdk:"dependencies"`
	Group            types.String          `tfsdk:"group"`
	ReadOnly         types.Bool            `tfsdk:"read_only"`
	Generated        types.Bool            `tfsdk:"generated"`
	InfoIcon         types.Bool            `tfsdk:"info_icon"`
	Step             types.Int64           `tfsdk:"step"`
	MaxLength        types.Int64           `tfsdk:"max_length"`
	Min              types.Int64           `tfsdk:"min"`
	Max              types.Int64           `tfsdk:"max"`
	Networks         []Networks            `tfsdk:"networks"`
	Options          []Options             `tfsdk:"options"`
	OptionsSortable  types.Bool            `tfsdk:"options_sortable"`
}

// Categories is the tfsdk model of Categories
type Categories struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	DeviceType  types.String `tfsdk:"device_type"`
}
