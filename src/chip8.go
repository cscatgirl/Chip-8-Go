package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

const HEIGHT = 32
const WIDTH = 64
const SCALE = 10

func init() {
	runtime.LockOSThread()
}
func getNullRenderState() graphics.SfRenderStates {
	return (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0))
}
func render(buff [32][64]byte, pixel *graphics.Struct_SS_sfRectangleShape, r_window *graphics.Struct_SS_sfRenderWindow) {
	graphics.SfRenderWindow_clear(*r_window, graphics.GetSfBlack())
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if buff[y][x] == 1 {
				vec := graphics.NewSfVector2f()
				vec.SetX(float32(x * 10))
				vec.SetY(float32(y * 10))
				graphics.SfRectangleShape_setPosition(*pixel, vec)
				graphics.SfRenderWindow_drawRectangleShape(*r_window, *pixel, getNullRenderState())
			}
		}
	}

	graphics.SfRenderWindow_display(*r_window)
}
func handle_input_events(cpu *CPU) {
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyNum1)) == 1 {
		cpu.keys[0x1] = 1
	} else {
		cpu.keys[0x1] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyNum2)) == 1 {
		cpu.keys[0x2] = 1
	} else {
		cpu.keys[0x2] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyNum3)) == 1 {
		cpu.keys[0x3] = 1
	} else {
		cpu.keys[0x3] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyNum4)) == 1 {
		cpu.keys[0xC] = 1
	} else {
		cpu.keys[0xC] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyQ)) == 1 {
		cpu.keys[0x4] = 1
	} else {
		cpu.keys[0x4] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyW)) == 1 {
		cpu.keys[0x5] = 1
	} else {
		cpu.keys[0x5] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyE)) == 1 {
		cpu.keys[0x6] = 1
	} else {
		cpu.keys[0x6] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyR)) == 1 {
		cpu.keys[0xD] = 1
	} else {
		cpu.keys[0xD] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyA)) == 1 {
		cpu.keys[0x7] = 1
	} else {
		cpu.keys[0x7] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyS)) == 1 {
		cpu.keys[0x8] = 1
	} else {
		cpu.keys[0x8] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyD)) == 1 {
		cpu.keys[0x9] = 1
	} else {
		cpu.keys[0x9] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyF)) == 1 {
		cpu.keys[0xE] = 1
	} else {
		cpu.keys[0xE] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyZ)) == 1 {
		cpu.keys[0xA] = 1
	} else {
		cpu.keys[0xA] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyX)) == 1 {
		cpu.keys[0] = 1
	} else {
		cpu.keys[0] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyC)) == 1 {
		cpu.keys[0xB] = 1
	} else {
		cpu.keys[0xB] = 0
	}
	if window.SfKeyboard_isKeyPressed(window.SfKeyCode(window.SfKeyV)) == 1 {
		cpu.keys[0xF] = 1
	} else {
		cpu.keys[0xF] = 0
	}
}

func main() {
	args := os.Args[1:]
	font_set := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	vm := window.NewSfVideoMode()
	defer window.DeleteSfVideoMode(vm)
	vm.SetWidth(WIDTH * SCALE)
	vm.SetHeight(HEIGHT * SCALE)
	vm.SetBitsPerPixel(32)

	cs := window.NewSfContextSettings()
	defer window.DeleteSfContextSettings(cs)
	w := graphics.SfRenderWindow_create(vm, "CHIP-8-EMU", uint(window.SfResize|window.SfClose), cs)
	defer window.SfWindow_destroy(w)
	graphics.SfRenderWindow_setFramerateLimit(w, 60)
	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)
	pixel := graphics.SfRectangleShape_create()
	scale := graphics.NewSfVector2f()
	scale.SetX(SCALE)
	scale.SetY(SCALE)
	graphics.SfRectangleShape_setSize(pixel, scale)
	graphics.SfRectangleShape_setFillColor(pixel, graphics.GetSfWhite())
	var cpu = new(CPU)
	cpu.Init(font_set)
	fmt.Println(args[0])
	code, err := cpu.readRom(args[0])
	if code != 0 {
		fmt.Print(err)
	}
	for window.SfWindow_isOpen(w) > 0 {
		/* Process events */
		for window.SfWindow_pollEvent(w, ev) > 0 {
			/* Close window: exit */
			if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
				return
			}
		}
		for i := 0; i < 10; i++ {
			cpu.readByte()
		}

		handle_input_events(cpu)
		if cpu.delayTimer > 0 {
			cpu.delayTimer--
		}
		///time.Sleep(1000 / 60)
		if cpu.draw() {
			render(cpu.Screen, &pixel, &w)
		}
	}
}
