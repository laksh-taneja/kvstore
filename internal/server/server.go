package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/laksh-taneja/kvstore/internal/cache"
)

type CommandType string
type ResponseType string

const (
	CmdGet CommandType = "GET"
	CmdSet CommandType = "SET"
)

const (
	RespNotFound   ResponseType = "ERR_NOT_FOUND\n"
	OK             ResponseType = "OK\n"
	RespBadCommand ResponseType = "ERR_BAD_SYNTAX\n"
)

type Server struct {
	cache *cache.Cache
	port  string
}

type Request struct {
	Op    CommandType
	Key   string
	Value any
}

// connection procedures
func parseRequest(r string) (*Request, error) {
	req := strings.Fields(r)
	if len(req) < 2 {
		return nil, fmt.Errorf("not enough args, need at least 2, got %v", len(req))
	}

	cmd := CommandType(req[0])

	switch cmd {
	case CmdGet:
		if len(req) != 2 {
			return nil, fmt.Errorf("Invalid arg count for GET, need 2 got %v", len(req))
		}
	case CmdSet:
		if len(req) != 3 {
			return nil, fmt.Errorf("Invalid arg count for SET, need 3 got %v", len(req))
		}
	default:
		return nil, fmt.Errorf("Invalid command SET | GET available %v", cmd)
	}

	reqObj := &Request{
		Op:  cmd,
		Key: req[1],
	}
	if cmd == CmdSet {
		reqObj.Value = req[2]
	}

	return reqObj, nil
}

func (s *Server) handleConnection(c net.Conn) {
	scanner := bufio.NewScanner(c)
	defer c.Close()

	for scanner.Scan() {
		r := scanner.Text()
		req, err := parseRequest(r)
		if err != nil {
			log.Printf("couldn't parse string: %v", err)
			c.Write([]byte(RespBadCommand))
			continue
		}
		if req.Op == "GET" {
			val, ok := s.cache.Access(req.Key)
			if !ok {
				c.Write([]byte(RespNotFound))
				continue
			}
			c.Write([]byte(OK))
			fmt.Fprintf(c, "%v\n", val)
			continue
		}
		if req.Op == "SET" {
			s.cache.Write(req.Key, req.Value)
			c.Write([]byte(OK))
			continue
		}
	}
	if scanner.Err() != nil {
		log.Printf("Can't scan request %v: %v", c.RemoteAddr(), scanner.Err())
		return
	}
}

// server api
func New(c *cache.Cache, port string) *Server {
	if port == "" {
		log.Printf("Invalid port, defaulting to :8080")
		port = ":8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	return &Server{
		cache: c,
		port:  port,
	}
}

func (s *Server) Run() error {

	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to handle connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}
