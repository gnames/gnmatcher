# gnmatcher

The app matches a list of scientific name-strings to canonical forms of
scientific names from various biodiversity datasets.

## Introduction

This project is a component of a scientific names verification
(reconciliation/resolution) service [gnames]. The purpose of verification is to
compare a list of apparent scientific name-strings to a comprehensive set of
scientific names collected from many external biodiversity sources. The
`gnmatcher` project receives a list of name-strings and returns back 0 or more
canonical forms of known names for each name-string.

The project aims to do such verification as fast and accurate as possible.
Quite often, humans or character-recognition software (OCR) introduce
misspellings in the name-strings. For this reason, `gnmatcher` uses
fuzzy-matching algorithms when no exact match exists.  Also, for cases where
full name-string does not have a match, `gnmatcher` tries to match it against
parts of names. For example, if name-string did not get a match on a subspecies
level, the algorithm will try to match it on species and genus levels.

Reconciliation is the normalization of lexical variations of the same name, and
comparison of them to normalized names from biodiversity data sources.

Resolution is a determination of how a nomenclaturally registered name can be
interpreted from the point of taxonomy. For example, a name can be an accepted
name for species, a synonym, or a discarded one.

The `gnmatcher` app functions as an HTTP service. An app can access it using
HTTP client libraries.  The API's methods and structures are described in
the [model] dir.

## Input and Output

A user calls HTTP resource `/match` sending an array of name-strings to the
service and gets back canonical forms, the match type, as well as other
metadata described as an `Output` message in the [protobuf] file.

The optimal size of the input is 5-10 thousand name-strings per array. Note
that 10,000 is the maximal size, and larger arrays will be truncated.

## Performance

For performance measurement we took [100,000 name-strings][testdata] where only
30% of them were 'real' names. On a modern CPU with 12 hyper threads and
`GNM_JOBS_NUM` environment variable set to 8, the service was able to process
about 8,000 name-strings per second. For 'clean' data where most of the names
are "real", you should see an even higher performance.

## Prerequisites

* You will need PostgreSQL with a restored dump of
   [`gnames` database][gnames dump].

* Docker service

## Usage

### Usage with docker

* Install docker `gnmatcher` image: ``docker pull gnames/gnmatcher``.

* Copy [.env.example] file on user's disk and change values
  of environment variables accordingly.

* Start the service:

    ```bash
    docker run -p 8080:8080 -d --env-file your-env-file \
    gnames/gnmatcher -- rest -p 8080`
    ```

  This command will set the service on port 8080 and will make it available
  through port 8080 on a local machine.

### Usage from command line

* Download the [latest verion] of `gnmatcher` binary, untar and put somewhere
  in `PATH`.

* Run `gnmatcher -V` to generate configuration at
  `~/.config/gnmatcher.yaml`

* Edit `~/.config/gnmatcher.yaml` accordingly.

* Run ``gnmatcher rest -p 1234``

The service will run on the given port.

## Client

A user can find an example of a client for the service in this
[test file][rest-client]

## Development

To run tests a developer needs to install [BDD] binary [ginkgo]

There is a docker-compose file that sets up HTTP service to run tests. To run
it to the following:

1. Copy `.env.example` file to the `.env` file in the project's root directory,
   change the settings accordingly.

2. Build the `gnmatcher` binary and docker image using ``make dc`` command.

3. Run docker-compose command ``docker compose``

4. Run tests via ``go test ./...`` or ``ginkgo ./...``

[gnames]: https://github.com/gnames/gnames
[gnames dump]: https://opendata.globalnames.org/dumps/gnames-latest.sql.gz
[model]: https://github.com/gnames/gnmatcher/tree/master/model
[.env.example]: https://raw.githubusercontent.com/gnames/gnmatcher/master/.env.example
[testdata]: https://github.com/gnames/gnmatcher/blob/master/testdata/testdata.csv
[rest-client]: https://github.com/gnames/gnmatcher/blob/master/rest/rest_test.go
[BDD]: https://en.wikipedia.org/wiki/Behavior-driven_development
[ginkgo]: https://github.com/onsi/ginkgo
