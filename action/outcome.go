package action

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type Outcome struct {
	Version             string               `json:"version"`
	Repository          string               `json:"repository"`
	Commit              *Commit              `json:"commit,omitempty"`
	Result              string               `json:"result"`
	Desc                string               `json:"desc"`
	VerificationDetails *VerificationDetails `json:"verification_details,omitempty"`
	Errors              []OutcomeError       `json:"errors"`
}

type Commit struct {
	CommitHash     string `json:"commit_hash"`
	TreeHash       string `json:"tree_hash"`
	ParentHash     string `json:"parent_hash"`
	Author         *Actor `json:"author"`
	Committer      *Actor `json:"committer"`
	Signed         bool   `json:"signed"`
	SignatureKeyID string `json:"signature_key_id,omitempty"`
}

type Actor struct {
	Name         string    `json:"name"`
	EmailAddress string    `json:"email_address"`
	Timestamp    time.Time `json:"timestamp"`
}

type VerificationDetails struct {
	VerifiedBy    string         `json:"verified_by"`
	EmailAddress  string         `json:"email_address,omitempty"`
	ThirdPartyKey *ThirdPartyKey `json:"third_party_key,omitempty"`
	BIManagedKey  *BIManagedKey  `json:"bim_managed_key,omitempty"`
}

type ThirdPartyKey struct {
	KeyID       string `json:"key_id"`
	Fingerprint string `json:"fingerprint"`
	UserID      string `json:"user_id"`
}

type BIManagedKey struct {
	KeyID        string `json:"key_id"`
	EmailAddress string `json:"email_address"`
}

type OutcomeError struct {
	Desc string `json:"desc"`
}

func (o *Outcome) SetResultAndDescription(result, desc string) {
	o.Result = result
	o.Desc = desc
}

func (o *Outcome) SetVerificationDetailsEmailAddress(emailAddress string) {
	o.VerificationDetails = &VerificationDetails{
		VerifiedBy:   "EMAIL_ADDRESS",
		EmailAddress: emailAddress,
	}
}

func (o *Outcome) SetVerificationDetailsThirdPartyKey(tpk *ThirdPartyKey) {
	o.VerificationDetails = &VerificationDetails{
		VerifiedBy:    "THIRD_PARTY_KEY",
		ThirdPartyKey: tpk,
	}
}

func (o *Outcome) SetVerificationDetailsBIManagedKey(keyID, emailAddress string) {
	o.VerificationDetails = &VerificationDetails{
		VerifiedBy: "BI_MANAGED_KEY",
		BIManagedKey: &BIManagedKey{
			KeyID:        keyID,
			EmailAddress: emailAddress,
		},
	}
}

func (o *Outcome) SetCommit(c *object.Commit) {
	o.Commit = &Commit{
		CommitHash: c.Hash.String(),
		TreeHash:   c.TreeHash.String(),
		ParentHash: c.ParentHashes[0].String(),
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

func NewOutcomeError(err error) OutcomeError {
	return OutcomeError{Desc: err.Error()}
}

func (o *Outcome) SetErrors(errs ...error) {
	for _, err := range errs {
		o.Errors = append(o.Errors, NewOutcomeError(err))
	}
}
