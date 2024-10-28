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

// PeerMdmDataSourceSchema defines the schema for PeerMdm datasource
var PeerMdmDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to read the Peer MDM entity of the PowerFlex Array.",
	MarkdownDescription: "This datasource is used to read the Peer MDM entity of the PowerFlex Array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "default datasource id",
			MarkdownDescription: "default datasource id",
			Computed:            true,
		},
		"peer_system_details": schema.ListNestedAttribute{
			Description:         "List of Peer MDMs",
			MarkdownDescription: "List of Peer MDMs",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: peerMdmDataAttributes,
			},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.PeerMdmFilter{})),
		},
	},
}

var peerMdmDataAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description:         "Unique identifier of the peer mdm instance.",
		MarkdownDescription: "Unique identifier of the peer mdm instance.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		Description:         "Name of the peer mdm instance.",
		MarkdownDescription: "Name of the peer mdm instance.",
		Computed:            true,
	},
	"port": schema.Int64Attribute{
		Description:         "Port of the peer mdm instance.",
		MarkdownDescription: "Port of the peer mdm instance.",
		Computed:            true,
	},
	"peer_system_id": schema.StringAttribute{
		Description:         "Unique identifier of the peer mdm system.",
		MarkdownDescription: "Unique identifier of the peer mdm system.",
		Computed:            true,
	},
	"system_id": schema.StringAttribute{
		Description:         "Unique identifier of the peer mdm system.",
		MarkdownDescription: "Unique identifier of the peer mdm system.",
		Computed:            true,
	},
	"software_version_info": schema.StringAttribute{
		Description:         "Software version details of the peer mdm instance.",
		MarkdownDescription: "Software version details of the peer mdm instance.",
		Computed:            true,
	},
	"membership_state": schema.StringAttribute{
		Description:         "membership state of the peer mdm instance.",
		MarkdownDescription: "membership state of the peer mdm instance.",
		Computed:            true,
	},
	"perf_profile": schema.StringAttribute{
		Description:         "Performance profile of the peer mdm instance.",
		MarkdownDescription: "Performance profile of the peer mdm instance.",
		Computed:            true,
	},
	"network_type": schema.StringAttribute{
		Description:         "Network type of the peer mdm system.",
		MarkdownDescription: "Network type of the peer mdm system.",
		Computed:            true,
	},
	"coupling_rc": schema.StringAttribute{
		Description:         "Coupling return code number of the peer mdm system.",
		MarkdownDescription: "Coupling return code number of the peer mdm system.",
		Computed:            true,
	},
	"ip_list": schema.ListNestedAttribute{
		Description:         "List of ips for the peer mdm system.",
		MarkdownDescription: "List of ips for the peer mdm system.",
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"ip": schema.StringAttribute{
					Description:         "Specifies the ip.",
					MarkdownDescription: "Specifies the ip.",
					Computed:            true,
				},
			},
		},
	},
}
