package timeat

// NTP client code

import (
	"encoding/binary"
	"net"
	"time"
)

type ntpTime struct {
	Seconds  uint32
	Fraction uint32
}

type ntpMsg struct {
	LiVnMode       byte // Leap Indicator (2) + Version (3) + Mode (3)
	Stratum        byte
	Poll           byte
	Precision      byte
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	ReferenceTime  ntpTime
	OriginTime     ntpTime
	ReceiveTime    ntpTime
	TransmitTime   ntpTime
}

func (t ntpTime) UTC() time.Time {
	nsec := uint64(t.Seconds)*1e9 + (uint64(t.Fraction) * 1e9 >> 32)
	return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nsec))
}

// NTPTime returns the current time (UTC) from NTP
// NTP code adapted from https://github.com/beevik/ntp/blob/master/ntp.go
func NTPTime() (time.Time, error) {
	conn, err := net.Dial("udp", "pool.ntp.org:123")
	if err != nil {
		return time.Time{}, nil
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(2 * time.Second))

	msg := ntpMsg{
		LiVnMode: 0x1B,
	}

	err = binary.Write(conn, binary.BigEndian, &msg)
	if err != nil {
		return time.Time{}, err
	}

	err = binary.Read(conn, binary.BigEndian, &msg)
	if err != nil {
		return time.Time{}, err
	}

	return msg.ReceiveTime.UTC(), nil
}
