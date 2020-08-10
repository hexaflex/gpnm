package pnm

import (
	"fmt"
	"image"
	"io"
	"math"
	"strings"
)

type PNMType int

// Known PNM file types.
const (
	BitmapAscii PNMType = iota
	BitmapBinary
	GraymapAscii
	GraymapBinary
	PixmapAscii
	PixmapBinary
)

// Encode writes the Image m to w in PPM format.
// The type of output file is determined by the given PNM type value.
// Any image can be encoded, but depending on the chosen type,
// the encoding may be lossy.
func Encode(w io.Writer, m image.Image, ptype PNMType) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("DecodeConfig: %v", x)
		}
	}()

	switch ptype {
	case BitmapAscii:
		encodeP1(w, m)
	case GraymapAscii:
		encodeP2(w, m)
	case PixmapAscii:
		encodeP3(w, m)
	case BitmapBinary:
		encodeP4(w, m)
	case GraymapBinary:
		encodeP5(w, m)
	case PixmapBinary:
		encodeP6(w, m)
	default:
		return fmt.Errorf("Invalid PPM type %d", ptype)
	}

	return
}

func encodeP1(w io.Writer, m image.Image) {
	b := m.Bounds()
	row := make([]string, b.Dx())

	write(w, "P1\n%d %d\n", b.Dx(), b.Dy())

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, _, _, _ := m.At(x, y).RGBA()

			if byte(r)/0xff == 1 {
				row[x] = "1"
			} else {
				row[x] = "0"
			}
		}

		write(w, "%s\n", strings.Join(row, " "))
	}
}

func encodeP2(w io.Writer, m image.Image) {
	b := m.Bounds()
	row := make([]string, b.Dx())

	write(w, "P2\n%d %d\n255\n", b.Dx(), b.Dy())

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, _, _, _ := m.At(x, y).RGBA()
			row[x] = fmt.Sprintf("%d", byte(r))
		}

		write(w, "%s\n", strings.Join(row, " "))
	}
}

func encodeP3(w io.Writer, m image.Image) {
	b := m.Bounds()
	row := make([]string, b.Dx())

	write(w, "P3\n%d %d\n255\n", b.Dx(), b.Dy())

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			row[x] = fmt.Sprintf("%d %d %d", byte(r), byte(g), byte(b))
		}

		write(w, "%s\n", strings.Join(row, " "))
	}
}

func encodeP4(w io.Writer, m image.Image) {
	var bit int

	b := m.Bounds()
	bytes := int(math.Ceil((float64(b.Dx()) / 8)))
	bits := newBitset(uint(bytes * b.Dy() * 8))
	pad := (bytes * 8) - b.Dx()

	write(w, "P4\n%d %d\n", b.Dx(), b.Dy())

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, _, _, _ := m.At(x, y).RGBA()

			if r > 0 {
				bits.Set(bit)
			}

			bit++
		}

		bit += pad
	}

	_, err := w.Write(bits)
	check(err)
}

func encodeP5(w io.Writer, m image.Image) {
	b := m.Bounds()
	data := make([]byte, 0, b.Dx()*b.Dy())

	write(w, "P5\n%d %d\n255\n", b.Dx(), b.Dy())

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, _, _, _ := m.At(x, y).RGBA()
			data = append(data, byte(r))
		}
	}

	_, err := w.Write(data)
	check(err)
}

func encodeP6(w io.Writer, m image.Image) {
	b := m.Bounds()
	data := make([]byte, 0, b.Dx()*b.Dy()*3)

	write(w, "P6\n%d %d\n255\n", b.Dx(), b.Dy())

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			data = append(data, byte(r), byte(g), byte(b))
		}
	}

	_, err := w.Write(data)
	check(err)
}

func write(w io.Writer, f string, argv ...interface{}) {
	_, err := fmt.Fprintf(w, f, argv...)
	check(err)
}
