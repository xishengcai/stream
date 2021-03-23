package stream

import (
	"reflect"
	"sort"
	"sync"
)

type handleData func(v []interface{}) []interface{}

var _ Streamer = &pipeline{}

// pipeline Streamer的子类
type pipeline struct {
	lock sync.Mutex

	previousStage *pipeline
	sourceStage   *pipeline
	nextStage     *pipeline

	data, temp []interface{}

	do func(nextStage *pipeline, v interface{})

	parallel, entered, stop bool
	depth                   int
}

func New(arr interface{}, parallel bool) Streamer {
	data := make([]interface{}, 0)
	dataValue := reflect.ValueOf(&data).Elem()

	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Kind() == reflect.Slice || arrValue.Kind() == reflect.Array {
		for i := 0; i < arrValue.Len(); i++ {
			dataValue.Set(reflect.Append(dataValue, arrValue.Index(i)))
		}
	}

	p := &pipeline{
		data:     data,
		parallel: parallel,
		depth:    0,
	}
	p.sourceStage = p
	return p
}

func (p *pipeline) Map(fun Function) Streamer {
	nextStage := &pipeline{
		previousStage: p,
		do: func(nextStage *pipeline, v interface{}) {
			nextStage.do(nextStage.nextStage, fun(v))
		},
		sourceStage: p.sourceStage,
		depth:       p.depth + 1,
	}
	p.nextStage = nextStage
	return nextStage

}

func (p *pipeline) Filter(fun Predicate) Streamer {
	nextStage := &pipeline{
		previousStage: p,
		do: func(nextStage *pipeline, v interface{}) {
			if fun(v) {
				nextStage.do(nextStage.nextStage, v)
			}
		},
		sourceStage: p.sourceStage,
		depth:       p.depth + 1,
	}
	p.nextStage = nextStage

	return nextStage
}

func (p *pipeline) Distinct(comparator Comparator) Streamer {
	handle := func(v []interface{}) []interface{} {
		return removeDuplicate(v, comparator)
	}
	return statefulSetOp(p, handle)
}

func (p *pipeline) FlatMap(fun Function) Streamer {
	nextStage := &pipeline{
		previousStage: p,
		do: func(nextStage *pipeline, v interface{}) {
			nextStage.do(nextStage.nextStage, fun(v))
		},
		sourceStage: p.sourceStage,
		depth:       p.depth + 1,
	}
	p.nextStage = nextStage

	return nextStage
}

func (p *pipeline) Sorted(comparator Comparator) Streamer {
	handle := func(v []interface{}) []interface{} {
		s := &sortInterface{data: v, comparator: comparator}
		sort.Sort(s)
		return v
	}
	return statefulSetOp(p, handle)
}

// limit 类似于SQL语句中的Limit
// 水平操作，需要拿到所有数据，也叫有状态操作
func (p *pipeline) Limit(limit int) Streamer {
	handle := func(v []interface{}) []interface{} {
		if limit > len(v) {
			return nil
		}
		return v[:limit]
	}
	return statefulSetOp(p, handle)
}

// Skip 类似于sql语句中的indexOf， 与limit 配合可实现分页操作
func (p *pipeline) Skip(index int) Streamer {

	handle := func(v []interface{}) []interface{} {
		if index > len(v)-1 {
			return nil
		}
		return v[index:]
	}
	return statefulSetOp(p, handle)
}

func (p *pipeline) Reduce(reduceFun ReduceFun) interface{} {
	handle := func(v []interface{}) []interface{} {
		return v
	}
	currentStage := statefulSetOp(p, handle)

	if len(currentStage.data) == 0 {
		return nil
	} else if len(currentStage.data) == 1 {
		return currentStage.data
	}
	result := currentStage.data[0]
	for _, item := range currentStage.data[1:] {
		result = reduceFun(result, item)
	}
	return result

}

func (p *pipeline) ForEach(consumer Consumer) {

	nextStage := &pipeline{
		previousStage: p,
		do: func(nextStage *pipeline, v interface{}) {
			consumer(v)
		},
		sourceStage: p.sourceStage,
		depth:       p.depth + 1,
	}
	p.nextStage = nextStage

	terminal{}.evaluate(p.sourceStage)

}

func (p *pipeline) ToSlice(target interface{}) {

	targetValue := reflect.ValueOf(&target)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}

	arrValue := reflect.ValueOf(p.data)
	if arrValue.Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Kind() == reflect.Slice || arrValue.Kind() == reflect.Array {
		for i := 0; i < arrValue.Len(); i++ {
			targetValue.Set(reflect.Append(targetValue, arrValue.Index(i)))
		}
	}
}

// statefulSetOp 装饰器，
// before： evaluate
// after：修改sourceStage && p.previousStage.nextStage
func statefulSetOp(p *pipeline, handle handleData) *pipeline {
	nextStage := &pipeline{
		previousStage: p,
		do: func(nextStage *pipeline, v interface{}) {
			p.temp = append(p.temp, v)
		},
		depth: p.depth + 1,
	}
	p.nextStage = nextStage
	terminal{}.evaluate(p.sourceStage)
	nextStage.data = handle(p.temp)
	nextStage.sourceStage = nextStage
	return nextStage
}

// removeDuplicate 数组元素去重
func removeDuplicate(arr []interface{}, comparator Comparator) []interface{} {
	if arr == nil {
		return nil
	}

	result := make([]interface{}, 0)
	for i := range arr {
		flag := true
		for j := range result {
			if comparator(arr[i], result[j]) {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, arr[i])
		}
	}
	return result

}
