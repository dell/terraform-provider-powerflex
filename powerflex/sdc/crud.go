package sdc

import (
	"context"
	"terraform-provider-powerflex/client"
	sdcmodels "terraform-provider-powerflex/models/sdc"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceSdcs function to carry sdc CRUD opration schema
func ResourceSdcs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSdcCreate,
		ReadContext:   resourceSdcRead,
		UpdateContext: resourceSdcsUpdate,
		DeleteContext: resourceSdcDelete,
		Schema:        sdcmodels.SDCCRUDModel,
	}
}

func resourceSdcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if !d.HasChange("name") {
		return resourceSdcRead(ctx, d, m)
	}

	systemID := d.Get("systemid").(string)
	sdcID := d.Get("sdcid").(string)
	sdcName := d.Get("name").(string)
	goscaleioClient := client.GoscaleioClient
	system, err := goscaleioClient.FindSystem(systemID, "", "")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Goscaleio single System[goscaleioClient.FindSystem]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][resourceSdcCreate] Unable to get Goscaleio single System[goscaleioClient.FindSystem]")
		tflog.Error(ctx, err.Error())
		return diags
	}

	nameChng, err := system.ChangeSdcName(sdcID, sdcName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to Change name of SDC [system.ChangeSdcName]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][resourceSdcCreate] Unable to Change name of SDC[system.ChangeSdcName]")
		tflog.Error(ctx, err.Error())
		return diags
	}

	tflog.Debug(ctx, "[PowerFlex][resourceSdcCreate] nameChng result - "+client.PrettyJSON(nameChng))
	resultMap := make([]map[string]interface{}, 0)
	resultSDC := make(map[string]interface{})

	resultSDC["id"] = sdcID
	resultSDC["systemid"] = systemID
	resultSDC["name"] = sdcName

	resultMap = append(resultMap, resultSDC)
	tflog.Debug(ctx, "[PowerFlex][resourceSdcCreate] sdc resultMap - "+client.PrettyJSON(resultMap))

	if err := d.Set("sdcs", resultMap); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("name").(string))
	return resourceSdcRead(ctx, d, m)
}

func resourceSdcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	systemID := d.Get("systemid").(string)
	sdcID := d.Get("sdcid").(string)

	goscaleioClient := client.GoscaleioClient

	system, err := goscaleioClient.FindSystem(systemID, "", "")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[PowerFlex][resourceSdcRead] Unable to get Goscaleio single System[goscaleioClient.FindSystem]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][resourceSdcRead] Unable to get Goscaleio single System[goscaleioClient.FindSystem]")
		tflog.Error(ctx, err.Error())
		return diags
	}

	diags = getSingleSdc(ctx, system, sdcID, d)

	return diags
}

func resourceSdcsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSdcCreate(ctx, d, m)
}

func resourceSdcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	// return resourceSdcRead(ctx, d, m)
	return diags
}
