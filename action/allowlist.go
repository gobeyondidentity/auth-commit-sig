package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"gopkg.in/yaml.v3"
)

// Allowlist is the struct containing two lists:
//
// 1. (EmailAddresses) Email addresses and the repositories that the
// email address can be used to bypass signature verification.
// 2. (ThirdPartyKeys) Third party keys and the repositories that the
// third party key can be used for signature verification.
//
// If the list of repositories is empty, the email address or third
// party key can be used for ALL repositories.
type Allowlist struct {
	// EmailAddresses is the list of EmailAddressEntries.
	EmailAddresses []EmailAddressEntry `yaml:"email_addresses"`
	// ThirdPartyKeys is the list of ThirdPartyKeyEntries.
	ThirdPartyKeys []ThirdPartyKeyEntry `yaml:"third_party_keys"`
}

// EmailAddressEntry is a struct containing an email address and a list of
// repositories for which the email address can bypass signature verification.
// If the list of repositories is empty, the email address can bypass all
// repositories.
type EmailAddressEntry struct {
	EmailAddress string   `yaml:"email_address"`
	Repositories []string `yaml:"repositories"`
}

// ThirdPartyKeyEntry is a struct containing a third party key and a list of
// repositories for which the third party key can be used for signature verification.
// If the list of repositories is empty, the third party key can be used for signature
// verification on all repositories.
type ThirdPartyKeyEntry struct {
	Key          string   `yaml:"key"`
	Repositories []string `yaml:"repositories"`
}

// LoadAllowlist verifies and parses the allowlist configuration from the allowlist
// file path. If filePath is empty, returns an Allowlist struct containing empty
// lists for email addresses and third party keys.
func LoadAllowlist(filePath string) (*Allowlist, error) {
	if filePath == "" {
		log.Println("No allowlist configured")
		return &Allowlist{}, nil
	}
	yfile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf(`failed to read allowlist yaml configuration file at '%s': %w`, filePath, err)
	}

	var allowlist *Allowlist
	err = yaml.Unmarshal(yfile, &allowlist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal allowlist yaml configuration file: %w", err)
	}

	return allowlist, nil
}

// RepoAllowlist is the struct containing the validated email addresses
// and third party keys from the Allowlist struct for the specified repository.
type RepoAllowlist struct {
	// EmailAddresses is the list of validated emails addresses allowed to
	// bypass commit signature verification for the specified repository.
	EmailAddresses []string
	// ThirdPartyKeys is an array of keyrings used to validate a PGP
	// signature for the specified repository.
	ThirdPartyKeys []openpgp.EntityList
}

// GetAllowlistForRepo parses the allowlist for valid email addresses and
// third party keys from the Allowlist struct for the specified repository.
// Returns any errors encountered while parsing.
func GetAllowlistForRepo(al *Allowlist, repo string) (*RepoAllowlist, []error) {
	emails, eaErrs := getValidEmailAddressesForRepo(al.EmailAddresses, repo)
	keyRings, tpkErrs := getValidThirdPartyKeysForRepo(al.ThirdPartyKeys, repo)

	repoAllowlist := &RepoAllowlist{
		EmailAddresses: emails,
		ThirdPartyKeys: keyRings,
	}

	return repoAllowlist, append(eaErrs, tpkErrs...)
}

// getValidEmailAddressesForRepo parses an array of EmailAddressEntries and returns a list
// of valid allowlist email addresses for the specified repository.
// Returns any errors encountered while parsing.
func getValidEmailAddressesForRepo(entries []EmailAddressEntry, repo string) ([]string, []error) {
	emails := []string{}
	errs := []error{}
	for _, e := range entries {
		emailAddress := e.EmailAddress
		if err := Email(e.EmailAddress); err != nil {
			errs = append(errs, err)
		} else {
			if len(e.Repositories) == 0 || containsRepo(repo, e.Repositories) {
				emails = append(emails, emailAddress)
			}
		}
	}

	return emails, errs
}

// getValidEmailAddressesForRepo parses an array of ThirdPartyKeyEntries and returns a list
// of keyRings used for PGP signature validation.
// Returns any errors encountered while parsing.
func getValidThirdPartyKeysForRepo(entries []ThirdPartyKeyEntry, repo string) ([]openpgp.EntityList, []error) {
	keyRings := []openpgp.EntityList{}
	errs := []error{}
	for _, e := range entries {
		keyRing, err := openpgp.ReadArmoredKeyRing(strings.NewReader(e.Key))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse third party key: %s\n with error: %v", e.Key, err))
		} else {
			if len(e.Repositories) == 0 || containsRepo(repo, e.Repositories) {
				keyRings = append(keyRings, keyRing)
			}
		}
	}
	return keyRings, errs
}

// containsRepo checks if the specified repository is within an array
// of repositories.
func containsRepo(repo string, repos []string) bool {
	for _, r := range repos {
		if r == repo {
			return true
		}
	}
	return false
}
