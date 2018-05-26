package system

const (
	CpuFrequency    = 1023000
	AudioSampleRate = 44100
)

var (
	PendingInterrupt        bool
	PendingNMI              bool
	RunningTests            bool
	RunningFunctionalTests  bool
	RunningInterruptTests   bool
	Cycles                  uint64
	FrameCycles             uint64
	AudioChannel            chan int16
	LastAudioValue          int16
	LastAudioCycles         uint64
	AudioAttenuationCounter uint64
)

var DriveState struct {
	Drive        uint8
	Spinning     bool
	Phase        int8
	Phases       uint8
	BytePosition int
	Q6           bool
	Q7           bool
}

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
