package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gtank/isaac"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
	"rsps/entity"
	"rsps/net/packet/outgoing/login"
	"rsps/repository"
	"strings"
)

type LoginHandler struct{}

func (l *LoginHandler) HandlePacket(c *TCPClient) {
	if c.loginState == HandshakeStage {
		var packet LoginZero
		err := binary.Read(c.connection, binary.BigEndian, &packet)
		if err != nil {
			log.Println("login stage 1 error: " + err.Error())
			stopLogin(c, login.LoginServerRejected)
			return
		}

		c.Enqueue(&login.LoginHandshakeResponse{
			LoginStatus:      login.MayProceed,
			ServerSessionKey: 12345678,
		})
		c.Enqueue(&flush{})
		c.loginState = LoginStage
	} else if c.loginState == LoginStage {
		var packet LoginPacket
		err := binary.Read(c.connection, binary.BigEndian, &packet)
		if err != nil {
			log.Printf("login stage 1 error: " + err.Error())
			stopLogin(c, login.LoginServerRejected)
			return
		}

		e, _ := new(big.Int).SetString("33280025241734061313051117678670856264399753710527826596057587687835856000539511539311834363046145710983857746766009612538140077973762171163294453513440619295457626227183742315140865830778841533445402605660729039310637444146319289077374748018792349647460850308384280105990607337322160553135806205784213241305", 10)
		m, _ := new(big.Int).SetString("91553247461173033466542043374346300088148707506479543786501537350363031301992107112953015516557748875487935404852620239974482067336878286174236183516364787082711186740254168914127361643305190640280157664988536979163450791820893999053469529344247707567448479470137716627440246788713008490213212272520901741443", 10)
		encrypted := big.NewInt(0)
		encrypted.SetBytes(packet.EncryptionBytes[:])
		var rs big.Int
		rs = *rs.Exp(encrypted, e, m)
		rsaBuffer := bytes.NewBuffer(rs.Bytes())

		var rsaPacket RsaPacket
		err = binary.Read(rsaBuffer, binary.BigEndian, &rsaPacket)
		if err != nil {
			fmt.Println(err)
			stopLogin(c, login.BadSessionId)
			return
		}
		name, err := rsaBuffer.ReadBytes(10)
		if err != nil {
			fmt.Println(err)
			stopLogin(c, login.BadSessionId)
			return
		}
		//log.Printf("%s", string(name))
		pass, err := rsaBuffer.ReadBytes(10) // password
		if err != nil {
			fmt.Println(err)
			stopLogin(c, login.BadSessionId)
			return
		}
		//log.Printf("%s", string(pass))

		playerName := strings.TrimSpace(string(name))
		// load player
		player := entity.NewPlayer()
		playerPosition, hashedPassword, err := c.PlayerRepository.Load(playerName)
		if err == repository.PlayerNotFoundError {
			hashedPassword, err = c.PlayerRepository.Create(playerName, pass, player.Position)
			if err != nil {
				stopLogin(c, login.LoginServerRejected)
			}
		} else if err != nil {
			log.Printf("load player err: %s", err.Error())
			stopLogin(c, login.InvalidCredentials)
			return
		} else {
			player.Position = playerPosition
		}

		// check password
		err = bcrypt.CompareHashAndPassword(hashedPassword, pass)
		if err != nil {
			stopLogin(c, login.InvalidCredentials)
			return
		}

		c.Player = player
		err = c.Player.LoadPlayer(playerName)
		if err != nil {
			log.Printf(err.Error())
			stopLogin(c, login.LoginServerRejected)
			return
		}
		c.World.AddPlayer(player)

		inC := isaac.ISAAC{}

		sessionKey := make([]uint32, 4)
		sessionKey[0] = uint32(rsaPacket.ClientSessionKey >> 32)
		sessionKey[1] = uint32(rsaPacket.ClientSessionKey)
		sessionKey[2] = uint32(rsaPacket.ServerSessionKey >> 32)
		sessionKey[3] = uint32(rsaPacket.ServerSessionKey)
		inC.Generate(sessionKey)
		c.Decryptor = &inC

		for i := 0; i < 4; i++ {
			sessionKey[i] += 50
		}
		outC := isaac.ISAAC{}
		outC.Generate(sessionKey)
		c.Encryptor = &outC

		c.Enqueue(&login.LoginResponse{
			ReturnCode:   login.LoginSuccess,
			PlayerRights: 3,
			Flagged:      0,
		})
		c.Enqueue(&flush{})

		c.loginState = Initialize
	}
}

func stopLogin(c *TCPClient, reason byte) {
	c.loginState = Disconnected
	c.Enqueue(&login.LoginResponse{
		ReturnCode:   reason,
		PlayerRights: 0,
		Flagged:      0,
	})
	c.Enqueue(&flush{})
}

type LoginPacket struct {
	Type              byte
	PacketSize        byte
	Magic             byte
	Version           int16
	LowMem            byte
	Unknown           [9]int32
	EncryptPacketSize uint8
	EncryptionBytes   [128]byte
}

type RsaPacket struct {
	Id               byte
	ClientSessionKey int64
	ServerSessionKey int64
	Uid              int32
}

type LoginZero struct {
	Protocol byte
	NameHash byte
}
