package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
// configuration repository. If repoName and filePath are not empty, returns an
// Allowlist struct containing empty lists for email addresses and third party keys.
func GetAllowlist(repoName, filePath string) (*Allowlist, error) {
	fullPath := path.Join(allowlistPrefix, repoName, filePath)
	log.Println("fullpath: " + fullPath)
	dir, _ := os.Getwd()
	log.Println("dir: " + dir + "/allowlist")
	files, err := ioutil.ReadDir(dir + "/allowlist")
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	log.Println("dir: " + "./allowlist")
	files, err = ioutil.ReadDir("./allowlist")
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	log.Println(dir + "/" + fullPath)
	yfile, err := ioutil.ReadFile(dir + "/" + fullPath)
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
