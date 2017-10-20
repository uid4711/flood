package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ufuchs/itplus/base/fcc"
	"github.com/ufuchs/zeroconf"
)

const (
	CONN_PORT    = 3333
	CONN_TYPE    = "tcp"
	SERVICE_NAME = "_flood._zyx._tcp"
)

//
//
//
func handleSignals(cancel context.CancelFunc, sigch <-chan os.Signal) {
	for {
		select {
		case <-sigch:
			fmt.Printf("\r")
			cancel()
		}
	}
}

//
//
//
func main() {

	var (
		err  error
		sigs = make(chan os.Signal, 2)
		//		mainWG      sync.WaitGroup
		hostname string
		server   *zeroconf.Server
	)

	if hostname, err = os.Hostname(); err != nil {
		fcc.Fatal(err)
	}

	if server, err = zeroconf.Register(hostname, SERVICE_NAME, "local.", CONN_PORT, []string{}, nil); err != nil {
		fcc.Fatalf("==> Zeroconf : Registering service '%v' failed - %v", SERVICE_NAME, err)
	}

	defer server.Shutdown()

	var hostAndPort = "192.168.178.107" + ":" + strconv.Itoa(CONN_PORT)

	l, err := net.Listen(CONN_TYPE, hostAndPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening on " + hostAndPort)

	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go handleSignals(cancel, sigs)
	go handleAccept(ctx, cancel, l)

	select {
	case <-ctx.Done():
		l.Close()
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("==> App stopped...")
	}

}

//
//
//
func handleAccept(ctx context.Context, cancel context.CancelFunc, l net.Listener) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Accept")
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				break
			}

			go handleRequest1(ctx, conn)
		}
	}
}

// Handles incoming requests.
func handleRequest(ctx context.Context, conn net.Conn) {

	fmt.Println("handle")

	s := NewChannelService()
	err := s.InitChannels()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {

		runtime.Gosched()

		//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		select {
		case <-ctx.Done():
			conn.Close()
			writeYML(s, "app.yml")
			fmt.Println("Done!")
			return
		default:

			for _, c := range s.Channels {

				// send the channel
				//				go func(ctx context.Context, cancel context.CancelFunc, c *Channel, conn net.Conn, wg *sync.WaitGroup) {

				_, err := conn.Write(c.GetRecord())

				if err != nil {
					conn.Close()
					writeYML(s, "app.yml")
					fmt.Println("Exit by error...")
					break
				}

				c.Success()

				//				time.Sleep(10 * time.Microsecond)

			}

		}

	}

}

// Handles incoming requests.
func handleRequest1(ctx context.Context, conn net.Conn) {

	fmt.Println("handle")

	var wg sync.WaitGroup

	s := NewChannelService()
	err := s.InitChannels()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctxL, cancelL := context.WithCancel(ctx)

	for {

		//runtime.Gosched()

		//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		select {
		case <-ctxL.Done():
			conn.Close()
			writeYML(s, "app.yml")
			fmt.Println("Done!!")
			return
		case <-ctx.Done():
			conn.Close()
			writeYML(s, "app.yml")
			fmt.Println("Done!")
			return
		default:

			for _, c := range s.Channels {

				wg.Add(1)

				// send the channel
				//				go func(ctx context.Context, cancel context.CancelFunc, c *Channel, conn net.Conn, wg *sync.WaitGroup) {
				go func(cancel context.CancelFunc, c *Channel, conn net.Conn, wg *sync.WaitGroup) {

					select {
					case <-ctxL.Done():
						wg.Done()
						return
					default:
						_, err := conn.Write(c.GetRecord())

						if err != nil {
							fmt.Println(err)
							wg.Done()
							cancelL()
							return
						}

						c.Success()

						wg.Done()

					}

				}(cancelL, c, conn, &wg)

				//time.Sleep(500 * time.Nanosecond)

			}
			wg.Wait()

		}

	}

}
