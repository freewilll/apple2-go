package system

// The system package is a dumping ground for globals that are shared between
// the packages.

const (
	CpuFrequency    = 1023000 // 6402 CPU frequency in Hz
	AudioSampleRate = 44100   // Audio sample rate in Hz
)

var (
	PendingInterrupt        bool       // Set when an interrupt has just happened
	PendingNMI              bool       // Set when a non maskable interrupt has just happened
	RunningTests            bool       // For testing
	RunningFunctionalTests  bool       // For testing
	RunningInterruptTests   bool       // For testing
	Cycles                  uint64     // Total CPU cycles executed
	FrameCycles             uint64     // CPU cycles executed in the current frame
	AudioChannel            chan int16 // Audio channel
	LastAudioValue          int16      // + or - value of the current square wave
	LastAudioCycles         uint64     // Last CPU cycle when audio was sent to the channel
	AudioAttenuationCounter uint64     // Counter to keep track of when the audio should be zeroed after inactivity
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

// Handle a write to a magic test address that triggers an interrupt and/or an NMI
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
