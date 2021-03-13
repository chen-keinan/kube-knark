package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/dropbox/goebpf"
)

var (
	ErrProgramNotFound = errors.New("program not found")
	ErrMapNotFound     = errors.New("map not found")
)

type Event_t struct {
	KtimeNs uint64   `json:"start_time,omitempty""`
	Pid     uint32   `json:"pid,omitempty"`
	Uid     uint32   `json:"uid,omitempty"`
	Gid     uint32   `json:"gid,omitempty"`
	Type    int32    `json:"type,omitempty"`
	Comm    [32]byte `json:"comm,omitempty"`
}

type Program struct {
	bpf goebpf.System
	pe  *goebpf.PerfEvents
	wg  sync.WaitGroup
}

func StartKnark() {
	// cleanup old probes
	if err := goebpf.CleanupProbes(); err != nil {
		log.Println(err)
	}

	// load ebpf program
	p, err := LoadProgram("kprobe.elf")
	if err != nil {
		log.Fatalf("LoadProgram() failed: %v", err)
	}
	p.ShowInfo()

	// attach ebpf kprobes
	if err := p.AttachProbes(); err != nil {
		log.Fatalf("AttachProbes() failed: %v", err)
	}
	defer p.DetachProbes()

	// wait until Ctrl+C pressed
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	<-ctrlC

	// display some stats
	fmt.Println()
	fmt.Printf("%d Event(s) Received\n", p.pe.EventsReceived)
	fmt.Printf("%d Event(s) lost (e.g. small buffer, delays in processing)\n", p.pe.EventsLost)
}

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

func (p *Program) startPerfEvents(events <-chan []byte) {
	p.wg.Add(1)
	go func(events <-chan []byte) {
		defer p.wg.Done()

		// print header
		fmt.Printf("\nTIME          PCOMM             PID    UID    GID    DESC\n\n")
		for {

			// receive exec events
			if b, ok := <-events; ok {

				// parse proc info
				var ev Event_t
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
				var desc string
				if len(args) > 0 {
					desc = args[0]
				}
				if len(args) > 2 {
					desc += " " + strings.Join(args[2:], " ")
				}
				// display process execution event
  				buff := new(bytes.Buffer)
				json.NewEncoder(buff).Encode(&ev)
				fmt.Println(buff.String())
			} else {
				break
			}
		}
	}(events)
}

func (p *Program) stopPerfEvents() {
	p.pe.Stop()
	p.wg.Wait()
}

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

func (p *Program) DetachProbes() error {
	p.stopPerfEvents()
	for _, prog := range p.bpf.GetPrograms() {
		prog.Detach()
		prog.Close()
	}
	return nil
}

func (p *Program) ShowInfo() {
	fmt.Println()
	fmt.Println("Maps:")
	for _, item := range p.bpf.GetMaps() {
		m := item.(*goebpf.EbpfMap)
		fmt.Printf("\t%s: %v, Fd %v\n", m.Name, m.Type, m.GetFd())
	}
	fmt.Println("\nPrograms:")
	for _, prog := range p.bpf.GetPrograms() {
		fmt.Printf("\t%s: %v (%s), size %d, license \"%s\"\n",
			prog.GetName(), prog.GetType(), prog.GetSection(), prog.GetSize(), prog.GetLicense(),
		)
	}
}
