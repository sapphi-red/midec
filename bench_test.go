package midec_test

import (
	"image/gif"
	"os"
	"testing"

	"github.com/sapphi-red/midec"
	_ "github.com/sapphi-red/midec/gif"
)

func loadGif() *os.File {
	fp, err := os.Open(testdataFolder + "gif/1.gif")
	if err != nil {
		panic(err)
	}
	return fp
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
