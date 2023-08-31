# foto

_Yet another another publishing tool for minimalist photographers._

Demo site: [https://foto.lhzhang.com](https://foto.lhzhang.com)

<img width="550" alt="" src="https://user-images.githubusercontent.com/480052/181399181-25d3657d-7fff-4ad0-97a0-392a8effb18e.png">

## Features

- `Simple` One binary, three commands. No database.
- `Customizable` Highly customizable by configuration and template files.

## Install

### macOS

```bash
$ brew tap waynezhang/tap
$ brew install foto
```

Or download the binary from [here](https://github.com/waynezhang/foto/releases)

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

## Customize

### Basic customize `foto.toml`

<details>
  <summary>Click to expand</summary>
  
  ```toml
  [site]
  # The title of the site
  title = "A new site"
  # The name of the author
  author = "Author Here"
  
  # Site navigation links
  # You can remove any navigation links or add more link by adding following lines
  #     [[site.nav]]
  #     icon = ""
  #     link = ""
  # Navigation links are added in the order encountered.
  
  [[site.nav]]
  icon = "assets/icons/home.svg"
  link = "https://"
  
  [[site.nav]]
  icon = "assets/icons/instagram.svg"
  link = "https://instagram.com/xxx"
  
  [[site.nav]]
  icon = "assets/icons/twitter.svg"
  link = "https://twitter.com/xxx"
  
  # Setttings for photo size
  [image]
  # Width for thumbnail images
  thumbnailWidth = 640
  # Width for enlarged images
  originalWidth = 2048

  # Layout for grids
  [layout]
  minColumn = 1
  maxColumn = 4
  minWidth = 200
  
  # Photo sections
  # You can remove or add more sections by adding following lines
  #     [[section]]
  #     title = "section title"
  #     text = "section description (HTML supported)"
  #     slug = "section-slug"
  #     folder = "folder of photos"
  #     ascending = false
  # Photo sections are added in the order encountered.
  [[section]]
  title = "Section 1"
  text = ""
  slug = "section-1"
  folder = "~/photos/section-1"
  ascending = false
  
  [[section]]
  title = "Section 2"
  text = ""
  slug = "section-2"
  folder = "~/photos/section-2"
  ascending = false
  
  # Other setings
  [others]
  # Folders that should be copied together when exporting sites
  folders = [ "assets", "media" ]
  # Show `Generated by foto` footer or not
  show_foto_footer = true
  ```
</details>

### Style customize

Template and CSS styles can be modified without chagning `foto` binary.

The template file is placed in `templates/template.html`.

It's also possible to add additional settings in `foto.toml` ([ref](https://toml.io/en)) and refer it in template file.
`foto` uses `html/template` package in Golang. Please refer to [this link](https://pkg.go.dev/html/template) for more information.

# LICENSE
 
See [LICENSE](./LICENSE)

# Credit

`foto` is highly inspried by [moul](https://moul.app).
