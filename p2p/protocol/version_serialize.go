package protocol

import (
	"encoding/binary"
)

const (
	nameLenSize     = 4
	partVersionSize = 2
)

func (version Version) toBytes() (bytes []byte, versionSize int) {
	nameLen := len(version.Name)
	versionSize = nameLenSize + nameLen + 2 + 2
	bytes = make([]byte, versionSize)

	index := 0
	binary.BigEndian.PutUint32(bytes[index:index+nameLenSize], uint32(nameLen))
	index += nameLenSize

	copy(bytes[index:index+len(version.Name)], version.Name)
	index += len(version.Name)

	binary.BigEndian.PutUint16(bytes[index:index+partVersionSize], version.Major)
	index += partVersionSize

	binary.BigEndian.PutUint16(bytes[index:index+partVersionSize], version.Minor)
	index += partVersionSize

	return
}

func versionsToBytes(versions []Version) (allBytes []byte) {
	fullSize := 0
	for _, version := range versions {
		bytes, versionSize := version.toBytes()
		allBytes = append(allBytes, bytes...)
		fullSize += versionSize
	}

	return allBytes
}

func newVersion(bytes []byte) (version Version, offset uint32) {

	size := binary.BigEndian.Uint32(bytes[offset : offset+nameLenSize])
	if size == 0 {
		return Version{}, 0
	}

	offset += nameLenSize

	version.Name = string(bytes[offset : offset+size])
	offset += size

	version.Major = binary.BigEndian.Uint16(bytes[offset : offset+partVersionSize])
	offset += partVersionSize

	version.Minor = binary.BigEndian.Uint16(bytes[offset : offset+partVersionSize])
	offset += partVersionSize

	return
}

func bytesToVersions(bytes []byte) (versions []Version) {
	index := 0
	for index != len(bytes) {
		version, size := newVersion(bytes[index:])
		if size == 0 {
			break
		}

		index += int(size)
		versions = append(versions, version)
	}

	return
}

// TODO: don't use this method!
func addSize(bytes []byte) []byte {
	size := len(bytes)
	ans := make([]byte, 2)
	binary.BigEndian.PutUint16(ans, uint16(size))
	ans = append(ans, bytes...)
	return ans
}
