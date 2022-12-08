package schemastructures

import (
	"context"
	"os"
	"terraform-provider-powerflex/client/helper"

	"github.com/AnshumanPradipPatil1506/goscaleio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var AUTH_SCHEMA map[string]*schema.Schema = map[string]*schema.Schema{
	"host": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_HOST", nil),
	},
	"username": {
		Type:        schema.TypeString,
		Description: "Add Powerflex Manager Username",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_USERNAME", nil),
	},
	"password": {
		Type:        schema.TypeString,
		Description: "Add Powerflex Manager Password",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_PASSWORD", nil),
	},
	"insecure": {
		Type:        schema.TypeString,
		Description: "Add Insecure Value[true/false]",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_INSECURE", "true"), // anshuman check default to set
	},
	"usecerts": {
		Type:        schema.TypeString,
		Description: "Add Use Certificates Value[true/false]",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_USECERTS", "true"), // anshuman check default to set
	},
	"pfxm_version": {
		Type:        schema.TypeString,
		Description: "Add Powerflex Manager Verion",
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("PFxM_VERSION", ""), // anshuman check default to set
	},
}

func AuthConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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

	helper.ENV.Username = username
	helper.ENV.Password = password
	helper.ENV.Host = host
	helper.ENV.Insecure = insecure
	helper.ENV.UseCerts = useCerts
	helper.ENV.Version = version

	client, err := goscaleio.NewClient()

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
	helper.GoscaleioClient = client
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

	_, err = helper.GoscaleioClient.Authenticate(&goscaleioConf)

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
