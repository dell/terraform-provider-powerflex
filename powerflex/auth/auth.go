package auth

import (
	"context"
	"os"
	"terraform-provider-powerflex/client"

	"github.com/AnshumanPradipPatil1506/goscaleio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GoscaleioAuth function for authenticating goscaleio
func GoscaleioAuth(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(string)
	useCerts := d.Get("usecerts").(string)
	version := d.Get("pfxm_version").(string)

	os.Setenv("GOSCALEIO_ENDPOINT", host)
	os.Setenv("GOSCALEIO_VERSION", version)
	os.Setenv("GOSCALEIO_INSECURE", insecure)
	os.Setenv("GOSCALEIO_USECERTS", useCerts)

	client.ENV.Username = username
	client.ENV.Password = password
	client.ENV.Host = host
	client.ENV.Insecure = insecure
	client.ENV.UseCerts = useCerts
	client.ENV.Version = version

	goscaleioclient, err := goscaleio.NewClient()

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
	tflog.Debug(ctx, "[PowerFlex][AuthConfigure] username "+username)
	tflog.Debug(ctx, "[PowerFlex][AuthConfigure] password "+password)
	tflog.Debug(ctx, "[PowerFlex][AuthConfigure] host here"+host)
	tflog.Debug(ctx, "[PowerFlex][AuthConfigure] insecure "+insecure)
	tflog.Debug(ctx, "[PowerFlex][AuthConfigure] useCerts "+useCerts)

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
