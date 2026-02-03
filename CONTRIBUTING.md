# Contributing to crx-registry

Thank you for your interest in contributing to crx-registry!

## Adding a New Extension

1. **Get the Extension ID from Chrome Web Store**
   ```
   https://chromewebstore.google.com/detail/extension-name/{ID}
                                                           ↑ 32 characters
   ```

2. **Create `pkgs/{name}.yaml`**
   ```yaml
   name: extension-name        # Must match filename (kebab-case)
   id: xxxxxxxx...             # 32-character Extension ID
   display_name: Extension Name
   description: A brief description of the extension
   homepage: https://...
   repository: https://github.com/...  # optional
   tags:
     - developer-tools
   ```

3. **Validate**
   ```bash
   make validate
   ```

4. **Submit a Pull Request**

## Schema Reference

### Required Fields

| Field | Description |
|-------|-------------|
| `name` | Package name (must match filename, kebab-case) |
| `id` | Chrome Extension ID (32 lowercase letters a-p) |
| `display_name` | Human-readable display name |

### Optional Fields

| Field | Description |
|-------|-------------|
| `description` | Brief description of the extension |
| `homepage` | Official website URL |
| `repository` | Source code repository URL |
| `tags` | List of tags (kebab-case) |

## Extension ID Format

Chrome Extension IDs consist of 32 lowercase letters (a-p).

Example: `dbepggeogbaibhgnhhndojpepiihcmeb`
