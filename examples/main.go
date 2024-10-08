package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/heywinit/gomine"
	"github.com/heywinit/gomine/packets"
	"github.com/heywinit/gomine/packets/models"
	"runtime"
)

type ClientBoundChatMessage struct {
	packets.MinecraftPacket

	JsonData string     `mc:"string"`
	Position byte       `mc:"inherit"`
	Sender   complex128 `mc:"inherit"`
}

type ServerBoundChatMessage struct {
	packets.MinecraftPacket

	Message string `mc:"string"`
}

type ChatMessage struct {
	Translate string `json:"translate"`

	With []struct {
		Text string `json:"text"`
	}
}

var (
	host     = flag.String("host", "127.0.0.1", "Server host")
	port     = flag.Uint("port", 25565, "Server port")
	username = flag.String("username", "Notch", "In-game username")
)

func main() {
	flag.Parse()

	client := gomine.Client{}
	err := client.Initialize(*host, uint16(*port), 766, *username)
	if err != nil {
		panic(err)
	}

	botInfoMessage := new(ServerBoundChatMessage)
	botInfoMessage.PacketID = 0x03
	botInfoMessage.Message = fmt.Sprintf("I'm running on %v, %v", runtime.GOOS, runtime.GOARCH)

	err = client.WritePacket(botInfoMessage)
	if err != nil {
		panic(err)
	}

	for {
		packet, err := client.ReceivePacket()
		if err != nil {
			panic(err)
		}

		switch packet.PacketID {
		case 0x0E:
			// receivedChatMessage := ClientBoundChatMessage{MinecraftPacket: packet}
			receivedChatMessage := new(ClientBoundChatMessage)
			err := packet.DeserializeData(receivedChatMessage)
			if err != nil {
				panic(err)
			}

			chatMessage := new(ChatMessage)
			err = json.Unmarshal([]byte(receivedChatMessage.JsonData), chatMessage)
			if err != nil {
				panic(err)
			}

			if chatMessage.Translate == "chat.type.text" {
				user := chatMessage.With[0].Text
				playerText := chatMessage.With[1].Text

				if user == *username {
					continue
				}

				fmt.Printf("<%s> %s\n", user, playerText)

				chatMessage := new(ServerBoundChatMessage)
				chatMessage.PacketID = 0x03
				chatMessage.Message = playerText

				err := client.WritePacket(chatMessage)
				if err != nil {
					panic(err)
				}
			}
		case 0x1F:
			// receivedKeepalive := models.KeepAlivePacket{MinecraftPacket: packet}
			receivedKeepalive := new(models.KeepAlivePacket)
			err := packet.DeserializeData(receivedKeepalive)
			if err != nil {
				panic(err)
			}

			serverBoundKeepalive := new(models.KeepAlivePacket)
			serverBoundKeepalive.PacketID = 0x10
			serverBoundKeepalive.KeepAliveID = receivedKeepalive.KeepAliveID

			err = client.WritePacket(serverBoundKeepalive)
			if err != nil {
				panic(err)
			}
		}
	}
}
