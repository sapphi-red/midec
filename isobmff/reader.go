// Package isobmff implements a Animated HEIF / Animated AVIF detector
package isobmff

import (
	"encoding/binary"
	"io"

	"github.com/sapphi-red/midec"
)

const isobmmfHeader = "????ftyp"

var animatedableBrands = []string{
	"mif1", // HEIF: structural brand image
	"msf1", // HEIF: structural brand image sequence
	"heic", // HEIF HEVC: general image
	"heix", // HEIF HEVC: general image
	"hevc", // HEIF HEVC: general image sequence
	"hevx", // HEIF HEVC: general image sequence
	"heim", // HEIF HEVC: multiview image
	"heis", // HEIF HEVC: scalable image
	"hevm", // HEIF HEVC: multiview image sequence
	"hevs", // HEIF HEVC: scalable image sequence
	"avci", // HEIF AVC: general image
	"avcs", // HEIF AVC: general image sequence
	"avif", // AVIF(HEIF AV1): image
	"avis", // AVIF(HEIF AV1): image sequence
}

type boxHeaderData struct {
	dataSize int64
	untilEnd bool
	boxType  string
}

type movieHeaderBoxData struct {
	duration int64
}

type handlerReferenceBoxData struct {
	handlerType string
}

type decoder struct {
	midec.ReadAdvancer
}

func (d *decoder) read(data interface{}) error {
	return binary.Read(d.Reader, binary.BigEndian, data)
}

func (d *decoder) decodeBoxHeader() (bhd boxHeaderData, err error) {
	var size int32
	err = d.read(&size)
	if err != nil {
		return
	}

	boxTypeBuf := make([]byte, 4)
	_, err = d.ReadFull(boxTypeBuf)
	if err != nil {
		return
	}

	boxType := string(boxTypeBuf)

	if size == 0 {
		return boxHeaderData{
			dataSize: 0,
			untilEnd: true,
			boxType:  boxType,
		}, nil
	}

	if size != 1 {
		return boxHeaderData{
			dataSize: int64(size - 4 - 4),
			untilEnd: false,
			boxType:  boxType,
		}, nil
	}

	var largeSize int64
	err = d.read(&largeSize)
	if err != nil {
		return
	}

	return boxHeaderData{
		dataSize: largeSize - 4 - 4 - 8,
		untilEnd: false,
		boxType:  boxType,
	}, nil
}

func (d *decoder) findBox(targetBoxType string, dataSizeP *int64) (bool, boxHeaderData, error) {
	dataSize := int64(0)
	if dataSizeP != nil {
		dataSize = *dataSizeP
	}

	for {
		bhd, err := d.decodeBoxHeader()
		if err != nil {
			return false, boxHeaderData{}, err
		}
		dataSize -= 4 + 4 + bhd.dataSize

		if bhd.boxType == targetBoxType {
			return true, bhd, nil
		}

		if bhd.untilEnd {
			return false, boxHeaderData{}, nil
		}

		if err := d.Advance(uint(bhd.dataSize)); err != nil {
			return false, boxHeaderData{}, err
		}

		if dataSizeP != nil && dataSize <= 0 {
			return false, boxHeaderData{}, err
		}
	}
}

func (d *decoder) decodeFileTypeBox() (bool, error) {
	var size int32
	if err := d.read(&size); err != nil {
		return false, err
	}

	err := d.Advance(
		4, // type
	)
	if err != nil {
		return false, err
	}

	brandBuff := make([]byte, 4)
	if _, err = d.ReadFull(brandBuff); err != nil {
		return false, err
	}

	err = d.Advance(uint(size) - 4 - 4 - 4)
	if err != nil {
		return false, err
	}

	brand := string(brandBuff)

	for _, b := range animatedableBrands {
		if brand == b {
			return true, nil
		}
	}
	return false, nil
}

func (d *decoder) decodeMovieHeaderBox(dataSize int64) (mhbd movieHeaderBoxData, err error) {
	var version byte
	err = d.read(&version)
	if err != nil {
		return
	}

	if version == 1 {
		err = d.Advance(
			8 + // creation_time
				8 + // modification_time
				4, // timescale
		)
		if err != nil {
			return
		}

		var duration int64
		err = d.read(&duration)
		if err != nil {
			return
		}

		err = d.Advance(uint(dataSize) - 1 - 8 - 8 - 4 - 8)
		if err != nil {
			return
		}

		return movieHeaderBoxData{
			duration: duration,
		}, nil
	}

	err = d.Advance(
		4 + // creation_time
			4 + // modification_time
			4, // timescale
	)
	if err != nil {
		return
	}

	var duration int32
	err = d.read(&duration)
	if err != nil {
		return
	}

	err = d.Advance(uint(dataSize) - 1 - 4 - 4 - 4 - 4)
	if err != nil {
		return
	}

	return movieHeaderBoxData{
		duration: int64(duration),
	}, nil
}

func (d *decoder) decodeHandlerReferenceBox(dataSize int64) (hrbd handlerReferenceBoxData, err error) {
	err = d.Advance(
		1 + // (FullBox) version
			3 + // (FullBox) flags
			4, // pre_defined
	)
	if err != nil {
		return
	}

	handlerTypeBuff := make([]byte, 4)
	_, err = d.ReadFull(handlerTypeBuff)
	if err != nil {
		return
	}

	err = d.Advance(uint(dataSize) - 1 - 3 - 1 - 4)
	if err != nil {
		return
	}

	return handlerReferenceBoxData{
		handlerType: string(handlerTypeBuff),
	}, nil
}

func (d *decoder) confirmTrackBoxIsPictureType(dataSize int64) (bool, error) {
	found, mdiaD, err := d.findBox("mdia", &dataSize)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	found, hdlrD, err := d.findBox("hdlr", &mdiaD.dataSize)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	hrbd, err := d.decodeHandlerReferenceBox(hdlrD.dataSize)
	if err != nil {
		return false, err
	}

	isPictureType := hrbd.handlerType == "pict"
	return isPictureType, nil
}

func (d *decoder) decode() (bool, error) {
	isAnimatable, err := d.decodeFileTypeBox()
	if err != nil {
		return false, err
	}
	if !isAnimatable {
		return false, nil
	}

	found, moovhd, err := d.findBox("moov", nil)
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	if !found {
		return false, nil
	}

	moovSize := moovhd.dataSize

	hasValidDuration := false
	hasPictTrak := false
	for {
		bhd, err := d.decodeBoxHeader()
		if err != nil {
			return false, err
		}
		moovSize -= 4 + 4 + bhd.dataSize

		switch bhd.boxType {
		case "mvhd":
			mhbd, err := d.decodeMovieHeaderBox(bhd.dataSize)
			if err != nil {
				return false, err
			}

			hasValidDuration = mhbd.duration > 0
			if !hasValidDuration {
				return false, nil
			}
		case "trak":
			isPictureType, err := d.confirmTrackBoxIsPictureType(bhd.dataSize)
			if err != nil {
				return false, err
			}

			if isPictureType {
				hasPictTrak = isPictureType
			}
		default:
			if err := d.Advance(uint(bhd.dataSize)); err != nil {
				return false, err
			}
		}

		if hasValidDuration && hasPictTrak {
			return true, nil
		}
		if moovSize <= 0 {
			return false, nil
		}
	}
}

func isAnimated(r io.Reader) (bool, error) {
	d := decoder{*midec.NewReadAdvancer(r)}
	return d.decode()
}

func init() {
	midec.RegisterFormat("isobmff", isobmmfHeader, isAnimated)
}
