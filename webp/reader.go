// Package webp implements a Animated WebP detector
package webp

import (
	"encoding/binary"
	"io"

	"github.com/sapphi-red/midec"
)

const webpHeader = "RIFF????WEBPVP8"

const (
	maskVP8XAnimation = 1 << 1
)

type chunkHeaderData struct {
	fourCC   string
	dataSize uint32
}

type decoder struct {
	midec.ReadAdvancer
}

func (d *decoder) readUint32() (u uint32, err error) {
	err = binary.Read(d.Reader, binary.LittleEndian, &u)
	return
}

func (d *decoder) skipHeader() error {
	return d.Advance(
		4 + // 'RIFF'
			4 + // File size
			4, // 'WEBP'
	)
}

func (d *decoder) decodeChunkHeader() (chd chunkHeaderData, err error) {
	fourCCBuf := make([]byte, 4)
	_, err = d.ReadFull(fourCCBuf)
	if err != nil {
		return
	}

	dataSize, err := d.readUint32()
	if err != nil {
		return
	}

	return chunkHeaderData{
		fourCC: string(fourCCBuf),
		dataSize: dataSize,
	}, nil
}

func (d *decoder) decodeVP8XChunk() (bool, error) {
	buf := make([]byte, 1)
	if _, err := d.ReadFull(buf); err != nil {
		return false, err
	}

	isAnimation := (buf[0]&maskVP8XAnimation) != 0

	err := d.Advance(
		3 + // Reserved
		3 + // Canvas Width Minus One
		3, // Canvas Height Minus One
	)
	if err != nil {
		return false, nil
	}

	return isAnimation, nil
}

func (d *decoder) skipThisChunk(dataSize uint32) error {
	return d.Advance(uint(dataSize))
}

func (d *decoder) decode() (bool, error) {
	if err := d.skipHeader(); err != nil {
		return false, err
	}

	chd, err := d.decodeChunkHeader()
	if err != nil {
		return false, err
	}

	if chd.fourCC != "VP8X" {
		return false, nil
	}

	isAnimation, err := d.decodeVP8XChunk()
	if err != nil {
		return false, err
	}
	if !isAnimation {
		return false, nil
	}

	frameCount := 0
	for {
		chd, err := d.decodeChunkHeader()
		if err != nil {
			if err != io.EOF {
				return false, nil
			}
			return false, err
		}

		if chd.fourCC == "ANMF" {
			frameCount++
			if frameCount >= 2 {
				return true, nil
			}
		}

		if err := d.skipThisChunk(chd.dataSize); err != nil {
			return false, err
		}
	}
}

func isAnimated(r io.Reader) (bool, error) {
	d := decoder{*midec.NewReadAdvancer(r)}
	return d.decode()
}

func init() {
	midec.RegisterFormat("webp", webpHeader, isAnimated)
}
