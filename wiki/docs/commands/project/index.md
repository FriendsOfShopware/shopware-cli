---
title: Project Commands
weight: 30
---

## shopware-cli project create [folder] [version]

Create a new Shopware 6 project from the choosen version

Arguments:

* `folder` - **Required**: Folder name for the installation
* `version` - Version to install

## shopware-cli project admin-build

Builds the Administration with all installed extensions

## shopware-cli project storefront-build

Builds the Storefront with all installed extensions

## shopware-cli project worker

Starts the Shopware worker in background and tails the log

Arguments:

* Worker amount - `shopware-cli project worker 5` starts 5 workers

## shopware-cli project dump [database]

Dumps the MySQL database as SQL

Parameters:

* `--host` - MySQL Host (default: 127.0.0.1)
* `--port` - MySQL Port (default: 3306)
* `--username` - MySQL Username (default: root)
* `--password` - MySQL Password (default: root)
* `--output` - Output file (default: `dump.sql`)
* `--clean` - Ignores content of following tables: `cart`, `customer_recovery`, `dead_message`, `enqueue`, `increment`, `elasticsearch_index_task`, `log_entry`, `message_queue_stats`, `notification`, `payment_token`, `refresh_token`, `version`, `version_commit`, `version_commit_data`, `webhook_event_log`
* `--skip-lock-tables` - Skips locking of tables
* `--anonymize` - Anonymizes known user data tables. See https://github.com/FriendsOfShopware/shopware-cli/blob/main/cmd/project/project_dump.go#L61 for the list

Examples:

- `shopware-cli project dump sw6 --host 127.0.0.1 --username root --password root --clean --anonymize`

## shopware-cli project admin-api [method] [path]

Run authentificated curl against the admin api

Arguments:

* `method` - **Required:** HTTP method
* `path` - **Required:** HTTP path

Parameters:

* `--output-token` - Outputs only the access token

Examples:

- `shopware-cli project admin-api POST "/search/tax" -- -d '{"limit": 1}' -H 'Accept: application/json' -H 'Content-Type: application/json'`
