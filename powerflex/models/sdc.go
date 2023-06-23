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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SdcDataSourceModel - for returning result to terraform.
type SdcDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Sdcs []SdcModel   `tfsdk:"sdcs"`
	Name types.String `tfsdk:"name"`
}

// SdcModel - MODEL for SDC data returned by goscaleio.
type SdcModel struct {
	ID                 types.String   `tfsdk:"id"`
	SystemID           types.String   `tfsdk:"system_id"`
	SdcIP              types.String   `tfsdk:"sdc_ip"`
	SdcApproved        types.Bool     `tfsdk:"sdc_approved"`
	OnVMWare           types.Bool     `tfsdk:"on_vmware"`
	SdcGUID            types.String   `tfsdk:"sdc_guid"`
	MdmConnectionState types.String   `tfsdk:"mdm_connection_state"`
	Name               types.String   `tfsdk:"name"`
	Links              []SdcLinkModel `tfsdk:"links"`
}

// SdcLinkModel - MODEL for SDC Links data returned by goscaleio.
type SdcLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

// SdcResourceSchemaDescriptions defines struct for SDC resource schema description
var SdcResourceSchemaDescriptions = struct {
	SdcResourceSchema  string
	LastUpdated        string
	SystemID           string
	SdcIP              string
	SdcApproved        string
	OnVMWare           string
	SdcGUID            string
	MdmConnectionState string
	Links              string
	LinksRel           string
	LinksHref          string
}{
	SdcResourceSchema:  "This resource can be used to manage Storage Data Clients on a PowerFlex array.",
	LastUpdated:        "The Last updated timestamp of the SDC.",
	SystemID:           "The System ID of the fetched SDC.",
	SdcIP:              "The IP of the fetched SDC.",
	SdcApproved:        "If the fetched SDC is approved.",
	OnVMWare:           "If the fetched SDC is on vmware.",
	SdcGUID:            "The GUID of the fetched SDC.",
	MdmConnectionState: "The MDM connection status of the fetched SDC.",
	Links:              "The Links of the fetched SDC.",
	LinksRel:           "The Links-Rel of the fetched SDC.",
	LinksHref:          "The Links-HREF of the fetched SDC.",
}

// SdcResourceModel struct for CSV Data Processing
type SdcResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	SDCDetails  types.List   `tfsdk:"sdc_details"`
	MdmPassword types.String `tfsdk:"mdm_password"`
	LiaPassword types.String `tfsdk:"lia_password"`
}

// SDCDetailDataModel defines the struct for CSV Parse Data
type SDCDetailDataModel struct {
	SDCID              types.String `tfsdk:"sdc_id"`
	IP                 types.String `tfsdk:"ip"`
	UserName           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	OperatingSystem    types.String `tfsdk:"operating_system"`
	IsMdmOrTb          types.String `tfsdk:"is_mdm_or_tb"`
	IsSdc              types.String `tfsdk:"is_sdc"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	SDCName            types.String `tfsdk:"name"`
	SystemID           types.String `tfsdk:"system_id"`
	SdcApproved        types.Bool   `tfsdk:"sdc_approved"`
	OnVMWare           types.Bool   `tfsdk:"on_vmware"`
	SdcGUID            types.String `tfsdk:"sdc_guid"`
	MdmConnectionState types.String `tfsdk:"mdm_connection_state"`
}

// CsvRow desfines the srtuct for the CSV Data
type CsvRow struct {
	IP                 string
	UserName           string
	Password           string
	OperatingSystem    string
	IsMdmOrTb          string
	IsSdc              string
	PerformanceProfile string
}

// SdcDatasourceSchemaDescriptions defines struct for SDC datasource schema description
var SdcDatasourceSchemaDescriptions = struct {
	LastUpdated        string
	ID                 string
	SystemID           string
	Name               string
	SdcIP              string
	SdcApproved        string
	OnVMWare           string
	SdcGUID            string
	MdmConnectionState string
	Links              string
	LinksRel           string
	LinksHref          string
}{
	LastUpdated:        "The Last updated timestamp of the fetched SDC.",
	ID:                 "The ID of the fetched SDC.",
	SystemID:           "The System ID of the fetched SDC.",
	Name:               "The name of the fetched SDC.",
	SdcIP:              "The IP of the fetched SDC.",
	SdcApproved:        "If the fetched SDC is approved.",
	OnVMWare:           "If the fetched SDC is on vmware.",
	SdcGUID:            "The GUID of the fetched SDC.",
	MdmConnectionState: "The MDM connection status of the fetched SDC.",
	Links:              "The Links of the fetched SDC.",
	LinksRel:           "The Links-Rel of the fetched SDC.",
	LinksHref:          "The Links-HREF of the fetched SDC.",
}
