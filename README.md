## gpnm

This package implements an encoder and decoder for `PNM` image formats.
It can be used with Go's image library. It covers all formats as defined
by the 'P1' to 'P6' specifications.


### File formats

Each format differs in what colors it is designed to represent:

* __PBM__ is for bitmaps (black and white, no grays).
* __PGM__ is for grayscale.
* __PPM__ is for "pixmaps" which represent full RGB color.

Each file starts with a two-byte magic number (in ASCII) that identifies the
type of file it is (PBM, PGM, and PPM) and its encoding (ASCII or binary).
The magic number is a capital P followed by a single-digit number.

Magic Number | Type    | Encoding
-------------|---------|-------------
P1           | bitmap  | ASCII
P2           | graymap | ASCII
P3           | pixmap  | ASCII
P4           | bitmap  | Binary
P5           | graymap | Binary
P6           | pixmap  | Binary


The ASCII formats allow for human readability and easy transfer to other
platforms (so long as those platforms understand ASCII), while the binary
formats are more efficient both in file size and in ease of parsing, due
to the absence of whitespace.

In the binary formats, PBM uses 1 bit per pixel, PGM uses 8 bits per pixel,
and PPM uses 24 bits per pixel: 8 for red, 8 for green, 8 for blue.


### Usage

    go get github.com/hexaflex/gpnm


    import "image"
    import _ "github.com/hexaflex/gpnm"

    ...
    img, format, err := image.Decode("myfile.pnm)
    ...

### References

* [Netpbm](http://en.wikipedia.org/wiki/Netpbm_format)


### License

Unless otherwise stated, all of the work in this project is subject to a
3-clause BSD license. Its contents can be found in the enclosed LICENSE file.

