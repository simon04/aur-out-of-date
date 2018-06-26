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
  -config string
        Config file (default "$XDG_CONFIG_HOME/aur-out-of-date/config.json")
  -devel
        Check -git/-svn/-hg packages
  -flag
        Flag out-of-date on AUR
  -json
        Generate JSON Text Sequences (RFC 7464)
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

The output can be switched to a machine-readable format – [JavaScript Object Notation (JSON) Text Sequences](https://tools.ietf.org/html/rfc7464) – using `-json`.

```json
$ aur-out-of-date -json -pkg nodejs-osmtogeojson spectre-meltdown-checker
{"type":"package","name":"nodejs-osmtogeojson","message":"Package nodejs-osmtogeojson should be updated from 2.2.12-1 to 3.0.0","version":"2.2.12-1","upstream":"3.0.0","status":"OUT-OF-DATE"}
{"type":"package","name":"spectre-meltdown-checker","message":"Package spectre-meltdown-checker 0.35-1 matches upstream version 0.35","version":"0.35-1","upstream":"0.35","status":"UP-TO-DATE"}
```

Summary statistics can be enabled using `-statistics`.

The option `-flag` flags out-of-date packages on AUR after a user prompt: "Should the package … be flagged out-of-date?"

The tool `aur-out-of-date` exists with code `4` if at least one out-of-date package has been found.

Principle
---------

For each package, the upstream URL and/or source URL is matched against supported platforms. For those platforms the latest release is obtained via an API/HTTP call.

* `github.com` or `github.io` → http://api.github.com/repos/…/…/releases/latest (provide a [personal access token](https://github.com/settings/tokens) in the environment variable `GITHUB_TOKEN` for [higher request limits](https://developer.github.com/v3/#rate-limiting))
* `registry.npmjs.org` → https://registry.npmjs.org/-/package/…/dist-tags
* `pypi.python.org` or `files.pythonhosted.org` → https://pypi.python.org/pypi/…/json
* `search.cpan.org` or `search.mcpan.org` → https://fastapi.metacpan.org/v1/release/…
* `rubygems.org` or `gems.rubyforge.org` → https://rubygems.org/api/v1/versions/….json
* `gitlab.com` or any self-hosted GitLab instance → http://gitlab.com/api/v4/…/…/repository/tags (provide a [personal access token](https://github.com/settings/tokens) in the environment variable `GITLAB_TOKEN` for [higher request limits](https://docs.gitlab.com/ee/api/#oauth2-tokens))

Configuration
-------------

The tool reads a configuration file from `$XDG_CONFIG_HOME/aur-out-of-date/config.json`. This allows to ignore certain package versions from being reported as out-of-date. The string `"*"` acts as a placeholder for all versions.

```json
{
  "ignore": {
    "foo": ["*"],
    "osmtogeojson": ["3.0.0-beta.3", "3.0.0-rc.1"]
  }
}
```

Running `aur-out-of-date -pkg osmtogeojson` yields:

```
[UNKNOWN] [osmtogeojson][3.0.0b3-2] ignoring package upgrade to 3.0.0-beta.3
```


Related projects
----------------

* https://github.com/repology/repology
* https://github.com/lilydjwg/nvchecker

License
-------

GNU General Public License v3.0
