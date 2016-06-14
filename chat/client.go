package main

import (
  "github.com/gorilla/websocket"
)

// "client" represents a person who are chatting.
type client struct {
  // "socket" means the websocket for myself.
  socket *websocket.Conn
  // "send" can get messages.
  send chan []byte
  // "room" means the room where this client joins at.
  room *room
}

func (c *client) read() {
  for {
    if _, msg, err := c.socket.ReadMessage(); err == nil {
      c.room.forward <- msg
    }else{
      break
    }
  }
  c.socket.Close()
}

func (c *client) write() {
  for msg := range c.send {
    if err := c.socket.WriteMessage(websocket.TextMessage, msg);
          err != nil {
        break
      }
  }
  c.socket.Close()
}
