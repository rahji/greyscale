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
	"image/color"
	"log"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "show filetype, color model, and dimensions of an image",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		m, filetype, err := readImage(infile)
		if err != nil {
			log.Fatal(err)
		}

		bounds := m.Bounds()
		numpix := bounds.Max.X * bounds.Max.Y

		var colorModelName string
		switch m.ColorModel() {
		case color.RGBAModel:
			colorModelName = "RGBA"
		case color.RGBA64Model:
			colorModelName = "RGBA64"
		case color.NRGBAModel:
			colorModelName = "NRGBA"
		case color.NRGBA64Model:
			colorModelName = "NRGBA64"
		case color.AlphaModel:
			colorModelName = "Alpha"
		case color.Alpha16Model:
			colorModelName = "Alpha16"
		case color.GrayModel:
			colorModelName = "Gray"
		case color.Gray16Model:
			colorModelName = "Gray16"
		case color.CMYKModel:
			colorModelName = "CMYK"
		default:
			colorModelName = "Unknown"
		}

		fmt.Printf("Filetype:      %s\n", filetype)
		fmt.Printf("Color Model:   %s\n", colorModelName)
		fmt.Printf("Min Bounds:    %d x %d\n", bounds.Min.X, bounds.Min.Y)
		fmt.Printf("Max Bounds:    %d x %d\n", bounds.Max.X, bounds.Max.Y)
		fmt.Printf("Total Pixels:  %d\n", numpix)
	},
}

func init() {
	showCmd.AddCommand(infoCmd)
}
