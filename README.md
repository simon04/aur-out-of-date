aur-out-of-date
==========

Iterates through a user's AUR (Arch User Repository) packages, and determines out-of-date packages w.r.t. their upstream version.

Installation
------------

```sh
$ go get github.com/simon04/aur-out-of-date
```

The tool is also available in AUR: [aur-out-of-date](https://aur.archlinux.org/packages/aur-out-of-date/)

Usage
-----

```
$ aur-out-of-date -user simon04
[OUT-OF-DATE] [python-mwclient] Package python-mwclient should be updated from 0.8.6-1 to 0.8.7-1
[OUT-OF-DATE] [nodejs-osmtogeojson] Package nodejs-osmtogeojson should be updated from 2.2.12-1 to 3.0.0
[OUT-OF-DATE] [spectre-meltdown-checker] Package spectre-meltdown-checker should be updated from 0.31-1 to 0.32

$ aur-out-of-date -pkg caddy -pkg dep -pkg aur-out-of-date
[UP-TO-DATE]  [aur-out-of-date] Package aur-out-of-date 0.4.0-1 matches upstream version 0.4.0
[UP-TO-DATE]  [caddy] Package caddy 0.10.10-3 matches upstream version 0.10.10
[UP-TO-DATE]  [dep] Package dep 0.3.2-2 matches upstream version 0.3.2
```

Principle
---------

For each package, the upstream URL and/or source URL is matched against supported platforms. For those platforms the latest release is obtained via an API/HTTP call.

* `github.com` or `github.io` → https://github.com/…/…/releases.atom
* `registry.npmjs.org` → https://registry.npmjs.org/-/package/…/dist-tags
* `pypi.python.org` or `files.pythonhosted.org` → https://pypi.python.org/pypi/…/json

License
-------

GNU General Public License v3.0
