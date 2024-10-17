package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// based on this edit you server tcp req and udp req and then re write your lookfortcp, gettcp and handleudp,
// also you will need to change the way you handle sending of response.
// 
type TCPRequest struct {
	Addr net.TCPAddr
	Data uint32
}

type UDPRequest struct {
	Addr net.UDPAddr
	Data uint32
}

func (s *Server) intToByte(n uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		s.ErrorCh <- err
	}
	return buf.Bytes()
}

type Server struct {
	logger          *log.Logger
	TCPListener     net.Listener
	UDPConn         *net.UDPConn
	TCPReq, TCPResp chan uint32
	UDPReq, UDPResp chan uint32
	ErrorCh         chan error
}

func NewServer(laddr, raddr string, lport, rport int, timeout time.Duration) (*Server, error) {
	logger := log.New(os.Stdout, "[SERVER]", log.Default().Flags())
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
		logger:      logger,
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
	go s.HandleUDPConnection()
}

func (s *Server) LookForTCPConnection() {
	for {
		conn, err := s.TCPListener.Accept()
		if err != nil {
			s.ErrorCh <- fmt.Errorf("[TCP] error in establishing connection\n")
			continue
		}
		s.logger.Println("established connection with ", conn.RemoteAddr().String())
		s.logger.Println("server will starting handeling requests for ", conn.RemoteAddr().String())
		go s.getTCPRequest(conn)
	}
}

func (s *Server) HandleUDPConnection() {
	for {
		buf := make([]byte, 1024)
		n, raddr, err := s.UDPConn.ReadFromUDP(buf)
		if err != nil {
			s.ErrorCh <- fmt.Errorf("[UDP] error in reading from %s\n", raddr.String())
			continue
		}
		req := buf[:n]

	}
}

func (s *Server) getTCPRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		s.ErrorCh <- fmt.Errorf("error in reading message from %s\n", conn.RemoteAddr().String())
	}
	req := buf[:n]
	s.TCPReq <- binary.BigEndian.Uint32(req)
}

func (s *Server) getUDPRequest(req uint32) {

}

// can we generalise send response ??
func (s *Server) SendResponse() {
	select {
	case <-s.TCPReq:

	case <-s.UDPReq:
	default:
		s.ErrorCh <- fmt.Errorf("unkown protocol is sending a request\n")
	}
}
