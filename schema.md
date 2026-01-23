# crx-registry YAML Schema

## Package Schema (pkgs/{name}.yaml)

```yaml
# Required fields
name: string          # Package name (must match filename, kebab-case)
id: string            # Chrome Web Store Extension ID (32 characters)
display_name: string  # Display name

# Optional fields
description: string   # Description
homepage: string      # Official website URL
repository: string    # Source code repository URL
tags: []string        # Tags for search
```

### Example

```yaml
name: ublock-origin
id: cjpalhdlnbpafiamejdnhcphjbkeiagm
display_name: uBlock Origin
description: An efficient blocker for Chromium and Firefox
homepage: https://ublockorigin.com/
repository: https://github.com/gorhill/uBlock
tags:
  - ad-blocker
  - privacy
  - security
```

## Registry Index (registry.yaml)

```yaml
version: 1
packages:
  - ublock-origin
  - vimium
  - react-devtools
  # ... list of package names
```

## User Config (~/.config/crx/config.yaml)

```yaml
# Registry settings
registries:
  - name: standard
    type: github
    repo: user/crx-registry
    ref: main  # optional, default: main

# Extensions to install
extensions:
  - ublock-origin
  - vimium

# Settings
settings:
  # Policy JSON output path
  policy_path: /Library/Google/Chrome/policies/managed

  # Install mode
  # force_install: Force install (cannot be removed)
  # normal_install: Auto install (can be disabled)
  mode: force_install
```

## Validation Rules

1. `name` must match the filename
2. `id` must be 32 alphanumeric characters
3. `tags` must be lowercase kebab-case
4. `homepage` and `repository` must be valid URLs

## Extension ID Format

Chrome Extension ID consists of 32 lowercase letters (a-p).

Example: `cjpalhdlnbpafiamejdnhcphjbkeiagm`

Can be obtained from Chrome Web Store URL:
```
https://chromewebstore.google.com/detail/ublock-origin/cjpalhdlnbpafiamejdnhcphjbkeiagm
                                                       ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
```
