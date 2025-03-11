package main

type Task struct {
	ID          int
	CPURequired float64
	MemRequired float64
	BWRequired  float64
	ExecTime    float64
}

func NewTask(id int, cpu, mem, bw, exec float64) *Task {
	return &Task{
		ID:          id,
		CPURequired: cpu,
		MemRequired: mem,
		BWRequired:  bw,
		ExecTime:    exec,
	}
}
