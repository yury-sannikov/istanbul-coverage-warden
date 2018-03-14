# istanbul-coverage-warden
Watch for source code coverage growth

commands:

copy [branch-name] [S3 bucket]

Copy coverage files to S3 bucket for later chech.

check [destination-branch] [S3 bucket] [coverage.xml]

check against copied coverage, fail build if coverage drop


# Why?
We are lazy, getting older, forget thing and sometimes hate unit tests :-) This tool aimed to help us not to forget improve code coverage for a new functionality we add.

# How to use
`istanbul-coverage-warden` is supposed to be executed during pull request build and during main branch build. Main brach build should use `copy` command to expose gathered code coverage staticstics to S3 bucket for later use.

Pull request builder shoud use `check` command, and specify branch it build against. `check` command then download main branch code coverage metrics and check against current code coverage. If any source code file has lower `LineRate` than main branch, it exit with non-zero code and fail pull request build.

### Copy metrics from main branch
```
istanbul-coverage-warden copy \
  -d my-project-codecoverage-check \
  -b master \
  -f ./cobertura-coverage.xml
```

### Check metrics during pull request build
```
istanbul-coverage-warden check \
  -d my-project-codecoverage-check \
  -b master \
  -f ./cobertura-coverage.xml
```

Pull request build job should provide branch name your PR made agains. For Jenkins Github Pull Request plugin you might use `ghprbTargetBranch` parameter.

Then check for return codes:
* 0 - Success
* 1 - Coverage dropped, build should fail
* 2 - File open error (might be first build)

# Installation

If you have ready binarry, just copy it to Jenkins. It should not require any dependencies. If you need to build it, follow the following steps:

### Install and configure Go
https://golang.org/doc/install

### Build package
go get -u github.com/yury-sannikov/istanbul-coverage-warden

It should download and build binary
