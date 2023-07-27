/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package provider

import (
	"terraform-provider-powerflex/powerflex/helper"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var ClusterReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to install the PowerFlex Cluster.",
	MarkdownDescription: "This resource can be used to install the PowerFlex Cluster.",
	Attributes:          ClusterResourceModelSchema(),
}

func ClusterResourceModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Optional:            true,
			Computed:            true,
		},

		"cluster": schema.ListNestedAttribute{
			MarkdownDescription: "Cluster Installation Details",
			Description:         "Cluster Installation Details",
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ClusterInstallationDetailsDataModelSchema()},
		},

		"storage_pools": schema.ListNestedAttribute{
			MarkdownDescription: "Storage Pool Details",
			Description:         "Storage Pool Details",
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: StoragePoolDetailsDataModelSchema()},
		},

		"mdm_password": schema.StringAttribute{
			MarkdownDescription: "MDM Password",
			Description:         "MDM Password",
			Required:            true,
		},

		"lia_password": schema.StringAttribute{
			MarkdownDescription: "Lia Password",
			Description:         "Lia Password",
			Required:            true,
		},

		"allow_non_secure_communication_with_mdm": schema.BoolAttribute{
			MarkdownDescription: "Allow Non Secure Communication With MDM",
			Description:         "Allow Non Secure Communication With MDM",
			Optional:            true,
		},

		"allow_non_secure_communication_with_lia": schema.BoolAttribute{
			MarkdownDescription: "Allow Non Secure Communication With lia",
			Description:         "Allow Non Secure Communication With lia",
			Optional:            true,
		},

		"disable_non_mgmt_components_auth": schema.BoolAttribute{
			MarkdownDescription: "Disable Non Mgmt Components Auth",
			Description:         "Disable Non Mgmt Components Auth",
			Optional:            true,
		},

		"mdm_list": schema.SetNestedAttribute{
			MarkdownDescription: "Cluster MDM Details",
			Description:         "Cluster MDM Details",
			Computed:            true,
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ClusterMDMDetailsDataModelSchema()},
		},

		"sds_list": schema.SetNestedAttribute{
			MarkdownDescription: "Cluster SDS Details",
			Description:         "Cluster SDS Details",
			Computed:            true,
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ClusterSDSDetailsDataModelSchema()},
		},

		"sdc_list": schema.SetNestedAttribute{
			MarkdownDescription: "Cluster SDC Details",
			Description:         "Cluster SDC Details",
			Computed:            true,
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ClusterSDCDetailsDataModelSchema()},
		},

		"sdr_list": schema.SetNestedAttribute{
			MarkdownDescription: "Cluster SDR Details",
			Description:         "Cluster SDR Details",
			Computed:            true,
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ClusterSDRDetailsDataModelSchema()},
		},

		"protection_domains": schema.ListNestedAttribute{
			MarkdownDescription: "Cluster Protection Domain Details",
			Description:         "Cluster Protection Domain Details",
			Computed:            true,
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ClusterProtectionDomainDetailsDataModelSchema()},
		},
	}
}

func ClusterInstallationDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"ips": schema.StringAttribute{
			MarkdownDescription: "IP address to be used for multiple purposes. Use this field to designate one IP address that will be assigned to all of the following: MDM IP, MDM Mgmt IP and SDS All IP. This option is provided for use cases where separate networks for data and management are not required.",
			Description:         "IP address to be used for multiple purposes. Use this field to designate one IP address that will be assigned to all of the following: MDM IP, MDM Mgmt IP and SDS All IP. This option is provided for use cases where separate networks for data and management are not required.",
			Optional:            true,
		},

		"username": schema.StringAttribute{
			MarkdownDescription: "Enter name either root or a non-root sudo user",
			Description:         "Enter name either root or a non-root sudo user",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("root"),
			},
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Password used to log in to the node.",
			Description:         "Password used to log in to the node.",
			Optional:            true,
		},

		"operating_system": schema.StringAttribute{
			MarkdownDescription: "Operating System",
			Description:         "Operating System",
			Optional:            true,
		},

		"is_mdm_or_tb": schema.StringAttribute{
			MarkdownDescription: "Is Mdm Or Tb",
			Description:         "Is Mdm Or Tb",
			Optional:            true,
		},

		"mdm_ips": schema.StringAttribute{
			MarkdownDescription: "MDM IP addresses used to communicate with other PowerFlex components in the storage network. This is required for all MDMs, Tiebreakers and Standbys.Leave this field blank for hosts that are not part of the MDM cluster.",
			Description:         "MDM IP addresses used to communicate with other PowerFlex components in the storage network. This is required for all MDMs, Tiebreakers and Standbys.Leave this field blank for hosts that are not part of the MDM cluster.",
			Optional:            true,
		},

		"mdm_mgmt_ip": schema.StringAttribute{
			MarkdownDescription: "The IP address for the management-only network.The management IP address is not required for: Tiebreaker, Standby Tiebreaker, and any host that is not an MDM. In such cases, leave this field blank.",
			Description:         "The IP address for the management-only network.The management IP address is not required for: Tiebreaker, Standby Tiebreaker, and any host that is not an MDM. In such cases, leave this field blank.",
			Optional:            true,
		},

		"mdm_name": schema.StringAttribute{
			MarkdownDescription: "MDMName",
			Description:         "MDMName",
			Optional:            true,
		},

		"perf_profile_for_mdm": schema.StringAttribute{
			MarkdownDescription: "Performance Profile For MDM",
			Description:         "Performance Profile For MDM",
			Optional:            true,
		},

		"virtual_ips": schema.StringAttribute{
			MarkdownDescription: "Virtual IPs",
			Description:         "Virtual IPs",
			Optional:            true,
		},

		"virtual_ip_nics": schema.StringAttribute{
			MarkdownDescription: "The NIC to which the virtual IP addresses are mapped.",
			Description:         "The NIC to which the virtual IP addresses are mapped.",
			Optional:            true,
		},

		"is_sds": schema.StringAttribute{
			MarkdownDescription: "Is Sds. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Description:         "Is Sds. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Yes",
				"No",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("No"),
			},
		},

		"sds_name": schema.StringAttribute{
			MarkdownDescription: "SDS Name",
			Description:         "SDS Name",
			Optional:            true,
		},

		"sds_all_ips": schema.StringAttribute{
			MarkdownDescription: "SDS IP addresses to be used for communication among all nodes.",
			Description:         "SDS IP addresses to be used for communication among all nodes.",
			Optional:            true,
		},

		"sds_to_sds_only_ips": schema.StringAttribute{
			MarkdownDescription: "SDS IP addresses to be used for communication among SDS nodes. When the replication feature is used, these addresses are also used for SDS-SDR communication.",
			Description:         "SDS IP addresses to be used for communication among SDS nodes. When the replication feature is used, these addresses are also used for SDS-SDR communication.",
			Optional:            true,
		},

		"sds_to_sdc_only_ips": schema.StringAttribute{
			MarkdownDescription: "SDS IP addresses to be used for communication among SDS and SDC nodes only.",
			Description:         "SDS IP addresses to be used for communication among SDS and SDC nodes only.",
			Optional:            true,
		},

		"protection_domain": schema.StringAttribute{
			MarkdownDescription: "Protection Domain",
			Description:         "Protection Domain",
			Optional:            true,
		},

		"fault_set": schema.StringAttribute{
			MarkdownDescription: "Fault Set",
			Description:         "Fault Set",
			Optional:            true,
		},

		"sds_storage_device_list": schema.StringAttribute{
			MarkdownDescription: "Storage devices to be added to an SDS. For more than one device, use a comma separated list, with no spaces.",
			Description:         "Storage devices to be added to an SDS. For more than one device, use a comma separated list, with no spaces.",
			Optional:            true,
		},

		"storage_pool_list": schema.StringAttribute{
			MarkdownDescription: "Sets Storage Pool names",
			Description:         "Sets Storage Pool names",
			Optional:            true,
		},

		"sds_storage_device_names": schema.StringAttribute{
			MarkdownDescription: "Sets names for devices.",
			Description:         "Sets names for devices.",
			Optional:            true,
		},

		"perf_profile_for_sds": schema.StringAttribute{
			MarkdownDescription: "Performance Profile For SDS",
			Description:         "Performance Profile For SDS",
			Optional:            true,
		},

		"is_sdc": schema.StringAttribute{
			MarkdownDescription: "Is Sdc. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Description:         "Is Sdc. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Yes",
				"No",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("No"),
			},
		},

		"perf_profile_for_sdc": schema.StringAttribute{
			MarkdownDescription: "Performance Profile For SDC",
			Description:         "Performance Profile For SDC",
			Optional:            true,
		},

		"sdc_name": schema.StringAttribute{
			MarkdownDescription: "SDC Name",
			Description:         "SDC Name",
			Optional:            true,
		},

		"is_rfcache": schema.StringAttribute{
			MarkdownDescription: "Is RFCache. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Description:         "Is RFCache. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Yes",
				"No",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("No"),
			},
		},

		"rf_cache_ssd_device_list": schema.StringAttribute{
			MarkdownDescription: "List of SSD devices to provide RFcache acceleration for Medium Granularity data layout Storage Pools.",
			Description:         "List of SSD devices to provide RFcache acceleration for Medium Granularity data layout Storage Pools.",
			Optional:            true,
		},

		"is_sdr": schema.StringAttribute{
			MarkdownDescription: "Is SDR. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Description:         "Is SDR. The acceptable values are `Yes` and `No`. Default value is `No`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Yes",
				"No",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("No"),
			},
		},

		"sdr_name": schema.StringAttribute{
			MarkdownDescription: "SDR Name",
			Description:         "SDR Name",
			Optional:            true,
		},

		"sdr_port": schema.StringAttribute{
			MarkdownDescription: "SDR Port",
			Description:         "SDR Port",
			Optional:            true,
		},

		"sdr_application_ips": schema.StringAttribute{
			MarkdownDescription: "The IP addresses through which the SDC communicates with the SDR.",
			Description:         "The IP addresses through which the SDC communicates with the SDR.",
			Optional:            true,
		},

		"sdr_storage_ips": schema.StringAttribute{
			MarkdownDescription: "The IP addresses through which the SDR communicates with the MDM for server side control communications.",
			Description:         "The IP addresses through which the SDR communicates with the MDM for server side control communications.",
			Optional:            true,
		},

		"sdr_external_ips": schema.StringAttribute{
			MarkdownDescription: "The IP addresses through which the SDR communicates with peer systems SDRs",
			Description:         "The IP addresses through which the SDR communicates with peer systems SDRs",
			Optional:            true,
		},

		"sdr_all_ips": schema.StringAttribute{
			MarkdownDescription: "SDR IP addresses to be used for communication among all nodes (including all three roles)",
			Description:         "SDR IP addresses to be used for communication among all nodes (including all three roles)",
			Optional:            true,
		},
		"perf_profile_for_sdr": schema.StringAttribute{
			MarkdownDescription: "Performance Profile For SDR",
			Description:         "Performance Profile For SDR",
			Optional:            true,
		},
	}
}

func StoragePoolDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"protection_domain": schema.StringAttribute{
			MarkdownDescription: "Protection Domain",
			Description:         "Protection Domain",
			Optional:            true,
		},

		"storage_pool": schema.StringAttribute{
			MarkdownDescription: "Storage Pool",
			Description:         "Storage Pool",
			Optional:            true,
		},

		"media_type": schema.StringAttribute{
			MarkdownDescription: "Media Type",
			Description:         "Media Type",
			Optional:            true,
		},

		"extern_alacceleration": schema.StringAttribute{
			MarkdownDescription: "External Acceleration",
			Description:         "External Acceleration",
			Optional:            true,
		},

		"data_layout": schema.StringAttribute{
			MarkdownDescription: "Data Layout",
			Description:         "Data Layout",
			Optional:            true,
		},

		"zero_padding": schema.StringAttribute{
			MarkdownDescription: "Zero Padding",
			Description:         "Zero Padding",
			Optional:            true,
		},

		"compression_method": schema.StringAttribute{
			MarkdownDescription: "Compression Method",
			Description:         "Compression Method",
			Optional:            true,
		},

		"replication_journal_capacity_percentage": schema.StringAttribute{
			MarkdownDescription: "Replication Journal Capacity Percentage",
			Description:         "Replication Journal Capacity Percentage",
			Optional:            true,
		},
	}
}

func ClusterMDMDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"ip": schema.StringAttribute{
			MarkdownDescription: "IP",
			Description:         "IP",
			Computed:            true,
		},

		"mdm_ip": schema.StringAttribute{
			MarkdownDescription: "MDM IP",
			Description:         "MDM IP",
			Computed:            true,
		},

		"mgmt_ip": schema.StringAttribute{
			MarkdownDescription: "MGMTIP",
			Description:         "MGMTIP",
			Computed:            true,
		},

		"virtual_ip_nic": schema.StringAttribute{
			MarkdownDescription: "Virtual IPNIC",
			Description:         "Virtual IPNIC",
			Computed:            true,
		},

		"role": schema.StringAttribute{
			MarkdownDescription: "Role",
			Description:         "Role",
			Computed:            true,
		},
		"mode": schema.StringAttribute{
			MarkdownDescription: "Mode",
			Description:         "Mode",
			Computed:            true,
		},
	}
}

func ClusterSDSDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"ip": schema.StringAttribute{
			MarkdownDescription: "IP",
			Description:         "IP",
			Computed:            true,
		},

		"all_ips": schema.StringAttribute{
			MarkdownDescription: "All IP",
			Description:         "All IP",
			Computed:            true,
		},

		"sds_only_ips": schema.StringAttribute{
			MarkdownDescription: "SDSOnly IP",
			Description:         "SDSOnly IP",
			Computed:            true,
		},

		"sds_sdc_ips": schema.StringAttribute{
			MarkdownDescription: "SDSSDCIP",
			Description:         "SDSSDCIP",
			Computed:            true,
		},

		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "Protection Domain Name",
			Description:         "Protection Domain Name",
			Computed:            true,
		},

		"protection_domain_name": schema.StringAttribute{
			MarkdownDescription: "Protection Domain Name",
			Description:         "Protection Domain Name",
			Computed:            true,
		},

		"fault_set": schema.StringAttribute{
			MarkdownDescription: "Fault Set",
			Description:         "Fault Set",
			Computed:            true,
		},

		"devices": schema.SetNestedAttribute{
			MarkdownDescription: "Devices",
			Description:         "Devices",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeviceDetailModelSchema()},
		},
	}
}

func DeviceDetailModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"path": schema.StringAttribute{
			MarkdownDescription: "Path",
			Description:         "Path",
			Computed:            true,
		},

		"storage_pool": schema.StringAttribute{
			MarkdownDescription: "Storage Pool Name",
			Description:         "Storage Pool Name",
			Computed:            true,
		},

		"max_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "Max Capacity",
			Description:         "Max Capacity",
			Computed:            true,
		},
	}
}

func ClusterSDCDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"gui_id": schema.StringAttribute{
			MarkdownDescription: "GUIID",
			Description:         "GUIID",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"ip": schema.StringAttribute{
			MarkdownDescription: "IP",
			Description:         "IP",
			Computed:            true,
		},
	}
}

func ClusterSDRDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"ip": schema.StringAttribute{
			MarkdownDescription: "IP",
			Description:         "IP",
			Computed:            true,
		},

		"port": schema.Int64Attribute{
			MarkdownDescription: "Port",
			Description:         "Port",
			Computed:            true,
		},

		"application_ips": schema.StringAttribute{
			MarkdownDescription: "Application IP",
			Description:         "Application IP",
			Computed:            true,
		},

		"storage_ips": schema.StringAttribute{
			MarkdownDescription: "Storage IP",
			Description:         "Storage IP",
			Computed:            true,
		},

		"external_ips": schema.StringAttribute{
			MarkdownDescription: "External IP",
			Description:         "External IP",
			Computed:            true,
		},

		"all_ips": schema.StringAttribute{
			MarkdownDescription: "All IP",
			Description:         "All IP",
			Computed:            true,
		},
	}
}

func ClusterProtectionDomainDetailsDataModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"storage_pool_list": schema.ListNestedAttribute{
			MarkdownDescription: "Storage Pools",
			Description:         "Storage Pools",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: StoragePoolDetailModelSchema()},
		},
	}
}

func StoragePoolDetailModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},

		"media_type": schema.StringAttribute{
			MarkdownDescription: "Media Type",
			Description:         "Media Type",
			Computed:            true,
		},

		"extern_alacceleration": schema.StringAttribute{
			MarkdownDescription: "External Acceleration",
			Description:         "External Acceleration",
			Computed:            true,
		},

		"data_layout": schema.StringAttribute{
			MarkdownDescription: "Data Layout",
			Description:         "Data Layout",
			Computed:            true,
		},

		"zero_padding": schema.StringAttribute{
			MarkdownDescription: "Zero Padding",
			Description:         "Zero Padding",
			Computed:            true,
		},

		"compression_method": schema.StringAttribute{
			MarkdownDescription: "Compression Method",
			Description:         "Compression Method",
			Computed:            true,
		},

		"replication_journal_capacity_percentage": schema.Int64Attribute{
			MarkdownDescription: "Replication Journal Capacity Percentage",
			Description:         "Replication Journal Capacity Percentage",
			Computed:            true,
		},
	}
}
