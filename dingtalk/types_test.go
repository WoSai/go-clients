package dingtalk

import (
	"fmt"
	"testing"
)

func TestUnixTimestamp_Time(t *testing.T) {
	u := &UnixTimestamp{ts: 1597573616828}
	fmt.Println(u.Time().String())
}
