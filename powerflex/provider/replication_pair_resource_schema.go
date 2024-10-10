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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

// ReplicationPairReourceSchema - variable holds schema for ReplicationPair resource
var ReplicationPairReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Replication Pairs entity of the PowerFlex Array. We can Create, Update and Delete the PowerFlex Replication Pairs using this resource. We can also Import an existing Replication Pairs from the PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Replication Pairs entity of the PowerFlex Array. We can Create, Update and Delete the PowerFlex Replication Pairs using this resource. We can also Import an existing Replication Pairs from the PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "Replication Pair Name",
			MarkdownDescription: "Replication Pair Name",
			Required:            true,
		},
		"source_volume_id": schema.StringAttribute{
			Description:         "Source Volume ID",
			MarkdownDescription: "Source Volume ID",
			Required:            true,
		},
		"destination_volume_id": schema.StringAttribute{
			Description:         "Destination Volume ID",
			MarkdownDescription: "Destination Volume ID",
			Required:            true,
		},
		"replication_consistency_group_id": schema.StringAttribute{
			Description:         "Replication Consistancy Group ID",
			MarkdownDescription: "Replication Consistancy Group ID",
			Required:            true,
		},
		"pause_initial_copy": schema.BoolAttribute{
			Description:         "Pause Copy of the replication pair instance.",
			MarkdownDescription: "Pause Copy of the replication pair instance.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"copy_type": schema.StringAttribute{
			Description:         "Copy Type for Replication Pairs only value is OnlineCopy",
			MarkdownDescription: "Copy Type for Replication Pairs only value is OnlineCopy",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the replication pair instance.",
			MarkdownDescription: "Unique identifier of the replication pair instance.",
			Computed:            true,
		},
		"remote_id": schema.StringAttribute{
			Description:         "Remote ID of the replication pair instance.",
			MarkdownDescription: "Remote ID of the replication pair instance.",
			Computed:            true,
		},
		"user_requested_pause_transmit_init_copy": schema.BoolAttribute{
			Description:         "User Requested Pause of the replication pair instance.",
			MarkdownDescription: "User Requested Pause of the replication pair instance.",
			Computed:            true,
		},
		"remote_capacity_in_mb": schema.Int64Attribute{
			Description:         "Remote Capacity in MB of the replication pair instance.",
			MarkdownDescription: "Remote Capacity in MB of the replication pair instance.",
			Computed:            true,
		},
		"local_volume_id": schema.StringAttribute{
			Description:         "Local Volume ID of the replication pair instance.",
			MarkdownDescription: "Local Volume ID of the replication pair instance.",
			Computed:            true,
		},
		"remote_volume_id": schema.StringAttribute{
			Description:         "Remote Volume ID of the replication pair instance.",
			MarkdownDescription: "Remote Volume ID of the replication pair instance.",
			Computed:            true,
		},
		"remote_volume_name": schema.StringAttribute{
			Description:         "Remote Volume Name of the replication pair instance.",
			MarkdownDescription: "Remote Volume Name of the replication pair instance.",
			Computed:            true,
		},
		"lifetime_state": schema.StringAttribute{
			Description:         "Lifetime State of the replication pair instance.",
			MarkdownDescription: "Lifetime State of the replication pair instance.",
			Computed:            true,
		},
		"peer_system_name": schema.StringAttribute{
			Description:         "Peer System Name of the replication pair instance.",
			MarkdownDescription: "Peer System Name of the replication pair instance.",
			Computed:            true,
		},
		"initial_copy_state": schema.StringAttribute{
			Description:         "Initial Copy State of the replication pair instance.",
			MarkdownDescription: "Initial Copy State of the replication pair instance.",
			Computed:            true,
		},
		"initial_copy_priority": schema.Int64Attribute{
			Description:         "Initial Copy Priority of the replication pair instance.",
			MarkdownDescription: "Initial Copy Priority of the replication pair instance.",
			Computed:            true,
		},
	},
}
