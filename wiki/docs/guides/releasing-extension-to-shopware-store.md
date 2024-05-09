---
title: Releasing a Extension to Shopware Store
weight: 10
---


In this Guide we will learn how Shopware CLI can make the extension deployment easier with the Shopware Store. 
First you need to [install](../install.md) the CLI.

## Login into your Shopware Account

First we need to login into our Shopware account. For this execute the command `shopware-cli account login` and login with your credentials. This credentials will be saved locally. With `shopware-cli account logout` can you logout again. For the CI you can set enviroment variables `SHOPWARE_CLI_ACCOUNT_EMAIL` and `SHOPWARE_CLI_ACCOUNT_PASSWORD` skip the login step.

## Optional: Change the active Shopware Account company

Your Shopware account can be in multiple companies. Use the command `shopware-cli account company list` to show all companys you have access.
With `shopware-cli account company use <id>` can you switch the current company access.

## Optional: Create the extension in the Account

To upload the zip later in the Store, you need to create the extension in the Shopware Store. If you haven't done this you can do this with

```
shopware-cli account producer extension create <Name> platform
```

possible values for the last parameter are: 
* `classic` (Shopware 5)
* `platform` (Shopware 6 Plugin system)
* `themes` (Shopware 6 App containg theme)
* `apps` (Shopware 6 App)

## Getting the Store Information into the Git repository

To edit the Store page locally, we need first to generate the local files based on the current store page. 
For this we can use the command `shopware-cli account producer extension info pull <extension-folder>`.
This command creates a `.shopware-extension.yml` config in the extension root folder with the current store page. The schema of the file can be found [here](../shopware-extension-yml-schema.md). Editors supporting SchemaStore, should have autocomplete out of the box like Jetbrains products, VSCode.

## Uploading local Store Information to the Store

After making changes on the `.shopware-extension.yml`, you have to run the `shopware-cli account producer extension info push <extension-folder>` to apply the changes on the store page. 

## Creating the zip for the Store

To create the zip, we can use `shopware-cli extension zip <path> --release`, this command creates a zip file from the latest Git tagged version. You can specify a specific Git commit / tag with `--git-commit 1.0.0`, or disable this behavior with `--disable-git` to take the directory as it is.

## Validating the zip

To save time, you can validate the zip with `shopware-cli extension validate <zip-path>` before uploading it to the Store. This command checks the most things of the extension store upload process.

## Uploading the extension zip to the Store

To upload the extension zip, you have to use the command `shopware-cli account producer extension upload <zip-path>`.
If the version already exists, it updates the zip. The Shopware version compatibility list is built from the `shopware/core` requirement in the composer.json. Changelogs are built from the `CHANGELOG_de-DE.md` and `CHANGELOG_en-GB.md` file or only `CHANGELOG.md` then for both languages: 

Here is an example content of an changelog file

```
# 0.1.0

* First release in Store
```

The changelog has to be written in both languages. 


## Generating the changelog

It's also possible to let shopware-cli generate the changelog for you. For this you need to create a `.shopware-extension.yml` with following content:

```yaml
changelog:
    enabled: true
```

After that, when you create a new zip using `shopware-cli extension zip <path> --release`, the changelog will be generated based on the git tags. It's important to set the flag `--release`, otherwise the changelog will not be generated.

### Configuration of the changelog generation

The changelog generation is very flexible, here is an full example of a Shopware plugin:

```yaml
changelog:
  enabled: true
  # only the commits matching to this regex will be used
  pattern: '^NEXT-\d+'
  # variables allows to extract metadata out of the commit message
  variables:
    ticket: '^(NEXT-\d+)\s'
  # go template for the changelog, it loops over all commits
  template: |
    {{range .Commits}}- [{{ .Message }}](https://issues.shopware.com/issues/{{ .Variables.ticket }})
    {{end}}
```

This example checks that all commits in the changelog needs to start with `NEXT-` in the beginning. The `variables` section allows to extract metadata out of the commit message. The `template` is a go template which loops over all commits and generates the changelog.
With the combination of `pattern`, `variables` and `template` we link the commit message to the Shopware ticket system.
