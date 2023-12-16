package pushback

import (
	"encoding/json"
	"io"
	"log"
)

type Decoder struct {
	decoder *json.Decoder
	buf     []json.Token
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		decoder: json.NewDecoder(r),
	}
}

func (d *Decoder) Token() (json.Token, error) {
	if len(d.buf) > 0 {
		t := d.buf[0]
		d.buf = d.buf[1:]
		return t, nil
	}
	return d.decoder.Token()
}

func (d *Decoder) Pushback(t json.Token) {
	d.buf = append(d.buf, t)
	log.Printf("pushback: %v (%d)", t, len(d.buf))
}
