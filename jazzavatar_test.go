package jazzavatar

import (
	"log"
	"testing"
)

func Test_InitJazzavatar(t *testing.T) {
	ja, err := new(Jazzavatar).Init("abendas", "60", "30")

	if err != nil {
		log.Fatal(err)
	}
	t.Log("Jazzavatar:", ja)
}
