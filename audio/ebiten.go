package audio

// This file contains the consumer part of the audio code. audio.go is responsible for producing to it

import (
	"errors"

	"github.com/freewilll/apple2-go/system"
	ebiten_audio "github.com/hajimehoshi/ebiten/audio"
)

var (
	audioContext *ebiten_audio.Context // Ebitem audio context
	player       *ebiten_audio.Player  // Ebitem stream player
	firstAudio   bool                  // True at startup
)

var (
	// Mute ensures no samples are output
	Mute bool

	// ClickWhenDriveHeadMoves makes the speaker click once every time the stepper motor magnets change
	ClickWhenDriveHeadMoves bool
)

// The streaming code is based on the ebiten sinewave example
type stream struct{}

// Read is called whenever the sound hardware wants some samples. Convert the
// 16 bit data in the sound buffer to 8 bit stereo values.
func (s *stream) Read(data []byte) (int, error) {
	dataLen := len(data)

	if firstAudio {
		// The first time, drain the audio queue and exit
		firstAudio = false

		for i := 0; i < len(system.AudioChannel); i++ {
			<-system.AudioChannel
		}
		return dataLen, nil
	}

	// Sanity test
	if dataLen%4 != 0 {
		return 0, errors.New("dataLen % 4 must be 0")
	}

	// Do nothing if we're muted, but ensure the channel keeps getting drained
	if Mute {
		firstAudio = true
		return dataLen, nil
	}

	// Consume the samples from the channel
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

// Close is called when the program exits
func (s *stream) Close() error {
	return nil
}

// InitEbiten initializes the audio sets up the ebiten output stream
func InitEbiten() {
	// Setup initial state
	firstAudio = true
	Mute = false
	ClickWhenDriveHeadMoves = false

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
