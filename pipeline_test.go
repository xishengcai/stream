package stream

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
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

func createStudents(num int) []student {
	names := []string{"Tom", "Kate", "Lucy", "Jim", "Jack", "King", "Lee", "Mask", "Mask", "Mask"}
	students := make([]student, num)
	rnd := func(start, end int) int { return rand.Intn(end-start) + start }
	for i := 0; i < num; i++ {
		students[i] = student{
			id:     i,
			name:   names[rnd(0, 9)],
			age:    rnd(0, 100),
			scores: []int{rnd(0, 100), rnd(0, 100), rnd(0, 100)},
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
	students := createStudents(10)

	stream := New(students, false)
	stream.
		Map(func(v interface{}) interface{} {
			s := v.(student)
			s.age += 100
			return s
		}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestFlatMap(t *testing.T) {
	stream := New([]string{"hello", "world"}, false).
		Map(func(v interface{}) interface{} {
			return strings.Split(v.(string), "")
		}).FlatMap(func(v interface{}) Streamer {
		return New(v, false)
	}).Distinct(func(s1 interface{}, s2 interface{}) bool {
		return reflect.DeepEqual(s1, s2)
	})
	stream.ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestFilter(t *testing.T) {
	students := createStudents(10)
	stream := New(students, false)
	stream.
		Filter(func(v interface{}) bool {
			return v.(student).age > 50
		}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestLimit(t *testing.T) {
	students := createStudents(10)
	stream := New(students, false)
	stream.Limit(5).ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestSkip(t *testing.T) {
	students := createStudents(10)
	stream := New(students, false)
	stream.Skip(3).ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestPage(t *testing.T) {
	students := createStudents(10)
	stream := New(students, false)

	// except last 5 item
	stream.
		Skip(5).
		Limit(5).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestDistinct(t *testing.T) {
	students := createStudents(10)
	students[3], students[5], students[9] = students[0], students[0], students[0]
	New(students, false).
		Map(func(v interface{}) interface{} {
			s := v.(student)
			s.age += 100
			return s
		}).
		Distinct(func(s1 interface{}, s2 interface{}) bool {
			return reflect.DeepEqual(s1, s2)
		}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestSorted(t *testing.T) {
	students := createStudents(10)
	New(students, false).
		Map(func(v interface{}) interface{} {
			s := v.(student)
			s.age += 100
			return s
		}).
		Sorted(func(s1 interface{}, s2 interface{}) bool {
			return s1.(student).age > s2.(student).age
		}).
		ForEach(func(v interface{}) { fmt.Println(v) })
}

func TestReduce(t *testing.T) {
	students := createStudents(10)
	sumAge := New(students, false).
		Map(func(v interface{}) interface{} {
			time.Sleep(time.Second * 1)
			return v
		}).
		Map(func(v interface{}) interface{} {
			return v.(student).age
		}).Reduce(func(i interface{}, j interface{}) interface{} {
		return i.(int) + j.(int)
	})

	t.Log(sumAge)
}

func TestReduceParallel(t *testing.T) {
	students := createStudents(100000)
	sumAge := New(students, true).
		Map(func(v interface{}) interface{} {
			time.Sleep(time.Second * 1)
			return v
		}).
		Map(func(v interface{}) interface{} {
			return v.(student).age
		}).Reduce(func(i interface{}, j interface{}) interface{} {
		return i.(int) + j.(int)
	})

	t.Log(sumAge)
}

func TestCount(t *testing.T) {
	students := createStudents(1000)

	count := New(students, false).
		Map(func(v interface{}) interface{} {
			return v.(student).age
		}).
		Skip(5000).
		Count()
	assert.Equal(t, count, 0)

	count = New(students, false).
		Map(func(v interface{}) interface{} {
			return v.(student).age
		}).
		Limit(100).
		Count()

	assert.Equal(t, count, 100)
}

func TestSlice(t *testing.T) {
	students := createStudents(10)
	sts := make([]student, 10)

	New(students, false).
		Skip(2).
		ToSlice(&sts)

	t.Log(sts)
}

func TestAnyMatch(t *testing.T) {
	students := createStudents(10)
	matched := New(students, false).AnyMatch(func(v interface{}) bool {
		return v.(student).age > 1
	})
	assert.Equal(t, matched, true)

	matched = New(students, false).AnyMatch(func(v interface{}) bool {
		return v.(student).age > 1000
	})

	assert.Equal(t, matched, false)
}
