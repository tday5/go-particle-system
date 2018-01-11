package psystem

import "testing"

func TestUpdate(t *testing.T) {

	particles := []particle{
		mover{genParticle{FloatPoint{1, 5}, 1}, FloatPoint{-1, 2}, FloatPoint{1.0, 0.0}, 1},
		mover{genParticle{FloatPoint{19, 5}, 1}, FloatPoint{1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{20, 2}, 1}, FloatPoint{-1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{6, 1}, 1}, FloatPoint{0, -2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{6, 19}, 2}, FloatPoint{0, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{9, 8}, 1}, FloatPoint{1, -1}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{9, 10}, 1}, FloatPoint{-1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{14, 14}, 1}, FloatPoint{1, 1}, FloatPoint{0.0, 0.0}, 1},
		static{genParticle{center: FloatPoint{15, 15}, radius: 1}}}

	afterOneUpdateElastic := []particle{
		mover{genParticle{FloatPoint{3, 7}, 1}, FloatPoint{2, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{18, 7}, 1}, FloatPoint{-1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{19, 4}, 1}, FloatPoint{-1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{6, 3}, 1}, FloatPoint{0, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{6, 17}, 2}, FloatPoint{0, -2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{10, 10}, 1}, FloatPoint{1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{8, 9}, 1}, FloatPoint{-1, -1}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{13.01980198019802, 13.01980198019802}, 1}, FloatPoint{-0.9801980198019802, -0.9801980198019802}, FloatPoint{0.0, 0.0}, 1},
		static{genParticle{center: FloatPoint{15, 15}, radius: 1}}}

	for i, p := range particles {
		assertEqual(p.update(20, 20, particles, i, 1.0), afterOneUpdateElastic[i], t)
	}
}

func TestCollision(t *testing.T) {

	movers := []mover{
		mover{genParticle{FloatPoint{9, 8}, 1}, FloatPoint{1, -1}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{9, 10}, 1}, FloatPoint{-1, 2}, FloatPoint{0.0, 0.0}, 1},
		mover{genParticle{FloatPoint{14, 14}, 1}, FloatPoint{1, 1}, FloatPoint{0.0, 0.0}, 1}}

	s := static{genParticle{center: FloatPoint{15, 15}, radius: 1}}

	assertEqual(collision(movers[0], movers[1], 1.0), FloatPoint{1, 2}, t)
	assertEqual(collision(movers[1], movers[0], 1.0), FloatPoint{-1, -1}, t)
	assertEqual(collision(movers[2], s, 1.0), FloatPoint{-0.9801980198019802, -0.9801980198019802}, t)
}

func TestAddForce(t *testing.T) {
	assertEqual(mover{velocity: FloatPoint{0.0, 0.0}, mass: 5.0}.addForce(
		FloatPoint{10.0, 50.0}).update(
		800, 600, []particle{}, 0, 1.0).parVelocity(), FloatPoint{2.0, 10.0}, t)

	assertEqual(static{}.addForce(FloatPoint{100.0, 100.0}).parVelocity(), FloatPoint{0.0, 0.0}, t)
}

func TestParMass(t *testing.T) {
	assertEqual(mover{mass: 54.5}.parMass(), 54.5, t)
	assertEqual(static{}.parMass(), 100.0, t)
}

func TestParCenter(t *testing.T) {
	assertEqual(genParticle{center: FloatPoint{40.5, 5.5}}.parCenter(),
		FloatPoint{40.5, 5.5}, t)
}

func TestParRadius(t *testing.T) {
	assertEqual(genParticle{radius: 3.2}.parRadius(), 3.2, t)
}

func TestParVelocity(t *testing.T) {
	assertEqual(mover{velocity: FloatPoint{50.0, 76239.8}}.parVelocity(),
		FloatPoint{50.0, 76239.8}, t)

	assertEqual(static{}.parVelocity(), FloatPoint{0.0, 0.0}, t)
}

func TestParImage(t *testing.T) {
	assertEqual(mover{}.parImage(), 0, t)
	assertEqual(static{}.parImage(), 1, t)
}
