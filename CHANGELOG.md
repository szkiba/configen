# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]


<a name="v0.2.0"></a>
## [v0.2.0] - 2021-04-23
### Build
- lower test coverage requirement temporarily

### Chore
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


[Unreleased]: https://github.com/szkiba/configen/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/szkiba/configen/compare/v0.1.0...v0.2.0
