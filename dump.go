package scummgo

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// PackageDump is a stable, JSON-serializable description of a ScummGo package.
// Both the Go and Rust dump tools emit this shape so dumps can be compared byte-for-byte
// (after canonical JSON encoding).
type PackageDump struct {
	FormatVersion uint16             `json:"format_version"`
	IndexMagic    string             `json:"index_magic"`
	DataMagic     string             `json:"data_magic"`
	Resources     []ResourceDump     `json:"resources"`
}

// ResourceDump describes one indexed resource and a decoded summary of its payload.
type ResourceDump struct {
	ID               string `json:"id"`
	Offset           uint32 `json:"offset"`
	Size             uint32 `json:"size"`
	Type             string `json:"type"`
	TypeCode         uint8  `json:"type_code"`
	Compression      string `json:"compression"`
	CompressionCode  uint8  `json:"compression_code"`
	PayloadSHA256    string `json:"payload_sha256"`
	DecodedSHA256    string `json:"decoded_sha256"`
	DecodedSize      int    `json:"decoded_size"`
	Detail           any    `json:"detail"`
}

// ScriptDumpDetail is the dump detail for a script resource.
type ScriptDumpDetail struct {
	Language     string `json:"language"`
	LanguageCode uint8  `json:"language_code"`
	CodeSHA256   string `json:"code_sha256"`
	CodeSize     int    `json:"code_size"`
	CodePreview  string `json:"code_preview,omitempty"`
}

// ImageDumpDetail is the dump detail for an image resource.
type ImageDumpDetail struct {
	PngSize   int    `json:"png_size"`
	PngSHA256 string `json:"png_sha256"`
}

// SpriteSheetDumpDetail is the dump detail for a sprite sheet resource.
type SpriteSheetDumpDetail struct {
	FrameWidth  uint16 `json:"frame_width"`
	FrameHeight uint16 `json:"frame_height"`
	PngSize     int    `json:"png_size"`
	PngSHA256   string `json:"png_sha256"`
}

// AudioDumpDetail is the dump detail for sound/music resources.
type AudioDumpDetail struct {
	Format    string `json:"format"`
	DataSize  int    `json:"data_size"`
	DataSHA256 string `json:"data_sha256"`
}

// CostumeDumpDetail is the dump detail for a costume resource.
type CostumeDumpDetail struct {
	Sprites     SpriteSheetDumpDetail `json:"sprites"`
	Animations  []AnimationDumpDetail `json:"animations"`
}

// AnimationDumpDetail is one costume animation entry.
type AnimationDumpDetail struct {
	Action uint8               `json:"action"`
	Flip   bool                `json:"flip"`
	Frames []FrameDumpDetail   `json:"frames"`
}

// FrameDumpDetail is one animation frame.
type FrameDumpDetail struct {
	Col     uint8  `json:"col"`
	Row     uint8  `json:"row"`
	DelayNs uint64 `json:"delay_ns"`
}

// DumpPackageFile loads {base}.idx / {base}.dat and returns a package dump.
// base may include a directory and must not include the .idx/.dat suffix.
func DumpPackageFile(base string) (*PackageDump, error) {
	idxPath := base + ".idx"
	datPath := base + ".dat"

	idxData, err := os.ReadFile(idxPath)
	if err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}
	datData, err := os.ReadFile(datPath)
	if err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}
	return DumpPackageBytes(idxData, datData)
}

// DumpPackageBytes inspects in-memory index/data file contents.
func DumpPackageBytes(idxData, datData []byte) (*PackageDump, error) {
	idx := bytes.NewReader(idxData)
	var idxHdr resourceFileHeader
	if err := BinaryDecode(idx, &idxHdr); err != nil {
		return nil, fmt.Errorf("decode index header: %w", err)
	}
	if !bytes.Equal(idxHdr.Magic[:], resourceIndexMagic[:]) {
		return nil, fmt.Errorf("bad index magic: %q", idxHdr.Magic)
	}

	dat := bytes.NewReader(datData)
	var datHdr resourceFileHeader
	if err := BinaryDecode(dat, &datHdr); err != nil {
		return nil, fmt.Errorf("decode data header: %w", err)
	}
	if !bytes.Equal(datHdr.Magic[:], resourceDataMagic[:]) {
		return nil, fmt.Errorf("bad data magic: %q", datHdr.Magic)
	}

	dump := &PackageDump{
		FormatVersion: idxHdr.Version,
		IndexMagic:    string(idxHdr.Magic[:]),
		DataMagic:     string(datHdr.Magic[:]),
	}

	for {
		var entry indexEntry
		if err := BinaryDecode(idx, &entry); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode index entry: %w", err)
		}
		rd, err := dumpResource(datData, entry)
		if err != nil {
			return nil, fmt.Errorf("resource %s: %w", entry.ID, err)
		}
		dump.Resources = append(dump.Resources, *rd)
	}

	sort.Slice(dump.Resources, func(i, j int) bool {
		return dump.Resources[i].ID < dump.Resources[j].ID
	})
	return dump, nil
}

// DumpPackageJSON returns canonical JSON for a package dump (stable key order via struct tags + sorted resources).
func DumpPackageJSON(dump *PackageDump) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(dump); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func dumpResource(datData []byte, entry indexEntry) (*ResourceDump, error) {
	if int(entry.Offset)+int(entry.Size) > len(datData) {
		return nil, fmt.Errorf("out of bounds offset=%d size=%d file=%d", entry.Offset, entry.Size, len(datData))
	}
	blob := datData[entry.Offset : entry.Offset+entry.Size]
	if len(blob) < resourceHeaderSize {
		return nil, fmt.Errorf("resource blob too small: %d", len(blob))
	}

	r := bytes.NewReader(blob)
	var h resourceHeader
	if err := BinaryDecode(r, &h); err != nil {
		return nil, err
	}
	payload := blob[resourceHeaderSize:]
	decoded, err := decompressPayload(h.Compression, payload)
	if err != nil {
		return nil, err
	}

	typeName, detail, err := inspectDecoded(h.Type, decoded)
	if err != nil {
		return nil, err
	}

	return &ResourceDump{
		ID:              entry.ID.String(),
		Offset:          entry.Offset,
		Size:            entry.Size,
		Type:            typeName,
		TypeCode:        uint8(h.Type),
		Compression:     compressionName(h.Compression),
		CompressionCode: uint8(h.Compression),
		PayloadSHA256:   sha256Hex(payload),
		DecodedSHA256:   sha256Hex(decoded),
		DecodedSize:     len(decoded),
		Detail:          detail,
	}, nil
}

func decompressPayload(c ResourceCompression, payload []byte) ([]byte, error) {
	switch c {
	case CompressionNone:
		return payload, nil
	case CompressionGzip:
		zr, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		defer zr.Close()
		return io.ReadAll(zr)
	default:
		return nil, fmt.Errorf("unsupported compression %d", c)
	}
}

func inspectDecoded(t resourceType, decoded []byte) (string, any, error) {
	r := bytes.NewReader(decoded)
	switch t {
	case resourceTypeScript:
		var lang ScriptLanguage
		var length uint32
		if err := BinaryDecode(r, &lang, &length); err != nil {
			return "", nil, err
		}
		code := make([]byte, length)
		if err := BinaryDecode(r, &code); err != nil {
			return "", nil, err
		}
		preview := string(code)
		if len(preview) > 80 {
			preview = preview[:80] + "..."
		}
		langName := "undefined"
		if lang == ScriptLua {
			langName = "lua"
		}
		return "script", ScriptDumpDetail{
			Language:     langName,
			LanguageCode: uint8(lang),
			CodeSHA256:   sha256Hex(code),
			CodeSize:     len(code),
			CodePreview:  preview,
		}, nil

	case resourceTypeImage:
		var length uint32
		if err := BinaryDecode(r, &length); err != nil {
			return "", nil, err
		}
		png := make([]byte, length)
		if err := BinaryDecode(r, &png); err != nil {
			return "", nil, err
		}
		return "image", ImageDumpDetail{PngSize: len(png), PngSHA256: sha256Hex(png)}, nil

	case resourceTypeSpriteSheet:
		detail, err := decodeSpriteSheetDump(r)
		if err != nil {
			return "", nil, err
		}
		return "spritesheet", detail, nil

	case resourceTypeSound:
		detail, err := decodeAudioDump(r)
		if err != nil {
			return "", nil, err
		}
		return "sound", detail, nil

	case resourceTypeMusic:
		detail, err := decodeAudioDump(r)
		if err != nil {
			return "", nil, err
		}
		return "music", detail, nil

	case resourceTypeCostume:
		sprites, err := decodeSpriteSheetDump(r)
		if err != nil {
			return "", nil, err
		}
		var count uint32
		if err := BinaryDecode(r, &count); err != nil {
			return "", nil, err
		}
		anims := make([]AnimationDumpDetail, 0, count)
		for i := uint32(0); i < count; i++ {
			var act byte
			var flip bool
			var frameCount uint32
			if err := BinaryDecode(r, &act, &flip, &frameCount); err != nil {
				return "", nil, err
			}
			frames := make([]FrameDumpDetail, frameCount)
			for j := uint32(0); j < frameCount; j++ {
				var col, row byte
				var delay uint64
				if err := BinaryDecode(r, &col, &row, &delay); err != nil {
					return "", nil, err
				}
				frames[j] = FrameDumpDetail{Col: col, Row: row, DelayNs: delay}
			}
			anims = append(anims, AnimationDumpDetail{Action: act, Flip: flip, Frames: frames})
		}
		return "costume", CostumeDumpDetail{Sprites: *sprites, Animations: anims}, nil

	default:
		return fmt.Sprintf("unknown(%d)", t), map[string]any{"raw_size": len(decoded)}, nil
	}
}

func decodeSpriteSheetDump(r io.Reader) (*SpriteSheetDumpDetail, error) {
	var w, h uint16
	var size uint32
	if err := BinaryDecode(r, &w, &h, &size); err != nil {
		return nil, err
	}
	png := make([]byte, size)
	if err := BinaryDecode(r, &png); err != nil {
		return nil, err
	}
	return &SpriteSheetDumpDetail{
		FrameWidth:  w,
		FrameHeight: h,
		PngSize:     len(png),
		PngSHA256:   sha256Hex(png),
	}, nil
}

func decodeAudioDump(r io.Reader) (*AudioDumpDetail, error) {
	var format [4]byte
	var length uint32
	if err := BinaryDecode(r, &format, &length); err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if err := BinaryDecode(r, &data); err != nil {
		return nil, err
	}
	return &AudioDumpDetail{
		Format:     string(bytes.TrimRight(format[:], "\x00")),
		DataSize:   len(data),
		DataSHA256: sha256Hex(data),
	}, nil
}

func compressionName(c ResourceCompression) string {
	switch c {
	case CompressionNone:
		return "none"
	case CompressionGzip:
		return "gzip"
	default:
		return fmt.Sprintf("unknown(%d)", c)
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// DumpPackageToFile writes the JSON dump of base.{idx,dat} to outPath.
func DumpPackageToFile(base, outPath string) error {
	dump, err := DumpPackageFile(base)
	if err != nil {
		return err
	}
	data, err := DumpPackageJSON(dump)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil && filepath.Dir(outPath) != "." {
		return err
	}
	return os.WriteFile(outPath, data, 0o644)
}
