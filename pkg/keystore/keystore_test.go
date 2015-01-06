// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The keystore tests require opengpg keys from the keystoretest package (keystoretest.KeyMap).
// The opengpg keys are auto generated by running the keygen.go command.
// keygen.go should not be run by an automated process. keygen.go is a helper to generate
// the keystoretest/keymap.go source file.
//
// If additional opengpg keys are need for testing, please use the following process:
//   * add a new key name to keygen.go
//   * cd keystore/keystoretest
//   * go run keygen.go
//   * check in the results
package keystore

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/coreos/rocket/Godeps/_workspace/src/golang.org/x/crypto/openpgp/errors"
	"github.com/coreos/rocket/pkg/keystore/keystoretest"
)

const tstprefix = "keystore-test"

func testKeyStoreConfig(dir string) (*Config, error) {
	c := &Config{
		RootPath:         path.Join(dir, "/etc/rkt/trustedkeys/root.d"),
		SystemRootPath:   path.Join(dir, "/usr/lib/rkt/trustedkeys/root.d"),
		PrefixPath:       path.Join(dir, "/etc/rkt/trustedkeys/prefix.d"),
		SystemPrefixPath: path.Join(dir, "/usr/lib/rkt/trustedkeys/prefix.d"),
	}
	for _, path := range []string{c.RootPath, c.SystemRootPath, c.PrefixPath, c.SystemPrefixPath} {
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func TestStoreTrustedKey(t *testing.T) {
	dir, err := ioutil.TempDir("", tstprefix)
	if err != nil {
		t.Fatalf("error creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	armoredPublicKey := keystoretest.KeyMap["example.com"].ArmoredPublicKey
	fingerprint := keystoretest.KeyMap["example.com"].Fingerprint

	keyStoreConfig, err := testKeyStoreConfig(dir)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	ks := New(keyStoreConfig)

	output, err := ks.StoreTrustedKeyPrefix("example.com/foo", bytes.NewBufferString(armoredPublicKey))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if filepath.Base(output) != fingerprint {
		t.Errorf("expected finger print %s, got %v", fingerprint, filepath.Base(output))
	}
	if err := ks.DeleteTrustedKeyPrefix("example.com/foo", fingerprint); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Errorf("unexpected error %v", err)
	}

	output, err = ks.MaskTrustedKeySystemPrefix("example.com/foo", fingerprint)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	fi, err := os.Lstat(output)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if fi.Size() != 0 {
		t.Errorf("expected empty file")
	}

	output, err = ks.StoreTrustedKeyRoot(bytes.NewBufferString(armoredPublicKey))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if filepath.Base(output) != fingerprint {
		t.Errorf("expected finger print %s, got %v", fingerprint, filepath.Base(output))
	}
	if err := ks.DeleteTrustedKeyRoot(fingerprint); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Errorf("unexpected error %v", err)
	}

	output, err = ks.MaskTrustedKeySystemRoot(fingerprint)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	fi, err = os.Lstat(output)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if fi.Size() != 0 {
		t.Errorf("expected empty file")
	}
}

func TestCheckSignature(t *testing.T) {
	trustedPrefixKeys := []string{
		"example.com/app",
		"acme.com/services",
		"acme.com/services/web/nginx",
	}
	trustedRootKeys := []string{
		"coreos.com",
	}
	trustedSystemRootKeys := []string{
		"acme.com",
	}

	dir, err := ioutil.TempDir("", tstprefix)
	if err != nil {
		t.Fatalf("error creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	keyStoreConfig, err := testKeyStoreConfig(dir)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	ks := New(keyStoreConfig)
	for _, key := range trustedPrefixKeys {
		if _, err := ks.StoreTrustedKeyPrefix(key, bytes.NewBufferString(keystoretest.KeyMap[key].ArmoredPublicKey)); err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}
	for _, key := range trustedRootKeys {
		if _, err := ks.StoreTrustedKeyRoot(bytes.NewBufferString(keystoretest.KeyMap[key].ArmoredPublicKey)); err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}
	for _, key := range trustedSystemRootKeys {
		dst := filepath.Join(ks.SystemRootPath, keystoretest.KeyMap[key].Fingerprint)
		if err := ioutil.WriteFile(dst, []byte(keystoretest.KeyMap[key].ArmoredPublicKey), 0644); err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	if _, err := ks.MaskTrustedKeySystemRoot(keystoretest.KeyMap["acme.com"].Fingerprint); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	checkSignatureTests := []struct {
		name    string
		key     string
		trusted bool
	}{
		{"coreos.com/etcd", "coreos.com", true},
		{"coreos.com/fleet", "coreos.com", true},
		{"coreos.com/flannel", "coreos.com", true},
		{"example.com/app", "example.com/app", true},
		{"acme.com/services/web/nginx", "acme.com/services/web/nginx", true},
		{"acme.com/services/web/auth", "acme.com/services", true},
		{"acme.com/etcd", "acme.com", false},
		{"acme.com/web/nginx", "acme.com", false},
		{"acme.com/services/web", "acme.com/services/web/nginx", false},
	}
	for _, tt := range checkSignatureTests {
		key := keystoretest.KeyMap[tt.key]
		message, signature, err := keystoretest.NewMessageAndSignature(key.ArmoredPrivateKey)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
			continue
		}
		signer, err := ks.CheckSignature(tt.name, message, signature)
		if tt.trusted {
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			if signer.PrimaryKey.KeyIdString() != key.Fingerprint {
				t.Errorf("expected fingerprint == %v, got %v", key.Fingerprint, signer.PrimaryKey.KeyIdString())
			}
			continue
		}
		if err == nil {
			t.Errorf("expected ErrUnknownIssuer error")
			continue
		}
		if err.Error() != errors.ErrUnknownIssuer.Error() {
			t.Errorf("unexpected error %v", err)
		}
	}
}
