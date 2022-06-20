//go:build windows
// +build windows

package fastime

import "time"

func (f *Fastime) now() time.Time {
	return time.Now().In(f.location)
}
