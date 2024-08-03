package gomine

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/heywinit/gomine/packets"
	"github.com/heywinit/gomine/packets/models"
	"net"
)

// LoginDisconnectError represents the disconnection of the client during the login state,
// during the initialization process
type LoginDisconnectError struct {
	// Reason is the json encoded Reason for which the client has been
	// disconnected from the server
	Reason string
}

func (disconnectError *LoginDisconnectError) Error() string {
	return fmt.Sprintf("gomine: server disconnected the client during initialization with message: %s", disconnectError.Reason)
}

// Connect initializes the connection to the server.
//
// host must have the format of "host:port" as a port has to be specified in order
// to open a connection. 25565 is not taken for granted
func (client *Client) Connect(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	client.connection = conn.(*net.TCPConn)

	return nil
}

// Initialize initializes the connection to the server by sending
// the handshake packet and the login packet
//
// host is the server fqdn or ip address of the server, port is the uint16 port where the server is listening on
// username is the in-game username the client will send to the server during handshaking. Might differ from the actual
// in-game username as the server sends a confirmation of it after the login state.
func (client *Client) Initialize(host string, port uint16, protocolVersion int32, username string) error {
	if err := client.Connect(fmt.Sprintf("%s:%v", host, port)); err != nil {
		return err
	}

	// Create handshake packet with the latest protocol version
	// and state information
	// NOTE: server address and server port are not used by the
	// notchian server, but it's best practice to fill them
	hp := models.HandshakePacket{
		MinecraftPacket: packets.MinecraftPacket{PacketID: 0x00},

		ProtocolVersion: protocolVersion,
		ServerAddress:   host,
		ServerPort:      port,
		NextState:       2,
	}

	err := client.WritePacket(&hp)
	if err != nil {
		return err
	}

	uuid := "069a79f444e94726a5befca90e38aaf5"
	//converting this hex string to a 128 bit integer
	// Remove "0x" prefix if present
	if len(uuid) >= 2 && uuid[:2] == "0x" {
		uuid = uuid[2:]
	}

	// Decode the hex string to bytes
	bytes, err := hex.DecodeString(uuid)
	if err != nil {
		fmt.Printf("invalid hexadecimal string: %v", err)
	}
	// Split the bytes into the most significant and least significant 64 bits
	msb := binary.BigEndian.Uint64(bytes[0:8])
	lsb := binary.BigEndian.Uint64(bytes[8:16])

	// Create a byte slice to hold the encoded UUID
	encodedUUID := make([]byte, 16)

	// Pack the most significant and least significant bits as unsigned 64-bit integers
	binary.BigEndian.PutUint64(encodedUUID[0:8], msb)
	binary.BigEndian.PutUint64(encodedUUID[8:16], lsb)

	loginPacket := models.LoginStartPacket{
		MinecraftPacket: packets.MinecraftPacket{PacketID: 0x00},

		Name: username,
		UUID: encodedUUID,
	}

	err = client.WritePacket(&loginPacket)
	if err != nil {
		return err
	}

	for {
		p, err := client.ReceivePacket()
		if err != nil {
			return err
		}

		switch p.PacketID {
		case 0x00: // disconnected
			disconnectPacket := new(models.DisconnectPacket)
			if err := p.DeserializeData(disconnectPacket); err != nil {
				return err
			}

			return &LoginDisconnectError{Reason: disconnectPacket.Reason}
		case 0x01: // client bound plugin message, we ignore it
			clientInfoPacket := models.ClientInformationPacket{
				MinecraftPacket:     packets.MinecraftPacket{PacketID: 0x00},
				Locale:              "en_US",
				ViewDistance:        10,
				ChatMode:            0,
				ChatColors:          false, //todo: add a way to enable chat colors
				DisplayedSkinParts:  0x7F,
				MainHand:            1, //right hand
				EnableTextFiltering: false,
				AllowServerListings: true,
			}

			err := client.WritePacket(&clientInfoPacket)

			return err
		case 0x03: // set compression
			setCompression := new(models.SetCompressionPacket)
			err := p.DeserializeData(setCompression)
			if err != nil {
				return err
			}

			if setCompression.Threshold < 0 {
				return errors.New("gomine: server sent a set compression packet with a negative threshold")
			}

			client.CompressionThreshold = setCompression.Threshold
		case 0x02: // login success
			loginSuccess := new(models.LoginSuccessPacket)
			err := p.DeserializeData(loginSuccess)

			if err != nil {
				return err
			}

			loginAckPacket := models.LoginAcknowledgedPacket{
				MinecraftPacket: packets.MinecraftPacket{PacketID: 0x03},
			}

			err = client.WritePacket(&loginAckPacket)
			if err != nil {
				return err
			}
		}
	}
}
