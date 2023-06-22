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

package helper

import (
	"context"

	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetPDResState is a helper function that marshalls API response into protectionDomainResourceModel
func GetPDResState(protectionDomain *scaleiotypes.ProtectionDomain) models.ProtectionDomainResourceModel {
	return models.ProtectionDomainResourceModel{
		Name:   types.StringValue(protectionDomain.Name),
		ID:     types.StringValue(protectionDomain.ID),
		Active: types.BoolValue(protectionDomain.ProtectionDomainState == "Active"),
		State:  types.StringValue(protectionDomain.ProtectionDomainState),

		// Network throttling params
		RebuildNetworkThrottlingInKbps:                  types.Int64Value(int64(protectionDomain.RebuildNetworkThrottlingInKbps)),
		RebalanceNetworkThrottlingInKbps:                types.Int64Value(int64(protectionDomain.RebalanceNetworkThrottlingInKbps)),
		OverallIoNetworkThrottlingInKbps:                types.Int64Value(int64(protectionDomain.OverallIoNetworkThrottlingInKbps)),
		VTreeMigrationNetworkThrottlingInKbps:           types.Int64Value(int64(protectionDomain.VTreeMigrationNetworkThrottlingInKbps)),
		ProtectedMaintenanceModeNetworkThrottlingInKbps: types.Int64Value(int64(protectionDomain.ProtectedMaintenanceModeNetworkThrottlingInKbps)),

		// Fine Granularity Params
		FglDefaultNumConcurrentWrites: types.Int64Value(int64(protectionDomain.FglDefaultNumConcurrentWrites)),
		FglMetadataCacheEnabled:       types.BoolValue(protectionDomain.FglMetadataCacheEnabled),
		FglDefaultMetadataCacheSize:   types.Int64Value(int64(protectionDomain.FglDefaultMetadataCacheSize)),

		// RfCache Params
		RfCacheEnabled:         types.BoolValue(protectionDomain.RfCacheEnabled),
		RfCacheAccpID:          types.StringValue(protectionDomain.RfCacheAccpID),
		RfCacheOperationalMode: types.StringValue(string(protectionDomain.RfCacheOperationalMode)),
		RfCachePageSizeKb:      types.Int64Value(int64(protectionDomain.RfCachePageSizeKb)),
		RfCacheMaxIoSizeKb:     types.Int64Value(int64(protectionDomain.RfCacheMaxIoSizeKb)),

		// Links
		Links: GetLinkTfList(protectionDomain.Links),
	}
}

// GetLinkTfList is a helper function that marshalls goscaleio links into types.List
func GetLinkTfList(links []*scaleiotypes.Link) types.List {
	sourceKeywordAttrTypes := map[string]attr.Type{
		"rel":  types.StringType,
		"href": types.StringType,
	}
	elemType := types.ObjectType{AttrTypes: sourceKeywordAttrTypes}
	objLinksList := []attr.Value{}

	for _, link := range links {
		obj := map[string]attr.Value{
			"rel":  types.StringValue(link.Rel),
			"href": types.StringValue(link.HREF),
		}
		objVal, _ := types.ObjectValue(sourceKeywordAttrTypes, obj)
		objLinksList = append(objLinksList, objVal)
	}
	listVal, _ := types.ListValue(elemType, objLinksList)
	return listVal
}

// GetLinksFromTfList is a helper function that unmarshalls goscaleio links from types.List
func GetLinksFromTfList(ctx context.Context, links types.List) ([]*scaleiotypes.Link, diag.Diagnostics) {
	var d diag.Diagnostics
	listVal := make([]*scaleiotypes.Link, 0)
	if links.IsNull() || links.IsUnknown() {
		return listVal, d
	}
	type source struct {
		Rel  types.String `tfsdk:"rel"`
		Href types.String `tfsdk:"href"`
	}
	sourceAttrTypes := []source{}
	diags := links.ElementsAs(ctx, &sourceAttrTypes, true)
	d.Append(diags...)

	for _, item := range sourceAttrTypes {
		listVal = append(listVal, &scaleiotypes.Link{
			Rel:  item.Rel.ValueString(),
			HREF: item.Href.ValueString(),
		})
	}
	return listVal, d
}
