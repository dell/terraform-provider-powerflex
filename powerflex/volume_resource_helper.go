package powerflex

import (
	"context"
	"fmt"

	"github.com/dell/goscaleio"
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

const (
	// MiKB to convert size in megabytes
	MiKB = 1024
	// GiKB to convert size in gigabytes
	GiKB = 1024 * MiKB
	// TiKB to convert size in terabytes
	TiKB = 1024 * GiKB
)

// VolumeResourceModel maps the resource schema data.
type VolumeResourceModel struct {
	ProtectionDomainName types.String `tfsdk:"protection_domain_name"`
	ProtectionDomainID   types.String `tfsdk:"protection_domain_id"`
	StoragePoolName      types.String `tfsdk:"storage_pool_name"`
	StoragePoolID        types.String `tfsdk:"storage_pool_id"`
	VolumeType           types.String `tfsdk:"volume_type"`
	UseRmCache           types.Bool   `tfsdk:"use_rm_cache"`
	CompressionMethod    types.String `tfsdk:"compression_method"`
	Size                 types.Int64  `tfsdk:"size"`
	CapacityUnit         types.String `tfsdk:"capacity_unit"`
	Name                 types.String `tfsdk:"name"`
	SizeInKb             types.Int64  `tfsdk:"size_in_kb"`
	ID                   types.String `tfsdk:"id"`
	AccessMode           types.String `tfsdk:"access_mode"`
	RemoveMode           types.String `tfsdk:"remove_mode"`
	SdcList              types.Set    `tfsdk:"sdc_list"`
}

// SDCItemize maps the sdc_list schema data
type SDCItemize struct {
	SdcID         string `tfsdk:"sdc_id"`
	LimitIops     int    `tfsdk:"limit_iops"`
	LimitBwInMbps int    `tfsdk:"limit_bw_in_mbps"`
	SdcName       string `tfsdk:"sdc_name"`
	AccessMode    string `tfsdk:"access_mode"`
}

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
	return SdcInfoAttrTypes
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

// covertToKB fucntion to convert size into kb
func convertToKB(capacityUnit string, size int64) int64 {
	var valInKiB int64
	switch capacityUnit {
	case "MB":
		valInKiB = size * MiKB
	case "TB":
		valInKiB = size * TiKB
	case "GB":
		valInKiB = size * GiKB
	}
	return valInKiB
}

// refreshState function to update the state of volume resource in terraform.tfstate file
func refreshVolumeState(vol *pftypes.Volume, state *VolumeResourceModel) (diags diag.Diagnostics) {
	state.StoragePoolID = types.StringValue(vol.StoragePoolID)
	state.UseRmCache = types.BoolValue(vol.UseRmCache)
	state.VolumeType = types.StringValue(vol.VolumeType)
	state.SizeInKb = types.Int64Value(int64(vol.SizeInKb))
	state.Name = types.StringValue(vol.Name)
	state.ID = types.StringValue(vol.ID)
	state.AccessMode = types.StringValue(vol.AccessModeLimit)
	state.CompressionMethod = types.StringValue(vol.CompressionMethod)
	mappedSdcInfoVal, diag2 := GetSdcSetValueFromInfo(vol.MappedSdcInfo)
	diags = append(diags, diag2...)
	state.SdcList = mappedSdcInfoVal
	return diags
}

// getStoragePoolInstance function to get storage pool from storage pool id and protection domain id
func getStoragePoolInstance(c *goscaleio.Client, spID string, pdID string) (*goscaleio.StoragePool, error) {
	getSystems, _ := c.GetSystems()
	sr := goscaleio.NewSystem(c)
	sr.System = getSystems[0]
	pdr := goscaleio.NewProtectionDomain(c)
	protectionDomain, err := sr.FindProtectionDomain(pdID, "", "")
	if err != nil {
		return nil, err
	}
	pdr.ProtectionDomain = protectionDomain
	spr := goscaleio.NewStoragePool(c)
	storagePool, err := pdr.FindStoragePool(spID, "", "")
	if err != nil {
		return nil, err
	}
	spr.StoragePool = storagePool
	return spr, nil
}

// Difference function to find the state difference b/w sdcs
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

// stringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type stringDefaultModifier struct {
	Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// PlanModifyString runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m stringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is unknown or known, do not set default value.
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.StringValue(m.Default)
	}
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringValue(m.Default)
	}
}

func stringDefault(defaultValue string) planmodifier.String {
	return stringDefaultModifier{
		Default: defaultValue,
	}
}

var sdcListSchema schema.SetNestedAttribute = schema.SetNestedAttribute{
	Description:         "mapped sdc info",
	Computed:            true,
	Optional:            true,
	MarkdownDescription: "mapped sdc info",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"sdc_id": schema.StringAttribute{
				Description:         "The ID of the SDC",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the SDC",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_name")),
				},
			},
			"sdc_name": schema.StringAttribute{
				Description:         "The Name of the SDC",
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The Name of the SDC",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("sdc_id")),
				},
			},
			"limit_iops": schema.Int64Attribute{
				Description:         "limit iops",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "limit iops",
			},
			"limit_bw_in_mbps": schema.Int64Attribute{
				Description:         "limit bw in mbps",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "limit bw in mbps",
			},
			"access_mode": schema.StringAttribute{
				Description:         "The Access Mode of the SDC",
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The Access Mode of the SDC",
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
