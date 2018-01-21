aur-out-of-date
==========

Iterates through a user's AUR (Arch User Repository) packages, and determines out-of-date packages w.r.t. their upstream version:

```sh
$ go get github.com/mikkeloscar/aur
$ go get github.com/mikkeloscar/gopkgbuild
$ go get github.com/mmcdole/gofeed
$ go get github.com/simon04/aur-out-of-date

$ aur-out-of-date -user simon04
Package python-mwclient should be updated from 0.8.6-1 to 0.8.7-1
Package python2-mwclient should be updated from 0.8.6-1 to 0.8.7-1
Package nodejs-osmtogeojson should be updated from 2.2.12-1 to 3.0.0-beta.3

$ aur-out-of-date -pkg caddy -pkg qgis
```

Principle
---------

For each package, the upstream URL and/or source URL is matched against supported platforms. For those platforms the latest release is obtained via an API/HTTP call.

* `github.com` or `github.io` → https://github.com/…/…/releases.atom

License
-------

GNU General Public License v3.0
