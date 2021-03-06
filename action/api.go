package action

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var (
	// version will be replaced at build time by a string like "1.2.3"
	version = "unknown"
	// userAgent will be sent with all API requests
	userAgent = "byndid/auth-commit-sig:" + version
)

// APIClient wraps an http.Client to provide access to the Beyond Identity Key
// Management API.
type APIClient struct {
	HTTPClient *http.Client
	APIToken   string
	APIBaseURL string
}

// BadResponseError is returned when an unexpected response is received from the
// API.
type BadResponseError struct {
	RequestMethod string
	RequestURL    *url.URL
	StatusCode    int
	Body          []byte
	Header        http.Header
	Cause         error
}

func (err BadResponseError) Error() string {
	return fmt.Sprintf("bad response from %s %s: %d %s: (body: %s): %v",
		err.RequestMethod, err.RequestURL, err.StatusCode, http.StatusText(err.StatusCode), string(err.Body), err.Cause)
}

// Authorization is returned by a successful GetAuthorization API call.
type Authorization struct {
	Authorized bool   `json:"authorized"`
	Message    string `json:"message"`
	GPGKey     GPGKey `json:"gpg_key"`
}

// PrettyPrint returns a nicely formatted JSON representation of the
// Authorization response.
func (a Authorization) PrettyPrint() string {
	bs, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bs)
}

//	GPGKey contains the GPG public key.
type GPGKey struct {
	// ID is the Beyond Identity ID of the GPG key.
	ID string `json:"id"`
	// Base64Key is the binary GPG "transferable public key message".
	Base64Key string `json:"base64_key"`
}

// GetAuthorization calls the Beyond Identity Key Management API to authorize a
// GPG key for git commit signing.
func (c APIClient) GetAuthorization(ctx context.Context, keyID, committerEmail string) (*Authorization, error) {
	u, err := url.Parse(c.APIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	u.Path = path.Join(u.Path, "v0", "gpg", "key", "authorization", "git-commit-signing")

	q := u.Query()
	q.Set("key_id", keyID)
	q.Set("committer_email", committerEmail)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+c.APIToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read api response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, BadResponseError{
			RequestMethod: req.Method,
			RequestURL:    req.URL,
			StatusCode:    resp.StatusCode,
			Body:          body,
			Header:        resp.Header,
			Cause:         fmt.Errorf("expected status %d", http.StatusOK),
		}
	}

	a := Authorization{}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return nil, BadResponseError{
			RequestMethod: req.Method,
			RequestURL:    req.URL,
			StatusCode:    resp.StatusCode,
			Body:          body,
			Header:        resp.Header,
			Cause:         err,
		}
	}

	return &a, nil
}
