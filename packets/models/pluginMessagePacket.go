package models

import "github.com/heywinit/gomine/packets"

type ServerboundPluginMessagePacket struct {
	packets.MinecraftPacket
	Channel string `mc:"string"` //identifier
	Data    []byte `mc:"bytearray" len:"32767"`
}

type ClientboundPluginMessagePacket struct {
	packets.MinecraftPacket
	Channel string `mc:"string"`                  //identifier
	Data    []byte `mc:"bytearray" len:"1048576"` //idk why the sizes are different but yes
}
