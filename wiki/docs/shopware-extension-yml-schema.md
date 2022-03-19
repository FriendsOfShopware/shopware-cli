---
title: 'Schema of .shopware-extension.yml'
---

**Title:** shopware-cli

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/Config                                                      |
|                           |                                                                           |

**Description:** shopware cli project configuration definition file

| Property                   | Pattern | Type   | Deprecated | Definition                | Title/Description        |
| -------------------------- | ------- | ------ | ---------- | ------------------------- | ------------------------ |
| - [url](#url )             | No      | string | No         | -                         | URL to Shopware instance |
| - [admin_api](#admin_api ) | No      | object | No         | In #/definitions/AdminApi | -                        |
| - [dump](#dump )           | No      | object | No         | In #/definitions/Dump     | -                        |
| - [sync](#sync )           | No      | object | No         | In #/definitions/Sync     | -                        |
|                            |         |        |            |                           |                          |

## <a name="url"></a>1. [Optional] Property `.shopware-project.yml > url`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** URL to Shopware instance

## <a name="admin_api"></a>2. [Optional] Property `.shopware-project.yml > admin_api`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/AdminApi                                                    |
|                           |                                                                           |

| Property                                     | Pattern | Type   | Deprecated | Definition | Title/Description            |
| -------------------------------------------- | ------- | ------ | ---------- | ---------- | ---------------------------- |
| - [client_id](#admin_api_client_id )         | No      | string | No         | -          | Client ID of integration     |
| - [client_secret](#admin_api_client_secret ) | No      | string | No         | -          | Client Secret of integration |
| - [username](#admin_api_username )           | No      | string | No         | -          | Username of admin user       |
| - [password](#admin_api_password )           | No      | string | No         | -          | Password of admin user       |
|                                              |         |        |            |            |                              |

### <a name="admin_api_client_id"></a>2.1. [Optional] Property `.shopware-project.yml > admin_api > client_id`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** Client ID of integration

### <a name="admin_api_client_secret"></a>2.2. [Optional] Property `.shopware-project.yml > admin_api > client_secret`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** Client Secret of integration

### <a name="admin_api_username"></a>2.3. [Optional] Property `.shopware-project.yml > admin_api > username`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** Username of admin user

### <a name="admin_api_password"></a>2.4. [Optional] Property `.shopware-project.yml > admin_api > password`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** Password of admin user

## <a name="dump"></a>3. [Optional] Property `.shopware-project.yml > dump`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/Dump                                                        |
|                           |                                                                           |

| Property                    | Pattern | Type            | Deprecated | Definition | Title/Description |
| --------------------------- | ------- | --------------- | ---------- | ---------- | ----------------- |
| - [rewrite](#dump_rewrite ) | No      | object          | No         | -          | -                 |
| - [nodata](#dump_nodata )   | No      | array of string | No         | -          | -                 |
| - [ignore](#dump_ignore )   | No      | array of string | No         | -          | -                 |
| - [where](#dump_where )     | No      | object          | No         | -          | -                 |
|                             |         |                 |            |            |                   |

### <a name="dump_rewrite"></a>3.1. [Optional] Property `.shopware-project.yml > dump > rewrite`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

### <a name="dump_nodata"></a>3.2. [Optional] Property `.shopware-project.yml > dump > nodata`

| Type                      | `array of string`                                                         |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be    | Description |
| ---------------------------------- | ----------- |
| [nodata items](#dump_nodata_items) | -           |
|                                    |             |

#### <a name="autogenerated_heading_2"></a>3.2.1. .shopware-project.yml > dump > nodata > nodata items

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

### <a name="dump_ignore"></a>3.3. [Optional] Property `.shopware-project.yml > dump > ignore`

| Type                      | `array of string`                                                         |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be    | Description |
| ---------------------------------- | ----------- |
| [ignore items](#dump_ignore_items) | -           |
|                                    |             |

#### <a name="autogenerated_heading_3"></a>3.3.1. .shopware-project.yml > dump > ignore > ignore items

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

### <a name="dump_where"></a>3.4. [Optional] Property `.shopware-project.yml > dump > where`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

## <a name="sync"></a>4. [Optional] Property `.shopware-project.yml > sync`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/Sync                                                        |
|                           |                                                                           |

| Property                                | Pattern | Type  | Deprecated | Definition | Title/Description |
| --------------------------------------- | ------- | ----- | ---------- | ---------- | ----------------- |
| - [config](#sync_config )               | No      | array | No         | -          | -                 |
| - [theme](#sync_theme )                 | No      | array | No         | -          | -                 |
| - [mail_template](#sync_mail_template ) | No      | array | No         | -          | -                 |
| - [entity](#sync_entity )               | No      | array | No         | -          | -                 |
|                                         |         |       |            |            |                   |

### <a name="sync_config"></a>4.1. [Optional] Property `.shopware-project.yml > sync > config`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be      | Description |
| ------------------------------------ | ----------- |
| [SyncConfigItem](#sync_config_items) | -           |
|                                      |             |

#### <a name="autogenerated_heading_4"></a>4.1.1. .shopware-project.yml > sync > config > SyncConfigItem

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/SyncConfigItem                                              |
|                           |                                                                           |

| Property                                             | Pattern | Type   | Deprecated | Definition | Title/Description      |
| ---------------------------------------------------- | ------- | ------ | ---------- | ---------- | ---------------------- |
| - [sales_channel](#sync_config_items_sales_channel ) | No      | string | No         | -          | Sales Channel to apply |
| + [settings](#sync_config_items_settings )           | No      | object | No         | -          | -                      |
|                                                      |         |        |            |            |                        |

##### <a name="sync_config_items_sales_channel"></a>4.1.1.1. Property `.shopware-project.yml > sync > config > System Config Sync > sales_channel`

**Title:** Sales Channel to apply

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_config_items_settings"></a>4.1.1.2. Property `.shopware-project.yml > sync > config > System Config Sync > settings`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

### <a name="sync_theme"></a>4.2. [Optional] Property `.shopware-project.yml > sync > theme`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be      | Description |
| ------------------------------------ | ----------- |
| [ThemeConfigItem](#sync_theme_items) | -           |
|                                      |             |

#### <a name="autogenerated_heading_5"></a>4.2.1. .shopware-project.yml > sync > theme > ThemeConfigItem

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/ThemeConfigItem                                             |
|                           |                                                                           |

| Property                                  | Pattern | Type   | Deprecated | Definition | Title/Description |
| ----------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| - [name](#sync_theme_items_name )         | No      | string | No         | -          | -                 |
| - [settings](#sync_theme_items_settings ) | No      | object | No         | -          | -                 |
|                                           |         |        |            |            |                   |

##### <a name="sync_theme_items_name"></a>4.2.1.1. Property `.shopware-project.yml > sync > theme > Theme Config Sync > name`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_theme_items_settings"></a>4.2.1.2. Property `.shopware-project.yml > sync > theme > Theme Config Sync > settings`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

### <a name="sync_mail_template"></a>4.3. [Optional] Property `.shopware-project.yml > sync > mail_template`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be               | Description |
| --------------------------------------------- | ----------- |
| [MailTemplateItem](#sync_mail_template_items) | -           |
|                                               |             |

#### <a name="autogenerated_heading_6"></a>4.3.1. .shopware-project.yml > sync > mail_template > MailTemplateItem

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/MailTemplateItem                                            |
|                           |                                                                           |

| Property                                                  | Pattern | Type   | Deprecated | Definition | Title/Description |
| --------------------------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| - [id](#sync_mail_template_items_id )                     | No      | string | No         | -          | -                 |
| - [translations](#sync_mail_template_items_translations ) | No      | array  | No         | -          | -                 |
|                                                           |         |        |            |            |                   |

##### <a name="sync_mail_template_items_id"></a>4.3.1.1. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > id`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_mail_template_items_translations"></a>4.3.1.2. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be                                             | Description |
| --------------------------------------------------------------------------- | ----------- |
| [MailTemplateItemTranslation](#sync_mail_template_items_translations_items) | -           |
|                                                                             |             |

##### <a name="autogenerated_heading_7"></a>4.3.1.2.1. .shopware-project.yml > sync > mail_template > Mail Template Sync > translations > MailTemplateItemTranslation

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/MailTemplateItemTranslation                                 |
|                           |                                                                           |

| Property                                                                       | Pattern | Type           | Deprecated | Definition | Title/Description |
| ------------------------------------------------------------------------------ | ------- | -------------- | ---------- | ---------- | ----------------- |
| - [language](#sync_mail_template_items_translations_items_language )           | No      | string         | No         | -          | -                 |
| - [sender_name](#sync_mail_template_items_translations_items_sender_name )     | No      | string         | No         | -          | -                 |
| - [subject](#sync_mail_template_items_translations_items_subject )             | No      | string         | No         | -          | -                 |
| - [html](#sync_mail_template_items_translations_items_html )                   | No      | string         | No         | -          | -                 |
| - [plain](#sync_mail_template_items_translations_items_plain )                 | No      | string         | No         | -          | -                 |
| - [custom_fields](#sync_mail_template_items_translations_items_custom_fields ) | No      | object or null | No         | -          | -                 |
|                                                                                |         |                |            |            |                   |

##### <a name="sync_mail_template_items_translations_items_language"></a>4.3.1.2.1.1. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations > Mail Template Single Translation > language`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_mail_template_items_translations_items_sender_name"></a>4.3.1.2.1.2. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations > Mail Template Single Translation > sender_name`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_mail_template_items_translations_items_subject"></a>4.3.1.2.1.3. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations > Mail Template Single Translation > subject`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_mail_template_items_translations_items_html"></a>4.3.1.2.1.4. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations > Mail Template Single Translation > html`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_mail_template_items_translations_items_plain"></a>4.3.1.2.1.5. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations > Mail Template Single Translation > plain`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_mail_template_items_translations_items_custom_fields"></a>4.3.1.2.1.6. Property `.shopware-project.yml > sync > mail_template > Mail Template Sync > translations > Mail Template Single Translation > custom_fields`

| Type                      | `object or null`                                                          |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

### <a name="sync_entity"></a>4.4. [Optional] Property `.shopware-project.yml > sync > entity`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be      | Description |
| ------------------------------------ | ----------- |
| [EntitySyncItem](#sync_entity_items) | -           |
|                                      |             |

#### <a name="autogenerated_heading_8"></a>4.4.1. .shopware-project.yml > sync > entity > EntitySyncItem

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/EntitySyncItem                                              |
|                           |                                                                           |

| Property                                 | Pattern | Type   | Deprecated | Definition | Title/Description |
| ---------------------------------------- | ------- | ------ | ---------- | ---------- | ----------------- |
| + [entity](#sync_entity_items_entity )   | No      | string | No         | -          | -                 |
| - [exists](#sync_entity_items_exists )   | No      | array  | No         | -          | -                 |
| + [payload](#sync_entity_items_payload ) | No      | object | No         | -          | API payload       |
|                                          |         |        |            |            |                   |

##### <a name="sync_entity_items_entity"></a>4.4.1.1. Property `.shopware-project.yml > sync > entity > Entity Sync > entity`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="sync_entity_items_exists"></a>4.4.1.2. Property `.shopware-project.yml > sync > entity > Entity Sync > exists`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be                     | Description |
| --------------------------------------------------- | ----------- |
| [EntitySyncFilter](#sync_entity_items_exists_items) | -           |
|                                                     |             |

##### <a name="autogenerated_heading_9"></a>4.4.1.2.1. .shopware-project.yml > sync > entity > Entity Sync > exists > EntitySyncFilter

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Defined in**            | #/definitions/EntitySyncFilter                                            |
|                           |                                                                           |

| Property                                                | Pattern | Type                                    | Deprecated | Definition | Title/Description |
| ------------------------------------------------------- | ------- | --------------------------------------- | ---------- | ---------- | ----------------- |
| + [type](#sync_entity_items_exists_items_type )         | No      | enum (of string)                        | No         | -          | filter type       |
| + [field](#sync_entity_items_exists_items_field )       | No      | string                                  | No         | -          | field             |
| - [value](#sync_entity_items_exists_items_value )       | No      | string, integer, array, boolean or null | No         | -          | value             |
| - [operator](#sync_entity_items_exists_items_operator ) | No      | enum (of string)                        | No         | -          | -                 |
| - [queries](#sync_entity_items_exists_items_queries )   | No      | array                                   | No         | -          | -                 |
|                                                         |         |                                         |            |            |                   |

| All of(Requirement)                                |
| -------------------------------------------------- |
| [item 0](#sync_entity_items_exists_items_allOf_i0) |
|                                                    |

##### <a name="sync_entity_items_exists_items_allOf_i0"></a>4.4.1.2.1.1. Property `.shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > allOf > item 0`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="autogenerated_heading_10"></a>4.4.1.2.1.1.1. If (type = "multi")

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

##### <a name="autogenerated_heading_11"></a>4.4.1.2.1.1.1.1. The following properties are required
* type
* queries

##### <a name="sync_entity_items_exists_items_type"></a>4.4.1.2.1.2. Property `.shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > type`

| Type                      | `enum (of string)`                                                        |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** filter type

Must be one of:
* "equals"
* "multi"
* "contains"
* "prefix"
* "suffix"
* "not"
* "range"
* "until"
* "equalsAll"
* "equalsAny"

##### <a name="sync_entity_items_exists_items_field"></a>4.4.1.2.1.3. Property `.shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > field`

| Type                      | `string`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** field

##### <a name="sync_entity_items_exists_items_value"></a>4.4.1.2.1.4. Property `.shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > value`

| Type                      | `string, integer, array, boolean or null`                                 |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** value

##### <a name="sync_entity_items_exists_items_operator"></a>4.4.1.2.1.5. Property `.shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > operator`

| Type                      | `enum (of string)`                                                        |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

Must be one of:
* "AND"
* "OR"
* "XOR"

##### <a name="sync_entity_items_exists_items_queries"></a>4.4.1.2.1.6. Property `.shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > queries`

| Type                      | `array`                                                                   |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

|                      | Array restrictions |
| -------------------- | ------------------ |
| **Min items**        | N/A                |
| **Max items**        | N/A                |
| **Items unicity**    | False              |
| **Additional items** | False              |
| **Tuple validation** | See below          |
|                      |                    |

| Each item of this array must be                                   | Description |
| ----------------------------------------------------------------- | ----------- |
| [EntitySyncFilter](#sync_entity_items_exists_items_queries_items) | -           |
|                                                                   |             |

##### <a name="autogenerated_heading_12"></a>4.4.1.2.1.6.1. .shopware-project.yml > sync > entity > Entity Sync > exists > Entity Sync Filter > queries > EntitySyncFilter

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
| **Same definition as**    | [Entity Sync Filter](#sync_entity_items_exists_items)                     |
|                           |                                                                           |

##### <a name="sync_entity_items_payload"></a>4.4.1.3. Property `.shopware-project.yml > sync > entity > Entity Sync > payload`

| Type                      | `object`                                                                  |
| ------------------------- | ------------------------------------------------------------------------- |
| **Additional properties** | [[Any type: allowed]](# "Additional Properties of any type are allowed.") |
|                           |                                                                           |

**Description:** API payload

----------------------------------------------------------------------------------------------------------------------------
