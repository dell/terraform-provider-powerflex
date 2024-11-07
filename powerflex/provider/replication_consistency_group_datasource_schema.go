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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// ReplicationConsistencyGroupDataSourceSchema defines the schema for RCG datasource
var ReplicationConsistencyGroupDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to read the Replication Consistency Group entity of the PowerFlex Array.",
	MarkdownDescription: "This datasource is used to read the Replication Consistency Group entity of the PowerFlex Array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "default datasource id",
			MarkdownDescription: "default datasource id",
			Computed:            true,
		},
		"replication_consistency_group_details": schema.ListNestedAttribute{
			Description:         "List of Replication Consistency Group",
			MarkdownDescription: "List of Replication Consistency Group",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: replicationConsistencyGroupAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.ReplicationConsistencyGroupFilter{})),
		},
	},
}

var replicationConsistencyGroupAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "Unique identifier of the replication consistency group instance.",
		MarkdownDescription: "Unique identifier of the replication consistency group instance.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name of the replication consistency group instance.",
		MarkdownDescription: "Name of the replication consistency group instance.",
		Computed:            true,
	},
	"remote_id": schema.StringAttribute{
		Description:         "Remote ID of the replication consistency group instance.",
		MarkdownDescription: "Remote ID of the replication consistency group instance.",
		Computed:            true,
	},
	"rpo_in_seconds": schema.Int64Attribute{
		Description:         "rpoInSeconds of the replication consistency group instance.",
		MarkdownDescription: "rpoInSeconds of the replication consistency group instance.",
		Computed:            true,
	},
	"protection_domain_id": schema.StringAttribute{
		Description:         "Protection Domain ID of the replication consistency group instance.",
		MarkdownDescription: "Protection Domain ID of the replication consistency group instance.",
		Computed:            true,
	},
	"remote_protection_domain_id": schema.StringAttribute{
		Description:         "Remote Protection Domain ID of the replication consistency group instance.",
		MarkdownDescription: "Remote Protection Domain ID of the replication consistency group instance.",
		Computed:            true,
	},
	"destination_system_id": schema.StringAttribute{
		Description:         "Destination System ID of the replication consistency group instance.",
		MarkdownDescription: "Destination System ID of the replication consistency group instance.",
		Computed:            true,
	},
	"peer_mdm_id": schema.StringAttribute{
		Description:         "Peer MDM ID of the replication consistency group instance.",
		MarkdownDescription: "Peer MDM ID of the replication consistency group instance.",
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
	"curr_consist_mode": schema.StringAttribute{
		Description:         "Consistency Mode of the replication consistency group instance.",
		MarkdownDescription: "Consistency Mode of the replication consistency group instance.",
		Computed:            true,
	},
	"freeze_state": schema.StringAttribute{
		Description:         "Freeze State of the replication consistency group instance.",
		MarkdownDescription: "Freeze State of the replication consistency group instance.",
		Computed:            true,
	},
	"pause_mode": schema.StringAttribute{
		Description:         "Pause Mode of the replication consistency group instance.",
		MarkdownDescription: "Pause Mode of the replication consistency group instance.",
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
	"target_volume_access_mode": schema.StringAttribute{
		Description:         "Target Volume Access Mode of the replication consistency group instance.",
		MarkdownDescription: "Target Volume Access Mode of the replication consistency group instance.",
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
	"local_activity_state": schema.StringAttribute{
		Description:         "Local Activity State of the replication consistency group instance.",
		MarkdownDescription: "Local Activity State of the replication consistency group instance.",
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
}
