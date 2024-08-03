package models

import "github.com/heywinit/gomine/packets"

type SetCompressionPacket struct {
	packets.MinecraftPacket

	Threshold int32 `mc:"varint"`
}
