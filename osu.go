package replay

import (
	"encoding/binary"
	"errors"
	"io"
)

var oneByte = make([]byte, 1)

type osuReader struct {
	r io.Reader
}

func (or *osuReader) Read(p []byte) (n int, err error) {
	return or.r.Read(p)
}

func (or *osuReader) ReadByte() (byte, error) {
	_, e := or.Read(oneByte)
	if e != nil {
		return 0, e
	}
	return oneByte[0], nil
}

func (or *osuReader) ReadTypes(vals ...interface{}) {
	var e error
	for _, val := range vals {
		switch val := val.(type) {
		case *string:
			switch b, _ := or.ReadByte(); b {
			case 0x0b:
				var size uint64
				size, e = binary.ReadUvarint(or)
				s := make([]byte, size)
				_, e = or.Read(s)
				*val = string(s)
			case 0x00:
				*val = ""
				break
			default:
				e = errors.New("osuReader: invalid string")
			}
		default:
			e = binary.Read(or, binary.LittleEndian, val)
		}

		if e != nil {
			panic(e)
		}
	}
}

type osuWriter struct {
	w io.Writer
}

/*
func PutUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}
 */

func (ow osuWriter) uvarint(x int) error {
	var e error
	for x >= 0x80 {
		e = ow.WriteByte(byte(x) | 0x80)
		if e != nil {
			return e
		}
		x >>= 7
	}
	ow.WriteByte(byte(x))
	if e != nil {
		return e
	}
	return nil
}

func (ow osuWriter) WriteTypes(vals ...interface{}) {
	var e error

	for _, val := range vals {
		switch val := val.(type) {
		case string:
			if val == "" {
				e = ow.WriteByte(0x00)
				break
			}
			e = ow.WriteByte(0x0b)
			if e != nil {
				break
			}
			e = ow.uvarint(len(val))
			if e != nil {
				break
			}
			_, e = ow.Write([]byte(val))
		default:
			e = binary.Write(ow, binary.LittleEndian, val)
		}

		if e != nil {
			panic(e)
		}
	}
}

func (ow osuWriter) WriteByte(c byte) error {
	oneByte[0] = c
	_, e := ow.Write(oneByte)
	return e
}

func (ow osuWriter) Write(p []byte) (n int, err error) {
	return ow.w.Write(p)
}


