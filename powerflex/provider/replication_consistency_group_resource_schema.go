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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ReplicationConsistencyGroupReourceSchema - variable holds schema for ReplicationConsistencyGroup resource
var ReplicationConsistencyGroupReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Replication Consistency Group entity of the PowerFlex Array. We can Create, Update and Delete the PowerFlex Replication Consistency Group using this resource. We can also Import an existing Replication Consistency Group from the PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Replication Consistency Group entity of the PowerFlex Array. We can Create, Update and Delete the PowerFlex Replication Consistency Group using this resource. We can also Import an existing Replication Consistency Group from the PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"protection_domain_id": schema.StringAttribute{
			Description:         "Replication Consistency Group Protection Domain ID",
			MarkdownDescription: "Replication Consistency Group Protection Domain ID",
			Required:            true,
		},
		"remote_protection_domain_id": schema.StringAttribute{
			Description:         "Replication Consistency Group Remote Protection Domain ID",
			MarkdownDescription: "Replication Consistency Group Remote Protection Domain ID",
			Required:            true,
		},
		"destination_system_id": schema.StringAttribute{
			Description:         "Replication Consistency Group Destination System ID",
			MarkdownDescription: "Replication Consistency Group Destination System ID",
			Required:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Replication Consistency Group Name",
			MarkdownDescription: "Replication Consistency Group Name",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"rpo_in_seconds": schema.Int64Attribute{
			Description:         "Replication Consistency Group RPO in Seconds (min: 15, default 15. max: 3600)",
			MarkdownDescription: "Replication Consistency Group RPO in Seconds (min: 15, default 15, max: 3600)",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(15),
			Validators: []validator.Int64{
				int64validator.AtLeast(15),
				int64validator.AtMost(3600),
			},
		},
		"curr_consist_mode": schema.StringAttribute{
			Description:         "Consistency Mode of the replication consistency group instance.",
			MarkdownDescription: "Consistency Mode of the replication consistency group instance.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("Consistent"),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Consistent",
				"Inconsistent",
			)},
		},
		"freeze_state": schema.StringAttribute{
			Description:         "Freeze State of the replication consistency group instance.",
			MarkdownDescription: "Freeze State of the replication consistency group instance.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("Unfrozen"),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Unfrozen",
				"Frozen",
			)},
		},
		"pause_mode": schema.StringAttribute{
			Description:         "Pause Mode of the replication consistency group instance.",
			MarkdownDescription: "Pause Mode of the replication consistency group instance.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("Resume"),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Resume",
				"Pause",
			)},
		},
		"target_volume_access_mode": schema.StringAttribute{
			Description:         "Target Volume Access Mode of the replication consistency group instance.",
			MarkdownDescription: "Target Volume Access Mode of the replication consistency group instance.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("NoAccess"),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"NoAccess",
				"ReadOnly",
			)},
		},
		"local_activity_state": schema.StringAttribute{
			Description:         "Local Activity State of the replication consistency group instance.",
			MarkdownDescription: "Local Activity State of the replication consistency group instance.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("Active"),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Active",
				"Terminated",
			)},
		},
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the replication consistency group instance.",
			MarkdownDescription: "Unique identifier of the replication consistency group instance.",
			Computed:            true,
		},
		"peer_mdm_id": schema.StringAttribute{
			Description:         "Replication Consistency Group Peer Mdm ID",
			MarkdownDescription: "Replication Consistency Group Peer Mdm ID",
			Computed:            true,
		},
		"remote_id": schema.StringAttribute{
			Description:         "Remote ID of the replication consistency group instance.",
			MarkdownDescription: "Remote ID of the replication consistency group instance.",
			Computed:            true,
		},
		"remote_mdm_id": schema.StringAttribute{
			Description:         "Remote MDM ID of the replication consistency group instance.",
			MarkdownDescription: "Remote MDM ID of the replication consistency group instance.",
			Computed:            true,
		},
		"replication_direction": schema.StringAttribute{
			Description:         "Replication Direction of the replication consistency group instance.",
			MarkdownDescription: "Replication Direction of the replication consistency group instance.",
			Computed:            true,
		},
		"lifetime_state": schema.StringAttribute{
			Description:         "Lifetime State of the replication consistency group instance.",
			MarkdownDescription: "Lifetime State of the replication consistency group instance.",
			Computed:            true,
		},
		"snap_creation_in_progress": schema.BoolAttribute{
			Description:         "Snap Creation In Progress of the replication consistency group instance.",
			MarkdownDescription: "Snap Creation In Progress of the replication consistency group instance.",
			Computed:            true,
		},
		"last_snap_group_id": schema.StringAttribute{
			Description:         "Last Snap Group ID of the replication consistency group instance.",
			MarkdownDescription: "Last Snap Group ID of the replication consistency group instance.",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			Description:         "Type of the replication consistency group instance.",
			MarkdownDescription: "Type of the replication consistency group instance.",
			Computed:            true,
		},
		"disaster_recovery_state": schema.StringAttribute{
			Description:         "Disaster Recovery State of the replication consistency group instance.",
			MarkdownDescription: "Disaster Recovery State of the replication consistency group instance.",
			Computed:            true,
		},
		"remote_disaster_recovery_state": schema.StringAttribute{
			Description:         "Remote Disaster Recovery State of the replication consistency group instance.",
			MarkdownDescription: "Remote Disaster Recovery State of the replication consistency group instance.",
			Computed:            true,
		},
		"failover_type": schema.StringAttribute{
			Description:         "Failover Type of the replication consistency group instance.",
			MarkdownDescription: "Failover Type of the replication consistency group instance.",
			Computed:            true,
		},
		"failover_state": schema.StringAttribute{
			Description:         "Failover State of the replication consistency group instance.",
			MarkdownDescription: "Failover State of the replication consistency group instance.",
			Computed:            true,
		},
		"active_local": schema.BoolAttribute{
			Description:         "Active Local of the replication consistency group instance.",
			MarkdownDescription: "Active Local of the replication consistency group instance.",
			Computed:            true,
		},
		"active_remote": schema.BoolAttribute{
			Description:         "Active Remote of the replication consistency group instance.",
			MarkdownDescription: "Active Remote of the replication consistency group instance.",
			Computed:            true,
		},
		"abstract_state": schema.StringAttribute{
			Description:         "Abstract State of the replication consistency group instance.",
			MarkdownDescription: "Abstract State of the replication consistency group instance.",
			Computed:            true,
		},
		"error": schema.Int64Attribute{
			Description:         "Error of the replication consistency group instance.",
			MarkdownDescription: "Error of the replication consistency group instance.",
			Computed:            true,
		},
		"remote_activity_state": schema.StringAttribute{
			Description:         "Remote Activity State of the replication consistency group instance.",
			MarkdownDescription: "Remote Activity State of the replication consistency group instance.",
			Computed:            true,
		},
		"inactive_reason": schema.Int64Attribute{
			Description:         "Inactive Reason of the replication consistency group instance.",
			MarkdownDescription: "Inactive Reason of the replication consistency group instance.",
			Computed:            true,
		},
	},
}
