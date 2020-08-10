// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package pnm

import "image"

func init() {
	image.RegisterFormat("pnm", "P?", Decode, DecodeConfig)
}

// check panics if the given error is non-nil.
// This is used by the encoder and decoder.
func check(err error) {
	if err != nil {
		panic(err)
	}
}
