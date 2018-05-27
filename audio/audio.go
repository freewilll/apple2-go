package audio

import (
	"apple2/system"
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
	// 1023000/44100=23.19 cycles per audio sample
	cyclesPerAudioSample := system.CpuFrequency / float64(system.AudioSampleRate)

	// Should be about 1023000/60=17050
	elapsedCycles := system.FrameCycles - system.LastAudioCycles

	// Should be about 17050/23.19=735 audio samples per frame
	audioSamples := uint64(float64(elapsedCycles) / cyclesPerAudioSample)

	for i := uint64(0); i < audioSamples; i++ {
		b := attenuate(system.LastAudioValue)
		system.AudioChannel <- b
	}
	system.LastAudioCycles = system.FrameCycles
}
