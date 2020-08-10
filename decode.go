package pnm

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"
	"strconv"
)

// DecodeConfig decodes image configuration/metadata.
func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("DecodeConfig: %v", x)
		}
	}()

	br := bufio.NewReader(r)
	format := readS(br)

	switch format {
	case "P1", "P2", "P3", "P4", "P5", "P6":
	default:
		panic("Unknown PPM format: " + format)
	}

	cfg.ColorModel = color.RGBAModel
	cfg.Width = int(readU(br))
	cfg.Height = int(readU(br))
	return
}

// Decode decodes an image from the given data.
func Decode(r io.Reader) (img image.Image, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Decode: %v", x)
		}
	}()

	br := bufio.NewReader(r)
	format := readS(br)
	width := readU(br)
	height := readU(br)
	rect := image.Rect(0, 0, int(width), int(height))

	switch format {
	case "P1":
		img = image.NewAlpha(rect)
		decodeP1(br, img.(draw.Image), width, height)
	case "P2":
		img = image.NewGray(rect)
		decodeP2(br, img.(draw.Image), width, height)
	case "P3":
		img = image.NewRGBA(rect)
		decodeP3(br, img.(draw.Image), width, height)
	case "P4":
		img = image.NewAlpha(rect)
		decodeP4(br, img.(draw.Image), width, height)
	case "P5":
		img = image.NewGray(rect)
		decodeP5(br, img.(draw.Image), width, height)
	case "P6":
		img = image.NewRGBA(rect)
		decodeP6(br, img.(draw.Image), width, height)
	default:
		panic("Unknown PPM format: " + format)
	}

	return
}

// decodeP1 reads an ASCII bitmap
func decodeP1(r *bufio.Reader, img draw.Image, width, height uint) {
	var x, y int

	for y = 0; y < int(height); y++ {
		for x = 0; x < int(width); x++ {
			img.Set(x, y, color.Alpha{uint8(readU(r)) * 0xff})
		}
	}
}

// decodeP2 reads an ASCII graymap
func decodeP2(r *bufio.Reader, img draw.Image, width, height uint) {
	var x, y int

	mask := readU(r)
	mul := 255 / mask

	for y = 0; y < int(height); y++ {
		for x = 0; x < int(width); x++ {
			img.Set(x, y, color.Gray{uint8((readU(r) & mask) * mul)})
		}
	}
}

// decodeP3 reads an ASCII pixmap
func decodeP3(r *bufio.Reader, img draw.Image, width, height uint) {
	var x, y int

	mask := readU(r)
	mul := 255 / mask

	for y = 0; y < int(height); y++ {
		for x = 0; x < int(width); x++ {
			img.Set(x, y, color.RGBA{
				uint8((readU(r) & mask) * mul),
				uint8((readU(r) & mask) * mul),
				uint8((readU(r) & mask) * mul),
				0xff,
			})
		}
	}
}

// decodeP4 reads a binary bitmap
func decodeP4(r *bufio.Reader, img draw.Image, width, height uint) {
	var x, y, bit int

	bytes := int(math.Ceil((float64(width) / 8)))
	bits := newBitset(uint(bytes) * height * 8)
	pad := (bytes * 8) - int(width)

	space(r)

	_, err := r.Read(bits)
	check(err)

	for y = 0; y < int(height); y++ {
		for x = 0; x < int(width); x++ {
			if bits.Test(bit) {
				img.Set(x, y, color.Alpha{0xff})
			} else {
				img.Set(x, y, color.Alpha{0x00})
			}

			bit++
		}

		bit += pad
	}
}

// decodeP5 reads a binary graymap
func decodeP5(r *bufio.Reader, img draw.Image, width, height uint) {
	var x, y, pix int

	mask := byte(readU(r))
	mul := 255 / mask
	data := make([]byte, width*height)

	space(r)

	_, err := r.Read(data)
	check(err)

	for y = 0; y < int(height); y++ {
		for x = 0; x < int(width); x++ {
			img.Set(x, y, color.Gray{(data[pix] & mask) * mul})
			pix++
		}
	}
}

// decodeP6 reads a binary pixmap
func decodeP6(r *bufio.Reader, img draw.Image, width, height uint) {
	var x, y, pix int

	mask := byte(readU(r))
	mul := 255 / mask
	data := make([]byte, width*height*3)

	space(r)

	_, err := r.Read(data)
	check(err)

	for y = 0; y < int(height); y++ {
		for x = 0; x < int(width); x++ {
			img.Set(x, y, color.RGBA{
				(data[pix] & mask) * mul,
				(data[pix+1] & mask) * mul,
				(data[pix+2] & mask) * mul,
				0xff,
			})

			pix += 3
		}
	}
}

// readU8 reads the next value as a byte
func readU8(r *bufio.Reader) uint8 {
	b, err := r.ReadByte()
	check(err)
	return b
}

// readU reads the next value as an int
func readU(r *bufio.Reader) uint {
	n, err := strconv.ParseUint(readS(r), 10, 32)
	check(err)
	return uint(n)
}

// readS reads the next value as a string.
func readS(r *bufio.Reader) string {
	return string(read(r))
}

// read reads the next value, ignoring whitespace.
func read(r *bufio.Reader) []byte {
	var comment bool
	space(r)

	buf := make([]byte, 0, 16)

loop:
	for {
		b, err := r.ReadByte()
		check(err)

		switch {
		case b == '#':
			comment = true
		case (comment && isNewline(b)) || (!comment && isSpace(b)):
			check(r.UnreadByte())
			break loop
		}

		if !comment {
			buf = append(buf, b)
		}
	}

	if len(buf) == 0 {
		return read(r)
	}

	return buf
}

// space reads whitespace from the given reader, until
// the stream ends or a non-whitespace character is found.
func space(r *bufio.Reader) {
	for {
		b, err := r.ReadByte()
		check(err)

		if !isSpace(b) {
			check(r.UnreadByte())
			return
		}
	}
}

// isSpace returns true if the given byte represents whitespace:
//
//    ' '   (0x20)	space (SPC)
//    '\t'	(0x09)	horizontal tab (TAB)
//    '\n'	(0x0a)	newline (LF)
//    '\v'	(0x0b)	vertical tab (VT)
//    '\f'	(0x0c)	feed (FF)
//    '\r'	(0x0d)	carriage return (CR)
//
func isSpace(b byte) bool {
	switch b {
	case 0x20, 0xa, 0x9, 0xb, 0xc, 0xd:
		return true
	}
	return false
}

// isNewline returns true if the given byte represents a newline:
//
//    '\n'	(0x0a)	newline (LF)
//    '\r'	(0x0d)	carriage return (CR)
//
func isNewline(b byte) bool {
	switch b {
	case 0xa, 0xd:
		return true
	}
	return false
}
