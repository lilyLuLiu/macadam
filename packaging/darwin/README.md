## How to build

```sh
$ make NO_CODESIGN=1 pkginstaller

# or to create signed pkg
$ make CODESIGN_IDENTITY=<ID> PRODUCTSIGN_IDENTITY=<ID> pkginstaller

# or to prepare a signed and notarized pkg for release
$ make CODESIGN_IDENTITY=<ID> PRODUCTSIGN_IDENTITY=<ID> NOTARIZE_USERNAME=<appleID> NOTARIZE_PASSWORD=<appleID-password> NOTARIZE_TEAM=<team-id> notarize
```

The generated pkg will be written to `out/macadam-macos-installer-universal.pkg`.
Currently the pkg installs `macadam`, `vfkit`, and `gvproxy` to `/opt/macadam`

## Uninstalling

```sh
$ sudo pkgutil --forget com.redhat.macadam
$ sudo rm -rf /opt/macadam
```