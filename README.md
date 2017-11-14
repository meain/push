# Push [![Build Status](https://travis-ci.org/meain/push.svg?branch=master)](https://travis-ci.org/meain/push)

> [pushbullet](https://www.pushbullet.com/) + [transfer.sh](https://transfer.sh/)

Just push a simple note to [pushbullet](https://www.pushbullet.com/) or send a file link to pushbullet after uploading it to [transfer.sh](https://transfer.sh/)

![screencast](https://i.imgur.com/EnCKJTE.gif)

## Installation

Download binary from [releases page](https://github.com/meain/push/releases)

**OR**

```sh
go get github.com/meain/push
cd $GOPATH/src/github.com/meain/push
go install cmd/push/push.go
```

## Configure

Once installed run 
```bash
push conf
```
It will guide you through the configuration. It will guide you through the configuration

## Usage

```sh
# push note
push note "Get Schwifty"

# push file
push file todo.md

# push folder ( will upload as zip )
push file todoDir

# Change configuration of push
push conf [deivice|token]
```
