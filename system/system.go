package system

// The system package is a dumping ground for globals that are shared between
// the packages.

const (
	// CPUFrequency is the 6502 CPU frequency in Hz
	CPUFrequency = 1023000

	// AudioSampleRate is the audio sample rate in Hz
	AudioSampleRate = 44100
)

var (
	// PendingInterrupt is set when an interrupt has just happened
	PendingInterrupt bool

	// PendingNMI is set when a non maskable interrupt has just happened
	PendingNMI bool

	// RunningTests is set when  tests are running
	RunningTests bool

	// RunningFunctionalTests is set when functional tests are running
	RunningFunctionalTests bool

	// RunningInterruptTests is set when interupt tests are running
	RunningInterruptTests bool

	// Cycles is the total CPU cycles executed
	Cycles uint64

	// FrameCycles is the CPU cycles executed in the current frame
	FrameCycles uint64

	// AudioChannel is the audio channel used to produce and consume audio samples between the cpu and audio packages
	AudioChannel chan int16

	// LastAudioValue is the + or - value of the current square wave
	LastAudioValue int16

	// LastAudioCycles is the ast CPU cycle when audio was sent to the channel
	LastAudioCycles uint64

	// AudioAttenuationCounter is a counter to keep track of when the audio should be zeroed after inactivity
	AudioAttenuationCounter uint64
)

// DriveState has the state of the disk drive
var DriveState struct {
	Drive        uint8 // What drive we're using. Currently only 1 is implemented
	Spinning     bool  // Is the motor spinning
	Phase        int8  // Phase of the stepper motor
	Phases       uint8 // the 4 lowest bits represent the 4 stepper motor magnet on/off states.
	BytePosition int   // Index of the position on the current track
	Q6           bool  // Q6 soft switch
	Q7           bool  // Q7 soft switch
}

// Init initializes the system-wide state
func Init() {
	Cycles = 0
	AudioChannel = make(chan int16, AudioSampleRate*4) // 1 second
	LastAudioValue = 0x2000
}

// WriteInterruptTestOpenCollector handles a write to a magic test address that triggers an interrupt and/or an NMI
func WriteInterruptTestOpenCollector(address uint16, oldValue uint8, value uint8) {
	oldInterrupt := (oldValue & 0x1) == 0x1
	oldNMI := (oldValue & 0x2) == 0x2

	interrupt := (value & 0x1) == 0x1
	NMI := (value & 0x2) == 0x2

	if oldInterrupt != interrupt {
		PendingInterrupt = interrupt
	}

	if oldNMI != NMI {
		PendingNMI = NMI
	}
}
