# crl-updater

CRL update daemon

## Installation

Check the [releases](https://github.com/na4ma4/crl-updater/releases) page.

## Configuration

`/etc/ssl/crl-updater.toml`

```toml
## Name of the collection of source/target/actions.
[main]
## Source of the CRL.
source="https://companyca.example.com/ca.crl"
## Target file containing the CRL.
target="/etc/service/ssl/ca-crl.pem"

## Name of the colletion plus ".actions"
[main.actions]
## This command is run before the target is checked.
precheck="/opt/bin/command-to-run-before-checking-target.sh"

## This command is run before the target is updated only if it requires updating.
preinstall="/opt/bin/command-to-run-before-updating-target.sh"

## This command is run after the target is updated and only if the target is updated.
postinstall="/opt/bin/command-to-run-after-updating-target.sh"

## This command is run regardless of any errors or if the target was updated.
post="/opt/bin/command-to-run-after-checking-and_or-updating.sh"
```
