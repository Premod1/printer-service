//go:build windows
// +build windows

package printer

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Win32 Print Spooler API declarations for Windows
var (
	winspool         = syscall.NewLazyDLL("winspool.drv")
	openPrinterW     = winspool.NewProc("OpenPrinterW")
	startDocPrinterW = winspool.NewProc("StartDocPrinterW")
	startPagePrinter = winspool.NewProc("StartPagePrinter")
	writePrinter     = winspool.NewProc("WritePrinter")
	endPagePrinter   = winspool.NewProc("EndPagePrinter")
	endDocPrinter    = winspool.NewProc("EndDocPrinter")
	closePrinter     = winspool.NewProc("ClosePrinter")
)

// DOC_INFO_1W structure for Win32 API
type DOC_INFO_1W struct {
	PDocName    *uint16
	POutputFile *uint16
	PDatatype   *uint16
}

// stringToUTF16Ptr converts a Go string to a UTF-16 pointer
func stringToUTF16Ptr(s string) (*uint16, error) {
	utf16Slice, err := syscall.UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	return &utf16Slice[0], nil
}

// PrintRawESCPOSWindows prints raw ESC/POS bytes using Win32 Print Spooler API
func PrintRawESCPOSWindows(printerName string, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to print")
	}

	// Convert printer name to UTF-16
	printerNameUTF16, err := stringToUTF16Ptr(printerName)
	if err != nil {
		return fmt.Errorf("failed to convert printer name to UTF-16: %v", err)
	}

	// Step 1: Open printer with OpenPrinterW
	var hPrinter uintptr
	ret, _, err := openPrinterW.Call(
		uintptr(unsafe.Pointer(printerNameUTF16)),
		uintptr(unsafe.Pointer(&hPrinter)),
		uintptr(0), // pDefault (NULL)
	)
	if ret == 0 {
		return fmt.Errorf("OpenPrinterW failed for printer '%s': %v", printerName, err)
	}
	defer func() {
		// Always close printer handle
		closePrinter.Call(hPrinter)
	}()

	// Step 2: Prepare document info for RAW printing
	docName, _ := stringToUTF16Ptr("ESC/POS Raw Print Job")
	dataType, _ := stringToUTF16Ptr("RAW")

	docInfo := DOC_INFO_1W{
		PDocName:    docName,
		POutputFile: nil,
		PDatatype:   dataType,
	}

	// Step 3: Start document with StartDocPrinterW
	ret, _, err = startDocPrinterW.Call(
		hPrinter,
		uintptr(1), // level
		uintptr(unsafe.Pointer(&docInfo)),
	)
	if ret == 0 {
		return fmt.Errorf("StartDocPrinterW failed: %v", err)
	}
	docID := ret
	defer func() {
		// Always end document
		endDocPrinter.Call(hPrinter)
	}()

	// Step 4: Start page with StartPagePrinter
	ret, _, err = startPagePrinter.Call(hPrinter)
	if ret == 0 {
		return fmt.Errorf("StartPagePrinter failed: %v", err)
	}
	defer func() {
		// Always end page
		endPagePrinter.Call(hPrinter)
	}()

	// Step 5: Write raw data with WritePrinter
	var bytesWritten uint32
	ret, _, err = writePrinter.Call(
		hPrinter,
		uintptr(unsafe.Pointer(&data[0])),      // pointer to data
		uintptr(len(data)),                     // data length
		uintptr(unsafe.Pointer(&bytesWritten)), // bytes written
	)
	if ret == 0 {
		return fmt.Errorf("WritePrinter failed: %v", err)
	}

	// Verify all bytes were written
	if int(bytesWritten) != len(data) {
		return fmt.Errorf("incomplete write: wrote %d bytes out of %d", bytesWritten, len(data))
	}

	fmt.Printf("Successfully printed %d bytes to Windows printer '%s' (Job ID: %d)\n",
		bytesWritten, printerName, docID)

	return nil
}
