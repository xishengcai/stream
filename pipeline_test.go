package stream

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

type student struct {
	id     int
	name   string
	age    int
	scores []int
}

func (s student) String() string {
	return fmt.Sprintf("{id:%d, name:%s, age:%d,scores:%v}", s.id, s.name, s.age, s.scores)
}

func createStudents() []student {
	names := []string{"Tom-0", "Kate-1", "Lucy-2", "Jim-3", "Jack-4", "King-5", "Lee-6", "Mask-7","Mask-8","Mask-9"}
	students := make([]student, 10)
	rnd := func(start, end int) int { return rand.Intn(end-start) + start }
	for i := 0; i <= 9; i++ {
		students[i] = student{
			id:     i,
			name:   names[i],
			age:    rnd(0, 100),
			scores: []int{rnd(60, 100), rnd(60, 100), rnd(60, 100)},
		}
	}
	return students
}

//https://blog.golang.org/laws-of-reflection
func TestReflect(tt *testing.T) {
	type T struct {
		A int
		B string
	}

	t := T{A: 1, B: "hello"}
	e := reflect.ValueOf(&t).Elem()
	typeOfT := e.Type()
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		tt.Logf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	e.Field(0).SetInt(77)
	e.Field(1).SetString("Sunset Strip")
	tt.Log("t is now", t)

	reflect.ValueOf(&t).Elem()

}

func TestMap(t *testing.T) {
	students := createStudents()

	stream := New(students, false)
	stream.
		Map(func(v interface{}) interface{} {
			s := v.(student)
			s.age += 100
			return s
		}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}


func TestFilter(t *testing.T) {
	students := createStudents()
	stream := New(students, false)
	stream.
		Filter(func(v interface{}) bool {
			return v.(student).age > 50
		}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestLimit(t *testing.T) {
	students := createStudents()
	stream := New(students, false)
	stream.Limit(5).ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestSkip(t *testing.T) {
	students := createStudents()
	stream := New(students, false)
	stream.Skip(3).ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestPage(t *testing.T){
	students := createStudents()
	stream := New(students, false)

	// except last 5 item
	stream.
		Skip(5).
		Limit(5).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestDistinct(t *testing.T){
	students := createStudents()
	students[3], students[5] = students[0],students[0]
	New(students, false).
		Map(func(v interface{}) interface{} {
			s := v.(student)
			s.age += 100
			return s
		}).
		Distinct().
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestSorted(t *testing.T){
	students := createStudents()
	New(students, false).
		Map(func(v interface{}) interface{} {
			s := v.(student)
			s.age += 100
			return s
		}).
		Sorted(func (s1 interface{}, s2 interface{})bool{
			return s1.(student).age > s2.(student).age}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}
