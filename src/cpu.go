package main

import (
	"bufio"
	"math/rand/v2"
	"os"
)

type CPU struct {
	RAM            [4096]byte
	Regs           [16]byte
	I              uint16
	Delay          byte
	Sound          byte
	PC             uint16
	SP             byte
	Stack          [16]uint16
	Screen         [32][64]byte
	shouldDraw     bool
	keys           [16]uint8
	soundTimer     uint8
	delayTimer     uint8
	lastkeyPressed int
	wasPressed     bool
}

func (cpu *CPU) draw() bool {
	draw := cpu.shouldDraw
	cpu.shouldDraw = false
	return draw
}
func (cpu *CPU) Init(fontset []byte) {
	cpu.PC = 0x200
	cpu.SP = 0
	for i := 0; i < len(fontset); i++ {
		cpu.RAM[i] = fontset[i]
	}
	cpu.shouldDraw = false
	cpu.wasPressed = false
}
func (cpu *CPU) readRom(fileName string) (rcode int, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 1, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return 1, err
	}
	size := int(fileinfo.Size())
	buff := make([]byte, size)
	reader := bufio.NewReader(file)
	_, err = reader.Read(buff)
	for i := 0; i < size; i++ {
		cpu.RAM[i+0x200] = buff[i]
	}
	return 0, nil
}
func (cpu *CPU) readByte() {
	var opt_code uint16 = uint16(cpu.RAM[cpu.PC])<<8 | uint16(cpu.RAM[cpu.PC+1])
	cpu.PC += 2
	var masked_off uint16 = opt_code & 0xF000
	switch masked_off {
	case 0x0000:
		var sub_mask uint16 = opt_code & 0x00FF
		switch sub_mask {
		case 0x00E0:
			for i := 0; i < len(cpu.Screen); i++ {
				for j := 0; j < len(cpu.Screen[i]); j++ {
					cpu.Screen[i][j] = 0x0
				}
			}
			cpu.shouldDraw = true
		case 0x00EE:
			cpu.SP = cpu.SP - 1
			cpu.PC = cpu.Stack[cpu.SP]
		}
	case 0x1000:
		cpu.PC = opt_code & 0x0FFF
	case 0x2000:
		cpu.Stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = opt_code & 0x0FFF
	case 0x3000:
		v := (opt_code & 0x0F00) >> 8
		kk := uint8(opt_code & 0x00FF)
		if cpu.Regs[v] == kk {
			cpu.PC += 2
		}
	case 0x4000:
		v := (opt_code & 0x0F00) >> 8
		kk := uint8(opt_code & 0x00FF)
		if cpu.Regs[v] != kk {
			cpu.PC += 2
		}
	case 0x5000:
		v_one := (opt_code & 0x0F00) >> 8
		v_two := (opt_code & 0x00F0) >> 8
		if cpu.Regs[v_one] == cpu.Regs[v_two] {
			cpu.PC += 2
		}
	case 0x6000:
		v := (opt_code & 0x0F00) >> 8
		kk := uint8(opt_code & 0x00FF)
		cpu.Regs[v] = kk
	case 0x7000:
		v := (opt_code & 0x0F00) >> 8
		vx := cpu.Regs[v]
		kk := uint8(opt_code & 0x00FF)
		cpu.Regs[v] = vx + kk
	case 0x8000:
		var sub_mask uint16 = opt_code & 0x000F
		switch sub_mask {
		case 0x0000:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			cpu.Regs[vx] = cpu.Regs[vy]
		case 0x0001:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			val_x := cpu.Regs[vx]
			val_y := cpu.Regs[vy]
			cpu.Regs[vx] = val_x | val_y
			cpu.Regs[15] = 0
		case 0x0002:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			val_x := cpu.Regs[vx]
			val_y := cpu.Regs[vy]
			cpu.Regs[vx] = val_x & val_y
			cpu.Regs[15] = 0

		case 0x0003:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			val_x := cpu.Regs[vx]
			val_y := cpu.Regs[vy]
			cpu.Regs[vx] = val_x ^ val_y
			cpu.Regs[15] = 0

		case 0x0004:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			val_x := cpu.Regs[vx]
			val_y := cpu.Regs[vy]
			sum := uint16(val_x) + uint16(val_y)
			cpu.Regs[vx] = uint8(sum & 0xFF)
			if sum > 0xFF {
				cpu.Regs[0xF] = 1
			} else {
				cpu.Regs[0xF] = 0
			}

		case 0x0005:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			val_x := cpu.Regs[vx]
			val_y := cpu.Regs[vy]
			cpu.Regs[vx] = val_x - val_y
			if val_y > val_x {
				cpu.Regs[15] = 0
			} else {
				cpu.Regs[15] = 1
			}

		case 0x0006:
			vx := (opt_code & 0x0F00) >> 8
			val_x := cpu.Regs[vx]
			cpu.Regs[vx] >>= 1
			cpu.Regs[15] = val_x & 1

		case 0x0007:
			vx := (opt_code & 0x0F00) >> 8
			vy := (opt_code & 0x00F0) >> 4
			val_x := cpu.Regs[vx]
			val_y := cpu.Regs[vy]
			cpu.Regs[vx] = val_y - val_x
			if val_x > val_y {
				cpu.Regs[15] = 0
			} else {
				cpu.Regs[15] = 1
			}

		case 0x000E:
			vx := (opt_code & 0x0F00) >> 8
			val_x := cpu.Regs[vx]
			cpu.Regs[vx] <<= 1
			cpu.Regs[15] = val_x >> 7

		}
	case 0x9000:
		vx := (opt_code & 0x0F00) >> 8
		vy := (opt_code & 0x00F0) >> 4
		val_x := cpu.Regs[vx]
		val_y := cpu.Regs[vy]
		if val_x != val_y {
			cpu.PC += 2
		}
	case 0xA000:
		cpu.I = opt_code & 0x0FFF
	case 0xB000:
		cpu.PC = (opt_code & 0x0FFF) + uint16(cpu.Regs[0])
	case 0xC000:
		vx := (opt_code & 0x0F00) >> 8
		cpu.Regs[vx] = uint8((opt_code & 0x0FFF)) & uint8(rand.IntN(255))
	case 0xD000:
		vx := (opt_code & 0x0F00) >> 8
		vy := (opt_code & 0x00F0) >> 4
		x := cpu.Regs[vx]
		y := cpu.Regs[vy]
		n := opt_code & 0x000F
		var row uint16 = 0
		var col uint16 = 0
		for row = 0; row < n; row++ {
			sprite_byte := cpu.RAM[cpu.I+row]
			for col = 0; col < 8; col++ {
				sprite_pixel := sprite_byte & (0x80 >> col)
				if sprite_pixel != 0 {
					y := (y + uint8(row)) % 32
					x := (x + uint8(col)) % 64
					curr_pixle := &cpu.Screen[y][x]
					if *curr_pixle == 1 {
						cpu.Regs[15] = 1
					} else {
						cpu.Regs[15] = 0
					}
					*curr_pixle ^= 1

				}
			}
		}
		cpu.shouldDraw = true
	case 0xE000:
		sub_mask := opt_code & 0x0FF
		switch sub_mask {
		case 0x09E:
			vx := (opt_code & 0x0F00) >> 8
			if cpu.keys[cpu.Regs[vx]] == 1 {
				cpu.PC += 2
			}
		case 0x00A1:
			vx := (opt_code & 0x0F00) >> 8
			if cpu.keys[cpu.Regs[vx]] == 0 {
				cpu.PC += 2
			}
		}
	case 0xF000:
		sub_mask := opt_code & 0x00FF
		switch sub_mask {
		case 0x0007:
			vx := (opt_code & 0x0F00) >> 8
			cpu.Regs[vx] = cpu.delayTimer
		case 0x000A:

			vx := (opt_code & 0x0F00) >> 8
			if !cpu.wasPressed {
				pressed := false
				for key := 0; key < len(cpu.keys); key++ {
					if cpu.keys[key] != 0 {
						cpu.Regs[vx] = byte(key)
						pressed = true
						cpu.wasPressed = true
						cpu.lastkeyPressed = key
						break
					}
				}
				if !pressed {
					cpu.PC -= 2
					return
				}
			} else {
				if cpu.keys[cpu.lastkeyPressed] != 0 {
					cpu.PC -= 2
					return

				}
				cpu.wasPressed = false
			}
		case 0x0015:
			vx := (opt_code & 0x0F00) >> 8
			cpu.delayTimer = cpu.Regs[vx]
		case 0x0018:
			vx := (opt_code & 0x0F00) >> 8
			cpu.soundTimer = cpu.Regs[vx]
		case 0x001E:
			vx := (opt_code & 0x0F00) >> 8
			cpu.I = cpu.I + uint16(cpu.Regs[vx])
		case 0x0029:
			vx := (opt_code & 0x0F00) >> 8
			cpu.I = uint16(cpu.Regs[vx]) * 0x5
		case 0x0033:
			vx := (opt_code & 0x0F00) >> 8
			x_val := cpu.Regs[vx]
			cpu.RAM[cpu.I] = x_val / 100
			cpu.RAM[cpu.I+1] = (x_val / 10) % 10
			cpu.RAM[cpu.I+2] = (x_val % 100) % 10
		case 0x0055:
			vx := (opt_code & 0x0F00) >> 8
			for i := 0; i < int(vx+1); i++ {
				cpu.RAM[cpu.I+uint16(i)] = cpu.Regs[i]
			}
			//cpu.I = vx + 1
			//cpu.I++
		case 0x0065:
			vx := (opt_code & 0x0F00) >> 8
			for i := 0; i < int(vx)+1; i++ {
				cpu.Regs[i] = cpu.RAM[cpu.I+uint16(i)]
			}
			//cpu.I = vx + 1
			//cpu.I++
		}
	}

}
