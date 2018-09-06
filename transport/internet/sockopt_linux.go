package internet

import "syscall"

func applySocketOptions(fd uintptr, config *SocketConfig) error {
	if config.Mark != 0 {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, config.Mark); err != nil {
			return err
		}
	}
	return nil
}
