package winrmclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	bytes []byte
)

func TestWinRMClientInvalidHost(t *testing.T) {

	winRMClient := &WinRMClient{}

	var contexts map[string]string

	_ = json.Unmarshal(bytes, &contexts)

	context := make(map[string]string)

	context["username"] = "administrator"

	context["password"] = "D@ngerous1!"

	context["host"] = "10.247.39.136"

	winRMClient.GetConnection(context, false)

	if winRMClient.Init() {

		ouptut := winRMClient.ExecuteCommand("get-date")

		fmt.Println(ouptut)

	}

	winRMClient.Destroy()

}

func TestWinRMClientFileUpload(t *testing.T) {

	winRMClient := &WinRMClient{}

	var contexts map[string]string

	_ = json.Unmarshal(bytes, &contexts)

	context := make(map[string]string)

	context["username"] = "administrator"

	context["password"] = "D@ngerous1!"

	context["host"] = "10.247.39.136"

	winRMClient.GetConnection(context, false)

	if winRMClient.Init() {

		ouptut := winRMClient.ExecuteCommand("get-date")

		fmt.Println(ouptut)

		//winRMClient.Upload("C:/EMC-ScaleIO-sdc.msi", "/root/powerflex_packages/EMC-ScaleIO-sdc-4.5-0.287.msi")

		ouptut = winRMClient.ExecuteCommand("msiexec.exe /x \"C:\\EMC-ScaleIO-sdc.msi\" MDM_IP=\"10.247.103.160,10.247.103.161\" /q")

		if ouptut == "SUCCESS" {
			fmt.Println(ouptut)
		}

	}

	winRMClient.Destroy()

}
