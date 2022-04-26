module github.com/chanzuckerberg/fogg

go 1.17

replace github.com/spf13/afero v1.2.2 => github.com/chanzuckerberg/afero v0.0.0-20190514223411-36a9495a9b51

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/antzucaro/matchr v0.0.0-20191224151129-ab6ba461ddec
	github.com/aws/aws-sdk-go v1.43.45
	github.com/blang/semver v3.5.1+incompatible
	github.com/chanzuckerberg/go-misc v0.0.0-20210209191033-ae96c0409e3f
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.10.0
	github.com/go-errors/errors v1.1.1
	github.com/google/go-github/v27 v27.0.6
	github.com/hashicorp/go-getter v1.5.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/hcl/v2 v2.12.0
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80
	github.com/hashicorp/terraform v0.14.9
	github.com/hashicorp/terraform-config-inspect v0.0.0-20200806211835-c481b8bfa41e
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/kr/pretty v0.2.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/segmentio/go-prompt v1.2.1-0.20161017233205-f0d19b6901ad
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.1
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/ini.v1 v1.66.4
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

require (
	cloud.google.com/go v0.65.0 // indirect
	cloud.google.com/go/storage v1.10.0 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-textseg v1.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/apparentlymart/go-versions v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bmatcuk/doublestar v1.2.4 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/klauspost/compress v1.11.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/zclconf/go-cty v1.8.0 // indirect
	github.com/zclconf/go-cty-yaml v1.0.2 // indirect
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.34.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
	google.golang.org/grpc v1.31.1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
