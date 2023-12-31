package middleware

import (
	"net"
	"net/http"

	"go-metricscol/internal/server"
)

func (s server.Server) TrustedSubnetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get("X-Real-IP")
		if len(headerValue) != 0 && len(s.Config.TrustedSubnet) != 0 {
			ip := net.ParseIP(headerValue)
			if ip == nil {
				http.Error(w, "Couldn't parse IP from header X-Real-IP", http.StatusBadRequest)
			}

			_, ipNet, err := net.ParseCIDR(s.Config.TrustedSubnet)
			if err != nil {
				http.Error(w, "Couldn't parse CIDR", http.StatusInternalServerError)
			}

			if !ipNet.Contains(ip) {
				http.Error(w, "X-Real-IP is not in trusted subnet", http.StatusForbidden)
			}

			r.Header.Del("X-Real-IP")
		}
		next.ServeHTTP(w, r)
	})
}
