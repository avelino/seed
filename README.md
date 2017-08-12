# seed

Library packaging and distribution for Golang

**status: WIP**

## Install

	curl https://raw.githubusercontent.com/goseed/install/master/initial.sh | sh


## Update

	seed get -u seed


## API

| command | alias | parameters | description |
|---|---|---|---|
| search | s | - | Find remote Seed to an Index Server |
| register | r | -f Seed.toml | The distutils command register is used to submit your distribution’s meta-data to an Seed Index Server |
| push | p | -f Seed.toml | The distutils command upload pushes the distribution files to Seed Index Server |
| get | g | -u / -f Seed.toml / -to [`gopath`, `vendor`] | Fetch from and integrate with remote repository to **GOPATH** or **vendor** (if exist folder vendor this path) |
| install | i | -u / -f Seed.toml (requires file) | Installs all packages from the toml file |
| list | l | -f Seed.toml | Shows your locally installed to **GOPATH** or **vendor** (if exist folder vendor this path) |


## Config

File: **~/.seedrc**

```
[seed]
verbose = false
sources = [
	"https://packages.goseed.io/",
]

[source]
[[packages.goseed.io]]
token = "my key"
# OR
# username = ""
# password = ""


```


## Package

### Seed.toml

```
[package]
name = "seed"
version = "0.1"
authors = ["Frist Last Name <mail@goseed.io>"]
description = "Package Manager"
homepage = "https://goseed.io/"
documentation = "https://docs.goseed.io/"
repository = "https://github.com/goseed/seed"
readme = "README.md"
keywords = ["package", "manager"]
categories = ["command-line-utilities", "network-programming"]
license = "MIT"
exclude = [
    "assets/*",
    "packageX/**/*.go",
	"vendor/*",
]
include = [
    "**/*.go",
    "Seed.toml",
]

[dependencies]
package_name = "0.1"
package_name2 = { git = "https://github.com/username/repo", branch = "master" }
```
