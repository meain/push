package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jhoonb/archivex"

	configure "github.com/meain/push/configure"
	pushbullet "github.com/meain/push/pushbullet"
	transfersh "github.com/meain/push/transfersh"
)

func createZip(source string) string {
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
		s := transfersh.TransferFile(file)
		return s, nil
	case mode.IsRegular():
		fmt.Println("Pushing file: ", os.Args[2])
		s := transfersh.TransferFile(name)
		return s, nil
	}

	return "", errors.New("Could not detect file type")
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

func getDeviceChoice(token string) (string, error) {
	fmt.Println("\n\nCHOOSE DEVICE")
	devices := pushbullet.GetDevices(token)
	var usableDevices []pushbullet.Device
	for _, v := range devices {
		if len(v.Nick) > 0 {
			usableDevices = append(usableDevices, v)
		}
	}
	for i, d := range usableDevices {
		if i < 9 {
			fmt.Println("0" + strconv.Itoa(i+1) + "| " + d.Nick + "  ( " + d.Iden + " )")
		} else { // Hopefully people are sane
			fmt.Println(strconv.Itoa(i+1) + "| " + d.Nick + "  ( " + d.Iden + " )")
		}
	}
	fmt.Print("Enter device number: ")
	var choice string
	fmt.Scanln(&choice)
	choiceNumber, err := strconv.Atoi(choice)
	if err != nil {
		log.Fatal(err)
	}
	iden := usableDevices[choiceNumber-1].Iden
	return iden, nil
}

func getUserToken() string {
	fmt.Println("ADD TOKEN")
	fmt.Println("Get your pushbullet api token from here:")
	fmt.Println("https://www.pushbullet.com/#settings/account")
	fmt.Print("Enter token: ")
	var token string
	fmt.Scanln(&token)
	return token
}

func printUsage() {
	fmt.Println(`Usage: push [command] ...
Available commands:

note - sens a note to your pushbullet device
	 - eg: 'push note "Get Schwifty"'
file - uploads a file to transfer.sh and sends link to pushbullet device
	 - eg 'push file "Assignment.docx"'
conf - configure push vriables ( writes to ~/.push.yaml file )
	- 'push conf' (configure both token and default device)
	- 'push conf token' (configure pushbullet token )
	- 'push conf device' (configure pushbullet device )`)
}

func main() {
	var conf configure.Conf
	confFile := "~/.push.yaml"
	err := conf.GetConf(confFile)
	if err != nil {
		fmt.Println("Could not read config file. Set it up here")
		fmt.Println("")

		token := getUserToken()
		conf.SetToken(confFile, token)

		iden, err := getDeviceChoice(conf.Token)
		if err != nil {
			log.Fatal(err)
		}
		conf.SetDefaultDevice(confFile, iden)
	}

	switch os.Args[1] {
	case "note":
		if len(os.Args) != 3 {
			printUsage()
			break
		}
		fmt.Println("Pushing note: ", os.Args[2])
		err := pushbullet.PushNote("Message", os.Args[2], conf.Token, conf.Device)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("✓ Push complete")
	case "file":
		if len(os.Args) != 3 {
			printUsage()
			break
		}
		err := transferPushFile(conf.Token, conf.Device, os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("✓ Push complete")
	case "conf":
		if len(os.Args) == 2 {
			token := getUserToken()
			conf.SetToken(confFile, token)

			iden, err := getDeviceChoice(conf.Token)
			if err != nil {
				log.Fatal(err)
			}
			conf.SetDefaultDevice(confFile, iden)
		} else if len(os.Args) == 3 {
			if os.Args[2] == "device" {
				iden, err := getDeviceChoice(conf.Token)
				if err != nil {
					log.Fatal(err)
				}
				conf.SetDefaultDevice(confFile, iden)
			} else if os.Args[2] == "token" {
				token := getUserToken()
				conf.SetToken(confFile, token)
			} else {
				printUsage()
			}
		} else {
			printUsage()
		}
	default:
		printUsage()
	}
}
