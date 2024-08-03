package gomine

import (
	"io"

	"github.com/heywinit/gomine/packets"
)

// SerializablePacket defines the standard methods that a struct should have
// in order to be serializable by the library
//
// You can actually create your own methods as long as they respect this standard
type SerializablePacket interface {
	// SerializeData takes an interface pointer as input and serializes all the fields in the
	// data buffer. It can and will return an error in case of invalid data
	SerializeData(inter interface{}) error

	SerializeCompressed(writer io.Writer, compressionThreshold int) error
	SerializeUncompressed(writer io.Writer) error
}

// WritePacket calls SerializeData and then calls WriteRawPacket
func (client *Client) WritePacket(packet SerializablePacket) error {
	client.writeMu.Lock()
	defer client.writeMu.Unlock()

	if err := packet.SerializeData(packet); err != nil {
		return err
	}

	if client.IsCompressionEnabled() {
		return packet.SerializeCompressed(client.connection, int(client.CompressionThreshold))
	} else {
		return packet.SerializeUncompressed(client.connection)
	}
}

// WriteRawPacket takes a rawpacket as input and serializes it in the connection
func (client *Client) WriteRawPacket(rawPacket *packets.MinecraftRawPacket) error {
	client.writeMu.Lock()
	defer client.writeMu.Unlock()

	if client.IsCompressionEnabled() {
		return rawPacket.WriteCompressed(client.connection)
	} else {
		return rawPacket.WriteUncompressed(client.connection)
	}
}

// ReceiveRawPacket reads a raw packet from the connection but doesn't deserialize
// neither uncompress it
func (client *Client) ReceiveRawPacket() (*packets.MinecraftRawPacket, error) {
	client.readMu.Lock()
	defer client.readMu.Unlock()

	if client.IsCompressionEnabled() {
		return packets.FromCompressedReader(client.connection)
	} else {
		return packets.FromUncompressedReader(client.connection)
	}
}

// ReceivePacket receives and deserializes a packet from the connection, uncompressing it
// if necessary
func (client *Client) ReceivePacket() (*packets.MinecraftPacket, error) {
	rawPacket, err := client.ReceiveRawPacket()
	if err != nil {
		return nil, err
	}

	return packets.FromRawPacket(rawPacket)
}
