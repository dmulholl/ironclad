
# Ironclad

A command line utility for creating and managing encrypted password databases.

    $ ironclad --help

    Usage: ironclad [FLAGS] [COMMAND]

      A utility for creating and managing encrypted password databases.

    Flags:
      --help            Print the application's help text and exit.
      --version         Print the application's version number and exit.

    Commands:
      add               Add a new entry to a database.
      config            Set or print a configuration option.
      decrypt           Decrypt a file.
      delete            Delete entries from a database.
      dump              Dump a database's internal JSON data store.
      edit              Edit an existing database entry.
      encrypt           Encrypt a file.
      export            Export entries from a database.
      gen               Generate a new random password.
      import            Import entries into a database.
      list              List database entries.
      new               Create a new database.
      pass              Copy a password to the clipboard.
      purge             Purge deleted entries from a database.
      tags              List database tags.
      user              Copy a username to the clipboard.

    Command Help:
      help <command>    Print the specified command's help text and exit.

Database files are encrypted using industry-standard cryptographic protocols.

* Data is encrypted using 256-bit AES in CBC mode.
* Padding is performed using the PKCS#7 padding scheme.
* Authentication is performed using the HMAC-SHA-256 algorithm.
* Encryption keys are generated using the PBKDF2 key derivation algorithm with an SHA-256 hash.

This application is a cross-platform utility written in a high-level, garbage-collected language. It has *not* been hardened against system-local threats, e.g. malicious code running with user-level privileges on the user's system, or attackers with physical access to the user's hardware.
