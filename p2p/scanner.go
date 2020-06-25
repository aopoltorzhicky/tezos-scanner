package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/aopoltorzhicky/tezos-scanner/p2p/ffi"
	"github.com/aopoltorzhicky/tezos-scanner/p2p/protocol"
)

// Scanner -
type Scanner struct {
	bootstrap []net.IP

	result     chan *protocol.Peer
	candidates chan *Node
	stop       chan struct{}

	dropAfter    int64
	threadsCount int64
	syncedTime   int64

	stopped bool

	attemptsDuration time.Duration
	identity         ffi.Identity

	proofedPeers sync.Map
	mutex        sync.Mutex
	wg           sync.WaitGroup
}

// IsStopped -
func (scanner *Scanner) IsStopped() bool {
	return scanner.stopped
}

func prepareBootstrap(bootstrap []string) ([]net.IP, error) {
	ips := make([]net.IP, 0)
	for _, address := range bootstrap {
		potentialIps, err := net.LookupIP(address)
		if err != nil {
			log.Printf("Lookup ip failed: %s", err)
			continue
		}

		ips = append(ips, potentialIps...)
	}
	return ips, nil
}

// NewScanner -
func NewScanner(bootstrap []string, identity ffi.Identity, opts ...ScannerOption) (*Scanner, error) {
	ips, err := prepareBootstrap(bootstrap)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("Empty bootstrap array")
	}
	scanner := &Scanner{
		bootstrap:  ips,
		identity:   identity,
		result:     make(chan *protocol.Peer, 1024),
		candidates: make(chan *Node, 1024),
		stopped:    true,
	}
	for _, opt := range opts {
		opt(scanner)
	}

	if scanner.syncedTime == 0 {
		scanner.syncedTime = 120
	}

	if scanner.threadsCount == 0 {
		scanner.threadsCount = 4
	}
	scanner.stop = make(chan struct{}, scanner.threadsCount)

	if scanner.attemptsDuration.Seconds() == 0 {
		scanner.attemptsDuration = 300 * time.Second
	}
	return scanner, nil
}

// Scan -
func (scanner *Scanner) Scan() {
	for _, ip := range scanner.bootstrap {
		peer := &protocol.Peer{
			Address: net.TCPAddr{
				Port: 9732,
				IP:   ip,
			},
		}
		scanner.candidates <- NewNode(peer, scanner.attemptsDuration, scanner.dropAfter, scanner.syncedTime)

	}

	for i := int64(0); i < scanner.threadsCount; i++ {
		scanner.wg.Add(1)
		go scanner.process()
	}

	scanner.stopped = false
}

// Listen -
func (scanner *Scanner) Listen() chan *protocol.Peer {
	return scanner.result
}

// Stop -
func (scanner *Scanner) Stop() {
	scanner.stopped = true
	log.Print("Stopping scanner...")
	for i := int64(0); i < scanner.threadsCount; i++ {
		scanner.stop <- struct{}{}
	}
	close(scanner.candidates)
	scanner.wg.Wait()

	close(scanner.result)
	close(scanner.stop)
}

func (scanner *Scanner) process() {
	defer scanner.wg.Done()

	for {
		select {
		case <-scanner.stop:
			log.Print("Thread has stopped")
			return

		case candidate := <-scanner.candidates:
			if err := scanner.processCandidate(candidate); err != nil {
				log.Printf("Error during process peer: %s", err)
			}
		}
	}
}

func (scanner *Scanner) processCandidate(candidate *Node) error {
	defer candidate.close()

	candidateIP := candidate.Peer.Address.IP.String()
	log.Printf("Scan %s. Attempt %d", candidateIP, candidate.attemptsCount)
	if _, ok := scanner.proofedPeers.Load(candidateIP); ok {
		log.Printf("Peer %s has already scanned", candidateIP)
		return nil
	}

	peer, err := candidate.getPeer(scanner.identity)
	if err != nil {
		if scanner.stopped {
			return nil
		}

		scanner.mutex.Lock()
		defer scanner.mutex.Unlock()

		if candidate.hasAttempts() {
			scanner.candidates <- candidate
			return nil
		}

		log.Printf("[getPeer] %s", err)
		candidate.setErrorState(err)
		scanner.proofedPeers.Store(candidateIP, candidate.Peer)
		scanner.result <- candidate.Peer
		return err
	}

	neighbors, err := peer.GetPeersAddresses()
	if err != nil {
		log.Printf("[GetPeersAddresses] %s", err)
		return err
	}

	if scanner.stopped {
		return nil
	}

	scanner.mutex.Lock()
	defer scanner.mutex.Unlock()

	scanner.proofedPeers.Store(peer.Address.IP.String(), peer)
	scanner.result <- peer
	for _, newPeer := range neighbors {
		scanner.candidates <- NewNode(newPeer, scanner.attemptsDuration, scanner.dropAfter, scanner.syncedTime)
	}
	return nil
}
