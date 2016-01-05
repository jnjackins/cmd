package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"

	"sigint.ca/clip"
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
	filename string
	editbuf  *text.Buffer
	scr      screen.Screen
	win      screen.Window
	winr     image.Rectangle
	bgColor  = color.White
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s file ...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
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
	editbuf.Clipboard = &clip.Clipboard{}

	if flag.NArg() == 1 {
		load(flag.Arg(0))
	} else if flag.NArg() > 1 {
		log.Println("multiple files not yet supported")
		flag.Usage()
		os.Exit(1)
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
			if e.Direction == key.DirPress &&
				e.Modifiers == key.ModMeta &&
				e.Code == key.CodeS {
				// meta-s
				save()
			}
			if e.Direction == key.DirPress || e.Direction == key.DirNone {
				editbuf.SendKeyEvent(e)
				win.Send(paint.Event{})
			}

		case mouse.Event:
			editbuf.SendMouseEvent(e)
			win.Send(paint.Event{})

		case paint.Event:
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

func load(s string) {
	filename = s
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Printf("error opening %q for reading: %v", filename, err)
		return
	}
	buf, err := ioutil.ReadFile(filename)
	editbuf.Load(buf)
	f.Close()
}

func save() {
	if filename == "" {
		log.Println("saving untitled file not yet supported")
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("error opening %q for writing: %v", filename, err)
	}
	r := bytes.NewBuffer(editbuf.Contents())
	if _, err := io.Copy(f, r); err != nil {
		log.Printf("error writing to %q: %v", filename, err)
	}
	f.Close()
}
