# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]


<a name="v0.3.2"></a>
## [v0.3.2] - 2021-04-26
### Fix
- experimental markdown processing flags


<a name="v0.3.1"></a>
## [v0.3.1] - 2021-04-26
### Chore
- prepare v0.3.1

### Feat
- experimental (hidden flag): render and open README.md

### Fix
- watch raw input directory changes too


<a name="v0.3.0"></a>
## [v0.3.0] - 2021-04-25
### Chore
- prepare v0.3.0

### Feat
- add support for raw file copy
- added "expr" template function (https://github.com/antonmedv/expr)


<a name="v0.2.0"></a>
## [v0.2.0] - 2021-04-23
### Build
- lower test coverage requirement temporarily

### Chore
- prepare v0.2.0
- temporary downgrade test coverage requirements
- disable switch lint

### Docs
- update usage

### Feat
- improve watching
- watch and generate onf input change
- validate template function (JSON Schema based)
- better JSON schema preloading
- new template functions: jp (jmespath), jq, jptr (JSON Pointer), uritpl (URI template), qsParse, qsJoin
- support jsonc, JSON with comments
- added console output handler template functions
- added assert and equal template functions

### Fix
- imporve YAML unmarshal to produce JSON like maps

### Refactor
- rename internal console struct to deferrer


<a name="v0.1.0"></a>
## v0.1.0 - 2021-04-19
### Chore
- prepare v0.1.0

### Ci
- fix workflow configuration

### Feat
- initial version


[Unreleased]: https://github.com/szkiba/configen/compare/v0.3.2...HEAD
[v0.3.2]: https://github.com/szkiba/configen/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/szkiba/configen/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/szkiba/configen/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/szkiba/configen/compare/v0.1.0...v0.2.0
