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

Additionally, an allowlist can be configured for

1. select committer email addresses to bypass signature verification
2. select third party keys used for signature verification

See section on [Allowlist](#allowlist).

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
        id: signature-verification
        uses: gobeyondidentity/auth-commit-sig@v1
        with:
          # API Token to call Beyond Identity cloud.
          # The token should be stored as an actions secret or fetched from a secrets manager.
          # https://docs.github.com/en/actions/security-guides/encrypted-secrets.
          api_token: ${{ secrets.BYNDID_KEY_MGMT_API_TOKEN }}
          # Repository which the signature verification action is performed on.
          # This is also used to match against repositories listed on the allowlist.
          repository: "gobeyondidentity/auth-commit-sig"
        # Echo outcome JSON to standard out.
        # Outcome can be used in other steps as ${{ steps.signature-verification.outputs.outcome }}.
      - name: Print outcome
        run: |
          echo '${{ steps.signature-verification.outputs.outcome }}'
    # Output of job. Can be accessed in other jobs as ${{ needs.auth-commit-sig.outputs.outcome }}.
    outputs:
      outcome: ${{ steps.signature-verification.outputs.outcome }}
```

The `BYNDID_KEY_MGMT_API_TOKEN` should be set as an
[actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
in the repository or fetched from a secrets manager.

## Allowlist

An allowlist can be configured for the github action to pass. The configuration contains the
following two lists:

1. Email addresses and the repositories the email address can be used to bypass signature verification.
2. Third party keys and the repositories that the third party key can be used for signature verification.

If the email address of the committer is on the allowlist, the action will bypass signature
verification. Otherwise, it will continue with the regular signature verification process.

If there are third party keys on the allowlist, the action will attempt to verify
the signature using all those keys. If verification succeeds, the action will
pass. If verification fails, the action continues with the regular signature verification
process.

The `repositories` parameter attached to the email addresses and third party keys is optional. If not provided,
the email address or third party key will be used on ALL repositories the action is run on.

### Actions Workflow

Some additional steps added to the action workflow.

1. In a repository, create an allowlist YAML file. The repository _can_ be the same repository
   where the action is registered, or it can be a separate repository. We recommend it be stored in
   seperate repo that has strict access controls.

   An example allowlist file called `allowlist_example.yaml` can be found within this repository and below.

### Example Allowlist YAML file - allowlist_example.yaml

```yaml
email_addresses:
  # `repositories` not defined, can bypass signature verification for _any_ repository.
  - email_address: user1@company.com

  # `repositories` defined, can bypass signature verification _only_ for repository_A and repository_B.
  - email_address: user2@company.com
    repositories:
      - repository_A
      - repository_B

third_party_keys:
  # `repositories` not defined, can be used for signature verification for _any_ repository.
  - key: |
      -----BEGIN PGP PUBLIC KEY BLOCK-----

      m3IEYr73lBMIKoZIzj01AQcCAwT1gCXHMjKP6EWgtzJNxpkfFWhpK4dsV1dfbzRz
      ...truncated

      -----END PGP PUBLIC KEY BLOCK-----

  # `repositories` defined, can be used for signature verification _only_ for repository_C and repository_D.
  - key: |
      -----BEGIN PGP PUBLIC KEY BLOCK-----

      t3IEYr738sG78d7SD01AQcCAwT1gCXHMjKP6EWgtzJNxpkfF01AQcCAwT11dfbzRx
      ...truncated

      -----END PGP PUBLIC KEY BLOCK-----
    repositories:
      - repository_C
      - repository_D
```

2. Add an additional `Checkout Allowlist repository` step the workflow configuration. This step
   checks out the repo where the allowlist YAML file is stored and stores the files from that repo locally
   in the action workspace under `path`. If your allowlist is in the same repository where the action
   is registered, this step can be skipped.

   If the allowlist repository is private, you should set up an SSH Key for cloning.
   See [Deploy Keys](https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys).

   This SSH Key should be set as an [actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
   called `ALLOWLIST_REPO_SSH_PRIV_KEY` in the repository where the action is run or fetched from a secrets manager.

3. Add additional variables to the `Authorize with Beyond Identity` step.

- `allowlist_config_file_path`: File path to the allowlist configuration file from actions workspace.
  - If it is from a different repository, this variable should be prefixed with `./` and the `path` from Step 2.
    - Ex. `./{path}/{path_to_allowlist}`
  - If it is the same repository as the action, this variable should be prefiexed with `./` and the path to the allowlist configuration file.
    - Ex. `./{path_to_allowlist}`

### Actions Workflow with Allowlist

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
        # configured in a separate repository. See README for details on how to set
        # up the allowlist configuration.
      - name: Checkout Allowlist repository
        uses: actions/checkout@v3
        with:
          # Path to allowlist repository in the format (owner/repo_name).
          # If it is the same repository as above, can skip this step.
          repository: "gobeyondidentity/commit-signing-allowlist"
          # SSH Key for private repositories. If repo is public, can omit.
          # See https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys.
          # The key should be stored as an actions secret or fetched from a secrets manager.
          # See https://docs.github.com/en/actions/security-guides/encrypted-secrets.
          ssh-key: ${{ secrets.ALLOWLIST_REPO_SSH_PRIV_KEY }}
          # Relative path under $GITHUB_WORKSPACE to place the repository.
          # NOTE: the repo_name will not be cloned as a folder. Only the files will
          # be checked out at this path.
          path: "allowlist-dir"
      - name: Authorize with Beyond Identity
        id: signature-verification
        uses: gobeyondidentity/auth-commit-sig@v1
        env:
          # If API_BASE_URL is omitted, defaults to our production API server
          # at https://api.byndid.com/key-mgmt.
          API_BASE_URL: "https://api.byndid.com/key-mgmt"
        with:
          # API Token to call Beyond Identity cloud.
          # The token should be stored as an actions secret or fetched from a secrets manager.
          # https://docs.github.com/en/actions/security-guides/encrypted-secrets.
          api_token: ${{ secrets.BYNDID_KEY_MGMT_API_TOKEN }}
          # Repository which the signature verification action is performed on.
          # This is also used to match against repositories listed on the allowlist.
          repository: "gobeyondidentity/auth-commit-sig"
          # If allowlist is configured, set this to the file path to the allowlist
          # configuration file starting from $GITHUB_WORKSPACE.
          # If it is from a different repository, this variable should be prefixed
          # with `./` and the `path` from Step 2.
          # Ex. `./{path}/{path_to_allowlist}`
          # If it is the same repository as the action, it is just the path to the
          # allowlist configuration file.
          # Ex. `./{path_to_allowlist}`
          # See README for details on how to format this file.
          allowlist_config_file_path: "./allowlist-dir/allowlist.yaml"
        # Echo outcome JSON to standard out.
        # Outcome can be used in other steps as ${{ steps.signature-verification.outputs.outcome }}.
      - name: Print outcome
        run: |
          echo '${{ steps.signature-verification.outputs.outcome }}'
    # Output of job. Can be accessed in other jobs as ${{ needs.auth-commit-sig.outputs.outcome }}.
    outputs:
      outcome: ${{ steps.signature-verification.outputs.outcome }}
```

## Outcome output

When the action is complete, it returns an output that is a JSON blob containing information about the
results of the action. It contains info about the commit, if signature verification was successful,
if an allowlist email address or third party key was used, errors, and other details.

The outcome is returned as an output of the step and can be accessed within the job as `${{ steps.signature-verification.outputs.outcome }}`.

The outcome is returned as an output of the job and can be used in other jobs as `${{ needs.auth-commit-sig.outputs.outcome }}`

See [Defining outputs for jobs](https://docs.github.com/en/actions/using-jobs/defining-outputs-for-jobs) for more
information.

### Example Outcomes

#### Passed with `BI_MANAGED_KEY`

```json
{
  "version": "1.0.0",
  "repository": "gobeyondidentity/auth-commit-sig",
  "commit": {
    "commit_hash": "7f1927513d02f546ed50579de7192c181836ea23",
    "tree_hash": "d742502776eb874cba1605b768ad438810ecf911",
    "parent_hashes": ["12cf190ced71c3d23107d90641b361b9edcc41f6"],
    "author": {
      "name": "John Doe",
      "email_address": "john@doe.com",
      "timestamp": "2022-09-05T13:58:12-04:00"
    },
    "committer": {
      "name": "John Doe",
      "email_address": "john@doe.com",
      "timestamp": "2022-09-05T13:58:12-04:00"
    },
    "signed": true,
    "signature_key_id": "87A2691085B4544E"
  },
  "result": "PASS",
  "desc": "Signature verified by Beyond Identity.",
  "verification_details": {
    "verified_by": "BI_MANAGED_KEY",
    "bi_managed_key": {
      "key_id": "87A2691085B4544E",
      "email_address": "john@doe.com"
    }
  },
  "errors": []
}
```

#### Passed with `THIRD_PARTY_KEY`

```json
{
  "version": "1.0.0",
  "repository": "gobeyondidentity/auth-commit-sig",
  "commit": {
    "commit_hash": "7f1927513d02f546ed50579de7192c181836ea23",
    "tree_hash": "d742502776eb874cba1605b768ad438810ecf911",
    "parent_hashes": ["12cf190ced71c3d23107d90641b361b9edcc41f6"],
    "author": {
      "name": "Jane Doe",
      "email_address": "jane@doe.com",
      "timestamp": "2022-09-05T13:58:12-04:00"
    },
    "committer": {
      "name": "Jane Doe",
      "email_address": "jane@doe.com",
      "timestamp": "2022-09-05T13:58:12-04:00"
    },
    "signed": true,
    "signature_key_id": "7723AD85B1221B3B"
  },
  "result": "PASS",
  "desc": "Signature verified by third party key on the allowlist.",
  "verification_details": {
    "verified_by": "THIRD_PARTY_KEY",
    "third_party_key": {
      "key_id": "7723AD85B1221B3B",
      "fingerprint": "VOcYYhfmV98gAMVoh6JpEIW0VE4=",
      "user_id": "Jane Doe <jane@doe.com>"
    }
  },
  "errors": []
}
```

#### Error

```json
{
  "version": "1.0.0",
  "repository": "gobeyondidentity/auth-commit-sig",
  "commit": {
    "commit_hash": "7f1927513d02f546ed50579de7192c181836ea23",
    "tree_hash": "d742502776eb874cba1605b768ad438810ecf911",
    "parent_hashes": ["12cf190ced71c3d23107d90641b361b9edcc41f6"],
    "author": {
      "name": "Jackie Doe",
      "email_address": "jackie@doe.com",
      "timestamp": "2022-09-05T13:58:12-04:00"
    },
    "committer": {
      "name": "Jackie Doe",
      "email_address": "jackie@doe.com",
      "timestamp": "2022-09-05T13:58:12-04:00"
    },
    "signed": false
  },
  "result": "FAIL",
  "desc": "Commit is not signed. See errors for details.",
  "errors": [
    {
      "desc": "commit is not signed"
    }
  ]
}
```
