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

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SnapshotPolicyResourceSchema - variable holds schema for Snapshot Policy
var SnapshotPolicyResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Snapshot Policy entity of the PowerFlex Array. We can Create, Update and Delete the snapshot policy using this resource. We can also import an existing snapshot policy from the PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Snapshot Policy entity of the PowerFlex Array. We can Create, Update and Delete the snapshot policy using this resource. We can also import an existing snapshot policy from the PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the Snapshot Policy.",
			MarkdownDescription: "ID of the Snapshot Policy.",
			Computed:            true,
		},
		"num_of_retained_snapshots_per_level": schema.ListAttribute{
			Description:         "List which represents the number of snapshots per retention level.",
			MarkdownDescription: "List which represents the number of snapshots per retention level.",
			Required:            true,
			ElementType:         types.Int64Type,
		},
		"auto_snapshot_creation_cadence_in_min": schema.Int64Attribute{
			Description:         "The interval in minutes between two snapshots in the policy.",
			MarkdownDescription: "The interval in minutes between two snapshots in the policy.",
			Required:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the Snapshot Policy",
			MarkdownDescription: "Name of the Snapshot Policy",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"paused": schema.BoolAttribute{
			Description:         "Indicates that the snapshot policy should paused or not. Default value is `false`.",
			MarkdownDescription: "Indicates that the snapshot policy should paused or not. Default value is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"volume_ids": schema.SetAttribute{
			Description:         "List which represents the volume ids which is to be assigned to the snapshot policy.",
			MarkdownDescription: "List which represents the volume ids which is to be assigned to the snapshot policy.",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"remove_mode": schema.StringAttribute{
			Description:         "When removing the source volume from the policy, user should choose how to handle the snapshots created by the policy. Valid values are 'Remove' and 'Detach'. Default value is `Detach`.",
			MarkdownDescription: "When removing the source volume from the policy, user should choose how to handle the snapshots created by the policy. Valid values are 'Remove' and 'Detach'. Default value is `Detach`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"Remove",
				"Detach",
			)},
			Default: stringdefault.StaticString("Detach"),
		},
		"secure_snapshots": schema.BoolAttribute{
			Description: "The auto snapshots will be created as secure. They cannot be edited or removed prior to their policy expiration time. Default value is `false`." +
				" Cannot be updated.",
			MarkdownDescription: "The auto snapshots will be created as secure. They cannot be edited or removed prior to their policy expiration time. Default value is `false`." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
		"snapshot_access_mode": schema.StringAttribute{
			Description: "The Access mode of auto snapshot. Valid values are 'ReadOnly' and 'ReadWrite'. Default value is `ReadOnly`" +
				" Cannot be updated.",
			MarkdownDescription: "The Access mode of auto snapshot. Valid values are 'ReadOnly' and 'ReadWrite'. Default value is `ReadOnly`" +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{stringvalidator.OneOf(
				"ReadOnly",
				"ReadWrite",
			)},
			Default: stringdefault.StaticString("ReadOnly"),
		},
	},
}
