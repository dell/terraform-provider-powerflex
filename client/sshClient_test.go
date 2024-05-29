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
	// srv, err := getServer()
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }
	// defer srv.Close()

	// t.Log("started ssh server")

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
		t.Fatalf(err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	op, err := sshP.ListDirUnix("/etc/testDir", true)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(op)
	assert.Equal(t, []string{"sbc", "lkm", "por"}, op)
}

func TestSshClientMReboot(t *testing.T) {

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
		t.Fatalf(err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	err = sshP.RebootUnix()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSshClientMUntar(t *testing.T) {

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
		t.Fatalf(err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	op, err := sshP.UntarUnix("testTarFile.tar", "/etc/testTarDir")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(op)
	assert.Equal(t, []string{"tfFile1", "tfFile2", "tfFile3"}, op)
}

func TestSshClientMScp(t *testing.T) {
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
		t.Fatalf(err.Error())
	}
	defer sshP.Close()

	t.Log("created ssh client")

	// upload sw
	scpProv := NewScpProvisioner(sshP)
	err = scpProv.Upload("/root/terraform-provider-powerflex/client/testFile.txt", "/tmp/testScpFile.txt", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	// read /tmp/testScpFile
	conts, err := os.ReadFile("/tmp/testScpFile")
	if err != nil {
		t.Fatalf(err.Error())
	}
	// convert conts to string
	contsStr := string(conts)
	assert.Equal(t, "Hello World!!!", contsStr)
}

func TestSshClientMWrongPass(t *testing.T) {
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
