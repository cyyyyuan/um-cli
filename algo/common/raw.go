package common

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RawDecoder struct {
	rd io.ReadSeeker

	audioExt string
}

func NewRawDecoder(rd io.ReadSeeker) Decoder {
	return &RawDecoder{rd: rd}
}

func (d *RawDecoder) Validate() error {
	header := make([]byte, 16)
	if _, err := io.ReadFull(d.rd, header); err != nil {
		return fmt.Errorf("read file header failed: %v", err)
	}
	if _, err := d.rd.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek file failed: %v", err)
	}

	for ext, sniffer := range snifferRegistry {
		if sniffer(header) {
			d.audioExt = strings.ToLower(ext)
			return nil
		}
	}
	return errors.New("audio doesn't recognized")
}

func (d *RawDecoder) Read(p []byte) (n int, err error) {
	return d.rd.Read(p)
}

func init() {
	RegisterDecoder("mp3", true, NewRawDecoder)
	RegisterDecoder("flac", true, NewRawDecoder)
	RegisterDecoder("ogg", true, NewRawDecoder)
	RegisterDecoder("m4a", true, NewRawDecoder)
	RegisterDecoder("wav", true, NewRawDecoder)
	RegisterDecoder("wma", true, NewRawDecoder)
	RegisterDecoder("aac", true, NewRawDecoder)
}
