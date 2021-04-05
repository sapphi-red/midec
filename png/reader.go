// Package png implements a APNG detector
package png

import (
	"encoding/binary"
	"io"

	"github.com/sapphi-red/midec"
)

const pngHeader = "\x89PNG\r\n\x1a\n"

type chunkHeaderData struct {
	length uint32
	typeId string
}

type decoder struct {
	midec.ReadAdvancer
}

func (d *decoder) readUint32() (u uint32, err error) {
	err = binary.Read(d.Reader, binary.BigEndian, &u)
	return
}

func (d *decoder) skipHeader() error {
	return d.Advance(8)
}

func (d *decoder) decodeChunkHeader() (chd chunkHeaderData, err error) {
	length, err := d.readUint32()
	if err != nil {
		return
	}

	typeIdBuf := make([]byte, 4)
	_, err = d.ReadFull(typeIdBuf)
	if err != nil {
		return
	}

	return chunkHeaderData{
		length: length,
		typeId: string(typeIdBuf),
	}, nil
}

func (d *decoder) decodeacTLChunk() (bool, error) {
	length, err := d.readUint32()
	if err != nil {
		return false, err
	}

	return length >= 2, nil
}

func (d *decoder) skipUnknownChunk(length uint32) error {
	return d.Advance(
		uint(length) + // data
		4, // CRC
	)
}

func (d *decoder) decode() (bool, error) {
	err := d.skipHeader()
	if err != nil {
		return false, err
	}

	for {
		chd, err := d.decodeChunkHeader()
		if err != nil {
			return false, err
		}

		switch chd.typeId {
		case "acTL":
			return d.decodeacTLChunk()
		case "IDAT":
			// acTL chunk must come before IDAT.
			// so if IDAT comes before acTL, it is not a apng.
			return false, nil
		default:
			err := d.skipUnknownChunk(chd.length)
			if err != nil {
				return false, err
			}
		}
	}
}

func isAnimated(r io.Reader) (bool, error) {
	d := decoder{*midec.NewReadAdvancer(r)}
	return d.decode()
}

func init() {
	midec.RegisterFormat("png", pngHeader, isAnimated)
}
