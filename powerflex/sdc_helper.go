/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	ID   types.String `tfsdk:"id"`
	Sdcs []sdcModel   `tfsdk:"sdcs"`
	Name types.String `tfsdk:"name"`
}

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

// getAllSdcState - function to return all sdc result from goscaleio.
func getAllSdcState(ctx context.Context, client goscaleio.Client, sdcs []scaleiotypes.Sdc) *[]sdcModel {
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

		response = append(response, sdcState)
	}

	return &response
}
