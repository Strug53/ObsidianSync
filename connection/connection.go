package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	local_addr string
	OS_Type    string = runtime.GOOS
	Hostname   string
	devices    []Device
	log        = slog.New(slog.NewTextHandler(os.Stdout, nil))

	stopFinding  bool      = false
	maintainChan chan bool = make(chan bool)
)

// mutex
type Device struct {
	Name       string
	IP         string
	OS         string
	last_octet int
}

func removeDuplicateStr(mut *sync.Mutex) {
	mut.Lock()
	allKeys := make(map[string]bool)
	list := []Device{}
	for _, item := range devices {
		if _, value := allKeys[item.Name]; !value {
			allKeys[item.Name] = true
			list = append(list, item)
		}
	}
	devices = list
	log.Debug("Removed Duplicates")
	fmt.Println(devices)
	mut.Unlock()
}

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

func getLastOctet() int {
	str := strings.Split(local_addr, ".")
	last_num, err := strconv.Atoi(str[len(str)-1])
	if err != nil {
		log.Error("ERROR IN CONVERTING. getLastOctet()")
	}
	return last_num
}

func maintainConnection() {
	//fmt.Println(<-maintainChan)

	last_octet_local := getLastOctet()

	str := strings.Split(devices[0].IP, ".")
	last_octet_remote, err := strconv.Atoi(str[len(str)-1])
	if err != nil {
		log.Error("ERROR IN CONVERTING. maintanConnection()")
	}

	if last_octet_local > last_octet_remote {
		conn, err := net.Dial("tcp4", devices[0].IP+":4827")
		if err != nil {
			log.Error("DIAL ERROR", slog.String("err_msg", err.Error()))
		}
		connbuf := bufio.NewReader(conn)

		go func() {
			fmt.Println("SENDING")
			for {
				_, err := conn.Write([]byte("PING\n"))
				if err != nil {
					fmt.Println("Error sending broadcast:", err)
				}
				time.Sleep(5 * time.Second)
			}
		}()
		for {
			str, err := connbuf.ReadString('\n')
			if err != nil {
				break
			}

			if len(str) > 0 {
				fmt.Println(str)
			}
		}
	} else {
		lst, err := net.Listen("tcp4", ":4827")
		if err != nil {
			log.Error("ERROR IN net.Listem", slog.String("err_msg", err.Error()))
		}
		defer lst.Close()

		for {
			conn, err := lst.Accept()
			if err != nil {
				log.Error("ERROR IN ACCCEPTING CONNECTION", slog.String("err_msg", err.Error()))
			}
			remoteAddr := conn.RemoteAddr().String()
			fmt.Println("Client connected from " + remoteAddr)

			scanner := bufio.NewScanner(conn)

			for {
				ok := scanner.Scan()
				if !ok {
					break
				}
				fmt.Println(scanner.Text())
				if scanner.Text() != "" {
					conn.Write([]byte("PONG\n"))
				}
			}

			fmt.Println("Client at " + remoteAddr + " disconnected.")

		}

	}

}

func GetDevices() []Device {
	return devices
}

func SendBroadcastUDP(conn *net.UDPConn, remote *net.UDPAddr) {
	for {
		if stopFinding {
			log.Debug("STOP FINDING. MAINTAINING CONNECTION")
			break
		}
		fmt.Println("SENDINING : ")
		//name := syscall.GetComputerName()
		last_octet := getLastOctet()

		_, err := conn.WriteToUDP([]byte(Hostname+"|"+OS_Type+"|"+strconv.Itoa(last_octet)), remote)
		if err != nil {
			fmt.Println("Error sending broadcast:", err)
		}
		time.Sleep(5 * time.Second) // Отправляем данные каждые 5 секунд
	}
}

func ReceiveBroadcastUDP(conn *net.UDPConn, remote *net.UDPAddr, mut *sync.Mutex) {

	for {
		if stopFinding {
			log.Debug("STOP FINDING. MAINTAINING CONNECTION")
			break
		}
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)

		body_resp := string(buf[:n])

		body_resp_arr := strings.Split(body_resp, "|")

		if err != nil {
			fmt.Println("Error receiving broadcast:", err)
			continue
		}
		if addr.IP.String() == local_addr {
			fmt.Println("Message from local")
			continue
		} else {
			//fmt.Printf("%s RECIEVED: %s | ", addr, buf[:n])
			last_octet, err := strconv.Atoi(body_resp_arr[2])
			if err != nil {
				log.Error("ERROR IN CONVERTING. RecieveBroadcastUDP()")
			}

			devices = append(devices, Device{body_resp_arr[0], addr.IP.String(), body_resp_arr[1], last_octet})
			//fmt.Println(devices)
			go removeDuplicateStr(mut)

		}
	}
}

func main() {
	Host_temp, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	Hostname = Host_temp
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

	var mut sync.Mutex
	go SendBroadcastUDP(conn, remote)
	go ReceiveBroadcastUDP(conn, remote, &mut)

	time.Sleep(10 * time.Second)
	stopFinding = true

	//maintainChan <- true

	maintainConnection()
}
