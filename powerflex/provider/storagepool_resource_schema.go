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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// StoragepoolReourceSchema - varible holds schema for Storagepool
var StoragepoolReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Storage Pool entity of PowerFlex Array. We can Create, Update and Delete the storage pool using this resource. We can also import an existing storage pool from PowerFlex array.",
	MarkdownDescription: "This resource is used to manage the Storage Pool entity of PowerFlex Array. We can Create, Update and Delete the storage pool using this resource. We can also import an existing storage pool from PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "ID of the Storage pool",
			MarkdownDescription: "ID of the Storage pool",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "ID of the Protection Domain under which the storage pool will be created." +
				" Conflicts with 'protection_domain_name'." +
				" Cannot be updated.",
			MarkdownDescription: "ID of the Protection Domain under which the storage pool will be created." +
				" Conflicts with `protection_domain_name`." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_name")),
			},
		},
		"protection_domain_name": schema.StringAttribute{
			Description: "Name of the Protection Domain under which the storage pool will be created." +
				" Conflicts with 'protection_domain_id'." +
				" Cannot be updated.",
			MarkdownDescription: "Name of the Protection Domain under which the storage pool will be created." +
				" Conflicts with `protection_domain_id`." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the Storage pool",
			MarkdownDescription: "Name of the Storage pool",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"media_type": schema.StringAttribute{
			Description:         "Media Type of the storage pool. Valid values are 'HDD', 'SSD' and 'Transitional'",
			MarkdownDescription: "Media Type of the storage pool. Valid values are `HDD`, `SSD` and `Transitional`",
			Required:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"HDD",
				"SSD",
				"Transitional",
			)},
		},
		"use_rmcache": schema.BoolAttribute{
			Description:         "Enable/Disable RMcache on a specific storage pool",
			MarkdownDescription: "Enable/Disable RMcache on a specific storage pool",
			Optional:            true,
			Computed:            true,
		},
		"use_rfcache": schema.BoolAttribute{
			Description:         "Enable/Disable RFcache on a specific storage pool",
			MarkdownDescription: "Enable/Disable RFcache on a specific storage pool",
			Optional:            true,
			Computed:            true,
		},
		"zero_padding_enabled": schema.BoolAttribute{
			Description:         "Enable/Disable padding policy on a specific storage pool",
			MarkdownDescription: "Enable/Disable padding policy on a specific storage pool",
			Optional:            true,
			Computed:            true,
		},
		"replication_journal_capacity": schema.Int64Attribute{
			Description:         "This defines the maximum percentage of Storage Pool capacity that can be used by replication for the journal. Before deleting the storage pool, this has to be set to 0.",
			MarkdownDescription: "This defines the maximum percentage of Storage Pool capacity that can be used by replication for the journal. Before deleting the storage pool, this has to be set to 0.",
			Optional:            true,
			Computed:            true,
		},
		"capacity_alert_high_threshold": schema.Int64Attribute{
			Description:         "Set the threshold for triggering capacity usage high-priority alert.",
			MarkdownDescription: "Set the threshold for triggering capacity usage high-priority alert.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"capacity_alert_critical_threshold": schema.Int64Attribute{
			Description:         "Set the threshold for triggering capacity usage critical-priority alert.",
			MarkdownDescription: "Set the threshold for triggering capacity usage critical-priority alert.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"protected_maintenance_mode_io_priority_policy": schema.StringAttribute{
			Description:         "Set the I/O priority policy for protected maintenance mode for a specific Storage Pool. Valid values are `unlimited`, `limitNumOfConcurrentIos` and `favorAppIos`",
			MarkdownDescription: "Set the I/O priority policy for protected maintenance mode for a specific Storage Pool. Valid values are `unlimited`, `limitNumOfConcurrentIos` and `favorAppIos`",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"unlimited",
				"limitNumOfConcurrentIos",
				"favorAppIos",
			)},
		},
		"protected_maintenance_mode_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			Description:         "The maximum number of concurrent protected maintenance mode migration I/Os per device",
			MarkdownDescription: "The maximum number of concurrent protected maintenance mode migration I/Os per device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.AtMost(20),
			},
		},
		"protected_maintenance_mode_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			Description:         "The maximum bandwidth of protected maintenance mode migration I/Os, in KB per second, per device",
			MarkdownDescription: "The maximum bandwidth of protected maintenance mode migration I/Os, in KB per second, per device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1024),
				int64validator.AtMost(1048576),
			},
		},
		"rebalance_enabled": schema.BoolAttribute{
			Description:         "Enable or disable rebalancing in the specified Storage Pool",
			MarkdownDescription: "Enable or disable rebalancing in the specified Storage Pool",
			Optional:            true,
			Computed:            true,
		},
		"rebalance_io_priority_policy": schema.StringAttribute{
			Description:         "Policy to use for rebalance I/O priority. Valid values are `unlimited`, `limitNumOfConcurrentIos` and `favorAppIos`",
			MarkdownDescription: "Policy to use for rebalance I/O priority. Valid values are `unlimited`, `limitNumOfConcurrentIos` and `favorAppIos`",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"unlimited",
				"limitNumOfConcurrentIos",
				"favorAppIos",
			)},
		},
		"rebalance_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			Description:         "The maximum number of concurrent rebalance I/Os per device",
			MarkdownDescription: "The maximum number of concurrent rebalance I/Os per device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.AtMost(20),
			},
		},
		"rebalance_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			Description:         "The maximum bandwidth of rebalance I/Os, in KB/s, per device",
			MarkdownDescription: "The maximum bandwidth of rebalance I/Os, in KB/s, per device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1024),
				int64validator.AtMost(1048576),
			},
		},
		"vtree_migration_io_priority_policy": schema.StringAttribute{
			Description:         "Set the I/O priority policy for V-Tree migration for a specific Storage Pool. Valid values are `limitNumOfConcurrentIos` and `favorAppIos`",
			MarkdownDescription: "Set the I/O priority policy for V-Tree migration for a specific Storage Pool. Valid values are `limitNumOfConcurrentIos` and `favorAppIos`",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"limitNumOfConcurrentIos",
				"favorAppIos",
			)},
		},
		"vtree_migration_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			Description:         "The maximum number of concurrent V-Tree migration I/Os per device",
			MarkdownDescription: "The maximum number of concurrent V-Tree migration I/Os per device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.AtMost(20),
			},
		},
		"vtree_migration_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			Description:         "The maximum bandwidth of V-Tree migration IOs, in KB per second, per device",
			MarkdownDescription: "The maximum bandwidth of V-Tree migration IOs, in KB per second, per device",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1024),
				int64validator.AtMost(25600),
			},
		},
		"spare_percentage": schema.Int64Attribute{
			Description:         "Sets the spare capacity reservation policy",
			MarkdownDescription: "Sets the spare capacity reservation policy",
			Optional:            true,
			Computed:            true,
		},
		"rm_cache_write_handling_mode": schema.StringAttribute{
			Description:         "Sets the Read RAM Cache write handling mode of the specified Storage Pool",
			MarkdownDescription: "Sets the Read RAM Cache write handling mode of the specified Storage Pool",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"Cached",
				"Passthrough",
			)},
		},
		"rebuild_enabled": schema.BoolAttribute{
			Description:         "Enable or disable rebuilds in the specified Storage Pool",
			MarkdownDescription: "Enable or disable rebuilds in the specified Storage Pool",
			Optional:            true,
			Computed:            true,
		},
		"rebuild_rebalance_parallelism": schema.Int64Attribute{
			Description:         "Maximum number of concurrent rebuild and rebalance activities on SDSs in the Storage Pool",
			MarkdownDescription: "Maximum number of concurrent rebuild and rebalance activities on SDSs in the Storage Pool",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.AtMost(10),
			},
		},
		"fragmentation": schema.BoolAttribute{
			Description:         "Enable or disable fragmentation in the Storage Pool",
			MarkdownDescription: "Enable or disable fragmentation in the Storage Pool",
			Optional:            true,
			Computed:            true,
		},
	},
}
