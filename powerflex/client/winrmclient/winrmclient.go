package winrmclient

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
)

const (
	Command = string("powershell -command \"$Host.UI.RawUI.BufferSize = New-Object Management.Automation.Host.Size (512,50);@@@ | Format-List\"")

	Host = "host"

	UserName = "username"

	Password = "password"

	Port = "port"

	Timeout = "timeout"

	Error = "error"

	Message = "message"
)

type WinRMClient struct {
	client *winrm.Client

	userName string

	password string

	target string

	port int

	timeout int

	shell *winrm.Shell

	errors []map[string]string
}

func (winRMClient *WinRMClient) GetErrors() []map[string]string {

	return winRMClient.errors
}

func (winRMClient *WinRMClient) Destroy() {

	if winRMClient.shell != nil {

		_ = winRMClient.shell.Close()

		// tflog.Debug(nil, string(fmt.Sprintf("disconnecting from %s : %s", winRMClient.target, winRMClient.port)))

	}
}

func (winRMClient *WinRMClient) setTarget(context map[string]string, host bool) *WinRMClient {

	if _, found := context[Host]; found {

		winRMClient.target = context[Host]

	}

	return winRMClient

}

func (winRMClient *WinRMClient) setPort(context map[string]string) *WinRMClient {

	if _, found := context[Port]; found {

		port, _ := strconv.Atoi(context[Port])

		winRMClient.port = port

	} else {

		winRMClient.port = 5985
	}

	return winRMClient

}

func (winRMClient *WinRMClient) setUsername(context map[string]string) *WinRMClient {

	if _, found := context[UserName]; found {

		winRMClient.userName = context[UserName]

	}

	return winRMClient
}

func (winRMClient *WinRMClient) setPassword(context map[string]string) *WinRMClient {

	if _, found := context[Password]; found {

		winRMClient.password = context[Password]

	}

	return winRMClient
}

func (winRMClient *WinRMClient) setTimeout(context map[string]string) *WinRMClient {

	if _, found := context[Timeout]; found {

		timeout, _ := strconv.Atoi(context[Timeout])

		winRMClient.timeout = timeout

	} else {

		winRMClient.timeout = 60
	}

	return winRMClient

}

func (winRMClient *WinRMClient) GetConnection(context map[string]string, host bool) *WinRMClient {

	return winRMClient.setTarget(context, host).setPort(context).setUsername(context).setPassword(
		context).setTimeout(context)

}

func (winRMClient *WinRMClient) ExecuteCommand(command string) string {

	output := ""

	err := ""

	output, err, _, _ = winRMClient.client.RunWithString(strings.ReplaceAll(Command, "@@@", command), "")

	if len(err) > 0 {

		winRMClient.errors = append(winRMClient.errors, map[string]string{
			Error:   err,
			Message: "failed to execute command [" + command + "] on target " + winRMClient.target,
		})

		// tflog.Warn(nil, string(fmt.Sprintf("failed to execute command :%s on target %s with error %s", command, winRMClient.target, err)))

	}

	return string(output)
}

func (winRMClient *WinRMClient) Init() (result bool) {

	result = false

	errorMessage := fmt.Sprintf("Failed to establish %s connection on %s:%d", "WinRM", winRMClient.target, winRMClient.port)

	endpoint := &winrm.Endpoint{Host: winRMClient.target, Port: winRMClient.port, HTTPS: false, Insecure: false, CACert: nil, Cert: nil, Key: nil, Timeout: time.Duration(winRMClient.timeout) * time.Second}

	if strings.Contains(winRMClient.userName, "\\") {

		params := winrm.Parameters{TransportDecorator: func() winrm.Transporter { return &winrm.ClientNTLM{} }}

		winRMClient.client, _ = winrm.NewClientWithParameters(endpoint, winRMClient.userName, winRMClient.password, &params)

	} else {

		winRMClient.client, _ = winrm.NewClient(endpoint, winRMClient.userName, winRMClient.password)

	}

	var err error

	winRMClient.shell, err = winRMClient.client.CreateShell()

	if winRMClient.shell != nil {

		result = true

		// tflog.Debug(nil, string(fmt.Sprintf("Connected to %s:%d", winRMClient.target, winRMClient.port)))

	} else if err != nil {

		// tflog.Warn(nil, string(fmt.Sprintf("error %v occurred for %v host...", err.Error(), winRMClient.target)))

		// tflog.Debug(nil, string(fmt.Sprintf("Failed to establish %s connection on %s:%d", "WinRM", winRMClient.target, winRMClient.port)))

		if strings.Contains(string(err.Error()), "connection refused") || strings.Contains(string(err.Error()), "invalid port") {

			errorMessage = fmt.Sprintf("Invalid port %d, Please verify that port %d is up", winRMClient.port, winRMClient.port)

		} else if strings.Contains(string(err.Error()), "i/o timeout") {

			errorMessage = fmt.Sprintf("%s Timed out for %s:%d", "WinRM", winRMClient.target, winRMClient.port)

		} else if strings.Contains(string(err.Error()), "http response error: 401") {

			errorMessage = fmt.Sprintf("Invalid Credentials %s:%d", winRMClient.target, winRMClient.port)
		} else {

			errorMessage = fmt.Sprintf("Invalid port %d, Please verify that port %d is up", winRMClient.port, winRMClient.port)

		}

		//tflog.Warn(nil, string(errorMessage))

		winRMClient.errors = append(winRMClient.errors, map[string]string{
			Error:   err.Error(),
			Message: errorMessage,
		})

	}

	return
}

func (winRMClient *WinRMClient) newCopyClient() (*winrmcp.Winrmcp, error) {
	addr := fmt.Sprintf("%s:%d", winRMClient.target, winRMClient.port)

	config := winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     winRMClient.userName,
			Password: winRMClient.password,
		},
		Https:                 false,
		Insecure:              false,
		OperationTimeout:      time.Duration(winRMClient.timeout) * time.Second,
		CACertBytes:           nil,
		MaxOperationsPerShell: 15, // lowest common denominator
	}

	config.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	return winrmcp.New(addr, &config)
}

func (winRMClient *WinRMClient) Upload(dstPath string, srcPath string) error {

	// Convert the encoded string to a byte slice
	fileContent, err := decodeString(srcPath)
	if err != nil {
		log.Fatalf("Failed to decode string: %v", err)
	}

	// Create a bytes.Reader from the file content
	input := strings.NewReader(string(fileContent))

	wcp, err := winRMClient.newCopyClient()
	if err != nil {
		return err
	}

	err = wcp.Write(dstPath, input)
	if err != nil {
		return err
	}
	return nil
}

// Function to decode base64 encoded string to byte slice
func decodeString(encodedString string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}
