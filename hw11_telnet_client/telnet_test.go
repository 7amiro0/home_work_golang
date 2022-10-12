package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("client closed twice", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:8000")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
			require.NoError(t, client.Connect())

			require.NoError(t, client.Close())
			require.Error(t, client.Close(), "close tcp 127.0.0.1:49232->127.0.0.1:8000: use of closed network connection")
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()
		}()

		wg.Wait()
	})

	t.Run("send big text", func(t *testing.T) {
		loremText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Eget nullam non nisi est sit. Urna id volutpat lacus laoreet non curabitur gravida arcu ac. Sapien eget mi proin sed libero enim sed faucibus turpis. Aenean pharetra magna ac placerat vestibulum lectus mauris ultrices eros. At risus viverra adipiscing at in. Erat imperdiet sed euismod nisi porta lorem mollis aliquam. Suscipit adipiscing bibendum est ultricies integer quis auctor. At erat pellentesque adipiscing commodo elit at. Orci dapibus ultrices in iaculis nunc. Nunc scelerisque viverra mauris in aliquam sem fringilla ut. Egestas dui id ornare arcu odio ut sem nulla pharetra. Lectus vestibulum mattis ullamcorper velit. Consequat ac felis donec et odio pellentesque. Mi quis hendrerit dolor magna eget est lorem.\n\nOdio eu feugiat pretium nibh ipsum consequat nisl vel pretium. Lacus sed viverra tellus in hac habitasse. Dolor sit amet consectetur adipiscing elit duis tristique"
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), &bytes.Buffer{})
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString(loremText)
			err = client.Send()
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, loremText, string(request)[:n])
		}()

		wg.Wait()
	})

	t.Run("received big text", func(t *testing.T) {
		loremText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Eget nullam non nisi est sit. Urna id volutpat lacus laoreet non curabitur gravida arcu ac. Sapien eget mi proin sed libero enim sed faucibus turpis. Aenean pharetra magna ac placerat vestibulum lectus mauris ultrices eros. At risus viverra adipiscing at in. Erat imperdiet sed euismod nisi porta lorem mollis aliquam. Suscipit adipiscing bibendum est ultricies integer quis auctor. At erat pellentesque adipiscing commodo elit at. Orci dapibus ultrices in iaculis nunc. Nunc scelerisque viverra mauris in aliquam sem fringilla ut. Egestas dui id ornare arcu odio ut sem nulla pharetra. Lectus vestibulum mattis ullamcorper velit. Consequat ac felis donec et odio pellentesque. Mi quis hendrerit dolor magna eget est lorem.\n\nOdio eu feugiat pretium nibh ipsum consequat nisl vel pretium. Lacus sed viverra tellus in hac habitasse. Dolor sit amet consectetur adipiscing elit duis tristique"
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(&bytes.Buffer{}), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, loremText, out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			n, err := conn.Write([]byte(loremText))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("connect to not exist server", func(t *testing.T) {
		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("127.0.0.1:8000", timeout, ioutil.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
		err = client.Connect()
		require.Errorf(t, err, "expected error failed connect but received %q", err)
	})

	t.Run("closed_twice", func(t *testing.T) {
		l, err := net.Listen("tcp", "0.0.0.0:3302")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, l.Close())
		}()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:3302", 10*time.Second, io.NopCloser(in), out)

		require.NoError(t, client.Connect())
		require.NoError(t, client.Connect())

		require.NoError(t, client.Close())
		require.NoError(t, client.Close())
	})
}
