package scummgo

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReducePixels(t *testing.T) {
	skin := color.RGBA{200, 120, 100, 255}
	dark := color.RGBA{20, 20, 20, 255}
	none := color.RGBA{}

	t.Run("A thin sprite does not vanish", func(t *testing.T) {
		// A 6x6 sprite with a one pixel wide vertical line, reduced to 2x2.
		pixels := make([]color.RGBA, 36)
		for y := 0; y < 6; y++ {
			pixels[y*6+2] = skin
		}
		reduced := reducePixels(pixels, NewSize(6, 6), NewSize(2, 2))
		assert.Equal(t, []color.RGBA{skin, none, skin, none}, reduced)
	})

	t.Run("The most common color of the area is used", func(t *testing.T) {
		pixels := []color.RGBA{skin, skin, dark, skin}
		reduced := reducePixels(pixels, NewSize(2, 2), NewSize(1, 1))
		assert.Equal(t, []color.RGBA{skin}, reduced)
	})

	t.Run("Mostly transparent areas stay transparent", func(t *testing.T) {
		pixels := make([]color.RGBA, 16)
		pixels[0] = skin
		reduced := reducePixels(pixels, NewSize(4, 4), NewSize(1, 1))
		assert.Equal(t, []color.RGBA{none}, reduced)
	})
}
