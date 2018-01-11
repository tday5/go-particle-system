package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"psystem"
)

// Global variables
var (
	winTitle = "Particle System"

	imgPaths = []string{
		"/img/red_circle.png",
		"/img/black_circle.png",
		"/img/donut.png",
		"/img/donut2.png",
	}

	winWidth, winHeight = 800, 600
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("all 3 arguments not provided")
	}

	maxParticles, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse argument 1 because %s", err)
	}
	repulseStrength, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatalf("Failed to parse argument 2 because %s", err)
	}
	cor, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		log.Fatalf("Failed to parse argument 3 because %s", err)
	}
	repulseStrength *= 1000.0

	// SDL setup
	var window *sdl.Window
	var mux sync.Mutex
	var renderer *sdl.Renderer
	// SETUP
	var particleSystem = psystem.CreateSystem(maxParticles, repulseStrength, cor)

	// Load window
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil || window == nil {
		log.Fatalf("Failed to create window: %s\n", sdl.GetError())
	}
	defer window.Destroy()

	// Load renderer
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil || renderer == nil {
		log.Fatalf("Failed to create renderer: %s\n", sdl.GetError())
	}
	defer renderer.Destroy()

	// Load images
	images := make([]*sdl.Surface, len(imgPaths))
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	for i, path := range imgPaths {
		images[i], err = img.Load(gopath + path)
		if err != nil || images[i] == nil {
			log.Fatalf("Failed to load image: %s\n", sdl.GetError())
		}
		defer images[i].Free()
	}

	// Load textures from images
	textures := make([]*sdl.Texture, len(images))
	for i := range images {
		textures[i], err = renderer.CreateTextureFromSurface(images[i])
		if err != nil || textures[i] == nil {
			log.Fatalf("Failed to load texture: %s\n", sdl.GetError())
		}
		defer textures[i].Destroy()
	}

	// Event loop
	running := true
	// Toggle for creating a static
	staticCreate := false
	createRadius := 14
	// helper for changing radius at which particles are created
	changeCreateRadius := func(radius *int, newRadius int) {
		if newRadius > 9 && newRadius < 31 {
			(*radius) = newRadius
		}
	}
	//Toggle for food
	donuts := false
	for running {
		// Poll events
		event := sdl.PollEvent()
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
			fmt.Println("Exiting")
			// Mouse handling
		case *sdl.MouseButtonEvent:
			if t.Type == sdl.MOUSEBUTTONDOWN {
				if t.Button == sdl.BUTTON_LEFT {
					if staticCreate {
						particleSystem.AddNewStatic(float64(t.X), float64(t.Y), float64(createRadius))
					} else {
						particleSystem.AddNewMover(float64(t.X), float64(t.Y), float64(createRadius), float64(createRadius/10))
					}
				} else if t.Button == sdl.BUTTON_RIGHT {
					err := particleSystem.Repulse(float64(t.X), float64(t.Y))
					if err != nil {
						log.Fatalf("Error in repulsion: %s", err.Error())
					}
					fmt.Printf("Repulsed at: %v, %v\n", t.X, t.Y)
				}
			}
			// Key handling
		case *sdl.KeyboardEvent:
			if t.Type == sdl.KEYDOWN {
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_S:
					staticCreate = true
				case sdl.SCANCODE_D:
					donuts = !donuts
				case sdl.SCANCODE_LEFT:
					changeCreateRadius(&createRadius, createRadius-1)
					window.SetTitle("Create radius is now: " + strconv.Itoa(createRadius))
				case sdl.SCANCODE_RIGHT:
					changeCreateRadius(&createRadius, createRadius+1)
					window.SetTitle("Create radius is now: " + strconv.Itoa(createRadius))
				case sdl.SCANCODE_C:
					particleSystem.ClearSystem()
				}
			} else if t.Type == sdl.KEYUP {
				if t.Keysym.Scancode == sdl.SCANCODE_S {
					staticCreate = false
				}
			}
		}

		// Update system
		particleSystem.UpdateSystem(winWidth, winHeight)
		// Draw code
		renderer.Clear()
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(winWidth), H: int32(winHeight)})
		// Draw particles
		toDraw := particleSystem.GetParticleInfo()
		mod := 0
		if donuts {
			mod = 2
		}

		wg := sync.WaitGroup{}
		for _, p := range toDraw {
			wg.Add(1)
			go func(images []*sdl.Surface, textures []*sdl.Texture, renderer *sdl.Renderer, mux *sync.Mutex, p psystem.DrawInfo, mod int) {
				sourceImg := sdl.Rect{X: 0, Y: 0, W: images[p.DrawTag+mod].W, H: images[p.DrawTag+mod].H}
				destImg := sdl.Rect{X: int32(p.X - (p.Radius / 2.0)),
					Y: int32(p.Y - (p.Radius / 2.0)),
					W: int32(p.Radius * 2), H: int32(p.Radius * 2)}

				mux.Lock()
				renderer.Copy(textures[p.DrawTag+mod], &sourceImg, &destImg)
				mux.Unlock()
				wg.Done()
			}(images, textures, renderer, &mux, p, mod)
		}
		wg.Wait()
		renderer.Present()
	}
}
