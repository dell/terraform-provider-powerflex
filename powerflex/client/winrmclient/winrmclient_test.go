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

		winRMClient.Upload("C:\\krunal.txt", "S3J1bmFsIFRoYWtrYXIKCgoKCgoKCgoKCgoKCgoKCgpHb29kIEJ5ZSEK")

	}

	winRMClient.Destroy()

}
