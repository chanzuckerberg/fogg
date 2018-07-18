# fogg tutorial

This tutorial will walk you through using fogg to create a new infrastructure repository.

Image your are a company named Acme Corporation and want to deploy staging and production versions of your website where each one consists of a single server (let's keep it simple). 

Note that fogg works by generating Terraform and Make files. It does not run any terraform commands for you.

Also note that we're not going to create the actual infrastructure here, just the scaffolding.

0. *install fogg*

    Go to https://github.com/chanzuckerberg/fogg/releases and down load the latest version for your OS/arch. Put it in your $PATH and make it executeable.

0. *create a working directory*

    Create a new directory for this tutorial and `cd` to it.

0. *setup git*

    `git init`

   Terraform depends on working from the root of a git repository, though it doesn't need to be pushed anywhere, so `git init` is enough for now. After this, your directory should look like this–

   ```
   $ tree .
   .

   0 directories, 0 files
   ```

0. *initialize fogg*
   
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
    And now your directory should look like this–

    ```
    $ tree
    .
    └── fogg.json

    0 directories, 1 file
    ```

    And if you look at `fogg.json`–

    ```
    $ cat fogg.json 
    {
    "defaults": {
        "aws_profile_backend": "acme-auth",
        "aws_profile_provider": "acme-auth",
        "aws_provider_version": "1.27.0",
        "aws_region_backend": "us-west-2",
        "aws_region_provider": "us-west-2",
        "aws_regions": [
            "ap-south-1",
            "eu-west-3",
            "eu-west-2",
            "eu-west-1",
            "ap-northeast-2",
            "ap-northeast-1",
            "sa-east-1",
            "ca-central-1",
            "ap-southeast-1",
            "ap-southeast-2",
            "eu-central-1",
            "us-east-1",
            "us-east-2",
            "us-west-1",
            "us-west-2"
        ],
        "infra_s3_bucket": "acme-infra",
        "owner": "infra@acme.example",
        "project": "acme",
        "shared_infra_version": "0.10.0",
        "terraform_version": "0.11.0"
    },
    "accounts": {},
    "envs": {},
    "modules": {}
    }
    ```

    Note that the questions you answered have all been filled into parts of this file and fogg supplied some additional defaults.

0. *build initial repo*

    As we said before fogg works by generating code (terraform, make and bash) and the general workflow is–

    0. update fogg.json
    0. run `fogg apply`

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

    Fogg has created 2 directories and put some files in them – `scripts` which exists to hold a handful of shell scripts useful to our operations and `terraform` which will include all our terraform code, both reusable modules and live infrastructure.

0. *configure a staging env*

    Fogg helps you organize your terraform code and the resources they create into separate environments. Thing 'staging' vs 'production' vs 'qa'. It is advisable to have they separate so that you can operate on them independently. Let's create a 'staging' environment.

    To create a new env, edit your fogg.json to look like this–

    ```
    {
        "defaults": {
            "aws_profile_backend": "acme-auth",
            "aws_profile_provider": "acme-auth",
            "aws_provider_version": "1.27.0",
            "aws_region_backend": "us-west-2",
            "aws_region_provider": "us-west-2",
            "aws_regions": [
            "ap-south-1",
            "eu-west-3",
            "eu-west-2",
            "eu-west-1",
            "ap-northeast-2",
            "ap-northeast-1",
            "sa-east-1",
            "ca-central-1",
            "ap-southeast-1",
            "ap-southeast-2",
            "eu-central-1",
            "us-east-1",
            "us-east-2",
            "us-west-1",
            "us-west-2"
            ],
            "infra_s3_bucket": "acme-infra",
            "owner": "infra@acme.example",
            "project": "acme",
            "shared_infra_version": "0.10.0",
            "terraform_version": "0.11.0"
        },
        "accounts": {},
        "envs": {
            "staging": {
                "type": "bare"
            }
        },
        "modules": {}
    }
    ```

    Note that we added `"staging": {}` in the `envs` object. [Future versions of fogg will add a command to do this.]

0. *create staging env*


* create component
* create module
* ?