---
title: Changelog
---


### 1.5.0.dev

* The `add` command can now automatically generate a password for a new entry.

* Rename the `setpass` command to `masterpass`.

* The default character pool for generating passwords now includes symbols.



### 1.4.2

* The config directory is now correctly stored in the user's `%appdata%` directory on Windows.

* Turn on support for ANSI color codes in the Windows console.



### 1.4.0

* Add colored terminal output.



### 1.3.0

* Add a `restore` command for restoring inactive entries.

* Add a `--deleted` flag to the `list` command for displaying inactive entries.

* Add a confirmation step to the `purge` command.



### 1.2.0

* Add a `url` command for copying an entry's url to the clipboard.

* Switch to using Go's official `dep` tool to manage vendored dependencies.



### 1.1.0

* Add a `setpass` command for changing a database's master password.

* Support importing from standard input.



### 1.0.0

* First stable release.
