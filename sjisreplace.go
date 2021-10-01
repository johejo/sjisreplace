// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on https://github.com/golang/text/blob/master/encoding/japanese/shiftjis.go

// Package sjisreplace provides a encoder to safely convert to Shift-JIS.
package sjisreplace

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// NewEncoder returns a new transformer that is very similar to the Shift-JIS encoder in golang.org/x/text/encoding/japanese
// The Shift-JIS encoder in golang.org/x/text/encoding/japanese returns an error if it finds a rune that cannot be converted to Shift-JIS.
// This encoder does not return an error in the same case, replaces the target rune with a pre-specified rune, and continues processing.
func NewEncoder(r rune) transform.Transformer {
	return encoder{r: r}
}

type encoder struct {
	transform.NopResetter
	r rune
}

// Transform imprements transform.Transformer.
func (e encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	r, size := rune(0), 0
loop:
	for ; nSrc < len(src); nSrc += size {
		r = rune(src[nSrc])

		// Decode a 1-byte rune.
		if r < utf8.RuneSelf {
			size = 1

		} else {
			// Decode a multi-byte rune.
			r, size = utf8.DecodeRune(src[nSrc:])
			if size == 1 {
				// All valid runes of size 1 (those below utf8.RuneSelf) were
				// handled above. We have invalid UTF-8 or we haven't seen the
				// full character yet.
				if !atEOF && !utf8.FullRune(src[nSrc:]) {
					err = transform.ErrShortSrc
					break loop
				}
			}

			// func init checks that the switch covers all tables.
			switch {
			case encode0Low <= r && r < encode0High:
				if r = rune(encode0[r-encode0Low]); r>>tableShift == jis0208 {
					goto write2
				}
			case encode1Low <= r && r < encode1High:
				if r = rune(encode1[r-encode1Low]); r>>tableShift == jis0208 {
					goto write2
				}
			case encode2Low <= r && r < encode2High:
				if r = rune(encode2[r-encode2Low]); r>>tableShift == jis0208 {
					goto write2
				}
			case encode3Low <= r && r < encode3High:
				if r = rune(encode3[r-encode3Low]); r>>tableShift == jis0208 {
					goto write2
				}
			case encode4Low <= r && r < encode4High:
				if r = rune(encode4[r-encode4Low]); r>>tableShift == jis0208 {
					goto write2
				}
			case encode5Low <= r && r < encode5High:
				if 0xff61 <= r && r < 0xffa0 {
					r -= 0xff61 - 0xa1
					goto write1
				}
				if r = rune(encode5[r-encode5Low]); r>>tableShift == jis0208 {
					goto write2
				}
			}

			r = e.r
			goto write1
		}
	write1:
		if nDst >= len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst] = uint8(r)
		nDst++
		continue

	write2:
		j1 := uint8(r>>codeShift) & codeMask
		j2 := uint8(r) & codeMask
		if nDst+2 > len(dst) {
			err = transform.ErrShortDst
			break loop
		}
		if j1 <= 61 {
			dst[nDst+0] = 129 + j1/2
		} else {
			dst[nDst+0] = 193 + j1/2
		}
		if j1&1 == 0 {
			dst[nDst+1] = j2 + j2/63 + 64
		} else {
			dst[nDst+1] = j2 + 159
		}
		nDst += 2
		continue
	}
	return nDst, nSrc, err
}
