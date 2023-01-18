# Buscarron v1.3.1 (2023-01-18)

### Features :sparkles:

* Shared rate limit across forms

### Bugfixes :bug:

* Do not ban internal IPs

### Misc :zzz:

* Updated deps
* Enabled vendoring

# Buscarron v1.3.0 (2022-10-25)

### Features :sparkles:

* Advanced email validation, using using actual connection to the sender's SMTP address
* Email spamlist with wildcards support
* Permanent ban list, based on IP hashes / IDs

### Breaking changes :warning:

* ARMv7 (arm32/aarch32) architecture is deprecated

### Misc :zzz:

* Updated deps

# Buscarron v1.2.1 (2022-08-23)

### Misc :zzz:

* Updated deps

# Buscarron v1.2.0 (2022-06-20)

### Features :sparkles:

* Spam checker can work with with email localpart list
* Add optional domain validation on the forms level
* Add `NOENCRYPTION` option

### Bugfixes :bug:

* Added CORS headers

### Misc :zzz:

* Updated deps

# Buscarron v1.1.0 (2022-05-10)

### Features :sparkles:

* Email confirmation after form submissions with [Postmark](https://postmarkapp.com)
* [Templates](https://pkg.go.dev/text/template) support for redirect URLs
* Automatically ban spammers, scanners and bots

# Buscarron v1.0.0 (2022-04-23)

### Features :sparkles:

* Initial release

### Bugfixes :bug:

N/A _that's the first version, bugs are creating on that step!_

### Breaking changes :warning:

Buscarron has been created.
