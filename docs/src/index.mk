---
title: Home
meta title: Ironclad &mdash; a command line password manager
meta description: >
    Ironclad is a command line utility for creating and managing encrypted
    password databases.
---

Ironclad is a command line utility for creating and managing encrypted password databases.



### Download

* [Linux][]
* [Mac][]
* [Windows][]

See the [releases](@root/releases//) page for file hashes.



### Usage

Run `ironclad --help` to view the application's command line help:

    Usage: ironclad [FLAGS] [COMMAND]

      Ironclad is a command line password manager.

    Flags:
      -h, --help        Print the help text and exit.
      -v, --version     Print the version number and exit.

    Basic Commands:
      add               Add a new entry to a database.
      delete            Delete entries from a database.
      edit              Edit an existing database entry.
      gen               Generate a new random password.
      init              Initialize a new password database.
      list              List database entries.
      new               Alias for 'add'.
      pass              Copy a password to the clipboard.
      show              Alias for 'list --verbose'.
      user              Copy a username to the clipboard.

    Additional Commands:
      config            Set or print a configuration option.
      decrypt           Decrypt a file.
      dump              Dump a database's JSON data store.
      encrypt           Encrypt a file.
      export            Export entries from a database.
      import            Import entries into a database.
      purge             Purge deleted entries from a database.
      setpass           Change a database's master password.
      tags              List database tags.

    Command Help:
      help <command>    Print a command's help text.

Run `ironclad help <command>` to view the help text for a specific command.

The [quickstart guide](@root/quickstart//) is a short tutorial for first-time users.



### Source

Ironclad is written in Go. If you have a Go compiler installed
you can run:

    $ go get github.com/dmulholland/ironclad/ironclad

This will download, compile, and install the latest version of the application
to your `$GOPATH/bin` directory.

You can find the source files on [Github][].

[github]: https://github.com/dmulholland/ironclad




### Security

Database files are encrypted using industry-standard cryptographic protocols.

* Data is encrypted using 256-bit AES in CBC mode.
* Padding is performed using the PKCS #7 padding scheme.
* Authentication is performed using the HMAC-SHA-256 protocol.
* Encryption keys are generated using 10,000 rounds of the PBKDF2 key derivation algorithm with an SHA-256 hash.

Encrypted files have no special markers and are indistinguishable from random data.

Note that the application itself is a cross-platform utility written in a high-level, garbage-collected language. It has *not* been hardened against system-local threats, e.g. malicious code running with user-level privileges on the user's system, or adversaries with physical access to the user's hardware.



### Password Caching

Ironclad caches the master password in memory for a default period of 15 minutes from its last use. You can set a custom timeout using the `config` command:

    $ ironclad config timeout <minutes>

Setting the timeout to `0` will disable caching altogether.



### File Encryption

Ironclad doubles as a simple file encryption utility using the `encrypt` and `decrypt` commands. Files are encrypted using the same 256-bit AES protocol as password databases. Original files are unaffected by either encryption or decryption.



### Rationale

I built this cross-platform utility as a prototype implementation of Ironclad's core idea --- an open-source password manager organised around a simple JSON data store.

Complexity is the enemy of security, so Ironclad is as uncomplicated as possible. A password database is a simple JSON file which you can view using the `dump` command:

    $ ironclad dump

This file is encrypted using 256-bit AES, an industry-standard protocol supported on all platforms and across all programming languages.

The Ironclad application itself is a cross-platform prototype. However, alternative native clients should be straightforward to implement and can take better advantage of the built-in security features offered by specific operating systems.



### Alternative Implementations

* Ryan Wynn's [llave](https://github.com/rwynn/llave) provides desktop and mobile interfaces for Ironclad databases.



### License

Ironclad is released under an MIT license.



[linux]:
  https://github.com/dmulholland/ironclad/releases/download/1.1.0/linux.zip

[mac]:
  https://github.com/dmulholland/ironclad/releases/download/1.1.0/mac.zip

[windows]:
  https://github.com/dmulholland/ironclad/releases/download/1.1.0/windows.zip
