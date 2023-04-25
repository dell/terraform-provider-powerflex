package powerflex

import (
	"context"
	"fmt"

	types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	frameworkTypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// ProtectionDomainDataSourceSchema defines the schema for Protection Domain datasource
var ProtectionDomainResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to manage protection domains on a PowerFlex array.",
	MarkdownDescription: "This resource can be used to manage protection domains on a PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the protection domain instance.",
			MarkdownDescription: "Unique identifier of the protection domain instance.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Unique name of the protection domain instance.",
			MarkdownDescription: "Unique name of the protection domain instance.",
			Required:            true,
		},
		"active": schema.BoolAttribute{
			Description:         "Whether the PD should be in 'Active' state.",
			MarkdownDescription: "Whether the PD should be in `Active` state.",
			Computed:            true,
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolDefault(true),
			},
		},
		"state": schema.StringAttribute{
			Description:         "State of the PD.",
			MarkdownDescription: "State of the PD.",
			Computed:            true,
		},
		"rf_cache_accp_id": schema.StringAttribute{
			Description:         "ID of the Rf Cache Acceleration Pool attached to the PD.",
			MarkdownDescription: "ID of the Rf Cache Acceleration Pool attached to the PD.",
			Computed:            true,
		},
		"rf_cache_enabled": schema.BoolAttribute{
			Description:         "Whether SDS Rf Cache is enabled or not.",
			MarkdownDescription: "Whether SDS Rf Cache is enabled or not.",
			Computed:            true,
			Optional:            true,
		},
		"rf_cache_operational_mode": schema.StringAttribute{
			Description:         "Operational Mode of the SDS RF Cache.",
			MarkdownDescription: "Operational Mode of the SDS RF Cache.",
			Computed:            true,
			Optional:            true,
			// Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(types.PDRCModeRead),
					string(types.PDRCModeWrite),
					string(types.PDRCModeReadAndWrite),
					string(types.PDRCModeWriteMiss),
				),
			},
		},
		"rf_cache_page_size_kb": schema.Int64Attribute{
			Description:         "Page size of the SDS RF Cache in KB.",
			MarkdownDescription: "Page size of the SDS RF Cache in KB.",
			Computed:            true,
			Optional:            true,
		},
		"rf_cache_max_io_size_kb": schema.Int64Attribute{
			Description:         "Maximum IO of the SDS RF Cache in KB.",
			MarkdownDescription: "Maximum IO of the SDS RF Cache in KB.",
			Computed:            true,
			Optional:            true,
		},
		"fgl_default_num_concurrent_writes": schema.Int64Attribute{
			Description:         "Fine Granularity default number of concurrent writes.",
			MarkdownDescription: "Fine Granularity default number of concurrent writes.",
			Computed:            true,
			// Optional:            true,
		},
		"fgl_metadata_cache_enabled": schema.BoolAttribute{
			Description:         "Whether Fine Granularity Metadata Cache is enabled or not.",
			MarkdownDescription: "Whether Fine Granularity Metadata Cache is enabled or not.",
			Computed:            true,
			Optional:            true,
		},
		"fgl_default_metadata_cache_size": schema.Int64Attribute{
			Description:         "Fine Granularity Metadata Cache size.",
			MarkdownDescription: "Fine Granularity Metadata Cache size.",
			Computed:            true,
			Optional:            true,
		},
		"protected_maintenance_mode_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for protected maintenance mode in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for protected maintenance mode in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
		},
		"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for rebuilding in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for rebuilding in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
		},
		"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for rebalancing in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for rebalancing in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
		},
		"vtree_migration_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for vtree migration in KBps." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for vtree migration in KBps." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
		},
		"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
			Description: "Maximum allowed IO for protected maintenance mode in KBps. Must be greater than any other network throttling parameter." +
				" The value '0' represents unlimited bandwidth. The default value is '0'.",
			MarkdownDescription: "Maximum allowed IO for protected maintenance mode in KBps. Must be greater than any other network throttling parameter." +
				" The value `0` represents unlimited bandwidth. The default value is `0`.",
			Computed: true,
			Optional: true,
		},
		"replication_capacity_max_ratio": schema.Int64Attribute{
			Description:         "Maximum Replication Capacity Ratio.",
			MarkdownDescription: "Maximum Replication Capacity Ratio.",
			Computed:            true,
			// Optional:            true,
		},
		"links": schema.ListNestedAttribute{
			Description:         "Underlying REST API links.",
			MarkdownDescription: "Underlying REST API links.",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description:         "Specifies the relationship with the Protection Domain.",
						MarkdownDescription: "Specifies the relationship with the Protection Domain.",
						Computed:            true,
					},
					"href": schema.StringAttribute{
						Description:         "Specifies the exact path to fetch the details.",
						MarkdownDescription: "Specifies the exact path to fetch the details.",
						Computed:            true,
					},
				},
			},
		},
	},
}

// boolDefaultModifier is a plan modifier that sets a default value for a
// types.BoolType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type boolDefaultModifier struct {
	Default bool
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %t", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%t`", m.Default)
}

// PlanModifyBool runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m boolDefaultModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// If the value is unknown or known, do not set default value.
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		resp.PlanValue = frameworkTypes.BoolValue(m.Default)
	}
}

func boolDefault(defaultValue bool) planmodifier.Bool {
	return boolDefaultModifier{
		Default: defaultValue,
	}
}
