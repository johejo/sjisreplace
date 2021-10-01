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
