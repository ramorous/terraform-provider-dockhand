# Terraform Provider Dockhand Registry Publishing Checklist

This file tracks the status and requirements for publishing the Terraform Provider for Dockhand to the Terraform Registry.

## Pre-Publishing Requirements

- [x] Provider code complete and tested
- [x] All 7 resources implemented
- [x] All 6 data sources implemented  
- [x] Provider binary successfully builds
- [x] Go module properly named: `github.com/ramorous/terraform-provider-dockhand`
- [x] Provider properly namespaced: `registry.terraform.io/ramorous/dockhand`
- [x] GitHub repository public and named `terraform-provider-dockhand`
- [x] Repository owner is the organization: `ramorous`

## Documentation Requirements

- [x] README.md with comprehensive provider documentation
- [x] QUICKSTART.md with getting started guide
- [x] TESTING.md with testing procedures
- [x] CONTRIBUTING.md with contribution guidelines
- [x] CHANGELOG.md with version history
- [x] Examples directory with sample configurations
- [x] Inline code documentation for resources/data sources
- [x] RELEASING.md with release process documentation

## Build and Release Setup

- [x] .goreleaser.yaml - Release configuration
- [x] GitHub Actions workflows:
  - [x] .github/workflows/release.yml - Automated releases
  - [x] .github/workflows/test.yml - Testing on PRs
  - [x] .github/workflows/registry-publish.yml - Documentation publishing
- [x] GNUmakefile with build targets
- [x] go.mod with proper module path
- [x] MIT License file

## GPG Signing Setup (Required for Registry)

- [x] GPG key created and exported
- [x] GPG_PRIVATE_KEY secret added to GitHub repository
- [x] GPG_PASSPHRASE secret added to GitHub repository
- [x] .goreleaser.yaml configured for signing

## Step-by-Step Publishing Instructions

### 1. Setup GPG Signing (One-time)

```bash
# Generate GPG key if you don't have one
gpg --gen-key

# List your keys
gpg --list-secret-keys

# Export your private key (replace KEY_ID with your key ID)
gpg --armor --export-secret-key KEY_ID > private-key.asc

# Get your fingerprint
gpg --list-secret-keys --keyid-format LONG
# Copy the fingerprint (typically 40 characters)
```

### 2. Add GitHub Secrets

1. Go to: https://github.com/ramorous/terraform-provider-dockhand/settings/secrets/actions
2. Add new repository secrets:
   - Name: `GPG_PRIVATE_KEY`
   - Value: [Contents of private-key.asc]
   - Name: `GPG_PASSPHRASE`
   - Value: [Your GPG key passphrase]

### 3. Create Initial Release

```bash
# Update version in code if needed
# Update CHANGELOG.md
# Commit changes
git add .
git commit -m "docs: prepare release v0.1.1"

# Create and push version tag
git tag -a v0.1.1 -m "Release version 0.1.1"
git push origin v0.1.1

# GitHub Actions will automatically:
# - Build binaries for all platforms
# - Generate checksums
# - Sign checksums with GPG
# - Create GitHub release with artifacts
```

### 4. Verify Release

Check: https://github.com/ramorous/terraform-provider-dockhand/releases

Verify:
- [ ] terraform-provider-dockhand_0.1.0_linux_amd64.zip
- [ ] terraform-provider-dockhand_0.1.0_linux_arm64.zip
- [ ] terraform-provider-dockhand_0.1.0_darwin_amd64.zip
- [ ] terraform-provider-dockhand_0.1.0_darwin_arm64.zip
- [ ] terraform-provider-dockhand_0.1.0_windows_amd64.zip
- [ ] terraform-provider-dockhand_0.1.0_SHA256SUMS file
- [ ] terraform-provider-dockhand_0.1.0_SHA256SUMS.sig signed file

### 5. Publish to Terraform Registry

1. Visit: https://registry.terraform.io
2. Sign in or create account
3. Click "Publish" → "Provider"
4. Select "GitHub" and authorize
5. Select "ramorous/terraform-provider-dockhand"
6. Click "Publish provider"

### 6. Verify Registry Publication

After 5-10 minutes:
- [ ] Provider appears at: https://registry.terraform.io/providers/ramorous/dockhand
- [ ] Documentation is accessible
- [ ] Versions are listed in registry
- [ ] Checksums are verified

## Post-Publication

- [ ] Announce release in GitHub discussion/issues
- [ ] Update website documentation if applicable
- [ ] Monitor registry for any issues
- [ ] Continue maintaining and releasing updates

## Registry Requirements

The Terraform Registry requires:

1. **Naming Convention**: `terraform-provider-{name}`
2. **Public Repository**: Must be publicly accessible on GitHub
3. **Proper Namespacing**: organization/provider-name format
4. **Version Tags**: Must use v{version} format
5. **Signed Releases**: Checksums must be GPG-signed
6. **Documentation**: Generated from code comments
7. **Valid License**: MIT is valid and acceptable

## Useful Links

- Terraform Registry: https://registry.terraform.io
- Provider Publishing: https://registry.terraform.io/docs/providers/publish-provider
- GPG Signing: https://docs.github.com/en/github/authenticating-to-github/managing-commit-signature-verification
- GoReleaser: https://goreleaser.com/
- Terraform Plugin Framework: https://developer.hashicorp.com/terraform/plugin/framework

## Troubleshooting

### Release Workflow Fails

Check GitHub Actions logs at:
https://github.com/ramorous/terraform-provider-dockhand/actions

Common issues:
- GPG key not properly exported
- GPG_PRIVATE_KEY secret not configured
- Incorrect key fingerprint in .goreleaser.yaml

### Registry Publication Fails

Verify:
- Release has checksum files (both .txt and .sig)
- Signatures are valid
- Binary naming matches format
- Repository is public
- GitHub Actions workflows complete successfully

## Contact

For questions about publishing:
- Terraform Registry Docs: https://registry.terraform.io/docs
- HashiCorp Support: support@hashicorp.com
