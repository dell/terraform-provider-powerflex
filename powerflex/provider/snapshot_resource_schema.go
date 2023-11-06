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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-powerflex/powerflex/helper"
)

// SnapshotResourceSchema variable to define schema for the snapshot resource
var SnapshotResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Snapshot of volumes on PowerFlex Array. We can Create, Update and Delete the snapshots using this resource. We can also import an existing snapshot from PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Snapshot of volumes on PowerFlex Array. We can Create, Update and Delete the snapshots using this resource. We can also import an existing snapshot from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the snapshot.",
			Required:            true,
			MarkdownDescription: "The name of the snapshot.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"volume_id": schema.StringAttribute{
			Description: "The ID of the volume from which snapshot is to be created." +
				" Conflicts with 'volume_name'." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "The ID of the volume from which snapshot is to be created." +
				" Conflicts with `volume_name`." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("volume_name")),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"volume_name": schema.StringAttribute{
			Description: "The volume name for which snapshot is created." +
				" Conflicts with 'volume_id'." +
				" Cannot be updated.",
			Optional: true,
			//Computed: true,
			MarkdownDescription: "The volume name for which snapshot is created." +
				" Conflicts with `volume_id`." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("volume_id")),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"access_mode": schema.StringAttribute{
			Description:         "The Access mode of snapshot. Valid values are 'ReadOnly' and 'ReadWrite'. Default value is 'ReadOnly'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The Access mode of snapshot. Valid values are `ReadOnly` and `ReadWrite`. Default value is `ReadOnly`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"ReadOnly",
				"ReadWrite",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("ReadOnly"),
			},
		},
		"id": schema.StringAttribute{
			Description:         "The ID of the snapshot.",
			Computed:            true,
			MarkdownDescription: "The ID of the snapshot.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"size": schema.Int64Attribute{
			Description: "Size of the snapshot. The unit of size is defined by 'capacity_unit'." +
				" The storage capacity of a snapshot must be a multiple of 8GB and cannot be decreased.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "Size of the snapshot. The unit of size is defined by `capacity_unit`." +
				" The storage capacity of a snapshot must be a multiple of 8GB and cannot be decreased.",
		},
		"capacity_unit": schema.StringAttribute{
			Description:         "Unit of capacity of the volume. Must be one of 'GB' and 'TB'. Default value is 'GB'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Unit of capacity of the volume. Must be one of `GB` and `TB`. Default value is `GB`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"GB",
				"TB",
			),
				stringvalidator.AlsoRequires(path.MatchRoot("size")),
			},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("GB"),
			},
		},
		"size_in_kb": schema.Int64Attribute{
			Description:         "Size in KB",
			Computed:            true,
			MarkdownDescription: "Size in KB",
		},
		"lock_auto_snapshot": schema.BoolAttribute{
			Description:         "lock auto snapshot",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "lock auto snapshot",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"desired_retention": schema.Int64Attribute{
			Description: "The minimum amount of time that the snapshot should be retained on the array starting at the time of apply." +
				" The unit is defined by 'retention_unit'." +
				" Cannot be decreased.",
			Optional: true,
			MarkdownDescription: "The minimum amount of time that the snapshot should be retained on the array starting at the time of apply." +
				" The unit is defined by `retention_unit`." +
				" Cannot be decreased.",
		},
		"retention_unit": schema.StringAttribute{
			Description:         "Retention unit of the snapshot. Valid values are 'hours' and 'days'. Default value is 'hours'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Retention unit of the snapshot. Valid values are `hours` and `days`. Default value is `hours`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"hours",
				"days",
			),
				stringvalidator.AlsoRequires(path.MatchRoot("desired_retention")),
			},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("hours"),
			},
		},
		"retention_in_min": schema.StringAttribute{
			Description:         "retention of snapshot in min",
			Computed:            true,
			MarkdownDescription: "retention of snapshot in min",
		},
		"remove_mode": schema.StringAttribute{
			Description:         "Remove mode of the snapshot. Valid values are 'ONLY_ME' and 'INCLUDING_DESCENDANTS'. Default value is 'ONLY_ME'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Remove mode of the snapshot. Valid values are `ONLY_ME` and `INCLUDING_DESCENDANTS`. Default value is `ONLY_ME`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"ONLY_ME",
				"INCLUDING_DESCENDANTS",
			)},
			PlanModifiers: []planmodifier.String{
				helper.StringDefault("ONLY_ME"),
			},
		},
	},
}
