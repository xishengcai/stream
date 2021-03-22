package stream

import "sort"

var _ sort.Interface = (*sortInterface)(nil)

type sortInterface struct {
	data []interface{}
	comparator  Comparator
}

func (s sortInterface) Len() int            { return len(s.data) }
func (s *sortInterface) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s *sortInterface) Less(i, j int) bool { return s.comparator(s.data[i], s.data[j]) }

