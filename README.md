go-statsd-client
================

The simplest statsd client in Go.

Usage
-----

```
conn, err := net.Dial("udp", "localhost:8125")
if err != nil {
    log.Println("WARNING: failed to created statsd client, using NOOP instead. Error:", err)
}
// if conn is nil, the Statter will do nothing and produce no error.
statter := statsd.Statter{conn}

// Send a counter
statter.Counter(1.0, "succeed", 1)
// Or with a tag
statter.Send(1.0, "succeed", 1, "c", []string{"#foo"})

```
