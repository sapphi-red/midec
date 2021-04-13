package midec_test

import (
	"image/gif"
	"os"
	"testing"

	"github.com/sapphi-red/midec"
	_ "github.com/sapphi-red/midec/gif"
	_ "github.com/sapphi-red/midec/png"
	_ "github.com/sapphi-red/midec/webp"
	_ "github.com/sapphi-red/midec/isobmff"
)

func loadFile(file string) *os.File {
	fp, err := os.Open(testdataFolder + file)
	if err != nil {
		panic(err)
	}
	return fp
}

func loadGif() *os.File {
	return loadFile("gif/1.gif")
}

func BenchmarkGIF_ImageGIF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fp := loadGif()
		b.StartTimer()

		decoded, err := gif.DecodeAll(fp)
		if err != nil {
			panic(err)
		}
		_ = len(decoded.Image)
	}
}

func BenchmarkGIF_Midec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fp := loadGif()
		b.StartTimer()

		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkPNG_Midec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fp := loadFile("png/1.png")
		b.StartTimer()

		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkWebP_Midec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fp := loadFile("webp/1.webp")
		b.StartTimer()

		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkHEIFAVIF_Midec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fp := loadFile("isobmff/1.avif")
		b.StartTimer()

		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}
	}
}
