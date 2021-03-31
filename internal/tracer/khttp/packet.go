package khttp

import (
	"errors"
	"fmt"
	"github.com/chen-keinan/kube-knark/external"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/hsiafan/glow/flagx"
	"os"
	"runtime"
	"strconv"
	"time"
)

func listenOneSource(handle *pcap.Handle) chan gopacket.Packet {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	return packets
}

// set packet capture filter, by ip and port
func setDeviceFilter(handle *pcap.Handle, filterIP string, filterPort uint16) error {
	var bpfFilter = "tcp"
	if filterPort != 0 {
		bpfFilter += " port " + strconv.Itoa(int(filterPort))
	}
	if filterIP != "" {
		bpfFilter += " ip host " + filterIP
	}
	return handle.SetBPFFilter(bpfFilter)
}

// adapter multi channels to one channel. used to aggregate multi devices data
func mergeChannel(channels []chan gopacket.Packet) chan gopacket.Packet {
	var channel = make(chan gopacket.Packet)
	for _, ch := range channels {
		go func(c chan gopacket.Packet) {
			for packet := range c {
				channel <- packet
			}
		}(ch)
	}
	return channel
}

func openSingleDevice(device string, filterIP string, filterPort uint16) (localPackets chan gopacket.Packet, err error) {
	defer func() {
		if msg := recover(); msg != nil {
			switch x := msg.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			localPackets = nil
		}
	}()
	handle, err := pcap.OpenLive(device, 65536, false, pcap.BlockForever)
	if err != nil {
		return
	}

	if err := setDeviceFilter(handle, filterIP, filterPort); err != nil {
		fmt.Fprintln(os.Stderr, "set capture filter failed, ", err)
	}
	localPackets = listenOneSource(handle)
	return
}

//StartNetListener invoke net listener for kernel http events
func StartNetListener(errChan chan error, netEventChan chan *netevent.HTTPNetData) {
	go func() {
		var option = &external.Option{}
		cmd, err := flagx.NewCommand("httpdump", "capture and dump http contents", option, func() error {
			return run(option, errChan, netEventChan)
		})
		if err != nil {
			errChan <- fmt.Errorf("failed to get command for net listener")
		}
		os.Args = append(os.Args, "-level=all")
		cmd.ParseOsArgsAndExecute()
	}()
}

//nolint:gocyclo
func run(option *external.Option, errChan chan error, netEventChan chan *netevent.HTTPNetData) error {
	if option.Port > 65536 {
		return fmt.Errorf("ignored invalid port %v", option.Port)
	}

	if option.Status != "" {
		statusSet, err := external.ParseIntSet(option.Status)
		if err != nil {
			errChan <- fmt.Errorf("status range not valid %v", option.Status)
		}
		option.StatusSet = statusSet
	}

	var packets chan gopacket.Packet
	if option.File != "" {
		//TODO: read file stdin
		// read from pcap file
		var handle, err = pcap.OpenOffline(option.File)
		if err != nil {
			errChan <- fmt.Errorf("open file %v error: %w", option.File, err)
		}
		packets = listenOneSource(handle)
	} else if option.Device == "any" && runtime.GOOS != "linux" {
		// capture all device
		// Only linux 2.2+ support any interface. we have to list all network device and listened on them all
		interfaces, err := pcap.FindAllDevs()
		if err != nil {
			errChan <- fmt.Errorf("find device error: %w", err)
		}

		var packetsSlice = make([]chan gopacket.Packet, len(interfaces))
		for _, itf := range interfaces {
			localPackets, err := openSingleDevice(itf.Name, option.Ip, uint16(option.Port))
			if err != nil {
				continue
			}
			packetsSlice = append(packetsSlice, localPackets)
		}
		packets = mergeChannel(packetsSlice)
	} else if option.Device != "" {
		// capture one device
		var err error
		packets, err = openSingleDevice(option.Device, option.Ip, uint16(option.Port))
		if err != nil {
			errChan <- fmt.Errorf("listen on device %v failed, error: %w", option.Device, err)
		}
	} else {
		errChan <- errors.New("no device or pcap file specified")
	}

	var handler = &external.HTTPConnectionHandler{
		Option: option,
		// TODO: stdout
		Printer: external.NewPrinter(netEventChan),
	}
	var assembler = external.NewTCPAssembler(handler)
	assembler.FilterIP = option.Ip
	assembler.FilterPort = uint16(option.Port)
	var ticker = time.NewTicker(time.Second * 10).C

outer:
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				break outer
			}
			// only assembly tcp/ip packets
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil ||
				packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}
			var tcp = packet.TransportLayer().(*layers.TCP)

			assembler.Assemble(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// flush connections that haven't been activity in the idle time
			assembler.FlushOlderThan(time.Now().Add(-option.Idle))
		}
	}
	assembler.FinishAll()
	external.WaitGroup.Wait()
	handler.Printer.Finish()
	external.PrinterWaitGroup.Wait()
	return nil
}
