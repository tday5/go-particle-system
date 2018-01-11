package psystem

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

// System represents a grouping of particles and manages their interactions
type System struct {
	particles []particle

	maxParticles, numParticles int

	repulseStrength, cor float64

	mux sync.Mutex
}

// CreateSystem creates a system with the given attributes
func CreateSystem(maxParticles int, repulseStrength, cor float64) System {
	return System{maxParticles: maxParticles, repulseStrength: repulseStrength, cor: cor}
}

// ClearSystem clears the sysem of all its particles
func (sys *System) ClearSystem() {
	sys.particles = []particle{}
}

// UpdateSystem updates the system of particles by one frame
func (sys *System) UpdateSystem(windowWidth, windowHeight int) {
	temp := make([]particle, len(sys.particles))
	for i, p := range sys.particles {
		temp[i] = p
	}
	wg := sync.WaitGroup{}
	for i, p := range sys.particles {
		wg.Add(1)
		go func(i int, p particle) {
			temp := p.update(windowHeight, windowWidth, temp, i, sys.cor)
			sys.mux.Lock()
			sys.particles[i] = temp
			sys.mux.Unlock()
			wg.Done()
		}(i, p)
	}
	wg.Wait()
}

// AddNewMover adds a new mover to the system at the given coordinates
func (sys *System) AddNewMover(x, y, radius, mass float64) {
	if sys.checkParticleCountPosAndAdd(x, y, radius) {
		var v FloatPoint
		if sys.cor >= 1.0 {
			v = FloatPoint{0, 0}
		} else {
			v = FloatPoint{rand.Float64(), rand.Float64()}
		}
		sys.particles = append(sys.particles, mover{genParticle{FloatPoint{x, y}, radius},
			v, FloatPoint{0.0, 0.0}, mass})
		fmt.Printf("Created particle at: %v, %v\n", x, y)
	}
}

// AddNewStatic adds a new static to the system at given coordinates
func (sys *System) AddNewStatic(x, y, radius float64) {
	if sys.checkParticleCountPosAndAdd(x, y, radius) {
		sys.particles = append(sys.particles, static{genParticle{FloatPoint{x, y}, radius}})
		fmt.Printf("Created static at: %v, %v\n", x, y)
	}
}

// helper for adding particles
func (sys *System) checkParticleCountPosAndAdd(x, y, radius float64) bool {
	if sys.numParticles < sys.maxParticles {
		for _, p := range sys.particles {
			if math.Abs(Dist(FloatPoint{x, y}, p.parCenter())) < radius+p.parRadius() {
				fmt.Println("Particle at that location already")
				return false
			}
		}
		sys.numParticles++
		return true
	}
	fmt.Println("Max particles reached")
	return false
}

// Repulse repulses all particles in the system accordingly
func (sys *System) Repulse(x, y float64) error {
	origin := FloatPoint{x, y}

	wait := make(chan error)
	for i, p := range sys.particles {
		go func(i int, p particle, sys *System, wait chan error) {
			// Repulsion equation
			dir := origin.Sub(p.parCenter())
			mag := Mag(dir)
			dir, err := Normalize(dir)
			if err != nil {
				wait <- err
				return
			}
			force := -1.0 * sys.repulseStrength / (mag * mag)
			dir = dir.Mul(force)

			sys.mux.Lock()
			sys.particles[i] = p.addForce(dir)
			sys.mux.Unlock()
			wait <- nil
		}(i, p, sys, wait)
	}
	for range sys.particles {
		if <-wait != nil {
			return ZeroMagError{}
		}
	}
	return nil
}

// DrawInfo represents necessary information for drawing a particle system
type DrawInfo struct {
	X, Y, Radius float64

	DrawTag int
}

// GetParticleInfo returns DrawInfo for drawing particles
func (sys *System) GetParticleInfo() []DrawInfo {
	temp := make([]DrawInfo, len(sys.particles))
	for i, p := range sys.particles {
		temp[i] = DrawInfo{p.parCenter().X, p.parCenter().Y, p.parRadius(), p.parImage()}
	}
	return temp
}
