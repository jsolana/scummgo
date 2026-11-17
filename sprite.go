package scummgo

import (
	"image/color"
	"io"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SpriteSheet represents a collection of sprites arranged in a grid-shaped sheet.
type SpriteSheet struct {
	raw       *rl.Image
	tex       rl.Texture2D
	frameSize Size

	reduced map[reducedSprite]rl.Texture2D
}

// reducedSprite identifies a sprite drawn at a lower resolution than its final size.
type reducedSprite struct {
	col, row uint
	size     Size
}

// pixelationRate controls how fast the pixels of a scaled sprite grow when it gets smaller.
const pixelationRate = 1.5

// pixelBlock returns the size in pixels of the blocks used to draw a sprite scaled by the given
// factor: 1 for sprites that are close, bigger for sprites that are far. The farther a sprite is,
// the less resolution it has, so it looks coarser instead of sharper.
func pixelBlock(scale float32) int {
	return 1 + max(0, int((1-scale)*pixelationRate))
}

// LoadSpriteSheetFromFile loads a sprite sheet from a image file.
func LoadSpriteSheetFromFile(path string, frameSize Size) *SpriteSheet {
	return &SpriteSheet{
		raw:       rl.LoadImage(path),
		frameSize: frameSize,
	}
}

// Release releases the resources used by the sprite sheet.
func (s *SpriteSheet) Release() {
	rl.UnloadTexture(s.tex)
	for _, tex := range s.reduced {
		rl.UnloadTexture(tex)
	}
	s.reduced = nil
}

// DrawSprite draws a sprite from the sprite sheet at the given position, scaled by the given
// factor. The scaled size is rounded to whole pixels and the texture is sampled without filtering.
// Scaled down sprites are drawn with bigger pixels the smaller they are (see pixelBlock), keeping
// their bottom center in place.
func (s *SpriteSheet) DrawSprite(col, row uint, pos Position, flip bool, scale float32) {
	src := Rectangle{
		Pos: Position{
			int(s.frameSize.W) * int(col),
			int(s.frameSize.H) * int(row),
		},
		Size: s.frameSize,
	}
	if flip {
		src.Size = src.Size.FlipH()
	}
	dst := Rectangle{Pos: pos, Size: s.frameSize.Scale(scale)}
	tex := s.texture()
	if block := pixelBlock(scale); block > 1 {
		// Draw a reduced version of the sprite, with one pixel per block, scaled up to the blocks.
		reduced := NewSize(max(1, dst.Size.W/block), max(1, dst.Size.H/block))
		tex = s.reducedTexture(col, row, reduced)
		src = Rectangle{Size: reduced}
		if flip {
			src.Size = src.Size.FlipH()
		}
		size := NewSize(reduced.W*block, reduced.H*block)
		dst.Pos = pos.Sub(NewPos((size.W-dst.Size.W)/2, size.H-dst.Size.H))
		dst.Size = size
	}
	rl.DrawTexturePro(tex, src.toRaylib(), dst.toRaylib(), rl.Vector2{}, 0, rl.White)
}

// BinaryEncode encodes the sprite sheet to a binary format. The encoded format is:
// - uint16: the width of each sprite.
// - uint16: the height of each sprite.
// - uint32: the length of the image bytes.
// - []byte: the image bytes in PNG format.
func (s *SpriteSheet) BinaryEncode(w io.Writer) (int, error) {
	bytes := rl.ExportImageToMemory(*s.raw, ".png")
	return BinaryEncode(w, uint16(s.frameSize.W), uint16(s.frameSize.H), uint32(len(bytes)), bytes)
}

// BinaryDecode decodes the sprite sheet from a binary format. See SpriteSheet.BinaryEncode for the
// format.
func (s *SpriteSheet) BinaryDecode(r io.Reader) error {
	var w, h uint16
	var size uint32
	if err := BinaryDecode(r, &w, &h, &size); err != nil {
		return err
	}
	bytes := make([]byte, size)
	if err := BinaryDecode(r, bytes); err != nil {
		return err
	}
	s.frameSize = Size{int(w), int(h)}
	s.raw = rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
	return nil
}

// reducePixels resizes the pixels of an image of size from to the smaller size to. Each resulting
// pixel takes the most common color of the area of the original image it covers, and it is
// transparent when less than a third of that area is opaque. Picking a single pixel instead would
// make thin parts of the sprite, or the whole sprite, vanish when it gets small.
func reducePixels(pixels []color.RGBA, from, to Size) []color.RGBA {
	reduced := make([]color.RGBA, to.W*to.H)
	for y := 0; y < to.H; y++ {
		for x := 0; x < to.W; x++ {
			x0, x1 := x*from.W/to.W, max((x+1)*from.W/to.W, x*from.W/to.W+1)
			y0, y1 := y*from.H/to.H, max((y+1)*from.H/to.H, y*from.H/to.H+1)

			counts := make(map[color.RGBA]int)
			var best color.RGBA
			opaque := 0
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					c := pixels[sy*from.W+sx]
					if c.A < 128 {
						continue
					}
					opaque++
					counts[c]++
					if counts[c] > counts[best] {
						best = c
					}
				}
			}
			if opaque*3 >= (x1-x0)*(y1-y0) {
				reduced[y*to.W+x] = best
			}
		}
	}
	return reduced
}

// reducedTexture returns the sprite at the given column and row resized to the given size.
func (s *SpriteSheet) reducedTexture(col, row uint, size Size) rl.Texture2D {
	key := reducedSprite{col, row, size}
	if tex, ok := s.reduced[key]; ok {
		return tex
	}

	frame := rl.ImageFromImage(*s.raw, rl.NewRectangle(
		float32(s.frameSize.W*int(col)), float32(s.frameSize.H*int(row)),
		float32(s.frameSize.W), float32(s.frameSize.H),
	))
	pixels := rl.LoadImageColors(&frame)
	reducedPixels := reducePixels(pixels, s.frameSize, size)
	rl.UnloadImageColors(pixels)
	rl.UnloadImage(&frame)

	reduced := rl.GenImageColor(size.W, size.H, rl.Blank)
	for i, c := range reducedPixels {
		if c.A != 0 {
			rl.ImageDrawPixel(reduced, int32(i%size.W), int32(i/size.W), c)
		}
	}
	tex := rl.LoadTextureFromImage(reduced)
	rl.UnloadImage(reduced)

	if s.reduced == nil {
		s.reduced = make(map[reducedSprite]rl.Texture2D)
	}
	s.reduced[key] = tex
	return tex
}

func (s *SpriteSheet) texture() rl.Texture2D {
	if !rl.IsTextureReady(s.tex) {
		s.tex = rl.LoadTextureFromImage(s.raw)
	}
	return s.tex
}
