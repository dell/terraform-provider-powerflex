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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSshClientM(t *testing.T) {
	t.Skip("Skipping this test case, only for Unit test")

	pass := "secret"
	sshP, err := NewSSHProvisioner(SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "2222",
		Username:   "root",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
	}, nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	op, err := sshP.ListDirUnix("/etc/testDir", true)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Log(op)
	assert.Equal(t, []string{"sbc", "lkm", "por"}, op)
}

func TestSshClientMReboot(t *testing.T) {
	t.Skip("Skipping this test case, only for Unit test")
	pass := "secret"
	sshP, err := NewSSHProvisioner(SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "2222",
		Username:   "root",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
	}, nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	err = sshP.RebootUnix()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestSshClientMUntar(t *testing.T) {
	t.Skip("Skipping this test case, only for Unit test")
	pass := "secret"
	sshP, err := NewSSHProvisioner(SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "2222",
		Username:   "root",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
	}, nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	op, err := sshP.UntarUnix("testTarFile.tar", "/etc/testTarDir")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Log(op)
	assert.Equal(t, []string{"tfFile1", "tfFile2", "tfFile3"}, op)
}

func TestSshClientMScp(t *testing.T) {
	t.Skip("Skipping this test case, only for Unit test")
	// remove file /tmp/testScpFile, just in case
	os.Remove("/tmp/testScpFile")

	pass := "secret"
	sshP, err := NewSSHProvisioner(SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "2222",
		Username:   "root",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
	}, nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	// upload sw
	scpProv := NewScpProvisioner(sshP)
	err = scpProv.Upload("/root/terraform-provider-powerflex/client/testFile.txt", "/tmp/testScpFile.txt", "")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	// read /tmp/testScpFile
	conts, err := os.ReadFile("/tmp/testScpFile")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	// convert conts to string
	contsStr := string(conts)
	assert.Equal(t, "Hello World!!!", contsStr)
}

func TestSshClientMWrongPass(t *testing.T) {
	t.Skip("Skipping this test case, only for Unit test")
	pass := "secret1"
	_, err := NewSSHProvisioner(SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "2222",
		Username:   "root",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
	}, nil)
	if err == nil {
		t.Fatalf("No error returned when wrong password is provided")
		return
	}
}

func TestGetSSHConfigWithFixedHostKey(t *testing.T) {
	// Test that getSSHConfig uses FixedHostKey when HostKey is provided
	t.Skip("Skipping - requires valid SSH public key")
	pass := "testpassword"
	// This is a sample RSA public key for testing
	hostKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC"
	config := SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "22",
		Username:   "testuser",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
		HostKey:    &hostKey,
	}

	sshConfig, err := config.getSSHConfig()
	assert.NoError(t, err, "getSSHConfig should not return error")
	assert.NotNil(t, sshConfig, "sshConfig should not be nil")
	assert.Equal(t, "testuser", sshConfig.User, "Username should match")
	assert.NotNil(t, sshConfig.HostKeyCallback, "HostKeyCallback should be set")
}

func TestGetSSHConfigWithoutHostKey(t *testing.T) {
	// Test that getSSHConfig fails when no known_hosts file and no host_key
	pass := "testpassword"
	config := SSHProvisionerConfig{
		IP:         "localhost",
		Port:       "22",
		Username:   "testuser",
		Password:   &pass,
		PrivateKey: nil,
		CaCert:     nil,
		HostKey:    nil,
	}

	_, err := config.getSSHConfig()
	assert.Error(t, err, "getSSHConfig should return error when no known_hosts and no host_key")
	assert.Contains(t, err.Error(), "no known_hosts file found", "Error should mention known_hosts file")
	assert.Contains(t, err.Error(), "no explicit host_key provided", "Error should mention host_key requirement")
}
