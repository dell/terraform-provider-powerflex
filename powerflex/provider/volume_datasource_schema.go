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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// VolumeDataSourceSchema is the schema for reading the volume data
var VolumeDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This data-source can be used to fetch information related to volumes from a PowerFlex array.",
	MarkdownDescription: "This data-source can be used to fetch information related to volumes from a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier of the volume instance." +
				"  Conflicts with 'name', 'storage_pool_id' and  'storage_pool_name'.",
			MarkdownDescription: "Unique identifier of the volume instance." +
				"  Conflicts with `name`, `storage_pool_id` and  `storage_pool_name`.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRoot("storage_pool_id"), path.MatchRoot("name"), path.MatchRoot("storage_pool_name")),
			},
		},
		"name": schema.StringAttribute{
			Description: "Name of the volume." +
				"  Conflicts with 'id', 'storage_pool_id' and  'storage_pool_name'.",
			MarkdownDescription: "Name of the volume." +
				"  Conflicts with `id`, `storage_pool_id` and  `storage_pool_name`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRoot("storage_pool_id"), path.MatchRoot("id"), path.MatchRoot("storage_pool_name")),
			},
		},
		"storage_pool_id": schema.StringAttribute{
			Description: "Specifies the unique identifier of the storage pool." +
				"  Conflicts with 'id', 'name' and  'storage_pool_name'.",
			MarkdownDescription: "Specifies the unique identifier of the storage pool." +
				"  Conflicts with `id`, `name` and  `storage_pool_name`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRoot("storage_pool_name"), path.MatchRoot("id"), path.MatchRoot("name")),
			},
		},
		"storage_pool_name": schema.StringAttribute{
			Description: "Specifies the unique identifier of the storage pool." +
				"  Conflicts with 'id', 'name' and 'storage_pool_id'.",
			MarkdownDescription: "Specifies the unique identifier of the storage pool." +
				"  Conflicts with `id`, `name` and `storage_pool_id`.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRoot("storage_pool_id"), path.MatchRoot("id"), path.MatchRoot("name")),
			},
		},
		"volumes": schema.ListNestedAttribute{
			Description:         "List of volumes.",
			MarkdownDescription: "List of volumes.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "Unique identifier of the volume instance.",
						MarkdownDescription: "Unique identifier of the volume instance.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Name of the volume.",
						MarkdownDescription: "Name of the volume.",
						Computed:            true,
					},
					"creation_time": schema.Int64Attribute{
						Description:         "Specifies the time of creation.",
						MarkdownDescription: "Specifies the time of creation.",
						Computed:            true,
					},
					"size_in_kb": schema.Int64Attribute{
						Description:         "Size of the volume in KB",
						MarkdownDescription: "Size of the volume in KB",
						Computed:            true,
					},
					"ancestor_volume_id": schema.StringAttribute{
						Description:         "The volume id to which the snapshot is linked to.",
						MarkdownDescription: "The volume id to which the snapshot is linked to.",
						Computed:            true,
					},
					"vtree_id": schema.StringAttribute{
						Description:         "Unique identifier of the VTree",
						MarkdownDescription: "Unique identifier of the VTree",
						Computed:            true,
					},
					"consistency_group_id": schema.StringAttribute{
						Description:         "The unique id for the consistency group.",
						MarkdownDescription: "The unique id for the consistency group.",
						Computed:            true,
					},
					"volume_type": schema.StringAttribute{
						Description:         "Specifies the type of that volume.",
						MarkdownDescription: "Specifies the type of that volume.",
						Computed:            true,
					},
					"use_rm_cache": schema.BoolAttribute{
						Description:         "Enable rm cache.",
						MarkdownDescription: "Enable rm cache.",
						Computed:            true,
					},
					"storage_pool_id": schema.StringAttribute{
						Description:         "Specifies the unique identifier of the storage pool.",
						MarkdownDescription: "Specifies the unique identifier of the storage pool.",
						Computed:            true,
					},
					"data_layout": schema.StringAttribute{
						Description:         "Specifies the layout for the data.",
						MarkdownDescription: "Specifies the layout for the data.",
						Computed:            true,
					},
					"not_genuine_snapshot": schema.BoolAttribute{
						Description:         "Specifies if not genuine snapshot.",
						MarkdownDescription: "Specifies if not genuine snapshot.",
						Computed:            true,
					},
					"access_mode_limit": schema.StringAttribute{
						Description:         "Specifies the access mode limit.",
						MarkdownDescription: "Specifies the access mode limit.",
						Computed:            true,
					},
					"secure_snapshot_exp_time": schema.Int64Attribute{
						Description:         "Specifies the secure snapshot expiry time.",
						MarkdownDescription: "Specifies the secure snapshot expiry time.",
						Computed:            true,
					},
					"managed_by": schema.StringAttribute{
						Description:         "Specifies by whom it's managed by.",
						MarkdownDescription: "Specifies by whom it's managed by.",
						Computed:            true,
					},
					"locked_auto_snapshot": schema.BoolAttribute{
						Description:         "Specifies if it's a locked auto snapshot.",
						MarkdownDescription: "Specifies if it's a locked auto snapshot.",
						Computed:            true,
					},
					"locked_auto_snapshot_marked_for_removal": schema.BoolAttribute{
						Description:         "Specifies if it's a locked auto snapshot marked for removal.",
						MarkdownDescription: "Specifies if it's a locked auto snapshot marked for removal.",
						Computed:            true,
					},
					"compression_method": schema.StringAttribute{
						Description:         "Specifies the compression method.",
						MarkdownDescription: "Specifies the compression method.",
						Computed:            true,
					},
					"time_stamp_is_accurate": schema.BoolAttribute{
						Description:         "Specifies if the time stamp is accurate.",
						MarkdownDescription: "Specifies if the time stamp is accurate.",
						Computed:            true,
					},
					"original_expiry_time": schema.Int64Attribute{
						Description:         "Specifies the original expiry time.",
						MarkdownDescription: "Specifies the original expiry time.",
						Computed:            true,
					},
					"volume_replication_state": schema.StringAttribute{
						Description:         "Specifies the volume replication state.",
						MarkdownDescription: "Specifies the volume replication state.",
						Computed:            true,
					},
					"replication_journal_volume": schema.BoolAttribute{
						Description:         "Specifies the replication journal volume.",
						MarkdownDescription: "Specifies the replication journal volume.",
						Computed:            true,
					},
					"replication_time_stamp": schema.Int64Attribute{
						Description:         "Specifies the replication time stamp.",
						MarkdownDescription: "Specifies the replication time stamp.",
						Computed:            true,
					},
					"links": schema.ListNestedAttribute{
						Description:         "Specifies the links asscociated for a volume.",
						MarkdownDescription: "Specifies the links asscociated for a volume.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description:         "Specifies the relationship with the volume.",
									MarkdownDescription: "Specifies the relationship with the volume.",
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
					"mapped_sdc_info": schema.ListNestedAttribute{
						Description:         "Specifies the list of sdc's mapped to a volume.",
						MarkdownDescription: "Specifies the list of sdc's mapped to a volume.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"sdc_id": schema.StringAttribute{
									Description:         "Unique identifier for sdc.",
									MarkdownDescription: "Unique identifier for sdc.",
									Computed:            true,
								},
								"sdc_ip": schema.StringAttribute{
									Description:         "Ip of the sdc.",
									MarkdownDescription: "Ip of the sdc.",
									Computed:            true,
								},
								"limit_iops": schema.Int64Attribute{
									Description:         "Specifies the IOPS limits.",
									MarkdownDescription: "Specifies the IOPS limits.",
									Computed:            true,
								},
								"limit_bw_in_mbps": schema.Int64Attribute{
									Description:         "Specifies the bandwidth limits in Mbps.",
									MarkdownDescription: "Specifies the bandwidth limits in Mbps.",
									Computed:            true,
								},
								"sdc_name": schema.StringAttribute{
									Description:         "Specifies the name of the sdc.",
									MarkdownDescription: "Specifies the name of the sdc.",
									Computed:            true,
								},
								"access_mode": schema.StringAttribute{
									Description:         "Specifies the access mode.",
									MarkdownDescription: "Specifies the access mode.",
									Computed:            true,
								},
								"is_direct_buffer_mapping": schema.BoolAttribute{
									Description:         "Specifies if it is direct buffer mapping.",
									MarkdownDescription: "Specifies if it is direct buffer mapping.",
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
