package main

import (
	"errors"
	"fmt"
	"github.com/fuzxxl/freefare/0.3/freefare"
	"github.com/fuzxxl/nfc/2.0/nfc"
	"github.com/jroimartin/gocui"
	"net"
	"os"
	"time"
)

//#include <stdlib.h>
import "C"

var targetChan = make(chan Tag)
var emptyChan = make(chan struct{})
var device nfc.Device

const targetLoopTimer = 50 * time.Millisecond
const displayHost = "localhost:1337"

type Exit struct{ Code int }

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e)
	}
}

func main() {
	defer handleExit()
	var err error

	// Set up NFC reader. Empty string = first found device
	device, err = nfc.Open("")
	if err != nil {
		fmt.Printf("Error opening NFC device:\n%s", err)

	}

	// "Initiator" mode means passive reader.
	if err = device.InitiatorInit(); err != nil {
		fmt.Printf("Error setting initiator mode:\n%s", err)
		panic(Exit{1})
	}

	if gui, err = setupDisplay(); err != nil {
		fmt.Printf("Error setting up GUI display:\n%s", err)
		panic(Exit{1})
	}
	defer gui.Close()

	addLogRow(fmt.Sprintf("Connected as initiator: %s", device.Connection()))
	addLogRow("Polling for tags now...")

	go pollDisplayState()
	go pollTags()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Printf("Error in GUI main loop:\n%s", err)
		panic(Exit{1})
	}
}

func pollDisplayState() {
	dirty := false
	for {
		select {
		case t := <-targetChan:
			if t.Confidential && t.Authorized {
				setDisplay('$')
			} else if t.Authorized {
				setDisplay('1')
			} else {
				setDisplay('2')
			}
			dirty = true
			addLogRow(fmt.Sprintf("%#v, allowed?: %t", t, t.Authorized))
		case <-emptyChan:
			if dirty {
				setDisplay('0')
				dirty = false
			}
		default:
		}
	}
}

func setDisplay(state rune) error {
	validStates := map[rune]string{
		'0': "clear",
		'1': "success",
		'2': "failure",
		'$': "DOLLA DOLLA BILLS, Y'ALL",
	}

	if _, contains := validStates[state]; !contains {
		return errors.New("state must be 0/1/2/$")
	}

	conn, err := net.Dial("tcp", displayHost)
	if err != nil {
		return err
	}

	defer conn.Close()

	conn.Write([]byte{byte(state)})
	return nil
}

func pollTags() {
	var lastUID = ""

	for {
		time.Sleep(targetLoopTimer)

		tags, err := freefare.GetTags(device)
		if err != nil {
			emptyChan <- struct{}{}
			addLogRow(fmt.Sprintf("Can't read tags: %s", err))
			continue
		}

		if len(tags) != 1 {
			lastUID = ""
			emptyChan <- struct{}{}
			continue
		}

		tag := tags[0]
		/*
			tag, success := tags[0].(freefare.ClassicTag)
			if !success {
				// Probably not a MiFare Classic. Nice try.
				continue
			}
		*/

		if lastUID == tag.UID() {
			// Skip, they're holding the tag in front of the reader. Otherwise
			// we get all sorts of weird shit happening.
			continue
		}

		if err = tag.Connect(); err != nil {
			// Can't activate the tag for some reason
			addLogRow(fmt.Sprintf("Can't activate/connect tag: %s\n", err))
			continue
		}

		auth, conf, block0 := allowAccess(tag)

		t := Tag{
			UID:          tag.UID(),
			Type:         tag.Type(),
			String:       tag.String(),
			Authorized:   auth,
			Confidential: conf,
			Block0:       block0,
		}

		targetChan <- t
		lastUID = t.UID

		if err = tag.Disconnect(); err != nil {
			// fmt.Printf("Can't deactivate/disconnect tag: %s\n", err)
			continue
		}
	}
}
