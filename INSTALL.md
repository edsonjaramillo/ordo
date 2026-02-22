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
```

Global package management:

```bash
ordo global install pnpm typescript
ordo global uninstall pnpm typescript
```

Global store lookup notes:

- `pnpm` global modules are resolved from `pnpm root -g` first, then from `PNPM_HOME/global/<layout>/node_modules` and default pnpm global locations.

Shell completion support:

- `ordo install <TAB>` suggests package names (registry-aware with local fallback).
- `ordo install --workspace <TAB>` suggests discovered workspace keys.
- `ordo global install <TAB>` suggests package managers first, then package names.
- `ordo global uninstall <TAB>` suggests package managers first, then installed global package names for the selected manager.
