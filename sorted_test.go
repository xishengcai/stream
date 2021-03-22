package stream

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {

	students := []interface{}{
		student{id: 1, name: "a1", age: 10},
		student{id: 2, name: "a2", age: 11},
		student{id: 3, name: "a3", age: 110},
		student{id: 4, name: "a4", age: 13},
		student{id: 5, name: "a5", age: 122},
		student{id: 6, name: "a6", age: 1},
	}
	c := &sortInterface{data: students, comparator: func(i interface{}, j interface{}) bool {
		return i.(student).age > j.(student).age
	}}
	sort.Sort(c)

	t.Logf("IS Sorted? %v", sort.IsSorted(c))

	t.Log(c.data)

}

func TestSortSlice(t *testing.T) {
	students := []student{
		{id: 1, name: "a1", age: 10},
		{id: 2, name: "a2", age: 11},
		{id: 3, name: "a3", age: 110},
		{id: 4, name: "a4", age: 13},
		{id: 5, name: "a5", age: 122},
		{id: 6, name: "a6", age: 1},
	}

	sort.Slice(students, func(i, j int) bool {
		return students[i].age > students[j].age
	})

	for _, s := range students {
		t.Log(s)
	}
}
