package main

type VM struct {
	ID            int
	CPUCapacity   float64
	MemCapacity   float64
	BWCapacity    float64
	CurrentLoad   float64
	AssignedTasks []*Task
}

func NewVM(id int, cpu, mem, bw float64) *VM {
	return &VM{
		ID:            id,
		CPUCapacity:   cpu,
		MemCapacity:   mem,
		BWCapacity:    bw,
		CurrentLoad:   0.0,
		AssignedTasks: make([]*Task, 0),
	}
}

func (vm *VM) CanHandleTask(t *Task) bool {
	cpuUsage := vm.getCPUUsage() + (t.CPURequired / vm.CPUCapacity * 100)
	memUsage := vm.getMemUsage() + (t.MemRequired / vm.MemCapacity * 100)
	bwUsage := vm.getBWUsage() + (t.BWRequired / vm.BWCapacity * 100)
	return cpuUsage <= 100 && memUsage <= 100 && bwUsage <= 100
}

func (vm *VM) AssignTask(t *Task) {
	vm.AssignedTasks = append(vm.AssignedTasks, t)
	vm.updateLoad()
}

func (vm *VM) getCPUUsage() float64 {
	total := 0.0
	for _, t := range vm.AssignedTasks {
		total += t.CPURequired
	}
	return (total / vm.CPUCapacity) * 100
}

func (vm *VM) getMemUsage() float64 {
	total := 0.0
	for _, t := range vm.AssignedTasks {
		total += t.MemRequired
	}
	return (total / vm.MemCapacity) * 100
}

func (vm *VM) getBWUsage() float64 {
	total := 0.0
	for _, t := range vm.AssignedTasks {
		total += t.BWRequired
	}
	return (total / vm.BWCapacity) * 100
}

func (vm *VM) updateLoad() {
	vm.CurrentLoad = (vm.getCPUUsage() + vm.getMemUsage() + vm.getBWUsage()) / 3
}
