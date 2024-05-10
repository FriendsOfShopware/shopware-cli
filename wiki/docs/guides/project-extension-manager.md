---
title: Project Extension Manager
weight: 13
---

shopware-cli has an extension manager to install and manage extensions in your Shopware project through the Shopware API. Kinda like the Extension Manager in the Shopware 6 Administration Panel, but for the CLI.

!!! note
    This functionality was designed for Shopware SaaS and should not be used for self-hosted installations. [The recommandation is to use the Deployment Helper and install all plugins via Composer.](https://developer.shopware.com/docs/guides/hosting/installation-updates/deployments/deployment-helper.html)

To use the extension manager, you need a `.shopware-project.yml`, this can be created with the command `shopware-cli project config init`.


## Commands

### List all extensions

```bash
shopware-cli project extension list
```

### Install an extension

```bash
shopware-cli project extension install <extension-name>
```

### Uninstall an extension

```bash
shopware-cli project extension uninstall <extension-name>
```

### Update an extension

```bash
shopware-cli project extension update <extension-name>
```

### Outdated extensions

Shows all extensions that have an update available.

```bash
shopware-cli project extension outdated
```

### Upload extension

Uploads an extension to the Shopware instance.

```bash
shopware-cli project extension upload <path-to-extension-zip>
```

### Delete extension

Deletes an extension from the Shopware instance.

```bash
shopware-cli project extension delete <extension-name>
```
