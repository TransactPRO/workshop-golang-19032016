package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// UI contains the UI data.
type UI struct {
	gui   *gocui.Gui
	sh    shutdownHandler
	msgCh chan string
}

type shutdownHandler func()

const (
	// UsersView defines a view containing a list of users.
	UsersView string = "users"
	// ChatView defines a view containing the main chat.
	ChatView string = "chat"
	// TextView defines a view containing an input text.
	TextView string = "text"
)

// DeployGUI deploys the GUI.
func DeployGUI(sh shutdownHandler, msgCh chan string) (ui *UI, err error) {
	// Initializing a new GUI.
	g := gocui.NewGui()
	err = g.Init()
	if err != nil {
		return
	}

	ui = &UI{
		gui:   g,
		sh:    sh,
		msgCh: msgCh,
	}

	// Setting the desired layout by passing the corresponding handler to the GUI method.
	g.SetLayout(ui.layout)

	// Setting Ctr+C binding.
	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return
	}

	// Setting Ctr+S binding.
	err = g.SetKeybinding(TextView, gocui.KeyCtrlS, gocui.ModNone, ui.processText)
	if err != nil {
		return
	}

	// We want the cursor to be visible.
	g.Cursor = true

	// Executing the main loop within the routine in order to perform
	// further actions.
	go func() {
		// This loop is going to run while the GUI is active (the same logic as in GTK+).
		err := g.MainLoop()
		if err != nil && err != gocui.ErrQuit {
			log.Fatal(err)
		}
		// Getting to this part only after the "quit" handler has been called and closed.
		g.Close()
		ui.sh()
	}()

	return
}

// WriteToView writes the message into the requested view.
func (u *UI) WriteToView(view, msg string) {
	u.gui.Execute(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}
		fmt.Fprintln(v, msg)
		return nil
	})
}

func (u *UI) layout(g *gocui.Gui) (err error) {
	// Retrieving the terminal's size.
	maxX, maxY := g.Size()

	// Setting the users list view.
	_, err = g.SetView(UsersView, 0, 0, maxX/5, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}

	// Setting the chat history view.
	var cView *gocui.View
	cView, err = g.SetView(ChatView, maxX/5+1, 0, maxX-1, maxY*4/5)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}
	cView.Autoscroll = true

	// Setting the text editor view.
	var tView *gocui.View
	tView, err = g.SetView(TextView, maxX/5+1, maxY*4/5+1, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}
	// The text ditor view must be editable.
	tView.Editable = true
	// Also it's better to start the execution with the focus given to the text view.
	err = g.SetCurrentView(TextView)
	if err != nil {
		return
	}

	return
}

func (u *UI) processText(g *gocui.Gui, v *gocui.View) (err error) {
	bufferStr := v.ViewBuffer()
	if bufferStr == "" {
		return
	}
	u.msgCh <- bufferStr[:len(bufferStr)-1]
	v.Clear()
	v.SetCursor(0, 0)
	return
}

// the quit handler is being called as we press the Ctr+C combination and returns
// the corresponding error to the running GUI.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
