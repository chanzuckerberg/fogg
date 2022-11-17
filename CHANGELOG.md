# Changelog

## [0.76.4](https://github.com/vincenthsh/fogg/compare/v0.76.3...v0.76.4) (2022-11-17)


### BugFixes

* change fogg ci gh actions to run on pr ([#28](https://github.com/vincenthsh/fogg/issues/28)) ([f0b6026](https://github.com/vincenthsh/fogg/commit/f0b602645e9dface75df7fa0320227c11d8292d3))

## [0.76.3](https://github.com/vincenthsh/fogg/compare/v0.76.2...v0.76.3) (2022-11-17)


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.136 to 1.44.139 ([#24](https://github.com/vincenthsh/fogg/issues/24)) ([a4e3987](https://github.com/vincenthsh/fogg/commit/a4e3987282baf4c76c1d3d902ee77df5e0f3f6be))


### BugFixes

* fogg ci gh actions tflint ([#25](https://github.com/vincenthsh/fogg/issues/25)) ([105b9c8](https://github.com/vincenthsh/fogg/commit/105b9c879ed3313221d422896cd01ea685331d86))
* fogg ci gh actions typo on tflint cache ([#27](https://github.com/vincenthsh/fogg/issues/27)) ([8ff5fef](https://github.com/vincenthsh/fogg/commit/8ff5feff18c20abe31404564006d450cbfb5f271))

## [0.76.2](https://github.com/vincenthsh/fogg/compare/v0.76.1...v0.76.2) (2022-11-17)


### BugFixes

* fogg ci gh actions tflint ([#23](https://github.com/vincenthsh/fogg/issues/23)) ([9db0eda](https://github.com/vincenthsh/fogg/commit/9db0eda6d5286e0e257137c83162f04c4187055d))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.122 to 1.44.126 ([#17](https://github.com/vincenthsh/fogg/issues/17)) ([3c7bdf8](https://github.com/vincenthsh/fogg/commit/3c7bdf8b306b20a593464feea822e73b5142fe38))
* bump github.com/aws/aws-sdk-go from 1.44.126 to 1.44.131 ([#19](https://github.com/vincenthsh/fogg/issues/19)) ([cd5edbe](https://github.com/vincenthsh/fogg/commit/cd5edbee4327544dfa1b180c477550b44e8db4c6))
* bump github.com/aws/aws-sdk-go from 1.44.131 to 1.44.136 ([#21](https://github.com/vincenthsh/fogg/issues/21)) ([b22e95c](https://github.com/vincenthsh/fogg/commit/b22e95cf5b537bcbf074c0c0bbd807431a342a2b))
* bump github.com/hashicorp/hcl/v2 from 2.14.1 to 2.15.0 ([#22](https://github.com/vincenthsh/fogg/issues/22)) ([36f8ed8](https://github.com/vincenthsh/fogg/commit/36f8ed8320652e9b18ef809e552763ee4ef3d578))

## [0.76.1](https://github.com/vincenthsh/fogg/compare/v0.76.0...v0.76.1) (2022-10-31)


### BugFixes

* Read terraformVersion from ComponentCommon.Common ([#15](https://github.com/vincenthsh/fogg/issues/15)) ([af6b9d4](https://github.com/vincenthsh/fogg/commit/af6b9d4d3dd2d74152d8ecc38f2e03dcf0536675))

## [0.76.0](https://github.com/vincenthsh/fogg/compare/v0.75.0...v0.76.0) (2022-10-31)


### Features

* Add pre-commit to Github Actions workflow ([#14](https://github.com/vincenthsh/fogg/issues/14)) ([dddef4c](https://github.com/vincenthsh/fogg/commit/dddef4c662f2c728048208d4444f7e087426ac18))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.120 to 1.44.121 ([#9](https://github.com/vincenthsh/fogg/issues/9)) ([32304bd](https://github.com/vincenthsh/fogg/commit/32304bde57d7ade05fe3f0f7e87d99f3ae879d51))
* bump github.com/aws/aws-sdk-go from 1.44.121 to 1.44.122 ([#13](https://github.com/vincenthsh/fogg/issues/13)) ([4a2d7dc](https://github.com/vincenthsh/fogg/commit/4a2d7dc0e856b1d4026bd89a317af878e45d3661))
* bump github.com/spf13/cobra from 1.6.0 to 1.6.1 ([#12](https://github.com/vincenthsh/fogg/issues/12)) ([db39a91](https://github.com/vincenthsh/fogg/commit/db39a91de81e01670d59458d0ad2c6f4b8d60899))
* bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#10](https://github.com/vincenthsh/fogg/issues/10)) ([3b213ab](https://github.com/vincenthsh/fogg/commit/3b213ab74240290cfe812d21a76666d637c45471))

## [0.75.0](https://github.com/vincenthsh/fogg/compare/v0.74.0...v0.75.0) (2022-10-21)


### ⚠ BREAKING CHANGES

* Use Snowflake-Lab version of Terraform Snowflake provider (#655)
* update tf fmt to 1.1.X; enforce lint and docs  (#590)

### Features

* add okta/okta provider ([#669](https://github.com/vincenthsh/fogg/issues/669)) ([d2a4dbf](https://github.com/vincenthsh/fogg/commit/d2a4dbff798c90952798170df85702590516d8a9))
* Add OpsGenie to fogg ([#712](https://github.com/vincenthsh/fogg/issues/712)) ([1b16ed9](https://github.com/vincenthsh/fogg/commit/1b16ed918090341817863dfbf56c21974793f801))
* Add Sops provider ([16acf92](https://github.com/vincenthsh/fogg/commit/16acf92ea040114dc0641aa513a1817e6b958019))
* Add support for atlantis repoCfg ([#8](https://github.com/vincenthsh/fogg/issues/8)) ([c26898f](https://github.com/vincenthsh/fogg/commit/c26898f2eae8527122df8014f546ae2d71477ef3))
* Add support for multiple modules per component ([afc2fa6](https://github.com/vincenthsh/fogg/commit/afc2fa6799dcade2fe9ce03d9e93a6f583c268ce))
* Add tf doc plugin to circle ci template ([#619](https://github.com/vincenthsh/fogg/issues/619)) ([9e5b8ba](https://github.com/vincenthsh/fogg/commit/9e5b8babb522ea4c632eab14bc875f12d863de17))
* adding assert library to fogg to make it easier to use in modules ([#642](https://github.com/vincenthsh/fogg/issues/642)) ([9451f02](https://github.com/vincenthsh/fogg/commit/9451f02fc7f4ae21a7f2191b94dcca2a7740fcda))
* allow fogg to authenticate with GH OAuth token ([#663](https://github.com/vincenthsh/fogg/issues/663)) ([a48a095](https://github.com/vincenthsh/fogg/commit/a48a095eb971df744f57bbf2fc80833cb6ddc59d))
* allow optional SSH keys in CI ([#726](https://github.com/vincenthsh/fogg/issues/726)) ([485303c](https://github.com/vincenthsh/fogg/commit/485303ccb0b88f92cca03f6db597f3db0901f583))
* allowing fogg to run with a github app token ([#679](https://github.com/vincenthsh/fogg/issues/679)) ([161a7c5](https://github.com/vincenthsh/fogg/commit/161a7c5d87bdd64086149c0b050de7c11b0a1ded))
* Automatic release process ([#621](https://github.com/vincenthsh/fogg/issues/621)) ([02584d2](https://github.com/vincenthsh/fogg/commit/02584d20ce1d6610cf9fd65acd7ada2ded825b53))
* CCIE-180 allow running 'fogg init' with flags instead of realtime user prompts ([#693](https://github.com/vincenthsh/fogg/issues/693)) ([07010d4](https://github.com/vincenthsh/fogg/commit/07010d43913b5194db33f6866920dd10f7d6b116))
* Default module variables ([#731](https://github.com/vincenthsh/fogg/issues/731)) ([cc29d58](https://github.com/vincenthsh/fogg/commit/cc29d58f83e4aa58ba37c62a2e559e93881cafc9))
* dont require bucket or profile ([#704](https://github.com/vincenthsh/fogg/issues/704)) ([b0614b7](https://github.com/vincenthsh/fogg/commit/b0614b7f84f412b5f67229c917071b4c15f3453d))
* fogg create tfe folder ([#730](https://github.com/vincenthsh/fogg/issues/730)) ([b259e9b](https://github.com/vincenthsh/fogg/commit/b259e9bbfa242e8ceec8538c7d6a1633d7b395df))
* make auth0 provider configurable ([#739](https://github.com/vincenthsh/fogg/issues/739)) ([4fd6aec](https://github.com/vincenthsh/fogg/commit/4fd6aec7507c64745047bed8f0b41474d8fa1749))
* **make:** Prefer local reviewdog binaries if they're present ([#580](https://github.com/vincenthsh/fogg/issues/580)) ([0939c08](https://github.com/vincenthsh/fogg/commit/0939c08375758c2e194f77487a818e95cc40b18a))
* no longer require root global backend ([#682](https://github.com/vincenthsh/fogg/issues/682)) ([ba46914](https://github.com/vincenthsh/fogg/commit/ba469143b996c7357a9326bcbe360c4843aac0f6))
* Optimize Github Actions CI ([#592](https://github.com/vincenthsh/fogg/issues/592)) ([34343a3](https://github.com/vincenthsh/fogg/commit/34343a3aac4c74eab3bec2c4174a8ce3c865f37f))
* Remove codecov from build and create conventional commits workflow ([#623](https://github.com/vincenthsh/fogg/issues/623)) ([7978a08](https://github.com/vincenthsh/fogg/commit/7978a08f19b7ff7469d62b6d87b4049ced348fac))
* update locals.tf.json on apply ([#707](https://github.com/vincenthsh/fogg/issues/707)) ([3a40011](https://github.com/vincenthsh/fogg/commit/3a40011a3f90f9706b0928fa5e0de2a197bb7938))
* updated ci that is faster and better ([#684](https://github.com/vincenthsh/fogg/issues/684)) ([3e78e17](https://github.com/vincenthsh/fogg/commit/3e78e17b4c0c2a8381012af7bb4f6875d6f7d65a))
* Use Snowflake-Lab version of Terraform Snowflake provider ([#655](https://github.com/vincenthsh/fogg/issues/655)) ([ad21af3](https://github.com/vincenthsh/fogg/commit/ad21af33b05780fee7d9f4b9755abfeea74aac5e))


### BugFixes

* assert provider getting into fogg ([#644](https://github.com/vincenthsh/fogg/issues/644)) ([02af482](https://github.com/vincenthsh/fogg/commit/02af4825c62b40ea030932ca8b7d196830770d58))
* Backend Workspaces is a block, not argument ([#591](https://github.com/vincenthsh/fogg/issues/591)) ([2065fe8](https://github.com/vincenthsh/fogg/commit/2065fe8f2861ec97685aff887430ed2998456515))
* bump failed release ([#705](https://github.com/vincenthsh/fogg/issues/705)) ([a8d0d6b](https://github.com/vincenthsh/fogg/commit/a8d0d6bbc73b1d93ce5f229eb2625a3184c610e2))
* cache breaking with github app creds ([#687](https://github.com/vincenthsh/fogg/issues/687)) ([5274906](https://github.com/vincenthsh/fogg/commit/5274906b6fe3b7eb9fe221575f2fe91a8b11cbbe))
* CircleCI config proper parallelism/bucketing ([#407](https://github.com/vincenthsh/fogg/issues/407)) ([0d244cc](https://github.com/vincenthsh/fogg/commit/0d244cc59cfb1fe8ecb049da0d0b32996117a715))
* don't list nested dirs in modules folders ci ([#719](https://github.com/vincenthsh/fogg/issues/719)) ([e3e6c42](https://github.com/vincenthsh/fogg/commit/e3e6c421e1a017eac4c3ea0ab7be23412ebce3cc))
* go 1.17 mods ([#589](https://github.com/vincenthsh/fogg/issues/589)) ([9225132](https://github.com/vincenthsh/fogg/commit/92251321a87c7d1b69a6ab9a3281877ce8f4d3d7))
* goreleaser: brews.github deprecated in favor of brews.tap ([#552](https://github.com/vincenthsh/fogg/issues/552)) ([d14797f](https://github.com/vincenthsh/fogg/commit/d14797f7b52faa97a067501edd10df0f991c6c2d))
* issue where the default version is not used ([#716](https://github.com/vincenthsh/fogg/issues/716)) ([2b45e36](https://github.com/vincenthsh/fogg/commit/2b45e36a8dad19b5e085332fdbecd1504b30ff51))
* missing quotes in allowed accounts ([#637](https://github.com/vincenthsh/fogg/issues/637)) ([2695b4e](https://github.com/vincenthsh/fogg/commit/2695b4eedd191ab4687a6d6fc985817949b40e69))
* Prevent including table if its an empty string ([#697](https://github.com/vincenthsh/fogg/issues/697)) ([14cdabf](https://github.com/vincenthsh/fogg/commit/14cdabf8871fe70ea2ff9238d0366b90e73b3ed2))
* release please pin ([#698](https://github.com/vincenthsh/fogg/issues/698)) ([28f9fad](https://github.com/vincenthsh/fogg/commit/28f9fad2c3edd3f8ded7794ffd35b0af5fc5d93a))
* remote apply default true ([#717](https://github.com/vincenthsh/fogg/issues/717)) ([282994b](https://github.com/vincenthsh/fogg/commit/282994b8222a6a2a72bfdd2e615dde26ab114811))
* remove test ([#699](https://github.com/vincenthsh/fogg/issues/699)) ([046a4b3](https://github.com/vincenthsh/fogg/commit/046a4b3c6cc964de4144347f3aa0ed04bf0cc351))
* treat AWS account IDs as strings ([#634](https://github.com/vincenthsh/fogg/issues/634)) ([d8f52bc](https://github.com/vincenthsh/fogg/commit/d8f52bc256b3eb2e8987c7bc2f83e7dfffd505e9))
* update locals struct to decode ([#727](https://github.com/vincenthsh/fogg/issues/727)) ([9d05d25](https://github.com/vincenthsh/fogg/commit/9d05d252a820850a9e790e1cc169c8ef7a877c0b))
* update tf fmt to 1.1.X; enforce lint and docs  ([#590](https://github.com/vincenthsh/fogg/issues/590)) ([ec73edb](https://github.com/vincenthsh/fogg/commit/ec73edb41a7eaeaabeef0ac28be6e88265b5d2f6))
* update the auth0 provider name ([#737](https://github.com/vincenthsh/fogg/issues/737)) ([cd92b36](https://github.com/vincenthsh/fogg/commit/cd92b36d70c5895e1c8b8a2a1b3863c84214540e))
* Upgrade terraform-config-inspect ([#729](https://github.com/vincenthsh/fogg/issues/729)) ([75a4032](https://github.com/vincenthsh/fogg/commit/75a4032858d69e4e3f1a2e82f9d8c60a07e2a2a7))
* Use github token ([d3ae53b](https://github.com/vincenthsh/fogg/commit/d3ae53ba7bd07a59069f146f436c1211608ae31b))
* version contraint in assert provider ([#646](https://github.com/vincenthsh/fogg/issues/646)) ([d27e54b](https://github.com/vincenthsh/fogg/commit/d27e54b3e766e4354c72f39f13c703e8b5ef1130))


### Misc

* Add darwin/arm64 to downloadable architectures ([#607](https://github.com/vincenthsh/fogg/issues/607)) ([11ef59b](https://github.com/vincenthsh/fogg/commit/11ef59b2f40f6831897674636a06dfe64f8cfa66))
* Add release please workflow ([364f856](https://github.com/vincenthsh/fogg/commit/364f8569f9b8bf7aeb3b4a05d9b48e676807e404))
* build CI use go-version-file ([b74a646](https://github.com/vincenthsh/fogg/commit/b74a6464d909a700505e1cd9f5387d19cdcad3eb))
* bump github.com/aws/aws-sdk-go from 1.44.100 to 1.44.104 ([#736](https://github.com/vincenthsh/fogg/issues/736)) ([73533f4](https://github.com/vincenthsh/fogg/commit/73533f4e1b2dd019b7f10bcf0d8534c74879966f))
* bump github.com/aws/aws-sdk-go from 1.44.104 to 1.44.105 ([#740](https://github.com/vincenthsh/fogg/issues/740)) ([ca757fc](https://github.com/vincenthsh/fogg/commit/ca757fc83fd0818c5d4ebfef796894b6a11b6da7))
* bump github.com/aws/aws-sdk-go from 1.44.105 to 1.44.120 ([#3](https://github.com/vincenthsh/fogg/issues/3)) ([77998dd](https://github.com/vincenthsh/fogg/commit/77998dd62c111f26431dac29d109b9538e808c80))
* bump github.com/aws/aws-sdk-go from 1.44.52 to 1.44.56 ([#681](https://github.com/vincenthsh/fogg/issues/681)) ([518f5b4](https://github.com/vincenthsh/fogg/commit/518f5b46e4fa152696fa0ad78e17f17868ce7fc2))
* bump github.com/aws/aws-sdk-go from 1.44.56 to 1.44.61 ([#689](https://github.com/vincenthsh/fogg/issues/689)) ([b35740f](https://github.com/vincenthsh/fogg/commit/b35740f7e8f00020e7cade6b4fe9fa244c97d2ed))
* bump github.com/aws/aws-sdk-go from 1.44.61 to 1.44.66 ([#692](https://github.com/vincenthsh/fogg/issues/692)) ([de5b495](https://github.com/vincenthsh/fogg/commit/de5b495259eaaecb31b6cd773d3ef51b9d4bb4cc))
* bump github.com/aws/aws-sdk-go from 1.44.66 to 1.44.67 ([#695](https://github.com/vincenthsh/fogg/issues/695)) ([c97e9c0](https://github.com/vincenthsh/fogg/commit/c97e9c08d48f51dd859f8b9d3447c4dc7e97c3e5))
* bump github.com/aws/aws-sdk-go from 1.44.67 to 1.44.70 ([#702](https://github.com/vincenthsh/fogg/issues/702)) ([4115db1](https://github.com/vincenthsh/fogg/commit/4115db135753d8f0368c605988de73f857863e4d))
* bump github.com/aws/aws-sdk-go from 1.44.70 to 1.44.76 ([#709](https://github.com/vincenthsh/fogg/issues/709)) ([c061961](https://github.com/vincenthsh/fogg/commit/c06196116b7949f06fabc40686ac064069a8e0a9))
* bump github.com/aws/aws-sdk-go from 1.44.76 to 1.44.81 ([#711](https://github.com/vincenthsh/fogg/issues/711)) ([ba6e36c](https://github.com/vincenthsh/fogg/commit/ba6e36c6d32ee522330c108a8754ee0ba7c4d962))
* bump github.com/aws/aws-sdk-go from 1.44.81 to 1.44.83 ([#714](https://github.com/vincenthsh/fogg/issues/714)) ([27124ee](https://github.com/vincenthsh/fogg/commit/27124ee1c77238bfc5b01bed31467acbea041584))
* bump github.com/aws/aws-sdk-go from 1.44.83 to 1.44.86 ([#715](https://github.com/vincenthsh/fogg/issues/715)) ([d365cd8](https://github.com/vincenthsh/fogg/commit/d365cd842760f0ec86fb7b0aaf6d748ae38bfad1))
* bump github.com/aws/aws-sdk-go from 1.44.86 to 1.44.91 ([#722](https://github.com/vincenthsh/fogg/issues/722)) ([d9eb39c](https://github.com/vincenthsh/fogg/commit/d9eb39c9861319ae61b031466ae44a177772cd1c))
* bump github.com/aws/aws-sdk-go from 1.44.91 to 1.44.95 ([#725](https://github.com/vincenthsh/fogg/issues/725)) ([8af988b](https://github.com/vincenthsh/fogg/commit/8af988b4d768952ad299f85245c8e177e09b7a6e))
* bump github.com/aws/aws-sdk-go from 1.44.95 to 1.44.100 ([#732](https://github.com/vincenthsh/fogg/issues/732)) ([18a1722](https://github.com/vincenthsh/fogg/commit/18a1722b37b57c96a315b4e466d32925254603bc))
* bump github.com/fatih/color from 1.10.0 to 1.13.0 ([#5](https://github.com/vincenthsh/fogg/issues/5)) ([07f85ef](https://github.com/vincenthsh/fogg/commit/07f85efe454833667d4c8171a47290250d6b96ad))
* bump github.com/go-errors/errors from 1.1.1 to 1.4.2 ([#4](https://github.com/vincenthsh/fogg/issues/4)) ([a811316](https://github.com/vincenthsh/fogg/commit/a8113167335deac0d9f8016ebac7b5d39450f421))
* bump github.com/hashicorp/hcl/v2 from 2.13.0 to 2.14.0 ([#723](https://github.com/vincenthsh/fogg/issues/723)) ([c175c08](https://github.com/vincenthsh/fogg/commit/c175c0831ad6fa84b327de2d5c82b283d7437f25))
* bump github.com/hashicorp/hcl/v2 from 2.14.0 to 2.14.1 ([#741](https://github.com/vincenthsh/fogg/issues/741)) ([b4b3269](https://github.com/vincenthsh/fogg/commit/b4b3269f213c5d9cd1b8c24b1ae293468a8933af))
* bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#690](https://github.com/vincenthsh/fogg/issues/690)) ([af48725](https://github.com/vincenthsh/fogg/commit/af487254263dd5dce0820b35f468448a98583eae))
* bump github.com/spf13/cobra from 1.5.0 to 1.6.0 ([#1](https://github.com/vincenthsh/fogg/issues/1)) ([63e6e0c](https://github.com/vincenthsh/fogg/commit/63e6e0c287c19cda82b13942f9c1c6e1b5318806))
* Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 ([#672](https://github.com/vincenthsh/fogg/issues/672)) ([ffbfc9c](https://github.com/vincenthsh/fogg/commit/ffbfc9c166d0d2edd7cb9932567ca7c7dc08e4dc))
* Bump go from 1.17 to 1.18 ([a7a8ae4](https://github.com/vincenthsh/fogg/commit/a7a8ae48e506615c46bcee62bf01c1ef3e242210))
* bump gopkg.in/ini.v1 from 1.66.6 to 1.67.0 ([#701](https://github.com/vincenthsh/fogg/issues/701)) ([f5ffb28](https://github.com/vincenthsh/fogg/commit/f5ffb28f0b426d8373fc2e63e473c71492624b21))
* Disable integration test ([bcb7532](https://github.com/vincenthsh/fogg/commit/bcb75323adb403e5ed88a54357bb59476143eb7b))
* Fix the dependabot PR commit prefix ([#677](https://github.com/vincenthsh/fogg/issues/677)) ([56930f1](https://github.com/vincenthsh/fogg/commit/56930f11b5817312e54dfff06f8db601d342175c))
* Have Git ignore .terraform.lock.hcl ([#657](https://github.com/vincenthsh/fogg/issues/657)) ([398f41e](https://github.com/vincenthsh/fogg/commit/398f41ef65f04c078dd0430b3ef00d7b090b3f4f))
* **main:** release 0.60.0 ([#622](https://github.com/vincenthsh/fogg/issues/622)) ([075e696](https://github.com/vincenthsh/fogg/commit/075e696a1702e1ec4eca7873889947a964f687ef))
* **main:** release 0.60.1 ([#635](https://github.com/vincenthsh/fogg/issues/635)) ([e0f3e35](https://github.com/vincenthsh/fogg/commit/e0f3e35ad9efa1f8a5701f7026fcf172f60791b2))
* **main:** release 0.60.2 ([#636](https://github.com/vincenthsh/fogg/issues/636)) ([155a74b](https://github.com/vincenthsh/fogg/commit/155a74b62a16299b1bce54609859c14785fd1976))
* **main:** release 0.60.3 ([#640](https://github.com/vincenthsh/fogg/issues/640)) ([2c8c1e1](https://github.com/vincenthsh/fogg/commit/2c8c1e1577c39fcd55af04a35ecf16054630d572))
* **main:** release 0.61.0 ([#643](https://github.com/vincenthsh/fogg/issues/643)) ([3900fec](https://github.com/vincenthsh/fogg/commit/3900fec02216636164f428d347fac04d798cf894))
* **main:** release 0.61.1 ([#645](https://github.com/vincenthsh/fogg/issues/645)) ([04a2836](https://github.com/vincenthsh/fogg/commit/04a283664a564dcfe012c851dbfed62e8148ceba))
* **main:** release 0.61.2 ([#647](https://github.com/vincenthsh/fogg/issues/647)) ([746a0a8](https://github.com/vincenthsh/fogg/commit/746a0a825dd93882f88bc2c1b519e2b0c062a2e8))
* **main:** release 0.62.0 ([#656](https://github.com/vincenthsh/fogg/issues/656)) ([bb0a966](https://github.com/vincenthsh/fogg/commit/bb0a9664b0476ab8301522b486b80fce8a7b57be))
* **main:** release 0.62.1 ([#658](https://github.com/vincenthsh/fogg/issues/658)) ([b79da43](https://github.com/vincenthsh/fogg/commit/b79da435012527407db2d657cb4f3d05364d5fbf))
* **main:** release 0.63.0 ([#666](https://github.com/vincenthsh/fogg/issues/666)) ([d4ed841](https://github.com/vincenthsh/fogg/commit/d4ed8415fce4b6bdf2b9fd8e6be2dd348c54380b))
* **main:** release 0.64.0 ([#670](https://github.com/vincenthsh/fogg/issues/670)) ([595dbe1](https://github.com/vincenthsh/fogg/commit/595dbe169f350fdca6af504541942e4f04fdbdec))
* **main:** release 0.64.1 ([#676](https://github.com/vincenthsh/fogg/issues/676)) ([46294ee](https://github.com/vincenthsh/fogg/commit/46294ee0880953c0597c25579a68e93ae4072174))
* **main:** release 0.65.0 ([#680](https://github.com/vincenthsh/fogg/issues/680)) ([464d5eb](https://github.com/vincenthsh/fogg/commit/464d5ebed8d018b29d5eade912da2c3c4723fd60))
* **main:** release 0.65.1 ([#683](https://github.com/vincenthsh/fogg/issues/683)) ([2b3afd8](https://github.com/vincenthsh/fogg/commit/2b3afd88b46513e61ba4b837de2225e5c0f4f597))
* **main:** release 0.66.0 ([#685](https://github.com/vincenthsh/fogg/issues/685)) ([c4062d1](https://github.com/vincenthsh/fogg/commit/c4062d117112ef61e0e60346b96a7a830b66c55c))
* **main:** release 0.66.1 ([#688](https://github.com/vincenthsh/fogg/issues/688)) ([2822b37](https://github.com/vincenthsh/fogg/commit/2822b37adc744ca8efe54a927f320ec63211306b))
* **main:** release 0.67.0 ([#691](https://github.com/vincenthsh/fogg/issues/691)) ([16f5b50](https://github.com/vincenthsh/fogg/commit/16f5b50a08e2234166e568ac0d9ae09fe29193b2))
* **main:** release 0.67.1 ([#694](https://github.com/vincenthsh/fogg/issues/694)) ([04508f2](https://github.com/vincenthsh/fogg/commit/04508f2eeec0ffdcdca92571e94bb3ad595ea2cc))
* **main:** release 0.67.2 ([#700](https://github.com/vincenthsh/fogg/issues/700)) ([ef2abec](https://github.com/vincenthsh/fogg/commit/ef2abecbf8cfb50f2903e1a937b44d552bb8af23))
* **main:** release 0.68.0 ([#703](https://github.com/vincenthsh/fogg/issues/703)) ([fbbcfc5](https://github.com/vincenthsh/fogg/commit/fbbcfc56802e4d327a0425119d7e8c4b47ada9ec))
* **main:** release 0.68.1 ([#706](https://github.com/vincenthsh/fogg/issues/706)) ([463cac5](https://github.com/vincenthsh/fogg/commit/463cac5c47ad407d605e40747657e2fc6aef54af))
* **main:** release 0.69.0 ([#708](https://github.com/vincenthsh/fogg/issues/708)) ([e741db7](https://github.com/vincenthsh/fogg/commit/e741db7ff701ac0d834f0f64b5c02860e4b850fd))
* **main:** release 0.70.0 ([#710](https://github.com/vincenthsh/fogg/issues/710)) ([a19ee00](https://github.com/vincenthsh/fogg/commit/a19ee009265e66f4efb1e8da831aa1faacfd73fb))
* **main:** release 0.70.1 ([#713](https://github.com/vincenthsh/fogg/issues/713)) ([c749891](https://github.com/vincenthsh/fogg/commit/c7498910c55d03493ac2994284a8ffecc6545320))
* **main:** release 0.70.2 ([#718](https://github.com/vincenthsh/fogg/issues/718)) ([3dfc621](https://github.com/vincenthsh/fogg/commit/3dfc621a96933956487171aa1f3bfcab10c36836))
* **main:** release 0.70.3 ([#720](https://github.com/vincenthsh/fogg/issues/720)) ([034d40e](https://github.com/vincenthsh/fogg/commit/034d40ebfab66215fc8134384d8e8d2e51e314a0))
* **main:** release 0.71.0 ([#724](https://github.com/vincenthsh/fogg/issues/724)) ([85be2ef](https://github.com/vincenthsh/fogg/commit/85be2ef7cb8cb9b8a9264d0ca28cebb61a62c85e))
* **main:** release 0.71.1 ([#728](https://github.com/vincenthsh/fogg/issues/728)) ([d581c46](https://github.com/vincenthsh/fogg/commit/d581c46a2dc08754d76abdd7b10b0318128d454b))
* **main:** release 0.72.0 ([#733](https://github.com/vincenthsh/fogg/issues/733)) ([e39ad20](https://github.com/vincenthsh/fogg/commit/e39ad20f238a9cae27cc9e0c9de1a43e5c6d9f05))
* **main:** release 0.73.0 ([#735](https://github.com/vincenthsh/fogg/issues/735)) ([62128a6](https://github.com/vincenthsh/fogg/commit/62128a62e49109c4d522b42f5ccebb4d6ff393dc))
* **main:** release 0.74.0 ([#738](https://github.com/vincenthsh/fogg/issues/738)) ([3015b1a](https://github.com/vincenthsh/fogg/commit/3015b1a36c0a044c1845f8b222dc14bda473cb29))
* mark fogg_ci.yml as auto generated ([#626](https://github.com/vincenthsh/fogg/issues/626)) ([313a89f](https://github.com/vincenthsh/fogg/commit/313a89fda7b3e3b044720529cacf0cc7682df905))
* Rebase ([afc2fa6](https://github.com/vincenthsh/fogg/commit/afc2fa6799dcade2fe9ce03d9e93a6f583c268ce))
* remove codeql; we don't use it ([#639](https://github.com/vincenthsh/fogg/issues/639)) ([5176a6d](https://github.com/vincenthsh/fogg/commit/5176a6dc1703836daef045305d7e0945f2940aa1))
* trigger release ([#675](https://github.com/vincenthsh/fogg/issues/675)) ([e869c6b](https://github.com/vincenthsh/fogg/commit/e869c6bc11e8f956e76ebb49a9d19dc019dcfd9a))
* update fork ([#7](https://github.com/vincenthsh/fogg/issues/7)) ([585234d](https://github.com/vincenthsh/fogg/commit/585234ddff8539b48564153bcd2c19bd343387c6))
* Update random provider to latest ([217ad49](https://github.com/vincenthsh/fogg/commit/217ad4986a8fd3310345cf88ebf9540516766f11))
* Update utility providers ([#743](https://github.com/vincenthsh/fogg/issues/743)) ([b4e6ac9](https://github.com/vincenthsh/fogg/commit/b4e6ac9ecd7c9affde72d525b36906273c55922b))

## [0.74.0](https://github.com/chanzuckerberg/fogg/compare/v0.73.0...v0.74.0) (2022-09-23)


### Features

* make auth0 provider configurable ([#739](https://github.com/chanzuckerberg/fogg/issues/739)) ([4fd6aec](https://github.com/chanzuckerberg/fogg/commit/4fd6aec7507c64745047bed8f0b41474d8fa1749))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.100 to 1.44.104 ([#736](https://github.com/chanzuckerberg/fogg/issues/736)) ([73533f4](https://github.com/chanzuckerberg/fogg/commit/73533f4e1b2dd019b7f10bcf0d8534c74879966f))


### BugFixes

* update the auth0 provider name ([#737](https://github.com/chanzuckerberg/fogg/issues/737)) ([cd92b36](https://github.com/chanzuckerberg/fogg/commit/cd92b36d70c5895e1c8b8a2a1b3863c84214540e))

## [0.73.0](https://github.com/chanzuckerberg/fogg/compare/v0.72.0...v0.73.0) (2022-09-23)


### Features

* Default module variables ([#731](https://github.com/chanzuckerberg/fogg/issues/731)) ([cc29d58](https://github.com/chanzuckerberg/fogg/commit/cc29d58f83e4aa58ba37c62a2e559e93881cafc9))


### BugFixes

* Upgrade terraform-config-inspect ([#729](https://github.com/chanzuckerberg/fogg/issues/729)) ([75a4032](https://github.com/chanzuckerberg/fogg/commit/75a4032858d69e4e3f1a2e82f9d8c60a07e2a2a7))

## [0.72.0](https://github.com/chanzuckerberg/fogg/compare/v0.71.1...v0.72.0) (2022-09-19)


### Features

* fogg create tfe folder ([#730](https://github.com/chanzuckerberg/fogg/issues/730)) ([b259e9b](https://github.com/chanzuckerberg/fogg/commit/b259e9bbfa242e8ceec8538c7d6a1633d7b395df))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.95 to 1.44.100 ([#732](https://github.com/chanzuckerberg/fogg/issues/732)) ([18a1722](https://github.com/chanzuckerberg/fogg/commit/18a1722b37b57c96a315b4e466d32925254603bc))

### [0.71.1](https://github.com/chanzuckerberg/fogg/compare/v0.71.0...v0.71.1) (2022-09-16)


### BugFixes

* update locals struct to decode ([#727](https://github.com/chanzuckerberg/fogg/issues/727)) ([9d05d25](https://github.com/chanzuckerberg/fogg/commit/9d05d252a820850a9e790e1cc169c8ef7a877c0b))

## [0.71.0](https://github.com/chanzuckerberg/fogg/compare/v0.70.3...v0.71.0) (2022-09-15)


### Features

* allow optional SSH keys in CI ([#726](https://github.com/chanzuckerberg/fogg/issues/726)) ([485303c](https://github.com/chanzuckerberg/fogg/commit/485303ccb0b88f92cca03f6db597f3db0901f583))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.86 to 1.44.91 ([#722](https://github.com/chanzuckerberg/fogg/issues/722)) ([d9eb39c](https://github.com/chanzuckerberg/fogg/commit/d9eb39c9861319ae61b031466ae44a177772cd1c))
* bump github.com/aws/aws-sdk-go from 1.44.91 to 1.44.95 ([#725](https://github.com/chanzuckerberg/fogg/issues/725)) ([8af988b](https://github.com/chanzuckerberg/fogg/commit/8af988b4d768952ad299f85245c8e177e09b7a6e))
* bump github.com/hashicorp/hcl/v2 from 2.13.0 to 2.14.0 ([#723](https://github.com/chanzuckerberg/fogg/issues/723)) ([c175c08](https://github.com/chanzuckerberg/fogg/commit/c175c0831ad6fa84b327de2d5c82b283d7437f25))

### [0.70.3](https://github.com/chanzuckerberg/fogg/compare/v0.70.2...v0.70.3) (2022-09-01)


### BugFixes

* don't list nested dirs in modules folders ci ([#719](https://github.com/chanzuckerberg/fogg/issues/719)) ([e3e6c42](https://github.com/chanzuckerberg/fogg/commit/e3e6c421e1a017eac4c3ea0ab7be23412ebce3cc))

### [0.70.2](https://github.com/chanzuckerberg/fogg/compare/v0.70.1...v0.70.2) (2022-08-31)


### BugFixes

* remote apply default true ([#717](https://github.com/chanzuckerberg/fogg/issues/717)) ([282994b](https://github.com/chanzuckerberg/fogg/commit/282994b8222a6a2a72bfdd2e615dde26ab114811))

### [0.70.1](https://github.com/chanzuckerberg/fogg/compare/v0.70.0...v0.70.1) (2022-08-29)


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.76 to 1.44.81 ([#711](https://github.com/chanzuckerberg/fogg/issues/711)) ([ba6e36c](https://github.com/chanzuckerberg/fogg/commit/ba6e36c6d32ee522330c108a8754ee0ba7c4d962))
* bump github.com/aws/aws-sdk-go from 1.44.81 to 1.44.83 ([#714](https://github.com/chanzuckerberg/fogg/issues/714)) ([27124ee](https://github.com/chanzuckerberg/fogg/commit/27124ee1c77238bfc5b01bed31467acbea041584))
* bump github.com/aws/aws-sdk-go from 1.44.83 to 1.44.86 ([#715](https://github.com/chanzuckerberg/fogg/issues/715)) ([d365cd8](https://github.com/chanzuckerberg/fogg/commit/d365cd842760f0ec86fb7b0aaf6d748ae38bfad1))


### BugFixes

* issue where the default version is not used ([#716](https://github.com/chanzuckerberg/fogg/issues/716)) ([2b45e36](https://github.com/chanzuckerberg/fogg/commit/2b45e36a8dad19b5e085332fdbecd1504b30ff51))

## [0.70.0](https://github.com/chanzuckerberg/fogg/compare/v0.69.0...v0.70.0) (2022-08-23)


### Features

* Add OpsGenie to fogg ([#712](https://github.com/chanzuckerberg/fogg/issues/712)) ([1b16ed9](https://github.com/chanzuckerberg/fogg/commit/1b16ed918090341817863dfbf56c21974793f801))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.70 to 1.44.76 ([#709](https://github.com/chanzuckerberg/fogg/issues/709)) ([c061961](https://github.com/chanzuckerberg/fogg/commit/c06196116b7949f06fabc40686ac064069a8e0a9))

## [0.69.0](https://github.com/chanzuckerberg/fogg/compare/v0.68.1...v0.69.0) (2022-08-12)


### Features

* update locals.tf.json on apply ([#707](https://github.com/chanzuckerberg/fogg/issues/707)) ([3a40011](https://github.com/chanzuckerberg/fogg/commit/3a40011a3f90f9706b0928fa5e0de2a197bb7938))

### [0.68.1](https://github.com/chanzuckerberg/fogg/compare/v0.68.0...v0.68.1) (2022-08-08)


### BugFixes

* bump failed release ([#705](https://github.com/chanzuckerberg/fogg/issues/705)) ([a8d0d6b](https://github.com/chanzuckerberg/fogg/commit/a8d0d6bbc73b1d93ce5f229eb2625a3184c610e2))

## [0.68.0](https://github.com/chanzuckerberg/fogg/compare/v0.67.2...v0.68.0) (2022-08-08)


### Features

* dont require bucket or profile ([#704](https://github.com/chanzuckerberg/fogg/issues/704)) ([b0614b7](https://github.com/chanzuckerberg/fogg/commit/b0614b7f84f412b5f67229c917071b4c15f3453d))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.67 to 1.44.70 ([#702](https://github.com/chanzuckerberg/fogg/issues/702)) ([4115db1](https://github.com/chanzuckerberg/fogg/commit/4115db135753d8f0368c605988de73f857863e4d))
* bump gopkg.in/ini.v1 from 1.66.6 to 1.67.0 ([#701](https://github.com/chanzuckerberg/fogg/issues/701)) ([f5ffb28](https://github.com/chanzuckerberg/fogg/commit/f5ffb28f0b426d8373fc2e63e473c71492624b21))

### [0.67.2](https://github.com/chanzuckerberg/fogg/compare/v0.67.1...v0.67.2) (2022-08-03)


### BugFixes

* remove test ([#699](https://github.com/chanzuckerberg/fogg/issues/699)) ([046a4b3](https://github.com/chanzuckerberg/fogg/commit/046a4b3c6cc964de4144347f3aa0ed04bf0cc351))

### [0.67.1](https://github.com/chanzuckerberg/fogg/compare/v0.67.0...v0.67.1) (2022-08-02)


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.61 to 1.44.66 ([#692](https://github.com/chanzuckerberg/fogg/issues/692)) ([de5b495](https://github.com/chanzuckerberg/fogg/commit/de5b495259eaaecb31b6cd773d3ef51b9d4bb4cc))
* bump github.com/aws/aws-sdk-go from 1.44.66 to 1.44.67 ([#695](https://github.com/chanzuckerberg/fogg/issues/695)) ([c97e9c0](https://github.com/chanzuckerberg/fogg/commit/c97e9c08d48f51dd859f8b9d3447c4dc7e97c3e5))


### BugFixes

* Prevent including table if its an empty string ([#697](https://github.com/chanzuckerberg/fogg/issues/697)) ([14cdabf](https://github.com/chanzuckerberg/fogg/commit/14cdabf8871fe70ea2ff9238d0366b90e73b3ed2))
* release please pin ([#698](https://github.com/chanzuckerberg/fogg/issues/698)) ([28f9fad](https://github.com/chanzuckerberg/fogg/commit/28f9fad2c3edd3f8ded7794ffd35b0af5fc5d93a))

## [0.67.0](https://github.com/chanzuckerberg/fogg/compare/v0.66.1...v0.67.0) (2022-08-01)


### Features

* CCIE-180 allow running 'fogg init' with flags instead of realtime user prompts ([#693](https://github.com/chanzuckerberg/fogg/issues/693)) ([07010d4](https://github.com/chanzuckerberg/fogg/commit/07010d43913b5194db33f6866920dd10f7d6b116))


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.56 to 1.44.61 ([#689](https://github.com/chanzuckerberg/fogg/issues/689)) ([b35740f](https://github.com/chanzuckerberg/fogg/commit/b35740f7e8f00020e7cade6b4fe9fa244c97d2ed))
* bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#690](https://github.com/chanzuckerberg/fogg/issues/690)) ([af48725](https://github.com/chanzuckerberg/fogg/commit/af487254263dd5dce0820b35f468448a98583eae))

## [0.66.1](https://github.com/chanzuckerberg/fogg/compare/v0.66.0...v0.66.1) (2022-07-19)


### BugFixes

* cache breaking with github app creds ([#687](https://github.com/chanzuckerberg/fogg/issues/687)) ([5274906](https://github.com/chanzuckerberg/fogg/commit/5274906b6fe3b7eb9fe221575f2fe91a8b11cbbe))

## [0.66.0](https://github.com/chanzuckerberg/fogg/compare/v0.65.1...v0.66.0) (2022-07-18)


### Features

* updated ci that is faster and better ([#684](https://github.com/chanzuckerberg/fogg/issues/684)) ([3e78e17](https://github.com/chanzuckerberg/fogg/commit/3e78e17b4c0c2a8381012af7bb4f6875d6f7d65a))

## [0.65.1](https://github.com/chanzuckerberg/fogg/compare/v0.65.0...v0.65.1) (2022-07-16)


### Misc

* bump github.com/aws/aws-sdk-go from 1.44.52 to 1.44.56 ([#681](https://github.com/chanzuckerberg/fogg/issues/681)) ([518f5b4](https://github.com/chanzuckerberg/fogg/commit/518f5b46e4fa152696fa0ad78e17f17868ce7fc2))
* Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 ([#672](https://github.com/chanzuckerberg/fogg/issues/672)) ([ffbfc9c](https://github.com/chanzuckerberg/fogg/commit/ffbfc9c166d0d2edd7cb9932567ca7c7dc08e4dc))

## [0.65.0](https://github.com/chanzuckerberg/fogg/compare/v0.64.1...v0.65.0) (2022-07-15)


### Features

* allowing fogg to run with a github app token ([#679](https://github.com/chanzuckerberg/fogg/issues/679)) ([161a7c5](https://github.com/chanzuckerberg/fogg/commit/161a7c5d87bdd64086149c0b050de7c11b0a1ded))
* no longer require root global backend ([#682](https://github.com/chanzuckerberg/fogg/issues/682)) ([ba46914](https://github.com/chanzuckerberg/fogg/commit/ba469143b996c7357a9326bcbe360c4843aac0f6))


### Misc

* Fix the dependabot PR commit prefix ([#677](https://github.com/chanzuckerberg/fogg/issues/677)) ([56930f1](https://github.com/chanzuckerberg/fogg/commit/56930f11b5817312e54dfff06f8db601d342175c))

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
