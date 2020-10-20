# Changelog

## 4.1.0

* Added `removeStroke` option.

## 4.0.0

* Option `path` was replaced with `paths`. Allows to specify multiple paths for URL resolve.
* Default `encode` function now escapes curly brackets.
* Dropped Node.js 6 support. Node.js 8.7.0 or greater is now required.
* Upgraded to PostCSS 7.

## 3.1.0

* Added `xmlns` option to add xmlns attribute if not present

## 3.0.0

* Upgrade to postcss 6

## 2.3.1

* Warn more descriptive messages

## 2.3.0

* Add parent prop in dependency message

## 2.2.0

* Add dependency message for bundlers

## 2.1.2

* Allow color functions in short syntax

## 2.1.1

* Customize `ENOENT` error
* Codebase refactoring

## 2.1.0

* Add `removeFill` option to force color in pure icons

## 2.0.1

* Stringify decl values

## 2.0.0

* Code and tests refactoring
* `encode` is not included in `transform` and runs before it
* `encode` is `function` or `false`
* `transform` can't return false
* Fix bug with quotes [#19](https://github.com/TrySound/postcss-inline-svg/issues/19)

## 1.4.0

* Support "=" separator in svg-load() definition

## 1.3.2

* Fix URI malformed error

## 1.3.1

* Add read-cache

## 1.3.0

* Refactoring
* Loaded files will be cached

## 1.2.1

* Fix relative path detecting

## 1.2.0

* Add transform option

## 1.1.1

* Improve documentation

## 1.1.0

* Fix ie strict data uri by adding custom url encode
* Add encode option

## 1.0.0

* Initial release
