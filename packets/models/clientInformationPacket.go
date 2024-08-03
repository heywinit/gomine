package models

import "github.com/heywinit/gomine/packets"

type ClientInformationPacket struct {
	packets.MinecraftPacket
	Locale              string `mc:"string"`
	ViewDistance        byte   `mc:"inherit"`
	ChatMode            int32  `mc:"varint"` // 0 = enabled, 1 = commands only, 2 = hidden
	ChatColors          bool   `mc:"inherit"`
	DisplayedSkinParts  byte   `mc:"inherit"`
	MainHand            int32  `mc:"varint"` // 0 = left, 1 = right
	EnableTextFiltering bool   `mc:"inherit"`
	AllowServerListings bool   `mc:"inherit"`
}
