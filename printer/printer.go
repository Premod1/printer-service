package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Printer struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Default bool   `json:"default"`
}

type PrintJob struct {
	PrinterName string `json:"printerName"`
	Content     string `json:"content"`
	JobID       string `json:"jobId"`
}

// DetectPrinters returns list of available printers
func DetectPrinters() ([]Printer, error) {
	var printers []Printer

	switch runtime.GOOS {
	case "windows":
		return detectWindowsPrinters()
	case "darwin":
		return detectMacPrinters()
	case "linux":
		return detectLinuxPrinters()
	default:
		return printers, fmt.Errorf("unsupported platform")
	}
}

func detectWindowsPrinters() ([]Printer, error) {
	var printers []Printer

	// Using PowerShell to get printer list
	cmd := exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_Printer | Select-Object Name, Default, Status | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		return printers, err
	}

	// Parse PowerShell JSON output
	var psPrinters []struct {
		Name    string `json:"Name"`
		Default bool   `json:"Default"`
		Status  string `json:"Status"`
	}

	if err := json.Unmarshal(output, &psPrinters); err != nil {
		// If array parse fails, try single object
		var singlePrinter struct {
			Name    string `json:"Name"`
			Default bool   `json:"Default"`
			Status  string `json:"Status"`
		}
		if err := json.Unmarshal(output, &singlePrinter); err == nil {
			printers = append(printers, Printer{
				Name:    singlePrinter.Name,
				Status:  statusToString(singlePrinter.Status),
				Default: singlePrinter.Default,
			})
		}
	} else {
		for _, p := range psPrinters {
			printers = append(printers, Printer{
				Name:    p.Name,
				Status:  statusToString(p.Status),
				Default: p.Default,
			})
		}
	}

	return printers, nil
}

func detectMacPrinters() ([]Printer, error) {
	var printers []Printer

	cmd := exec.Command("lpstat", "-a")
	output, err := cmd.Output()
	if err != nil {
		return printers, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			printerName := parts[0]
			printers = append(printers, Printer{
				Name:   printerName,
				Status: "Ready",
			})
		}
	}

	return printers, nil
}

func detectLinuxPrinters() ([]Printer, error) {
	var printers []Printer

	cmd := exec.Command("lpstat", "-a")
	output, err := cmd.Output()
	if err != nil {
		return printers, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			printerName := parts[0]
			printers = append(printers, Printer{
				Name:   printerName,
				Status: "Ready",
			})
		}
	}

	return printers, nil
}

func statusToString(status string) string {
	switch status {
	case "3", "OK":
		return "Ready"
	case "4":
		return "Printing"
	case "5":
		return "Warmup"
	default:
		return "Unknown"
	}
}

// PrintText sends text to specified printer
func PrintText(printerName string, content string) error {
	switch runtime.GOOS {
	case "windows":
		return printWindows(printerName, content)
	case "darwin", "linux":
		return printUnix(printerName, content)
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func printWindows(printerName string, content string) error {
	// Create temporary file
	tempFile := "print_temp.txt"
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile)

	// Use Notepad to print (simple approach)
	cmd := exec.Command("notepad", "/p", tempFile)
	return cmd.Run()
}

func printUnix(printerName string, content string) error {
	cmd := exec.Command("lp", "-d", printerName)
	cmd.Stdin = strings.NewReader(content)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Print command failed: %v\nOutput: %s\n", err, string(output))
		return fmt.Errorf("print failed: %v, output: %s", err, string(output))
	}

	fmt.Printf("Print command successful. Output: %s\n", string(output))
	return nil
}

// PrintEscPos sends ESC/POS commands to specified printer
func PrintEscPos(printerName string, escPosData string) error {
	switch runtime.GOOS {
	case "windows":
		return printEscPosWindows(printerName, escPosData)
	case "darwin", "linux":
		return printEscPosUnix(printerName, escPosData)
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func printEscPosWindows(printerName string, escPosData string) error {
	// Create temporary file with ESC/POS data
	tempFile := "print_escpos_temp.bin"
	err := os.WriteFile(tempFile, []byte(escPosData), 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile)

	// Use copy command to send raw data to printer
	cmd := exec.Command("cmd", "/c", "copy", "/b", tempFile, printerName)
	return cmd.Run()
}

func printEscPosUnix(printerName string, escPosData string) error {
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw")
	cmd.Stdin = strings.NewReader(escPosData)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ESC/POS print command failed: %v\nOutput: %s\n", err, string(output))
		return fmt.Errorf("escpos print failed: %v, output: %s", err, string(output))
	}

	fmt.Printf("ESC/POS print command successful. Output: %s\n", string(output))
	return nil
}
