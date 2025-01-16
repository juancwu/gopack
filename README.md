# GoPack

Simple CLI to search and install packages with `gpk get <package-name> <...more>` or interactively by just running `gpk`.

# Installation

```
1. Clone repository
git clone https://github.com/juancwu/gopack

2. CD into repository and build
cd gopack && make build

3. Move binary to your desired location
mv ./build/gpk <your-path>
```

# Usage

There is a `help` command, but there is also this.

## Single Package

Following the naming scheme for Go packages, GoPack assumes the prefix `github.com/`.

You can directly use the package name or any query string you would use in [https://pkg.go.dev](https://pkg.go.dev), such as `something/package` and GoPack will search the Go directory
and select the first match of the results and install it.

Example: `gpk get package` will install `github.com/something/package`.

You can make GoPack to show the search results and select manually by passing the option `-select`.

## Multiple Packages

GoPack accepts multiple query strings separated by spaces. The default behaviour is the same as [Single Package](#single-package) and you can select the matches for each query
by using the option `-select`.

Example: `gpk get package something/else` will install `github.com/something/package` and `github.com/something/else`.

## Removing Packages

Just `go mod tidy`.
