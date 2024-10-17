package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

func (c *Client) intToByte(n uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		c.ErrorCh <- errors.New("error in convert int to []byte")
	}
	return buf.Bytes()
}

type Client struct {
	TCPConn          net.Conn
	UDPConn          net.Conn
	TCPReq, UDPReq   chan uint32
	TCPResp, UDPResp chan uint32
	ErrorCh          chan error
}

func NewClient(laddr, raddr string, lport, rport int, timeout time.Duration) (*Client, error) {
	// server ip
	remoteIP := net.ParseIP(raddr)
	// client ip
	localIP := net.ParseIP(laddr)

	// create tcp local addr, remote addr also create its dialer
	tcpLocalAddr := &net.TCPAddr{
		IP:   localIP,
		Port: lport,
	}
	tcpRemoteAddr := &net.TCPAddr{
		IP:   remoteIP,
		Port: rport,
	}
	tcpDialer := net.Dialer{
		LocalAddr: tcpLocalAddr,
		Timeout:   timeout,
	}

	// create udp local addr, remote addr also create its dialer
	udpLocalAddr := &net.UDPAddr{
		IP:   localIP,
		Port: lport,
	}
	udpRemoteAddr := &net.UDPAddr{
		IP:   remoteIP,
		Port: rport,
	}
	udpDialer := net.Dialer{
		LocalAddr: udpLocalAddr,
		Timeout:   timeout,
	}

	// tcpConn
	tcpConn, err := tcpDialer.Dial("tcp", tcpRemoteAddr.String())
	if err != nil {
		return nil, err
	}

	// udpConn
	udpConn, err := udpDialer.Dial("udp", udpRemoteAddr.String())
	if err != nil {
		return nil, err
	}

	return &Client{
		TCPConn: tcpConn,
		UDPConn: udpConn,
		TCPResp: make(chan uint32),
		TCPReq:  make(chan uint32),
		UDPResp: make(chan uint32),
		UDPReq:  make(chan uint32),
		ErrorCh: make(chan error),
	}, nil
}

func (c *Client) TCPResponse() {
	buf := make([]byte, 1024)
	n, err := c.TCPConn.Read(buf)
	if err != nil {
		c.ErrorCh <- errors.New("[TCP] error reading response")
	}
	resp := buf[:n]
	c.TCPResp <- binary.BigEndian.Uint32(resp)
}

func (c *Client) UDPResponse() {
	buf := make([]byte, 1024)
	n, err := c.UDPConn.Read(buf)
	if err != nil {
		c.ErrorCh <- errors.New("[UDP] error reading response")
	}
	resp := buf[:n]
	c.UDPResp <- binary.BigEndian.Uint32(resp)
}

func (c *Client) TCPRequest() {
	req := <-c.TCPReq
	data := c.intToByte(req)
	_, err := c.TCPConn.Write(data)
	if err != nil {
		c.ErrorCh <- errors.New("[TCP] error sending request")
	}
}

func (c *Client) UDPRequest() {
	req := <-c.UDPReq
	data := c.intToByte(req)

	_, err := c.UDPConn.Write(data)
	if err != nil {
		c.ErrorCh <- errors.New("[UDP] error sending request")
	}
}
