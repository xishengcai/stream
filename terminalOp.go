package stream

import "sync"

// terminalOperator 串、并行接口
type terminalOperator interface {
	sequential(streamer *Streamer)
	evaluateParallel(streamer *Streamer)
}

type terminal struct{}

func (t terminal) evaluate(streamer Streamer) {
	sourceStreamer := streamer.(*pipeline)
	if sourceStreamer.nextStage == nil {
		return
	}
	if sourceStreamer.parallel {
		t.evaluateParallel(sourceStreamer)
	} else {
		t.evaluateSequential(sourceStreamer)
	}
}

func (t *terminal) evaluateSequential(sourceStreamer *pipeline) {
	headStage := sourceStreamer.nextStage
	for _, item := range sourceStreamer.data {
		// 能够向链条一样执行到terminal operation
		// headStage -> secondStage -> thirdStage
		headStage.do(headStage.nextStage, item)
	}
}

func (t *terminal) evaluateParallel(sourceStreamer *pipeline) {
	wait := sync.WaitGroup{}
	wait.Add(len(sourceStreamer.data))
	headStage := sourceStreamer.nextStage

	for _, item := range sourceStreamer.data {
		go func(item interface{}) {
			defer wait.Done()
			headStage.do(headStage.nextStage, item)
		}(item)
	}

	wait.Wait()
}
