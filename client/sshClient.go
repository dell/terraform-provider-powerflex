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
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHProvisionerConfig ssh provisioner config
type SSHProvisionerConfig struct {
	IP         string
	Username   string
	Password   *string
	PrivateKey *string
	CaCert     *string
	HostKey    *string
}

// getSSHConfig returns ssh config
func (config *SSHProvisionerConfig) getSSHConfig() (*ssh.ClientConfig, error) {
	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		HostKeyCallback: ssh.FixedHostKey(nil),
	}

	// password or private key
	if config.PrivateKey != nil {
		// if private key is specified, use it
		privateKey, err := decodeString(*config.PrivateKey)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
		if config.CaCert != nil {
			// if CA cert is specified, use it
			certBytes, err := decodeString(*config.CaCert)
			if err != nil {
				return nil, err
			}
			pk, err := ssh.ParsePublicKey(certBytes)
			if err != nil {
				return nil, err
			}
			cert := pk.(*ssh.Certificate)
			signer, err = ssh.NewCertSigner(cert, signer)
			if err != nil {
				return nil, err
			}
		}
		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if config.Password != nil {
		// if password is specified, use it
		pwd := *config.Password
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.Password(pwd),
			ssh.KeyboardInteractive(PasswordOnlyKIC(pwd)),
		}
	} else {
		return nil, fmt.Errorf("password or private key must be specified")
	}

	// use fixed host key if provided
	if config.HostKey != nil {
		hostKey, err := decodeString(*config.HostKey)
		if err != nil {
			return nil, err
		}
		hostKeyPub, _, _, _, err := ssh.ParseAuthorizedKey(hostKey)
		if err != nil {
			return nil, err
		}
		sshConfig.HostKeyCallback = ssh.FixedHostKey(hostKeyPub)
	}
	sshConfig.SetDefaults()
	return sshConfig, nil
}

// Logger is an interface for logging
type Logger interface {
	Printf(string, ...any)
	Println(...any)
}

// SSHProvisioner ssh provisioner struct
type SSHProvisioner struct {
	sshClient *ssh.Client
	logger    Logger

	// in case we need to reconnect
	config *ssh.ClientConfig
	ip     string
}

// Close closes ssh connection
func (p *SSHProvisioner) Close() error {
	return p.sshClient.Close()
}

// Run runs ssh command
func (p *SSHProvisioner) Run(cmd string) (string, error) {
	p.logger.Printf("Running command: %s", cmd)
	session, err := p.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("failed to run command: %w", err)
	}
	return string(output), nil
}

// RunWithDir runs ssh command with directory
func (p *SSHProvisioner) RunWithDir(dir, cmd string) (string, error) {
	return p.Run(fmt.Sprintf("cd %s && %s", dir, cmd))
}

// RebootUnix reboots host
func (p *SSHProvisioner) RebootUnix() error {
	cmd := "reboot"
	p.logger.Printf("Running command: %s", cmd)
	session, err := p.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	err = session.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to run reboot command: %w", err)
	}
	p.logger.Printf("Reboot Started")
	time.Sleep(10 * time.Second)
	_ = session.Close()
	p.logger.Printf("Waiting for host IP to be available")
	err = p.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping host IP after reboot: %w", err)
	}
	time.Sleep(10 * time.Second)

	p.logger.Printf("Connecting to %s via ssh", p.ip)
	client, err := ssh.Dial("tcp", net.JoinHostPort(p.ip, "22"), p.config)
	if err != nil {
		return fmt.Errorf("failed to dial remote host: %w", err)
	}
	p.logger.Println("Connected")
	p.sshClient = client
	return nil
}

// GetLinesUnix get lines from string
func GetLinesUnix(op string) []string {
	lines := lineBreakRegex.Split(strings.TrimSpace(op), -1)
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines
}

// UntarUnix untars file
func (p *SSHProvisioner) UntarUnix(filename, dir string) ([]string, error) {
	op, err := p.Run(fmt.Sprintf("cd %s && tar -xvf %s", dir, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to untar file: %w: %s", err, op)
	}
	p.logger.Printf("Untar output: %s", op)
	lines := GetLinesUnix(op)
	return lines, nil
}

// ListDirUnix lists directory
func (p *SSHProvisioner) ListDirUnix(dir string, logOp bool) ([]string, error) {
	op, err := p.Run(fmt.Sprintf("ls %s", dir))
	if err != nil {
		return nil, fmt.Errorf("failed to run list directory command: %w: %s", err, op)
	}
	if logOp {
		p.logger.Printf("List Directory output: %s", op)
	}
	lines := GetLinesUnix(op)
	return lines, nil
}

// Ping pings host
func (p *SSHProvisioner) Ping() error {
	hostIP := p.ip
	start := time.Now()
	for time.Since(start) < 10*time.Minute {
		p.logger.Printf("Checkinging for host IP to be available...")
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(hostIP, "22"), 5*time.Second)
		if err == nil {
			p.logger.Printf("Host IP is available.\n")
			if err := conn.Close(); err != nil {
				p.logger.Println("[WARN] Failed to close TCP connection", err.Error())
			}
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("failed to reach host IP %s within timeout", hostIP)
}

// NewSSHProvisioner creates new ssh provisioner
func NewSSHProvisioner(config SSHProvisionerConfig, logger Logger) (*SSHProvisioner, error) {
	if logger == nil {
		logger = log.Default()
	}
	logger.Printf("Parsing configuration")
	sshConfig, err := config.getSSHConfig()
	if err != nil {
		return nil, fmt.Errorf("error parsing ssh configuration: %w", err)
	}
	logger.Printf("Connecting to %s", config.IP)
	client, err := ssh.Dial("tcp", net.JoinHostPort(config.IP, "22"), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial remote host: %w", err)
	}
	logger.Println("Connected")
	return &SSHProvisioner{
		sshClient: client,
		logger:    logger,
		config:    sshConfig,
		ip:        config.IP,
	}, nil
}

// PasswordOnlyKIC - An ssh.KeyboardInteractiveChallenge that returns the password for every question
func PasswordOnlyKIC(password string) ssh.KeyboardInteractiveChallenge {
	return func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		// Just send the password back for all questions
		answers := make([]string, len(questions))
		for i := range questions {
			answers[i] = password
		}

		return answers, nil
	}
}

// GetLines get lines from string
func GetLines(s string) []string {
	return strings.Split(s, "\n")
}
