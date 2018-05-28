package audio

// Very simple implementation of audio. Every frame, a channel is filled with
// all the audio samples from the last frame. Each time the speaker clicks, the
// channel is filled with the last audio samples. The channel is also
// filled at the end of the frame.

import "github.com/freewilll/apple2/system"

// Click handles a speaker click
func Click() {
	ForwardToFrameCycle()
	system.AudioAttenuationCounter = 400
	system.LastAudioValue = ^system.LastAudioValue
}

// attenuate makes sure the audio goes down to zero after a period of inactivity
func attenuate(sample int16) int16 {
	if system.AudioAttenuationCounter == 0 {
		return 0
	} else {
		system.AudioAttenuationCounter--
		return sample
	}
}

// ForwardToFrameCycle calculates how many audio samples need to be written to
// the channel based on how many CPU cycles have been executed since the last
// flush and shove them into the channel.
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
