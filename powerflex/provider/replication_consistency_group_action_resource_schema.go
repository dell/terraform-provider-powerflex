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
	"terraform-provider-powerflex/powerflex/constants"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ReplicationConsistencyGroupActionReourceSchema - variable holds schema for ReplicationConsistencyGroupAction resource
var ReplicationConsistencyGroupActionReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to execute actions on the Replication Consistency Group entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.",
	MarkdownDescription: "This resource is used to execute actions on the Replication Consistency Group entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Replication Consistency Group ID",
			MarkdownDescription: "Replication Consistency Group ID",
			Required:            true,
		},
		"action": schema.StringAttribute{
			Description:         "Replication Consistency Group Action",
			MarkdownDescription: "Replication Consistency Group Action",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(constants.Sync),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				constants.Sync,
				constants.Restore,
				constants.Failover,
				constants.Reverse,
				constants.Switchover,
				constants.Snapshot,
			)},
		},
	},
}
