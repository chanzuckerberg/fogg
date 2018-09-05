# fogg FAQ

## Where did the name come from?

When looking for a name we searched for people and things related to the concept of terraforming planets, which led us to [Martyn Fogg](https://en.wikipedia.org/wiki/Martyn_J._Fogg). 'Fogg' seemed like a nice short name, and there didn't appear to be other software with this name. So here we are.

## What is sicc?

SICC was a predecessor to fogg that we used only internally at CZI.


## fogg.json

### What are modules
[Terraform modules](https://www.terraform.io/intro/getting-started/modules.html) are a collection of terraform resources (including module and data).

### What are accounts
Accounts correspond to AWS accounts. These can configure singleton resources such as IAM users and groups.

### What are components
Components are instantiations of terraform modules and resources. These can be things like RDS databases and their corresponding security groups, S3 website hosting, ECS clusters, etc.

### Managing Multiple AWS Accounts
A typical multi account scenario is having your staging and prod environments split into two different aws accounts.
A `fogg.json` that might address this use-case:
```json
{
  "defaults": {
    "aws_region_provider": "us-west-2",
    "aws_region_backend": "us-west-2",
    "aws_regions": [
      "us-east-2",
      "us-west-2"
    ],
    "aws_profile": "fogg-profile-staging",
    "aws_profile_backend": "fogg-profile-staging",
    "aws_profile_provider": "fogg-profile-staging",
    "infra_s3_bucket": "my-fogg-bucket",
    "project": "fogg-example",
    "terraform_version": "0.11.8",
    "owner": "fogg@example.com",
    "aws_provider_version": "1.33.0"
  },
  "modules": {
    "aurora": {},
    "redis": {},
    "some-shiny-new-tech": {}
  },
  "accounts": {
    "aws-staging-account": {
      "account_id": 000000000000,
      "regions": [
        "us-west-2",
        "us-east-1"
      ]
    },
    "aws-prod-account": {
      "account_id": 11111111111,
      "aws_profile_backend": "fogg-profile-prod",
      "aws_profile_provider": "fogg-profile-prod",
      "components": {
        "iam-users":{}
      }
    }
  },
  "envs": {
    "prod": {
      "account_id": 11111111111,
      "aws_profile_provider": "fogg-profile-prod",
      "components": {
        "redis": {},
        "security-alerts": {
          "account_id": 22222222222 ,
          "aws_profile_backend": "fogg-profile-prod",
          "aws_profile_provider": "fogg-profile-security"
        },
      }
    },
    "staging": {
      "components": {
        "redis": {},
        "security-alerts": {
          "account_id": 22222222222 ,
          "aws_profile_provider": "fogg-profile-security"
        },
      }
    },

    ....
```
In the above example I've configured a couple of interesting things. We have three accounts in play: Staging (account_id: 000000000000), Prod (account_id: 111111111111), Security (account_id: 222222222222). Just by looking at `fogg.json` I can see how both staging and prod have redis and security-alerts components. I can also see how these security alerts are centralized into a security account. Fogg's modularity and hierarchical configuration naturally allow us to do powerful things like centralized logging, centralized alerting, centralized identity management, separation of concerns along aws account boundary lines.
