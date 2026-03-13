package logger

import (
	"net"
	"os"
	"time"
)

// logstashWriter sends log lines (JSON) to Logstash TCP.
// Enabled when LOGSTASH_HOST is set. Fire-and-forget to avoid blocking.
type logstashWriter struct {
	host string
	port string
}

func (w *logstashWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if w.host == "" {
		return n, nil
	}
	// Fire-and-forget: don't block the logger
	go func() {
		// Copy payload (caller may reuse buffer). Logstash TCP json codec needs newline delimiter.
		payload := make([]byte, len(p)+1)
		copy(payload, p)
		payload[len(p)] = '\n'

		addr := net.JoinHostPort(w.host, w.port)
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write(payload)
	}()
	return n, nil
}

func newLogstashWriter() *logstashWriter {
	host := os.Getenv("LOGSTASH_HOST")
	if host == "" {
		return &logstashWriter{}
	}
	port := os.Getenv("LOGSTASH_PORT")
	if port == "" {
		port = "5001"
	}
	return &logstashWriter{host: host, port: port}
}

func initLogstashOutput() *logstashWriter {
	return newLogstashWriter()
}
