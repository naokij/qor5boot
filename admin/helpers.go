package admin

import (
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type DataTableHeader struct {
	Text     string `json:"text"`
	Value    string `json:"value"`
	Width    string `json:"width"`
	Sortable bool   `json:"sortable"`
}

func getStringHash(v string, len int) string {
	h := sha256.New()
	h.Write([]byte(v))
	return fmt.Sprintf("%x", h.Sum(nil))[:len]
}

func ip(r *http.Request) string {
	if p := proxy(r); len(p) > 0 {
		return strings.TrimSpace(p[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func proxy(r *http.Request) []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}

	return nil
}
