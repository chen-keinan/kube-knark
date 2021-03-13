package trace

import (
	"bytes"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	bpf "github.com/iovisor/gobpf/elf"
	"github.com/pkg/errors"
	"github.com/prometheus/procfs"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"sync"
	"time"
)

// clockHz is the clock tick rate (/proc/PID/stat gives start time in ticks).
const clockHz = 100

// probe and map names in the eBPF program.
const (
	execveProbe       = "kprobe/SyS_execve"
	execveReturnProbe = "kretprobe/SyS_execve"
	execveMap         = "execve_events"
	doExitProbe       = "kprobe/do_exit"
)

// Process Monitor

type ProcessMonitor struct {
	// eBPF
	module        *bpf.Module
	execvePerfMap *bpf.PerfMap
	bpfEvents     chan []byte
	lostBPFEvents chan uint64
	lostCount     uint64

	// internal state
	bootTime     time.Time
	procfs       procfs.FS
	processTable map[int]*process
	warnOnce     sync.Once

	output chan interface{}
	done   <-chan struct{}
}

type eventSource int

const (
	sourceBPF eventSource = iota + 1
	sourceProcFS
)

type processState int

const (
	stateStarted processState = iota + 1
	stateError
	stateExited
)

type process struct {
	processData
	State     processState
	Source    eventSource
	EndTime   time.Time
	ErrorCode int32
}

func NewMonitor() (*ProcessMonitor, error) {
	procfs, err := procfs.NewFS(procfs.DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	// Fetch and cache the boot time.
	stat, err := procfs.NewStat()
	if err != nil {
		return nil, err
	}

	return &ProcessMonitor{
		bootTime:     time.Unix(int64(stat.BootTime), 0),
		procfs:       procfs,
		processTable: map[int]*process{},
	}, nil
}

func (m *ProcessMonitor) Start(done <-chan struct{}) (<-chan interface{}, error) {
	if err := m.initBPF(); err != nil {
		return nil, err
	}
	m.output = make(chan interface{}, 1)

	go func() {
		defer close(m.output)
		defer m.execvePerfMap.PollStop()
		defer m.module.Close()

		allProcs, err := m.readProcs()
		if err != nil {
		} else {
			for _, p := range allProcs {
				m.processTable[p.PID] = p
				m.publish(p)
			}
		}

		for {
			select {
			case data := <-m.bpfEvents:
				m.handleBPFData(data)
			case count := <-m.lostBPFEvents:
				m.lostCount += count

			case <-done:
				return
			}
		}
	}()

	return m.output, nil
}

func (m *ProcessMonitor) initBPF() error {
	f, err := os.Open("exec.o")
	if err != nil {
		return errors.Wrap(err, "failed to open elf file")
	}
	defer f.Close()
	// Load module to kernel.

	m.module = bpf.NewModuleFromReader(f)

	err = m.module.Load(nil)

	if err != nil {
		return errors.Wrap(err, "failed to load ebpf module to kernel")
	}

	// Setup our perf event readers.
	m.bpfEvents = make(chan []byte, 64)
	m.lostBPFEvents = make(chan uint64, 1)
	m.execvePerfMap, err = bpf.InitPerfMap(m.module, execveMap, m.bpfEvents, m.lostBPFEvents)
	if err != nil {
		m.module.Close()
		return errors.Wrapf(err, "failed to initialize %v perf map", execveMap)
	}

	// Enable the kprobes.
	if err := m.module.EnableKprobe(execveProbe, 0); err != nil {
		m.module.Close()
		return errors.Wrapf(err, "failed to enable %v probe", execveProbe)
	}

	if err := m.module.EnableKprobe(execveReturnProbe, 0); err != nil {
		m.module.Close()
		return errors.Wrapf(err, "failed to enable %v probe", execveReturnProbe)
	}

	if err := m.module.EnableKprobe(doExitProbe, 0); err != nil {
		m.module.Close()
		return errors.Wrapf(err, "failed to enable %v probe", doExitProbe)
	}

	m.execvePerfMap.PollStart()
	return nil
}

func (m *ProcessMonitor) handleBPFData(data []byte) {
	switch len(data) {
	case sizeofExecveData:
		event, err := unmarshalData(data)
		if err != nil {
			return
		}

		// Process already exists in the table.
		if _, exists := m.processTable[event.PID]; exists {
			return
		}

		// Sanity check the RealStartTimeNS value.
		if absDuration(event.RealStartTimeNS-event.KTimeNS) > 10*time.Second {
			event.RealStartTimeNS = event.KTimeNS

			// task_struct data is probably garbage so fall-back to /proc/PID
			status, err := m.procStatus(event.PID)
			if err == nil {
				event.PPID = status.PID
			}

			m.warnOnce.Do(func() {

			})
		}

		m.processTable[event.PID] = &process{
			State:  stateStarted,
			Source: sourceBPF,
			processData: processData{
				StartTime:  m.bootTime.Add(event.RealStartTimeNS),
				PPID:       event.PPID,
				ParentComm: utils.NullTerminatedString(event.Comm[:]),
				PID:        event.PID,
				UID:        event.UID,
				GID:        event.GID,
			},
		}
	case sizeofExecveArg:
		event, err := unmarshalArg(data)
		if err != nil {
			return
		}

		p, found := m.processTable[event.PID]
		if !found || p.Source == sourceProcFS {
			return
		}

		// The first argument sent is the exe.
		arg := utils.NullTerminatedString(event.Arg[:])
		if len(p.Exe) == 0 {
			p.Exe = arg
			return
		}
		p.Args = append(p.Args, arg)
	case sizeofExecveRtn:
		event, err := unmarshalRtn(data)
		if err != nil {
			return
		}

		p, found := m.processTable[event.PID]
		if !found || p.Source == sourceProcFS {
			return
		}
		if event.ReturnCode != 0 {
			p.State = stateError
			p.ErrorCode = event.ReturnCode
		}

		m.publish(p)
	case sizeofExitData:
		event, err := unmarshalExitData(data)
		if err != nil {
			return
		}

		p, found := m.processTable[event.PID]
		if !found || p.ErrorCode != 0 {
			return
		}
		p.State = stateExited
		p.EndTime = m.bootTime.Add(time.Duration(event.KTime))
		delete(m.processTable, event.PID)
		m.publish(p)
	}
}

func (m *ProcessMonitor) publish(p *process) {
	var event interface{}
	switch p.State {
	case stateStarted:
		event = ProcessStarted{
			Type:        "started",
			processData: p.processData,
		}
	case stateExited:
		event = ProcessExited{
			Type:        "exited",
			processData: p.processData,
			EndTime:     p.EndTime,
			RunningTime: p.EndTime.Sub(p.StartTime),
		}
	case stateError:
		event = ProcessError{
			Type:        "error",
			processData: p.processData,
			ErrorCode:   p.ErrorCode,
		}
	default:
		return
	}

	// Output event.
	select {
	case <-m.done:
	case m.output <- event:
	}
}

func (m *ProcessMonitor) procStatus(pid int) (*procfs.ProcStatus, error) {
	p, err := m.procfs.NewProc(int(pid))
	if err != nil {
		return nil, err
	}

	status, err := p.NewStatus()
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (m *ProcessMonitor) readProcs() ([]*process, error) {
	procs, err := m.procfs.AllProcs()
	if err != nil {
		return nil, err
	}

	out := make([]*process, 0, len(procs))
	for _, p := range procs {
		process, err := readProc(p, m.bootTime)
		if err != nil {

			continue
		}
		out = append(out, process)
	}

	return out, nil
}

func readProc(p procfs.Proc, bootTime time.Time) (*process, error) {
	stat, err := p.Stat()
	if err != nil {
		return nil, err
	}

	status, err := p.NewStatus()
	if err != nil {
		return nil, err
	}

	args, err := p.CmdLine()
	if err != nil {
		return nil, err
	}

	exe, err := p.Executable()
	if err != nil {
		return nil, err
	}
	var uid int
	if len(status.UIDs[0]) > 0 {
		uid, err = strconv.Atoi(status.UIDs[0])
		if err != nil {
			uid = 0
		}
	}
	var gid int
	if len(status.GIDs[0]) > 0 {
		gid, err = strconv.Atoi(status.GIDs[0])
		if err != nil {
			gid = 0
		}
	}
	process := &process{
		State:  stateStarted,
		Source: sourceProcFS,
		processData: processData{
			StartTime: bootTime.Add(ticksToNanos(stat.Starttime)),
			PPID:      status.PID,
			PID:       stat.PPID,
			UID:       uid,
			GID:       gid,
			Comm:      status.Name,
			Exe:       exe,
			Args:      args,
		},
	}

	return process, nil
}

func ticksToNanos(ticks uint64) time.Duration {
	return time.Duration(ticks) * time.Second / clockHz
}

func absDuration(v time.Duration) time.Duration {
	if v < 0 {
		return -v
	}
	return v
}
