package workers

import "fmt"

//PacketMatches instance which match packet data to specific pattern
type PacketMatches struct {
	numOfWorkers int
	pmc          chan string
}

//NewMatches return new packet instance
func NewPacketMatches(numOfWorkers int, pmc chan string) *PacketMatches {
	return &PacketMatches{numOfWorkers: numOfWorkers, pmc: pmc}
}

//Invoke invoke packet matches workers
func (pm *PacketMatches) Invoke() {
	for i := 0; i < pm.numOfWorkers; i++ {
		go func() {
			for k := range pm.pmc {
				fmt.Print(k)
			}
		}()
	}
}
