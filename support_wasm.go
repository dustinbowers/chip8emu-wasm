package main

import (
	"fmt"
	"github.com/dustinbowers/chip8emu/chip8"
	"syscall/js"
)

func addSupport() {
	js.Global().Set("pause", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Emu.Pause()
		return nil
	}))
	js.Global().Set("resume", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Emu.Resume()
		return nil
	}))
	js.Global().Set("resetChip8", resetChip8Wrapper())
}

func resetChip8(romBytes []byte) {
	g.Emu = chip8.NewChip8()
	g.Emu.Reset()
	g.Emu.SetBeepHandler(func(t bool) { fmt.Println("BEEEP") })
	g.Emu.LoadRomBytes(romBytes)
	fmt.Printf("Loaded %d bytes into memory\n", len(romBytes))
	g.Running = true
	fmt.Printf("g: %+v\n", g)

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



