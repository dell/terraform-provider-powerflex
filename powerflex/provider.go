package powerflex

import (
	"context"
	"os"
	pd "terraform-provider-powerflex/powerflex/protectionDomain"
	volumedatasource "terraform-provider-powerflex/powerflex/volume"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &powerflexProvider{}
)

// New returns the powerflex provider
func New() provider.Provider {
	return &powerflexProvider{}
}

type powerflexProvider struct{}

type powerflexProviderModel struct {
	Host             types.String `tfsdk:"host"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	Insecure         types.String `tfsdk:"insecure"`
	UseCerts         types.String `tfsdk:"usecerts"`
	PowerflexVersion types.String `tfsdk:"powerflex_version"`
}

func (p *powerflexProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerflex"
}

func (p *powerflexProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Sensitive:   true,
			},
			"powerflex_version": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"usecerts": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"insecure": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *powerflexProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring powerflex client")

	var config powerflexProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown powerflex API Host",
			"The provider cannot create the powerflex API client as there is an unknown configuration value for the powerflex API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the powerflex_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown powerflex API Username",
			"The provider cannot create the powerflex API client as there is an unknown configuration value for the powerflex API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the powerflex_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown powerflex API Password",
			"The provider cannot create the powerflex API client as there is an unknown configuration value for the powerflex API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the powerflex_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("POWERFLEX_HOST")
	username := os.Getenv("POWERFLEX_USERNAME")
	password := os.Getenv("POWERFLEX_PASSWORD")

	insecure := os.Getenv("POWERFLEX_HOST")
	usecerts := os.Getenv("POWERFLEX_USERNAME")
	powerflexVersion := os.Getenv("POWERFLEX_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueString()
	}

	if !config.UseCerts.IsNull() {
		usecerts = config.UseCerts.ValueString()
	}

	if !config.PowerflexVersion.IsNull() {
		powerflexVersion = config.PowerflexVersion.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing powerflex API Host",
			"The provider cannot create the powerflex API client as there is a missing or empty value for the powerflex API host. "+
				"Set the host value in the configuration or use the powerflex_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing powerflex API Username",
			"The provider cannot create the powerflex API client as there is a missing or empty value for the powerflex API username. "+
				"Set the username value in the configuration or use the powerflex_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing powerflex API Password",
			"The provider cannot create the powerflex API client as there is a missing or empty value for the powerflex API password. "+
				"Set the password value in the configuration or use the powerflex_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "powerflex_host", host)
	ctx = tflog.SetField(ctx, "powerflex_username", username)
	ctx = tflog.SetField(ctx, "powerflex_password", password)

	ctx = tflog.SetField(ctx, "insecure", insecure)
	ctx = tflog.SetField(ctx, "usecerts", usecerts)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "powerflex_password")

	tflog.Debug(ctx, "Creating powerflex client")

	// Create a new powerflex client using the configuration values
	client, err := goscaleio.NewClientWithArgs(host, powerflexVersion, true, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create powerflex API Client",
			"An unexpected error occurred when creating the powerflex API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"powerflex Client Error: "+err.Error(),
		)
		return
	}

	var goscaleioConf goscaleio.ConfigConnect = goscaleio.ConfigConnect{}
	goscaleioConf.Endpoint = host
	goscaleioConf.Username = username
	goscaleioConf.Version = powerflexVersion
	goscaleioConf.Password = password

	_, err = client.Authenticate(&goscaleioConf)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Auth Goscaleio API Client",
			"An unexpected error occurred when creating the Unable to Auth Goscaleio API Client. "+
				"Unable to Auth Goscaleio API Client.\n\n"+
				"powerflex Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured powerflex client", map[string]any{"success": true})
}

func (p *powerflexProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		volumedatasource.VolumeDataSource,
		pd.ProtectionDomainDataSource,
	}
}

func (p *powerflexProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
