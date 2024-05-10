---
title: Project Config Synchronization
weight: 12
---

shopware-cli can synchronize the project configurations between different environments. This is useful for example to keep the configuration in the development and production environment in sync.

Following things are possible to synchronize:

- Theme Configuration
- System Configuration (including extension configuration)
- Mail Templates
- Entity

## Setup

To synchronize the project, you need to create a `.project.yml` file in the root of your project. This file contains the configuration for the synchronization.

You can use the command `shopware-cli project config init` to create a new `shopware-project.yml` file. Make sure that you configure the API access too as this is required for the synchronization.

## Initial pulling

To pull the configuration from the Shopware instance, you can use the command `shopware-cli project config pull`. This command pulls the configuration from the Shopware instance and stores it in the local `shopware-project.yml` file.

## Pushing the configuration

After you made the changes in the local `shopware-project.yml` file, you can push the changes to the Shopware instance with the command `shopware-cli project config push`.

This shows the difference between your local and the remote configuration and asks you if you want to push the changes.

## Entity synchronization

With Entity synchronization, you can synchronize any kind of entity using directly the Shopware API.

```yaml
sync:
  entity:
      - entity: tax
        payload:
          name: 'Tax'
          taxRate: 19
```

This example synchronizes a new tax entity with the name `Tax` and the tax rate `19`.

The further synchronizations will create the same entity again, you may want to fixed the entity ID to avoid duplicates.

You can also add an existence check, so it will be only created if an entity has been found:

```yaml
sync:
  entity:
    - entity: tax
      # build a criteria to check that the entity already exists. when exists this will be skipped
      exists:
        - type: equals
          field: name
          value: 'Tax'
      # actual api payload to create something
      payload:
        name: 'Tax'
      taxRate: 19
```
