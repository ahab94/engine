package engine

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

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

func nTasks(n int) []Executable {
	tasks := make([]Executable, 0)
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().Unix())
		tasks = append(tasks, &testTask{
			ID:    i,
			Delay: fmt.Sprintf("%dms", rand.Intn(999)),
		})
	}
	return tasks
}
