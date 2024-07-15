/*
Copyright © 2024 Rob Duarte <me@robduarte.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var html bool
var x int
var y int

// colorsCmd represents the colors command
var pickCmd = &cobra.Command{
	Use:   "pick",
	Short: "show the exact greyscale color of a pixel",
	Long: `
This produces an a hex value for the color of the specifed pixel
Again, it assumes the image is actually greyscale. 
`,
	Run: func(cmd *cobra.Command, args []string) {

		m, _, err := readImage(infile)
		if err != nil {
			log.Fatal(err)
		}

		bounds := m.Bounds()

		if x > bounds.Max.X {
			log.Fatal(fmt.Errorf("x value is too large"))
		}
		if y > bounds.Max.Y {
			log.Fatal(fmt.Errorf("y value is too large"))
		}

		r, _, _, _ := m.At(x, y).RGBA()
		grey := r >> 8 // a right-shift of 8 turns 65535 max to 255 max
		if html {
			fmt.Printf("#%x%x%x\n", grey, grey, grey)
		} else {
			fmt.Println(grey)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(pickCmd)
	pickCmd.PersistentFlags().StringVarP(&infile, "infile", "i", "", "input file (required)")
	pickCmd.PersistentFlags().IntVarP(&x, "x", "x", 0, "x value for the pixel to be examined (required)")
	pickCmd.PersistentFlags().IntVarP(&y, "y", "y", 0, "y value for the pixel to be examined (required)")
	pickCmd.PersistentFlags().BoolVar(&html, "html", false, "output as an HTML hex string")
	pickCmd.MarkPersistentFlagRequired("infile")
	pickCmd.MarkPersistentFlagRequired("x")
	pickCmd.MarkPersistentFlagRequired("y")
}
