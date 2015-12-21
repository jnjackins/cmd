package main

import (
	"image"
	"image/color"
	"log"

	"sigint.ca/graphics/text"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	editbuf *text.Buffer
	scr     screen.Screen
	win     screen.Window
	winr    image.Rectangle
	bgColor = color.White
)

func init() {
	log.SetFlags(0)
}

func main() {
	size := image.Pt(800, 600)
	font := basicfont.Face7x13
	height := font.Height
	var err error
	editbuf = text.NewBuffer(size, font, height, text.AcmeYellowTheme)
	if err != nil {
		log.Fatal(err)
	}

	driver.Main(func(s screen.Screen) {
		scr = s
		if w, err := scr.NewWindow(nil); err != nil {
			log.Fatal(err)
		} else {
			win = w
			defer win.Release()
			if err := eventLoop(); err != nil {
				log.Print(err)
			}
		}
	})
}

func eventLoop() error {
	for e := range win.Events() {
		switch e := e.(type) {
		case key.Event:
			if e.Direction == key.DirPress || e.Direction == key.DirNone {
				editbuf.SendKey(e)
				win.Send(paint.Event{})
			}

		case mouse.Event:
			editbuf.SendMouseEvent(e)
			win.Send(paint.Event{})

		case paint.Event:
			win.Fill(image.Rect(0, 0, 2000, 2000), color.RGBA{G: 255, A: 255}, screen.Src)
			win.Upload(image.ZP, editbuf, editbuf.Bounds())
			win.Publish()

		case size.Event:
			winr = e.Bounds()
			editbuf.Resize(e.Size())
			win.Send(paint.Event{})

		case lifecycle.Event:

		default:
			log.Printf("unhandled %T: %[1]v", e)
		}

	}
	return nil
}
