package models

import "github.com/heywinit/gomine/packets"

type DisconnectPacket struct {
	packets.MinecraftPacket

	Reason string `mc:"string"`
}
