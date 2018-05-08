package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"mos6502go/cpu"
	"mos6502go/mmu"
	"mos6502go/vid"
)

const (
	screenSizeFactor = 1     // Factor by which the whole screen is resized
	textVideoMemory  = 0x400 // Base location of page 1 text video memory
	flashFrames      = 8     // Number of frames when FLASH mode is toggled
)

var (
	charMap      *ebiten.Image
	flashCounter int
	flashOn      bool
)

var cpuState cpu.State
var showInstructions *bool
var disableBell *bool

func update(screen *ebiten.Image) error {
	cpu.Run(&cpuState, *showInstructions, nil, *disableBell, 1024000/60)
	return vid.DrawTextScreen(cpuState.MemoryMap, screen, charMap)
}

func main() {
	showInstructions = flag.Bool("show-instructions", false, "Show instructions code while running")
	disableBell = flag.Bool("disable-bell", false, "Disable bell")
	flag.Parse()

	cpu.InitDisasm()
	memory := mmu.InitRAM()
	mmu.InitApple2eROM(memory)

	cpuState.Memory = memory
	cpuState.MemoryMap = &memory.MemoryMap
	cpuState.Init()

	bootVector := 0xfffc
	lsb := (*cpuState.MemoryMap)[uint8(bootVector>>8)][uint8(bootVector&0xff)] // TODO move readMemory to mmu
	msb := (*cpuState.MemoryMap)[uint8((bootVector+1)>>8)][uint8((bootVector+1)&0xff)]
	cpuState.PC = uint16(lsb) + uint16(msb)<<8

	var err error
	charMap, _, err = ebitenutil.NewImageFromFile("./pr-latin1.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.Run(update, 280*screenSizeFactor, 192*screenSizeFactor, 2, "Apple //e")
}
