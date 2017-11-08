# Push

> [pushbullet](https://www.pushbullet.com/) + [transfer.sh](https://transfer.sh/)

Just push a simple note to [pushbullet](https://www.pushbullet.com/) or send a file link to pushbullet after uploading it to [transfer.sh](https://transfer.sh/)

## Installation

Download binary from [releases page](https://github.com/meain/push/releases)

**OR**

```sh
go get github.com/meain/push
cd $GOPATH/src/github.com/meain/push
go install
```

## Configure

Add a yaml file to your home directory named `~/.push.yaml` with contents like below

```yaml
token: "<YOUR_PUSHBULLET_TOKEN>"
device: "<YOUR_DEVICE_NAME>"
```

You can get pushbullet api token from https://www.pushbullet.com/#settings/account

## Usage

```sh
# push note
push note "Get Schwifty"

# push file
push file todo.md
```
