package trace

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/events"
	"github.com/dropbox/goebpf"
	"sync"
)

var (
	//ErrProgramNotFound program not found error
	ErrProgramNotFound = errors.New("program not found")
	//ErrMapNotFound map not found error
	ErrMapNotFound = errors.New("map not found")
)

//Program object
type Program struct {
	bpf goebpf.System
	pe  *goebpf.PerfEvents
	wg  sync.WaitGroup
}

//LoadProgram load ebpf program
func LoadProgram(filename string) (*Program, error) {
	// create system
	bpf := goebpf.NewDefaultEbpfSystem()
	// load compiled ebpf elf file
	if err := bpf.LoadElf(filename); err != nil {
		return nil, err
	}
	// load programs
	for _, prog := range bpf.GetPrograms() {
		if err := prog.Load(); err != nil {
			return nil, err
		}
	}
	return &Program{bpf: bpf}, nil
}

//startPerfEvents pull ebpf events
func (p *Program) startPerfEvents(kevents <-chan []byte) {
	p.wg.Add(1)
	go func(kevents <-chan []byte) {
		defer p.wg.Done()

		// print header
		fmt.Printf("\nTIME          PCOMM             PID    UID    GID    DESC\n\n")
		for {

			// receive exec events
			if b, ok := <-kevents; ok {

				// parse proc info
				var ev EventKprobe
				buf := bytes.NewBuffer(b)
				if err := binary.Read(buf, binary.LittleEndian, &ev); err != nil {
					fmt.Printf("error: %v\n", err)
					continue
				}

				// parse args
				tokens := bytes.Split(buf.Bytes(), []byte{0x00})
				var args []string
				for _, arg := range tokens {
					if len(arg) > 0 {
						args = append(args, string(arg))
					}
				}

				// build display strings
				ts := goebpf.KtimeToTime(ev.KtimeNs)
				comm := goebpf.NullTerminatedStringToString(ev.Comm[:])
				ke := &events.KprobeEvent{
					StartTime: ts.Format("15:04:05.000"),
					Comm:      comm,
					Pid:       ev.Pid,
					Uid:       ev.UID,
					Gid:       ev.Gid,
					Args:      args,
				}
				// display process execution event
				kwriter := new(bytes.Buffer)
				err := json.NewEncoder(kwriter).Encode(&ke)
				if err != nil {
					continue
				}
				fmt.Println(kwriter.String())
			} else {
				break
			}
		}
	}(kevents)
}

func (p *Program) stopPerfEvents() {
	p.pe.Stop()
	p.wg.Wait()
}

//AttachProbes attach ebpf program to kernel
func (p *Program) AttachProbes() error {
	// attach all probe programs
	for _, prog := range p.bpf.GetPrograms() {
		if err := prog.Attach(nil); err != nil {
			return err
		}
	}

	// get handles to perf event map
	m := p.bpf.GetMapByName("events")
	if m == nil {
		return ErrMapNotFound
	}

	// create perf events
	var err error
	p.pe, err = goebpf.NewPerfEvents(m)
	if err != nil {
		return err
	}
	events, err := p.pe.StartForAllProcessesAndCPUs(4096)
	if err != nil {
		return err
	}

	// start event listeners
	p.wg = sync.WaitGroup{}
	p.startPerfEvents(events)

	return nil
}

//DetachProbes detach ebpf program from kernel
func (p *Program) DetachProbes() error {
	p.stopPerfEvents()
	for _, prog := range p.bpf.GetPrograms() {
		err := prog.Detach()
		if err != nil {
			return err
		}
		err = prog.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//ShowInfo how kprobe program info
func (p *Program) ShowInfo() {
	for _, item := range p.bpf.GetMaps() {
		m, ok := item.(*goebpf.EbpfMap)
		if ok {
			fmt.Printf("\t%s: %v, Fd %v\n", m.Name, m.Type, m.GetFd())
		}
	}
	fmt.Println("\nPrograms:")
	for _, prog := range p.bpf.GetPrograms() {
		fmt.Printf("\t%s: %v (%s), size %d, license \"%s\"\n",
			prog.GetName(), prog.GetType(), prog.GetSection(), prog.GetSize(), prog.GetLicense(),
		)
	}
}

//EventsLost return num of ebpf event lost
func (p *Program) EventsLost() int {
	return p.pe.EventsLost
}

//EventsReceived return num of ebpf events received
func (p *Program) EventsReceived() int {
	return p.pe.EventsReceived
}
