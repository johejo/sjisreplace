# sjisreplace

[![ci](https://github.com/johejo/sjisreplace/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/johejo/sjisreplace/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/johejo/sjisreplace.svg)](https://pkg.go.dev/github.com/johejo/sjisreplace)
[![codecov](https://codecov.io/gh/johejo/sjisreplace/branch/main/graph/badge.svg)](https://codecov.io/gh/johejo/sjisreplace)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/sjisreplace)](https://goreportcard.com/report/github.com/johejo/sjisreplace)

Package sjisreplace provides a encoder to safely convert to Shift-JIS.

The Shift-JIS encoder in golang.org/x/text/encoding/japanese returns an error if it finds a rune that cannot be converted to Shift-JIS.

This encoder does not return an error in the same case, replaces the target rune with a pre-specified rune, and continues processing.

## Example

```go
package sjisreplace_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/johejo/sjisreplace"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Example() {
	const (
		emoji    = "😋"
		hiragana = "あ"
		katakana = "ア"
		kanji    = "亜"
	)
	cases := []string{"a", emoji, hiragana, katakana, "1", "%", kanji}

	b := new(bytes.Buffer)
	w := transform.NewWriter(b, sjisreplace.NewEncoder('?'))
	for _, s := range cases {
		if _, err := w.Write([]byte(s)); err != nil {
			panic(err)
		}
	}
	if _, err := w.Write([]byte(strings.Join(cases, ""))); err != nil {
		panic(err)
	}

	r := transform.NewReader(b, japanese.ShiftJIS.NewDecoder())
	got, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(got)) // output: a?あア1%亜a?あア1%亜
}
```
