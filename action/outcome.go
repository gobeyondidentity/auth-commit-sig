package action

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// Outcome represents the outcome of the action.
type Outcome struct {
	Version             string               `json:"version"`
	Repository          string               `json:"repository"`
	Commit              *Commit              `json:"commit,omitempty"`
	Result              string               `json:"result"`
	Desc                string               `json:"desc"`
	VerificationDetails *VerificationDetails `json:"verification_details,omitempty"`
	Errors              []OutcomeError       `json:"errors"`
}

// Commit contains information about a commit.
type Commit struct {
	CommitHash     string   `json:"commit_hash"`
	TreeHash       string   `json:"tree_hash"`
	ParentHashes   []string `json:"parent_hashes"`
	Author         *Actor   `json:"author"`
	Committer      *Actor   `json:"committer"`
	Signed         bool     `json:"signed"`
	SignatureKeyID string   `json:"signature_key_id,omitempty"`
}

// Actor represents a commit actor.
type Actor struct {
	Name         string    `json:"name"`
	EmailAddress string    `json:"email_address"`
	Timestamp    time.Time `json:"timestamp"`
}

// VerificationDetails contains information about how the commit
// signature was verified.
type VerificationDetails struct {
	VerifiedBy    string         `json:"verified_by"`
	EmailAddress  string         `json:"email_address,omitempty"`
	ThirdPartyKey *ThirdPartyKey `json:"third_party_key,omitempty"`
	BIManagedKey  *BIManagedKey  `json:"bi_managed_key,omitempty"`
}

// ThirdPartyKey represents a third party key that was used to
// sign a commit.
type ThirdPartyKey struct {
	KeyID       string `json:"key_id"`
	Fingerprint string `json:"fingerprint"`
	UserID      string `json:"user_id"`
}

// BIManagedKey represents a Beyond Identity managed key that was
// used to sign a commit.
type BIManagedKey struct {
	KeyID        string `json:"key_id"`
	EmailAddress string `json:"email_address"`
}

// OutcomeError represents an error that occurred during the action.
type OutcomeError struct {
	Desc string `json:"desc"`
}

// SetResultAndDescription sets the result and description of an Outcome.
func (o *Outcome) SetResultAndDescription(result, desc string) {
	o.Result = result
	o.Desc = desc
}

// SetVerificationDetailsEmailAddress sets the verification details with
// a commit verified by an allowlist email address.
func (o *Outcome) SetVerificationDetailsEmailAddress(emailAddress string) {
	o.VerificationDetails = &VerificationDetails{
		VerifiedBy:   "EMAIL_ADDRESS",
		EmailAddress: emailAddress,
	}
}

// SetVerificationDetailsThirdPartyKey sets the verification details with
// a commit signed by a third party key.
func (o *Outcome) SetVerificationDetailsThirdPartyKey(tpk *ThirdPartyKey) {
	o.VerificationDetails = &VerificationDetails{
		VerifiedBy:    "THIRD_PARTY_KEY",
		ThirdPartyKey: tpk,
	}
}

// SetVerificationDetailsBIManagedKey sets the verification details with
// a commit signed by a Beyond Identity managed key.
func (o *Outcome) SetVerificationDetailsBIManagedKey(keyID, emailAddress string) {
	o.VerificationDetails = &VerificationDetails{
		VerifiedBy: "BI_MANAGED_KEY",
		BIManagedKey: &BIManagedKey{
			KeyID:        keyID,
			EmailAddress: emailAddress,
		},
	}
}

// SetCommit sets the Commit field within the Outcome.
func (o *Outcome) SetCommit(c *object.Commit) {
	pHashes := []string{}
	for _, ph := range c.ParentHashes {
		if ph.String() != "" {
			pHashes = append(pHashes, ph.String())
		}
	}

	o.Commit = &Commit{
		CommitHash:   c.Hash.String(),
		TreeHash:     c.TreeHash.String(),
		ParentHashes: pHashes,
		Author: &Actor{
			Name:         c.Author.Name,
			EmailAddress: c.Author.Email,
			Timestamp:    c.Author.When,
		},
		Committer: &Actor{
			Name:         c.Committer.Name,
			EmailAddress: c.Committer.Email,
			Timestamp:    c.Committer.When,
		},
		Signed: len(c.PGPSignature) > 0,
	}
}

// NewOutcomeError converts an error into an OutcomeError.
func NewOutcomeError(err error) OutcomeError {
	return OutcomeError{Desc: err.Error()}
}

// SetErrors sets errors on the OutcomeError.
func (o *Outcome) SetErrors(errs ...error) {
	for _, err := range errs {
		o.Errors = append(o.Errors, NewOutcomeError(err))
	}
}
