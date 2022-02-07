---
title: Project Commands
weight: 30
---

## shopware-cli project create

Create a new Shopware 6 project from the choosen version

## shopware-cli project admin-build

Builds the Administration with all installed extensions

## shopware-cli project storefront-build

Builds the Storefront with all installed extensions

## shopware-cli project worker

Starts the Shopware worker in background and tails the log

Arguments:

* Worker amount - `shopware-cli project worker 5` starts 5 workers

## shopware-cli project dump

Dumps the MySQL database as SQL

Parameters:

* `--host` - MySQL Host (default: 127.0.0.1)
* `--port` - MySQL Port (default: 3306)
* `--username` - MySQL Username (default: root)
* `--password` - MySQL Password (default: root)
* `--output` - Output file (default: `dump.sql`)
* `--clean` - Ignores content of following tables: `cart`, `customer_recovery`, `dead_message`, `enqueue`, `increment`, `elasticsearch_index_task`, `log_entry`, `message_queue_stats`, `notification`, `payment_token`, `refresh_token`, `version`, `version_commit`, `version_commit_data`, `webhook_event_log`
* `--skip-lock-tables` - Skips locking of tables