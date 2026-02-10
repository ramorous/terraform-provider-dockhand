# Publishing Terraform Provider Dockhand to the Registry

This guide provides complete instructions for publishing the Terraform Provider for Dockhand to the official Terraform Registry.

## What This Provider Offers

The Terraform Provider for Dockhand allows you to manage Docker infrastructure through Terraform, including:
- Docker containers
- Docker Compose stacks
- Docker networks and volumes
- Docker images
- Environments

## Prerequisites Checklist

Before publishing, ensure you have:

- [x] Provider code complete and tested
- [x] All resources and data sources documented
- [x] Public GitHub repository named `terraform-provider-dockhand`
- [x] Repository organization/namespace: `ramorous`
- [x] MIT license included
- [x] Proper git tags in `v{version}` format
- [x] GitHub Actions workflows configured for CI/CD
- [x] GoReleaser configuration for automated builds
- [x] `go.mod` with correct module path: `github.com/ramorous/terraform-provider-dockhand`
- [x] Provider address in code: `registry.terraform.io/ramorous/dockhand`

## Files Created for Registry Publishing

### Build and Release Configuration

1. **.goreleaser.yaml** - Automated build configuration
   - Builds binaries for Linux, macOS, Windows
   - Supports amd64 and arm64 architectures
   - Generates checksums
   - GPG signs releases
   - Creates GitHub releases automatically

2. **.github/workflows/release.yml** - Release automation
   - Triggered on version tags (v*)
   - Imports GPG key from secrets
   - Runs GoReleaser
   - Creates GitHub release with all artifacts

3. **.github/workflows/test.yml** - Continuous Integration
   - Runs on PRs and pushes to main
   - Runs tests and linting
   - Builds for all supported platforms
   - Ensures code quality before merge

4. **.github/workflows/registry-publish.yml** - Documentation publication
   - Generates documentation for each release
   - Auto-commits generated docs
   - Keeps registry documentation in sync

### Configuration and Automation

5. **GNUmakefile** - Makefile with enhanced targets
   - `make build` - Build provider binary
   - `make install` - Install locally for testing
   - `make test` - Run all tests
   - `make lint` - Run code linter
   - `make fmt` - Format code
   - `make doc` - Generate documentation
   - `make release` - Create releases locally (optional)
   - `make release-snapshot` - Test release without publishing

### Documentation Files

6. **RELEASING.md** - Release process guide
   - Step-by-step release instructions
   - GPG setup procedures
   - GitHub secrets configuration
   - Registry publication steps
   - Troubleshooting guide

7. **DEVELOPMENT.md** - Developer setup guide
   - System requirements
   - Local setup instructions
   - Development workflow
   - Testing procedures
   - IDE configuration

8. **REGISTRY_PUBLISH_CHECKLIST.md** - Publishing checklist
   - Pre-publishing requirements
   - GPG setup steps
   - Detailed release instructions
   - Registry publication guide
   - Post-publication verification

9. **CONTRIBUTING.md** - Updated with detailed guidelines
   - Development setup
   - Branching strategy
   - Code style standards
   - Testing requirements
   - Resource/data source templates

10. **CHANGELOG.md** - Updated with version history
    - Semantic versioning format
    - Lists all resources and data sources
    - Upgrade guidelines
    - Contributing guidelines

11. **.gitignore** - Updated for Go and provider projects
    - Binaries, build outputs
    - IDE files, OS files
    - Terraform and test artifacts
    - Environment files

## Step-by-Step Publishing Process

### Step 1: Prepare GitHub Secrets (One-Time)

You need GPG keys for signing releases:

```bash
# Generate GPG key if needed
gpg --gen-key

# Export private key
gpg --list-secret-keys --keyid-format LONG
gpg --armor --export-secret-key [KEY_ID] > private-key.asc

# Remember your:
# - Key fingerprint (40 chars)
# - Passphrase
```

Add to GitHub repository secrets:
1. Visit: `https://github.com/ramorous/terraform-provider-dockhand/settings/secrets/actions`
2. Add `GPG_PRIVATE_KEY` with contents of private-key.asc
3. Add `GPG_PASSPHRASE` with your GPG passphrase

### Step 2: Create Initial Release

```bash
# Update code if needed
# Update CHANGELOG.md with v0.1.0 changes
# Commit changes
git add CHANGELOG.md
git commit -m "docs: prepare release v0.1.0"

# Create version tag
git tag -a v0.1.0 -m "Release version 0.1.0"

# Push to trigger workflow
git push origin main
git push origin v0.1.0
```

### Step 3: Verify GitHub Release

GitHub Actions will automatically:
1. Build binaries for all platforms (Linux, macOS, Windows)
2. Generate checksums (SHA256SUMS)
3. Sign checksums with GPG
4. Create GitHub release with all artifacts

Check: `https://github.com/ramorous/terraform-provider-dockhand/releases`

Verify these files exist:
- `terraform-provider-dockhand_0.1.0_linux_amd64.zip`
- `terraform-provider-dockhand_0.1.0_linux_arm64.zip`
- `terraform-provider-dockhand_0.1.0_darwin_amd64.zip`
- `terraform-provider-dockhand_0.1.0_darwin_arm64.zip`
- `terraform-provider-dockhand_0.1.0_windows_amd64.zip`
- `terraform-provider-dockhand_0.1.0_SHA256SUMS`
- `terraform-provider-dockhand_0.1.0_SHA256SUMS.sig`

### Step 4: Publish to Terraform Registry

1. Go to: https://registry.terraform.io
2. Sign in (create account if needed)
3. Click "Publish" → "Provider"
4. Select "GitHub" and authorize
5. Select repository: `ramorous/terraform-provider-dockhand`
6. Click "Publish provider"

Registry will verify:
- [ ] Repository is public
- [ ] Binary naming conventions are correct
- [ ] Checksums file is present
- [ ] Checksums are GPG-signed
- [ ] Release is properly tagged

### Step 5: Verify Registry Publication

After 5-10 minutes, your provider will be available at:
```
https://registry.terraform.io/providers/ramorous/dockhand
```

Verify:
- [ ] Provider appears in registry
- [ ] All versions are listed
- [ ] Documentation is generated
- [ ] Resources and data sources are documented
- [ ] Examples are accessible

## Using the Published Provider

Once published, users can use your provider with:

```hcl
terraform {
  required_providers {
    dockhand = {
      source  = "ramorous/dockhand"
      version = "~> 0.1"
    }
  }
}

provider "dockhand" {
  api_url = var.dockhand_api_url
  api_key = var.dockhand_api_key
}

resource "dockhand_container" "example" {
  name  = "my-container"
  image = "nginx:latest"
  # ... more configuration
}
```

## Future Releases

For subsequent releases:

```bash
# Update code and documentation
# Update CHANGELOG.md
git add .
git commit -m "feat: add new feature

Describes the new feature added in this release.
"

# Create new version tag
git tag -a v0.2.0 -m "Release version 0.2.0"

# Push to trigger release workflow
git push origin main
git push origin v0.2.0

# Registry automatically picks up new release within minutes
```

## Troubleshooting

### Release Workflow Fails

**Check GitHub Actions logs**: https://github.com/ramorous/terraform-provider-dockhand/actions

**Common issues**:
- GPG key not properly configured
- GPG_PRIVATE_KEY secret contains incorrect key
- Go version mismatch
- Missing dependencies

**Fix**:
```bash
# Verify GPG key locally
gpg --list-secret-keys
gpg --verify terraform-provider-dockhand_0.1.0_SHA256SUMS.sig

# Rebuild Go dependencies
go mod tidy
go mod download

# Test build locally
make build
```

### Provider Not Appearing in Registry

**Wait 10-15 minutes**: Registry caches updates

**Verify release is complete**:
- [ ] All binaries present
- [ ] Checksums file exists
- [ ] Signature file (.sig) exists
- [ ] Release is not marked as draft
- [ ] Repository is public

**Check registry logs**: https://registry.terraform.io/docs/providers/requirements

### Checksums Verification Fails

Ensure GoReleaser is using correct GPG key:

```bash
# In .goreleaser.yaml, verify:
# GPG_FINGERPRINT matches your key fingerprint

# Test GPG signing locally
gpg --armor --detach-sign --default-key [FINGERPRINT] terraform-provider-dockhand_0.1.0_SHA256SUMS
```

## Best Practices

1. **Follow Semantic Versioning**: MAJOR.MINOR.PATCH
   - MAJOR: Breaking changes
   - MINOR: New features (backward compatible)
   - PATCH: Bug fixes only

2. **Keep CHANGELOG Updated**: Document all changes

3. **Test Before Release**: Run full test suite
   ```bash
   make fmt && make lint && make test
   ```

4. **Review Release Notes**: Ensure docs are clear

5. **Announce Releases**: Update documentation, announce in discussions

6. **Monitor Issues**: Address problems quickly

## Support and Resources

- **Terraform Registry**: https://registry.terraform.io
- **Provider Publishing Guide**: https://registry.terraform.io/docs/providers/publish-provider
- **Terraform Forum**: https://discuss.hashicorp.com/c/terraform-core/
- **GitHub Issues**: https://github.com/ramorous/terraform-provider-dockhand/issues
- **Release Guide**: See [RELEASING.md](./RELEASING.md)
- **Development Guide**: See [DEVELOPMENT.md](./DEVELOPMENT.md)

## Quick Reference Commands

```bash
# Local development
make build          # Build binary
make test          # Run tests
make install       # Install locally
make fmt           # Format code
make lint          # Lint code
make doc           # Generate docs

# Releasing
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# Check status
git tag -l             # List all tags
git log --oneline      # View commit history
```

## Next Steps

1. Follow the step-by-step process above
2. Monitor GitHub Actions for build completion
3. Verify release artifacts
4. Publish to Terraform Registry
5. Announce to users
6. Continue maintaining and releasing updates
