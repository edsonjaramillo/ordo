# Install `ordo`

`ordo` is currently distributed through GitHub Releases.

## Quick install (latest release)

```bash
curl -fsSL https://github.com/edsonjaramillo/ordo/releases/latest/download/install.sh | sh
```

By default this installs to `~/.local/bin`.

## Install a pinned version

```bash
curl -fsSL https://github.com/edsonjaramillo/ordo/releases/latest/download/install.sh | sh -s -- --version v0.2.0
```

## Install to a custom directory

```bash
curl -fsSL https://github.com/edsonjaramillo/ordo/releases/latest/download/install.sh | sh -s -- --install-dir /usr/local/bin
```

If `/usr/local/bin` needs elevated permissions:

```bash
curl -fsSL https://github.com/edsonjaramillo/ordo/releases/latest/download/install.sh -o /tmp/ordo-install.sh
sudo ORDO_INSTALL_DIR=/usr/local/bin sh /tmp/ordo-install.sh
```

## Supported platforms

- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`

## Troubleshooting

- `ordo: command not found`
  - Ensure your install directory is on `PATH`:
    - `export PATH="$HOME/.local/bin:$PATH"`
- Unsupported OS/architecture
  - Only Linux/macOS on `amd64`/`arm64` are supported today.

## Smart dependency install in monorepos

Once `ordo` is installed, you can install dependencies in root or workspace:

```bash
ordo install zod
ordo install eslint typescript --dev --exact
ordo install react --workspace ui
ordo update react
ordo update ui/clsx
ordo preset prettier devDependencies
ordo preset prettier devDependencies prettier-plugin-tailwindcss --workspace ui
ordo catalog presets prettier devDependencies
ordo catalog presets prettier devDependencies prettier-plugin-tailwindcss --workspace ui --force
```

Global package management:

```bash
ordo global install pnpm typescript
ordo global uninstall pnpm typescript
ordo global update pnpm typescript
```

Global store lookup notes:

- `pnpm` global modules are resolved from `pnpm root -g` first, then from `PNPM_HOME/global/<layout>/node_modules` and default pnpm global locations.

Shell completion support:

- `ordo install <TAB>` suggests package names (registry-aware with local fallback).
- `ordo install --workspace <TAB>` suggests discovered workspace keys.
- `ordo global install <TAB>` suggests package managers found on `PATH` first (falls back to all supported managers if none are detected), then package names.
- `ordo global uninstall <TAB>` suggests package managers found on `PATH` first (falls back to all supported managers if none are detected), then installed global package names for the selected manager.
- `ordo update <TAB>` suggests installed dependency targets.
- `ordo global update <TAB>` suggests package managers found on `PATH` first (falls back to all supported managers if none are detected), then installed global package names for the selected manager.
- `ordo init --defaultPackageManager <TAB>` suggests package managers found on `PATH` (falls back to all supported managers if none are detected).
- `ordo preset <TAB>` suggests preset names from `ordo.json`.
- `ordo preset <preset> <TAB>` suggests preset buckets that have packages.
- `ordo preset <preset> <bucket> <TAB>` suggests package names for that preset bucket.
- `ordo preset --workspace <TAB>` suggests discovered workspace keys.
- `ordo catalog presets <TAB>` suggests preset names from `ordo.json`.
- `ordo catalog presets <preset> <TAB>` suggests preset buckets that have packages.
- `ordo catalog presets <preset> <bucket> <TAB>` suggests package names for that preset bucket.
- `ordo catalog presets --workspace <TAB>` suggests discovered workspace keys.
