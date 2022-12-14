package auth

import (
	"context"
	"terraform-provider-powerflex/client"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GoscaleioAuth function for authenticating goscaleio
func GoscaleioAuth(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)
	insecure := d.Get("insecure").(string)
	useCerts := d.Get("usecerts").(string)
	version := d.Get("powerflex_version").(string)

	if username == "" || password == "" || host == "" || insecure == "" || useCerts == "" || version == "" {
		tflog.Debug(ctx, "[PowerFlex][AuthConfigure] username not available")
		tflog.Debug(ctx, "[PowerFlex][AuthConfigure] password not available")
		tflog.Debug(ctx, "[PowerFlex][AuthConfigure] host not available "+host)
		tflog.Debug(ctx, "[PowerFlex][AuthConfigure] insecure not available "+insecure)
		tflog.Debug(ctx, "[PowerFlex][AuthConfigure] useCerts not available "+useCerts)

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "username, password, host, insecure, useCerts, version are required parameters.",
			Detail:   "Required parameters should be passed either in terraform provider or in ENV variables please see documentation.",
		})
		tflog.Error(ctx, "[PowerFlex][AuthConfigure] Unable to create Goscaleio client because of less parameters provided.")
		return nil, diags
	}

	client.ENV.Username = username
	client.ENV.Password = password
	client.ENV.Host = host
	client.ENV.Insecure = insecure
	client.ENV.UseCerts = useCerts
	client.ENV.Version = version
	// func NewClientWithArgs(
	// 	endpoint string,
	// 	version string,
	// 	insecure,
	// 	useCerts bool
	goscaleioclient, err := goscaleio.NewClientWithArgs(host, version, insecure == "true", useCerts == "true")

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Goscaleio client",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][AuthConfigure] Unable to create Goscaleio client")
		tflog.Error(ctx, err.Error())
		return nil, diags
	}

	tflog.Info(ctx, "[PowerFlex][AuthConfigure] Successfuly Created ScaleIO Client")
	client.GoscaleioClient = goscaleioclient

	var goscaleioConf goscaleio.ConfigConnect = goscaleio.ConfigConnect{}
	goscaleioConf.Endpoint = host
	goscaleioConf.Username = username
	goscaleioConf.Version = version
	goscaleioConf.Password = password

	_, err = client.GoscaleioClient.Authenticate(&goscaleioConf)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error authenticating Goscaleio client",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][AuthConfigure] Error authenticating Goscaleio client")
		tflog.Error(ctx, err.Error())
		return nil, diags
	}
	tflog.Info(ctx, "[PowerFlex][AuthConfigure] Successfuly logged in to ScaleIO Gateway")

	return nil, diags
}
