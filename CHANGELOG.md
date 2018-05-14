# Release History

## 1.3.0 (2018-05-14)

* GitHub: prefer tag_name over release name
* Clean version string for all providers (strip `v` prefix)

## 1.2.0 (2018-05-03)

* GitHub: support dots in repository names
* Add support for rubygems.org

## 1.1.0 (2018-03-29)

* NPM: support `@scoped/packages`
* GitHub: fall back to `tag_name` when release does not have a `name`

## 1.0.0 (2018-03-06)

* Use [GitHub releases API](https://developer.github.com/v3/repos/releases/) to skip pre-releases, release drafts
* Cache HTTP requests using `github.com/gregjones/httpcache`
* Provide machine-readable format: JSON Text Sequences ([RFC 7464](https://tools.ietf.org/html/rfc7464))
* Exit with code `4` if at least one out-of-date package has been found
* Fix error on Unicode characters in package version

## 0.8.0 (2018-02-24)

* Print summary statistics
* Flag AUR package out-of-date
* Fix checking huge number of packages

## 0.7.0 (2018-02-06)

* Read local .SRCINFO files
* Handle split packages correctly

## 0.6.0 (2018-01-27)

* Add flag to handle VCS packages only or skip them
* Add support for cpan.org

## 0.5.0 (2018-01-21)

* Add support for registry.npmjs.org
* Add support for pypi.python.org

## 0.4.0 (2018-01-21)

* Initial release including support for GitHub releases
