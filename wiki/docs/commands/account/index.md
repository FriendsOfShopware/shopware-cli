---
title: Account Commands
weight: 20
---

### account login

This command can be used to log in into your Shopware account. If you are in multiple companies, see `Account Company Use` command

### account logout

Logout from your Account

### account company list

List all your companies

### account company use

Switch the active company.

Parameters:

* Company ID - Can be obtained by \`account company list\`

### account producer info

Lists some basic information about the logged in producer

### account producer extension list

Lists all your extensions in the account

### account producer extension create

Creates an extension

Parameters:

* name - Your extension name
* generation - classic, platform, apps, themes

### account producer extension delete

Deletes an extension

Parameters:

* name - Your extension name

### account producer extension info pull

Downloads the store page information to the given extension

Parameters:

* path - Extension folder path

### account producer extension info push

Uploads the local store page information

Parameters:

* path - Extension folder path
