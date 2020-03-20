# Engine

A library that makes it very easy to write concurrent code in golang.
`Engine` lets you spawn `n` number of workers that are ready to execute your code.

## Parallelism vs Concurrency

Concurrency means multiple tasks which start, run, and complete in overlapping time periods, in no specific order.
Parallelism is when multiple tasks OR several part of a unique task literally run at the same time, e.g. on a multi-core processor.
Remember that Concurrency and parallelism are NOT the same thing.
[Must watch Rob Pike's talk](https://www.youtube.com/watch?v=cN_DpYBzKso).

## Concurrency is Hard

As wonderfully summed up by [dm03514](https://medium.com/dm03514-tech-blog/golang-candidates-and-contexts-a-heuristic-approach-to-race-condition-detection-e2b230e70d08)...
```
"The Go language provides absolutely amazing concurrency primitives and truly achieves making concurrency a first class citizen. 
Unfortunately ensuring concurrent correctness requires the combination of many different techniques in order to minimize the chances of concurrency related errors."
```

## How Engine makes concurrent execution easy and robust?
Engine is designed to spawn `n` workers that each maintain a `goroutine`.
It then dispatches work that run concurrently on these available workers as simple as `engine.Do(...)`.
It handles queueing of execution and error handling of each execution specifically without affecting other executions.
 
 
Engine achieves this by formulating a pattern for any execution. Any execution must adhere to the following interface...
```GO
type Executable interface {
	Execute() error
	IsCompleted() bool
	OnSuccess()
	OnFailure(err error)
}
```

An example of an executable...
```GO

type testTask struct {
	ID     int
	Fail   bool
	Delay  string
	Status string
}

func (t *testTask) Execute() error {
	duration, err := time.ParseDuration(t.Delay)
	if err != nil {
		logger.Warnf("parse duration error... overriding Delay as 1 second")
		duration = time.Second
	}
	time.Sleep(duration)

	if t.Fail {
		return errors.New("some error")
	}
	return nil
}

func (t *testTask) OnFailure(err error) {
	t.Status = "failed"
}

func (t *testTask) OnSuccess() {
	t.Status = "completed"
}

func (t *testTask) IsCompleted() bool {
	return t.Status == "completed"
}
```

```GO
e := NewEngine(context.TODO()) // spawn engine's instance...
e.Start(20) // start your engine with 20 workers...
done1 := d.Do(task1) // executes task1...
done2 := d.Do(task2) // executes task2...
done3 := d.Do(task3) // executes task3...
...
done100 := d.Do(task100) // executes task100...
<-done55 // optional: ability to block logic until the task is completed
```

This not only lets you run your work concurrently but also allows you to handle error specifically on each type of work with a simple task pattern.

## Issues?

If you are using this library and encounter any issues please feel free to open an issue.