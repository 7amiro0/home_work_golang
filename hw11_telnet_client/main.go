package main

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"net"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

var timeOut time.Duration

func init() {
	pflag.DurationVar(&timeOut, "timeout", 10*time.Second, "time out for connect to server")
}

func validateArgs() (string, error) {
	args := pflag.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "not enougt args")
		return "", errors.New("not enougt args")
	}

	host, port := args[0], args[1]

	isHost := regexp.MustCompile(`(localhost)|((\b25[0-5]|\b2[0-4]\d|\b[01]?\d\d?)(\.(25[0-5]|2[0-4]\d|[01]?\d\d?)){3})$`)
	isPort := regexp.MustCompile(`^((6553[0-5])|(655[0-2]\d)|(65[0-4]\d{2})|(6[0-4]\d{3})|([1-5]\d{4})|([0-5]{0,5})|(\d{1,4}))$`)

	if !isHost.MatchString(host) || !isPort.MatchString(port) {
		_, _ = fmt.Fprintf(os.Stderr, "not correct host and port")
		return "", errors.New("not correct host and port")
	}

	return net.JoinHostPort(host, port), nil
}

func received(client TelnetClient) {
	for {
		err := client.Receive()

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error in received time: %q\n", err)
			return
		}
	}
}

func sender(sig chan os.Signal, client TelnetClient) {
	for {
		select {
		case <-sig:
			return
		default:
			err := client.Send()

			if err != nil {
				if err == io.EOF {
					_, _ = fmt.Fprint(os.Stderr, "EOF\n")
					return
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "Error in send time: %q\n", err)
					return
				}
			}
		}
	}
}

func sendReceive(sig chan os.Signal, client TelnetClient) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		received(client)
	}()

	go func() {
		sender(sig, client)
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	pflag.Parse()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	address, err := validateArgs()
	client := NewTelnetClient(address, timeOut, os.Stdin, os.Stdout)

	_, _ = fmt.Fprintf(os.Stderr, "Connected to %v\n", address)
	err = client.Connect()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error in connect time: %q\n", err)
		return
	}

	sendReceive(sig, client)
	fmt.Println("Bye, bye client")

	err = client.Close()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error while closing: %q\n", err)
		return
	}
}
