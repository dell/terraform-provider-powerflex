package helper

import (
	"context"
	"fmt"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ client.Logger = &provisionerLogger{}
)

type provisionerLogger struct {
	ctx context.Context
}

func (l *provisionerLogger) Printf(format string, v ...any) {
	tflog.Info(l.ctx, fmt.Sprintf(format, v...))
}

func (l *provisionerLogger) Println(v ...any) {
	tflog.Info(l.ctx, fmt.Sprint(v...))
}

// SdcHostResource - helper for SDC host resource
type SdcHostResource struct {
	System *goscaleio.System
}

func (r *SdcHostResource) getSshProvisioner(ctx context.Context, plan models.SdcHostModel) (*client.SshProvisioner, string, error) {
	var remote models.SdcHostRemoteModel
	plan.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	dir := ""
	if remote.Dir == nil {
		dir = "/tmp"
	} else {
		dir = *remote.Dir
	}
	prov, err := client.NewSshProvisioner(client.SshProvisionerConfig{
		IP:         plan.Host.ValueString(),
		Username:   remote.User,
		Password:   remote.Password,
		PrivateKey: remote.PrivateKey,
	}, &provisionerLogger{ctx: ctx})
	return prov, dir, err
}

func (r *SdcHostResource) GetMdmIps(ctx context.Context, plan models.SdcHostModel) []string {
	var mdmIps []string
	plan.MdmIPs.ElementsAs(ctx, &mdmIps, true)
	return mdmIps
}

func (r *SdcHostResource) ReadSDCHost(ctx context.Context, state models.SdcHostModel) (models.SdcHostModel, error) {
	// get SDC by IP
	tflog.Info(ctx, "Finding SDC by IP")
	sdcData, err := r.System.FindSdc("SdcIP", state.Host.ValueString())
	if err != nil {
		return state, fmt.Errorf("error finding SDC by IP %s: %w", state.Host.ValueString(), err)
	}
	tflog.Info(ctx, "Found SDC by IP")
	state.ID = types.StringValue(sdcData.Sdc.ID)
	state.PerformanceProfile = types.StringValue(sdcData.Sdc.PerfProfile)
	state.Name = types.StringValue(sdcData.Sdc.Name)
	return state, nil
}

// SetSDCParams - function to set SDC parameters
func (r *SdcHostResource) SetSDCParams(ctx context.Context, plan, state models.SdcHostModel) error {
	// set name
	if plan.Name.ValueString() != state.Name.ValueString() && !plan.Name.IsUnknown() {
		tflog.Info(ctx, "Setting SDC name")
		if _, err := r.System.ChangeSdcName(state.ID.ValueString(), plan.Name.ValueString()); err != nil {
			return fmt.Errorf("error setting SDC name: %w", err)
		}
		tflog.Info(ctx, "SDC name set")
	}

	// set Performance Profile
	if plan.PerformanceProfile.ValueString() != state.PerformanceProfile.ValueString() && !plan.PerformanceProfile.IsUnknown() {
		tflog.Info(ctx, "Setting SDC performance profile")
		if _, err := r.System.ChangeSdcPerfProfile(state.ID.ValueString(), plan.PerformanceProfile.ValueString()); err != nil {
			return fmt.Errorf("error setting SDC performance profile: %w", err)
		}
		tflog.Info(ctx, "SDC performance profile set")
	}

	return nil
}
