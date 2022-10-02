// many of code copied from `time` package
// so keep it credited:

// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// modified parts are governed by a license described in LICENSE file.

package flextime

import (
	"github.com/ngicks/type-param-common/iterator"
)

const (
	tokenIdMask = uint(0b11111111) // 8 bit = 256 total
	// otherMask   = tokenIdMask << 8 // unused.
	lenMask = tokenIdMask << 16
	// otherOtherMask = tokenIdMask << 24 // unused.
)

type layoutToken uint

func (i layoutToken) T() layoutToken {
	return i & layoutToken(tokenIdMask)
}

func (i layoutToken) Is(v layoutToken) bool {
	return i.T()&v.T() == v.T()
}

func (i layoutToken) SetLen(l uint) layoutToken {
	return (i & layoutToken(^lenMask)) | layoutToken((tokenIdMask&l)<<16)
}

// Len is byte-length of token.
func (i layoutToken) Len() uint {
	return uint((i & layoutToken(lenMask)) >> 16)
}

func (i layoutToken) String() string {
	id := i.T()
	if id == GoFracSecond0 || id == GoFracSecond9 || id == IsoFracSecond {
		var start int
		var initial string
		if id != IsoFracSecond {
			start = 1
			initial = "."
		}

		var frac string
		if id == GoFracSecond0 {
			frac = "0"
		} else if id == GoFracSecond9 {
			frac = "9"
		} else {
			frac = "S"
		}

		return iterator.Fold[int, string](
			iterator.FromRange(start, int(i.Len())),
			func(accumulator string, _ int) string {
				return accumulator + frac
			},
			initial,
		)
	}
	for k, v := range goStrToNum {
		if v == id {
			return k
		}
	}
	for k, v := range isoStrToNum {
		if v == id {
			return k
		}
	}

	return "invalid"
}
