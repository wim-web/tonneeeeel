package port

import "net"

func AvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")

	if err != nil {
		return -1, err
	}

	return l.Addr().(*net.TCPAddr).Port, nil
}
