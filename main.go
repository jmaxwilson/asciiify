package main

import (
	"fmt"
	"image/color"
	"image/gif"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rthornton128/goncurses"
	"github.com/theMomax/asciiify/asciiif"
	"github.com/tomnomnom/xtermcolor"
)

func main() {

	var loopCount int
	var gifURL string
	var err error

	loopCount = -1

	// Process Commandline Args
	if len(os.Args) == 1 {
		fmt.Println("Usage: asciiify [GIF URL] [LOOP COUNT]")
		return
	}
	if len(os.Args) > 1 {
		gifURL = os.Args[1]
	}
	if len(os.Args) > 2 {
		loopCount, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid Loop Count")
			return
		}
		if loopCount <= 0 {
			fmt.Println("Loop Count must be greater than 0")
			return
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigs {

		}
	}()

	// Initialize NCurses
	stdscr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	defer goncurses.End()

	goncurses.Cursor(0)
	stdscr.ScrollOk(false)

	// Read the Gif from the URL specified on the commandline
	httpClient := http.Client{}
	response, err := httpClient.Get(gifURL)
	if err != nil {
		panic(err)
	}

	// Decode the GIF data from the response body
	gif, err := gif.DecodeAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Initialize the Loop count
	gif.LoopCount = loopCount

	// Play the GIF as an ASCII NCurses Video
	vid := asciiif.DecodeGIFStreamed(gif)

	goncurses.StartColor()

	for i := int16(0); i < 256; i++ {
		goncurses.InitPair(i, i, 0)
	}

	for img := range vid {
		start := time.Now()
		for y, row := range img.Image {
			for x, pix := range row {
				stdscr.ColorOn(int16(xtermcolor.FromColor(color.RGBA{pix.R, pix.G, pix.B, pix.A})))
				stdscr.MoveAddChar(y, x, goncurses.Char(pix.Char))
			}
		}
		stdscr.Refresh()
		took := time.Since(start)
		delay := time.Duration(img.Delay) * 10 * time.Millisecond
		if took < delay {
			time.Sleep(delay - took)
		}
	}
}
