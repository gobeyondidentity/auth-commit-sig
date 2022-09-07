package action

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	// PASS is the string representation of "PASS".
	PASS = "PASS"
	// FAIL is the string representation of "FAIL".
	FAIL = "FAIL"
)

// Run the action. Returns an Outcome that caputure's the results of the action.
//
// The action can pass in one of the following ways:
//
// 1. Bypassing signature verification through an email address on the allowlist (if configured).
// 2. Properly signed by a third party key on the allowlist (if configured).
// 3. Properly signed by a GPG key authorized for the committer.
func Run(ctx context.Context, cfg Config) *Outcome {
	o := &Outcome{Version: "1.0.0", Repository: cfg.Repository, Errors: []OutcomeError{}}
	errs := cfg.Validate()
	if len(errs) > 0 {
		o.SetErrors(errs...)
		o.SetResultAndDescription(FAIL, "Invalid config. See errors for details.")
		return o
	}

	log.Printf("Verifying commit with ref %q in %q", cfg.CommitRef, cfg.RepoPath)

	commit, err := GetCommit(cfg.RepoPath, cfg.CommitRef)
	if err != nil {
		o.SetErrors(err)
		o.SetResultAndDescription(FAIL, "Failed to get commit. See errors for details.")
		return o
	}
	o.SetCommit(commit)

	log.Printf("\nCommit:\n================\n%s\n================\n\n", PrettyPrintCommit(commit))

	// Load the allowlist.
	allowlist, err := LoadAllowlist(cfg.AllowlistConfigFilePath)
	if err != nil {
		o.SetErrors(err)
		o.SetResultAndDescription(FAIL, "Failed to load the allowlist. See errors for details.")
		return o
	}

	// Parse out valid allowlist email addresses and keys for the specified
	// repository. Adds any parsing errors to the outcome but does not return
	// at this step.
	repoAllowlist, errs := GetAllowlistForRepo(allowlist, cfg.Repository)
	if len(errs) > 0 {
		o.SetErrors(errs...)
	}

	committerEmail := commit.Committer.Email

	// If the repo allowlist contains email addresses, attempt to bypass signature verification
	// through the email address.
	if len(repoAllowlist.EmailAddresses) > 0 {
		log.Printf("Checking signature verification bypass with email address from the allowlist\n\n")
		verified := verifyCommitByEmailAddress(committerEmail, repoAllowlist.EmailAddresses)
		if verified {
			log.Printf("Committer email: \"%s\" is on email address allowlist, bypassing signature verification\n\n", committerEmail)
			o.SetVerificationDetailsEmailAddress(committerEmail)
			o.SetResultAndDescription(PASS, "Bypassed signature verification with an email address from the allowlist.")
			return o
		}
		log.Printf("Committer email: \"%s\" is not on email address allowlist, continuing signature verification\n\n", committerEmail)
	}

	// Validate that a signature exists for third party key validation and BI cloud verification.
	if commit.PGPSignature == "" {
		o.SetErrors(errors.New("commit is not signed"))
		o.SetResultAndDescription(FAIL, "Commit is not signed. See errors for details.")
		return o
	}

	issuerKeyID, err := ParseSignatureIssuerKeyID(commit.PGPSignature)
	if err != nil {
		o.SetErrors(err)
		o.SetResultAndDescription(FAIL, "Failed to parse signature. See errors for details.")
		return o
	}
	o.Commit.SignatureKeyID = issuerKeyID

	payload, err := EncodedCommitWithoutSignature(commit)
	if err != nil {
		o.SetErrors(err)
		o.SetResultAndDescription(FAIL, "Failed to encode commit. See errors for details.")
		return o
	}

	// If the repo allowlist contains third party keys, attempt to verify the signature through the keys.
	if len(repoAllowlist.ThirdPartyKeys) > 0 {
		log.Printf("Verifying commit signature with third party keys from the allowlist\n\n")
		tpk, pass := verifyCommitSignatureByThirdPartyKeys(repoAllowlist.ThirdPartyKeys, payload, commit)
		if pass {
			log.Printf("Commit is signed by authorized third party key\n\n")
			o.SetVerificationDetailsThirdPartyKey(tpk)
			o.SetResultAndDescription(PASS, "Signature verified by a third party key from the allowlist.")
			return o
		}
		log.Printf("No third party keys validated signature, continuing signature verification\n\n")
	}

	log.Printf("Getting authorization for GPG key %q with committer email address %q\n\n", issuerKeyID, committerEmail)

	// Attempt to verify signature through BI cloud.
	authorization, err := APIClient{
		HTTPClient: http.DefaultClient,
		APIToken:   cfg.APIToken,
		APIBaseURL: cfg.APIBaseURL,
	}.GetAuthorization(ctx, issuerKeyID, committerEmail)
	if err != nil {
		o.SetErrors(fmt.Errorf("failed to get authorization to BI cloud: %w", err))
		o.SetResultAndDescription(FAIL, "Failed to get authorization to BI cloud. See errors for details.")
		return o
	}

	log.Printf("\nAPI response:\n================\n%s\n================\n\n", authorization.PrettyPrint())

	err = Verify(commit, authorization)
	if err != nil {
		o.SetErrors(fmt.Errorf("failed to verify commit with authorization: %w", err))
		o.SetResultAndDescription(FAIL, "Failed to verify commit. See errors for details.")
		return o
	}

	log.Println("Commit is signed by an authorized Beyond Identity user")
	o.SetVerificationDetailsBIManagedKey(issuerKeyID, committerEmail)
	o.SetResultAndDescription(PASS, "Signature verified by Beyond Identity.")
	return o
}
