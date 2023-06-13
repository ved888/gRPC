package sample

import (
	"github.com/google/uuid"
	"grpc/pb/pb"
	"math/rand"
	"time"
)

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}
func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomCPUBrand() string {
	return randomStringFromSet("intel", "AMD")
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"RTX 2080",
		)
	}
	return randomStringFromSet(
		"RX 590",
		"RX 580",
	)
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomCPUName(brand string) string {
	if brand == "intel" {
		return randomStringFromSet(
			"xeon E 22864",
		)
	}
	return randomStringFromSet(
		"Regen 7",
	)
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	resolution := &pb.Screen_Resolution{
		Height: uint32(height),
		Width:  uint32(width),
	}
	return resolution
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomId() string {
	return uuid.NewString()
}

func randomLaptopBrand() string {
	return randomStringFromSet("dell", "lenovo", "Apple")
}
func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("macbook Air", "macbook pro")
	case "dell":
		return randomStringFromSet("latitude", "xps")
	default:
		return randomStringFromSet("thin pad", "thin pad x1")

	}
}
