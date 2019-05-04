package ds18s20

import (
	"encoding/hex"
	"fmt"
	"github.com/schmidtw/dae/onewire"
	"time"
)

const (
	CMD_CONVERT_T         = 0x44
	CMD_READ_SCRATCHPAD   = 0xbe
	CMD_WRITE_SCRATCHPAD  = 0x4e
	CMD_COPY_SCRATCHPAD   = 0x48
	CMD_READ_POWER_SUPPLY = 0xb4
)

type Ds18s20 struct {
	address *onewire.Address
	net     onewire.Adapter
}

func New(adapter onewire.Adapter, addr *onewire.Address) (*Ds18s20, error) {
	d := &Ds18s20{
		net:     adapter,
		address: addr,
	}

	return d, nil
}

func (d *Ds18s20) ConvertAll() {
	d.net.Reset()
	d.net.TxRx([]byte{0xcc, 0x44}, nil)
	time.Sleep(time.Millisecond * 750)
	//d.net.Reset()

	return
}

func (d *Ds18s20) readScratchPad() ([]byte, error) {
	addr := append([]byte{0x55}, d.address.Bytes()...)
	cmd := []byte{CMD_READ_SCRATCHPAD,
		0xff, 0xff, 0xff,
		0xff, 0xff, 0xff,
		0xff, 0xff, 0xff}
	tx := append(addr, cmd...)
	rx := make([]byte, len(tx))

	d.net.Reset()
	fmt.Printf("tx: %s\n", hex.Dump(tx))
	if err := d.net.TxRx(tx, rx); nil != err {
		return nil, err
	}
	fmt.Printf("rx: %s\n", hex.Dump(rx))

	data := rx[len(rx)-9:]
	if data[8] != onewire.Crc8(data[:8]) {
		return nil, fmt.Errorf("CRC didn't match")
	}

	return data, nil
}

// Returns the last measured temperature in degrees C
func (d *Ds18s20) LastTemp() (float64, error) {
	buf, err := d.readScratchPad()
	if nil != err {
		return 0.0, err
	}

	rv := float64(int(int8(buf[1]))<<8+int(buf[0])) / 2

	return rv, nil
}