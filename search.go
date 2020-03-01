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

// buildIndex indexes whole file
// Index helps to find all occurences of template
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

// Count function counts quantity of occurences of template in the file
func (s *Search) Count(tmpl *Template) int {
	s.buildIndex(tmpl)

	count := 0

	indexlen := len(s.index)
	// Move by indexed file
	for line := 0; line < indexlen; line++ {
		// Loop by different tmpl matches in one line
		// Coordinates - list of indexes(start,end) our tmpl in line.
		for tmplStr, coords := range s.index[line] {
			// Checks that indexed line contains first line of template
			if tmplStr != tmpl.GetString(0) {
				continue
			}

			// Iterate by coords
			for ci := 0; ci < len(coords); ci++ {
				matched := 1

				// Move line-by-line to the file end to match whole template
				ln := line + 1
				for i := 0; i < tmpl.Len()-1; i++ {
					if ln >= len(s.index) {
						break
					}

					nextOccurs := s.index[ln]

					nextCoords, ok := nextOccurs[tmpl.GetString(matched)]
					if !ok {
						break
					}

					if len(nextCoords) == 1 && coords[ci] != nextCoords[0] {
						break
					}

					for nci := 0; nci < len(nextCoords); nci++ {
						if !(coords[ci] == nextCoords[nci]) {
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

// This function implemented to use instead of regexp.Regexp.FindAllStringIndex function.
// That reduced memory usage by 30%
func findAllIndexes(s, substr []byte) []int {
	// Index of current match substr in s
	idx := bytes.Index(s, substr)
	if idx == -1 {
		return nil
	}

	lens := len(s)
	lensubstr := len(substr)
	// Preallocate memory to decrease allocations count
	res := make([]int, 0, lens/(2*lensubstr)+1)

	res = append(res, idx)

	// i – shift from beginning of string
	// prevIdx – index of previous match substr in s
	prevIdx := idx
	for i := prevIdx + lensubstr; i < lens; i = prevIdx + lensubstr {
		idx = bytes.Index(s[i:], substr)
		if idx == -1 {
			break
		}
		idx += i
		res = append(res, idx)
		prevIdx = idx
	}

	return res
}
