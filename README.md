# Push

Pushbullet + Transfer.sh

Just push a simple note to pushbullet or send a file link to pushbullet after uploading it to transfer.sh

## Installation

```sh
go get github.com/meain/push
cd $GOPATH/src/github.com/meain/push
go install
```

## Usage

```sh
# push note
push note "Get Schwifty"

# push file
push file todo.md
```
