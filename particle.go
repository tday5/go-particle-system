package psystem

import "math"

// Q: overwriting methods?
// embedded structs
// better errors
// error handling with goroutines
// mutex
// concurrent draw code

// particle represents a single particle in the system with physical attributes
type particle interface {
	// Per frame update
	update(windowHeight, windowWidth int, particles []particle, avoid int, cor float64) particle

	// addForce adds the given vector of force to the particle's velocity
	addForce(force FloatPoint) particle

	parMass() float64

	parCenter() FloatPoint

	parRadius() float64

	parVelocity() FloatPoint

	parImage() int
}

type genParticle struct {
	center FloatPoint

	radius float64
}

func (par genParticle) parCenter() FloatPoint {
	return par.center
}

func (par genParticle) parRadius() float64 {
	return par.radius
}

// mover represents a moving 2D particle
type mover struct {
	genParticle

	velocity, acceleration FloatPoint

	mass float64
}

// Update updates the mover position and velocity, checks for collisions
func (par mover) update(windowHeight, windowWidth int, particles []particle, avoid int, cor float64) particle {
	futureCenter := par.center.Add(par.velocity)
	// Check for collision with walls
	if futureCenter.X-par.radius <= 0.0 && futureCenter.Sub(par.center).X < 0.0 ||
		futureCenter.X+par.radius >= float64(windowWidth) && futureCenter.Sub(par.center).X > 0.0 {
		par.velocity.X = -par.velocity.X
	}
	if futureCenter.Y-par.radius <= 0.0 && futureCenter.Sub(par.center).Y < 0.0 ||
		futureCenter.Y+par.radius >= float64(windowHeight) && futureCenter.Sub(par.center).Y > 0.0 {
		par.velocity.Y = -par.velocity.Y
	}
	// Check for collisions with other particles
	for i, p := range particles {
		if i != avoid && math.Abs(Dist(par.center, p.parCenter())) <=
			par.radius+p.parRadius() {
			par.velocity = collision(par, p, cor)
		}
	}

	// Update velocity and move the particle
	par.velocity = par.velocity.Add(par.acceleration)
	par.acceleration = par.acceleration.Mul(0.0)
	par.center = par.center.Add(par.velocity)
	return par
}

// collision aids in determining velocity for mover a after a collision with
// particle b
func collision(a mover, b particle, cor float64) FloatPoint {
	m := (2 * b.parMass()) / (a.mass + b.parMass())
	aSubB := a.center.Sub(b.parCenter())
	v := DotProduct(a.velocity.Sub(b.parVelocity()), aSubB) / DotProduct(aSubB, aSubB)

	// equation for particle collisions
	vx := a.velocity.Sub(aSubB.Mul(m * v))
	return vx.Mul(cor)
}

func (par mover) addForce(force FloatPoint) particle {
	par.acceleration = par.acceleration.Add(force.Div(par.mass))
	return par
}

// Helpers to define phyical attributes
func (par mover) parMass() float64 {
	return par.mass
}

func (par mover) parVelocity() FloatPoint {
	return par.velocity
}

func (par mover) parImage() int {
	return 0
}

// static represents a static non-moving particle with large mass
type static struct {
	genParticle
}

func (par static) update(windowHeight, windowWidth int, particles []particle, avoid int, cor float64) particle {
	return par
}

func (par static) addForce(force FloatPoint) particle {
	return par
}

func (par static) parMass() float64 {
	return 100.0
}

func (par static) parVelocity() FloatPoint {
	return FloatPoint{0.0, 0.0}
}

func (par static) parImage() int {
	return 1
}
