# GoPack

Simple CLI to search and install packages with `gop get <package-name> <...more>` or interactively by just running `gop`.

# Installation

```
1. Clone repository
git clone https://github.com/juancwu/gopack

2. CD into repository and build
cd gopack && make build

3. Move binary to your desired location
mv ./build/gop <your-path>
```

# Usage

There is a `help` command (`gop --help`), but here's a comprehensive guide to all commands.

## Get Command

The `get` command allows you to search and install Go packages.

### Single Package

Following the naming scheme for Go packages, GoPack assumes the prefix `github.com/`.

You can directly use the package name or any query string you would use in [https://pkg.go.dev](https://pkg.go.dev), such as `something/package` and GoPack will search the Go directory
and select the first match of the results and install it.

Example: `gop get package` will install `github.com/something/package`.

You can make GoPack show the search results and select manually by passing the option `-select` or `-s`.

### Multiple Packages

GoPack accepts multiple query strings separated by spaces. The default behaviour is the same as [Single Package](#single-package) and you can select the matches for each query
by using the option `-select` or `-s`.

Example: `gop get package something/else` will install `github.com/something/package` and `github.com/something/else`.

## List Command

The `list` command displays all installed packages with their installation paths and versions.

Usage: `gop list`

## Run Command

The `run` command allows you to execute scripts defined in your `gopack.json` configuration file.

Usage examples:
- `gop run build` - Runs the "build" script defined in your configuration
- `gop run --list` or `gop run -l` - Lists all available scripts
- `gop run --init` or `gop run -i` - Initializes a new gopack.json configuration file
- `gop run --config custom.json` or `gop run -c custom.json` - Specifies a custom configuration file path

You can also run scripts directly from the root command:
- `gop build` - Equivalent to `gop run build`

## Update Command

The `update` command updates GoPack to the latest version from GitHub.

Usage: `gop update`

Features:
- Checks for the latest release on GitHub
- Downloads appropriate binary for your OS/architecture
- Replaces the current executable with the new version

## Version Command

The `version` command displays the current version of GoPack.

Usage: `gop version`

## Configuration

GoPack uses a `gopack.json` file to define scripts that can be run with the `run` command.

Example `gopack.json`:
```json
{
  "scripts": {
    "build": "go build -o gop",
    "test": "go test ./...",
    "lint": "go vet ./...",
    "format": "go fmt ./..."
  }
}
```

## Removing Packages

Just `go mod tidy`.

## Download All Dependencies

Just `go mod download`
