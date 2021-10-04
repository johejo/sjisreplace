package sjisreplace_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/johejo/sjisreplace"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Test(t *testing.T) {
	cases := []struct {
		s       string
		replace bool
	}{
		{s: "a", replace: false},
		{s: "あ", replace: false},
		{s: "ア", replace: false},
		{s: "阿", replace: false},
		{s: "😋", replace: true},
	}
	for _, cc := range cases {
		t.Run(cc.s, func(t *testing.T) {
			c := cc.s
			for i := 0; i < 10; i++ {
				n := int(math.Pow(4.0, float64(i)) + 1)
				t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
					long := strings.Repeat(c, n)
					b := new(bytes.Buffer)
					bw := bufio.NewWriter(b)
					w := transform.NewWriter(bw, sjisreplace.NewEncoder('?'))
					for _, e := range []byte(long) {
						if _, err := w.Write([]byte{e}); err != nil {
							t.Fatal(err)
						}
					}
					if err := bw.Flush(); err != nil {
						t.Fatal(err)
					}
					r := transform.NewReader(b, japanese.ShiftJIS.NewDecoder())
					_got, err := io.ReadAll(r)
					if err != nil {
						t.Fatal(err)
					}
					got := string(_got)
					var count int
					if cc.replace {
						count = strings.Count(got, "?")
					} else {
						count = strings.Count(got, c)
					}
					if count != n {
						t.Errorf("The number of chars is different before and after transformation. want=%d, got=%d", n, count)
					}
				})
			}
		})
	}
}
