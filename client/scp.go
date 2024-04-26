package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
)

type ScpProvisioner struct {
	logger Logger
	client *scp.Client
}

func NewScpProvisioner(prov *SshProvisioner) *ScpProvisioner {
	scpClient := scp.NewConfigurer("", nil).SSHClient(prov.sshClient).Create()
	return &ScpProvisioner{
		client: &scpClient,
		logger: prov.logger,
	}
}

func (p *ScpProvisioner) Upload(srcB64, dst, perm string) error {
	p.logger.Printf("Decoding input file")
	fileContent, err := decodeString(srcB64)
	if err != nil {
		return fmt.Errorf("failed to decode input file as base64: %w", err)
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

// Function to decode base64 encoded string to byte slice
func decodeString(encodedString string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}
