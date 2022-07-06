# Changelog

## [0.64.1](https://github.com/chanzuckerberg/fogg/compare/v0.64.0...v0.64.1) (2022-07-06)


### Misc

* trigger release ([#675](https://github.com/chanzuckerberg/fogg/issues/675)) ([e869c6b](https://github.com/chanzuckerberg/fogg/commit/e869c6bc11e8f956e76ebb49a9d19dc019dcfd9a))

## [0.64.0](https://github.com/chanzuckerberg/fogg/compare/v0.63.0...v0.64.0) (2022-06-29)


### Features

* add okta/okta provider ([#669](https://github.com/chanzuckerberg/fogg/issues/669)) ([d2a4dbf](https://github.com/chanzuckerberg/fogg/commit/d2a4dbff798c90952798170df85702590516d8a9))

## [0.63.0](https://github.com/chanzuckerberg/fogg/compare/v0.62.1...v0.63.0) (2022-06-13)


### Features

* allow fogg to authenticate with GH OAuth token ([#663](https://github.com/chanzuckerberg/fogg/issues/663)) ([a48a095](https://github.com/chanzuckerberg/fogg/commit/a48a095eb971df744f57bbf2fc80833cb6ddc59d))

## [0.62.1](https://github.com/chanzuckerberg/fogg/compare/v0.62.0...v0.62.1) (2022-06-02)


### Misc

* Have Git ignore .terraform.lock.hcl ([#657](https://github.com/chanzuckerberg/fogg/issues/657)) ([398f41e](https://github.com/chanzuckerberg/fogg/commit/398f41ef65f04c078dd0430b3ef00d7b090b3f4f))

## [0.62.0](https://github.com/chanzuckerberg/fogg/compare/v0.61.2...v0.62.0) (2022-06-02)


### ⚠ BREAKING CHANGES

* Use Snowflake-Lab version of Terraform Snowflake provider (#655)

### Features

* Use Snowflake-Lab version of Terraform Snowflake provider ([#655](https://github.com/chanzuckerberg/fogg/issues/655)) ([ad21af3](https://github.com/chanzuckerberg/fogg/commit/ad21af33b05780fee7d9f4b9755abfeea74aac5e))

### [0.61.2](https://github.com/chanzuckerberg/fogg/compare/v0.61.1...v0.61.2) (2022-05-05)


### BugFixes

* version contraint in assert provider ([#646](https://github.com/chanzuckerberg/fogg/issues/646)) ([d27e54b](https://github.com/chanzuckerberg/fogg/commit/d27e54b3e766e4354c72f39f13c703e8b5ef1130))

### [0.61.1](https://github.com/chanzuckerberg/fogg/compare/v0.61.0...v0.61.1) (2022-05-05)


### BugFixes

* assert provider getting into fogg ([#644](https://github.com/chanzuckerberg/fogg/issues/644)) ([02af482](https://github.com/chanzuckerberg/fogg/commit/02af4825c62b40ea030932ca8b7d196830770d58))

## [0.61.0](https://github.com/chanzuckerberg/fogg/compare/v0.60.3...v0.61.0) (2022-05-04)


### Features

* adding assert library to fogg to make it easier to use in modules ([#642](https://github.com/chanzuckerberg/fogg/issues/642)) ([9451f02](https://github.com/chanzuckerberg/fogg/commit/9451f02fc7f4ae21a7f2191b94dcca2a7740fcda))

### [0.60.3](https://github.com/chanzuckerberg/fogg/compare/v0.60.2...v0.60.3) (2022-04-26)


### Misc

* remove codeql; we don't use it ([#639](https://github.com/chanzuckerberg/fogg/issues/639)) ([5176a6d](https://github.com/chanzuckerberg/fogg/commit/5176a6dc1703836daef045305d7e0945f2940aa1))


### BugFixes

* missing quotes in allowed accounts ([#637](https://github.com/chanzuckerberg/fogg/issues/637)) ([2695b4e](https://github.com/chanzuckerberg/fogg/commit/2695b4eedd191ab4687a6d6fc985817949b40e69))

### [0.60.2](https://github.com/chanzuckerberg/fogg/compare/v0.60.1...v0.60.2) (2022-04-26)


### Misc

* mark fogg_ci.yml as auto generated ([#626](https://github.com/chanzuckerberg/fogg/issues/626)) ([313a89f](https://github.com/chanzuckerberg/fogg/commit/313a89fda7b3e3b044720529cacf0cc7682df905))

### [0.60.1](https://github.com/chanzuckerberg/fogg/compare/v0.60.0...v0.60.1) (2022-04-26)


### BugFixes

* treat AWS account IDs as strings ([#634](https://github.com/chanzuckerberg/fogg/issues/634)) ([d8f52bc](https://github.com/chanzuckerberg/fogg/commit/d8f52bc256b3eb2e8987c7bc2f83e7dfffd505e9))

## [0.60.0](https://github.com/chanzuckerberg/fogg/compare/v0.59.2...v0.60.0) (2022-03-21)


### Features

* Add tf doc plugin to circle ci template ([#619](https://github.com/chanzuckerberg/fogg/issues/619)) ([9e5b8ba](https://github.com/chanzuckerberg/fogg/commit/9e5b8babb522ea4c632eab14bc875f12d863de17))
* Automatic release process ([#621](https://github.com/chanzuckerberg/fogg/issues/621)) ([02584d2](https://github.com/chanzuckerberg/fogg/commit/02584d20ce1d6610cf9fd65acd7ada2ded825b53))
* Remove codecov from build and create conventional commits workflow ([#623](https://github.com/chanzuckerberg/fogg/issues/623)) ([7978a08](https://github.com/chanzuckerberg/fogg/commit/7978a08f19b7ff7469d62b6d87b4049ced348fac))


### Misc

* Add darwin/arm64 to downloadable architectures ([#607](https://github.com/chanzuckerberg/fogg/issues/607)) ([11ef59b](https://github.com/chanzuckerberg/fogg/commit/11ef59b2f40f6831897674636a06dfe64f8cfa66))
