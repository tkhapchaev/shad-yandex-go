//go:build !solution

package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	height = 12
	width  = 8
	colon  = 4
)

func Digit(digit int) string {
	switch digit {
	case 0:
		return Zero
	case 1:
		return One
	case 2:
		return Two
	case 3:
		return Three
	case 4:
		return Four
	case 5:
		return Five
	case 6:
		return Six
	case 7:
		return Seven
	case 8:
		return Eight
	case 9:
		return Nine
	default:
		return ""
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var t string
	var k int

	urlStr := "http://" + r.Host + r.URL.String()
	u, err := url.Parse(urlStr)

	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()

	if len(q["time"]) != 0 && q["time"][0] != "" {
		t = q["time"][0]
	} else {
		t = time.Now().Format("15:04:05")
	}

	if len(t) != 8 {
		http.Error(w, "invalid time format", 400)
	}

	_, err = time.Parse("15:04:05", t)

	if err != nil {
		http.Error(w, "invalid time format", 400)
	}

	if len(q["k"]) != 0 && q["k"][0] != "" {
		k, err = strconv.Atoi(q["k"][0])
	} else {
		k = 1
	}

	if err != nil || k < 1 || k > 30 {
		http.Error(w, "invalid k value", 400)
	}

	img := Draw(t, k)
	err = png.Encode(w, img)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(200)
}

func DrawDigit(img *image.RGBA, digit int, x int, y int, k int) (*image.RGBA, int, int) {
	i := x
	j := y
	draw := Digit(digit - '0')

	for d, s := range draw {
		if draw[d] == 10 {
			continue
		}

		img = DrawPixel(img, string(s), i*k, j*k, k)

		if i-x == width-1 {
			i = x
			j++
		} else {
			i++
		}
	}

	i = x + width
	j = 0

	return img, i, j
}

func DrawColon(img *image.RGBA, x int, y int, k int) (*image.RGBA, int, int) {
	i := x
	j := y

	for c, s := range Colon {
		if Colon[c] == 10 {
			continue
		}

		img = DrawPixel(img, string(s), i*k, j*k, k)

		if i-x == colon-1 {
			i = x
			j++
		} else {
			i++
		}
	}

	i = x + colon
	j = 0

	return img, i, j
}

func DrawPixel(img *image.RGBA, digit string, x int, y int, k int) *image.RGBA {
	for i := y; i < y+k; i++ {
		for j := x; j < x+k; j++ {
			if digit == "1" {
				img.Set(j, i, Cyan)
			} else {
				img.Set(j, i, color.White)
			}
		}
	}

	return img
}

func Draw(time string, k int) *image.RGBA {
	s := 2*colon + 6*width
	w := s * k
	h := height * k

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	x := 0
	y := 0

	for _, t := range time {
		if t == ':' {
			img, x, y = DrawColon(img, x, y, k)
		} else {
			img, x, y = DrawDigit(img, int(t), x, y, k)
		}
	}

	return img
}

func main() {
	port := flag.Int("port", 800, "port string")
	flag.Parse()
	http.HandleFunc("/", Handle)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(*port), nil))
}
