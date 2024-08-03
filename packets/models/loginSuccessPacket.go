package models

import "github.com/heywinit/gomine/packets"

type LoginSuccessPacket struct {
	packets.MinecraftPacket

	UUID               []byte `mc:"bytes" len:"16"`
	Username           string `mc:"string"`
	NumberOfProperties int    `mc:"varint"`
	Properties         []struct {
		Name      string `mc:"string"`
		Value     string `mc:"string"`
		IsSigned  bool   `mc:"bool"`
		Signature string `mc:"string"`
	}
}
