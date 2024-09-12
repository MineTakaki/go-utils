package ioutilx

import (
	"bufio"
	"encoding/csv"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// ReadTextLines テキストを行単位に分割して取得します
func ReadTextLines(rd io.Reader, noSkip bool, fn func(string) string) (lines []string, err error) {
	r := bufio.NewReader(rd)
	for {
		var line string
		line, err = r.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				err = errors.WithStack(err)
				return
			}
			err = nil
			break
		}

		if fn != nil {
			line = fn(line)
		}
		if line == "" && !noSkip {
			continue
		}
		lines = append(lines, line)
	}
	return
}

// ReadCSVTextLines CSVテキストから簡易的に全体を読み込みます
func ReadCSVTextLines(rd io.Reader, sep, comment rune) (lines [][]string, err error) {
	r := csv.NewReader(rd)
	r.Comma = sep
	r.Comment = comment
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true

	commentStr := string(comment)

	for {
		var cols []string
		cols, err = r.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				err = errors.WithStack(err)
				return
			}
			err = nil
			break
		}

		n := len(cols)
		if n == 0 {
			continue
		}

		for i, s := range cols {
			if strings.HasPrefix(s, commentStr) {
				cols = cols[:i]
				break
			}
		}

		if len(cols) == 0 {
			continue
		}

		lines = append(lines, cols)
	}

	return
}
