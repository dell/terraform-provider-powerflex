package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var storagepoolResourceSchemaDescriptions = struct {
	storagepoolReourceSchema string
	LastUpdated              string
	SystemID                 string
	Name                     string
	ProtectionDomainID       string
	MediaType                string
	UseRmcache               string
	UseRfCache               string
	Links                    string
	LinksRel                 string
	LinksHref                string
}{
	storagepoolReourceSchema: "",
	LastUpdated:              "",
	SystemID:                 "",
	Name:                     "",
	ProtectionDomainID:       "",
	MediaType:                "",
	UseRmcache:               "",
	UseRfCache:               "",
	Links:                    "",
	LinksRel:                 "",
	LinksHref:                "",
}

// StoragepoolReourceSchema - varible holds schema for Storagepool
var StoragepoolReourceSchema schema.Schema = schema.Schema{
	Description: storagepoolResourceSchemaDescriptions.storagepoolReourceSchema,
	Attributes: map[string]schema.Attribute{
		"last_updated": schema.StringAttribute{
			Description: storagepoolResourceSchemaDescriptions.LastUpdated,
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "",
			Computed:    true,
		},
		"systemid": schema.StringAttribute{
			Description: storagepoolResourceSchemaDescriptions.SystemID,
			Computed:    true,
		},
		"protection_domain_id": schema.StringAttribute{
			Description: storagepoolResourceSchemaDescriptions.ProtectionDomainID,
			Required:    true,
		},
		"name": schema.StringAttribute{
			Description: storagepoolResourceSchemaDescriptions.Name,
			Required:    true,
		},
		"media_type": schema.StringAttribute{
			Description: storagepoolResourceSchemaDescriptions.MediaType,
			Required:    true,
		},
		"use_rmcache": schema.BoolAttribute{
			Description: storagepoolResourceSchemaDescriptions.UseRmcache,
			Computed:    true,
		},
		"use_rfcache": schema.BoolAttribute{
			Description: storagepoolResourceSchemaDescriptions.UseRfCache,
			Computed:    true,
		},
		"links": schema.ListNestedAttribute{
			Description: storagepoolResourceSchemaDescriptions.Links,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Description: storagepoolResourceSchemaDescriptions.LinksRel,
						Computed:    true,
					},
					"href": schema.StringAttribute{
						Description: storagepoolResourceSchemaDescriptions.LinksHref,
						Computed:    true,
					},
				},
			},
		},
	},
}
