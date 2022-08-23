package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/yaml.v3"
)

// Allowlist is the struct containing two lists:
//
// 1. Email addresses allowed to bypass signature verification.
// 2. Third party keys used for signature verification.
type Allowlist struct {
	// EmailAddresses is the list of emails addresses allowed to bypass commit
	// signature verification.
	EmailAddresses []string `yaml:"email_addresses"`
	// ThirdPartyKeys is the list of third party public keys used for signature
	// verification. It is a string array containing armored PGP keys.
	//
	// Armored simply means the format of the key is base64 encoded data, alongside
	// a plaintext header + footer:
	//
	// -----BEGIN PGP PUBLIC KEY BLOCK-----
	//
	// 9T6cSwE9PGVUwxYRFvrOVfEdtW2rGpQf46blrSRtTrc=
	// ...truncated
	//
	// -----END PGP PUBLIC KEY BLOCK-----
	ThirdPartyKeys []string `yaml:"third_party_keys"`
}

var allowlistPrefix = "allowlist"

// GetAllowlist verifies and parses the allowlist configuration from the allowlist
// file path. If filePath is empty, returns an Allowlist struct containing empty
// lists for email addresses and third party keys.
func GetAllowlist(filePath string) (*Allowlist, error) {
	if filePath == "" {
		log.Println("No allowlist configured")
		return &Allowlist{}, nil
	}
	fullPath := path.Join(allowlistPrefix, filePath)
	yfile, err := ioutil.ReadFile("./" + fullPath)
	if err != nil {
		return nil, fmt.Errorf(`failed to read allowlist yaml configuration file at '%s': %w`, fullPath, err)
	}

	var allowlist *Allowlist
	err = yaml.Unmarshal(yfile, &allowlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal allowlist yaml configuration file: %w", err)
	}

	return allowlist, nil
}
