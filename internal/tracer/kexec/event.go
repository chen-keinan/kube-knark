package kexec

//EventKprobe Object
type EventKprobe struct {
	KtimeNs uint64
	Pid     uint32
	UID     uint32
	Gid     uint32
	Type    int32
	Comm    [32]byte
}
