---
 title: 'Schema of .shopware-project.yml' 
---

Any configuration field is optional. When you create a `.shopware-project.yml`, you get also IDE autocompletion for all fields.


```yaml
# .shopware-project.yml

# URL to Shopware instance, required for admin API calls (clear cache, sync stuff)
url: 'http://localhost'
admin_api:
    # For integration use these two fields
    client_id:
    client_secret:
    # For normal user use these two fields
    username:
    password:
    # When your server doesn't have a valid SSL certificate, you can disable the SSL check
    disable_ssl_check: false

# used only for project CI command
build:
  # deletes all public source folders of all extensions, can only be used when /bundles is served from local and not external CDN
  remove_extension_assets: false
  # skips the bin/console asset:install part
  disable_asset_copy: false
  # when enabled src/Resources/app/{storefront/administration} folder will be preserved and not deleted.
  # If your plugin requires, you should move the files out of src/Resources which need to be accessed by PHP and JS
  keep_extension_source: false
  # delete additional paths after build
  cleanup_paths:
    - path
  # change the browserslist of the storefront build, see https://browsersl.ist for the syntax as string (example: defaults, not dead)
  browserslist: ''
  # exclude extensions to be built by shopware-cli, only their PHP code will be shipped without any CSS/JS
  exclude_extensions:
    - name

# used for MySQL dump creation
dump:
    # rewrite columns
    rewrite:
        table:
            column: "'new-value'"
            column2: "faker.Internet().Email()" # Uses faker data. See https://github.com/jaswdr/faker
    # ignore table content
    nodata:
        - this-table-is-dumped-without-rows
    # ignore complete tables
    ignore:
        - this-table-is-not-dumped
    # add where conditions to tables
    where:
        my_table: "id > 10"

# you can use shopware-cli project config pull, to get your current shop state
sync:
    # Sync system config to your remote shop using admin API
    config:
        # can be also null for default value
        - sales_channel: yourSalesChannelId
          settings:
            my_config: myValue
    # Sync theme config to your remote shop using admin API
    theme:
        - name: ThemeName
          settings:
            my_config: myValue

    mail_template:
        - id: mailTemplateId
          translations:
            - language: de-DE
              sender_name: 'Sender Name'
              subject: 'Subject'
              html: relativeFilePath
              plain: relativeFilePath
              custom_fields: null
    entity:
        - entity: tax
          # optional: build a criteria to check that the entity already exists. when exists this will be skipped
          exists:
            - type: equals
              field: name
              value: 'Tax'
          # actual api payload to create something
          payload:
            name: 'Tax'
            taxRate: 19
```

## Advanced usage

### Configuration includes

You can include one or more `.shopware-project.yml` files to reuse or override configurations. This is useful when you have multiple projects with the same configuration. You can also use this to create a base configuration for your stages or teams and extend it for your own needs.
This also can be used to toggle specific plugin configurations for different stages e.g. enabling/disabling the Paypal sandbox mode depending on the environment.

Parent `.shopware-project.base.yml`:

```yaml
url: 'http://localhost'
admin_api:
  # For integration use this both fields
  client_id: 'client id'
  client_secret: 'client secret'
```

Child `.shopware-project.dev.yml`:

```yaml
include:
  - '.shopware-project.base.yml'
url: 'http://dev.localhost.test'
sync:
  config:
    - settings:
        SwagPayPal.settings.sandbox: true
```

Child `.shopware-project.prod.yml`:

```yaml
include:
  - '.shopware-project.base.yml'
url: 'http://prod.localhost.test'
sync:
  config:
    - settings:
        SwagPayPal.settings.sandbox: false
```

You would apply them using the `--project-config` option:

```bash
# for development
shopware-cli project --project-config='.shopware-project.dev.yml' config push

# for production
shopware-cli project --project-config='.shopware-project.prod.yml' config push
```

### Environment Variables

You can use environment variables in your `.shopware-project.yml` file. This is useful for example when you want to use the same configuration for multiple projects or environments. This can also be useful to apply secrets without committing them in the config directly.

```yaml
# .shopware-project.yml

url: 'http://localhost'
admin_api:
    # there are two valid environment variable syntax
    client_id: ${SHOPWARE_CLI_CLIENT_ID}
    client_secret: $SHOPWARE_CLI_CLIENT_SECRET
```