package main

import (
	"math"
)

type Swarm struct {
	Particles    []*Particle
	GBest        []int
	GBestFitness float64
	VMs          []*VM
	Tasks        []*Task
}

// Create a new swarm
func NewSwarm(numParticles, numTasks, numVMs int, vms []*VM, tasks []*Task) *Swarm {
	particles := make([]*Particle, numParticles)
	for i := 0; i < numParticles; i++ {
		particles[i] = NewParticle(numTasks, numVMs)
	}

	return &Swarm{
		Particles:    particles,
		GBest:        append([]int{}, particles[0].Position...),
		GBestFitness: math.Inf(1),
		VMs:          vms,
		Tasks:        tasks,
	}
}

// Calculate fitness based on load balance
func (s *Swarm) CalculateFitness(p *Particle) float64 {
	// Reset VM loads
	for _, vm := range s.VMs {
		vm.AssignedTasks = make([]*Task, 0)
		vm.CurrentLoad = 0.0
	}

	// Assign tasks to VMs
	for taskID, vmID := range p.Position {
		if vmID >= 0 && vmID < len(s.VMs) && taskID < len(s.Tasks) {
			if s.VMs[vmID].CanHandleTask(s.Tasks[taskID]) {
				s.VMs[vmID].AssignTask(s.Tasks[taskID])
			}
		}
	}

	// Compute load difference
	minLoad := math.Inf(1)
	maxLoad := math.Inf(-1)
	for _, vm := range s.VMs {
		if vm.CurrentLoad < minLoad {
			minLoad = vm.CurrentLoad
		}
		if vm.CurrentLoad > maxLoad {
			maxLoad = vm.CurrentLoad
		}
	}
	return maxLoad - minLoad // smaller is better
}

// Main PSO optimization loop
func (s *Swarm) Optimize(iterations int, w, c1, c2 float64) {
	for iter := 0; iter < iterations; iter++ {
		for _, p := range s.Particles {
			fitness := s.CalculateFitness(p)

			// Update particle's personal best
			if fitness < p.PBestFitness || p.PBestFitness == 0.0 {
				p.PBestFitness = fitness
				p.PBest = append([]int{}, p.Position...)
			}

			// Update global best
			if fitness < s.GBestFitness || s.GBestFitness == math.Inf(1) {
				s.GBestFitness = fitness
				s.GBest = append([]int{}, p.Position...)
			}

			// Move the particle
			p.UpdateVelocity(s.GBest, w, c1, c2)
			p.UpdatePosition(len(s.VMs))
		}
	}
}

// Return the best assignment found so far
func (s *Swarm) GetBestAssignment() []int {
	return s.GBest
}

// Return the best (lowest) fitness found so far
func (s *Swarm) GetBestFitness() float64 {
	return s.GBestFitness
}
