package util

import (
	"encoding/json"
	"net"
)

type Client struct {
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
}

func (client *Client) GetConn() net.Conn {
	return client.connection
}

func (client *Client) SetConn(newConn net.Conn) {
	client.connection = newConn
}

func (client *Client) GetEncoder() json.Encoder {
	return *client.encoder
}

func (client *Client) SetEncoder(newEncoder *json.Encoder) {
	client.encoder = newEncoder
}

func (client *Client) GetDecoder() json.Decoder {
	return *client.decoder
}

func (client *Client) SetDecoder(newDecoder *json.Decoder) {
	client.decoder = newDecoder
}
