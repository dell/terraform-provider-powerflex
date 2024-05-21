/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FaultSetDataSourceSchema defines the schema for Fault Set datasource
var FirmwareRepositoryDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing firmware repository from PowerFlex array. The information fetched from this datasource can be used for getting the necessary details regarding the bundles and their components in that firmware repository.",
	MarkdownDescription: "This datasource is used to query the existing firmware repository from PowerFlex array. The information fetched from this datasource can be used for getting the necessary details regarding the bundles and their components in that firmware repository.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"firmware_repository_ids": schema.SetAttribute{
			Description:         "List of firmware repository IDs",
			MarkdownDescription: "List of firmware repository IDs",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.ConflictsWith(
					path.MatchRoot("firmware_repository_names"),
				),
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"firmware_repository_names": schema.SetAttribute{
			Description:         "List of firmware repository names",
			MarkdownDescription: "List of firmware repository names",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"firmware_repository_details": schema.ListNestedAttribute{
			Description:         "Firmware Repository details",
			MarkdownDescription: "Firmware Repository details",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "ID of the Firmware Repository",
						Description:         "ID of the Firmware Repository",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Firmware Repository name",
						Description:         "Firmware Repository name",
						Computed:            true,
					},
					"source_location": schema.StringAttribute{
						MarkdownDescription: "Source Location",
						Description:         "Source Location",
						Computed:            true,
					},
					"source_type": schema.StringAttribute{
						MarkdownDescription: "Source Type",
						Description:         "Source Type",
						Computed:            true,
					},
					"disk_location": schema.StringAttribute{
						MarkdownDescription: "Disk Location",
						Description:         "Disk Location",
						Computed:            true,
					},
					"filename": schema.StringAttribute{
						MarkdownDescription: "Filename",
						Description:         "Filename",
						Computed:            true,
					},
					"username": schema.StringAttribute{
						MarkdownDescription: "Username",
						Description:         "Username",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Password",
						Description:         "Password",
						Computed:            true,
					},
					"download_status": schema.StringAttribute{
						MarkdownDescription: "Download Status",
						Description:         "Download Status",
						Computed:            true,
					},
					"created_date": schema.StringAttribute{
						MarkdownDescription: "Created Date",
						Description:         "Created Date",
						Computed:            true,
					},
					"created_by": schema.StringAttribute{
						MarkdownDescription: "Created By",
						Description:         "Created By",
						Computed:            true,
					},
					"updated_date": schema.StringAttribute{
						MarkdownDescription: "Updated Date",
						Description:         "Updated Date",
						Computed:            true,
					},
					"updated_by": schema.StringAttribute{
						MarkdownDescription: "Updated By",
						Description:         "Updated By",
						Computed:            true,
					},
					"default_catalog": schema.BoolAttribute{
						MarkdownDescription: "Default Catalog",
						Description:         "Default Catalog",
						Computed:            true,
					},
					"embedded": schema.BoolAttribute{
						MarkdownDescription: "Embedded",
						Description:         "Embedded",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						MarkdownDescription: "State",
						Description:         "State",
						Computed:            true,
					},
					"software_components": schema.ListNestedAttribute{
						MarkdownDescription: "Software Components",
						Description:         "Software Components",
						Computed:            true,
						NestedObject:        schema.NestedAttributeObject{Attributes: ComponentSchema()},
					},
					// "software_bundles": schema.ListNestedAttribute{
					// 	MarkdownDescription: "Software Bundles",
					// 	Description:         "Software Bundles",
					// 	Computed:            true,
					// 	NestedObject:        schema.NestedAttributeObject{Attributes: BundleSchema()},
					// },
					"bundle_count": schema.Int64Attribute{
						MarkdownDescription: "Bundle Count",
						Description:         "Bundle Count",
						Computed:            true,
					},
					"component_count": schema.Int64Attribute{
						MarkdownDescription: "Component Count",
						Description:         "Component Count",
						Computed:            true,
					},
					"user_bundle_count": schema.Int64Attribute{
						MarkdownDescription: "User Bundle Count",
						Description:         "User Bundle Count",
						Computed:            true,
					},
					"minimal": schema.BoolAttribute{
						MarkdownDescription: "Minimal",
						Description:         "Minimal",
						Computed:            true,
					},
					"download_progress": schema.Int64Attribute{
						MarkdownDescription: "Download Progress",
						Description:         "Download Progress",
						Computed:            true,
					},
					"extract_progress": schema.Int64Attribute{
						MarkdownDescription: "Extract Progress",
						Description:         "Extract Progress",
						Computed:            true,
					},
					"file_size_in_gigabytes": schema.Float64Attribute{
						MarkdownDescription: "File Size In Gigabytes",
						Description:         "File Size In Gigabytes",
						Computed:            true,
					},
					"signature": schema.StringAttribute{
						MarkdownDescription: "Signature",
						Description:         "Signature",
						Computed:            true,
					},
					"custom": schema.BoolAttribute{
						MarkdownDescription: "Custom",
						Description:         "Custom",
						Computed:            true,
					},
					"needs_attention": schema.BoolAttribute{
						MarkdownDescription: "Needs Attention",
						Description:         "Needs Attention",
						Computed:            true,
					},
					"job_id": schema.StringAttribute{
						MarkdownDescription: "Job ID",
						Description:         "Job ID",
						Computed:            true,
					},
					"rcmapproved": schema.BoolAttribute{
						MarkdownDescription: "Rcmapproved",
						Description:         "Rcmapproved",
						Computed:            true,
					},
				},
			},
		},
	},
}

// ComponentSchema is a function that returns the schema for Component
func ComponentSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},
		"package_id": schema.StringAttribute{
			MarkdownDescription: "Package ID",
			Description:         "Package ID",
			Computed:            true,
		},
		"dell_version": schema.StringAttribute{
			MarkdownDescription: "Dell Version",
			Description:         "Dell Version",
			Computed:            true,
		},
		"vendor_version": schema.StringAttribute{
			MarkdownDescription: "Vendor Version",
			Description:         "Vendor Version",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "Component ID",
			Description:         "Component ID",
			Computed:            true,
		},
		"device_id": schema.StringAttribute{
			MarkdownDescription: "Device ID",
			Description:         "Device ID",
			Computed:            true,
		},
		"sub_device_id": schema.StringAttribute{
			MarkdownDescription: "Sub Device ID",
			Description:         "Sub Device ID",
			Computed:            true,
		},
		"vendor_id": schema.StringAttribute{
			MarkdownDescription: "Vendor ID",
			Description:         "Vendor ID",
			Computed:            true,
		},
		"sub_vendor_id": schema.StringAttribute{
			MarkdownDescription: "Sub Vendor ID",
			Description:         "Sub Vendor ID",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "Created Date",
			Description:         "Created Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "Created By",
			Description:         "Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "Updated Date",
			Description:         "Updated Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "Updated By",
			Description:         "Updated By",
			Computed:            true,
		},
		// "path": schema.StringAttribute{
		// 	MarkdownDescription: "Path",
		// 	Description:         "Path",
		// 	Computed:            true,
		// },
		// "hash_md5": schema.StringAttribute{
		// 	MarkdownDescription: "Hash Md5",
		// 	Description:         "Hash Md5",
		// 	Computed:            true,
		// },
		// "name": schema.StringAttribute{
		// 	MarkdownDescription: "Name",
		// 	Description:         "Name",
		// 	Computed:            true,
		// },
		// "category": schema.StringAttribute{
		// 	MarkdownDescription: "Category",
		// 	Description:         "Category",
		// 	Computed:            true,
		// },
		// "component_type": schema.StringAttribute{
		// 	MarkdownDescription: "Component Type",
		// 	Description:         "Component Type",
		// 	Computed:            true,
		// },
		// "operating_system": schema.StringAttribute{
		// 	MarkdownDescription: "Operating System",
		// 	Description:         "Operating System",
		// 	Computed:            true,
		// },
		// "system_ids": schema.ListAttribute{
		// 	MarkdownDescription: "System IDs",
		// 	Description:         "System IDs",
		// 	Computed:            true,
		// 	ElementType:         types.StringType,
		// },
		// "custom": schema.BoolAttribute{
		// 	MarkdownDescription: "Custom",
		// 	Description:         "Custom",
		// 	Computed:            true,
		// },
		// "needs_attention": schema.BoolAttribute{
		// 	MarkdownDescription: "Needs Attention",
		// 	Description:         "Needs Attention",
		// 	Computed:            true,
		// },
		// "ignore": schema.BoolAttribute{
		// 	MarkdownDescription: "Ignore",
		// 	Description:         "Ignore",
		// 	Computed:            true,
		// },
		// "original_component_id": schema.StringAttribute{
		// 	MarkdownDescription: "Original Component ID",
		// 	Description:         "Original Component ID",
		// 	Computed:            true,
		// },
		// "firmware_repo_name": schema.StringAttribute{
		// 	MarkdownDescription: "Firmware Repo Name",
		// 	Description:         "Firmware Repo Name",
		// 	Computed:            true,
		// },
	}
}

// // BundleSchema is a function that returns the schema for Bundle
// func BundleSchema() map[string]schema.Attribute {
// 	return map[string]schema.Attribute{
// 		"id": schema.StringAttribute{
// 			MarkdownDescription: "ID",
// 			Description:         "ID",
// 			Computed:            true,
// 		},
// 		"name": schema.StringAttribute{
// 			MarkdownDescription: "Name",
// 			Description:         "Name",
// 			Computed:            true,
// 		},
// 		"version": schema.StringAttribute{
// 			MarkdownDescription: "Version",
// 			Description:         "Version",
// 			Computed:            true,
// 		},
// 		"bundle_date": schema.StringAttribute{
// 			MarkdownDescription: "Bundle Date",
// 			Description:         "Bundle Date",
// 			Computed:            true,
// 		},
// 		"created_date": schema.StringAttribute{
// 			MarkdownDescription: "Created Date",
// 			Description:         "Created Date",
// 			Computed:            true,
// 		},
// 		"created_by": schema.StringAttribute{
// 			MarkdownDescription: "Created By",
// 			Description:         "Created By",
// 			Computed:            true,
// 		},
// 		"updated_date": schema.StringAttribute{
// 			MarkdownDescription: "Updated Date",
// 			Description:         "Updated Date",
// 			Computed:            true,
// 		},
// 		"updated_by": schema.StringAttribute{
// 			MarkdownDescription: "Updated By",
// 			Description:         "Updated By",
// 			Computed:            true,
// 		},
// 		"description": schema.StringAttribute{
// 			MarkdownDescription: "Description",
// 			Description:         "Description",
// 			Computed:            true,
// 		},
// 		"user_bundle": schema.BoolAttribute{
// 			MarkdownDescription: "User Bundle",
// 			Description:         "User Bundle",
// 			Computed:            true,
// 		},
// 		"user_bundle_path": schema.StringAttribute{
// 			MarkdownDescription: "User Bundle Path",
// 			Description:         "User Bundle Path",
// 			Computed:            true,
// 		},
// 		"device_type": schema.StringAttribute{
// 			MarkdownDescription: "Device Type",
// 			Description:         "Device Type",
// 			Computed:            true,
// 		},
// 		"device_model": schema.StringAttribute{
// 			MarkdownDescription: "Device Model",
// 			Description:         "Device Model",
// 			Computed:            true,
// 		},
// 		"fw_repository_id": schema.StringAttribute{
// 			MarkdownDescription: "Fw Repository ID",
// 			Description:         "Fw Repository ID",
// 			Computed:            true,
// 		},
// 		"bundle_type": schema.StringAttribute{
// 			MarkdownDescription: "Bundle Type",
// 			Description:         "Bundle Type",
// 			Computed:            true,
// 		},
// 		"custom": schema.BoolAttribute{
// 			MarkdownDescription: "Custom",
// 			Description:         "Custom",
// 			Computed:            true,
// 		},
// 		"needs_attention": schema.BoolAttribute{
// 			MarkdownDescription: "Needs Attention",
// 			Description:         "Needs Attention",
// 			Computed:            true,
// 		},
// 		"software_components": schema.ListNestedAttribute{
// 			MarkdownDescription: "Software Components",
// 			Description:         "Software Components",
// 			Computed:            true,
// 			NestedObject:        schema.NestedAttributeObject{Attributes: ComponentSchema()},
// 		},
// 	}
// }
