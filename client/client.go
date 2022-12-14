package client

import (
	"bytes"
	"encoding/json"

	"github.com/dell/goscaleio"
)

// GoscaleioClient object for client to carry goscaleio object sigular
var GoscaleioClient *goscaleio.Client

// ENVVars is a structure for env vars passed to provider
type ENVVars struct {
	Host     string
	Username string
	Password string
	Insecure string
	UseCerts string
	Version  string
}

// ENV Variable for env vars
var ENV *ENVVars = &ENVVars{
	Host:     "",
	Username: "",
	Password: "",
	Insecure: "",
	UseCerts: "",
	Version:  "",
}

// PrettyJSON functio for Print json pretty format
func PrettyJSON(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}
