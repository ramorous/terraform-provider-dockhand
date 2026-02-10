# Releasing the Terraform Provider

This document describes how to release a new version of the Terraform Provider for Dockhand.

## Prerequisites

Before releasing, ensure you have:

1. **GPG Key Setup**
   - A GPG key pair for signing releases
   - The private key exported in ASCII format
   - The key's fingerprint noted

2. **GitHub Secrets Configuration**
   - Go to https://github.com/ramorous/terraform-provider-dockhand/settings/secrets/actions
   - Add the following secrets:
     - `GPG_PRIVATE_KEY`: Your GPG private key (ASCII armored format)
     - `GPG_PASSPHRASE`: Your GPG key passphrase

3. **GoReleaser Installation**
   ```bash
   brew install goreleaser
   # or
   curl -sL https://git.io/goreleaser | bash
   ```

## Version Numbering

This project uses semantic versioning:
- `MAJOR.MINOR.PATCH` (e.g., `1.2.3`)
- Pre-releases use `-alpha`, `-beta`, `-rc` suffixes (e.g., `1.0.0-beta.1`)

## Release Process

### 1. Prepare the Release

Update version information:

```bash
# Update version in main.go if needed
# Update CHANGELOG.md with changes for this release
# Update README.md if needed

git add CHANGELOG.md README.md
git commit -m "docs: prepare release v1.0.0"
```

### 2. Create a Release Tag

```bash
# Create an annotated tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push the tag to trigger the release workflow
git push origin v1.0.0
```

### 3. Automated Release

The GitHub Actions workflow will automatically:
- Build binaries for all platforms (Linux, macOS, Windows)
- Generate checksums
- Sign checksums with GPG
- Create a GitHub release with artifacts
- Generate documentation

### 4. Manual Release (if needed)

If you prefer to create releases locally:

```bash
# Build snapshot
make release-snapshot

# Release
make release
```

After release, verify in GitHub:
- Check the release page: https://github.com/ramorous/terraform-provider-dockhand/releases
- Verify binaries are present and correctly named
- Verify checksums file is signed

## Publishing to Terraform Registry

After a release is created on GitHub:

1. **Verify Release Assets**
   - Checksums file: `terraform-provider-dockhand_1.0.0_SHA256SUMS`
   - Signed checksums: `terraform-provider-dockhand_1.0.0_SHA256SUMS.sig`
   - Platform binaries: `terraform-provider-dockhand_1.0.0_linux_amd64.zip`, etc.

2. **Publish to Registry**
   - Visit https://registry.terraform.io
   - Sign in or create an account
   - Navigate to "Publish" → "Provider"
   - Select "GitHub" as source
   - Authorize and follow prompts
   - Select the `ramorous/terraform-provider-dockhand` repository

3. **Verify Publication**
   - Provider appears at: https://registry.terraform.io/providers/ramorous/dockhand
   - Documentation is auto-generated from provider code
   - Versions are listed in the registry

## Troubleshooting

### GPG Signing Issues

If GoReleaser fails to sign:

```bash
# Verify GPG key is available
gpg --list-secret-keys

# Export private key in ASCII format
gpg --armor --export-secret-key <KEY_ID> > private-key.asc

# Test signing with the key
gpg --sign --detach-sign --armor private-key.asc
```

### Binary Naming

Ensure binaries follow the pattern:
```
terraform-provider-dockhand_v1.0.0_OS_ARCH.zip
```

Supported OS values: `darwin`, `linux`, `windows`
Supported ARCH values: `amd64`, `arm64`

### Registry Publication Fails

Common issues:
- Release not properly signed (verify checksums file has `.sig` companion)
- Repository not public
- Binary naming doesn't match expected format
- Documentation generation failed (ensure `go:generate` tags are present)

## Documentation

Documentation is automatically generated from:
- Provider schema comments
- Resource/DataSource schema comments
- Examples in `examples/` directory

After release, verify docs appear correctly at:
https://registry.terraform.io/providers/ramorous/dockhand/latest/docs

## Support

For issues with the release process:
- Check GitHub Actions logs
- Review GoReleaser documentation: https://goreleaser.com
- Check Terraform Registry documentation: https://registry.terraform.io/docs/providers
