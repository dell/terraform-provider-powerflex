package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SDSResourceSchema variable to define schema for the SDS resource
var SDSResourceSchema schema.Schema = schema.Schema{
	Description: "Manages SDS resource",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "The id of the SDS",
			Computed:            true,
			MarkdownDescription: "The id of the SDS",
		},
		"protection_domain_id": schema.StringAttribute{
			Description:         "Protection domain id",
			Required:            true,
			MarkdownDescription: "Protection domain id",
		},
		"name": schema.StringAttribute{
			Description:         "Name of SDS",
			Optional:            true,
			MarkdownDescription: "Name of SDS",
		},
		"ip_list": schema.ListAttribute{
			Description:         "IP list of SDS",
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "IP list of SDS",
		},
		"ip_role": schema.ListAttribute{
			Description:         "IP role list of SDS",
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "IP role list of SDS",
		},
		"drl_mode": schema.StringAttribute{
			Description:         "DRL mode of SDS",
			Computed:            true,
			MarkdownDescription: "DRL mode of SDS",
		},
		"rmcache_frozen": schema.BoolAttribute{
			Description:         "RMcache frozon state of SDS",
			Computed:            true,
			MarkdownDescription: "RMcache frozon state of SDS",
		},
		"fault_set_id": schema.StringAttribute{
			Description:         "Fault set id of SDS",
			Computed:            true,
			MarkdownDescription: "Fault set id of SDS",
		},
		"rmcache_memory_allocation_state": schema.StringAttribute{
			Description:         "Rmcache memory allocation state of SDS",
			Computed:            true,
			MarkdownDescription: "Rmcache memory allocation state of SDS",
		},
		"port": schema.Int64Attribute{
			Description:         "Port of SDS",
			Computed:            true,
			MarkdownDescription: "Port mode of SDS",
		},
		"membership_state": schema.StringAttribute{
			Description:         "Membership state of SDS",
			Computed:            true,
			MarkdownDescription: "Membership state of SDS",
		},
		"rmcache_enabled": schema.BoolAttribute{
			Description:         "Rmcache enabled state of SDS",
			Computed:            true,
			MarkdownDescription: "Rmcache enabled state of SDS",
		},
		"is_on_vmware": schema.BoolAttribute{
			Description:         "Is on vmware state of SDS",
			Computed:            true,
			MarkdownDescription: "Is on vmware state of SDS",
		},
		"mdm_connection_state": schema.StringAttribute{
			Description:         "Mdm connection state of SDS",
			Computed:            true,
			MarkdownDescription: "Mdm connection state of SDS",
		},
		"rmcache_size_in_kb": schema.Int64Attribute{
			Description:         "Rmcache size in kb of SDS",
			Computed:            true,
			MarkdownDescription: "Rmcache size in kb of SDS",
		},
		"num_of_io_buffers": schema.Int64Attribute{
			Description:         "Number of io buffers of SDS",
			Computed:            true,
			MarkdownDescription: "Number of io buffers of SDS",
		},
		"sds_state": schema.StringAttribute{
			Description:         "State of SDS",
			Computed:            true,
			MarkdownDescription: "State of SDS",
		},
	},
}
