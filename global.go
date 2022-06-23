package fastime

import (
	"context"
	"sync"
	"time"
)

//nolint:gochecknoglobals
var (
	once     sync.Once
	instance Fastime
)

func init() {
	once.Do(func() {
		instance = New().StartTimerD(context.Background(), time.Millisecond*5)
	})
}

func IsDaemonRunning() bool {
	return instance.IsDaemonRunning()
}

func GetLocation() *time.Location {
	return instance.GetLocation()
}

func GetFormat() string {
	return instance.GetFormat()
}

// SetLocation replaces time location
func SetLocation(location *time.Location) Fastime {
	return instance.SetLocation(location)
}

// SetFormat replaces time format
func SetFormat(format string) Fastime {
	return instance.SetFormat(format)
}

// Now returns current time
func Now() time.Time {
	return instance.Now()
}

// Stop stops stopping time refresh daemon
func Stop() {
	instance.Stop()
}

// UnixNow returns current unix time
func UnixNow() int64 {
	return instance.UnixNow()
}

// UnixUNow returns current unix time
func UnixUNow() uint32 {
	return instance.UnixUNow()
}

// UnixNanoNow returns current unix nano time
func UnixNanoNow() int64 {
	return instance.UnixNanoNow()
}

// UnixUNanoNow returns current unix nano time
func UnixUNanoNow() uint32 {
	return instance.UnixUNanoNow()
}

// FormattedNow returns formatted byte time
func FormattedNow() []byte {
	return instance.FormattedNow()
}

// StartTimerD provides time refresh daemon
func StartTimerD(ctx context.Context, dur time.Duration) Fastime {
	return instance.StartTimerD(ctx, dur)
}
