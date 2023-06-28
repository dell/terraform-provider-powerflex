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

package provider

import (
	"context"

	"terraform-provider-powerflex/powerflex/helper"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (d *protectionDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

func (d *protectionDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ProtectionDomainDataSourceModel
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Info(ctx, "[POWERFLEX] protectionDomainDataSourceModel"+helper.PrettyJSON((state)))

	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}

	var protectionDomains []*scaleiotypes.ProtectionDomain

	if !state.ID.IsNull() {
		// Fetch protection domain of given id
		var protectionDomain *scaleiotypes.ProtectionDomain
		protectionDomain, err = system.FindProtectionDomain(state.ID.ValueString(), "", "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex ProtectionDomain by ID",
				err.Error(),
			)
			return
		}
		protectionDomains = append(protectionDomains, protectionDomain)
	} else if !state.Name.IsNull() {
		// Fetch protection domain of given name
		var protectionDomain *scaleiotypes.ProtectionDomain
		protectionDomain, err = system.FindProtectionDomain("", state.Name.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex ProtectionDomain by name",
				err.Error(),
			)
			return
		}
		protectionDomains = append(protectionDomains, protectionDomain)
		// this is required for acceptance testing
		state.ID = types.StringValue(protectionDomain.ID)
	} else {
		// Fetch all protection domains
		protectionDomains, err = system.GetProtectionDomain("")
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex ProtectionDomains",
				err.Error(),
			)
			return
		}
		// this is required for acceptance testing
		state.ID = types.StringValue("DummyID")
	}

	state.ProtectionDomains = helper.GetAllProtectionDomainState(protectionDomains)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
