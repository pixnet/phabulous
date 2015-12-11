![Phabulous](http://i.imgur.com/0ezr6XZ.png)

![Example](http://i.imgur.com/Uv4nVJa.png)

Phabulous is a Slack bot for Phabricator.

> **Node.js version:** The Javascript version of this project has been replaced
with a rewrite from scratch in Go. The code for the Javascript version is
available at the **legacy** branch, but it won't be actively maintained.

[![wercker status](https://app.wercker.com/status/f135144cca86098cea87c00a36dfc564/s/master "wercker status")](https://app.wercker.com/project/bykey/f135144cca86098cea87c00a36dfc564)

## Features

- Post URLs to revisions and tasks right from Slack.
- Summon (mention) reviewers of a revision.
- Route specific events (Tasks, Revisions, Commits) into specific channels.
- Push all feed events into a single channel (This may flood a channel if your
  organization is big enough).
- Pretty icons ;)

## Guides

- [Getting Started](http://phabricator.chromabits.com/w/phabulous/start/):
A guide on how to setup Phabulous for the first time.
- [Upgrade Notes](http://phabricator.chromabits.com/w/phabulous/upgrade/):
Instructions on how to upgrade to newer versions of Phabulous.
- [Command Reference](http://phabricator.chromabits.com/w/phabulous/commands/):
Reference for all available bot commands.
- [Help & Troubleshooting](http://phabricator.chromabits.com/w/phabulous/faq/):
Tips and answers to common problems.
- [Wiki](http://phabricator.chromabits.com/w/phabulous/): More articles and
information about Phabulous.

## Compiling from source

To compile Phabulous, you need a recent version of Go:

```
go get github.com/etcinit/phabulous

// or, for cross-compiling:

go get github.com/mitchellh/gox
git clone git@github.com:etcinit/phabulous.git
cd phabulous
make
```
## Custom make for linux machines
### Fetch go dependencies
```
go get github.com/mitchellh/gox
go get github.com/kr/pretty
go get github.com/jacobstr/confer
go get github.com/Sirupsen/logrus
go get github.com/codegangsta/cli
go get github.com/etcinit/phabulous/app
go get github.com/facebookgo/inject
```
### Build
git clone git@github.com:etcinit/phabulous.git
cd phabulous && make
