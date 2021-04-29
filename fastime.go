package fastime

import (
	"context"
	"math"
	"sync/atomic"
	"time"
	"unsafe"
)

// Fastime is fastime's base struct, it's stores atomic time object
type Fastime struct {
	running       *atomic.Value
	t             *atomic.Value
	ut            int64
	unt           int64
	uut           uint32
	uunt          uint32
	ft            *atomic.Value
	format        *atomic.Value
	cancel        context.CancelFunc
	correctionDur time.Duration
	location      *time.Location
	dur           int64
}

// New returns Fastime. Returned instance is not updated automatically.
// Call Fastime.StartTimerD to start update in the background.
func New() *Fastime {
	return NewWithLocation(time.UTC)
}

// NewWithLocation returns Fastime that will return time in the passed location.
// Nil location will panic. Returned instance is not updated automatically.
// Call Fastime.StartTimerD to start update in the background.
func NewWithLocation(loc *time.Location) *Fastime {
	running := new(atomic.Value)
	running.Store(false)
	f := &Fastime{
		running: running,
		t:       new(atomic.Value),
		ut:      math.MaxInt64,
		unt:     math.MaxInt64,
		uut:     math.MaxUint32,
		uunt:    math.MaxUint32,
		ft: func() *atomic.Value {
			av := new(atomic.Value)
			av.Store(make([]byte, 0, len(time.RFC3339))[:0])
			return av
		}(),
		format: func() *atomic.Value {
			av := new(atomic.Value)
			av.Store(time.RFC3339)
			return av
		}(),
		correctionDur: time.Millisecond * 100,
		location:      loc,
	}
	return f.refresh()
}

func (f *Fastime) update() *Fastime {
	return f.store(f.Now().Add(time.Duration(atomic.LoadInt64(&f.dur))))
}

func (f *Fastime) refresh() *Fastime {
	return f.store(f.now())
}

func (f *Fastime) store(t time.Time) *Fastime {
	f.t.Store(t)
	ut := t.Unix()
	unt := t.UnixNano()
	atomic.StoreInt64(&f.ut, ut)
	atomic.StoreInt64(&f.unt, unt)
	atomic.StoreUint32(&f.uut, *(*uint32)(unsafe.Pointer(&ut)))
	atomic.StoreUint32(&f.uunt, *(*uint32)(unsafe.Pointer(&unt)))
	form := f.format.Load().(string)
	f.ft.Store(t.AppendFormat(make([]byte, 0, len(form)), form))
	return f
}

// SetFormat replaces time format.
func (f *Fastime) SetFormat(format string) *Fastime {
	f.format.Store(format)
	f.refresh()
	return f
}

// Now returns current time.
func (f *Fastime) Now() time.Time {
	return f.t.Load().(time.Time)
}

// Stop stops stopping time refresh daemon.
func (f *Fastime) Stop() {
	if f.running.Load().(bool) {
		f.cancel()
		atomic.StoreInt64(&f.dur, 0)
		return
	}
}

// UnixNow returns current unix time.
func (f *Fastime) UnixNow() int64 {
	return atomic.LoadInt64(&f.ut)
}

// UnixUNow returns current unix time as uint32.
func (f *Fastime) UnixUNow() uint32 {
	return atomic.LoadUint32(&f.uut)
}

// UnixNanoNow returns current unix nano time
func (f *Fastime) UnixNanoNow() int64 {
	return atomic.LoadInt64(&f.unt)
}

// UnixUNanoNow returns current unix nano time as uint32.
func (f *Fastime) UnixUNanoNow() uint32 {
	return atomic.LoadUint32(&f.uunt)
}

// FormattedNow returns formatted byte time
func (f *Fastime) FormattedNow() []byte {
	return f.ft.Load().([]byte)
}

// StartTimerD provides time refresh daemon
func (f *Fastime) StartTimerD(ctx context.Context, dur time.Duration) *Fastime {
	if f.running.Load().(bool) {
		f.Stop()
	}
	f.refresh()

	var ct context.Context
	ct, f.cancel = context.WithCancel(ctx)
	go func() {
		f.running.Store(true)
		f.dur = math.MaxInt64
		atomic.StoreInt64(&f.dur, dur.Nanoseconds())
		ticker := time.NewTicker(time.Duration(atomic.LoadInt64(&f.dur)))
		ctick := time.NewTicker(f.correctionDur)
		for {
			select {
			case <-ct.Done():
				f.running.Store(false)
				ticker.Stop()
				ctick.Stop()
				return
			case <-ticker.C:
				f.update()
			case <-ctick.C:
				select {
				case <-ct.Done():
					f.running.Store(false)
					ticker.Stop()
					ctick.Stop()
					return
				case <-ticker.C:
					f.refresh()
				}
			}
		}
	}()

	return f
}
