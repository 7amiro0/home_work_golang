package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	timeOut time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	address string
}

func (c *Client) Connect() (err error) {
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeOut)
	return err
}

func (c *Client) Close() (err error) {
	return c.conn.Close()
}

func (c Client) Send() (err error) {
	_, err = io.Copy(c.conn, c.in)
	return err
}

func (c Client) Receive() (err error) {
	_, err = io.Copy(c.out, c.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) (telnet TelnetClient) {
	return &Client{
		address: address,
		timeOut: timeout,
		in:      in,
		out:     out,
	}
}
