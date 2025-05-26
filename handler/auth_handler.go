package handler

import (
	"bytes"
	"encoding/binary"
	"fleare-cli/comm"
	"fmt"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

type Connection struct {
	// Define fields for connection details
	Conn     net.Conn
	ClientID string
}

func NewConnection(conn net.Conn, clientID string) *Connection {
	return &Connection{
		Conn:     conn,
		ClientID: clientID,
	}
}

func init() {
	// Initialize any necessary resources or configurations here
}

func ConnectWithPassword(host string, port int, username string, password string) (*Connection, error) {
	// Input validation
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port number: %d", port)
	}
	// if username == "" {
	// 	return nil, fmt.Errorf("username cannot be empty")
	// }
	// Password can be empty in some cases, so we don't validate it

	// Connect to the server
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", addr, err)
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	connection := NewConnection(conn, "")

	// Create authentication command
	cmd := &comm.Command{
		Command: "auth",
		Args:    []string{username, password},
	}

	// Send authentication command
	if err := connection.Write(cmd); err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to write authentication command: %w", err)
	}

	// Read server response
	var resp comm.Response
	if err = connection.Read(&resp); err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to read authentication response: %w", err)
	}

	// Check authentication result
	if resp.Status != "Ok" {
		connection.Close()
		if len(resp.Result) > 0 {
			return nil, fmt.Errorf("%s", string(resp.Result))
		}
		return nil, fmt.Errorf("authentication failed with unknown error")
	}
	connection.ClientID = resp.ClientId

	fmt.Printf("%s %s\n", resp.Status, string(resp.Result))
	return connection, nil
}

func (c *Connection) Read(msg proto.Message) error {

	var length uint32

	// Read the length prefix
	if err := binary.Read(c.Conn, binary.BigEndian, &length); err != nil {
		return err
	}

	// Read the message data
	data := make([]byte, length)
	_, err := c.Conn.Read(data)
	if err != nil {
		return err
	}

	// Unmarshal the data
	return proto.Unmarshal(data, msg)
}

func (c *Connection) Write(msg proto.Message) error {
	// Encode protobuf
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Create a buffer
	var buf bytes.Buffer

	// Write the length prefix (4 bytes, big endian)
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}

	// Write the protobuf data
	if _, err := buf.Write(data); err != nil {
		return err
	}

	// Send everything
	_, err = c.Conn.Write(buf.Bytes())
	return err
}

func (c *Connection) Close() error {
	return c.Conn.Close()
}
