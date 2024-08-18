package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}

	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("Nice! I will call you %s", c.name))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.name))
	c.msg(fmt.Sprintf("Welcome to the room %s!", r.name))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("Available rooms are: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}

	c.room.broadcast(c, c.name+": "+strings.Join(args[1:], " "))
}

func (s *server) quit(c *client) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("Sad to see you go :(  Bye...")

	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}
