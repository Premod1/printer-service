//go:build !windows
// +build !windows

package printer

import "fmt"

// PrintRawESCPOSWindows is a stub for non-Windows platforms
func PrintRawESCPOSWindows(printerName string, data []byte) error {
	return fmt.Errorf("Win32 Print Spooler API only available on Windows")
}
