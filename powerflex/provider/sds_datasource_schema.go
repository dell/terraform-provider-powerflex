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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SdsDataSourceSchema is the schema for reading the sds data
var SdsDataSourceSchema schema.Schema = schema.Schema{
	Description: "This data-source can be used to fetch information related to Storage Data Servers from a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder identifier attribute.",
			MarkdownDescription: "Placeholder identifier attribute.",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "Protection Domain ID." +
				" Conflicts with 'protection_domain_name'.",
			MarkdownDescription: "Protection Domain ID." +
				" Conflicts with `protection_domain_name`.",
			Optional: true,
		},
		"protection_domain_name": schema.StringAttribute{
			Description: "Protection Domain Name." +
				" Conflicts with 'protection_domain_id'.",
			MarkdownDescription: "Protection Domain Name." +
				" Conflicts with `protection_domain_id`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_id")),
			},
		},
		"sds_ids": schema.ListAttribute{
			Description: "List of SDS IDs." +
				" Conflicts with 'sds_names'.",
			MarkdownDescription: "List of SDS IDs." +
				" Conflicts with `sds_names`.",
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
		"sds_names": schema.ListAttribute{
			Description: "List of SDS names." +
				" Conflicts with 'sds_ids'.",
			MarkdownDescription: "List of SDS names." +
				" Conflicts with `sds_ids`.",
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.List{
				listvalidator.ConflictsWith(path.MatchRoot("sds_ids")),
				listvalidator.SizeAtLeast(1),
			},
		},
		"sds_details": schema.ListNestedAttribute{
			Description:         "List of fetched SDS.",
			MarkdownDescription: "List of fetched SDS.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "SDS ID.",
						MarkdownDescription: "SDS ID.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "SDS name.",
						MarkdownDescription: "SDS name.",
						Computed:            true,
					},
					"ip_list": schema.ListNestedAttribute{
						Description:         "List of IPs associated with SDS.",
						MarkdownDescription: "List of IPs associated with SDS.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"ip": schema.StringAttribute{
									Description:         "SDS IP.",
									MarkdownDescription: "SDS IP.",
									Computed:            true,
								},
								"role": schema.StringAttribute{
									Description:         "SDS IP role.",
									MarkdownDescription: "SDS IP role.",
									Computed:            true,
								},
							},
						},
					},
					"port": schema.Int64Attribute{
						Description:         "SDS port.",
						MarkdownDescription: "SDS port.",
						Computed:            true,
					},
					"sds_state": schema.StringAttribute{
						Description:         "SDS state.",
						MarkdownDescription: "SDS state.",
						Computed:            true,
					},
					"membership_state": schema.StringAttribute{
						Description:         "Membership state.",
						MarkdownDescription: "Membership state.",
						Computed:            true,
					},
					"mdm_connection_state": schema.StringAttribute{
						Description:         "MDM connection state.",
						MarkdownDescription: "MDM connection state.",
						Computed:            true,
					},
					"drl_mode": schema.StringAttribute{
						Description:         "DRL mode.",
						MarkdownDescription: "DRL mode.",
						Computed:            true,
					},
					"rmcache_enabled": schema.BoolAttribute{
						Description:         "Whether RM cache is enabled or not.",
						MarkdownDescription: "Whether RM cache is enabled or not.",
						Computed:            true,
					},
					"rmcache_size": schema.Int64Attribute{
						Description:         "Indicates the size of Read RAM Cache on the specified SDS in KB.",
						MarkdownDescription: "Indicates the size of Read RAM Cache on the specified SDS in KB.",
						Computed:            true,
					},
					"rmcache_frozen": schema.BoolAttribute{
						Description:         "Indicates whether the Read RAM Cache is currently temporarily not in use.",
						MarkdownDescription: "Indicates whether the Read RAM Cache is currently temporarily not in use.",
						Computed:            true,
					},
					"on_vmware": schema.BoolAttribute{
						Description:         "Presence on VMware.",
						MarkdownDescription: "Presence on VMware.",
						Computed:            true,
					},
					"faultset_id": schema.StringAttribute{
						Description:         "Fault set ID.",
						MarkdownDescription: "Fault set ID.",
						Computed:            true,
					},
					"num_io_buffers": schema.Int64Attribute{
						Description:         "Number of IO buffers.",
						MarkdownDescription: "Number of IO buffers.",
						Computed:            true,
					},
					"rmcache_memory_allocation_state": schema.StringAttribute{
						Description:         "Indicates the state of the memory allocation process. Can be one of 'in progress' and 'done'.",
						MarkdownDescription: "Indicates the state of the memory allocation process. Can be one of `in progress` and `done`.",
						Computed:            true,
					},
					"performance_profile": schema.StringAttribute{
						Description:         "Performance profile.",
						MarkdownDescription: "Performance profile.",
						Computed:            true,
					},
					"software_version_info": schema.StringAttribute{
						Description:         "Software version information.",
						MarkdownDescription: "Software version information.",
						Computed:            true,
					},
					"configured_drl_mode": schema.StringAttribute{
						Description:         "Configured DRL mode.",
						MarkdownDescription: "Configured DRL mode.",
						Computed:            true,
					},
					"rfcache_enabled": schema.BoolAttribute{
						Description:         "Whether RF cache is enabled or not.",
						MarkdownDescription: "Whether RF cache is enabled or not.",
						Computed:            true,
					},
					"maintenance_state": schema.StringAttribute{
						Description:         "Maintenance state.",
						MarkdownDescription: "Maintenance state.",
						Computed:            true,
					},
					"maintenance_type": schema.StringAttribute{
						Description:         "Maintenance type.",
						MarkdownDescription: "Maintenance type.",
						Computed:            true,
					},
					"rfcache_error_low_resources": schema.BoolAttribute{
						Description:         "RF cache error for low resources.",
						MarkdownDescription: "RF cache error for low resources.",
						Computed:            true,
					},
					"rfcache_error_api_version_mismatch": schema.BoolAttribute{
						Description:         "RF cache error for API version mismatch.",
						MarkdownDescription: "RF cache error for API version mismatch.",
						Computed:            true,
					},
					"rfcache_error_inconsistent_cache_configuration": schema.BoolAttribute{
						Description:         "RF cache error for inconsistent cache configuration.",
						MarkdownDescription: "RF cache error for inconsistent cache configuration.",
						Computed:            true,
					},
					"rfcache_error_inconsistent_source_configuration": schema.BoolAttribute{
						Description:         "RF cache error for inconsistent source configuration.",
						MarkdownDescription: "RF cache error for inconsistent source configuration.",
						Computed:            true,
					},
					"rfcache_error_invalid_driver_path": schema.BoolAttribute{
						Description:         "RF cache error for invalid driver path.",
						MarkdownDescription: "RF cache error for invalid driver path.",
						Computed:            true,
					},
					"rfcache_error_device_does_not_exist": schema.BoolAttribute{
						Description:         "RF cache error for device does not exist.",
						MarkdownDescription: "RF cache error for device does not exist.",
						Computed:            true,
					},
					"authentication_error": schema.StringAttribute{
						Description:         "Authentication error.",
						MarkdownDescription: "Authentication error.",
						Computed:            true,
					},
					"fgl_num_concurrent_writes": schema.Int64Attribute{
						Description:         "FGL concurrent writes.",
						MarkdownDescription: "FGL concurrent writes.",
						Computed:            true,
					},
					"fgl_metadata_cache_state": schema.StringAttribute{
						Description:         "FGL metadata cache state.",
						MarkdownDescription: "FGL metadata cache state.",
						Computed:            true,
					},
					"fgl_metadata_cache_size": schema.Int64Attribute{
						Description:         "FGL metadata cache size.",
						MarkdownDescription: "FGL metadata cache size.",
						Computed:            true,
					},
					"num_restarts": schema.Int64Attribute{
						Description:         "Number of restarts.",
						MarkdownDescription: "Number of restarts.",
						Computed:            true,
					},
					"last_upgrade_time": schema.Int64Attribute{
						Description:         "Last time SDS was upgraded.",
						MarkdownDescription: "Last time SDS was upgraded.",
						Computed:            true,
					},
					"sds_decoupled": schema.SingleNestedAttribute{
						Description:         "SDS decoupled windows.",
						MarkdownDescription: "SDS decoupled windows.",
						Computed:            true,
						Attributes:          getSdsAllWindowParamsSchema(),
					},
					"sds_configuration_failure": schema.SingleNestedAttribute{
						Description:         "SDS configuration failure windows.",
						MarkdownDescription: "SDS configuration failure windows.",
						Computed:            true,
						Attributes:          getSdsAllWindowParamsSchema(),
					},
					"sds_receive_buffer_allocation_failures": schema.SingleNestedAttribute{
						Description:         "SDS receive buffer allocation failure windows.",
						MarkdownDescription: "SDS receive buffer allocation failure windows.",
						Computed:            true,
						Attributes:          getSdsAllWindowParamsSchema(),
					},
					"certificate_info": schema.SingleNestedAttribute{
						Description:         "Certificate Information.",
						MarkdownDescription: "Certificate Information.",
						Computed:            true,
						Attributes:          getCertificateInfoSchema(),
					},
					"raid_controllers": schema.ListNestedAttribute{
						Description:         "RAID controllers information.",
						MarkdownDescription: "RAID controllers information.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: getRaidControllersSchema(),
						},
					},
					"links": schema.ListNestedAttribute{
						Description:         "Specifies the links asscociated with SDS.",
						MarkdownDescription: "Specifies the links asscociated with SDS.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description:         "Specifies the relationship with the SDS.",
									MarkdownDescription: "Specifies the relationship with the SDS.",
									Computed:            true,
								},
								"href": schema.StringAttribute{
									Description:         "Specifies the exact path to fetch the details.",
									MarkdownDescription: "Specifies the exact path to fetch the details.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	},
}

// getSdsAllWindowParamsSchema defines the schema for different window types
func getSdsAllWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			Description:         "Short Window Parameters.",
			MarkdownDescription: "Short Window Parameters.",
			Computed:            true,
			Attributes:          getSdsWindowParamsSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			Description:         "Medium Window Parameters.",
			MarkdownDescription: "Medium Window Parameters.",
			Computed:            true,
			Attributes:          getSdsWindowParamsSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			Description:         "Long Window Parameters.",
			MarkdownDescription: "Long Window Parameters.",
			Computed:            true,
			Attributes:          getSdsWindowParamsSchema(),
		},
	}
}

// getSdsWindowParamsSchema defines the schema for windows parameters
func getSdsWindowParamsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			Description:         "Threshold.",
			MarkdownDescription: "Threshold.",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			Description:         "Window Size in seconds.",
			MarkdownDescription: "Window Size in seconds.",
			Computed:            true,
		},
		"last_oscillation_count": schema.Int64Attribute{
			Description:         "Last oscillation count.",
			MarkdownDescription: "Last oscillation count.",
			Computed:            true,
		},
		"last_oscillation_time": schema.Int64Attribute{
			Description:         "Last oscillation time.",
			MarkdownDescription: "Last oscillation time.",
			Computed:            true,
		},
		"max_failures_count": schema.Int64Attribute{
			Description:         "Maximum failures count.",
			MarkdownDescription: "Maximum failures count.",
			Computed:            true,
		},
	}
}

// getCertificateInfoSchema defines the schema for certificates
func getCertificateInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"subject": schema.StringAttribute{
			Description:         "Certificate subject.",
			MarkdownDescription: "Certificate subject.",
			Computed:            true,
		},
		"issuer": schema.StringAttribute{
			Description:         "Certificate issuer.",
			MarkdownDescription: "Certificate issuer.",
			Computed:            true,
		},
		"valid_from": schema.StringAttribute{
			Description:         "The start date of the certificate validity.",
			MarkdownDescription: "The start date of the certificate validity.",
			Computed:            true,
		},
		"valid_to": schema.StringAttribute{
			Description:         "The end date of the certificate validity.",
			MarkdownDescription: "The end date of the certificate validity.",
			Computed:            true,
		},
		"thumbprint": schema.StringAttribute{
			Description:         "Certificate thumbprint.",
			MarkdownDescription: "Certificate thumbprint.",
			Computed:            true,
		},
		"valid_from_asn1_format": schema.StringAttribute{
			Description:         "The start date of the Asn1 format.",
			MarkdownDescription: "The start date of the Asn1 format.",
			Computed:            true,
		},
		"valid_to_asn1_format": schema.StringAttribute{
			Description:         "The end date of the Asn1 format.",
			MarkdownDescription: "The end date of the Asn1 format.",
			Computed:            true,
		},
	}
}

// getRaidControllersSchema defines the schema for RAID controller
func getRaidControllersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"serial_number": schema.StringAttribute{
			Description:         "Serial number.",
			MarkdownDescription: "Serial number.",
			Computed:            true,
		},
		"model_name": schema.StringAttribute{
			Description:         "Model name.",
			MarkdownDescription: "Model name.",
			Computed:            true,
		},
		"vendor_name": schema.StringAttribute{
			Description:         "Vendor name.",
			MarkdownDescription: "Vendor name.",
			Computed:            true,
		},
		"firmware_version": schema.StringAttribute{
			Description:         "Firmware version.",
			MarkdownDescription: "Firmware version.",
			Computed:            true,
		},
		"driver_version": schema.StringAttribute{
			Description:         "Driver version.",
			MarkdownDescription: "Driver version.",
			Computed:            true,
		},
		"driver_name": schema.StringAttribute{
			Description:         "Driver name.",
			MarkdownDescription: "Driver name.",
			Computed:            true,
		},
		"pci_address": schema.StringAttribute{
			Description:         "PCI address.",
			MarkdownDescription: "PCI address.",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			Description:         "RAID status.",
			MarkdownDescription: "RAID status.",
			Computed:            true,
		},
		"battery_status": schema.StringAttribute{
			Description:         "Battery status.",
			MarkdownDescription: "Battery status",
			Computed:            true,
		},
	}
}
