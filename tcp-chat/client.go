package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	MaxMessageLength     = 200 // Example length
	MaxMessagesPerMinute = 10
	MessageInterval      = time.Minute
)

type client struct {
	conn            net.Conn
	name            string
	room            *room
	commands        chan<- command
	lastMessageTime time.Time
	messageCount    int
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		if err := c.validateMessage(msg); err != nil {
			c.handleInvalidMessage(err)
			continue
		}

		if !c.canSendMessage() {
			c.err(fmt.Errorf("rate limit exceeded, please wait before sending another message"))
			continue
		}

		if containsForbiddenWords(msg) {
			c.err(fmt.Errorf("message contains inappropriate content"))
			continue
		}

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/name":
			c.commands <- command{
				id:     CMD_NAME,
				client: c,
				args:   args,
			}

		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}

		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}

		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}

		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}

		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) validateMessage(msg string) error {
	if len(msg) > MaxMessageLength {
		return fmt.Errorf("message exceeds maximum length of %d characters", MaxMessageLength)
	}
	return nil
}

func (c *client) canSendMessage() bool {
	if time.Since(c.lastMessageTime) > MessageInterval {
		c.messageCount = 0
		c.lastMessageTime = time.Now()
	}
	if c.messageCount >= MaxMessagesPerMinute {
		return false
	}
	c.messageCount++
	return true
}

func (c *client) handleInvalidMessage(err error) {
	c.err(fmt.Errorf("failed to send message: %v", err))
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func containsForbiddenWords(msg string) bool {
	var forbiddenWords = []string{"badword1", "badword2"}
	for _, word := range forbiddenWords {
		if strings.Contains(msg, word) {
			return true
		}
	}
	return false
}
