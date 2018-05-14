package audio

import (
	"errors"
	"mos6502go/system"

	ebiten_audio "github.com/hajimehoshi/ebiten/audio"
)

var (
	audioContext *ebiten_audio.Context
	player       *ebiten_audio.Player
	firstAudio   bool
	Mute         bool
)

type stream struct{}

func (s *stream) Read(data []byte) (int, error) {
	dataLen := len(data)

	if firstAudio {
		// The first time, drain the audio queue
		firstAudio = false

		for i := 0; i < len(system.AudioChannel); i++ {
			<-system.AudioChannel
		}
		return dataLen, nil
	}

	if dataLen%4 != 0 {
		return 0, errors.New("dataLen % 4 must be 0")
	}

	if Mute {
		return dataLen, nil
	}

	samples := dataLen / 4

	for i := 0; i < dataLen; i++ {
		data[i] = 0
	}

	for i := 0; i < samples; i++ {
		b := <-system.AudioChannel

		data[4*i] = byte(b)
		data[4*i+1] = byte(b >> 8)
		data[4*i+2] = byte(b)
		data[4*i+3] = byte(b >> 8)
	}

	return dataLen, nil
}

func (s *stream) Close() error {
	return nil
}

func Init() {
	system.AudioCycles = 0
	firstAudio = true
	Mute = false

	var err error
	audioContext, err = ebiten_audio.NewContext(system.AudioSampleRate)
	if err != nil {
		panic(err)
	}

	// Pass the (infinite) stream to audio.NewPlayer.
	// After calling Play, the stream never ends as long as the player object lives.
	// var err error
	player, err = ebiten_audio.NewPlayer(audioContext, &stream{})
	if err != nil {
		panic(err)
	}
	player.Play()
}

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
