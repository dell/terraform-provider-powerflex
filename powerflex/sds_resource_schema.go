package powerflex

import (
	"fmt"

	"github.com/dell/goscaleio"
	types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SDSResourceSchema variable to define schema for the SDS resource
var SDSResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to manage Storage Device Servers on a PowerFlex array.",
	MarkdownDescription: "This resource can be used to manage Storage Device Servers on a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "The id of the SDS",
			Computed:            true,
			MarkdownDescription: "The id of the SDS",
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "ID of the Protection Domain under which the SDS will be created." +
				" Conflicts with 'protection_domain_name'." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "ID of the Protection Domain under which the SDS will be created." +
				" Conflicts with `protection_domain_name`." +
				" Cannot be updated.",
		},
		"protection_domain_name": schema.StringAttribute{
			Description: "Name of the Protection Domain under which the SDS will be created." +
				" Conflicts with 'protection_domain_id'." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "Name of the Protection Domain under which the SDS will be created." +
				" Conflicts with `protection_domain_id`." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_id")),
			},
		},
		"name": schema.StringAttribute{
			Description:         "Name of SDS.",
			Required:            true,
			MarkdownDescription: "Name of SDS.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"ip_list": schema.SetNestedAttribute{
			Description: "List of IPs to be assigned to the SDS." +
				fmt.Sprintf(" There must be atleast one IP with '%s' role or ", goscaleio.RoleAll) +
				fmt.Sprintf("atleast two IPs, one with role '%s' ", goscaleio.RoleSdcOnly) +
				fmt.Sprintf("and the other with role '%s'.", goscaleio.RoleSdsOnly),
			MarkdownDescription: "List of IPs to be assigned to the SDS." +
				fmt.Sprintf(" There must be atleast one IP with `%s` role or ", goscaleio.RoleAll) +
				fmt.Sprintf("atleast two IPs, one with role `%s` ", goscaleio.RoleSdcOnly) +
				fmt.Sprintf("and the other with role `%s`.", goscaleio.RoleSdsOnly),
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						Description:         "IP address to be assigned to the SDS.",
						MarkdownDescription: "IP address to be assigned to the SDS.",
						Required:            true,
					},
					"role": schema.StringAttribute{
						Description: "Role to be assigned to the IP address." +
							fmt.Sprintf(" Valid values are '%s', '%s' and '%s'.",
								goscaleio.RoleAll,
								goscaleio.RoleSdcOnly,
								goscaleio.RoleSdsOnly,
							),
						MarkdownDescription: "Role to be assigned to the IP address." +
							fmt.Sprintf(" Valid values are `%s`, `%s` and `%s`.",
								goscaleio.RoleAll,
								goscaleio.RoleSdcOnly,
								goscaleio.RoleSdsOnly,
							),
						Required: true,
						Validators: []validator.String{stringvalidator.OneOf(
							goscaleio.RoleAll,
							goscaleio.RoleSdcOnly,
							goscaleio.RoleSdsOnly,
						)},
					},
				},
			},
		},
		"drl_mode": schema.StringAttribute{
			Description:         "DRL mode of SDS",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "DRL mode of SDS",
		},
		"rmcache_frozen": schema.BoolAttribute{
			Description:         "RMcache frozen state of SDS",
			Computed:            true,
			MarkdownDescription: "RMcache frozen state of SDS",
		},
		"fault_set_id": schema.StringAttribute{
			Description:         "Fault set id of SDS",
			Computed:            true,
			MarkdownDescription: "Fault set id of SDS",
		},
		"rmcache_memory_allocation_state": schema.StringAttribute{
			Description:         "Rmcache memory allocation state of SDS.",
			Computed:            true,
			MarkdownDescription: "Rmcache memory allocation state of SDS.",
		},
		"port": schema.Int64Attribute{
			Description:         "Port of SDS",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Port of SDS",
		},
		"membership_state": schema.StringAttribute{
			Description:         "Membership state of SDS",
			Computed:            true,
			MarkdownDescription: "Membership state of SDS",
		},
		"rmcache_enabled": schema.BoolAttribute{
			Description:         "Rmcache enabled state of SDS",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Rmcache enabled state of SDS",
		},
		"performance_profile": schema.StringAttribute{
			Description: "Performance Profile of SDS. " +
				fmt.Sprintf("Valid values are '%s' and '%s'.", types.PerformanceProfileCompact, types.PerformanceProfileHigh) +
				" Default value is determined by array settings.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "Performance Profile of SDS. " +
				fmt.Sprintf("Valid values are `%s` and `%s`.", types.PerformanceProfileCompact, types.PerformanceProfileHigh) +
				" Default value is determined by array settings.",
			Validators: []validator.String{stringvalidator.OneOf(
				types.PerformanceProfileHigh,
				types.PerformanceProfileCompact,
			)},
		},
		"rfcache_enabled": schema.BoolAttribute{
			Description:         "Rfcache enabled state of SDS",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Rfcache enabled state of SDS",
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
		"rmcache_size_in_mb": schema.Int64Attribute{
			Description:         "Read RAM cache size in MB of SDS. Can be set only when 'rmcache_enabled' is true.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Read RAM cache size in MB of SDS. Can be set only when `rmcache_enabled` is true.",
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
