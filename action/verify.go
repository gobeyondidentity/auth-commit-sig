package action

import (
	"fmt"

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
