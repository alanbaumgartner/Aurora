package internal

import (
	"encoding/json"
	"net"
)

type Client struct {
	conn net.Conn
	enc  *json.Encoder
	dec  *json.Decoder
}

func (client *Client) GetConn() net.Conn {
	return client.conn
}

func (client *Client) GetEncoder() json.Encoder {
	return *client.enc
}

func (client *Client) GetDecoder() json.Decoder {
	return *client.dec
}

type List struct {
	Clients []Client
}

func NewList() List {
	self := List{}
	self.Clients = []Client{}
	return self
}

func (list *List) Get(index int) Client {
	if index < len(list.Clients) {
		return list.Clients[index]
	}
	return Client{}
}

func (list *List) Add(conn net.Conn) {
	for _, client := range list.Clients {
		if client.GetConn() == conn {
			return
		}
	}
	newConnData := Client{conn, json.NewEncoder(conn), json.NewDecoder(conn)}
	list.Clients = append(list.Clients, newConnData)
}

func (list *List) Remove(conn net.Conn) {
	for index, client := range list.Clients {
		if client.GetConn() == conn {
			list.Clients = append(list.Clients[:index], list.Clients[index+1:]...)
			return
		}
	}
}

func (list *List) All() []Client {
	return list.Clients
}

func (list *List) Clear() {
	list.Clients = []Client{}
}

func (list *List) isEmpty() bool {
	if len(list.Clients) == 0 {
		return true
	}
	return false
}
