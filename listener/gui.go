package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui

func setupDisplay() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}

	// Set manager/keybinding
	//g.Cursor = true
	g.Highlight = true
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.SetManagerFunc(layout)

	if err = keybindings(g); err != nil {
		return nil, err
	}

	return g, nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("log", -1, -1, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Wrap = true
		v.Title = "Log"
		//v.SelBgColor = gocui.ColorGreen
		//v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("settings", maxX/2, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = true
		v.Title = "Settings"

		fmt.Fprintln(v, "Errrrr")
		fmt.Fprintln(v, "There are meant to be widgets here. TBA.")

	}
	return nil
}

func addLogRow(s string) {
	gui.Update(func(g *gocui.Gui) error {
		v, _ := g.View("log")
		fmt.Fprintln(v, s)
		return nil
	})
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "log" {
		_, err := g.SetCurrentView("settings")
		return err
	}
	_, err := g.SetCurrentView("log")
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("log", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("settings", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}
