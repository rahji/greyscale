// assumes an image with 256 greys

package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

var scale = []string{
	"Black",
	"Very Dark Gray",
	"Dark Gray",
	"Medium Dark Gray",
	"Slate Gray",
	"Dim Gray",
	"Light Slate Gray",
	"Gray",
	"Light Gray",
	"Gainsboro",
	"Silver",
	"Light Silver",
	"Very Light Gray",
	"Near White",
	"Off White",
	"White",
}

func main() {
	// get the filename as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: greyscale <filename>")
		return
	}
	filename := os.Args[1]

	// read the image file
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	bounds := m.Bounds()
	numpix := bounds.Max.X * bounds.Max.Y

	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	var histogram [16]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := m.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 12 reduces this to the range [0, 15].
			histogram[r>>12]++
		}
	}

	var out strings.Builder
	out.WriteString("|Color|Pixels|Percent|\n")
	out.WriteString("|----:|-----:|------:|\n")
	for i, x := range histogram {
		pct := float64(x) / float64(numpix) * 100
		out.WriteString(fmt.Sprintf("|%s|%d|%.02f%%|\n", scale[i], x, pct))
	}
	md, err := glamour.Render(out.String(), "dark")
	fmt.Print(md)
}
