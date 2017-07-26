# istanbul-coverage-warden
Watch for source code coverage growth
# Why?
We are lazy, getting older, forget thing and sometimes hate unit tests :-) This tool aimed to help us not to forget improve code coverage for a new functionality we add.

# How to use
Run istanbul with cobertura reporter `istanbul report cobertura`

Then run istanbul-coverage-warden and supply current cobertura report and previous report from previous build:
```istanbul-coverage-warden --current ./workspace/cobertura-coverage.xml --previous ./workspace-prev/cobertura-coverage-prev.xml```

Then check for return codes:
* 0 - Success
* 1 - Coverage dropped, build should fail
* 2 - File open error (might be first build)

# Installation

If you have ready binarry, just copy it to Jenkins. It should not require any dependencies. If you need to build it, follow the following steps:

### Install and configure Go
https://golang.org/doc/install

### Build package
go get -u https://github.com/Fonteva/istanbul-coverage-warden

It should download and build binary
