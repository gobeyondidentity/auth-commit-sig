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

## Usage

The recommended usage is to add a workflow to your repository that runs on all
pull requests targeting the default branch (e.g. `main`).

```
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
      - uses: actions/checkout@v2
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
          api_token: ${{ secrets.BYNDID_KEY_MGMT_API_TOKEN }}
```

The `BYNDID_KEY_MGMT_API_TOKEN` should be set as a secret in the repository.
