package powerflex

import (
	"context"
	"fmt"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	SdcID         types.String `tfsdk:"sdc_id"`
	LimitIops     types.Int64  `tfsdk:"limit_iops"`
	LimitBwInMbps types.Int64  `tfsdk:"limit_bw_in_mbps"`
	SdcName       types.String `tfsdk:"sdc_name"`
	AccessMode    types.String `tfsdk:"access_mode"`
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

// DifferenceMap function to find the state difference b/w sdcs
func DifferenceMap(a, b map[string]string) map[string]string {
	m := make(map[string]bool)
	diff := make(map[string]string)

	for item := range b {
		m[item] = true
	}

	for item := range a {
		if _, ok := m[item]; !ok {
			diff[item] = item
		}
	}
	return diff
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
