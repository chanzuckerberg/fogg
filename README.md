# fogg

Fogg is an opinionated tool for managing infrastructure-as-code repositories using Terraform.

It is the result of a desire to standardize and automate best practices. Some of the things it standardizes–

* repository layout
* remote state
* resource naming
* resource isolation

It makes life easy for folks working with cloud infrastructure. We've been using fogg and its predecessor internally at CZI for ~10 months. It has made it possible for many developers without terraform experience to roll new infrastructure with little stress and high quality.

"I hope one day you might consider open sourcing `fogg`, i really love it. This would of saved me so much time in the past." - @lenn0x

## Install

Binaries are available on the releases page. Download one for your architecture, put it in your path and make it executeable.

A homebrew tap will be coming soon.

## Usage

Fogg works entirely by code generation. It will generate diretories and files to organize and standardize your repo and then it gets out of your way for you to use terraform and make to manage your infrastructure.

The basic workflow is –

1. update fogg.json
2. run `fogg apply` to code generate
3. use the generated Makefiles to run your Terraform commands

## Design Principles

### Convention over Configuration

Much like Ruby on Rails, we prefer to use conventions to organize our repos rather than a plethora of configuration. Our opinions might not be exactly the way you would do things, but our hope is that be having a set of clear opinions that are thoroughly applied will be productive.

### Transparency

Fogg tries to stay out of your way– it will do the work it needs to by generating Terraform and make files, then it steps aside for you to manage your infrastructure. Everything that could effect your infra is right there on disk in your directory for you to read and understand.

There is no magic.

And if you ever decide to stop using it, you have a working repo you can take in a different direction, just stop running `fogg apply`.


## Copyright

Copyright 2017-2018, Chan Zuckerberg Initiative, LLC

For license, see [LICENSE].