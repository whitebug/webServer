package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	server3()
}

func server3() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		lissajous(writer, 5)
	})
	http.HandleFunc("/cycles", lissajousParameters)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajousParameters(writer http.ResponseWriter, request *http.Request) {
	for key, value := range request.Form {
		if key == "cycles" {
			if val, err := strconv.ParseFloat(value[0], 64); err == nil {
				fmt.Println(err)
			} else {
				lissajous(writer, val)
			}
		}
	}
}

func handler3(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "%s %s %s\n", request.Method, request.URL, request.Proto)
	for k, v := range request.Header {
		fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(writer, "Host = %q\n", request.Host)
	fmt.Fprintf(writer, "RemoteAddr = %q\n", request.RemoteAddr)
	if err := request.ParseForm(); err != nil {
		log.Print(err)
	}
	for key, value := range request.Form {
		fmt.Fprintf(writer, "Form[%q] = %q\n", key, value)
	}
}

var green = color.RGBA{R: 70, G: 180, B: 125, A: 1}
var someColor = color.RGBA{R: 30, G: 20, B: 40, A: 1}
var palette = []color.Color{color.White, color.Black, green, someColor}

const (
	whiteIndex = 0
	blackIndex = 1
	greenIndex = 2
)

func lissajous(out io.Writer, cycles float64) {
	const (
		res     = 0.001
		size    = 100
		nFrames = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nFrames}
	phase := 0.0
	finalSize := 2*size + 1
	for i := 0; i < nFrames; i++ {
		rect := image.Rect(0, 0, finalSize, finalSize)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

var mutex sync.Mutex
var count int

func server2() {
	http.HandleFunc("/", handler2)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler2(writer http.ResponseWriter, request *http.Request) {
	mutex.Lock()
	count++
	mutex.Unlock()
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}

func counter(writer http.ResponseWriter, request *http.Request) {
	mutex.Lock()
	fmt.Fprintf(writer, "Count %d\n", count)
	mutex.Unlock()
}

func server1() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Returns path component from URL request
func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}
