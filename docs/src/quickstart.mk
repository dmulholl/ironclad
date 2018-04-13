---
title: Getting Started
---

Download the latest version of Ironclad from the [releases][] page, unzip the package and place the binary somewhere on your system's [$PATH][].

[releases]: @root/releases//
[$PATH]: https://alistapart.com/article/the-path-to-enlightenment



### Initialize a new database

Use the `init` command to initialize a new password database. This command takes a single argument specifying the location of the database file.

    $ ironclad init topsecret/passwords.idb

Ironclad will create a new database file at the specified location.
You can call the file anything you like --- Ironclad doesn't require any particular file extension for databases. (In fact, database files have no special markers and are indistinguishable from random data.)

You'll be prompted to enter a master password which will be used to encrypt the file. You'll be prompted to re-enter this password when you next run a command on the database; Ironclad will then default to caching the password in memory for 15 minutes from its last use.



### Add entries

Use the `add` command to add a new entry to your database.

    $ ironclad add

You'll be prompted to supply values for the entry's fields --- simply press return to leave any unwanted field empty.



### List entries

Use the `list` command to see a list of all the database's entries.

    $ ironclad list

Use the `show` command to view the full details of each entry in the list.

    $ ironclad show

(The `show` comand is an alias for `list --verbose`.)



### Copy a password or username

Use the `pass` command to copy a password to the clipboard. The clipboard is automatically overwritten after a default timeout of 10 seconds.

    $ ironclad pass entry

You can specify the entry by its ID or by any unique case-insensitive substring of its title.

The `user` command works identically for usernames.



### Command line help

Use the `--help` flag to view Ironclad's command line help, including a list of all available commands.

    $ ironclad --help

Use the `help` command to view detailed help for any individual command.

    $ ironclad help <command>
