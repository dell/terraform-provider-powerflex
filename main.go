package main

import (
	"terraform-provider-powerflex/powerflex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return powerflex.Provider()
		},
	})
}

// func main() {
// 	var debug bool
// 	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
// 	flag.Parse()
// 	err := providerserver.Serve(context.Background(), powerflex.Provider, providerserver.ServeOpts{
// 		Address: "dell.com/ses/powerflex",
// 		Debug:   debug,
// 	})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// }
