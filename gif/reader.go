// Package gif implements a Animated GIF detector
package gif

import (
	"errors"
	"io"

	"github.com/sapphi-red/midec"
)

const gifHeader = "GIF8?a"

// ErrUnknownBlock indicates that detecting encountered an unknown block.
var ErrUnknownBlock = errors.New("midec: (gif) unknown block")

const (
	maskGlobalColorTableFlag = 1 << 7
	maskSizeGlobalColorTable = 0b111

	maskLocalColorTableFlag = 1 << 7
	maskSizeLocalColorTable = 0b111
)

type colorTableData struct {
	flag bool // Global/Local Color Table Flag
	size uint // Size of Global/Local Color Table
}

type blockType uint

const (
	blockTypeUnknown blockType = iota
	blockTypeTerminator
	blockTypeImageBlock
	blockTypeGraphicControlExtension
	blockTypeCommentExtension
	blockTypePlainTextExtension
	blockTypeApplicationExtension
)

type decoder struct {
	midec.ReadAdvancer
}

func (d *decoder) readOneByte() (byte, error) {
	buf := make([]byte, 1)
	if _, err := d.ReadFull(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (d *decoder) skipHeader() error {
	err := d.Advance(
		3 + // Signature
			3 + // Version
			2 + // Logical Screen Width
			2, // Logical Screen Height
	)
	if err != nil {
		return err
	}

	gctd, err := d.decodeHeaderPackedFields()
	if err != nil {
		return err
	}

	gctdLen := uint(0)
	if gctd.flag {
		gctdLen = 3 * (1 << (1 + gctd.size))
	}

	err = d.Advance(
		1 + // Background Color Index
			1 + // Pixel Aspect Ratio
			gctdLen, // Global Color Table
	)
	return err
}

func (d *decoder) decodeHeaderPackedFields() (ctd colorTableData, err error) {
	packedFields, err := d.readOneByte()
	if err != nil {
		return
	}

	ctd.flag = packedFields&maskGlobalColorTableFlag != 0
	ctd.size = uint(packedFields & maskSizeGlobalColorTable)
	return
}

func (d *decoder) parseBlockType() (blockType, error) {
	imageSeparator, err := d.readOneByte()
	if err != nil {
		return blockTypeUnknown, err
	}

	if imageSeparator == 0x2c {
		return blockTypeImageBlock, nil
	}
	if imageSeparator == 0x3b {
		return blockTypeTerminator, nil
	}
	if imageSeparator != 0x21 {
		return blockTypeUnknown, ErrUnknownBlock
	}

	typeId, err := d.readOneByte()
	if err != nil {
		return blockTypeUnknown, err
	}

	switch typeId {
	case 0xf9:
		return blockTypeGraphicControlExtension, nil
	case 0xfe:
		return blockTypeCommentExtension, nil
	case 0x01:
		return blockTypePlainTextExtension, nil
	case 0xff:
		return blockTypeApplicationExtension, nil
	}
	return blockTypeUnknown, ErrUnknownBlock
}

func (d *decoder) skipImageBlock() error {
	err := d.Advance(
		2 + // Image Left Position
			2 + // Image Top Position
			2 + // Image Width
			2, // Image Height
	)
	if err != nil {
		return err
	}

	lctd, err := d.decodeImageBlockPackedFields()
	if err != nil {
		return err
	}

	lctdLen := uint(0)
	if lctd.flag {
		lctdLen = 3 * (1 << (1 + lctd.size))
	}

	err = d.Advance(
		lctdLen + // Local Color Table
			1, // LZW Minimum Code Size
	)
	if err != nil {
		return err
	}

	for {
		blockSize, err := d.decodeImageBlockBlockSize()
		if err != nil {
			return err
		}
		// when block terminator
		if blockSize == 0 {
			break
		}

		err = d.Advance(
			blockSize, // Image Data
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *decoder) skipGraphicControlExtensionBlock() error {
	err := d.Advance(
		1 + // Block Size
			1 + // Packed Fields
			2 + // Delay Time
			1 + // Transparent Color Index
			1, // Block Terminator
	)
	return err
}

func (d *decoder) skipCommentExtensionBlock() error {
	blockSize, err := d.readOneByte()
	if err != nil {
		return err
	}

	err = d.Advance(
		uint(blockSize) + // Comment Data
			1, // Block Terminator
	)
	return err
}

func (d *decoder) skipPlainTextExtensionBlock() error {
	err := d.Advance(
		1 + // Block Size
			2 + // Text Grid Left Position
			2 + // Text Grid Top Position
			2 + // Text Grid Width
			2 + // Text Grid Height
			1 + // Character Cell Width
			1 + // Character Cell Height
			1 + // Text Foreground Color Index
			1, // Text Background Color Index
	)
	if err != nil {
		return err
	}

	blockSize, err := d.readOneByte()
	if err != nil {
		return err
	}

	err = d.Advance(
		uint(blockSize) + // Plain Text Data
			1, // Block Terminator
	)
	return err
}

func (d *decoder) skipApplicationExtensionBlock() error {
	err := d.Advance(
		1 + // Block Size
			8 + // Application Identifier
			3, // Application Authentication Code
	)
	if err != nil {
		return err
	}

	blockSize, err := d.readOneByte()
	if err != nil {
		return err
	}

	err = d.Advance(
		uint(blockSize) + // Application Data
			1, // Block Terminator
	)
	return err
}

func (d *decoder) decodeImageBlockPackedFields() (ctd colorTableData, err error) {
	packedFields, err := d.readOneByte()
	if err != nil {
		return
	}

	ctd.flag = packedFields&maskLocalColorTableFlag != 0
	ctd.size = uint(packedFields & maskSizeLocalColorTable)
	return
}

func (d *decoder) decodeImageBlockBlockSize() (uint, error) {
	sizeByte, err := d.readOneByte()
	if err != nil {
		return 0, err
	}

	return uint(sizeByte), nil
}

func (d *decoder) decodeBlocks() (bool, error) {
	imageBlockCount := 0

	for {
		blockType, err := d.parseBlockType()
		if err != nil {
			return false, err
		}

		switch blockType {
		case blockTypeTerminator:
			return imageBlockCount >= 2, nil

		case blockTypeImageBlock:
			if err := d.skipImageBlock(); err != nil {
				return false, err
			}
			imageBlockCount++
			if imageBlockCount >= 2 {
				return true, nil
			}

		case blockTypeGraphicControlExtension:
			if err := d.skipGraphicControlExtensionBlock(); err != nil {
				return false, err
			}

		case blockTypeCommentExtension:
			if err := d.skipCommentExtensionBlock(); err != nil {
				return false, err
			}

		case blockTypePlainTextExtension:
			if err := d.skipPlainTextExtensionBlock(); err != nil {
				return false, err
			}

		case blockTypeApplicationExtension:
			if err := d.skipApplicationExtensionBlock(); err != nil {
				return false, err
			}

		}
	}
}

func (d *decoder) decode() (bool, error) {
	if err := d.skipHeader(); err != nil {
		return false, err
	}

	return d.decodeBlocks()
}

func isAnimated(r io.Reader) (bool, error) {
	d := decoder{*midec.NewReadAdvancer(r)}
	return d.decode()
}

func init() {
	midec.RegisterFormat("gif", gifHeader, isAnimated)
}
