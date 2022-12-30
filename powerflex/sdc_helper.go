package powerflex

import (
	"context"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	ID    types.String `tfsdk:"id"`
	Sdcs  []sdcModel   `tfsdk:"sdcs"`
	SdcID types.String `tfsdk:"sdc_id"`
	// SystemID types.String `tfsdk:"systemid"`
	Name types.String `tfsdk:"name"`
}

// sdcStatistics - MODEL for SDC statistics.
type sdcStatistics struct {
	NumOfMappedVolumes      types.Int64   `tfsdk:"numofmappedvolumes"`
	VolumeIds               VolumeIdsList `tfsdk:"volumeids"`
	UserDataReadBwc         sdcBwc        `tfsdk:"userdatareadbwc"`
	UserDataWriteBwc        sdcBwc        `tfsdk:"userdatawritebwc"`
	UserDataTrimBwc         sdcBwc        `tfsdk:"userdatatrimbwc"`
	UserDataSdcReadLatency  sdcBwc        `tfsdk:"userdatasdcreadlatency"`
	UserDataSdcWriteLatency sdcBwc        `tfsdk:"userdatasdcwritelatency"`
	UserDataSdcTrimLatency  sdcBwc        `tfsdk:"userdatasdctrimlatency"`
}

// sdcBwc - MODEL for SDC statistics BWC.
type sdcBwc struct {
	TotalWeightInKb types.Int64 `tfsdk:"totalweightinkb"`
	NumOccured      types.Int64 `tfsdk:"numoccured"`
	NumSeconds      types.Int64 `tfsdk:"numseconds"`
}

// VolumeIdsList - MODEL for SDC statistics Volume Id List.
type VolumeIdsList []types.String

// sdcModel - MODEL for SDC data returned by goscaleio.
type sdcModel struct {
	ID                 types.String   `tfsdk:"id"`
	SystemID           types.String   `tfsdk:"system_id"`
	SdcIP              types.String   `tfsdk:"sdc_ip"`
	SdcApproved        types.Bool     `tfsdk:"sdc_approved"`
	OnVMWare           types.Bool     `tfsdk:"on_vmware"`
	SdcGUID            types.String   `tfsdk:"sdc_guid"`
	MdmConnectionState types.String   `tfsdk:"mdm_connection_state"`
	Name               types.String   `tfsdk:"name"`
	Links              []sdcLinkModel `tfsdk:"links"`
	// Statistics         sdcStatistics  `tfsdk:"statistics"`
}

// sdcLinkModel - MODEL for SDC Links data returned by goscaleio.
type sdcLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// getFilteredSdcState - function to filter sdc result from goscaleio.
func getFilteredSdcState(sdcs *[]sdcModel, method string, name string, id string) *[]sdcModel {
	response := []sdcModel{}
	for _, sdcValue := range *sdcs {
		if method == sdcFilterType.ByName && name == sdcValue.Name.ValueString() {
			response = append(response, sdcValue)
		}
		if method == sdcFilterType.ByID && id == sdcValue.ID.ValueString() {
			response = append(response, sdcValue)
		}
	}
	return &response
}

// func createStaticsObject(stats scaleiotypes.SdcStatistics) (ret sdcStatistics) {
// 	VolumeIdsAll := VolumeIdsList{}
// 	for _, v := range stats.VolumeIds {
// 		VolumeIdsAll = append(VolumeIdsAll, types.StringValue(v))
// 	}
// 	ret.NumOfMappedVolumes = types.Int64Value(int64(stats.NumOfMappedVolumes))
// 	ret.VolumeIds = VolumeIdsAll

// 	ret.UserDataReadBwc = sdcBwc{
// 		TotalWeightInKb: types.Int64Value(int64(stats.UserDataReadBwc.TotalWeightInKb)),
// 		NumOccured:      types.Int64Value(int64(stats.UserDataReadBwc.NumOccured)),
// 		NumSeconds:      types.Int64Value(int64(stats.UserDataReadBwc.NumSeconds)),
// 	}
// 	ret.UserDataWriteBwc = sdcBwc{
// 		TotalWeightInKb: types.Int64Value(int64(stats.UserDataWriteBwc.TotalWeightInKb)),
// 		NumOccured:      types.Int64Value(int64(stats.UserDataWriteBwc.NumOccured)),
// 		NumSeconds:      types.Int64Value(int64(stats.UserDataWriteBwc.NumSeconds)),
// 	}
// 	ret.UserDataTrimBwc = sdcBwc{
// 		TotalWeightInKb: types.Int64Value(int64(stats.UserDataTrimBwc.TotalWeightInKb)),
// 		NumOccured:      types.Int64Value(int64(stats.UserDataTrimBwc.NumOccured)),
// 		NumSeconds:      types.Int64Value(int64(stats.UserDataTrimBwc.NumSeconds)),
// 	}
// 	ret.UserDataSdcReadLatency = sdcBwc{
// 		TotalWeightInKb: types.Int64Value(int64(stats.UserDataSdcReadLatency.TotalWeightInKb)),
// 		NumOccured:      types.Int64Value(int64(stats.UserDataSdcReadLatency.NumOccured)),
// 		NumSeconds:      types.Int64Value(int64(stats.UserDataSdcReadLatency.NumSeconds)),
// 	}
// 	ret.UserDataSdcWriteLatency = sdcBwc{
// 		TotalWeightInKb: types.Int64Value(int64(stats.UserDataSdcWriteLatency.TotalWeightInKb)),
// 		NumOccured:      types.Int64Value(int64(stats.UserDataSdcWriteLatency.NumOccured)),
// 		NumSeconds:      types.Int64Value(int64(stats.UserDataSdcWriteLatency.NumSeconds)),
// 	}
// 	ret.UserDataSdcTrimLatency = sdcBwc{
// 		TotalWeightInKb: types.Int64Value(int64(stats.UserDataSdcTrimLatency.TotalWeightInKb)),
// 		NumOccured:      types.Int64Value(int64(stats.UserDataSdcTrimLatency.NumOccured)),
// 		NumSeconds:      types.Int64Value(int64(stats.UserDataSdcTrimLatency.NumSeconds)),
// 	}

// 	return
// }

// getAllSdcState - function to return all sdc result from goscaleio.
func getAllSdcState(ctx context.Context, client goscaleio.Client, sdcs []scaleiotypes.Sdc) (*[]sdcModel, error) {
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

		// sdcc := goscaleio.NewSdc(&client, &sdcValue)
		// stats, err := sdcc.GetStatistics()
		// if err != nil {
		// 	tflog.Debug(ctx, "[POWERFLEX] err in GetStatistics "+helper.PrettyJSON(err))
		// 	return nil, err // Sometimes unable to find link is error we get in 4.0
		// }

		// tflog.Debug(ctx, "[POWERFLEX] stats in GetStatistics "+helper.PrettyJSON(stats))

		// sdcState.Statistics = createStaticsObject(*stats)

		response = append(response, sdcState)
	}

	return &response, nil
}
