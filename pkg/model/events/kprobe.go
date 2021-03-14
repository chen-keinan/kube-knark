package events

//KprobeEvent define external kprobe event
type KprobeEvent struct {
	StartTime string
	Pid       uint32
	UID       uint32
	Gid       uint32
	Comm      string
	Args      []string
}
