//Hi! Thanks for reading! Here's some info:
//This is my first ever (hopefully) finished Go project :D
//No AI code was used.
//I did make use of other online resources, all of which I've tried to reference and credit at the bottom of this file and/or the README.
//The only use of AI was learning, research and I occasionally asked it about optimisation tips.

package main

import (
	"image/color"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//Define the particles stored in the Game struct :)
type Game struct{
	particles []Particle
}

var size int = 30
var largerSize int = 70
var gravity float32 = 0.5
//A struct containing the structure of each particle. 
type Particle struct {
	PosX float32 `json:"x"`
	PosY float32 `json:"y"`
	VelY float32 `json:"velY"`
	VelX float32 `json:"velX"`
}

//Globally accessible :D
var smallCircle *ebiten.Image
var otherCircle *ebiten.Image
func init () {
//Define a new "screen"
	smallCircle = ebiten.NewImage(size, size)
	otherCircle = ebiten.NewImage(size * 3, size * 3)
//Make a centered circle with a diameter of 10 pixels
	vector.FillCircle(smallCircle, float32(size / 2), float32(size / 2), float32(size / 2), color.White, true)
	vector.FillCircle(otherCircle, float32(largerSize / 2),   float32(largerSize / 2),   float32(largerSize / 2), color.RGBA{R: 255, G: 121, B: 104, A: 1}, true)
}

func (g *Game) handleCollisions(i int) {
	for j := i + 1; j < len(g.particles); j++ {
		//Make sure the particle doesn't check if it's colliding with itself!
		if i >= j {
			continue
		}
		//Calculate individual distances
		distanceX := g.particles[i].PosX - g.particles[j].PosX 
		distanceY := g.particles[i].PosY - g.particles[j].PosY
		//Calculate full squared distance
		distanceSqrd := (distanceX * distanceX) + (distanceY * distanceY)
		//Square root distance
		distance := math.Sqrt(float64(distanceSqrd));
		//Find normalised distances
		vectorX := distanceX / float32(distance)
		vectorY := distanceY / float32(distance)
		//Find overlap
		var overlap float32 = 0
		if i == 0 {
			overlap = float32(largerSize) - float32(distance)
			//TODO: Fix the larger size
		} else {
			overlap = float32(size) - float32(distance)
		}
		
		var damping float32 =  0.25;
		
		if math.Sqrt(float64(distanceSqrd)) < float64(overlap + float32(distance)) && overlap > 0 {
			//Multiply the vector by the overlap
			//Then add a damping value.
			g.particles[j].VelX -= vectorX * overlap * damping;
			g.particles[j].VelY -= vectorY * overlap * damping;
			g.particles[i].VelX += vectorX * overlap * damping;
			g.particles[i].VelY += vectorY * overlap * damping;
		}
	}
}

func (g *Game) Update() error {
	//If the particle list is empty, populate it with randomly placed particles
	if len(g.particles) == 0 {
		for i := 0; i < 2000; i++ {
			tempPart := Particle{
				
				//PosX: float32(200 + (i * 10)),
				//PosY: float32(200 + (i * 90)),

				//Random positions :D
				PosX: float32(rand.IntN(1280) + 40),
				PosY: float32(rand.IntN(880) + 40),
			}
			//Append the newly created particle to the list!
			g.particles = append(g.particles, tempPart)
			
		}
	}
//Loop through every particle, checking "am I at the bottom of the screen?". If not, continue!
	for i := range g.particles {
		//Increase the position by the velocity
		g.particles[i].PosY += float32(g.particles[i].VelY);
		g.particles[i].PosX += float32(g.particles[i].VelX);

		if(!(g.particles[i].PosY > float32(960 - size))) {
			//Accelerate with "gravity"
			g.particles[i].VelY += 0.2
		} else {
			//Make sure particles don't surpass the limits of the screen
			g.particles[i].PosY = float32(960 - size);
			
			//Add some bounce!
			g.particles[i].VelY = (g.particles[i].VelY * -1) / 2;
			
		}

		if(g.particles[i].PosY >= (960 - float32(size) )) {
			g.particles[i].VelX *= 0.90
		}
	
	
		if((g.particles[i].PosX > float32(1280 + size * 2))){
			//Make sure particles don't surpass the limits of the screen
			g.particles[i].PosX = float32(1280 + size * 2);
			g.particles[i].VelX = (g.particles[i].VelX * -1) / 2;
			
		} else if((g.particles[i].PosX < float32(0 - size))) {
			g.particles[i].PosX = float32(0 - size);
			g.particles[i].VelX = (g.particles[i].VelX * -1) / 2;
		} 

		
		g.handleCollisions(i);
	}

	x, y := ebiten.CursorPosition()
	g.particles[0].PosX = float32(x) - float32(size * 3) / 2
	g.particles[0].PosY = float32(y) - float32(size * 3) / 2
	g.particles[0].VelX = 0
	g.particles[0].VelY = 0

	return nil
}

//For the sake of transparency,
//The below was optimised with research help partially from AI and the official Ebitengine docs.
//To clarify, the code was made and typed manually. I used Gemini for information about how the option parameters work.

func (g *Game) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options2 := &ebiten.DrawImageOptions{}
	for i, part := range g.particles {
		options.GeoM.Reset()
		options.GeoM.Translate(float64(part.PosX), float64(part.PosY))

		options2.GeoM.Reset()
		options2.GeoM.Translate(float64(part.PosX), float64(part.PosY))

		if(i == 0) {
			screen.DrawImage(otherCircle, options2)
		} else {
			screen.DrawImage(smallCircle, options)
		}
		
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 960
}

func main() {
	ebiten.SetWindowTitle("AquaCapy")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

//Resources used:
//https://ebitengine.org/en/documents/install.html - This is the boilerplate code used to get the window open!
//https://ebitengine.org/en/documents/cheatsheet.html#ebiten.DrawImageOptions
//https://ebitengine.org/en/tour/image.html