package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

func (s *Server) intToByte(n uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		s.ErrorCh <- err
	}
	return buf.Bytes()
}

type Server struct {
	TCPListener     net.Listener
	UDPConn         *net.UDPConn
	TCPReq, TCPResp chan uint32
	UDPReq, UDPResp chan uint32
	ErrorCh         chan error
}

func NewServer(laddr, raddr string, lport, rport int, timeout time.Duration) (*Server, error) {
	localIP := net.ParseIP(laddr)

	tcpLocalAddr := &net.TCPAddr{
		IP:   localIP,
		Port: lport,
	}
	udpLocalAddr := &net.UDPAddr{
		IP:   localIP,
		Port: lport,
	}

	tcpListener, err := net.ListenTCP("tcp", tcpLocalAddr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", udpLocalAddr)
	if err != nil {
		return nil, err
	}
	return &Server{
		TCPListener: tcpListener,
		UDPConn:     udpConn,
		TCPReq:      make(chan uint32),
		TCPResp:     make(chan uint32),
		UDPReq:      make(chan uint32),
		UDPResp:     make(chan uint32),
	}, nil
}

func (s *Server) Start() {
	go s.LookForTCPConnection()
	go s.LookForUDPConnection()
}

func (s *Server) LookForTCPConnection() {
	
}

func (s *Server) LookForUDPConnection() {

}

func (s *Server) getTCPRequest(conn net.Conn) {

}

func (s *Server) getUDPRequest() {

}

// can we generalise send response ??
func (s *Server) SendResponse() {
	select {
	case <-s.TCPReq:

	case <-s.UDPReq:
	default:
		s.ErrorCh <- errors.New("unkown protocol is sending a request")
	}
}
