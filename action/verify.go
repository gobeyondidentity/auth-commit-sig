package action

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Verify verifies a commit.
func Verify(commit *object.Commit, authorization *Authorization) error {
	if !authorization.Authorized {
		return fmt.Errorf("authorization denied: %s", authorization.Message)
	}

	err := VerifyCommitSignature(authorization.GPGKey.Base64Key, commit)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

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

// verifyCommitSignatureByThirdPartyKeys accepts a commit object and a list of
// keyRings. Returns true if the signature attached to the commit object is
// validated by a key within the list; otherwise returns false.
func verifyCommitSignatureByThirdPartyKeys(keyRings []openpgp.EntityList, payload string, commit *object.Commit) (*ThirdPartyKey, bool) {
	for _, key := range keyRings {
		signer, err := openpgp.CheckArmoredDetachedSignature(key, strings.NewReader(payload), strings.NewReader(commit.PGPSignature), nil)
		if err == nil {
			keyID := fmt.Sprintf("%X", signer.PrimaryKey.KeyId)
			fp := base64.StdEncoding.EncodeToString(signer.PrimaryKey.Fingerprint)
			userID := signer.PrimaryIdentity().Name
			log.Printf("Signature made using RSA key %s\nand fingerprint %s\nfrom %s\n\n", keyID, fp, userID)
			return &ThirdPartyKey{
				KeyID:       keyID,
				Fingerprint: fp,
				UserID:      userID,
			}, true
		}
	}
	return nil, false
}
