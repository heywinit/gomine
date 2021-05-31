package mcprot

import (
	"errors"

	"github.com/BRA1L0R/go-mcprot/packets"
	"github.com/BRA1L0R/go-mcprot/varint"
)

func (mc *McProt) ReceiveUncompressedPacket() (*packets.UncompressedPacket, error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.Connection)
	if err != nil {
		return nil, err
	}

	packetId, packetIdLen, err := varint.DecodeReaderVarInt(mc.Connection)
	if err != nil {
		return nil, err
	}

	packetContent := make([]byte, packetLength-packetIdLen)
	mc.Connection.Read(packetContent)

	packet := new(packets.UncompressedPacket)

	packet.Length = packetLength
	packet.PacketID = packetId
	packet.Data.Write(packetContent)

	return packet, nil
}

func (mc *McProt) ReceivePacket() (*packets.StandardPacket, error) {
	packetLength, _, err := varint.DecodeReaderVarInt(mc.Connection)
	if err != nil {
		return nil, err
	}

	dataLength, dLenLen, err := varint.DecodeReaderVarInt(mc.Connection)
	if err != nil {
		return nil, err
	}

	if dataLength != 0 {
		// drain the remaining packet
		drain := make([]byte, packetLength-dLenLen)
		mc.Connection.Read(drain)

		return nil, errors.New("compressed packet received, unable to process such packet at the moment")
	}

	packetId, pIdLen, err := varint.DecodeReaderVarInt(mc.Connection)
	if err != nil {
		return nil, err
	}

	newPacket := new(packets.StandardPacket)
	newPacket.Length = packetLength
	newPacket.DataLength = dataLength
	newPacket.PacketID = packetId

	remainingDataLen := packetLength - dLenLen - pIdLen

	data := make([]byte, remainingDataLen)
	mc.Connection.Read(data)

	newPacket.Data.Write(data)

	return newPacket, nil
}
