[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](https://github.com/tterb/atomic-design-ui/blob/master/LICENSEs)
# Shopware CLI

A cli which contains handy helpful commands for daily Shopware tasks

## Features

- Manage your Shopware account extensions in the CLI
- Build and validate Shopware extensions
## All commands

```
shopware-cli account login
shopware-cli account logout
shopware-cli account company list
shopware-cli account company use [companyId]
shopware-cli account producer
shopware-cli account producer info
shopware-cli account producer extension create [name] [generation]
shopware-cli account producer extension delete [name]
shopware-cli account producer extension list
shopware-cli extension validate [folder or zip path]
shopware-cli-cli extension prepare [folder]
shopware-cli extension zip [folder]
```
## Installation

There are many options to install shopware-cli. The binary file itself can be found in the latest GitHub release. 
The releases contain also packages for Debian, Red Hat and Alpine.

For Homebrew use `brew install FriendsOfShopware/tap/shopware-cli`.

## Develop Locally

Clone the project

```bash
  git clone https://github.com/FriendsOfShopware/shopware-cli.git
```

Go to the project directory

```bash
  cd shopware-cli
```

Run the cli

```bash
    go run . account login
```
## Contributing

Contributions are always welcome!
