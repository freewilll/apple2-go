package main

import (
	"flag"
	"mos6502go/cpu"
	"mos6502go/mmu"
)

func main() {
	showInstructions := flag.Bool("show-instructions", false, "Show instructions code while running")
	disableBell := flag.Bool("disable-bell", false, "Disable bell")
	flag.Parse()

	cpu.InitDisasm()
	memory := mmu.InitRAM()
	mmu.InitApple2eROM(memory)

	var s cpu.State
	s.Memory = memory
	s.MemoryMap = &memory.MemoryMap
	s.Init()

	bootVector := 0xfffc
	lsb := (*s.MemoryMap)[uint8(bootVector>>8)][uint8(bootVector&0xff)] // TODO move readMemory to mmu
	msb := (*s.MemoryMap)[uint8((bootVector+1)>>8)][uint8((bootVector+1)&0xff)]
	s.PC = uint16(lsb) + uint16(msb)<<8
	cpu.Run(&s, *showInstructions, nil, *disableBell, 0)
}
