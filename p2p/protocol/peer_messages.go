package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"log"
)

// PeerMessageType -
type PeerMessageType = uint16

// messgae tags
const (
	UnknownTag          PeerMessageType = 0x0
	AdvertiseTag                        = 0x03
	GetCurrentBranchTag                 = 0x10
	CurrentHeadTag                      = 0x14
)

// sizes
const (
	peerMessageLenSize       = 4
	peerMessageTypeSize      = 2
	peerMessageStringLenSize = 4
	chainIDLen               = 4
	blockHeaderSize          = 4
	levelSize                = 4
	byteSize                 = 1
	predecessorSize          = 32
	timestampSize            = 8
	operationHashSize        = 32
	fitnessBlockSize         = 4
	contextSize              = 32
)

// BootstrapMsg -
type BootstrapMsg struct {
}

// AdvertiseMsg -
type AdvertiseMsg struct {
	Addresses []string
}

// GetCurrentBranchMsg -
type GetCurrentBranchMsg struct {
	Branch string
}

// GetCurrentHeadMsg -
type GetCurrentHeadMsg struct {
	ChainID string
}

type byteSlice = []byte

// BlockHeader -
type BlockHeader struct {
	Level          uint32
	Proto          byte
	Predecessor    []byte
	Timestamp      int64
	ValidationPass byte
	OperationHash  []byte
	Fitness        []byteSlice
	Context        []byte
	ProtocolData   []byte
}

// CurrentHeadMsg -
type CurrentHeadMsg struct {
	ChainID            string
	CurrentBlockHeader BlockHeader
}

func (BootstrapMsg) toBytes() []byte {
	// "000000020002"
	data := []byte{0, 0, 0, 2, 0, 2}
	return data
}

func (msg CurrentHeadMsg) toBytes() []byte {
	bytes := make([]byte, peerMessageLenSize+peerMessageTypeSize+chainIDLen+blockHeaderSize+levelSize+byteSize+
		predecessorSize+timestampSize+byteSize+operationHashSize+fitnessBlockSize+contextSize)

	offset := 0
	binary.BigEndian.PutUint32(bytes[offset:offset+peerMessageLenSize], 124)
	offset += peerMessageLenSize

	binary.BigEndian.PutUint16(bytes[offset:offset+peerMessageTypeSize], 17)
	offset += peerMessageTypeSize

	chainID, err := hex.DecodeString(msg.ChainID)
	if err != nil {
		log.Println("serialize error: ", err)
		return bytes
	}
	copy(bytes[offset:offset+chainIDLen], chainID)
	offset += chainIDLen

	binary.BigEndian.PutUint32(bytes[offset:offset+blockHeaderSize], 114)
	offset += blockHeaderSize

	return bytes
}

func (branch GetCurrentHeadMsg) toBytes() []byte {
	b, _ := hex.DecodeString(branch.ChainID)
	size := len(b)
	bytes := make([]byte, peerMessageLenSize+peerMessageTypeSize+size)

	index := 0
	binary.BigEndian.PutUint32(bytes[:peerMessageLenSize], uint32(size+2))
	index += peerMessageLenSize

	// тэг для запроса currentHead
	copy(bytes[index:index+peerMessageTypeSize], []byte{0, 19})
	index += peerMessageTypeSize

	copy(bytes[index:index+size], b)
	return bytes
}

func parseString(data []byte) (str string, offset uint32) {
	offset = 0
	length := binary.BigEndian.Uint32(data[offset : offset+peerMessageStringLenSize])
	offset += peerMessageStringLenSize

	str = string(data[offset : offset+length])
	offset += length
	return
}

func newAdvertiseMsg(data []byte, messageLen uint32) (message AdvertiseMsg) {
	var offset uint32 = 0
	for offset+peerMessageLenSize+peerMessageTypeSize < messageLen {
		str, off := parseString(data[offset:])
		message.Addresses = append(message.Addresses, str)
		offset += off
	}
	return
}

func newGetCurrentBranchMsg(data []byte) (message GetCurrentBranchMsg) {
	message.Branch = hex.EncodeToString(data)
	return
}

// NewCurrentHeadMsg -
func NewCurrentHeadMsg(data []byte) (message CurrentHeadMsg) {
	offset := 0

	message.ChainID = hex.EncodeToString(data[offset : offset+chainIDLen])
	offset += chainIDLen

	sizeBlockHeader := binary.BigEndian.Uint32(data[offset : offset+blockHeaderSize])
	offset += blockHeaderSize

	message.CurrentBlockHeader.Level = binary.BigEndian.Uint32(data[offset : offset+levelSize])
	offset += levelSize

	message.CurrentBlockHeader.Proto = data[offset]
	offset += byteSize

	message.CurrentBlockHeader.Predecessor = append(message.CurrentBlockHeader.Predecessor, data[offset:offset+predecessorSize]...)
	offset += predecessorSize

	message.CurrentBlockHeader.Timestamp = int64(binary.BigEndian.Uint64(data[offset : offset+timestampSize]))
	offset += timestampSize

	message.CurrentBlockHeader.ValidationPass = data[offset]
	offset += byteSize

	message.CurrentBlockHeader.OperationHash = append(message.CurrentBlockHeader.OperationHash, data[offset:offset+operationHashSize]...)
	offset += operationHashSize

	sizeFitnessBlock := binary.BigEndian.Uint32(data[offset : offset+fitnessBlockSize])
	offset += fitnessBlockSize

	if sizeFitnessBlock != 0 {
		for sizeFitnessBlock != 0 {
			sizeFitness := binary.BigEndian.Uint32(data[offset : offset+4])
			offset += 4
			sizeFitnessBlock -= 4
			fitness := make([]byte, sizeFitness)
			fitness = data[offset : offset+int(sizeFitness)]
			offset += int(sizeFitness)
			sizeFitnessBlock -= sizeFitness
			message.CurrentBlockHeader.Fitness = append(message.CurrentBlockHeader.Fitness, fitness)
		}
	}

	message.CurrentBlockHeader.Context = append(message.CurrentBlockHeader.Context, data[offset:offset+contextSize]...)
	offset += contextSize

	message.CurrentBlockHeader.ProtocolData = append(message.CurrentBlockHeader.ProtocolData, data[offset:sizeBlockHeader+8]...)

	return
}

func parseMessage(data []byte) (obj interface{}, messageType PeerMessageType) {
	obj = nil
	messageType = UnknownTag

	if len(data) < peerMessageLenSize {
		return
	}

	var offset uint32 = 0
	length := binary.BigEndian.Uint32(data[offset : offset+peerMessageLenSize])
	offset += peerMessageLenSize

	if length+peerMessageLenSize != uint32(len(data)) {
		return
	}

	messageType = binary.BigEndian.Uint16(data[offset : offset+peerMessageTypeSize])
	offset += peerMessageTypeSize

	switch messageType {
	case GetCurrentBranchTag:
		obj = newGetCurrentBranchMsg(data[offset : offset+length-peerMessageTypeSize])

	case AdvertiseTag:
		obj = newAdvertiseMsg(data[offset:offset+length-peerMessageTypeSize], length)

	case CurrentHeadTag:
		obj = NewCurrentHeadMsg(data[offset : offset+length-peerMessageTypeSize])
	}

	return
}
