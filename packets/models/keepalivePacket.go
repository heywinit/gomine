package models

import "github.com/heywinit/gomine/packets"

type KeepAlivePacket struct {
	packets.MinecraftPacket

	KeepAliveID int64 `mc:"inherit"`
}
