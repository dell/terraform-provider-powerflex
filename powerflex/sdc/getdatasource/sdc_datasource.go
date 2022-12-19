package getdatasource

import (
	"context"

	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sdcDataSource{}
	_ datasource.DataSourceWithConfigure = &sdcDataSource{}
)

// SDCDataSource - function used to return SDC DataSource provider with singleton values.
func SDCDataSource() datasource.DataSource {
	return &sdcDataSource{}
}

// sdcFilterType - Enum structure for filter types.
var sdcFilterType = struct {
	All    string
	ByName string
	ByID   string
}{
	All:    "All",
	ByName: "ByName",
	ByID:   "ByID",
}

// sdcDataSource - for returning singleton holder with goscaleio client.
type sdcDataSource struct {
	client *goscaleio.Client
}

// sdcDataSourceModel - for returning result to terraform.
type sdcDataSourceModel struct {
	Sdcs     []sdcModel   `tfsdk:"sdcs"`
	ID       types.String `tfsdk:"sdcid"`
	SystemID types.String `tfsdk:"systemid"`
	Name     types.String `tfsdk:"name"`
}

type SdcStatistics struct {
	NumOfMappedVolumes      types.Int64   `tfsdk:"numofmappedvolumes"`
	VolumeIds               VolumeIdsList `tfsdk:"volumeids"`
	UserDataReadBwc         SdcBwc        `tfsdk:"userdatareadbwc"`
	UserDataWriteBwc        SdcBwc        `tfsdk:"userdatawritebwc"`
	UserDataTrimBwc         SdcBwc        `tfsdk:"userdatatrimbwc"`
	UserDataSdcReadLatency  SdcBwc        `tfsdk:"userdatasdcreadlatency"`
	UserDataSdcWriteLatency SdcBwc        `tfsdk:"userdatasdcwritelatency"`
	UserDataSdcTrimLatency  SdcBwc        `tfsdk:"userdatasdctrimlatency"`
}
type SdcBwc struct {
	TotalWeightInKb types.Int64 `tfsdk:"totalweightinkb"`
	NumOccured      types.Int64 `tfsdk:"numoccured"`
	NumSeconds      types.Int64 `tfsdk:"numseconds"`
}

type VolumeIdsList []types.String

// sdcModel - MODEL for SDC data returned by goscaleio.
type sdcModel struct {
	ID                 types.String   `tfsdk:"id"`
	SystemID           types.String   `tfsdk:"systemid"`
	SdcIP              types.String   `tfsdk:"sdcip"`
	SdcApproved        types.Bool     `tfsdk:"sdcapproved"`
	OnVMWare           types.Bool     `tfsdk:"onvmware"`
	SdcGUID            types.String   `tfsdk:"sdcguid"`
	MdmConnectionState types.String   `tfsdk:"mdmconnectionstate"`
	Name               types.String   `tfsdk:"name"`
	Links              []sdcLinkModel `tfsdk:"links"`
	Statistics         SdcStatistics  `tfsdk:"statistics"`
}

// sdcLinkModel - MODEL for SDC Links data returned by goscaleio.
type sdcLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// Metadata - function used to define datasource metadata[referance in tf file].
func (d *sdcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

// GetSchema - function used to return SDC datasource schema.
func (d *sdcDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SDCDataSourceScheme
}

// Configure - function to call initial configurations before resource execution.
func (d *sdcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

// Read - function to read sdc values from goscaleio.
func (d *sdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sdcDataSourceModel
	diags := req.Config.Get(ctx, &state)
	tflog.Info(ctx, "[POWERFLEX] sdcDataSourceModel"+helper.PrettyJSON((state)))
	system, err := d.client.FindSystem(state.SystemID.ValueString(), "", "")
	// tflog.Debug(ctx, "system"+helper.PrettyJSON(system))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex systems sdcs",
			err.Error(),
		)
		return
	}

	sdcs, err := system.GetSdc()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex sdcs",
			err.Error(),
		)
		return
	}
	// Set state
	searchFilter := sdcFilterType.All
	if state.Name.ValueString() != "" {
		searchFilter = sdcFilterType.ByName
	}
	if state.ID.ValueString() != "" {
		searchFilter = sdcFilterType.ByID
	}
	if state.Name.ValueString() != "" && state.ID.ValueString() != "" {
		searchFilter = sdcFilterType.ByID
	}

	if searchFilter == sdcFilterType.All {
		state.Sdcs = getAllSdcState(sdcs)
	} else {
		filterResult, err := getFilteredSdcState(*d.client, ctx, sdcs, searchFilter, state.Name.ValueString(), state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Statics for sdc id = "+state.ID.ValueString()+", name = "+state.Name.ValueString(),
				err.Error(),
			)
			return
		}
		state.Sdcs = *filterResult
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// getFilteredSdcState - function to filter sdc result from goscaleio.
func getFilteredSdcState(client goscaleio.Client, ctx context.Context, sdcs []scaleiotypes.Sdc, method string, name string, id string) (*[]sdcModel, error) {
	tflog.Debug(ctx, "[POWERFLEX] searchFilter getFilteredSdcState method "+method+" name "+name+" id "+id)
	response := []sdcModel{}
	for _, sdcValue := range sdcs {
		sdcState := sdcModel{
			ID:                 types.StringValue(sdcValue.ID),
			Name:               types.StringValue(sdcValue.Name),
			SdcGUID:            types.StringValue(sdcValue.SdcGUID),
			SdcApproved:        types.BoolValue(sdcValue.SdcApproved),
			OnVMWare:           types.BoolValue(sdcValue.OnVMWare),
			SystemID:           types.StringValue(sdcValue.SystemID),
			SdcIP:              types.StringValue(sdcValue.SdcIP),
			MdmConnectionState: types.StringValue(sdcValue.MdmConnectionState),
		}

		for _, link := range sdcValue.Links {
			sdcState.Links = append(sdcState.Links, sdcLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}

		sdcc := goscaleio.NewSdc(&client, &sdcValue)
		stats, err := sdcc.GetStatistics()
		if err != nil {
			tflog.Debug(ctx, "[POWERFLEX] err in GetStatistics "+helper.PrettyJSON(err))
			return nil, err // Sometimes unable to find link is error we get in 4.0
		}

		tflog.Debug(ctx, "[POWERFLEX] stats in GetStatistics "+helper.PrettyJSON(stats))

		sdcState.Statistics = createStaticsObject(*stats)

		if method == sdcFilterType.ByName && name == sdcValue.Name {

			response = append(response, sdcState)
		}
		if method == sdcFilterType.ByID && id == sdcValue.ID {
			response = append(response, sdcState)
		}

	}

	return &response, nil
}

func createStaticsObject(stats scaleiotypes.SdcStatistics) (ret SdcStatistics) {
	VolumeIdsAll := VolumeIdsList{}
	for _, v := range stats.VolumeIds {
		VolumeIdsAll = append(VolumeIdsAll, types.StringValue(v))
	}
	ret.NumOfMappedVolumes = types.Int64Value(int64(stats.NumOfMappedVolumes))
	ret.VolumeIds = VolumeIdsAll

	ret.UserDataReadBwc = SdcBwc{
		TotalWeightInKb: types.Int64Value(int64(stats.UserDataReadBwc.TotalWeightInKb)),
		NumOccured:      types.Int64Value(int64(stats.UserDataReadBwc.NumOccured)),
		NumSeconds:      types.Int64Value(int64(stats.UserDataReadBwc.NumSeconds)),
	}
	ret.UserDataWriteBwc = SdcBwc{
		TotalWeightInKb: types.Int64Value(int64(stats.UserDataWriteBwc.TotalWeightInKb)),
		NumOccured:      types.Int64Value(int64(stats.UserDataWriteBwc.NumOccured)),
		NumSeconds:      types.Int64Value(int64(stats.UserDataWriteBwc.NumSeconds)),
	}
	ret.UserDataTrimBwc = SdcBwc{
		TotalWeightInKb: types.Int64Value(int64(stats.UserDataTrimBwc.TotalWeightInKb)),
		NumOccured:      types.Int64Value(int64(stats.UserDataTrimBwc.NumOccured)),
		NumSeconds:      types.Int64Value(int64(stats.UserDataTrimBwc.NumSeconds)),
	}
	ret.UserDataSdcReadLatency = SdcBwc{
		TotalWeightInKb: types.Int64Value(int64(stats.UserDataSdcReadLatency.TotalWeightInKb)),
		NumOccured:      types.Int64Value(int64(stats.UserDataSdcReadLatency.NumOccured)),
		NumSeconds:      types.Int64Value(int64(stats.UserDataSdcReadLatency.NumSeconds)),
	}
	ret.UserDataSdcWriteLatency = SdcBwc{
		TotalWeightInKb: types.Int64Value(int64(stats.UserDataSdcWriteLatency.TotalWeightInKb)),
		NumOccured:      types.Int64Value(int64(stats.UserDataSdcWriteLatency.NumOccured)),
		NumSeconds:      types.Int64Value(int64(stats.UserDataSdcWriteLatency.NumSeconds)),
	}
	ret.UserDataSdcTrimLatency = SdcBwc{
		TotalWeightInKb: types.Int64Value(int64(stats.UserDataSdcTrimLatency.TotalWeightInKb)),
		NumOccured:      types.Int64Value(int64(stats.UserDataSdcTrimLatency.NumOccured)),
		NumSeconds:      types.Int64Value(int64(stats.UserDataSdcTrimLatency.NumSeconds)),
	}

	return
}

// getAllSdcState - function to return all sdc result from goscaleio.
func getAllSdcState(sdcs []scaleiotypes.Sdc) (response []sdcModel) {
	for _, sdcValue := range sdcs {
		sdcState := sdcModel{
			ID:                 types.StringValue(sdcValue.ID),
			Name:               types.StringValue(sdcValue.Name),
			SdcGUID:            types.StringValue(sdcValue.SdcGUID),
			SdcApproved:        types.BoolValue(sdcValue.SdcApproved),
			OnVMWare:           types.BoolValue(sdcValue.OnVMWare),
			SystemID:           types.StringValue(sdcValue.SystemID),
			SdcIP:              types.StringValue(sdcValue.SdcIP),
			MdmConnectionState: types.StringValue(sdcValue.MdmConnectionState),
		}

		for _, link := range sdcValue.Links {
			sdcState.Links = append(sdcState.Links, sdcLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		response = append(response, sdcState)
	}

	return
}
