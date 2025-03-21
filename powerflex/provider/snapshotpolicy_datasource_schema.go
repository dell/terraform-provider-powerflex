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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SnapshotPolicyDataSourceSchema is the schema for reading the snapshot policy data
var SnapshotPolicyDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing snapshot policies from the PowerFlex array. The information fetched from this datasource can be used for getting the details.",
	MarkdownDescription: "This datasource is used to query the existing snapshot policies from the PowerFlex array. The information fetched from this datasource can be used for getting the details.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Unique identifier of the snapshot policy instance to fetch." +
				" Conflicts with 'name'.",
			MarkdownDescription: "Unique identifier of the snapshot policy instance to fetch." +
				" Conflicts with `name`.",
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"snapshotpolicies": schema.ListNestedAttribute{
			Description:         "List of snapshot policies.",
			MarkdownDescription: "List of snapshot policies.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "Unique identifier of the snapshot policy instance.",
						MarkdownDescription: "Unique identifier of the snapshot policy instance.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Name of the snapshot policy.",
						MarkdownDescription: "Name of the snapshot policy.",
						Computed:            true,
					},
					"snapshot_policy_state": schema.StringAttribute{
						Description:         "Specifies the current state of the snapshot policy.",
						MarkdownDescription: "Specifies the current state of the snapshot policy.",
						Computed:            true,
					},
					"auto_snapshot_creation_cadence_in_min": schema.Int64Attribute{
						Description:         "Auto snapshot creation cadence in min.",
						MarkdownDescription: "Auto snapshot creation cadence in min.",
						Computed:            true,
					},
					"max_vtree_auto_snapshots": schema.Int64Attribute{
						Description:         "Max vtree auto snapshots.",
						MarkdownDescription: "Max vtree auto snapshots.",
						Computed:            true,
					},
					"num_of_source_volumes": schema.Int64Attribute{
						Description:         "Number of source Volumes.",
						MarkdownDescription: "Number of source Volumes.",
						Computed:            true,
					},
					"num_of_expired_but_locked_snapshots": schema.Int64Attribute{
						Description:         "Number of expired but locked snapshots.",
						MarkdownDescription: "Number of expired but locked snapshots.",
						Computed:            true,
					},
					"num_of_creation_failures": schema.Int64Attribute{
						Description:         "Number of creation failures.",
						MarkdownDescription: "Number of creation failures.",
						Computed:            true,
					},
					"num_of_retained_snapshots_per_level": schema.ListAttribute{
						ElementType:         types.Int64Type,
						Description:         "Number of retained snapshots per level.",
						MarkdownDescription: "Number of retained snapshots per level.",
						Computed:            true,
					},
					"snapshot_access_mode": schema.StringAttribute{
						Description:         "Snapshot Access Mode.",
						MarkdownDescription: "Snapshot Access Mode.",
						Computed:            true,
					},
					"secure_snapshots": schema.BoolAttribute{
						Description:         "Secure snapshots.",
						MarkdownDescription: "Secure snapshots.",
						Computed:            true,
					},
					"time_of_last_auto_snapshot": schema.Int64Attribute{
						Description:         "Time of last auto snapshot.",
						MarkdownDescription: "Time of last auto snapshot.",
						Computed:            true,
					},
					"next_auto_snapshot_creation_time": schema.Int64Attribute{
						Description:         "Next auto snapshot creation time.",
						MarkdownDescription: "Next auto snapshot creation time.",
						Computed:            true,
					},
					"time_of_last_auto_snapshot_creation_failure": schema.Int64Attribute{
						Description:         "Time of last auto snapshot creation failure.",
						MarkdownDescription: "Time of last auto snapshot creation failure.",
						Computed:            true,
					},
					"last_auto_snapshot_creation_failure_reason": schema.StringAttribute{
						Description:         "Last auto snapshot creation failure reason.",
						MarkdownDescription: "Last auto snapshot creation failure reason.",
						Computed:            true,
					},
					"last_auto_snapshot_failure_in_first_level": schema.BoolAttribute{
						Description:         "Last auto snapshot failure in first level.",
						MarkdownDescription: "Last auto snapshot failure in first level.",
						Computed:            true,
					},
					"num_of_auto_snapshots": schema.Int64Attribute{
						Description:         "Number of auto snapshots.",
						MarkdownDescription: "Number of auto snapshots.",
						Computed:            true,
					},
					"num_of_locked_snapshots": schema.Int64Attribute{
						Description:         "Number of locked snapshots.",
						MarkdownDescription: "Number of locked snapshots.",
						Computed:            true,
					},
					"system_id": schema.StringAttribute{
						Description:         "System Identifier.",
						MarkdownDescription: "System Identifier.",
						Computed:            true,
					},
					"links": schema.ListNestedAttribute{
						Description:         "Specifies the links associated for a snapshot policy.",
						MarkdownDescription: "Specifies the links associated for a snapshot policy.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description:         "Specifies the relationship with the snapshot policy.",
									MarkdownDescription: "Specifies the relationship with the snapshot policy.",
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
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.SnapshotPolicyFilter{})),
		},
	},
}
