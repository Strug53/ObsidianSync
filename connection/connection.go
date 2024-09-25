package main

import (
	"fmt"
	"net"
	"time"
)

var (
	local_addr string
)

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP address found")
}

func SendBroadcastUDP(conn *net.UDPConn, remote *net.UDPAddr) {
	for {
		fmt.Println("SENDINING : ")
		_, err := conn.WriteToUDP([]byte("data"), remote)
		if err != nil {
			fmt.Println("Error sending broadcast:", err)
		}
		time.Sleep(5 * time.Second) // Отправляем данные каждые 5 секунд
	}
}

func ReceiveBroadcastUDP(conn *net.UDPConn) {

	for {

		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving broadcast:", err)
			continue
		}
		if addr.IP.String() == local_addr {
			fmt.Println("Message from local")
			continue
		} else {
			fmt.Printf("%s RECIEVED: %s | ", addr, buf[:n])
			fmt.Println(conn.LocalAddr().Network())
		}
	}
}

func main() {
	local, err := net.ResolveUDPAddr("udp4", ":4826")
	if err != nil {
		panic(err)
	}

	remote, err := net.ResolveUDPAddr("udp4", "255.255.255.255:4826")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", local)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	local_addr, err = getLocalIP()
	fmt.Println(local_addr)
	if err != nil {
		panic(err)
	}

	go SendBroadcastUDP(conn, remote)
	ReceiveBroadcastUDP(conn)
}
