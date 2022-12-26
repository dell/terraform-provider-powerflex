package powerflex

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

var ProviderConfigForTesting = ``

func init() {
	err := godotenv.Load("POWERFLEX_TERRAFORM_TEST.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ProviderConfigForTesting = `
		provider "powerflex" {
			username = "` + os.Getenv("GOSCALEIO_USERNAME") + `"
			password = "` + os.Getenv("GOSCALEIO_PASSWORD") + `"
			host = "` + os.Getenv("GOSCALEIO_ENDPOINT") + `"
			insecure = ""
			usecerts = ""
			powerflex_version = ""
		}
	`
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerflex": providerserver.NewProtocol6WithError(New()),
}
