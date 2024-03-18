# Changelog

## [ Upcoming Release - vx.x.x ]

### Changed/Added

### Fixed

### Deprecated

## [ 2024/03/18 - v0.1.7 ]

### Changed

* CHART: Enabled overriding of images to deploy
* CHART: Set default images to SHA256 instead of tags to improve security

## [ 2024/03/15 - v0.1.6 ]

### Changed

* Removed a load of branding
* Switched S3 uploader over to my own

## [ 2024/01/12 - v0.1.5 ]

### Fixed

* Outputted results file now has quotes around the completed result

## [ 2024/01/12 - v0.1.4 ]

### Fixed

* Updated test tracking to enable the test itself to create the tracker as some data was being missed.

## [ 2024/01/11 - v0.1.3 ]

### Changed/Added

* Added test tracking to write a results file out

## [ 2024/01/09 - v0.1.2 ]

### Changed/Added

* Updated process to fetch deployment before updating it to prevent an error "Operation cannot be fulfilled"
* minor pipeline update to add skip-existing to chart release

## [ 2024/01/03 - v0.1.1 ]

### Changed/Added

* Minor updates to test processing for better detection and to prevent failures of the whole app if one fails

## [ 2024/01/03 - v0.1.0 ]

### Changed/Added

* Rewritten (nearly) the entire codebase
* Added new helm chart to use as the source of the resources to deploy
* Added ability to build new chart on release - chart is manually version defined for now while using the
  chart-releaser-action
* Improved some checks around the tests that are run as part of the refactor

### Fixed

### Deprecated

* built in workloads in favor of helm chart

## [ 2023/11/27 - v0.1.0-beta.5 ]

### Changed/Added

* Updated container scanner
* Updated go modules
* Switched to chainguard Golang image
* Switched from alpine to WolfiOS.

### Fixed

### Deprecated

## [ 2023/08/30 - v0.1.0-beta.4 ]

### Changed/Added

* improved metric naming and start/stop data points

## [ 2023/08/30 - v0.1.0-beta.3 ]

### Changed/Added

* updated workflows
* Adding Metrics to Tests

## [ 2023/04/28 - v0.1.0-beta.2 ]

### Changed/Added

* dockerfile and pipeline additions
* removed golang precommit due to inactivity on repo - new one will be found and used soon
* changed logo to prevent getting a telling off ;-)

## [ 2023/01/20 - v0.1.0-beta.1 ]

### Changed/Added

* New process for working with DogKat which doesn't require an external helm chart

## [ Last Helm chart supported version ]

* Up to this point, there has been no changelog supplied for previous versions as it was a rapid iterative process.
