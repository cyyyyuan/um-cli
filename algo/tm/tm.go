package tm

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/unlock-music/cli/algo/common"
)

var replaceHeader = []byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70}
var magicHeader = []byte{0x51, 0x51, 0x4D, 0x55} //0x15, 0x1D, 0x1A, 0x21

type Decoder struct {
	raw    io.ReadSeeker
	offset int
	audio  io.Reader
}

func (d *Decoder) Validate() error {
	header := make([]byte, 8)
	if _, err := io.ReadFull(d.raw, header); err != nil {
		return fmt.Errorf("tm read header: %w", err)
	}
	if !bytes.Equal(magicHeader, header[:len(magicHeader)]) {
		return errors.New("tm: valid magic header")
	}

	d.audio = io.MultiReader(bytes.NewReader(replaceHeader), d.raw)
	return nil
}

func (d *Decoder) Read(buf []byte) (int, error) {
	return d.audio.Read(buf)
}

func NewTmDecoder(rd io.ReadSeeker) common.Decoder {
	return &Decoder{raw: rd}

}

func init() {
	// QQ Music IOS M4a
	common.RegisterDecoder("tm2", false, NewTmDecoder)
	common.RegisterDecoder("tm6", false, NewTmDecoder)
	// QQ Music IOS Mp3
	common.RegisterDecoder("tm0", false, common.NewRawDecoder)
	common.RegisterDecoder("tm3", false, common.NewRawDecoder)
}
