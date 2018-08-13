# fogg design

Fogg is designed around a few principles– convention over configuration and transparency.

## Principles

### Convention over Configuration

Much time in software engineering is spend arguing over style and other issues with roughly equivalent solutions.

With fogg we provide a set of answers to questions like "how should I organize my terraform code"? Much like Ruby on Rails, we've identified some patterns and made it easy if you follow them.

Some things we standardize–

* repository layout
* remote state (locking coming soon)
* resource naming
* resource isolation

### Transparency

A significant challenge in practicing infrastructure-as-code is getting feedback on what your code will do.

Tools like Terraform make this better by having an explicit plan stage which will tell you what's about to happen.

As we've built tools around Terraform our goal has been to avoid making this worse – it should not be harder to understand what's going on.

To that end, one of the significant decisions we made was to have this tool work via code generation. That means that you can always read the code we've generated to understand what's going on (you can even tweak it if you need to temporarily work around a limitation).

This tranparency should make it easier to try out fogg– you can always see the code for yourself and if you even decide to stop using it, you already have a working repository.

## Terraform Best Practices

There are a number of terraform best practices which reduce operational risk and improve collaboration.

### Remote State

Terraform relies on a state file, which maps resources it manages to their state and metadata in cloud providers. Think of it like a join table.

You can either store this state file locally (wherever you are running Terraform) or remotely (on something like S3).

Storing it remotely prevents data loss, and improves collaboration. Fogg will automatically configure remote state in S3 for all terraform scopes it sets up.

### limiting blast radius

Blast radius is a metaphorical term for software operations wherein we consider "if this thing blows up how much else will go with it".

That is – if this tools fails here how far will the damage go.

This is relevant in Terrafrom because of its powerful feature of being able to update many resources at once.

To limit potential damage it is best to avoid putting all your terraform managed resources in one scope (ie one state file). Others have argued that you [need 1 per environment](https://charity.wtf/2016/03/30/terraform-vpc-and-why-you-want-a-tfstate-file-per-env/). Fogg goes further and allows you to have multiple scopes (we call them components) in an environment which are all loosely tied together by [remote state configurations](https://www.terraform.io/docs/providers/terraform/d/remote_state.html).

Now, setting up multiple scopes for terraform is as easy as creating directories. But what fogg provides is a way to maintain high engineer productivity when also factoring code into multiple components.

### tagging / naming conventions

It is best to tag and name resources consistently. Fogg helps with this be generating variables with names and tags for you to use.

