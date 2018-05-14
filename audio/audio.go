package audio

import (
	"mos6502go/system"
)

func Click() {
	ForwardToFrameCycle()
	system.AudioAttenuationCounter = 400
	system.LastAudioValue = ^system.LastAudioValue
}

func attenuate(sample int16) int16 {
	if system.AudioAttenuationCounter == 0 {
		return 0
	} else {
		system.AudioAttenuationCounter--
		return sample
	}
}

func ForwardToFrameCycle() {
	ratio := float64(system.AudioSampleRate) / system.CpuFrequency

	samples := uint64(ratio * float64(system.FrameCycles-system.LastAudioCycles))
	var i uint64
	for i = 0; i < samples; i++ {
		b := attenuate(system.LastAudioValue)
		system.AudioChannel <- b
	}
	system.LastAudioCycles = system.FrameCycles
}
