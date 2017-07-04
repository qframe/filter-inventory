package qtypes_messages

import (
	"time"
	"crypto/sha1"
	"fmt"
	"github.com/qnib/qframe-types"
)

const (
	version = "0.0.0"
)

type Base struct {
	BaseVersion string
	ID				string
	Time			time.Time
	SourceID		int
	SourcePath		[]string
	SourceSuccess 	bool
	Tags 			map[string]string // Additional KV
	Data			interface{}
}

func NewBase(src string) Base {
	return NewTimedBase(src, time.Now())
}

func NewTimedBase(src string, t time.Time) Base {
	b := Base {
		BaseVersion: version,
		ID: "",
		Time: t,
		SourceID: 0,
		SourcePath: []string{src},
		SourceSuccess: true,
		Tags: map[string]string{},
	}
	return b
}

func NewBaseFromBase(src string, b Base) Base {
	return Base {
		BaseVersion: b.BaseVersion,
		ID: b.ID,
		Time: b.Time,
		SourceID: b.SourceID,
		SourcePath: append(b.SourcePath, src),
		SourceSuccess: b.SourceSuccess,
		Tags: b.Tags,
		Data: b.Data,
	}
}

func NewBaseFromOldBase(src string, b qtypes.Base) Base {
	return Base {
		BaseVersion: b.BaseVersion,
		ID: b.ID,
		Time: b.Time,
		SourceID: b.SourceID,
		SourcePath: append(b.SourcePath, src),
		SourceSuccess: b.SourceSuccess,
		Tags: b.Data,
	}
}


// GenDefaultID uses "<source>-<time.UnixNano()>" and does a sha1 hash.
func (b *Base) GenDefaultID() string {
	s := fmt.Sprintf("%s-%d", b.GetLastSource(), b.Time.UnixNano())
	return Sha1HashString(s)
}

func (b *Base) GetMessageDigest() string {
	return b.ID[:13]
}

func (b *Base) GetTimeRFC() string {
	return b.Time.Format("2006-01-02T15:04:05.999999-07:00")
}

func (b *Base) GetTimeUnix() int64 {
	return b.Time.Unix()
}

func (b *Base) GetTimeUnixNano() int64 {
	return b.Time.UnixNano()
}

func (b *Base) AppendSource(src string) {
	b.SourcePath = append(b.SourcePath, src)
}

func (b *Base) GetLastSource() string {
	return b.SourcePath[len(b.SourcePath)-1]
}

func (b *Base) IsLastSource(src string) bool {
	return b.SourcePath[len(b.SourcePath)-1] == src
}

func (b *Base) InputsMatch(inputs []string) bool {
	for _, inp := range inputs {
		if b.IsLastSource(inp) {
			return true
		}

	}
	return false
}

func Sha1HashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
