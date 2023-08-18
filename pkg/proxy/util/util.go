package util

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	// "math/rand"
	"net"
	"strconv"
	"strings"
)

// ToIP4AddressStr converts the integer IP4 Address to the octet format
func ToIP4AddressStr(ip uint32) string {
	// convert the IP address to a 32-bit binary number
	ipBinary := fmt.Sprintf("%032b", ip)
	// fmt.Printf("This is the value of the ipBinary:%v and this is the value of the ip:%v", ipBinary, ip)

	// divide the binary number into four 8-bit segments
	firstByte, _ := strconv.ParseUint(ipBinary[0:8], 2, 64)
	secondByte, _ := strconv.ParseUint(ipBinary[8:16], 2, 64)
	thirdByte, _ := strconv.ParseUint(ipBinary[16:24], 2, 64)
	fourthByte, _ := strconv.ParseUint(ipBinary[24:32], 2, 64)

	// concatenate the four decimal segments with a dot separator to form the dot-decimal string
	return fmt.Sprintf("%d.%d.%d.%d", firstByte, secondByte, thirdByte, fourthByte)
}

func ToIPv6AddressStr(ip [4]uint32) string {
	// construct a byte slice
	ipBytes := make([]byte, 16) // IPv6 address is 128 bits or 16 bytes long
	for i := 0; i < 4; i++ {
		// for each uint32, extract its four bytes and put them into the byte slice
		ipBytes[i*4] = byte(ip[i] >> 24)
		ipBytes[i*4+1] = byte(ip[i] >> 16)
		ipBytes[i*4+2] = byte(ip[i] >> 8)
		ipBytes[i*4+3] = byte(ip[i])
	}
	// net.IP is a byte slice, so it can be directly used to construct an IPv6 address
	ipv6Addr := net.IP(ipBytes)
	return ipv6Addr.String()
}

func PeekBytes(reader *bufio.Reader) ([]byte, error) {
	var buffer []byte

	// for {
	// Create a temporary buffer to hold the incoming bytes
	buf := make([]byte, 1024)

	// Read bytes from the Reader
	reader.Peek(1)
	buf, err := reader.Peek(reader.Buffered())
	// fmt.Println("read bytes: ", n, ", err: ", err)
	if err != nil && err != io.EOF && err != bufio.ErrBufferFull {
		return nil, err
	}

	// Append the bytes to the buffer
	buffer = append(buffer, buf[:]...)

	// If we've reached the end of the input stream, break out of the loop
	// if err == bufio.ErrBufferFull || len(buf) != 1024 {
	// 	break
	// }
	// }

	return buffer, nil
}

// ReadBytes function is utilized to read the complete message from the reader until the end of the file (EOF).
// It returns the content as a byte array.
func ReadBytes(reader io.Reader) ([]byte, error) {
	var buffer []byte

	for {
		// Create a temporary buffer to hold the incoming bytes
		buf := make([]byte, 1024)
		// rand.Seed(time.Now().UnixNano())

		// Read bytes from the Reader
		n, err := reader.Read(buf)
		// fmt.Println("read bytes: ", n, ", err: ", err)
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Append the bytes to the buffer
		buffer = append(buffer, buf[:n]...)

		// If we've reached the end of the input stream, break out of the loop
		if err == io.EOF || n != 1024 {
			break
		}
	}

	return buffer, nil
}

func ReadBytes1(reader io.Reader) ([]byte, string, error) {
	logStr := ""
	var buffer []byte

	for {
		// Create a temporary buffer to hold the incoming bytes
		buf := make([]byte, 1024)

		// Read bytes from the Reader
		n, err := reader.Read(buf)
		// fmt.Println("read bytes: ", n , ", err: ", err)
		logStr += fmt.Sprintln("read bytes: ", n, ", err: ", err)
		if err != nil && err != io.EOF {
			return nil, logStr, err
		}

		// Append the bytes to the buffer
		buffer = append(buffer, buf[:n]...)

		// If we've reached the end of the input stream, break out of the loop
		// if err == io.EOF || n != 1024 {
		if n != 1024 {
			break
		}
	}

	return buffer, logStr, nil
}

func GetLocalIPv4() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP, nil
			}
		}
	}

	return nil, fmt.Errorf("No valid IP address found")
}

func ConvertToIPV4(ip net.IP) (uint32, bool) {
	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0, false // Return 0 or handle the error accordingly
	}

	return uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3]), true
}

func GetLocalIPv6() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() == nil && ipNet.IP.To16() != nil {
				return ipNet.IP, nil
			}
		}
	}

	return nil, fmt.Errorf("No valid IPv6 address found")
}

func ConvertIPv6ToUint32Array(ip net.IP) ([4]uint32, error) {
	ip = ip.To16()
	if ip == nil {
		return [4]uint32{}, errors.New("invalid IPv6 address")
	}

	return [4]uint32{
		binary.BigEndian.Uint32(ip[0:4]),
		binary.BigEndian.Uint32(ip[4:8]),
		binary.BigEndian.Uint32(ip[8:12]),
		binary.BigEndian.Uint32(ip[12:16]),
	}, nil
}

func IPToDotDecimal(ip net.IP) string {
	ipStr := ip.String()
	if ip.To4() != nil {
		ipStr = ip.To4().String()
	}
	return ipStr
}

// It checks if the cmd is related to docker or not, it also returns if its a docker compose file
func IsDockerRelatedCommand(cmd string) (bool, string) {
	// Check for Docker command patterns
	dockerCommandPatterns := []string{
		"docker-compose ",
		"sudo docker-compose ",
		"docker compose ",
		"sudo docker compose ",
		"docker ",
		"sudo docker ",
	}

	for _, pattern := range dockerCommandPatterns {
		if strings.HasPrefix(strings.ToLower(cmd), pattern) {
			if strings.Contains(pattern, "compose") {
				return true, "docker-compose"
			}
			return true, "docker"
		}
	}

	// Check for Docker Compose file extension
	dockerComposeFileExtensions := []string{".yaml", ".yml"}
	for _, extension := range dockerComposeFileExtensions {
		if strings.HasSuffix(strings.ToLower(cmd), extension) {
			return true, "docker-compose"
		}
	}

	return false, ""
}
