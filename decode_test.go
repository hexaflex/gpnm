package pnm

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

type testCase struct {
	Name   string
	Data   []byte
	Width  int
	Height int
}

var tests = []testCase{
	{
		Name: "P1",
		Data: []byte(`P1
		# This is an example bitmap of the letter "J"
		6 10
		0 0 0 0 1 0
		0 0 0 0 1 0
		0 0 0 0 1 0
		0 0 0 0 1 0
		0 0 0 0 1 0
		0 0 0 0 1 0
		1 0 0 0 1 0
		0 1 1 1 0 0
		0 0 0 0 0 0
		0 0 0 0 0 0
		`),
		Width:  6,
		Height: 10,
	},
	{
		Name: "P2",
		Data: []byte(`P2
		# Shows the word "FEEP" (example from Netpbm man page on PGM)
		24 7
		15
		0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
		0  3  3  3  3  0  0  7  7  7  7  0  0 11 11 11 11  0  0 15 15 15 15  0
		0  3  0  0  0  0  0  7  0  0  0  0  0 11  0  0  0  0  0 15  0  0 15  0
		0  3  3  3  0  0  0  7  7  7  0  0  0 11 11 11  0  0  0 15 15 15 15  0
		0  3  0  0  0  0  0  7  0  0  0  0  0 11  0  0  0  0  0 15  0  0  0  0
		0  3  0  0  0  0  0  7  7  7  7  0  0 11 11 11 11  0  0 15  0  0  0  0
		0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
		`),
		Width:  24,
		Height: 7,
	},
	{
		Name: "P3",
		Data: []byte(`P3
		# P3 means colors are in ASCII, then 3 columns and 2 rows,
		# then 255 for max color, then RGB triplets
		3 2
		255
		255   0   0     0 255   0     0   0 255
		255 255   0   255 255 255     0   0   0
		`),
		Width:  3,
		Height: 2,
	},
	{
		Name: "P4",
		Data: []byte(`P4
		# This is an example bitmap of the letter "J"
		6 10 ` + "\x08\x08\x08\x08\x08\x08\x88\x70\x00\x00"),
		Width:  6,
		Height: 10,
	},
	{
		Name: "P5",
		Data: []byte(`P5
		# Shows the word "FEEP" (example from Netpbm man page on PGM)
		24 7
		15 ` +
			"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x03\x03\x03\x03\x00\x00\x07\x07\x07\x07\x00\x00\x0b\x0b\x0b\x0b\x00\x00\x0f\x0f\x0f\x0f\x00" +
			"\x00\x03\x00\x00\x00\x00\x00\x07\x00\x00\x00\x00\x00\x0b\x00\x00\x00\x00\x00\x0f\x00\x00\x0f\x00" +
			"\x00\x03\x03\x03\x00\x00\x00\x07\x07\x07\x00\x00\x00\x0b\x0b\x0b\x00\x00\x00\x0f\x0f\x0f\x0f\x00" +
			"\x00\x03\x00\x00\x00\x00\x00\x07\x00\x00\x00\x00\x00\x0b\x00\x00\x00\x00\x00\x0f\x00\x00\x00\x00" +
			"\x00\x03\x00\x00\x00\x00\x00\x07\x07\x07\x07\x00\x00\x0b\x0b\x0b\x0b\x00\x00\x0f\x00\x00\x00\x00" +
			"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
		Width:  24,
		Height: 7,
	},
	{
		Name: "P6",
		Data: []byte(`P6
		# P6 means colors are in binary, then 3 columns and 2 rows,
		# then 255 for max color, then RGB triplets
		3 2
		255 ` + "\xff\x00\x00\x00\xff\x00\x00\x00\xff\xff\xff\x00\xff\xff\xff\x00\x00\x00"),
		Width:  3,
		Height: 2,
	},
}

func TestDecode(t *testing.T) {
	for _, tc := range tests {
		img, fmt, err := image.Decode(bytes.NewBuffer(tc.Data))
		if err != nil {
			t.Fatal(err)
		}

		if fmt != "pnm" {
			t.Fatalf("%s: Format mismatch; expected \"ppm\", have %q", tc.Name, fmt)
		}

		saveImg(img, tc.Name)
	}
}

func TestDecodeConfig(t *testing.T) {
	for _, tc := range tests {
		cfg, fmt, err := image.DecodeConfig(bytes.NewBuffer(tc.Data))
		if err != nil {
			t.Fatal(err)
		}

		if fmt != "pnm" {
			t.Fatalf("%s: Format mismatch; expected \"ppm\", have %q", tc.Name, fmt)
		}

		if cfg.Width != tc.Width {
			t.Fatalf("%s: Width mismatch; expected %d, have %d", tc.Name, tc.Width, cfg.Width)
		}

		if cfg.Height != tc.Height {
			t.Fatalf("%s: Height mismatch; expected %d, have %d", tc.Name, tc.Height, cfg.Height)
		}
	}
}

func saveImg(img image.Image, name string) {
	fd, err := os.Create(filepath.Join("testdata", name+".png"))
	if err != nil {
		return
	}

	defer fd.Close()
	png.Encode(fd, img)
}
