package competo

import (
	"bufio"
	"bytes"
	"io"
)

type Search struct {
	index []map[*string][]int
	r     io.ReadCloser
	tmpl  *Template
}

func NewSearch(r io.ReadCloser) *Search {
	return &Search{
		index: make([]map[*string][]int, 0, 0),
		r:     r,
	}
}

func (s *Search) buildIndex(tmpl *Template) {
	scanner := bufio.NewScanner(s.r)
	li := 0

	for scanner.Scan() {
		str := scanner.Bytes()
		occ := make(map[*string][]int)
		for ei, expr := range tmpl.Strings {
			indexes := findAllIndexes(str, tmpl.Bytes[ei])

			if indexes == nil {
				continue
			}
			occ[expr] = indexes
		}
		s.index = append(s.index, occ)
		li++
	}

	_ = s.r.Close()
}

func (s *Search) Count(tmpl *Template) int {
	s.buildIndex(tmpl)

	count := 0

	for ln := 0; ln < len(s.index); ln++ {
		for tmplStr, coordinates := range s.index[ln] {
			if tmplStr != tmpl.GetString(0) {
				continue
			}

			for ci := 0; ci < len(coordinates); ci++ {
				matched := 1
				ln := ln + 1
				for i := 0; i < tmpl.Len()-1; i++ {
					if ln >= len(s.index) {
						break
					}

					nextOccurences := s.index[ln]

					nextCoordinates, ok := nextOccurences[tmpl.GetString(matched)]
					if !ok {
						break
					}

					if len(nextCoordinates) == 1 && coordinates[ci] != nextCoordinates[0] {
						break
					}

					for nci := 0; nci < len(nextCoordinates); nci++ {
						if !(coordinates[ci] == nextCoordinates[nci]) {
							continue
						}
						matched++
						ln++
					}
				}

				if matched == tmpl.Len() {
					count++
				}
			}
		}
	}

	return count
}

// По сравнению с регулярками нужно на 40% меньше памяти
func findAllIndexes(s, substr []byte) []int {
	// Индекс предыдущего элемента в общей строке
	index := bytes.Index(s, substr)
	if index == -1 {
		return nil
	}

	lens := len(s)
	lensubstr := len(substr)
	result := make([]int, 0, lens/(2*lensubstr)+1)

	result = append(result, index)

	// i - сдвиг относительно начала исходной строки
	// index - индекс текущего элемента
	// lastIndex - индекс предыдущего элемента
	prevIndex := index
	for i := prevIndex + lensubstr; i < lens; i = prevIndex + lensubstr {
		index = bytes.Index(s[i:], substr)
		if index == -1 {
			break
		}
		index += i
		result = append(result, index)
		prevIndex = index
	}

	return result
}
