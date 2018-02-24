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
$ aur-out-of-date
Usage of aur-out-of-date:
  -devel
        Check -git/-svn/-hg packages
  -flag
        Flag out-of-date on AUR
  -local
        Local .SRCINFO files
  -pkg
        AUR package name(s)
  -statistics
        Print summary statistics
  -user string
        AUR username
```

AUR packages can be obtained …

- for a given AUR user (using `-user simon04`; specify `-devel` to include VCS packages), or
- from a list of packages via [AUR RPC](https://aur.archlinux.org/rpc.php) (using `-pkg package1 package2 …`), or
- from local `.SRCINFO` files (using `-local packages/*/.SRCINFO`).

```
$ aur-out-of-date -user simon04
[OUT-OF-DATE] [python-mwclient] Package python-mwclient should be updated from 0.8.6-1 to 0.8.7-1
[OUT-OF-DATE] [nodejs-osmtogeojson] Package nodejs-osmtogeojson should be updated from 2.2.12-1 to 3.0.0
[OUT-OF-DATE] [spectre-meltdown-checker] Package spectre-meltdown-checker should be updated from 0.31-1 to 0.32

$ aur-out-of-date -user simon04 -devel
[UP-TO-DATE]  [ocproxy-git] Package ocproxy-git 1.60.r8.g8f15425-3 matches upstream version 1.60

$ aur-out-of-date -pkg caddy dep aur-out-of-date
[UP-TO-DATE]  [aur-out-of-date] Package aur-out-of-date 0.4.0-1 matches upstream version 0.4.0
[UP-TO-DATE]  [caddy] Package caddy 0.10.10-3 matches upstream version 0.10.10
[UP-TO-DATE]  [dep] Package dep 0.3.2-2 matches upstream version 0.3.2

$ aur-out-of-date -local packages/*/.SRCINFO
```

Summary statistics can be enabled using `-statistics`. The option `-flag` flags out-of-date packages on AUR after a user prompt: "Should the package … be flagged out-of-date?"

Principle
---------

For each package, the upstream URL and/or source URL is matched against supported platforms. For those platforms the latest release is obtained via an API/HTTP call.

* `github.com` or `github.io` → https://github.com/…/…/releases.atom
* `registry.npmjs.org` → https://registry.npmjs.org/-/package/…/dist-tags
* `pypi.python.org` or `files.pythonhosted.org` → https://pypi.python.org/pypi/…/json
* `search.cpan.org` or `search.mcpan.org` -> https://fastapi.metacpan.org/v1/release/…

License
-------

GNU General Public License v3.0
