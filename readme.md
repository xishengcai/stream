## 前言
golang copy java 轮子(java.util.stream)。
Stream是java SE8 API添加的用于增加集合的操作接口，可以让你以一种声明的方式处理集合数据。将要处理的集合看作一种流的创建者，将集合内部的元素转换为流并且在管道中传输，并且可以在管道的节点上处理，比如筛选，排序，聚合等。元素流在管道内经过中间操作（intermediate operation）的处理，最后由终端操作（terminal operation）得到前面处理的结果。

sample example of stream
```
import java.util.stream.Stream;

public class Test {

    public static void main(String[] args) {
        Stream.of(1, 2, 3, 4, 5, 6, 7, 8, 9)
                .map(item -> item * 2)
                .forEach(item -> System.out.print(item + " "));

        System.out.println();

        Stream.of(1, 2, 5, 3, 4, 5, 6, 7, 8, 9)
                .distinct()
                .parallel()
                .map(item -> item * 2)
                .forEach(item -> System.out.print(item + " "));

        System.out.println();

        Stream.of(1, 2, 3, 4, 5, 6, 7, 8, 9)
                .parallel()
                .map(item -> item * 2)
                .forEachOrdered(item -> System.out.print(item + " "));
    }
}
```

result
```
sandbox> exited with status 0
2 4 6 8 10 12 14 16 18 
8 12 16 18 14 10 6 2 4 
2 4 6 8 10 12 14 16 18 
sandbox> exited with status 0
```

## java stream 源码解析
接口总浏


接口和实体类介绍


**stream pipeline**

Stream的执行过程被抽象出一个pipeline的概念，每个pipeline会包含不同的阶段（stage） 

- 起始阶段（source stage），有且只有一个，Stream创建的时候即被创建，比如：通过Stream.of接口创建时，会实例化ReferencePipeline.Head作为Stream的起始阶段，

- 过程阶段(intermediate stage)，0个或多个，如下例中包含两个过程阶段：distinct、map，parallel是一个特殊的存在，它只是设置一个标志位， 并不是一个过程阶段。对于过程阶段的各种操作，又有无状态操作(StatelessOp)和有状态操作(StatefulOp)之分, 比如：对于distinct、dropWhile、sorted需要在执行过程种记录执行状态，即有状态操作，而map、filter则属于无状态操作; 

- 终结阶段(terminal stage)，有且仅有一个，用于结果计算或者执行一些副作用，如下例中的forEach

```
Stream.of(1, 2, 5, 3, 4, 5, 6, 7, 8, 9)
                .distinct()
                .parallel()
                .map(item -> item * 2)
                .forEach(item -> System.out.print(item + " "));
```

上例中，最终构造的pipeline如图所示，pipeline数据结构是一个双向链表，每个节点分别存储上一阶段，下一阶段，及起始阶段。终端操作前均为lazy操作，所有操作并未真正执行。而终端操作会综合所有阶段执行计算。
![stream_pipeline](./image/stream_pipeline.png)

**思考**
1. 终结操作前均为lazy操作，所有操作并未真正执行
    
2. 数据是如何传递的
```
    sourceArray[0] -> headStage(sourceElement) => result
                        -> secondStage(result) => result
                           -> thirdStage(result) => result
                            -> statefulSet(result), 该阶段的do 方法内不会有 nextStage.do() ,因为该阶段会变成sourceStage


    sourceArray[1] -> headStage(sourceElement) => result
                        -> secondStage(result) => result
                           -> thirdStage(result) => result

    sourceArray[2] -> headStage(sourceElement) => result
                            -> secondStage(result) => result
                            -> thirdStage(result) => result
```

3. 为什么需要有状态操作

有状态操作是为了解决某些阶段需要全量数据才能处理，所以在类似于sort 阶段需要将上一步的数据都存到 sort 阶段的temp中，
然后在使用sort 逻辑处理所有temp数据， 再保存到 data中。

4. 为什么要引入操作标识位？

5. nextStage.do(nextStage.do, dateElement) ？

一个元素会被每个阶段依次处理，直到遇到状态操作。然后将有状态操作阶段设置为sourceStage。再继续重复，直到终止操作



**Stream 源码操作**

**Stream的创建**
Stream、IntStream、LongStream、DoubleStream接口都提供了静态方法of，用于便捷地创建Stream，分别用于创建引用类型、int、long、double的Stream。以Stream接口为例，传入一系列有序元素，比如Stream.of(1, 2, 3, 4, 5, 6, 7, 8, 9)。如下所示，Stream.of方法通过调用Arrays.stream实现。所以Stream.of(1, 2, 3, 4, 5, 6, 7, 8, 9)等同于Arrays.stream(new int[]{1, 2, 3, 4, 5, 6, 7, 8, 9})

```
//Stream.java
public static<T> Stream<T> of(T... values) {
        return Arrays.stream(values);
    }
//Arrays.java
public static <T> Stream<T> stream(T[] array) {
        return stream(array, 0, array.length);
    }

public static <T> Stream<T> stream(T[] array, int startInclusive, int endExclusive) {
        return StreamSupport.stream(spliterator(array, startInclusive, endExclusive), false);
    }
```

StreamSupport是一个工具类用于创建顺序或并行Stream。StreamSupport.stream需要两个参数： - spliterator，Spliterator类型，可通过Spliterators.spliterator创建，用于遍历/拆分数组， - parallel，boolean类型，标示是否并行，默认为false。

ReferencePipeline实现了stream接口，是实现Stream过程阶段或起始阶段的抽象基类。ReferencePipeline.Head是原始阶段的实现。

```
//StreamSupport.java
public static <T> Stream<T> stream(Spliterator<T> spliterator, boolean parallel) {
        Objects.requireNonNull(spliterator);
        return new ReferencePipeline.Head<>(spliterator,
                                            StreamOpFlag.fromCharacteristics(spliterator),
                                            parallel);
    }
```

至此，Stream创建的简单流程就完成了，IntStream、LongStream、DoubleStream的创建也是类似的。



## 代码阅读
### file：BaseStream.java
```java
public interface BaseStream<T, S extends BaseStream<T, S>>
        extends AutoCloseable {
    
    Iterator<T> iterator();
    Spliterator<T> spliterator();
    boolean isParallel();
    S sequential();
    S parallel();
    S unordered();
    S onClose(Runnable closeHandler);
}
```

### file: collector.java
```java

public interface Collector<T, A, R> {

    Supplier<A> supplier();
    BiConsumer<A, T> accumulator();
    BinaryOperator<A> combiner();
    Function<A, R> finisher();
    Set<Characteristics> characteristics();
    public static<T, R> Collector<T, R, R> of(Supplier<R> supplier,
                                              BiConsumer<R, T> accumulator,
                                              BinaryOperator<R> combiner,
                                              Characteristics... characteristics) {
        Objects.requireNonNull(supplier);
        Objects.requireNonNull(accumulator);
        Objects.requireNonNull(combiner);
        Objects.requireNonNull(characteristics);
        Set<Characteristics> cs = (characteristics.length == 0)
                                  ? Collectors.CH_ID
                                  : Collections.unmodifiableSet(EnumSet.of(Collector.Characteristics.IDENTITY_FINISH,
                                                                           characteristics));
        return new Collectors.CollectorImpl<>(supplier, accumulator, combiner, cs);
    }
    public static<T, A, R> Collector<T, A, R> of(Supplier<A> supplier,
                                                 BiConsumer<A, T> accumulator,
                                                 BinaryOperator<A> combiner,
                                                 Function<A, R> finisher,
                                                 Characteristics... characteristics) {
        Objects.requireNonNull(supplier);
        Objects.requireNonNull(accumulator);
        Objects.requireNonNull(combiner);
        Objects.requireNonNull(finisher);
        Objects.requireNonNull(characteristics);
        Set<Characteristics> cs = Collectors.CH_NOID;
        if (characteristics.length > 0) {
            cs = EnumSet.noneOf(Characteristics.class);
            Collections.addAll(cs, characteristics);
            cs = Collections.unmodifiableSet(cs);
        }
        return new Collectors.CollectorImpl<>(supplier, accumulator, combiner, finisher, cs);
    }

    enum Characteristics {
        CONCURRENT,
        UNORDERED,
        IDENTITY_FINISH
    }
}
```



## 构思
1. 定义接口
2. 定义实体类
    定义 属性

3. 实现逻辑
    链表数据结构，terminal operation 向上执行

    source
        next.do
            next.do
    execute（stream)
    nextStage.do(nextStage, func(v))
        nextStage.do(nextStage, func(v))

    Map，filter

    ToSlice


link
[小蓝鲸](https://club.perfma.com/article/116123)
