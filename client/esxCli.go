package client

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var lineBreakRegex = regexp.MustCompile("\r?\n")

// EsxCli - esxcli client
type EsxCli struct {
	client *SshProvisioner
}

// NewEsxCli - creates new esxcli client
func NewEsxCli(prov *SshProvisioner) *EsxCli {
	return &EsxCli{
		client: prov,
	}
}

// EsxCliSw - esxcli software model
type EsxCliSw struct {
	Name             string
	Version          string
	Vendor           string
	Type             string
	InstallationDate string
}

// NewEsxCliSw - creates new esxcli software model from esxcli output
func NewEsxCliSw(txt string) (*EsxCliSw, error) {
	sw := EsxCliSw{}
	fields := strings.Fields(txt)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid software list item :%s", txt)
	}
	sw.Name = fields[0]
	sw.Version = fields[1]
	sw.Vendor = fields[2]
	sw.Type = fields[3]
	sw.InstallationDate = fields[4]
	return &sw, nil
}

// SoftwareList - returns list of installed software
func (e *EsxCli) SoftwareList() ([]*EsxCliSw, error) {
	op, err := e.client.Run("esxcli software vib list")
	if err != nil {
		return nil, err
	}
	ret := make([]*EsxCliSw, 0)
	installedSw := lineBreakRegex.Split(op, -1)
	// the first two lines contain formatting for humans
	// ignore them
	installedSw = installedSw[2:]
	for _, sw := range installedSw {
		sw = strings.TrimSpace(sw)
		if sw == "" {
			// ignore empty lines
			continue
		}
		det, merr := NewEsxCliSw(sw)
		if merr != nil {
			err = errors.Join(err, merr)
			continue
		}
		ret = append(ret, det)
	}
	return ret, err
}

// GetSoftwareByNameRegex - returns software by name using regex patterns
func (e *EsxCli) GetSoftwareByNameRegex(name *regexp.Regexp) (*EsxCliSw, error) {
	sw, err := e.SoftwareList()
	if err != nil {
		return nil, err
	}
	ind := slices.IndexFunc(sw, func(s *EsxCliSw) bool {
		return name.MatchString(s.Name)
	})
	if ind == -1 {
		return nil, fmt.Errorf("software not found")
	}
	return sw[ind], nil
}

// VinbInstallCommand - esxcli software vib install command model
type VibInstallCommand struct {
	ZipFile  string
	SigCheck bool
}

// SoftwareInstall - installs software on the esxi host
func (e *EsxCli) SoftwareInstall(vib VibInstallCommand) (string, error) {
	command := fmt.Sprintf("esxcli software vib install -d %s", vib.ZipFile)
	if !vib.SigCheck {
		command = fmt.Sprintf("%s --no-sig-check", command)
	}
	return e.client.Run(command)
}

// SetModuleParameters - sets module parameters on esxi host
func (e *EsxCli) SetModuleParameters(module string, params map[string]string) (string, error) {
	lparams := make([]string, 0)
	for k, v := range params {
		lparams = append(lparams, fmt.Sprintf("%s=%s", k, v))
	}
	sparams := strings.Join(lparams, " ")

	return e.client.Run(fmt.Sprintf(`esxcli system module parameters set -m %s -p "%s"`, module, sparams))
}

// SoftwareRmv - removes software on the esxi host
func (e *EsxCli) SoftwareRmv(name string) (string, error) {
	return e.client.Run(fmt.Sprintf("esxcli software vib remove -n %s", name))
}
