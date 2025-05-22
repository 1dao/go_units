package share

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type XPTPack struct {
	ProtoType uint16 // 协议类型
	//Encrypted uint8  // 加密标识
	packSeq uint16 // 包序号
	// MessageID   uint8  // 消息ID
	// ChunkSeq    uint16 // 分片序号
	// TotalChunks uint16 // 总分片数
	RawData []byte // 原始数据
	//Decrypted []byte // 解密数据
}

var protoSpecialMap map[uint16]bool = make(map[uint16]bool)

func init() {
	protoSpecialMap[Protocol_Heartbeat] = true
	//protoSpecialMap[Protocol_Login] = true
	//protoSpecialMap[Protocol_Logout] = true
}

func (p *XPTPack) IsSpecialProtocol() bool {
	return protoSpecialMap[p.ProtoType]
}

func Unmarshal(data []byte) (*XPTPack, error) {
	p := &XPTPack{}
	fmt.Println("special protocol recv packet", data)

	var curser uint16 = 0
	var n uint16 = binary.BigEndian.Uint16(data[curser : curser+2])
	curser += 2
	if len(data) != int(n) {
		return nil, errors.New("invalid data")
	}

	p.ProtoType = binary.BigEndian.Uint16(data[curser : curser+2])
	curser += 2

	if protoSpecialMap[p.ProtoType] {
		return p, nil
	}

	p.packSeq = binary.BigEndian.Uint16(data[curser : curser+2])
	curser += 2

	var l uint16 = binary.BigEndian.Uint16(data[curser : curser+2])
	curser += 2

	p.RawData = data[curser : curser+l]

	return p, nil
}

func Marshal(pt uint16, seq uint16, data []byte) ([]byte, error) {
	curser := 0
	// pack len
	var packet []byte
	if protoSpecialMap[pt] {
		packet = make([]byte, 2+2)
		binary.BigEndian.PutUint16(packet[curser:curser+2], uint16(2+2))
		curser += 2

		// pack proto type
		binary.BigEndian.PutUint16(packet[curser:curser+2], pt)
		curser += 2

		fmt.Println("special protocol send packet", packet)
		return packet, nil
	} else {
		packet = make([]byte, 2+2+4+2+2+len(data))
	}
	binary.BigEndian.PutUint16(packet[curser:curser+2], uint16(2+2+2+2+len(data)))
	curser += 2

	// pack proto type
	binary.BigEndian.PutUint16(packet[curser:curser+2], pt)
	curser += 2

	//  pack seq
	binary.BigEndian.PutUint16(packet[curser:curser+2], seq) // TODO : seq
	curser += 2

	// pack data len
	binary.BigEndian.PutUint16(packet[curser:curser+2], uint16(len(data)))
	curser += 2

	// pack data
	copy(packet[curser:], data)
	curser += len(data)

	return packet[:curser], nil
}
