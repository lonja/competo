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

// Нужно связать строки и регулярки таким образом, чтобы я мог получить строки по регулярке
type Template struct {
	Strings []*string
	Bytes   [][]byte
	mapping map[int]int // map string -> regexp
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

	strTmpl := string(rawTmpl)
	splitedTmpl := strings.Split(strTmpl, "\n")

	splitedTmpl = deleteEmptyRows(splitedTmpl)

	us := uniqueStrings(splitedTmpl)

	lenus := len(us)

	result := make([]*string, lenus)
	bytes := make([][]byte, lenus)
	for i := 0; i < lenus; i++ {
		result[i] = &us[i]
		bytes[i] = []byte(us[i])
	}

	lm := make(map[int]int)
	for i, s := range splitedTmpl {
		for j, expr := range result {
			if *expr != s {
				continue
			}

			lm[i] = j
		}
	}

	return &Template{Strings: result, Bytes: bytes, mapping: lm, len: len(splitedTmpl)}, nil
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
