package a2s

import (
	"errors"

	"github.com/handsomematt/go-valve/util"
)

var errBadPacketHeader = errors.New("bad packet header")
var errGoldSrc = errors.New("goldsrc server not supported")

const theShipAppID = 2400

// Info about a Source server
type Info struct {
	ProtocolVersion byte   // protocol version used by the server
	ServerName      string // name of the server
	Map             string // map the server has currently loaded
	Folder          string // name of the folder containing the game files
	Game            string // full name of the game
	AppID           uint16 // steam application id of game
	Players         byte   // number of players on the server
	MaxPlayers      byte   // max number of players the server reports it can hold
	Bots            byte   // number of bots on the server
	ServerType      byte   // type of server: d(edicated) l(isten) p(roxy)
	Platform        byte   // operating system of the server l(inux) w(indows) m(ac)/o(sx)
	Passworded      bool   // whether the server requires a password
	VAC             bool   // whether the server uses VAC
	Version         string // version of the game installed on the server

	// only if the server is running The Ship
	TheShip struct {
		Mode      uint8
		Witnesses uint8
		Duration  uint8
	}

	// only if the server is running Source TV
	SourceTV struct {
		Port uint16
		Name string
	}

	// Extra Data Fields
	Port     uint16 // the server's game port number
	SteamID  uint64 // server's steamid
	Keywords string // tags that describe the game according to the server
	GameID   uint64 // the server's 64 bit gameid
}

// QueryInfo ...
func (querier *Querier) QueryInfo() (*Info, error) {
	message := []byte("\xFF\xFF\xFF\xFFTSource Engine Query\x00")

	_, err := querier.udpConnection.Write(message)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1400)
	n, _, err := querier.udpConnection.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	reader := util.NewBinaryReader(buffer[:n])
	reader.Read(4) // A2SHeader

	header := reader.ReadUInt8()
	if header == 'm' {
		return nil, errGoldSrc
	}
	if header != 'I' {
		return nil, errBadPacketHeader
	}

	var info Info
	// buffer[:4] (A2S Header (assume 0xFFFFFFF))
	info.ProtocolVersion = reader.ReadUInt8()
	info.ServerName = reader.ReadCString()
	info.Map = reader.ReadCString()
	info.Folder = reader.ReadCString()
	info.Game = reader.ReadCString()
	info.AppID = reader.ReadUInt16()
	info.Players = reader.ReadUInt8()
	info.MaxPlayers = reader.ReadUInt8()
	info.Bots = reader.ReadUInt8()
	info.ServerType = reader.ReadUInt8()
	info.Platform = reader.ReadUInt8()
	info.Passworded = reader.ReadBool()
	info.VAC = reader.ReadBool()

	// The Ship
	if info.AppID == 2400 {
		info.TheShip.Mode = reader.ReadUInt8()
		info.TheShip.Witnesses = reader.ReadUInt8()
		info.TheShip.Duration = reader.ReadUInt8()
	}

	info.Version = reader.ReadCString()

	if !reader.More() {
		return &info, nil
	}

	edf := reader.ReadUInt8()

	if edf&0x80 != 0 {
		info.Port = reader.ReadUInt16()
	}

	if edf&0x10 != 0 {
		info.SteamID = reader.ReadUInt64()
	}

	if edf&0x40 != 0 {
		info.SourceTV.Port = reader.ReadUInt16()
		info.SourceTV.Name = reader.ReadCString()
	}

	if edf&0x20 != 0 {
		info.Keywords = reader.ReadCString()
	}

	if edf&0x01 != 0 {
		info.GameID = reader.ReadUInt64()
	}

	return &info, nil
}
