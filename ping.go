package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	ProtocolICMP = 1
)

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readAnswer(c *icmp.PacketConn, id int, peer_from net.Addr, send_time time.Time, reply chan int) {
	var (
		n    int
		peer net.Addr
		err  error
		hdr  *ipv4.Header
		rb   []byte
	)
	for {
		rb = make([]byte, 1500)
		n, peer, err = c.ReadFrom(rb)
		check_err(err)
		if peer.String() != peer_from.String() {
			continue
		}
		hdr, err = icmp.ParseIPv4Header(rb)
		check_err(err)
		if hdr.ID == id {
			break
		}
	}
	rm, err := icmp.ParseMessage(ProtocolICMP, rb[:n])
	check_err(err)
	duration := time.Now().Sub(send_time)
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("%d bytes from %s: id=%d icmp_seq=%d ttl=%d time=%.1f ms\n", n, peer, hdr.ID, hdr.FragOff, hdr.TTL, duration.Seconds()*1000)
		reply <- 1
	default:
		fmt.Printf("got %+v; want echo reply\n", rm)
		reply <- 0
	}
}

func main() {
	var (
		sended, recieved, loss, seq, exit_code int
		send_time                              time.Time
	)
	id, err := strconv.Atoi(os.Args[2])
	check_err(err)
	addr := &net.IPAddr{IP: net.ParseIP(os.Args[1])}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	check_err(err)
	defer c.Close()

	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, os.Interrupt)
	tick_chan := time.Tick(time.Second)
	reply_chan := make(chan int)

Loop:
	for {
		select {
		case <-sig_chan:
			break Loop
		case <-tick_chan:
			seq++
			wm := icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{
					ID:   id,
					Seq:  seq,
					Data: []byte("HELLO-R-U-THERE"),
				},
			}
			wb, err := wm.Marshal(nil)
			check_err(err)
			_, err = c.WriteTo(wb, addr)
			check_err(err)
			send_time = time.Now()
			sended++
			go readAnswer(c, id, addr, send_time, reply_chan)
		case success_count := <-reply_chan:
			recieved += success_count
		}
	}
	if sended > recieved {
		exit_code = 1
	}
	fmt.Println("")
	fmt.Printf("--- %s ping statistics ---\n", addr)
	if sended > 0 {
		loss = (sended - recieved) * 100 / sended
	}
	fmt.Printf("%d packets transmitted, %d received, %d%% packet loss\n", sended, recieved, loss)
	os.Exit(exit_code)
}
