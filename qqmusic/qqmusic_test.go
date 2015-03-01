package qqmusic

import (
	"testing"
)

func TestGetUrl(t *testing.T) {
	url, err := GetSongUrlByID(101803519)
	t.Log(url, err)
	url, err = GetSongUrlByMID("003uJT681inGC5")
	t.Log(url, err)
}
