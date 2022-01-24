---
title: Extension Store Quick Start
weight: 10
---

In this Guide we will learn how Shopware CLI can make the extension deployment easier with the Shopware Store. 
First you need to [install](../install) the CLI.

## Login into your Shopware Account

First we need to login into our Shopware account. For this execute the command `shopware-cli account login` and login with your credentials. This credentials will be saved locally. With `shopware-cli account logout` can you logout again. For the CI you can set enviroment variables `SHOPWARE_CLI_ACCOUNT_EMAIL` and `SHOPWARE_CLI_ACCOUNT_PASSWORD` skip the login step.


## Optional: Change the active Shopware Account company

Your Shopware account can be in multiple companies. Use the command `shopware-cli account company list` to show all companys you have access.
With `shopware-cli account company use <id>` can you switch the current company access.

## Optional: Create the extension in the Account

To upload the zip later in the Store, you need to create the extension in the Account. If you haven't done this you can do this with

```
shopware-cli extension create <Name> platform
```

possible values are: `classic` (Shopware 5) `platform` (Shopware 6 Plugin system) `themes` (Shopware 6 App containg theme) `apps` (Shopware 6 App)

## Getting the Store Information into the Git repository

To edit the Store page locally, we need first to generate the local files based on the current store page. 
For this we can use the command `shopware-cli account producer extension info pull <extension-folder>`.
This command creates a `.shopware-extension.yml` config in the root folder with the current store page. The schema of the file can be found [here](../shopware-extensions.yml-schema/). Editors supporting SchemaStore, should have autocomplete out of the box like Jetbrains products, VSCode.

## Uploading local Store Information to the Store

After making changes on the `.shopware-extension.yml`, you have to run the `shopware-cli account producer extension info push <extension-folder>` to apply the changes on the store page. 

## Uploading the extension zip to the Store

To upload the extension zip, you have to use the command `shopware-cli account producer extension upload <extension-zip>`.
If the version already exists, it updates the zip. The Shopware version compatibility list is built from the `shopware/core` requirement in the composer.json. Changelogs are built from the `CHANGELOG_de-DE.md` and `CHANGELOG_en-GB.md` file. 

Here is an example content of an changelog file

```
# 0.1.0

* First release in Store
```

The changelog has to be written in both languages. 