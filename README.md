<p align="center">
   <br/>
   <a href="https://developers.beyondidentity.com" target="_blank"><img src="https://user-images.githubusercontent.com/238738/178780350-489309c5-8fae-4121-a20b-562e8025c0ee.png" width="150px" ></a>
   <h3 align="center">Beyond Identity</h3>
   <p align="center">Universal Passkeys for Developers</p>
   <p align="center">
   All devices. Any protocol. Zero shared secrets. 
   </p>
</p>

# auth-commit-sig

GitHub Action for authorizing a Beyond Identity user to sign git commits and
verifying those signatures. The action enforces that all commits are signed by
authorized users in your organization's Beyond Identity directory.

Additionally, an allowlist can be configured for select users to bypass signature
verification and also for verification by third party keys. See section on
[Allowlist](#allowlist).

## Usage

The recommended usage is to add a workflow to your repository that runs on all
pull requests targeting the default branch (e.g. `main`).

### Actions Workflow

```yaml
name: Authorize Commit Signing

on:
  pull_request:
    # Optionally limit to only pull requests targeting the default branch.
    # Remove this to run on all pull requests.
    branches: [main]

jobs:
  auth-commit-sig:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          # Critical: check out the head commit on the branch. By default,
          # actions/checkout will check out a merge commit built for the pull
          # request and signed by Github itself. Using the pull-request HEAD
          # allows the action to check the latest commit on the pull request,
          # which must be signed by an authorized user before it can be merged.
          ref: ${{ github.event.pull_request.head.sha }}
      - name: Authorize with Beyond Identity
        uses: gobeyondidentity/auth-commit-sig@v0
        with:
          # API Token to call Beyond Identity cloud.
          # The token should be stored as an actions secret.
          # https://docs.github.com/en/actions/security-guides/encrypted-secrets.
          api_token: ${{ secrets.BYNDID_KEY_MGMT_API_TOKEN }}
```

The `BYNDID_KEY_MGMT_API_TOKEN` should be set as an 
[actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets) 
in the repository.

## Allowlist

An allowlist can be configured for the github action to pass. The configuration contains the
following two lists:

1. Email addresses allowed to bypass signature verification.
2. Third party keys used for signature verification.

If the email address of the committer is on the allowlist, the action will bypass signature
verification. Otherwise, it will continue with the regular signature verification process.

If there are third party keys on the allowlist, the action will attempt to verify
the signature using all the configured keys. If verification succeeds, the action will
pass. If verification fails, the action continues with the regular signature verification
process.

### Actions Workflow

Some additional steps added to the action workflow.

1. In a repository, create an allowlist YAML file. The repository _can_ be the same repository
   where the action is registered, or it can be a separate repository.

### Example Allowlist YAML file - allowlist.yaml

```yaml
email_addresses:
  - user1@company.com
  - user2@company.com
third_party_keys:
  - |
    -----BEGIN PGP PUBLIC KEY BLOCK-----

    m3IEYr73lBMIKoZIzj01AQcCAwT1gCXHMjKP6EWgtzJNxpkfFWhpK4dsV1dfbzRz
    ...truncated

    -----END PGP PUBLIC KEY BLOCK-----
  - |
    -----BEGIN PGP PUBLIC KEY BLOCK-----

    t3IEYr738sG78d7SD01AQcCAwT1gCXHMjKP6EWgtzJNxpkfF01AQcCAwT11dfbzRx
    ...truncated

    -----END PGP PUBLIC KEY BLOCK-----
```

2. Add the additional `Checkout Allowlist repository` step the workflow configuration. This step
   checks out repo where the allowlist YAML file is stored.

   If the repository is private, should set up an SSH Key for cloning. 
   See [Deploy Keys](https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys).

   This action step will also incorporate the private key as `ALLOWLIST_REPO_SSH_PRIV_KEY`. 
   This should be set as an [actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets) 
   in the repository.

3. Add additional variables to the `Authorize with Beyond Identity` step.

- `allowlist_config_file_path`: File path to the allowlist configuration file from root of the allowlist repository.

### Example Workflow Configuration

```yaml
name: Authorize Commit Signing

on:
  pull_request:
    branches: [main]

jobs:
  auth-commit-sig:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout commit
        uses: actions/checkout@v3
        with:
          # Critical: check out the head commit on the branch. By default,
          # actions/checkout will check out a merge commit built for the pull
          # request and signed by Github itself. Using the pull-request HEAD
          # allows the action to check the latest commit on the pull request,
          # which must be signed by an authorized user before it can be merged.
          ref: ${{ github.event.pull_request.head.sha }}
        # If allowlist is configured, checkout allowlist repository. Otherwise, do
        # not include this step.
        # The allowlist _can_ be in the same repository as the step above, or can be
        # configured in a seperate repository. See README for details on how to set
        # up the allowlist configuration.
      - name: Checkout Allowlist repository
        uses: actions/checkout@v3
        with:
          # Path to allowlist repository in the format (owner/repo_name).
          repository: "gobeyondidentity/commit-signing-allowlist"
          # SSH Key for private repositories. If repo is public, can omit.
          # See https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys.
          # The key should be stored as an actions secret.
          # See https://docs.github.com/en/actions/security-guides/encrypted-secrets.
          ssh-key: ${{ secrets.ALLOWLIST_REPO_SSH_PRIV_KEY }}
          # Do not modify. Our action uses this path to fetch the allowlist file.
          path: "allowlist"
      - name: Authorize with Beyond Identity
        uses: gobeyondidentity/auth-commit-sig@v0
        env:
          # If API_BASE_URL is omitted, defaults to our production API server
          # at https://api.byndid.com/key-mgmt.
          API_BASE_URL: "https://api.byndid.com/key-mgmt"
        with:
          # API Token to call Beyond Identity cloud.
          # The token should be stored as an actions secret.
          # https://docs.github.com/en/actions/security-guides/encrypted-secrets.
          api_token: ${{ secrets.BYNDID_KEY_MGMT_API_TOKEN }}
          # If allowlist is configured, set this to the file path to the allowlist
          # configuration file starting from the root of the allowlist repository.
          # See README for details on how to format this file.
          allowlist_config_file_path: "allowlist.yaml"
```
