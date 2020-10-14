package engine

import (
	"context"
	"testing"
)

func TestEngine_Stop(t *testing.T) {
	type fields struct {
		stop chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{

		{
			name:   "wantSuccess - close worker",
			fields: fields{stop: make(chan struct{})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewEngine(context.TODO())
			d.Start(0)
			d.stop = tt.fields.stop
			d.Stop()
		})
	}
}

func TestEngine_dispatch(t *testing.T) {
	d := NewEngine(context.TODO())
	w := NewWorker(context.TODO(), d.pool)
	d.workers = append(d.workers, w)

	t.Parallel()
	type fields struct {
		task *testTask
	}
	tests := []struct {
		name        string
		fields      fields
		want        string
		wantSuccess bool
	}{
		{
			name:        "wantSuccess - work test task 1 - Fail false - Delay 100ms",
			fields:      fields{task: &testTask{ID: 1, Fail: false, Delay: "100ms"}},
			want:        "completed",
			wantSuccess: true,
		},
		{
			name:        "Fail - work test task 2 - Delay 200ms - Fail true",
			fields:      fields{task: &testTask{ID: 2, Delay: "200ms", Fail: true}},
			want:        "failed",
			wantSuccess: false,
		},
		{
			name:        "wantSuccess - work test task 3 - Fail false - Delay 300ms",
			fields:      fields{task: &testTask{ID: 3, Fail: false, Delay: "300ms"}},
			want:        "completed",
			wantSuccess: true,
		},
		{
			name:        "Fail - work test task 4 - Fail true - Delay 400ms",
			fields:      fields{task: &testTask{ID: 4, Fail: true, Delay: "400ms"}},
			want:        "failed",
			wantSuccess: false,
		},
		{
			name:        "wantSuccess - work test task 5 - Fail false - Delay 500ms",
			fields:      fields{task: &testTask{ID: 5, Fail: false, Delay: "500ms"}},
			want:        "completed",
			wantSuccess: true,
		},
		{
			name:        "Fail - work test task 6 - Delay 600ms - Fail true",
			fields:      fields{task: &testTask{ID: 6, Delay: "600ms", Fail: true}},
			want:        "failed",
			wantSuccess: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.Start(0)
			done := make(chan bool)
			d.input <- work{
				Executable: tt.fields.task,
				success:    done,
			}

			if success := <-done; success != tt.wantSuccess {
				t.Errorf("worker <- task failed wanted: %v got %v", tt.wantSuccess, success)
			}

			if tt.fields.task.Status != tt.want {
				t.Errorf("worker <- task failed wanted: %s got %s", tt.want, tt.fields.task.Status)
			}
		})
	}
}

func TestEngine_Do(t *testing.T) {
	d := NewEngine(context.TODO())
	w := NewWorker(context.TODO(), d.pool)
	d.workers = append(d.workers, w)

	t.Parallel()
	type fields struct {
		task *testTask
	}
	tests := []struct {
		name        string
		fields      fields
		want        string
		wantSuccess bool
	}{
		{
			name:        "wantSuccess - work test task 1 - Fail false - Delay 100ms",
			fields:      fields{task: &testTask{ID: 1, Fail: false, Delay: "100ms"}},
			want:        "completed",
			wantSuccess: true,
		},
		{
			name:        "Fail - work test task 2 - Delay 200ms - Fail true",
			fields:      fields{task: &testTask{ID: 2, Delay: "200ms", Fail: true}},
			want:        "failed",
			wantSuccess: false,
		},
		{
			name:        "wantSuccess - work test task 3 - Fail false - Delay 300ms",
			fields:      fields{task: &testTask{ID: 3, Fail: false, Delay: "300ms"}},
			want:        "completed",
			wantSuccess: true,
		},
		{
			name:        "Fail - work test task 4 - Fail true - Delay 400ms",
			fields:      fields{task: &testTask{ID: 4, Fail: true, Delay: "400ms"}},
			want:        "failed",
			wantSuccess: false,
		},
		{
			name:        "wantSuccess - work test task 5 - Fail false - Delay 500ms",
			fields:      fields{task: &testTask{ID: 5, Fail: false, Delay: "500ms"}},
			want:        "completed",
			wantSuccess: true,
		},
		{
			name:        "Fail - work test task 6 - Delay 600ms - Fail true",
			fields:      fields{task: &testTask{ID: 6, Delay: "600ms", Fail: true}},
			want:        "failed",
			wantSuccess: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.Start(0)
			done := d.Do(tt.fields.task)

			if success := <-done; success != tt.wantSuccess {
				t.Errorf("worker <- task failed wanted: %v got %v", tt.wantSuccess, success)
			}

			if tt.fields.task.Status != tt.want {
				t.Errorf("worker <- task failed wanted: %s got %s", tt.want, tt.fields.task.Status)
			}
		})
	}
}
