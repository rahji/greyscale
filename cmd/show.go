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
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"os"

	"github.com/spf13/cobra"
)

var infile string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show greyscale info about an image",
	Long:  "the subcommands colors and info do the actual work",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'show colors' or 'show info'")
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.PersistentFlags().StringVarP(&infile, "infile", "i", "", "input file (required)")
	showCmd.MarkPersistentFlagRequired("infile")
}

// reads a file named f and returns
// a decoded image, file format, and err
func readImage(f string) (image.Image, string, error) {
	reader, err := os.Open(f)
	if err != nil {
		return nil, "", fmt.Errorf("os.open: %w", err)
	}
	defer reader.Close()

	return image.Decode(reader)
}
