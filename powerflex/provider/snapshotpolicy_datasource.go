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
)

var (
	_ datasource.DataSource              = &snapshotPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &snapshotPolicyDataSource{}
)

// SnapshotPolicyDataSource returns the snapshot policy data source
func SnapshotPolicyDataSource() datasource.DataSource {
	return &snapshotPolicyDataSource{}
}

type snapshotPolicyDataSource struct {
	client *goscaleio.Client
}

func (d *snapshotPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_policy"
}

func (d *snapshotPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SnapshotPolicyDataSourceSchema
}

func (d *snapshotPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

func (d *snapshotPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan models.SnapshotPolicyDataSourceModel
	var sps []*scaleiotypes.SnapshotPolicy
	var err error

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Read the snapshot policies based on snapshot policy id/name and if nothing
	//is mentioned , then returns all the snapshot policies
	if plan.Name.ValueString() != "" {
		sps, err = d.client.GetSnapshotPolicy(plan.Name.ValueString(), "")
	} else if plan.ID.ValueString() != "" {
		sps, err = d.client.GetSnapshotPolicy("", plan.ID.ValueString())
	} else {
		sps, err = d.client.GetSnapshotPolicy("", "")
	}
	//check if there is any error while getting the snapshot policy
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Snapshot Policy",
			err.Error(),
		)
		return
	}
	plan.SnapshotPolicies = helper.UpdateSnapshotPolicyState(sps)
	plan.ID = types.StringValue("dummyID")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
