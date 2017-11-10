package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jhoonb/archivex"

	configure "github.com/meain/push/configure"
	pushbullet "github.com/meain/push/pushbullet"
	transfersh "github.com/meain/push/transfersh"
)

func createZip(source string) string{
	zipfile := "/tmp/" + filepath.Base(source) + ".zip"
	zip := new(archivex.ZipFile)
	zip.Create(zipfile)
	zip.AddAll(source, true)
	zip.Close()
	return zipfile
}

func transferFileOrFolder(name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		file := createZip(name)
		fmt.Println("Pushing folder: ", os.Args[2])
		return transfersh.TransferFile(file), nil
	case mode.IsRegular():
		fmt.Println("Pushing file: ", os.Args[2])
		return transfersh.TransferFile(name), nil
	}

	return "", errors.New("Could not upload file to transfer.sh")
}

func transferPushFile(token string, device string, fileName string) error {
	s, err := transferFileOrFolder(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	err = pushbullet.PushNote("Message", s, token, device)
	if err != nil {
		return err
	}
	return nil
}

func printUsage() {
	fmt.Println(`Usage: push [command] ...
Available commands:
	note - sens a note to your pushbullet device
		 - eg: 'push note "Get Schwifty"'
	file - uploads a file to transfer.sh and sends link to pushbullet device
		 - eg 'push file "Assignment.docx"'`)
}

func main() {
	var conf configure.Conf
	conf.GetConf("~/.push.yaml")

	if len(os.Args) != 3 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "note":
		fmt.Println("Pushing note: ", os.Args[2])
		err := pushbullet.PushNote("Message", os.Args[2], conf.Token, conf.Device)
		if err != nil {
			log.Fatal(err)
		}
	case "file":
		err := transferPushFile(conf.Token, conf.Device, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	default:
		printUsage()
	}
}
