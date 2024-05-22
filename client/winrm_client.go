/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://mozilla.org/MPL/2.0/

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
)

// Constants
const (
	Command = string("powershell -command @@@")

	Host = "host"

	UserName = "username"

	Password = "password"

	Port = "port"

	Timeout = "timeout"

	Error = "error"

	Message = "message"
)

// WinRMClient struct
type WinRMClient struct {
	Client *winrm.Client

	UserName string

	Password string

	Target string

	Port int

	Timeout int

	Shell *winrm.Shell

	Errors []map[string]string
}

// GetErrors returns errors
func (winRMClient *WinRMClient) GetErrors() []map[string]string {
	return winRMClient.Errors
}

// Destroy closes the shell
func (winRMClient *WinRMClient) Destroy() {

	if winRMClient.Shell != nil {

		_ = winRMClient.Shell.Close()

	}
}

// setTarget sets the target
func (winRMClient *WinRMClient) setTarget(context map[string]string, host bool) *WinRMClient {

	if _, found := context[Host]; found {

		winRMClient.Target = context[Host]

	}

	return winRMClient

}

// setPort sets the port
func (winRMClient *WinRMClient) setPort(context map[string]string) *WinRMClient {

	if _, found := context[Port]; found {

		port, _ := strconv.Atoi(context[Port])

		winRMClient.Port = port

	} else {

		winRMClient.Port = 5985
	}

	return winRMClient

}

// setUsername sets the username
func (winRMClient *WinRMClient) setUsername(context map[string]string) *WinRMClient {

	if _, found := context[UserName]; found {

		winRMClient.UserName = context[UserName]

	}

	return winRMClient
}

// setPassword sets the password
func (winRMClient *WinRMClient) setPassword(context map[string]string) *WinRMClient {

	if _, found := context[Password]; found {

		winRMClient.Password = context[Password]

	}

	return winRMClient
}

// setTimeout sets the timeout
func (winRMClient *WinRMClient) setTimeout(context map[string]string) *WinRMClient {

	if _, found := context[Timeout]; found {

		timeout, _ := strconv.Atoi(context[Timeout])

		winRMClient.Timeout = timeout

	} else {

		winRMClient.Timeout = 60
	}

	return winRMClient

}

// GetConnection returns the connection
func (winRMClient *WinRMClient) GetConnection(context map[string]string, host bool) *WinRMClient {

	return winRMClient.setTarget(context, host).setPort(context).setUsername(context).setPassword(
		context).setTimeout(context)

}

// ExecuteCommand executes the command
func (winRMClient *WinRMClient) ExecuteCommand(command string) (string, error) {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	output, err, _, _ := winRMClient.Client.RunWithContextWithString(ctx, strings.ReplaceAll(Command, "@@@", command), "")

	if len(err) > 0 {

		output = "FAIL"

		return output, fmt.Errorf("failed to execute command [%s] on target %s", command, winRMClient.Target)

	}

	if output == "" {
		output = "SUCCESS"
	}

	return output, nil
}

// Init initializes the connection
func (winRMClient *WinRMClient) Init() (bool, error) {

	errorMessage := ""

	endpoint := &winrm.Endpoint{Host: winRMClient.Target, Port: winRMClient.Port, HTTPS: false, Insecure: false, CACert: nil, Cert: nil, Key: nil, Timeout: time.Duration(winRMClient.Timeout) * time.Second}

	if strings.Contains(winRMClient.UserName, "\\") {

		params := winrm.Parameters{TransportDecorator: func() winrm.Transporter { return &winrm.ClientNTLM{} }}

		winRMClient.Client, _ = winrm.NewClientWithParameters(endpoint, winRMClient.UserName, winRMClient.Password, &params)

	} else {

		winRMClient.Client, _ = winrm.NewClient(endpoint, winRMClient.UserName, winRMClient.Password)

	}

	var err error

	winRMClient.Shell, err = winRMClient.Client.CreateShell()

	if winRMClient.Shell != nil {

		return true, nil

	}

	errorMessage = fmt.Sprintf("Failed to establish %s connection on %s:%d", "WinRM", winRMClient.Target, winRMClient.Port)

	if strings.Contains(string(err.Error()), "connection refused") || strings.Contains(string(err.Error()), "invalid port") {

		errorMessage = fmt.Sprintf("Invalid port %d, Please verify that port %d is up", winRMClient.Port, winRMClient.Port)

	} else if strings.Contains(string(err.Error()), "i/o timeout") {

		errorMessage = fmt.Sprintf("%s Timed out for %s:%d", "Connection", winRMClient.Target, winRMClient.Port)

	} else if strings.Contains(string(err.Error()), "http response error: 401") {

		errorMessage = fmt.Sprintf("Invalid Credentials %s:%d", winRMClient.Target, winRMClient.Port)
	}

	return false, fmt.Errorf(errorMessage)
}

// newCopyClient creates a new copy client
func (winRMClient *WinRMClient) newCopyClient() (*winrmcp.Winrmcp, error) {
	addr := fmt.Sprintf("%s:%d", winRMClient.Target, winRMClient.Port)

	config := winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     winRMClient.UserName,
			Password: winRMClient.Password,
		},
		Https:                 false,
		Insecure:              false,
		OperationTimeout:      time.Duration(winRMClient.Timeout) * time.Second,
		CACertBytes:           nil,
		MaxOperationsPerShell: 15, // lowest common denominator
	}

	config.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	return winrmcp.New(addr, &config)
}

// Upload uploads a file
func (winRMClient *WinRMClient) Upload(dstPath string, srcPath string) error {

	input, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer input.Close()

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
