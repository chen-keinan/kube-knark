package trace

type Event_t struct {
	KtimeNs uint64
	Pid     uint32
	Uid     uint32
	Gid     uint32
	Type    int32
	Comm    [32]byte
}
