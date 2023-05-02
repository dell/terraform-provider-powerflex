package powerflex

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

type protectionDomainDataSourceModel struct {
	ProtectionDomains []protectionDomainModel `tfsdk:"protection_domains"`
	ID                types.String            `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
}

type protectionDomainModel struct {
	SystemID                    types.String    `tfsdk:"system_id"`
	SdrSdsConnectivityInfo      pdConnInfoModel `tfsdk:"sdr_sds_connectivity"`
	ReplicationCapacityMaxRatio types.Int64     `tfsdk:"replication_capacity_max_ratio"`

	// Network throttling params
	RebuildNetworkThrottlingInKbps                   types.Int64 `tfsdk:"rebuild_network_throttling_in_kbps"`
	RebalanceNetworkThrottlingInKbps                 types.Int64 `tfsdk:"rebalance_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingInKbps                 types.Int64 `tfsdk:"overall_io_network_throttling_in_kbps"`
	VTreeMigrationNetworkThrottlingInKbps            types.Int64 `tfsdk:"vtree_migration_network_throttling_in_kbps"`
	ProtectedMaintenanceModeNetworkThrottlingInKbps  types.Int64 `tfsdk:"protected_maintenance_mode_network_throttling_in_kbps"`
	OverallIoNetworkThrottlingEnabled                types.Bool  `tfsdk:"overall_io_network_throttling_enabled"`
	RebuildNetworkThrottlingEnabled                  types.Bool  `tfsdk:"rebuild_network_throttling_enabled"`
	RebalanceNetworkThrottlingEnabled                types.Bool  `tfsdk:"rebalance_network_throttling_enabled"`
	VTreeMigrationNetworkThrottlingEnabled           types.Bool  `tfsdk:"vtree_migration_network_throttling_enabled"`
	ProtectedMaintenanceModeNetworkThrottlingEnabled types.Bool  `tfsdk:"protected_maintenance_mode_network_throttling_enabled"`

	// Fine Granularity Params
	FglDefaultNumConcurrentWrites types.Int64 `tfsdk:"fgl_default_num_concurrent_writes"`
	FglMetadataCacheEnabled       types.Bool  `tfsdk:"fgl_metadata_cache_enabled"`
	FglDefaultMetadataCacheSize   types.Int64 `tfsdk:"fgl_default_metadata_cache_size"`

	// RfCache Params
	RfCacheEnabled         types.Bool   `tfsdk:"rf_cache_enabled"`
	RfCacheAccpID          types.String `tfsdk:"rf_cache_accp_id"`
	RfCacheOperationalMode types.String `tfsdk:"rf_cache_opertional_mode"`
	RfCachePageSizeKb      types.Int64  `tfsdk:"rf_cache_page_size_kb"`
	RfCacheMaxIoSizeKb     types.Int64  `tfsdk:"rf_cache_max_io_size_kb"`

	// Counter Params
	SdsConfigurationFailureCP            pdCounterModel `tfsdk:"sds_configuration_failure_counter"`
	SdsDecoupledCP                       pdCounterModel `tfsdk:"sds_decoupled_counter"`
	MdmSdsNetworkDisconnectionsCP        pdCounterModel `tfsdk:"mdm_sds_network_disconnections_counter"`
	SdsSdsNetworkDisconnectionsCP        pdCounterModel `tfsdk:"sds_sds_network_disconnections_counter"`
	SdsReceiveBufferAllocationFailuresCP pdCounterModel `tfsdk:"sds_receive_buffer_allocation_failures_counter"`

	State types.String                `tfsdk:"state"`
	Name  types.String                `tfsdk:"name"`
	ID    types.String                `tfsdk:"id"`
	Links []protectionDomainLinkModel `tfsdk:"links"`
}

type windowModel struct {
	Threshold       types.Int64 `tfsdk:"threshold"`
	WindowSizeInSec types.Int64 `tfsdk:"window_size_in_sec"`
}

type pdCounterModel struct {
	ShortWindow  windowModel `tfsdk:"short_window"`
	MediumWindow windowModel `tfsdk:"medium_window"`
	LongWindow   windowModel `tfsdk:"long_window"`
}

func pdCounterModelValue(p scaleiotypes.PDCounterParams) pdCounterModel {
	return pdCounterModel{
		ShortWindow: windowModel{
			Threshold:       types.Int64Value(int64(p.ShortWindow.Threshold)),
			WindowSizeInSec: types.Int64Value(int64(p.ShortWindow.WindowSizeInSec)),
		},
		MediumWindow: windowModel{
			Threshold:       types.Int64Value(int64(p.MediumWindow.Threshold)),
			WindowSizeInSec: types.Int64Value(int64(p.MediumWindow.WindowSizeInSec)),
		},
		LongWindow: windowModel{
			Threshold:       types.Int64Value(int64(p.LongWindow.Threshold)),
			WindowSizeInSec: types.Int64Value(int64(p.LongWindow.WindowSizeInSec)),
		},
	}
}

type pdConnInfoModel struct {
	ClientServerConnStatus types.String `tfsdk:"client_server_conn_status"`
	DisconnectedClientID   types.String `tfsdk:"disconnected_client_id"`
	DisconnectedClientName types.String `tfsdk:"disconnected_client_name"`
	DisconnectedServerID   types.String `tfsdk:"disconnected_server_id"`
	DisconnectedServerName types.String `tfsdk:"disconnected_server_name"`
	DisconnectedServerIP   types.String `tfsdk:"disconnected_server_ip"`
}

func pdConnInfoModelValue(p scaleiotypes.PDConnInfo) pdConnInfoModel {
	pdconninfo := pdConnInfoModel{
		ClientServerConnStatus: types.StringValue(p.ClientServerConnStatus),
	}
	if v := p.DisconnectedClientID; v != nil {
		pdconninfo.DisconnectedClientID = types.StringValue(*v)
	} else {
		pdconninfo.DisconnectedClientID = types.StringNull()
	}
	if v := p.DisconnectedClientName; v != nil {
		pdconninfo.DisconnectedClientName = types.StringValue(*v)
	} else {
		pdconninfo.DisconnectedClientName = types.StringNull()
	}
	if v := p.DisconnectedServerID; v != nil {
		pdconninfo.DisconnectedServerID = types.StringValue(*v)
	} else {
		pdconninfo.DisconnectedServerID = types.StringNull()
	}
	if v := p.DisconnectedServerName; v != nil {
		pdconninfo.DisconnectedServerName = types.StringValue(*v)
	} else {
		pdconninfo.DisconnectedServerName = types.StringNull()
	}
	if v := p.DisconnectedServerIP; v != nil {
		pdconninfo.DisconnectedServerIP = types.StringValue(*v)
	} else {
		pdconninfo.DisconnectedServerIP = types.StringNull()
	}
	return pdconninfo
}

type protectionDomainLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
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
	var state protectionDomainDataSourceModel
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Info(ctx, "[POWERFLEX] protectionDomainDataSourceModel"+helper.PrettyJSON((state)))

	system, err := getFirstSystem(d.client)
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

	state.ProtectionDomains = getAllProtectionDomainState(protectionDomains)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getAllProtectionDomainState(protectionDomains []*scaleiotypes.ProtectionDomain) (response []protectionDomainModel) {
	for _, protectionDomainValue := range protectionDomains {
		protectionDomainState := protectionDomainModel{
			SystemID:               types.StringValue(protectionDomainValue.SystemID),
			SdrSdsConnectivityInfo: pdConnInfoModelValue(protectionDomainValue.SdrSdsConnectivityInfo),

			// Network throttling params
			RebuildNetworkThrottlingInKbps:                   types.Int64Value(int64(protectionDomainValue.RebuildNetworkThrottlingInKbps)),
			RebalanceNetworkThrottlingInKbps:                 types.Int64Value(int64(protectionDomainValue.RebalanceNetworkThrottlingInKbps)),
			OverallIoNetworkThrottlingInKbps:                 types.Int64Value(int64(protectionDomainValue.OverallIoNetworkThrottlingInKbps)),
			VTreeMigrationNetworkThrottlingInKbps:            types.Int64Value(int64(protectionDomainValue.VTreeMigrationNetworkThrottlingInKbps)),
			ProtectedMaintenanceModeNetworkThrottlingInKbps:  types.Int64Value(int64(protectionDomainValue.ProtectedMaintenanceModeNetworkThrottlingInKbps)),
			OverallIoNetworkThrottlingEnabled:                types.BoolValue(protectionDomainValue.OverallIoNetworkThrottlingEnabled),
			RebuildNetworkThrottlingEnabled:                  types.BoolValue(protectionDomainValue.RebuildNetworkThrottlingEnabled),
			RebalanceNetworkThrottlingEnabled:                types.BoolValue(protectionDomainValue.RebalanceNetworkThrottlingEnabled),
			VTreeMigrationNetworkThrottlingEnabled:           types.BoolValue(protectionDomainValue.VTreeMigrationNetworkThrottlingEnabled),
			ProtectedMaintenanceModeNetworkThrottlingEnabled: types.BoolValue(protectionDomainValue.ProtectedMaintenanceModeNetworkThrottlingEnabled),

			// Fine Granularity Params
			FglDefaultNumConcurrentWrites: types.Int64Value(int64(protectionDomainValue.FglDefaultNumConcurrentWrites)),
			FglMetadataCacheEnabled:       types.BoolValue(protectionDomainValue.FglMetadataCacheEnabled),
			FglDefaultMetadataCacheSize:   types.Int64Value(int64(protectionDomainValue.FglDefaultMetadataCacheSize)),

			// RfCache Params
			RfCacheEnabled:         types.BoolValue(protectionDomainValue.RfCacheEnabled),
			RfCacheAccpID:          types.StringValue(protectionDomainValue.RfCacheAccpID),
			RfCacheOperationalMode: types.StringValue(string(protectionDomainValue.RfCacheOperationalMode)),
			RfCachePageSizeKb:      types.Int64Value(int64(protectionDomainValue.RfCachePageSizeKb)),
			RfCacheMaxIoSizeKb:     types.Int64Value(int64(protectionDomainValue.RfCacheMaxIoSizeKb)),

			// Counter Params
			SdsConfigurationFailureCP:            pdCounterModelValue(protectionDomainValue.SdsConfigurationFailureCP),
			SdsDecoupledCP:                       pdCounterModelValue(protectionDomainValue.SdsDecoupledCP),
			MdmSdsNetworkDisconnectionsCP:        pdCounterModelValue(protectionDomainValue.MdmSdsNetworkDisconnectionsCP),
			SdsSdsNetworkDisconnectionsCP:        pdCounterModelValue(protectionDomainValue.SdsSdsNetworkDisconnectionsCP),
			SdsReceiveBufferAllocationFailuresCP: pdCounterModelValue(protectionDomainValue.SdsReceiveBufferAllocationFailuresCP),

			State: types.StringValue(protectionDomainValue.ProtectionDomainState),
			Name:  types.StringValue(protectionDomainValue.Name),
			ID:    types.StringValue(protectionDomainValue.ID),
		}

		if v := protectionDomainValue.ReplicationCapacityMaxRatio; v != nil {
			protectionDomainState.ReplicationCapacityMaxRatio = types.Int64Value(int64(*v))
		} else {
			protectionDomainState.ReplicationCapacityMaxRatio = types.Int64Null()
		}

		for _, link := range protectionDomainValue.Links {
			protectionDomainState.Links = append(protectionDomainState.Links, protectionDomainLinkModel{
				Rel:  types.StringValue(link.Rel),
				HREF: types.StringValue(link.HREF),
			})
		}

		response = append(response, protectionDomainState)
	}

	return
}
