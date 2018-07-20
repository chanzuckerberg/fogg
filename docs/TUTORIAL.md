# fogg tutorial

* **Author:** @ryanking
* **State**: Rough Draft - just trying to get all the content out

This tutorial will walk you through using fogg to create a new infrastructure repository.

Image that you are a company named Acme Corporation and want to deploy staging and production versions of your website where each one consists of a single server (let's keep it simple).

Note that fogg works by generating Terraform and Make files. It does not run any terraform commands for you.

Also note that we're not going to create the actual infrastructure here, just the scaffolding.

1. *install fogg*

    Go to https://github.com/chanzuckerberg/fogg/releases and download the latest version for your OS/arch. Put it in your $PATH and make it executeable.

1. *create a working directory*

    Create a new directory for this tutorial and `cd` into it.

1. *setup git*

    `git init`

   Terraform depends on working from the root of a git repository, though it doesn't need to be pushed anywhere, so `git init` is enough for now. After this, your directory should look like this–

   ```bash
   $ tree .
   .

   0 directories, 0 files
   ```

1. *initialize fogg*

   `fogg init`

   Fogg uses a `fogg.json` file to define the structure of your repository. This command exists to help bootstrap this configuration file by asking some simple questions.

   ```
   $ fogg init
    project name?: acme
    aws region?: us-west-2
    infra bucket name?: acme-infra
    auth profile?: acme-auth
    owner?: infra@acme.example
    ```

    A bit about those questions–

    * *project name*: we've got to name things around here. This is a high level name for your site, infrastructure or product.
    * *aws region* our setup is super flexible to run things in any and/or multiple regions. To get started we need a single region that we will configure as a default. This is also the region for the s3 bucket that will hold state files, so maybe think about it a little bit.
    * *infra bucket name* - we are going to store terraform's state files here. Currently fogg does not create this bucket, so you will need to do that ahead of time. Note that it should be in the same region you said above.
    * *auth profile* - we use aws authentication profiles, use this to specify the one to be used as a default
    * *owner* we make it easy to tag all your resources with their owner. If you put this here will will drop variables everywhere with the owner in it.

    And now your directory should look like this–

    ```
    $ tree
    .
    └── fogg.json

    0 directories, 1 file
    ```

    And if you look at `fogg.json`–

    ```json
    $ cat fogg.json
    {
    "defaults": {
        "aws_profile_backend": "acme-auth",
        "aws_profile_provider": "acme-auth",
        "aws_provider_version": "1.27.0",
        "aws_region_backend": "us-west-2",
        "aws_region_provider": "us-west-2",
        "infra_s3_bucket": "acme-infra",
        "owner": "infra@acme.example",
        "project": "acme",
        "shared_infra_version": "0.10.0",
        "terraform_version": "0.11.7"
    },
    "accounts": {},
    "envs": {},
    "modules": {}
    }
    ```

    Note that the questions you answered have all been filled into parts of this file and fogg supplied some additional defaults.

1. *build initial repo*

    As we said before fogg works by generating code (terraform, make and bash) and the general workflow is–

    1. update fogg.json
    2. run `fogg apply`

    Apply is the command that actually writes out all the changes we've specified in fogg.json.

    So now that we've written an initial `fogg.json` let's do an apply–

    ```
    $ fogg apply
    INFO templating .fogg-version
    INFO templating .gitattributes
    INFO templating .gitignore
    INFO templating Makefile
    INFO touching README.md
    INFO copying scripts/docker-ssh-forward.sh
    INFO copying scripts/docker-ssh-mount.sh
    INFO copying scripts/install_or_update.sh
    INFO templating scripts/ssh_config
    INFO templating Makefile
    INFO touching README.md
    INFO touching main.tf
    INFO touching outputs.tf
    INFO templating sicc.tf
    INFO touching variables.tf
    ```

    You'll see some output about that fogg is doing and now we have some structure to our repository–

    ```
    $ tree .
    .
    ├── Makefile
    ├── README.md
    ├── fogg.json
    ├── scripts
    │   ├── docker-ssh-forward.sh
    │   ├── docker-ssh-mount.sh
    │   ├── install_or_update.sh
    │   └── ssh_config
    └── terraform
        └── global
            ├── Makefile
            ├── README.md
            ├── main.tf
            ├── outputs.tf
            ├── sicc.tf
            └── variables.tf

    3 directories, 13 files
    ```

    Fogg has created 3 directories and put some files in them – `scripts` which exists to hold a handful of shell scripts useful to our operations and `terraform` which will include all our terraform code, both reusable modules and live infrastructure.

    Before go on – a bit about how fogg organizes repos –

    Fogg applies an opinionated way to organize terraform repos. This serves to both make it easy to factor your terraform code into many scopes / state files and also provide some consistency (make working on a team a bit easier).

    Fogg organizes tf code into `global`, `accounts`, `envs` and `components`.

    * `global` - things are trying global across all your infrastructure. A good example is a Route53 zone, to which you want to add recrords from everywhere in your infra.
    * `accounts` - things that are relavant at the account level (aws here) - most aws iam stuff goes here. Note that we make it easy to have multiple accounts which configs for each in `terraform/accounts/account-name`.
    * `envs` - think staging vs prod here. fogg makes it easy to keep your tf separate for each one
    * `components` - in addition to separating environments we do one step further and make it easy to have multiple state files for each environment. In fogg we call those components. Each env can have many components and they all get their own statefile. On top of that each gets a `terrafom_remote_state` data source for all the other components in the same env.

    With that in mind, let's create a new env.

1. *configure a staging env*

    Fogg helps you organize your terraform code and the resources they create into separate environments. Think 'staging' vs 'production'. It is advisable to have them separate so that you can operate on them independently. Let's create a 'staging' environment.

    To create a new env, edit your fogg.json to look like this–

    ```json
    {
        "defaults": {
            "aws_profile_backend": "acme-auth",
            "aws_profile_provider": "acme-auth",
            "aws_provider_version": "1.27.0",
            "aws_region_backend": "us-west-2",
            "aws_region_provider": "us-west-2",
            "infra_s3_bucket": "acme-infra",
            "owner": "infra@acme.example",
            "project": "acme",
            "shared_infra_version": "0.10.0",
            "terraform_version": "0.11.7"
        },
        "accounts": {},
        "envs": {
            "staging": {},
        }
        "modules": {}
    }
    ```

    Note that we added `"staging": {}` in the `envs` object.

    `fogg apply`

    ```
    $ fogg apply
    INFO templating .fogg-version
    INFO templating .gitattributes
    INFO templating .gitignore
    INFO templating Makefile
    INFO skipping touch on existing file README.md
    INFO copying scripts/docker-ssh-forward.sh
    INFO copying scripts/docker-ssh-mount.sh
    INFO copying scripts/install_or_update.sh
    INFO templating scripts/ssh_config
    INFO templating Makefile
    INFO touching README.md
    INFO templating Makefile
    INFO skipping touch on existing file README.md
    INFO templating fogg.tf
    INFO removing sicc.tf
    INFO error removing sicc.tf. ignoring
    INFO skipping touch on existing file main.tf
    INFO skipping touch on existing file outputs.tf
    INFO skipping touch on existing file variables.tf
    ```

    And your directory should look like:

    ```
    $ tree .
    .
    ├── Makefile
    ├── README.md
    ├── fogg.json
    ├── scripts
    │   ├── docker-ssh-forward.sh
    │   ├── docker-ssh-mount.sh
    │   ├── install_or_update.sh
    │   └── ssh_config
    └── terraform
        ├── envs
        │   └── staging
        │       ├── Makefile
        │       └── README.md
        └── global
            ├── Makefile
            ├── README.md
            ├── fogg.tf
            ├── main.tf
            ├── outputs.tf
            └── variables.tf

    5 directories, 15 files
    ```

    A terraform/staging directory has been creaded with a few files in it, but nowhere to put terraform files yet– those go in components which are nested in envs. So let's create a component–

1. create vpc?
1. create component
1. create prod
1. refactor to module
1. create module