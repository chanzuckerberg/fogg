# fogg

[![Build Status](https://travis-ci.com/chanzuckerberg/fogg.svg?token=JNM9vNLfRsFzCcF1uEgc&branch=master)](https://travis-ci.com/chanzuckerberg/fogg) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/chanzuckerberg/fogg)

Fogg is an opinionated tool for managing infrastructure-as-code repositories using Terraform.

Terraform is a powerful tool for managing infrastructure– great when things go right, but dangerous when they don't. Best practices are emerging for reducing this risk, but they require significant work and knowledge to apply consistently.

We built fogg to automate these practices and scale to a larger pool of engineers who don't have to be terraform experts to use it safely.

A few of the things fogg standardizes–

* repository layout
* remote state (locking coming soon)
* resource naming
* resource isolation

It makes life easy for folks working with cloud infrastructure. We've been using fogg and its predecessor internally at CZI for ~10 months. It has made it possible for many developers without terraform experience to safely roll new infrastructure with less stress and higher quality.

"I hope one day you might consider open sourcing `fogg`, i really love it. This would have saved me so much time in the past." - @lenn0x

## Getting Help

If you need help getting started with fogg, either open a github issue or join our [gitter chat room](https://gitter.im/chanzuckerberg/fogg).

## Install

## Mac

You can use homebrew to install fogg –

```
brew tap chanzuckerberg/homebrew-fogg
brew install fogg
```

## Linux, Windows, etc.

Binaries are available on the releases page. Download one for your architecture, put it in your path and make it executable.

## Usage

Fogg works entirely by generating code (terraform and make). It will generate directories and files to organize and standardize your repo and then it gets out of your way for you to use terraform and make to manage your infrastructure.

The basic workflow is –

1. update fogg.json
2. run `fogg apply` to code generate
3. use the generated Makefiles to run your Terraform commands

## Design Principles

### Convention over Configuration

Much like Ruby on Rails, we prefer to use conventions to organize our repos rather than a plethora of configuration. Our opinions might not be exactly the way you would do things, but our hope is that be having a set of clear opinions that are thoroughly applied will be productive.

### Transparency

Fogg tries to stay out of your way– it will do its work by generating Terraform and Make files, and then it step aside for you to manage your infrastructure. Everything that could effect your infrastructure is right there in your repository for you to read and understand.

There is no magic.

And if you ever decide to stop using it, you have a working repo you can take in a different direction, just stop running `fogg apply` and go your own way.

## Contributing

We use standard go tools + makefiles to build fogg. Getting started should be as simple as-

0. install go
1. go get github.com/chanzuckerberg/fogg
2. cd $GOPATH/src/github.com/chanzuckerberg/fogg
3. make

If you would like to contribute some code, fork this repo and send a pull request.

## Copyright

Copyright 2017-2018, Chan Zuckerberg Initiative, LLC

For license, see [LICENSE](LICENSE).
