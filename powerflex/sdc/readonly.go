package sdc

import (
	"context"
	"strconv"

	"terraform-provider-powerflex/client"
	sdcmodels "terraform-provider-powerflex/models/sdc"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceSdcs function to carry sdc resource read schema
func DataSourceSdcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSdcsRead,
		Schema:      sdcmodels.SDCReadOnlyModel,
	}
}

func dataSourceSdcsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	systemID := d.Get("systemid").(string)
	sdcID := d.Get("sdcid").(string)

	goscaleioClient := client.GoscaleioClient

	if systemID == "" {
		tflog.Debug(ctx, "[PowerFlex][dataSourceSdcsRead] Empty value passed for systemid so getting first system available.")
		allSystems, err := goscaleioClient.GetSystems()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Goscaleio systems[goscaleioClient.GetSystems]",
				Detail:   err.Error(),
			})
			tflog.Error(ctx, "[PowerFlex][dataSourceSdcsRead] Unable to get Goscaleio systems[goscaleioClient.C.GetSystems]")
			tflog.Error(ctx, err.Error())
			return diags
		}
		systemID = allSystems[0].ID
	}

	system, err := goscaleioClient.FindSystem(systemID, "", "")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Goscaleio single System[goscaleioClient.FindSystem]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][dataSourceSdcsRead] Unable to get Goscaleio single System[goscaleioClient.FindSystem]")
		tflog.Error(ctx, err.Error())
		return diags
	}

	if sdcID == "" {
		diags = getAllSdcs(ctx, system, sdcID, d)
	} else {
		diags = getSingleSdc(ctx, system, sdcID, d)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
