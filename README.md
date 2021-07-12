# auth-commit-sig

GitHub Action for authorizing a Beyond Identity user to sign git commits and
verifying those signatures. The action enforces that all commits are signed by
authorized users in your organization's Beyond Identity directory.

## Usage

The recommended usage is to add a workflow to your repository that runs on all
pull requests targeting the default branch (e.g. `main`).

```
name: Verify Commit Signature

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
      - name: Verify with Beyond Identity
        uses: byndid/auth-commit-sig@v0
        with:
          api_token: ${{ secrets.BYNDID_KEY_MGMT_API_TOKEN }}
```

The `BYNDID_KEY_MGMT_API_TOKEN` should be set as a secret in the repository.


## Development

### Running the action locally

Set environment variables:

```
export API_TOKEN=<redacted>
export API_BASE_URL=https://api.rolling.byndid.run/key-mgmt
```

Build the image:

```
docker build . -t byndid/auth-commit-sig
```

Run the image:

```
docker run --rm --workdir /workspace -v $(pwd):/workspace \
    -e API_TOKEN -e API_BASE_URL \
    byndid/auth-commit-sig -ref=HEAD
```

`HEAD` can be replaced with any commit reference.
