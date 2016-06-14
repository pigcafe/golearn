package main

import (
  "github.com/gorilla/websocket"
  "net/http"
  "log"
  "github.com/pigcafe/chat/trace"
)

type room struct {
  // "forward" has messages for forwarding to other clients.
  forward chan []byte
  // "join" has clients who try to join to this chat room(self).
  join chan *client
  // "leave" has clients who try to leave from this chat room(self).
  leave chan *client
  // "clients" keeps all clients who are existing in this chat room(self).
  clients map[*client]bool
  // "tracer" receives the operation logs in chat room.
  tracer trace.Tracer
}

func newRoom() *room {
  return &room {
    forward: make(chan []byte),
    join: make(chan *client),
    leave: make(chan *client),
    clients: make(map[*client]bool),
    tracer: trace.Off(),
  }
}

func (r *room) run()  {
  for {
    select {
    case client := <-r.join:
      // Join
      r.clients[client] = true
      r.tracer.Trace("New client comes in.")
    case client := <-r.leave:
      // Leave
      delete(r.clients, client)
      close(client.send)
      r.tracer.Trace("A client goes out.")
    case msg := <-r.forward:
      r.tracer.Trace("Recieved a message: ", string(msg))
      // Forward the message to all clients
      for client := range r.clients {
        select {
        case client.send <- msg:
          // Send message
          r.tracer.Trace(" -- sent to client.")
        default:
          // Fail to send meesage
          delete(r.clients, client)
          close(client.send)
          r.tracer.Trace(" -- Failed to send. clean up the client.")
        }
      }
    }
  }
}

const (
  socketBufferSize = 1024
  messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
                  ReadBufferSize: socketBufferSize,
                  WriteBufferSize: socketBufferSize,
                }

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
  socket, err := upgrader.Upgrade(w, req, nil)
  if err != nil {
    log.Fatal("ServeHTTP:", err)
    return
  }
  client := &client {
    socket: socket,
    send: make(chan []byte, messageBufferSize),
    room: r,
  }
  r.join <- client
  defer func() { r.leave <- client }()
  go client.write()
  client.read()
}
