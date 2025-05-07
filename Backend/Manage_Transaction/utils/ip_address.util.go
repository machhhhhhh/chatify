package utils

import (
	"net"

	"github.com/gofiber/fiber/v2"
)

func GetIPAdress(ctx *fiber.Ctx) string {
	var master_ip string = ctx.IP()

	ifaces, err := net.Interfaces()
	if err != nil || len(ifaces) == 0 {
		return master_ip
	}

	for i := range ifaces {
		if ifaces[i].Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if ifaces[i].Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		address, err := ifaces[i].Addrs()
		if err != nil {
			return master_ip
		}

		for i := range address {
			var ip net.IP

			switch v := address[i].(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return master_ip + " | " + ip.String()
		}
	}

	return master_ip
}
