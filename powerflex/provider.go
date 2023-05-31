package powerflex

import (
	"context"
	"os"
	"strconv"

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

// New - returns new provider struct definition.
func New() provider.Provider {
	return &powerflexProvider{}
}

type powerflexProvider struct{}

// powerflexProviderModel - provider input struct.
type powerflexProviderModel struct {
	EndPoint         types.String `tfsdk:"endpoint"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	Insecure         types.Bool   `tfsdk:"insecure"`
	Timeout          types.Int64  `tfsdk:"timeout"`
	GatewayInstaller types.Bool   `tfsdk:"gatewayinstaller"`
}

// Metadata - provider metadata AKA name.
func (p *powerflexProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerflex"
}

// GetSchema - provider schema.
func (p *powerflexProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Terraform provider for Dell PowerFlex " +
			"can be used to interact with a Dell PowerFlex array in order to manage the array resources.",
		MarkdownDescription: "The Terraform provider for Dell PowerFlex " +
			"can be used to interact with a Dell PowerFlex array in order to manage the array resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description:         "The PowerFlex Gateway server URL (inclusive of the port).",
				MarkdownDescription: "The PowerFlex Gateway server URL (inclusive of the port).",
				Required:            true,
			},
			"username": schema.StringAttribute{
				Description:         "The username required for authentication.",
				MarkdownDescription: "The username required for authentication.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				Description:         "The password required for the authentication.",
				MarkdownDescription: "The password required for the authentication.",
				Required:            true,
				Sensitive:           true,
			},
			"insecure": schema.BoolAttribute{
				Description:         "Specifies if the user wants to skip SSL verification.",
				MarkdownDescription: "Specifies if the user wants to skip SSL verification.",
				Optional:            true,
			},
			"timeout": schema.Int64Attribute{
				Description:         "HTTPS timeout.",
				MarkdownDescription: "HTTPS timeout.",
				Optional:            true,
			},
			"gatewayinstaller": schema.BoolAttribute{
				Description:         "Specifies if the user wants to use Gateway Installer.",
				MarkdownDescription: "Specifies if the user wants to use Gateway Installer.",
				Optional:            true,
			},
		},
	}
}

// Configure - provider pre-initiate calle function.
func (p *powerflexProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring powerflex client")

	var config powerflexProviderModel
	var timeout int
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.EndPoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown powerflex API EndPoint",
			"The provider cannot create the powerflex API client as there is an unknown configuration value for the powerflex API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the POWERFLEX_ENDPOINT environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown powerflex API Username",
			"The provider cannot create the powerflex API client as there is an unknown configuration value for the powerflex API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the POWERFLEX_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown powerflex API Password",
			"The provider cannot create the powerflex API client as there is an unknown configuration value for the powerflex API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the POWERFLEX_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if config.GatewayInstaller.String() == "true" {

		endpoint := os.Getenv("GATEWAY_ENDPOINT")
		username := os.Getenv("GATEWAY_USERNAME")
		password := os.Getenv("GATEWAY_PASSWORD")
		insecure := os.Getenv("GATEWAY_INSECURE") == "true"

		if !config.EndPoint.IsNull() {
			endpoint = config.EndPoint.ValueString()
		}

		if !config.Username.IsNull() {
			username = config.Username.ValueString()
		}

		if !config.Password.IsNull() {
			password = config.Password.ValueString()
		}
		if !config.Insecure.IsNull() {
			insecure = config.Insecure.ValueBool()
		}
		if !config.Timeout.IsNull() {
			timeout = int(config.Timeout.ValueInt64())
		}

		if endpoint == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("endpoint"),
				"Missing gateway API Endpoint",
				"The provider cannot create the gateway API client as there is a missing or empty value for the gateway API endpoint. "+
					"Set the endpoint value in the configuration or use the GATEWAY_ENDPOINT environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Missing gateway API Username",
				"The provider cannot create the gateway API client as there is a missing or empty value for the gateway API username. "+
					"Set the username value in the configuration or use the GATEWAY_USERNAME environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Missing gateway API Password",
				"The provider cannot create the gateway API client as there is a missing or empty value for the gateway API password. "+
					"Set the password value in the configuration or use the GATEWAY_PASSWORD environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if resp.Diagnostics.HasError() {
			return
		}

		ctx = tflog.SetField(ctx, "gateway_endpoint", endpoint)
		ctx = tflog.SetField(ctx, "gateway_username", username)
		ctx = tflog.SetField(ctx, "gateway_password", password)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "gateway_password")
		ctx = tflog.SetField(ctx, "insecure", insecure)
		tflog.Debug(ctx, "Creating gateway client")

		// Create a new powerflex client using the configuration values
		Client, err := goscaleio.NewGateway(endpoint, username, password, insecure, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create gateway API Client",
				"An unexpected error occurred when creating the gateway API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"gateway Client Error: "+err.Error(),
			)
			return
		}

		var goscaleioConf goscaleio.ConfigConnect = goscaleio.ConfigConnect{}
		goscaleioConf.Endpoint = endpoint
		goscaleioConf.Username = username
		goscaleioConf.Version = ""
		goscaleioConf.Password = password

		resp.DataSourceData = Client
		resp.ResourceData = Client

	} else {
		endpoint := os.Getenv("POWERFLEX_ENDPOINT")
		username := os.Getenv("POWERFLEX_USERNAME")
		password := os.Getenv("POWERFLEX_PASSWORD")
		insecure := os.Getenv("POWERFLEX_INSECURE") == "true"
		if os.Getenv("POWERFLEX_TIMEOUT") != "" {
			var err error
			timeout, err = strconv.Atoi(os.Getenv("POWERFLEX_TIMEOUT"))
			if err != nil {
				resp.Diagnostics.AddError("Invalid POWERFLEX_TIMEOUT", err.Error())
			}
		}

		if !config.EndPoint.IsNull() {
			endpoint = config.EndPoint.ValueString()
		}

		if !config.Username.IsNull() {
			username = config.Username.ValueString()
		}

		if !config.Password.IsNull() {
			password = config.Password.ValueString()
		}
		if !config.Insecure.IsNull() {
			insecure = config.Insecure.ValueBool()
		}
		if !config.Timeout.IsNull() {
			timeout = int(config.Timeout.ValueInt64())
		}

		if endpoint == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("endpoint"),
				"Missing powerflex API Endpoint",
				"The provider cannot create the powerflex API client as there is a missing or empty value for the powerflex API endpoint. "+
					"Set the endpoint value in the configuration or use the POWERFLEX_ENDPOINT environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Missing powerflex API Username",
				"The provider cannot create the powerflex API client as there is a missing or empty value for the powerflex API username. "+
					"Set the username value in the configuration or use the POWERFLEX_USERNAME environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Missing powerflex API Password",
				"The provider cannot create the powerflex API client as there is a missing or empty value for the powerflex API password. "+
					"Set the password value in the configuration or use the POWERFLEX_PASSWORD environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if resp.Diagnostics.HasError() {
			return
		}

		ctx = tflog.SetField(ctx, "powerflex_endpoint", endpoint)
		ctx = tflog.SetField(ctx, "powerflex_username", username)
		ctx = tflog.SetField(ctx, "powerflex_password", password)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "powerflex_password")
		ctx = tflog.SetField(ctx, "insecure", insecure)
		ctx = tflog.SetField(ctx, "timeout", timeout)
		tflog.Debug(ctx, "Creating powerflex client")

		// Create a new powerflex client using the configuration values
		Client, err := goscaleio.NewClientWithArgs(endpoint, "", int64(timeout), insecure, true)
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
		goscaleioConf.Endpoint = endpoint
		goscaleioConf.Username = username
		goscaleioConf.Version = ""
		goscaleioConf.Password = password

		_, err = Client.Authenticate(&goscaleioConf)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Authenticate Goscaleio API Client",
				"An unexpected error occurred when authenticating the Goscaleio API Client. "+
					"Unable to Authenticate Goscaleio API Client.\n\n"+
					"powerflex Client Error: "+err.Error(),
			)
			return
		}

		resp.DataSourceData = Client
		resp.ResourceData = Client
	}

	tflog.Info(ctx, "Configured powerflex client", map[string]any{"success": true})
}

// DataSources - returns array of all datasources.
func (p *powerflexProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		VolumeDataSource,
		SDCDataSource,
		ProtectionDomainDataSource,
		StoragePoolDataSource,
		SnapshotPolicyDataSource,
		SDSDataSource,
		DeviceDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *powerflexProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProtectionDomainResource,
		NewSDSResource,
		NewVolumeResource,
		NewSnapshotResource,
		SDCResource,
		StoragepoolResource,
		NewSDCVolumesMappingResource,
		NewDeviceResource,
		NewPackageResource,
		NewSDCExpansionResource,
	}
}
