package models

import (
	"github.com/heywinit/gomine/packets"
)

type LoginStartPacket struct {
	packets.MinecraftPacket

	Name string `mc:"string"`
	UUID []byte `mc:"bytes" len:"16"`
}
