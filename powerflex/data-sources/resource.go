package datasources

import (
	"context"
	"terraform-provider-powerflex/client/helper"
	schemastructures "terraform-provider-powerflex/powerflex/schema-structures"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSdcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceSdcsUpdate,
		Schema:      schemastructures.SDC_UPDATE_RESOURCE_NANE_SCHEMA,
	}
}

func resourceSdcsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	goscaleioClient := helper.GoscaleioClient

	systemID := d.Get("systemid").(string)
	sdcID := d.Get("id").(string)
	newName := d.Get("name").(string)

	system, err := goscaleioClient.FindSystem(systemID, "", "")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Goscaleio single System[goscaleioClient.FindSystem] with id == " + systemID,
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][dataSourceSdcsRead] Unable to get Goscaleio single System[goscaleioClient.FindSystem] with id == "+systemID)
		tflog.Error(ctx, err.Error())
		return diags
	}

	nameChngedSdc, err := system.ChangeSdcName(sdcID, newName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to change SDC Name[goscaleioClient.resourceSdcsUpdate] with id == " + systemID,
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][dataSourceSdcsRead] Unable to change SDC Name[goscaleioClient.resourceSdcsUpdate] with id == "+systemID)
		tflog.Error(ctx, err.Error())
		return diags
	}

	resultMap := make([]map[string]interface{}, 0)
	resultSdcMap := schemastructures.NameChangedSdcToMap(*nameChngedSdc.Sdc)

	tflog.Debug(ctx, "[PowerFlex][NameChangedSdcToMap] sdc - "+helper.PrettyJson(resultSdcMap))
	resultMap = append(resultMap, resultSdcMap)

	if err := d.Set("sdcs", resultMap); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// func getNameChangedSdc(system *goscaleio.System, sdcID string, ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	singleSdc, err := system.FindSdc("ID", sdcID)
// 	if err != nil {
// 		diags = append(diags, diag.Diagnostic{
// 			Severity: diag.Error,
// 			Summary:  "Unable to get Goscaleio single SDC[system.FindSdc]",
// 			Detail:   err.Error(),
// 		})
// 		tflog.Error(ctx, "[PowerFlex][getSingleSdc] Unable to get Goscaleio single SDC[system.FindSdc]")
// 		tflog.Error(ctx, err.Error())
// 		return diags
// 	}
// 	tflog.Debug(ctx, "[PowerFlex][getSingleSdc] sdc - "+helper.PrettyJson(singleSdc))

// 	resultMap := make([]map[string]interface{}, 0)
// 	resultSdcMap := schemastructures.NameChangedSdcToMap(*singleSdc.Sdc)

// 	tflog.Debug(ctx, "[PowerFlex][Anshuman] sdc - "+helper.PrettyJson(resultSdcMap))
// 	resultMap = append(resultMap, resultSdcMap)

// 	if err := d.Set("sdcs", resultMap); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return diags
// }
