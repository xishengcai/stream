## 前言
golang copy java 轮子(java.util.stream)。
Stream是java SE8 API添加的用于增加集合的操作接口，可以让你以一种声明的方式处理集合数据。将要处理的集合看作一种流的创建者，将集合内部的元素转换为流并且在管道中传输，并且可以在管道的节点上处理，比如筛选，排序，聚合等。元素流在管道内经过中间操作（intermediate operation）的处理，最后由终端操作（terminal operation）得到前面处理的结果。

## wath
stream表示包含着一系列元素的集合，我们可以对其做不同类型的操作，用来对这些元素执行计算。使用流操作，可以使代码更简洁，减少了大量的for循环代码。

## 概念
- 中间操作，对数组元素操作完毕后，将结果传递给下一个操作。

- 终止操作， 是对流操作的一个结束动作。

- 有状态，有数据存储功能，线程不安全

- 无状态，不保存数据。线程安全

## function
> 中间态 无状态操作

- [x] Map
- [x] Filter
- [ ] FlatMap
- [ ] peek

> 中间态 有状态操作

- [x] Distinct
- [x] Sorted
- [x] Limit
- [x] Skip

> 终止态操作

- [x] Reduce
- [x] ForEach
- [x] ToSlice
- [x] count
- [] max
- [] min
- [] forEachOrdered
- [] anyMatch
- [] allMatch
- [] noneMatch
- [] findFirst
- [] findAny

## 待改进
- reduce 并发计算，是否有
- distinc 去重计算

## example
