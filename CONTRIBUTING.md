# Contributing

Thanks for your interest in contributing to Gonginx!

## Code of Conduct

Help us keep Gonginx open and inclusive. Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

* submit a ticket for your issue, assuming one does not already exist
  * clearly describe the issue including steps to reproduce when it is a bug
  * identify specific versions of the binaries and client libraries
* fork the repository on GitHub

## Making Changes

* create a branch from where you want to base your work
  * we typically name branches according to the following format: `fix/<issue_number>`
* make commits of logical units
* make sure your commit messages are in a clear and readable format, example:
  
```
fixed bug in http context
  
* fixed parsing locations
* cleanup variable replacing
* ...
```

* if you're fixing a bug or adding functionality you have to write a test that fails without your fix
* make sure to run `make test` in the root of the repo to ensure that your code is
  properly formatted and that tests pass.
    * test must prove the feature presents and fail when you revert your changes
* for parser/dumper safety changes, also run:
  * `go test -race ./...`
  * focused parser fuzz smoke run (optional but recommended): `go test -run=^$ -fuzz=FuzzParserParseNoPanic -fuzztime=5s ./parser`

## Submitting Changes

* push your changes to your branch in your fork of the repository
* submit a pull request against gongix' repository
