package defuse

import (
	"testing"
	"time"
)

// Defuse is the means by which a bomb may be defused.
type Defuse struct {
	c chan struct{}
}

// Defuse defuses the bomb.
func (d *Defuse) Defuse() {
	d.c <- struct{}{}
}

// Bomb explodes if not defused within the given duration. The effect of the
// explosion needs to be provided. Returns a Defuse that is used to defuse the
// bomb before it goes off.
func Bomb(dur time.Duration, explosion func()) Defuse {
	defuse := Defuse{
		c: make(chan struct{}, 1),
	}

	go func(defuse Defuse) {
		fuse := time.NewTimer(dur)
		select {
		case <-fuse.C:
			explosion()
		case <-defuse.c:
		}
		fuse.Stop()
	}(defuse)

	return defuse
}

// DoOrElse sets up the bomb for you. The given function must be completed
// within the given duration or the given explosion will occur.
func DoOrElse(do func(), within time.Duration, or func()) {
	defuse := Bomb(within, or)
	do()
	defuse.Defuse()
}

// FailTest returns an explosion that calls Fatal on the given testing.T,
// terminating the test. Includes the given string in the failure message.
func FailTest(t *testing.T, msg string) func() {
	return func() {
		t.Fatal("BOOM: " + msg)
	}
}

// Panic returns an explosion that panics with the given message.
func Panic(msg string) func() {
	return func() {
		panic(msg)
	}
}
