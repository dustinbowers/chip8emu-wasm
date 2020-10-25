//+build wasm

package main

import (
	"fmt"
	"github.com/dustinbowers/chip8emu/chip8"
	"log"
	"syscall/js"
	"time"
)

const (
	width  = 800
	height = 400
)

var (
	document js.Value
	canvas js.Value
	context js.Value
	emu chip8.Chip8
	running bool
)

func init() {
	document = js.Global().Get("document")
}

func main() {
	fmt.Println("WASM loaded.")

	js.Global().Set("setTarget", setTargetWrapper())
	js.Global().Set("resetChip8", resetChip8Wrapper())

	<-make(chan bool)
}

func resetChip8(romBytes []byte) {
	emu.Initialize()
	emu.LoadRomBytes(romBytes)
	fmt.Printf("Loaded %d bytes into memory\n", len(romBytes))

	runEmu(700)
}
func resetChip8Wrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		array := args[0]
		buf := make([]byte, array.Get("length").Int())
		n := js.CopyBytesToGo(buf, args[0])
		fmt.Printf("Bytes copied: %d\n", n)
		resetChip8(buf)
		return nil
	})
}

func runEmu(hz int64) {
	running = true
	delay := time.Duration(1000 / hz)
	go func() {
		for {
			_, err := emu.EmulateCycle()
			if err != nil {
				log.Fatalf("emu.EmulateCycle: %v", err)
			}
			if running == false {
				return
			}
			time.Sleep(time.Microsecond * delay * 1000) // ~700 Hz
		}
	}()

	go func() {
		for {
			if emu.DrawFlag {
				emu.DrawFlag = false

				// Render Screen
				context.Call("clearRect", 0, 0, width, height)
				blockWidth := int32(width / 64)
				blockHeight := int32(height / 32)
				for c, col := range emu.Screen {
					for r, on := range col {
						xPos := int32(c) * blockWidth
						yPos := int32(r) * blockHeight

						if on == 1 {
							context.Set("fillStyle", "#FFFFFF")
						} else {
							context.Set("fillStyle", "#000000")
						}
						context.Call("fillRect", xPos, yPos, blockWidth, blockHeight)
					}
				}
			}
			time.Sleep(time.Microsecond * 16700)
		}
	}()
}

func setTarget(target string) {
	canvas = js.
		Global().
		Get("document").
		Call("getElementById", target)

	context = canvas.Call("getContext", "2d")

	// reset
	canvas.Set("height", height)
	canvas.Set("width", width)
	context.Call("clearRect", 0, 0, width, height)

}
func setTargetWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface {} {
		setTarget(args[0].String())
		return nil
	})
}
