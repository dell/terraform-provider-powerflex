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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PeerSystemReourceSchema - variable holds schema for PeerSystemReource resource
var PeerSystemReourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage the Peer System entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. We can Create, Update and Delete the PowerFlex Peer Systems using this resource. We can also Import an existing Peer Systems from the PowerFlex array. Peer system refers to the setup where multiple MDM nodes work together as peers to provide redundancy and high availability. This means that if one MDM node fails, other peer MDM nodes can take over its responsibilities, ensuring continuous operation without disruptions",
	MarkdownDescription: "This resource is used to manage the Peer System entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. We can Create, Update and Delete the PowerFlex Peer Systems using this resource. We can also Import an existing Peer Systems from the PowerFlex array. Peer system refers to the setup where multiple MDM nodes work together as peers to provide redundancy and high availability. This means that if one MDM node fails, other peer MDM nodes can take over its responsibilities, ensuring continuous operation without disruptions",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "Name of the peer mdm instance.",
			MarkdownDescription: "Name of the peer mdm instance.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"peer_system_id": schema.StringAttribute{
			Description:         "Unique identifier of the peer mdm system.",
			MarkdownDescription: "Unique identifier of the peer mdm system.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"port": schema.Int64Attribute{
			Description:         "Port of the peer mdm instance.",
			MarkdownDescription: "Port of the peer mdm instance.",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(7611),
		},
		"perf_profile": schema.StringAttribute{
			Description:         "Performance profile of the peer mdm instance.",
			MarkdownDescription: "Performance profile of the peer mdm instance.",
			Computed:            true,
			Optional:            true,
			Default:             stringdefault.StaticString("HighPerformance"),
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Compact",
				"HighPerformance",
			)},
		},
		"ip_list": schema.SetAttribute{
			Description:         "List of ips for the peer mdm system.",
			MarkdownDescription: "List of ips for the peer mdm system.",
			ElementType:         types.StringType,
			Required:            true,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"add_certificate": schema.BoolAttribute{
			Description:         "Flag that if set to true will attempt to add certificate of the peer mdm destination to source. This flag is only used during create.",
			MarkdownDescription: "Flag that if set to true will attempt to add certificate of the peer mdm destination to source. This flag is only used during create.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"source_primary_mdm_information": schema.SingleNestedAttribute{
			Description:         "Only used if add_certificate is set to true during create. The source primary mdm information to get the root certificate.",
			MarkdownDescription: "Only used if add_certificate is set to true during create. The source primary mdm information to get the root certificate.",
			Optional:            true,
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"ip": schema.StringAttribute{
					Description:         "ip of the primary source mdm instance.",
					MarkdownDescription: "ip of the primary source mdm instance.",
					Optional:            true,
				},
				"ssh_port": schema.StringAttribute{
					Description:         "port of the primary source mdm instance.",
					MarkdownDescription: "port of the primary source mdm instance.",
					Optional:            true,
					Computed:            true,
					Default:             stringdefault.StaticString("22"),
				},
				"ssh_username": schema.StringAttribute{
					Description:         "ssh username of the source primary mdm instance.",
					MarkdownDescription: "ssh username of the source primary mdm instance.",
					Optional:            true,
				},
				"ssh_password": schema.StringAttribute{
					Description:         "ssh password of the source primary mdm instance.",
					MarkdownDescription: "ssh password of the source primary mdm instance.",
					Optional:            true,
					Sensitive:           true,
				},
				"management_ip": schema.StringAttribute{
					Description:         "ip of the source management instance.",
					MarkdownDescription: "ip of the source management instance.",
					Optional:            true,
				},
				"management_username": schema.StringAttribute{
					Description:         "ssh username of the source management instance.",
					MarkdownDescription: "ssh username of the source management instance.",
					Optional:            true,
				},
				"management_password": schema.StringAttribute{
					Description:         "password of the source instance.",
					MarkdownDescription: "password of the source instance.",
					Optional:            true,
					Sensitive:           true,
				},
			},
		},
		"destination_primary_mdm_information": schema.SingleNestedAttribute{
			Description:         "Only used if add_certificate is set to true during create. The destination primary mdm information to get the root certificate.",
			MarkdownDescription: "Only used if add_certificate is set to true during create. The destination primary mdm information to get the root certificate.",
			Optional:            true,
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"ip": schema.StringAttribute{
					Description:         "ip of the primary destination mdm instance.",
					MarkdownDescription: "ip of the primary destination mdm instance.",
					Optional:            true,
				},
				"ssh_port": schema.StringAttribute{
					Description:         "port of the primary destination mdm instance.",
					MarkdownDescription: "port of the primary destination mdm instance.",
					Optional:            true,
					Computed:            true,
					Default:             stringdefault.StaticString("22"),
				},
				"ssh_username": schema.StringAttribute{
					Description:         "ssh username of the destination primary mdm instance.",
					MarkdownDescription: "ssh username of the destination primary mdm instance.",
					Optional:            true,
				},
				"ssh_password": schema.StringAttribute{
					Description:         "ssh password of the primary destination mdm instance.",
					MarkdownDescription: "ssh password of the primary destination mdm instance.",
					Optional:            true,
					Sensitive:           true,
				},
				"management_ip": schema.StringAttribute{
					Description:         "ip of the destination management instance.",
					MarkdownDescription: "ip of the destination management instance.",
					Optional:            true,
				},
				"management_username": schema.StringAttribute{
					Description:         "ssh username of the destination management instance.",
					MarkdownDescription: "ssh username of the destination management instance.",
					Optional:            true,
				},
				"management_password": schema.StringAttribute{
					Description:         "password of the management instance.",
					MarkdownDescription: "password of the management instance.",
					Optional:            true,
					Sensitive:           true,
				},
			},
		},
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the peer mdm instance.",
			MarkdownDescription: "Unique identifier of the peer mdm instance.",
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
	},
}
