package action

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// Config configures a run of the action.
type Config struct {
	// RepoPath is a path to a directory containing a clone of the git repository.
	// Required.
	RepoPath string
	// CommitRef is the commit reference to verify (e.g. "HEAD" or
	// "cf2d2127c69c57bef0232b553146c418e1cba43a").
	// Required.
	CommitRef string
	// APIToken is used as a Bearer token for the Beyond Identity Key Management
	// API.
	// Required.
	APIToken string
	// APIBaseURL is the base URL of the Beyond Identity Key Management API.
	// Required.
	APIBaseURL string
	// AllowlistConfigRepoName is a path to the directory containing a clone of the
	// git repository containing the allowlist configuration, if configured.
	// Both AllowlistConfigRepoName and AllowlistConfigFilePath must be set if
	// allowlist is configured.
	AllowlistConfigRepoName string
	// AllowlistConfigFilePath is a path to the file containing the allowlist
	// configuration, if configured.
	// Both AllowlistConfigRepoName and AllowlistConfigFilePath must be set if
	// allowlist is configured.
	AllowlistConfigFilePath string
}

// MissingConfigFieldError is returned from Config.Validate() if any required
// fields are missing.
type MissingConfigFieldError string

func (e MissingConfigFieldError) Error() string {
	return fmt.Sprintf("missing config field: %s", string(e))
}

// Validate checks that the Config is valid.
func (c Config) Validate() error {
	if c.RepoPath == "" {
		return MissingConfigFieldError("RepoPath")
	}
	if c.CommitRef == "" {
		return MissingConfigFieldError("CommitRef")
	}
	if c.APIToken == "" {
		return MissingConfigFieldError("APIToken")
	}
	if c.APIBaseURL == "" {
		return MissingConfigFieldError("APIBaseURL")
	}
	if c.AllowlistConfigRepoName == "" && c.AllowlistConfigFilePath != "" {
		return MissingConfigFieldError("AllowlistConfigRepoName")
	}
	if c.AllowlistConfigRepoName != "" && c.AllowlistConfigFilePath == "" {
		return MissingConfigFieldError("AllowlistConfigFilePath")
	}
	return nil
}

// Run the action. Returns nil if the commit was properly signed by a GPG key
// authorized for the committer.
func Run(ctx context.Context, cfg Config) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Verifying commit with ref %q in %q", cfg.CommitRef, cfg.RepoPath)

	commit, err := GetCommit(cfg.RepoPath, cfg.CommitRef)
	if err != nil {
		return fmt.Errorf("failed to get commit: %w", err)
	}

	log.Printf("\nCommit:\n================\n%s\n================\n\n", PrettyPrintCommit(commit))

	committerEmail := commit.Committer.Email

	// Fetch allowlist.
	allowlist, err := GetAllowlist(cfg.AllowlistConfigRepoName, cfg.AllowlistConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to get allowlist: %w", err)
	}

	// If allowlist email addresses exist, verify commit by email address.
	if len(allowlist.EmailAddresses) > 0 {
		log.Printf("Checking email address allowlist for bypassing signature verification\n")
		onAllowlistEmails := verifyCommitByEmailAddress(committerEmail, allowlist.EmailAddresses)
		if onAllowlistEmails {
			log.Printf("Committer email: \"%s\" is on email address allowlist, bypassing signature verification\n\n", committerEmail)
			return nil
		}
		log.Printf("Committer email: \"%s\" is not on email address allowlist, continuing signature verification\n\n", committerEmail)
	}

	// If allowlist third party keys exists, verify signature throuh keys.
	if len(allowlist.ThirdPartyKeys) > 0 {
		log.Printf("Verifying commit signature with third party keys\n\n")
		pass, err := verifyCommitSignatureByThirdPartyKeys(allowlist.ThirdPartyKeys, commit)
		if err != nil {
			log.Printf("Failed to verify commit by third party keys: %v\n\n", err)
		}
		if pass {
			log.Printf("Commit is signed by authorized third party key\n\n")
			return nil
		}
		log.Printf("No third party keys validated signature, continuing signature verification\n\n")
	}

	// Signature verification.
	if commit.PGPSignature == "" {
		return errors.New("commit is not signed")
	}

	issuerKeyID, err := ParseSignatureIssuerKeyID(commit.PGPSignature)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	log.Printf("Getting authorization for GPG key %q with committer email %q", issuerKeyID, committerEmail)

	authorization, err := APIClient{
		HTTPClient: http.DefaultClient,
		APIToken:   cfg.APIToken,
		APIBaseURL: cfg.APIBaseURL,
	}.GetAuthorization(ctx, issuerKeyID, committerEmail)
	if err != nil {
		return fmt.Errorf("failed to get authorization: %w", err)
	}

	log.Printf("\nAPI response:\n================\n%s\n================\n\n", authorization.PrettyPrint())

	err = Verify(commit, authorization)
	if err != nil {
		return fmt.Errorf("failed to verify commit with authorization: %w", err)
	}

	log.Println("Commit is signed by an authorized Beyond Identity user")
	return nil
}

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
