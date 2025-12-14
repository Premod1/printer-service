/*
Package main demonstrates usage of PrintRawESCPOSWindows function.

Example ESC/POS Commands Usage:

 1. Basic Text Printing with Paper Cut:
    data := []byte{
    0x1B, 0x40,                     // ESC @ - Initialize printer
    'H', 'e', 'l', 'l', 'o', '!',   // Text content
    0x0A,                           // Line feed
    0x1D, 0x56, 0x42, 0x00,         // GS V B - Paper cut
    }
    err := printer.PrintRawESCPOSWindows("YourPrinterName", data)

 2. Receipt with Formatting:
    data := []byte{
    0x1B, 0x40,          // Initialize
    0x1B, 0x61, 0x01,    // Center align
    0x1B, 0x21, 0x18,    // Double size
    'S', 'T', 'O', 'R', 'E', 0x0A, 0x0A,
    0x1B, 0x21, 0x00,    // Normal size
    0x1B, 0x61, 0x00,    // Left align
    'I', 't', 'e', 'm', ' ', '1', ' ', '$', '1', '0', 0x0A,
    'T', 'o', 't', 'a', 'l', ' ', '$', '1', '0', 0x0A, 0x0A,
    0x1D, 0x56, 0x42, 0x00,  // Paper cut
    }

 3. Barcode Printing (Code128):
    data := []byte{
    0x1B, 0x40,                              // Initialize
    0x1B, 0x61, 0x01,                        // Center align
    0x1D, 0x6B, 0x49, 0x0C,                  // Code128, 12 chars
    '1','2','3','4','5','6','7','8','9','0','1','2',
    0x0A, 0x0A,
    0x1D, 0x56, 0x42, 0x00,                  // Paper cut
    }

 4. Text Formatting:
    data := []byte{
    0x1B, 0x40,          // Initialize
    0x1B, 0x45, 0x01,    // Bold on
    'B', 'O', 'L', 'D', 0x0A,
    0x1B, 0x45, 0x00,    // Bold off
    0x1B, 0x2D, 0x01,    // Underline on
    'U', 'n', 'd', 'e', 'r', 'l', 'i', 'n', 'e', 0x0A,
    0x1B, 0x2D, 0x00,    // Underline off
    0x1D, 0x42, 0x01,    // Reverse on
    'I', 'n', 'v', 'e', 'r', 't', 'e', 'd',
    0x1D, 0x42, 0x00,    // Reverse off
    0x0A, 0x0A,
    0x1D, 0x56, 0x42, 0x00,  // Paper cut
    }

Common ESC/POS Commands:
- 0x1B, 0x40: Initialize printer
- 0x0A: Line feed
- 0x0D: Carriage return
- 0x1B, 0x61, 0x00/0x01/0x02: Left/Center/Right align
- 0x1B, 0x21, 0x00: Normal font
- 0x1B, 0x21, 0x08: Double width
- 0x1B, 0x21, 0x10: Double height
- 0x1B, 0x21, 0x18: Double width and height
- 0x1B, 0x45, 0x01/0x00: Bold on/off
- 0x1B, 0x2D, 0x01/0x00: Underline on/off
- 0x1D, 0x42, 0x01/0x00: Reverse print on/off
- 0x1D, 0x56, 0x42, 0x00: Paper cut
- 0x1D, 0x6B, 0x49: Code128 barcode

Build and run instructions:
1. For Windows: go build -tags windows .
2. For other platforms: go build .

The Windows implementation uses Win32 Print Spooler API with proper error handling
and validates that all bytes are written successfully.
*/
package main
