package powerflex

import (
	"fmt"

	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// SDCItem maps the sdc_list schema data
type SDCItem struct {
	SdcID         types.String `tfsdk:"sdc_id"`
	LimitIops     types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName       types.String `tfsdk:"sdc_name"`
	AccessMode    types.String `tfsdk:"access_mode"`
}

// GetType returns the terraform Type of SDCItem
func (si SDCItem) GetType() map[string]attr.Type {
	return map[string]attr.Type{
		"sdc_id":           types.StringType,
		"limit_iops":       types.Int64Type,
		"limit_bw_in_mbps": types.Int64Type,
		"sdc_name":         types.StringType,
		"access_mode":      types.StringType,
	}
}

func validateSdcSet(sl []SDCItem) []error {
	errs := make([]error, 0)
	sm := make(map[string]int)
	for _, si := range sl {
		if si.SdcID.IsUnknown() && si.SdcName.IsUnknown() {
			continue
		}
		id := fmt.Sprintf("{id:%s, name:%s}", si.SdcID.ValueString(), si.SdcName.ValueString())
		sm[id]++
	}
	// raise errors for SDCs that have multiple entries in the set
	for id, count := range sm {
		if count == 1 {
			continue
		}
		errs = append(errs, fmt.Errorf("the SDC %s is found %d times in the list, but only 1 time expected", id, count))
	}
	return errs
}

// GetValue returns the terraform Value of SDCItem
func (si *SDCItem) GetValue() (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(si.GetType(), map[string]attr.Value{
		"sdc_id":           si.SdcID,
		"sdc_name":         si.SdcName,
		"access_mode":      si.AccessMode,
		"limit_iops":       si.LimitIops,
		"limit_bw_in_mbps": si.LimitBwInMbps,
	})
}

// GetSdcSetValueFromItems marshalls list of SDCItem to a terraform set
func GetSdcSetValueFromItems(sl []SDCItem) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	sdcInfoElemType := types.ObjectType{
		AttrTypes: SDCItem{}.GetType(),
	}
	if len(sl) == 0 {
		return types.SetNull(sdcInfoElemType), diags
	}
	objectSdcInfos := []attr.Value{}
	for _, si := range sl {
		objVal, dgs := si.GetValue()
		diags = append(diags, dgs...)
		objectSdcInfos = append(objectSdcInfos, objVal)
	}
	mappedSdcInfoVal, dgs := types.SetValue(sdcInfoElemType, objectSdcInfos)
	diags = append(diags, dgs...)
	return mappedSdcInfoVal, diags
}

// GetSdcSetValueFromInfo marshalls list of *pftypes.MappedSdcInfo to a terraform set
func GetSdcSetValueFromInfo(sl []*pftypes.MappedSdcInfo) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	sdcInfoElemType := types.ObjectType{
		AttrTypes: SDCItem{}.GetType(),
	}
	if len(sl) == 0 {
		return types.SetNull(sdcInfoElemType), diags
	}
	objectSdcInfos := []attr.Value{}
	for _, msi := range sl {
		si := SDCItem{
			SdcID:         types.StringValue(msi.SdcID),
			LimitIops:     types.Int64Value(int64(msi.LimitIops)),
			LimitBwInMbps: types.Int64Value(int64(msi.LimitBwInMbps)),
			SdcName:       types.StringValue(msi.SdcName),
			AccessMode:    types.StringValue(msi.AccessMode),
		}
		objVal, dgs := si.GetValue()
		diags = append(diags, dgs...)
		objectSdcInfos = append(objectSdcInfos, objVal)
	}
	mappedSdcInfoVal, dgs := types.SetValue(sdcInfoElemType, objectSdcInfos)
	diags = append(diags, dgs...)
	return mappedSdcInfoVal, diags
}

var sdcListSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	Computed:            true,
	Optional:            true,
	MarkdownDescription: "List of SDCs to be mapped to the volume. Exactly one of `sdc_id` or `sdc_name` must be specified.",
	DeprecationMessage:  "Please use sdc_volumes_mapping resource for mapping SDCs to the volumes/snapshots. This attribute will be removed in future release.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"sdc_id": schema.StringAttribute{
				Description: "The ID of the SDC." +
					" Conflicts with 'sdc_name'." +
					" Cannot be updated.",
				Optional: true,
				Computed: true,
				MarkdownDescription: "The ID of the SDC." +
					" Conflicts with `sdc_name`." +
					" Cannot be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_name")),
				},
			},
			"sdc_name": schema.StringAttribute{
				Description: "The Name of the SDC." +
					" Conflicts with 'sdc_id'." +
					" Cannot be updated.",
				Computed: true,
				Optional: true,
				MarkdownDescription: "The Name of the SDC." +
					" Conflicts with `sdc_id`." +
					" Cannot be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_id")),
				},
			},
			"limit_iops": schema.Int64Attribute{
				Description: "IOPS limit of the SDC." +
					" Valid values are 0 or integers greater than 10." +
					" 0 represents unlimited IOPS." +
					" Default value is 0.",
				Optional: true,
				Computed: true,
				MarkdownDescription: "IOPS limit of the SDC." +
					" Valid values are `0` or integers greater than `10`." +
					" `0` represents unlimited IOPS." +
					" Default value is `0`.",
			},
			"limit_bw_in_mbps": schema.Int64Attribute{
				Description: "Bandwidth limit in MBPS of the SDC." +
					" 0 represents unlimited IOPS." +
					" Default value is 0.",
				Optional: true,
				Computed: true,
				MarkdownDescription: "Bandwidth limit in MBPS of the SDC." +
					" `0` represents unlimited IOPS." +
					" Default value is `0`.",
			},
			"access_mode": schema.StringAttribute{
				Description: "The Access Mode of the SDC." +
					" Valid values are 'ReadOnly', 'ReadWrite' and 'NoAccess'." +
					" Default value is 'ReadOnly'",
				Computed: true,
				Optional: true,
				MarkdownDescription: "The Access Mode of the SDC." +
					" Valid values are `ReadOnly`, `ReadWrite` and `NoAccess`." +
					" Default value is `ReadOnly`",
				Validators: []validator.String{stringvalidator.OneOf(
					"ReadOnly",
					"ReadWrite",
					"NoAccess",
				)},
				PlanModifiers: []planmodifier.String{
					stringDefault("ReadOnly"),
				},
			},
		},
	},
}
