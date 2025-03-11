package main

import "math/rand"

type Particle struct {
	Position     []int
	Velocity     []float64
	PBest        []int
	PBestFitness float64
}

func NewParticle(numTasks, numVMs int) *Particle {
	position := make([]int, numTasks)
	velocity := make([]float64, numTasks)
	for i := 0; i < numTasks; i++ {
		position[i] = rand.Intn(numVMs)
		velocity[i] = rand.Float64() * 0.1
	}
	return &Particle{
		Position:     position,
		Velocity:     velocity,
		PBest:        append([]int{}, position...),
		PBestFitness: 0.0,
	}
}

func (p *Particle) UpdateVelocity(gBest []int, w, c1, c2 float64) {
	for i := range p.Velocity {
		r1, r2 := rand.Float64(), rand.Float64()
		cognitive := c1 * r1 * float64(p.PBest[i]-p.Position[i])
		social := c2 * r2 * float64(gBest[i]-p.Position[i])
		p.Velocity[i] = w*p.Velocity[i] + cognitive + social
	}
}

func (p *Particle) UpdatePosition(numVMs int) {
	for i := range p.Position {
		p.Position[i] = (p.Position[i] + int(p.Velocity[i])) % numVMs
		if p.Position[i] < 0 {
			p.Position[i] += numVMs
		}
	}
}
