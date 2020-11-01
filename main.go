//+build wasm, darwin

package main

import (
	"github.com/dustinbowers/chip8emu/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

const (
	screenWidth  = 64
	screenHeight = 32
	hz           = 700
)

type Game struct {
	Screen  *image.RGBA
	Emu     *chip8.Chip8
	Running bool
}

var KeyMap = map[ebiten.Key]uint8{
	ebiten.Key1: 0x1,
	ebiten.Key2: 0x2,
	ebiten.Key3: 0x3,
	ebiten.Key4: 0xC,
	ebiten.KeyQ: 0x4,
	ebiten.KeyW: 0x5,
	ebiten.KeyE: 0x6,
	ebiten.KeyR: 0xD,
	ebiten.KeyA: 0x7,
	ebiten.KeyS: 0x8,
	ebiten.KeyD: 0x9,
	ebiten.KeyF: 0xE,
	ebiten.KeyZ: 0xA,
	ebiten.KeyX: 0x0,
	ebiten.KeyC: 0xB,
	ebiten.KeyV: 0xF,
}

func (g *Game) HandleInput() {
	for k, v := range KeyMap {
		if ebiten.IsKeyPressed(k) {
			g.Emu.KeyDown(v)
		} else {
			g.Emu.KeyUp(v)
		}
	}
}

func (g *Game) Update() error {
	if g.Running {
		g.HandleInput()
		_, err := g.Emu.EmulateCycle()
		if err != nil {
			log.Fatalf("emu.EmulateCycle: %v", err)
		}
	}
	if g.Emu.DrawFlag {
		g.Emu.DrawFlag = false

		// Render Screen
		for c, col := range g.Emu.Screen {
			for r, on := range col {
				offset := 4*((r*screenWidth)+c)
				if offset >= 8188 {
					offset = 8188
				}
				if on >= 1 {
					g.Screen.Pix[offset+0] = 0xFF
					g.Screen.Pix[offset+1] = 0xFF
					g.Screen.Pix[offset+2] = 0xFF
					g.Screen.Pix[offset+3] = 0xFF
				} else {
					g.Screen.Pix[offset+0] = 0x00
					g.Screen.Pix[offset+1] = 0x00
					g.Screen.Pix[offset+2] = 0x00
					g.Screen.Pix[offset+3] = 0x00
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.Screen.Pix)
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var g *Game
func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Chip8 Emulator")
	ebiten.SetMaxTPS(hz)
	g = &Game{
		Screen: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		Emu: chip8.NewChip8(),
		Running: false,
	}
	addSupport()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
