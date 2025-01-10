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

// ReplicationPairsDataSourceSchema defines the schema for ReplicationPairs datasource
var ReplicationPairsDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to read the Replication Pairs entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. An RCG is a collection of multiple replication pairs that need to be replicated together in a coordinated and consistent manner. The key purpose is to ensure that all the data within the group is replicated in a consistent state. Applies to a group of data that needs to be kept consistent across the source and target",
	MarkdownDescription: "This datasource is used to read the Replication Pairs entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. An RCG is a collection of multiple replication pairs that need to be replicated together in a coordinated and consistent manner. The key purpose is to ensure that all the data within the group is replicated in a consistent state. Applies to a group of data that needs to be kept consistent across the source and target",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "default datasource id",
			MarkdownDescription: "default datasource id",
			Computed:            true,
		},
		"replication_pair_details": schema.ListNestedAttribute{
			Description:         "List of Replication Pairs",
			MarkdownDescription: "List of Replication Pairs",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: replicationPairDataAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.ReplicationPairFilter{})),
		},
	},
}

var replicationPairDataAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "Unique identifier of the replication pair instance.",
		MarkdownDescription: "Unique identifier of the replication pair instance.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name of the replication pair instance.",
		MarkdownDescription: "Name of the replication pair instance.",
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
	"replication_consistency_group_id": schema.StringAttribute{
		Description:         "Replication Consistency Group Id of the replication pair instance.",
		MarkdownDescription: "Replication Consistency Group Id of the replication pair instance.",
		Computed:            true,
	},
	"copy_type": schema.StringAttribute{
		Description:         "Copy Type of the replication pair instance.",
		MarkdownDescription: "Copy Type of the replication pair instance.",
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
}
