---
 title: 'Schema of .shopware-project.yml' 
---

Any configuration field is optional. When you create a `.shopware-project.yml`, you get also IDE autocompletion for all fields.


```yaml
# .shopware-project.yml

# URL to Shopware instance, required for admin api calls (clear cache, sync stuff)
url: 'http://localhost'
admin_api:
    # For integration use this both fields
    client_id:
    client_secret:
    # For normal user use this both fields
    username:
    password:
    # When your server don't have a valid SSL certificate, you can disable the SSL check
    disable_ssl_check: false

# used only for project ci command
build:
  # deletes all public source folders of all extensions, can be only used when /bundles is served from local and not external CDN
  remove_extension_assets: false
  # skips the bin/console asset:install part
  disable_asset_copy: false
  # when enabled src/Resources/app/{storefront/administration} folder will be preserved and not deleted.
  # If your plugin requires, you should move the files out of src/Resources which needs to be accessed by php and js
  keep_extension_source: false
  # delete additional paths after build
  cleanup_paths:
    - path
  # change the browserslist of the storefront build, see https://browsersl.ist for the syntax as string (example: defaults, not dead)
  browserslist: ''

# used for mysql dump creation
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
    # Sync system config to your remote shop using admin-api
    config:
        # can be also null for default value
        - sales_channel: yourSalesChannelid
          settings:
            my_config: myValue
    # Sync theme config to your remote shop using admin-api
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