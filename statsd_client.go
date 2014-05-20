package statsd

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Statter struct {
	io.Writer
}

func (s Statter) Counter(sampling float64, bucket string, n int) error {
	return s.Send(sampling, bucket, n, "c", nil)
}

func (s Statter) Timing(sampling float64, bucket string, v time.Duration) error {
	return s.Send(sampling, bucket, v.Nanoseconds()/1000000, "ms", nil)
}

func (s Statter) Gauge(sampling float64, bucket string, v interface{}) error {
	return s.Send(sampling, bucket, v, "g", nil)
}

func (s Statter) Histogram(sampling float64, bucket string, v interface{}) error {
	return s.Send(sampling, bucket, v, "h", nil)
}

func (s Statter) Set(sampling float64, bucket string, v interface{}) error {
	return s.Send(sampling, bucket, v, "s", nil)
}

func (s Statter) Send(sampling float64, bucket string, v interface{}, t string, optionals []string) error {
	if s.Writer == nil || !maybe(sampling) {
		return nil
	}

	var val string
	switch v := v.(type) {
	case time.Duration:
		val = fmt.Sprintf("%.3f", v.Seconds())
	case float32, float64:
		val = fmt.Sprintf("%.3f", v)
	default:
		val = fmt.Sprintf("%v", v)
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s:%s|%s", bucket, val, t)
	if sampling < 1.0 {
		fmt.Fprintf(buf, "|@%f", sampling)
	}
	for _, o := range optionals {
		fmt.Fprintf(buf, "|%s", o)
	}
	buf.WriteByte('\n')
	_, err := buf.WriteTo(s)
	return err

}

func maybe(r float64) bool {
	if r >= 1.0 {
		return true
	}
	return rand.Float64() <= r
}
