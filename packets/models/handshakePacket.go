package models

import "github.com/heywinit/gomine/packets"

type HandshakePacket struct {
	packets.MinecraftPacket

	ProtocolVersion int32  `mc:"varint"`
	ServerAddress   string `mc:"string"`
	ServerPort      uint16 `mc:"inherit"`
	NextState       int32  `mc:"varint"`
}
