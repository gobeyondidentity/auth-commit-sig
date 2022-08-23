package action

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// VerifyCommitSignature accepts a commit object and a base64-encoded PGP public
// key. Parses the key into a temporary key ring, then checks that the signature
// attached to the commit is valid.
func VerifyCommitSignature(base64Key string, commit *object.Commit) error {
	payload, err := EncodedCommitWithoutSignature(commit)
	if err != nil {
		return fmt.Errorf("failed to encode commit: %w", err)
	}

	err = CheckSignatureByKey(base64Key, commit.PGPSignature, payload)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// verifyCommitByEmail accepts a committer email address and a list of allowed email
// addresses that can bypass commit signature validation. Returns true if the email
// is found on the allowlist; otherwise returns false.
func verifyCommitByEmailAddress(committerEmailAddress string, allowlistEmailAddresses []string) bool {
	want := strings.ToLower(committerEmailAddress)
	for _, email := range allowlistEmailAddresses {
		if want == strings.ToLower(email) {
			return true
		}
	}

	return false
}

// verifyCommitSignatureByThirdPartyKeys accepts a commit object and a list of armored
// PGP public keys. Parses the keys one by one into a temporary key ring, then checks
// that the signature attached to the commit is valid. Returns true if the signature is
// validated by a key within the list.
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
func verifyCommitSignatureByThirdPartyKeys(publicKeys []string, commit *object.Commit) (bool, error) {
	payload, err := EncodedCommitWithoutSignature(commit)
	if err != nil {
		return false, fmt.Errorf("failed to encode commit: %w", err)
	}

	for _, key := range publicKeys {
		pubKeyBase64 := base64.StdEncoding.EncodeToString([]byte(key))
		keyRing, err := openpgp.ReadArmoredKeyRing(base64.NewDecoder(base64.StdEncoding, strings.NewReader(pubKeyBase64)))
		if err != nil {
			log.Printf("Failed to parse public key: \n\n%s\n\nwith error: %v", key, err)
		} else {
			_, err = openpgp.CheckArmoredDetachedSignature(keyRing, strings.NewReader(payload), strings.NewReader(commit.PGPSignature), nil)
			if err == nil {
				log.Printf("Commit validated by public key: \n\n%s\n\n", key)
				return true, nil
			}
		}
	}

	return false, nil
}
