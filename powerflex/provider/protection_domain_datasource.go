/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package provider

import (
	"context"
	"fmt"

	"terraform-provider-powerflex/powerflex/helper"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &protectionDomainDataSource{}
	_ datasource.DataSourceWithConfigure = &protectionDomainDataSource{}
)

// ProtectionDomainDataSource returns the datasource for protection domain
func ProtectionDomainDataSource() datasource.DataSource {
	return &protectionDomainDataSource{}
}

type protectionDomainDataSource struct {
	client *goscaleio.Client
}

func (d *protectionDomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protection_domain"
}

func (d *protectionDomainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProtectionDomainDataSourceSchema
}

func (d *protectionDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
}

func (d *protectionDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ProtectionDomainDataSourceModel
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	protectionDomains, err := helper.GetProtectionDomains(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex ProtectionDomains",
			err.Error(),
		)
		return
	}

	if state.ProtectionDomainFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.ProtectionDomainFilter, protectionDomains)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering protection domains: %v please validate the filter", state.ProtectionDomainFilter), err.Error(),
			)
			return
		}
		filteredPds := []scaleiotypes.ProtectionDomain{}
		for _, val := range filtered {
			filteredPds = append(filteredPds, val.(scaleiotypes.ProtectionDomain))
		}
		protectionDomains = filteredPds
	}

	state.ID = types.StringValue("protection-domain-datasoure-id")
	state.ProtectionDomains = helper.GetAllProtectionDomainState(protectionDomains)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
