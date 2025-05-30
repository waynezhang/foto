![foto-cover](https://github.com/waynezhang/foto/assets/480052/13e77201-c680-49f5-8ce7-9ba0c73e6ddc)

# foto

![](https://github.com/waynezhang/foto/actions/workflows/release.yml/badge.svg) ![](https://github.com/waynezhang/foto/actions/workflows/test.yml/badge.svg)

_Yet another another publishing tool for minimalist photographers._

Demo site: [https://foto.lhzhang.com](https://foto.lhzhang.com)

## Features

- `Simple` One binary, three commands. No database required.
- `Customizable` Highly customizable through configuration and template files.
- `Fast` Files are processed concurrently for optimal performance.

## Installation

### macOS

```bash
$ brew tap waynezhang/tap
$ brew install foto
```

Or download the binary from [here](https://github.com/waynezhang/foto/releases)

### Nix/NixOS

For Nix users, a Flake is provided. It can be used to run the application
directly or add the package to your configuration as flake input.

It also allows you to try out foto without permanent installation.

```sh
nix run github:waynezhang/foto
```

Consult the [Nix
manual](https://nix.dev/manual/nix/2.25/command-ref/new-cli/nix3-flake.html) for
details.

### Other platforms

Download the binary from [here](https://github.com/waynezhang/foto/releases)

## Usage

### Create a new site

```bash
~ $ foto create my_site
~ $ tree my_site
my_site
├── assets
│   ├── icons
│   │   ├── home.svg
│   │   ├── instagram.svg
│   │   └── twitter.svg
│   └── style.css
├── foto.toml # Configuration file, see below for more details.
├── media
│   └── avatar.jpg # Placeholder image for avatar.
└── templates
    └── template.html # Template file
```

### Preview

```bash
~/my_site $ foto preview
Creating Preview...
Listening on 5000...
```

The default port number is `5000`. It can be changed by `-p` flag.

### Export

```bash
~/my_site $ foto export -o ~/site_docs
Exprorting sites to /Users/xxx/site_docs...
```

### Clear cache

```bash
foto clear-cache
```

## Customization

### Basic configuration with `foto.toml`

See [foto.toml](./fs/static/foto.toml)

### Style customization

Template and CSS styles can be modified without changing the `foto` binary.

The template file is located at `templates/template.html`.

You can also add additional settings in `foto.toml` ([ref](https://toml.io/en)) and reference them in the template file.
`foto` uses the `html/template` package from Go. Please refer to [this link](https://pkg.go.dev/html/template) for more information. Besides, EXIF information is supported. Refer to [EXIF](https://exiftool.org/TagNames/EXIF.html) for all EXIF tags.

## Changelogs

See [CHANGELOG](./CHANGELOG.md)

## LICENSE
 
See [LICENSE](./LICENSE)

## Credit

`foto` is highly inspired by [moul](https://moul.app).
