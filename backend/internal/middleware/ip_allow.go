// Package middleware provides implementation for middleware
//
// File: ip_allow.go
// Description: implementation for middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package middleware нь HTTP middleware-уудыг агуулна.

Энэ файл нь IP whitelist middleware-ийг тодорхойлно.
Зөвхөн зөвшөөрөгдсөн IP хаягуудаас хандахыг зөвшөөрнө.

Use cases:
  - /metrics endpoint (зөвхөн internal network)
  - Admin panel (зөвхөн office IP)
  - Health check (зөвхөн load balancer)

Default allowed CIDRs:
  - 127.0.0.1/32 (localhost IPv4)
  - ::1/128 (localhost IPv6)
  - 10.0.0.0/8 (private class A)
  - 172.16.0.0/12 (private class B)
  - 192.168.0.0/16 (private class C)
  - fc00::/7 (IPv6 unique local)
*/
package middleware

import (
	"net"     // IP parsing
	"strings" // String manipulation

	"github.com/gofiber/fiber/v2" // Web framework
)

// ============================================================
// DEFAULT ALLOWED CIDRS
// ============================================================

// defaultAllowedCIDRs нь default-аар зөвшөөрөгдөх IP range-ууд.
// Private network болон localhost.
var defaultAllowedCIDRs = []string{
	// Localhost
	"127.0.0.1/32", // IPv4 localhost
	"::1/128",      // IPv6 localhost

	// Private networks (RFC 1918)
	"10.0.0.0/8",     // Class A private
	"172.16.0.0/12",  // Class B private
	"192.168.0.0/16", // Class C private

	// IPv6 private
	"fc00::/7", // Unique local addresses
}

// ============================================================
// IP ALLOW MIDDLEWARE
// ============================================================

// IPAllow нь зөвхөн зөвшөөрөгдсөн IP-уудаас хандахыг зөвшөөрөх middleware буцаана.
//
// Parameters:
//   - allowed: Зөвшөөрөгдсөн CIDR жагсаалт (хоосон бол default ашиглана)
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Response (хориглогдсон бол):
//   - 403 Forbidden
//
// CIDR format:
//   - IPv4: "192.168.1.0/24"
//   - IPv6: "2001:db8::/32"
//   - Single IP: "192.168.1.1/32"
//
// Жишээ:
//
//	// Default (localhost + private networks)
//	app.Get("/metrics", middleware.IPAllow(nil), handler.Metrics)
//
//	// Custom whitelist
//	app.Get("/admin", middleware.IPAllow([]string{
//	    "203.0.113.0/24",  // Office network
//	    "198.51.100.5/32", // VPN server
//	}), handler.Admin)
//
//	// Config-оос авах
//	app.Get("/metrics", middleware.IPAllow(cfg.Server.MetricsAllowCIDRs), handler.Metrics)
func IPAllow(allowed []string) fiber.Handler {
	// ============================================================
	// PARSE CIDRS
	// ============================================================
	var cidrs []*net.IPNet

	// Хоосон бол default ашиглах
	use := allowed
	if len(use) == 0 {
		use = defaultAllowedCIDRs
	}

	// CIDR-уудыг parse хийх
	for _, c := range use {
		_, n, err := net.ParseCIDR(strings.TrimSpace(c))
		if err == nil {
			cidrs = append(cidrs, n)
		}
		// Parse error бол алгасна (log хийж болно)
	}

	// ============================================================
	// MIDDLEWARE FUNCTION
	// ============================================================
	return func(c *fiber.Ctx) error {
		// Client IP авах
		// Proxy trust тохируулсан бол X-Forwarded-For ашиглана
		ipStr := c.IP()

		// IP parse хийх
		ip := net.ParseIP(ipStr)
		if ip == nil {
			// Invalid IP format
			return fiber.ErrForbidden
		}

		// Зөвшөөрөгдсөн CIDR-д байгаа эсэх шалгах
		for _, n := range cidrs {
			if n.Contains(ip) {
				// Зөвшөөрөгдсөн
				return c.Next()
			}
		}

		// Зөвшөөрөгдөөгүй
		return fiber.ErrForbidden
	}
}
