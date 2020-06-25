package main

import (
	"encoding/csv"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/aopoltorzhicky/tezos-scanner/p2p"
	"github.com/aopoltorzhicky/tezos-scanner/p2p/ffi"
)

func main() {
	cfg, err := getConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	identity, err := ffi.GetIdentity()
	if err != nil {
		panic(err)
	}
	scanner, err := p2p.NewScanner(
		cfg.Bootstrap,
		identity,
		p2p.WithAttemptsDuration(cfg.Attempts.Timeout),
		p2p.WithDropAfter(cfg.Attempts.Count),
		p2p.WithSyncedTime(cfg.SyncedTime),
		p2p.WithThreadsCount(cfg.ThreadsCount),
	)
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-stop
		log.Print("[!] Ctrl+C pressed in Terminal")
		scanner.Stop()
	}()

	header := [][]string{
		{"ip", "port", "id", "private", "rpc", "synced", "disable_mempool", "neighbors", "versions"},
	}

	recordFile, err := os.Create("peers.csv")
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(recordFile)

	// write header
	if err = writer.WriteAll(header); err != nil {
		panic(err)
	}

	scanner.Scan()
	for {
		if scanner.IsStopped() {
			break
		}
		select {
		case peer := <-scanner.Listen():
			if peer == nil {
				continue
			}
			log.Printf("Found peer: %s", peer.Address.IP)
			versions := make([]string, len(peer.Versions))
			for i := range peer.Versions {
				versions[i] = peer.Versions[i].Name
			}

			neighbors := make([]string, len(peer.Neighbors))
			for i := range peer.Neighbors {
				host, _, err := net.SplitHostPort(peer.Neighbors[i])
				if err != nil {
					panic(err)
				}
				neighbors[i] = host
			}

			row := [][]string{
				{
					peer.Address.IP.String(),
					strconv.Itoa(peer.Address.Port),
					peer.ID,
					strconv.FormatBool(peer.PrivateNode),
					strconv.FormatBool(peer.RPC),
					strconv.FormatBool(peer.Synced),
					strconv.FormatBool(peer.DisableMempool),
					strings.Join(neighbors, "|"),
					strings.Join(versions, "|"),
				},
			}
			if err = writer.WriteAll(row); err != nil {
				panic(err)
			}
		}
	}

	log.Print("Stopped")
	close(stop)
}
