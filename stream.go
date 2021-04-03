package stream

type BaseStreamer interface {
	evaluateParallel()
	parallel()
}

// Streamer 数据流接口
type Streamer interface {
	/*
		中间态操作，返回Streamer
		垂直操作：map，filter，
		水平操作：distinct，sorted，limit，skip， sourceStage 会改变
	*/
	Map(function Function) Streamer
	Filter(predicate Predicate) Streamer
	//FindFirst(predicate Predicate) Streamer
	FlatMap(beStreamer BeStreamer) Streamer

	Distinct(comparator Comparator) Streamer
	Sorted(comparator Comparator) Streamer
	Limit(limit int) Streamer
	Skip(index int) Streamer

	/*
		终止操作，触发流水线真正的执行动作，无返回值。
	*/
	Reduce(reduceFun ReduceFun) interface{}
	ForEach(consumer Consumer)
	ToSlice(targetSlice interface{})
	Count() int
}

// PredicateOp 断言
type Predicate func(v interface{}) bool

// Consumer 消费
type Consumer func(v interface{})

// Comparator 比较
type Comparator func(i, j interface{}) bool

// ReduceFun 归约
type ReduceFun func(i, j interface{}) interface{}

// Function 普通方法
type Function func(v interface{}) interface{}

// Stream 构造流方法
type BeStreamer func(v interface{}) Streamer
