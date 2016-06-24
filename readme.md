
# Ironclad

A prototype command line utility for managing Ironclad password databases.


## Security

Database files are encrypted using industry-standard cryptographic protocols.

* Data is encrypted using 256-bit AES in CBC mode.
* Padding is performed using the PKCS#7 padding scheme.
* Authentication is performed using the HMAC-SHA-256 algorithm.

This application is a cross-platform utility written in a high-level, garbage-collected language. It has *not* been hardened against system-local threats, e.g. malicious code running with user-level privileges on the user's system, or attackers with physical access to the user's hardware.

In particular, the password-caching feature means that malicious software with read-access to the user's home directory will have access to the master password.


## License

MIT
