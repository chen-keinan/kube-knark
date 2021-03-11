package trace

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"time"
	"unsafe"
)

var (
	sizeofExecveData = int(unsafe.Sizeof(ExecveData{}))
	sizeofExecveArg  = int(unsafe.Sizeof(ExecveArg{}))
	sizeofExecveRtn  = int(unsafe.Sizeof(ExecveRtn{}))
	sizeofExitData   = int(unsafe.Sizeof(ExitData{}))
)

type ExecveData struct {
	KTimeNS         time.Duration
	RealStartTimeNS time.Duration
	PID             int
	UID             int
	GID             int
	PPID            int
	Comm            [16]byte
}

func (e ExecveData) String() string {
	return fmt.Sprintf("ktime:%d, real_start_time:%d, pid:%d, uid:%d, gid:%d, ppid:%d, comm:%s",
		e.KTimeNS, e.RealStartTimeNS, e.PID, e.UID, e.GID, e.PPID, utils.NullTerminatedString(e.Comm[:]))
}

type ExecveArg struct {
	PID int
	_   int
	Arg [256]byte
}

func (e ExecveArg) String() string {
	return fmt.Sprintf("pid:%d, arg:%s", e.PID, utils.NullTerminatedString(e.Arg[:]))
}

type ExecveRtn struct {
	PID        int
	ReturnCode int32
}

func (e ExecveRtn) String() string {
	return fmt.Sprintf("pid:%d, rtn:%d", e.PID, e.ReturnCode)
}

type ExitData struct {
	KTime uint64
	PID   int
}

func (e ExitData) String() string {
	return fmt.Sprintf("ktime:%d, pid:%d", e.KTime, e.PID)
}

func unmarshalData(data []byte) (ExecveData, error) {
	var event ExecveData
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event)
	return event, err
}

func unmarshalArg(data []byte) (ExecveArg, error) {
	var event ExecveArg
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event)
	return event, err
}

func unmarshalRtn(data []byte) (ExecveRtn, error) {
	var event ExecveRtn
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event)
	return event, err
}

func unmarshalExitData(data []byte) (ExitData, error) {
	var event ExitData
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event)
	return event, err
}
