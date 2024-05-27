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
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
)

// ScpProvisioner - scp client
type ScpProvisioner struct {
	logger Logger
	client *scp.Client
}

// NewScpProvisioner - creates new scp client
func NewScpProvisioner(prov *SSHProvisioner) *ScpProvisioner {
	scpClient := scp.NewConfigurer("", nil).SSHClient(prov.sshClient).Create()
	return &ScpProvisioner{
		client: &scpClient,
		logger: prov.logger,
	}
}

// Upload Function to upload file to remote host
func (p *ScpProvisioner) Upload(src, dst, perm string) error {
	p.logger.Printf("Reading input file")
	// read src file
	fileContent, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if perm == "" {
		perm = "0655"
		p.logger.Printf("Using default file permission %s on remote host", perm)
	}

	// Create a bytes.Reader from the file content
	input := strings.NewReader(string(fileContent))

	p.logger.Printf("Uploading file to %s", dst)
	// the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
	return p.client.CopyFile(context.Background(), input, dst, perm)
}

// decodeString Function to decode base64 encoded string to byte slice
func decodeString(encodedString string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}
