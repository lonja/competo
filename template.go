package competo

import (
	"errors"
	"io/ioutil"
	"sort"
	"strings"
)

var (
	errEmptyFilePath = errors.New("file path is empty")
)

type Template struct {
	Strings []*string   // Addresses of unique strings in a template. Pointers used to reduce memory of the app.
	Bytes   [][]byte    // Bytes presentation of unique strings in a template
	mapping map[int]int // mapping from original template strings to unique string addresses
	len     int
}

func (t *Template) Len() int {
	return t.len
}

func (t *Template) GetString(i int) *string {
	return t.Strings[t.mapping[i]]
}

func ReadTemplateFromFile(filePath string) (*Template, error) {
	if filePath == "" {
		return nil, errEmptyFilePath
	}

	rawTmpl, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	splittedTmpl := strings.Split(string(rawTmpl), "\n")

	splittedTmpl = deleteEmptyRows(splittedTmpl)

	us := uniqueStrings(splittedTmpl)

	lenus := len(us)

	result := make([]*string, lenus)
	bytes := make([][]byte, lenus)
	for i := 0; i < lenus; i++ {
		result[i] = &us[i]
		bytes[i] = []byte(us[i])
	}

	lm := make(map[int]int)
	for i, s := range splittedTmpl {
		for j, expr := range result {
			if *expr != s {
				continue
			}

			lm[i] = j
		}
	}

	return &Template{Strings: result, Bytes: bytes, mapping: lm, len: len(splittedTmpl)}, nil
}

func deleteEmptyRows(tmpl []string) []string {
	for i := 0; i < len(tmpl); i++ {
		if len(tmpl[i]) != 0 {
			continue
		}

		// delete empty row
		copy(tmpl[:i], tmpl[i+1:])
		tmpl[len(tmpl)-1] = ""
		tmpl = tmpl[:len(tmpl)-1]
	}

	return tmpl
}

func uniqueStrings(s []string) []string {
	tmpResult := make([]string, len(s))
	copy(tmpResult, s)

	sort.Strings(tmpResult)
	j := 0
	for i := 1; i < len(tmpResult); i++ {
		if tmpResult[j] == tmpResult[i] {
			continue
		}
		j++
		tmpResult[j] = tmpResult[i]
	}
	return tmpResult[:j+1]
}
