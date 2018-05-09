package system

var (
	PendingInterrupt bool
	PendingNMI       bool

	RunningTests           bool
	RunningFunctionalTests bool
	RunningInterruptTests  bool
)

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
