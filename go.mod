module github.com/chanzuckerberg/fogg

go 1.18

replace github.com/spf13/afero v1.8.2 => github.com/chanzuckerberg/afero v0.0.0-20190514223411-36a9495a9b51

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/antzucaro/matchr v0.0.0-20191224151129-ab6ba461ddec
	github.com/aws/aws-sdk-go v1.44.207
	github.com/blang/semver v3.5.1+incompatible
	github.com/chanzuckerberg/go-misc v0.0.0-20210209191033-ae96c0409e3f
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/fatih/color v1.13.0
	github.com/go-errors/errors v1.4.2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/google/go-github/v27 v27.0.6
	github.com/hashicorp/go-getter v1.6.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/hcl/v2 v2.16.1
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80
	github.com/hashicorp/terraform v0.14.9
	github.com/hashicorp/terraform-config-inspect v0.0.0-20211115214459-90acf1ca460f
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/runatlantis/atlantis v0.20.1
	github.com/segmentio/go-prompt v1.2.1-0.20161017233205-f0d19b6901ad
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/afero v1.8.2
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	golang.org/x/exp v0.0.0-20220916125017-b168a2c6b86b
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/ini.v1 v1.67.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v1.6.1 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/storage v1.14.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-textseg v1.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/apparentlymart/go-versions v1.0.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bmatcuk/doublestar v1.2.4 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-ozzo/ozzo-validation v0.0.0-20170913164239-85dcd8368eba // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/compress v1.11.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	github.com/zclconf/go-cty-yaml v1.0.2 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220517005047-85d78b3ac167 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	google.golang.org/api v0.81.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/grpc v1.46.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)
