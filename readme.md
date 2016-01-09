
# Ironclad

A prototype command line utility for managing Ironclad password databases.

Usage:

    Usage: iron [FLAGS] [COMMAND]

      Command line utility for managing Ironclad password databases.

    Flags:
      --help            Print the application's help text and exit.
      --version         Print the application's version number and exit.

    Commands:
      add               Add a new entry to a database.
      delete            Delete an entry from a database.
      dump              Dump a database's internal JSON data store.
      edit              Edit an existing database entry.
      export            Export data from a database.
      gen               Generate a random password.
      import            Import data into a database.
      list              List database entries.
      new               Create a new database.
      pass              Print a password.
      purge             Purge deleted entries from a database.
      tags              List database tags.
      user              Print a username.

    Command Help:
      help <command>    Print the specified command's help text and exit.


## Security

Database files are encrypted using industry-standard cryptographic protocols.

* Data is encrypted using 256-bit AES in CBC mode.
* Padding is performed using the PKCS#7 padding scheme.
* Authentication is performed using the HMAC-SHA-256 algorithm.

This application is a cross-platform utility written in a high-level, garbage-collected language. It has *not* been hardened against system-local threats, e.g. malicious code running with user-level privileges on the user's system, or attackers with physical access to the user's hardware.

In particular, the password-caching feature means that malicious software with read-access to the user's home directory will have access to the master password.


## License

MIT
