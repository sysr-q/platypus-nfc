package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"net"
	"os"
	"time"
)

var screen tcell.Screen
var currentState = Nothing
var dirty = false

var (
	Nothing = tcell.
		StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlack)
	Success = tcell.
		StyleDefault.
		Background(tcell.ColorGreen).
		Foreground(tcell.ColorGreen)
	Failure = tcell.
		StyleDefault.
		Background(tcell.ColorRed).
		Foreground(tcell.ColorRed)
	DOLLARS = tcell.
		StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorYellowGreen)
)

const DOLLADOLLABILLSYALL = '$'

func main() {
	var err error
	if screen, err = tcell.NewScreen(); err != nil {
		fmt.Printf("Could not start tcell, NewScreen() gave error:\n%s", err)
		os.Exit(1)
	}

	if err = screen.Init(); err != nil {
		fmt.Printf("Could not start tcell, Init() gave error:\n%s", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		fmt.Printf("Could not listen on port :1337:\n%s", err)
		os.Exit(1)
	}

	// Set up to poll events in tcell.
	quit := make(chan struct{})
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					dirty = true
					screen.Sync()
				}
			case *tcell.EventResize:
				dirty = true
				screen.Sync()
			}
		}
	}()

	screen.HideCursor()
	screen.SetStyle(Nothing)
	screen.Clear()

	defer screen.Fini()

	// Listen for sockets to change the state.
	go listen(ln)

	// Flush the display if it's dirty
	displayFlush := 50 * time.Millisecond
	go func() {
		for {
			time.Sleep(displayFlush)
			if dirty {
				screen.SetStyle(currentState)

				if currentState == DOLLARS {
					xs, ys := screen.Size()
					for y := 0; y < ys; y++ {
						for x := 0; x < xs; x++ {
							screen.SetContent(x, y, DOLLADOLLABILLSYALL, nil, currentState)
						}
					}
				} else {
					screen.Clear()
				}

				screen.Show()
				dirty = false
			}
		}
	}()

	// Stop the program from exiting without us asking it to.
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
	}
}

func listen(ln net.Listener) {
	timeout := 5 * time.Second
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// Fuck latency.
		conn.SetDeadline(time.Now().Add(timeout))

		// I literally don't care if multiple sockets are open at once.
		// You're just fighting yourself to change the state. Idiot.
		go changeState(conn)
	}
}

func changeState(conn net.Conn) {
	defer conn.Close()
	var p [1]byte
	_, err := conn.Read(p[:])
	if err != nil {
		return
	}
	if p[0] == '0' {
		currentState = Nothing
	} else if p[0] == '1' {
		currentState = Success
	} else if p[0] == '2' {
		currentState = Failure
	} else if p[0] == '$' {
		currentState = DOLLARS
	}

	dirty = true
}
