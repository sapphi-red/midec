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
	return loadFile("gif/animated.gif")
}

func BenchmarkGIF_ImageGIF(b *testing.B) {
	fp := loadGif()
	for i := 0; i < b.N; i++ {
		decoded, err := gif.DecodeAll(fp)
		if err != nil {
			panic(err)
		}
		_ = len(decoded.Image)

		b.StopTimer()
		_, _ = fp.Seek(0, 0)
		b.StartTimer()
	}
}

func BenchmarkGIF_Midec(b *testing.B) {
	fp := loadGif()
	for i := 0; i < b.N; i++ {
		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}

		b.StopTimer()
		_, _ = fp.Seek(0, 0)
		b.StartTimer()
	}
}

func BenchmarkPNG_Midec(b *testing.B) {
	fp := loadFile("png/animated.png")
	for i := 0; i < b.N; i++ {
		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}

		b.StopTimer()
		_, _ = fp.Seek(0, 0)
		b.StartTimer()
	}
}

func BenchmarkWebP_Midec(b *testing.B) {
	fp := loadFile("webp/animated.webp")
	for i := 0; i < b.N; i++ {
		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}

		b.StopTimer()
		_, _ = fp.Seek(0, 0)
		b.StartTimer()
	}
}

func BenchmarkHEIFAVIF_Midec(b *testing.B) {
	fp := loadFile("isobmff/animated.avif")
	for i := 0; i < b.N; i++ {
		_, err := midec.IsAnimated(fp)
		if err != nil {
			panic(err)
		}

		b.StopTimer()
		_, _ = fp.Seek(0, 0)
		b.StartTimer()
	}
}
