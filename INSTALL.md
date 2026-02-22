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
