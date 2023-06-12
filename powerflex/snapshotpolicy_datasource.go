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

type snapshotPolicyDataSourceModel struct {
	SnapshotPolicies []snapshotPolicyModel `tfsdk:"snapshotpolicies"`
	ID               types.String          `tfsdk:"id"`
	Name             types.String          `tfsdk:"name"`
}

type snapshotPolicyModel struct {
	ID                                    types.String              `tfsdk:"id"`
	Name                                  types.String              `tfsdk:"name"`
	SnapshotPolicyState                   types.String              `tfsdk:"snapshot_policy_state"`
	AutoSnapshotCreationCadenceInMin      types.Int64               `tfsdk:"auto_snapshot_creation_cadence_in_min"`
	MaxVTreeAutoSnapshots                 types.Int64               `tfsdk:"max_vtree_auto_snapshots"`
	NumOfSourceVolumes                    types.Int64               `tfsdk:"num_of_source_volumes"`
	NumOfExpiredButLockedSnapshots        types.Int64               `tfsdk:"num_of_expired_but_locked_snapshots"`
	NumOfCreationFailures                 types.Int64               `tfsdk:"num_of_creation_failures"`
	NumOfRetainedSnapshotsPerLevel        []types.Int64             `tfsdk:"num_of_retained_snapshots_per_level"`
	SnapshotAccessMode                    types.String              `tfsdk:"snapshot_access_mode"`
	SecureSnapshots                       types.Bool                `tfsdk:"secure_snapshots"`
	TimeOfLastAutoSnapshot                types.Int64               `tfsdk:"time_of_last_auto_snapshot"`
	NextAutoSnapshotCreationTime          types.Int64               `tfsdk:"next_auto_snapshot_creation_time"`
	TimeOfLastAutoSnapshotCreationFailure types.Int64               `tfsdk:"time_of_last_auto_snapshot_creation_failure"`
	LastAutoSnapshotCreationFailureReason types.String              `tfsdk:"last_auto_snapshot_creation_failure_reason"`
	LastAutoSnapshotFailureInFirstLevel   types.Bool                `tfsdk:"last_auto_snapshot_failure_in_first_level"`
	NumOfAutoSnapshots                    types.Int64               `tfsdk:"num_of_auto_snapshots"`
	NumOfLockedSnapshots                  types.Int64               `tfsdk:"num_of_locked_snapshots"`
	SystemID                              types.String              `tfsdk:"system_id"`
	Links                                 []snapshotPolicyLinkModel `tfsdk:"links"`
}

type snapshotPolicyLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
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
	var plan snapshotPolicyDataSourceModel
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
	plan.SnapshotPolicies = updateSnapshotPolicyState(sps)
	plan.ID = types.StringValue("dummyID")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// updateSnapshotPolicyState iterates over the snapshotpolicy list and update the state
func updateSnapshotPolicyState(sps []*scaleiotypes.SnapshotPolicy) (response []snapshotPolicyModel) {
	for _, sp := range sps {
		spState := snapshotPolicyModel{
			ID:                                    types.StringValue(sp.ID),
			Name:                                  types.StringValue(sp.Name),
			SnapshotPolicyState:                   types.StringValue(sp.SnapshotPolicyState),
			AutoSnapshotCreationCadenceInMin:      types.Int64Value((int64)(sp.AutoSnapshotCreationCadenceInMin)),
			MaxVTreeAutoSnapshots:                 types.Int64Value((int64)(sp.MaxVTreeAutoSnapshots)),
			NumOfSourceVolumes:                    types.Int64Value((int64)(sp.NumOfSourceVolumes)),
			NumOfExpiredButLockedSnapshots:        types.Int64Value((int64)(sp.NumOfExpiredButLockedSnapshots)),
			NumOfCreationFailures:                 types.Int64Value((int64)(sp.NumOfCreationFailures)),
			SnapshotAccessMode:                    types.StringValue(sp.SnapshotAccessMode),
			SecureSnapshots:                       types.BoolValue(sp.SecureSnapshots),
			TimeOfLastAutoSnapshot:                types.Int64Value((int64)(sp.TimeOfLastAutoSnapshot)),
			NextAutoSnapshotCreationTime:          types.Int64Value((int64)(sp.NextAutoSnapshotCreationTime)),
			TimeOfLastAutoSnapshotCreationFailure: types.Int64Value((int64)(sp.TimeOfLastAutoSnapshotCreationFailure)),
			LastAutoSnapshotCreationFailureReason: types.StringValue(sp.LastAutoSnapshotCreationFailureReason),
			LastAutoSnapshotFailureInFirstLevel:   types.BoolValue(sp.LastAutoSnapshotFailureInFirstLevel),
			NumOfAutoSnapshots:                    types.Int64Value((int64)(sp.NumOfAutoSnapshots)),
			NumOfLockedSnapshots:                  types.Int64Value((int64)(sp.NumOfLockedSnapshots)),
			SystemID:                              types.StringValue(sp.SystemID),
		}

		for _, link := range sp.Links {
			spState.Links = append(spState.Links, snapshotPolicyLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}
		for _, rspl := range sp.NumOfRetainedSnapshotsPerLevel {
			spState.NumOfRetainedSnapshotsPerLevel = append(spState.NumOfRetainedSnapshotsPerLevel, types.Int64Value((int64)(rspl)))
		}
		response = append(response, spState)
	}
	return
}
