# 🖨️ Printer Service

A comprehensive cross-platform service for printing from web applications to local system printers. Features both REST API and WebSocket communication with full ESC/POS thermal printing support.

## 🚀 Quick Start

```bash
# Run the service
go run main.go

# Test the interfaces
# REST API: http://localhost:8081/api
# WebSocket: ws://localhost:8081/ws  
# Test UI: http://localhost:8081/
```

## ✨ Features

### Core Functionality
- 🌐 **Dual Communication**: REST API + WebSocket support
- 🖥️ **Cross-platform**: Windows, macOS, Linux compatibility
- 🧾 **Multiple Formats**: Plain text and ESC/POS thermal printing
- 📋 **Job Management**: Queue system with status tracking
- 🔄 **Auto Detection**: System printer discovery

### Enterprise Features
- 📊 **Progress Tracking**: Real-time job monitoring
- 🎯 **80mm Optimization**: Perfect thermal receipt formatting
- 🔒 **Error Handling**: Comprehensive error reporting
- 🚀 **High Performance**: Concurrent job processing
- 📱 **Web Interface**: Built-in test UI

## 📋 Table of Contents

- [🚀 Quick Start](#-quick-start)
- [📦 Installation](#-installation)
- [🔌 REST API](#-rest-api)
- [🌐 WebSocket API](#-websocket-api)
- [🧾 ESC/POS Printing](#-escpos-printing)
- [🔧 Integration Examples](#-integration-examples)
- [🖥️ Platform Support](#%EF%B8%8F-platform-support)
- [🛠️ Troubleshooting](#%EF%B8%8F-troubleshooting)
- [📚 Advanced Usage](#-advanced-usage)

## ⚡ Installation & Setup

### Prerequisites

- **Go 1.21+** for building from source
- **System printer drivers** installed
- **CUPS** (Linux/macOS) or **Windows Print Spooler**

### Option 1: Download Pre-built Executable

```bash
# Windows
curl -L -o printer-service.exe \
  https://github.com/Premod1/printer-service/releases/latest/download/printer-service.exe
./printer-service.exe

# Linux
curl -L -o printer-service \
  https://github.com/Premod1/printer-service/releases/latest/download/printer-service-linux
chmod +x printer-service
./printer-service

# macOS
curl -L -o printer-service \
  https://github.com/Premod1/printer-service/releases/latest/download/printer-service-mac
chmod +x printer-service
./printer-service
```

### Option 2: Build from Source

```bash
git clone https://github.com/Premod1/printer-service.git
cd printer-service
go mod download
go run main.go

# Or build executable
go build -o printer-service
```

### Verify Installation

```bash
# Health check
curl http://localhost:8081/health
# Response: "Printer service is running"

# Check available endpoints
curl http://localhost:8081/api/printers     # REST API
# WebSocket: ws://localhost:8081/ws         # WebSocket
# Test UI: http://localhost:8081/           # Web Interface
```

## 🔌 REST API

The REST API provides HTTP endpoints for printer management and print job submission with comprehensive job tracking.

**Base URL**: `http://localhost:8081/api`

### API Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/printers` | Get available printers |
| `POST` | `/api/print/text` | Queue text print job |
| `POST` | `/api/print/escpos` | Queue ESC/POS print job |
| `GET` | `/api/jobs` | List all print jobs |
| `GET` | `/api/jobs/{jobId}` | Get specific job status |

### 1. Get Available Printers

```bash
curl -X GET http://localhost:8081/api/printers
```

**Response**:
```json
{
  "success": true,
  "printers": [
    {
      "name": "HP LaserJet Pro",
      "status": "Ready",
      "default": true
    },
    {
      "name": "Thermal Receipt Printer",
      "status": "Ready", 
      "default": false
    }
  ]
}
```

### 2. Print Text Content

```bash
curl -X POST http://localhost:8081/api/print/text \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "HP LaserJet Pro",
    "content": "Hello World!\nThis is a test print.\nTimestamp: 2025-12-11 10:30:00"
  }'
```

**Response**:
```json
{
  "success": true,
  "jobId": "job_1765360872587975241",
  "message": "Print job queued successfully"
}
```

### 3. Print ESC/POS Content

```bash
curl -X POST http://localhost:8081/api/print/escpos \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "Thermal Receipt Printer",
    "rawData": "\\x1b@\\x1ba\\x01\\x1b!\\x30RECEIPT\\n\\x1b!\\x00\\x1ba\\x00Item 1\\x09\\x09$10.00\\nItem 2\\x09\\x09$15.00\\n\\x1b\\x64\\x02\\x1b\\x45\\x01TOTAL: $25.00\\x1b\\x45\\x00\\n\\x1b\\x64\\x05\\x1b\\x69"
  }'
```

### 4. Monitor Job Status

```bash
# Get all jobs
curl -X GET http://localhost:8081/api/jobs

# Get specific job
curl -X GET http://localhost:8081/api/jobs/job_1765360872587975241

# Filter jobs by status
curl -X GET "http://localhost:8081/api/jobs?status=completed&limit=5"
```

**Job Response**:
```json
{
  "success": true,
  "job": {
    "jobId": "job_1765360872587975241",
    "printerName": "HP LaserJet Pro",
    "content": "Hello World Test Print",
    "status": "completed",
    "createdAt": "2025-12-11T10:31:12.587977806+05:30",
    "completedAt": "2025-12-11T10:31:12.665607362+05:30",
    "error": "",
    "jobType": "text",
    "progress": 100
  }
}
```

### Job Status States

| Status | Description |
|--------|-------------|
| `pending` | Job queued, waiting to be processed |
| `printing` | Job currently being sent to printer |
| `completed` | Job finished successfully |
| `failed` | Job failed with error |

## 🌐 WebSocket API

For real-time applications, use WebSocket communication at `ws://localhost:8081/ws`.

### Connection Example

```javascript
const ws = new WebSocket('ws://localhost:8081/ws');

ws.onopen = () => {
  console.log('Connected to printer service');
  // Get printers immediately
  ws.send(JSON.stringify({ type: 'get_printers' }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleMessage(message);
};
```

### Message Types

#### Get Printers
```javascript
// Send
ws.send(JSON.stringify({
  type: "get_printers"
}));

// Receive
{
  "type": "printers_list",
  "payload": [
    {
      "name": "HP LaserJet Pro",
      "status": "Ready",
      "default": true
    }
  ]
}
```

#### Print Text
```javascript
// Send  
ws.send(JSON.stringify({
  type: "print",
  payload: {
    printerName: "HP LaserJet Pro",
    content: "Hello World!\nThis is a test print.",
    jobId: "job_123456"
  }
}));

// Receive
{
  "type": "print_success",
  "payload": {
    "jobId": "job_123456"
  }
}
```

#### Print ESC/POS
```javascript
// Send
ws.send(JSON.stringify({
  type: "print_raw_escpos",
  payload: {
    printerName: "Receipt Printer",
    jobId: "escpos_001", 
    rawData: "\x1b@\x1ba\x01RECEIPT\n\x1dV\x41\x00"
  }
}));

// Receive
{
  "type": "raw_escpos_print_success",
  "payload": {
    "jobId": "escpos_001"
  }
}
```

#### Error Handling
```javascript
{
  "type": "error",
  "payload": {
    "message": "Printer not found: Invalid Printer"
  }
}
```

## 🧾 ESC/POS Printing

Professional thermal receipt printing with 80mm paper optimization.

### Quick ESC/POS Example

```bash
# Complete thermal receipt
curl -X POST http://localhost:8081/api/print/escpos \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "POS-80",
    "rawData": "\\x1b@\\x1ba\\x01\\x1b!\\x30TECH STORE\\n\\x1b!\\x00\\x1ba\\x01123 Main Street\\nSilicon Valley, CA\\n\\x1b\\x64\\x02\\x1ba\\x00================================\\nDATE: 2025-12-11 10:30:00\\nCASHIER: Alice Johnson\\nRECEIPT #: R001234\\n================================\\n\\x1b\\x45\\x01ITEM\\x09\\x09QTY\\x09PRICE\\x1b\\x45\\x00\\n--------------------------------\\nLaptop\\x09\\x091\\x09$999.99\\nMouse\\x09\\x092\\x09$50.00\\nKeyboard\\x09\\x091\\x09$75.00\\n--------------------------------\\n\\x1b\\x45\\x01SUBTOTAL:\\x09\\x09$1124.99\\nTAX:\\x09\\x09\\x09$89.99\\nTOTAL:\\x09\\x09\\x09$1214.98\\x1b\\x45\\x00\\n================================\\n\\x1ba\\x01PAYMENT: CREDIT CARD\\nAMOUNT: $1214.98\\n\\x1b\\x64\\x02\\x1ba\\x01Thank you for shopping!\\nPlease visit us again!\\n\\x1b\\x64\\x05\\x1b\\x69"
  }'
```

### ESC/POS Command Reference

| Command | Code | Description | Example |
|---------|------|-------------|---------|
| **Initialize** | `\\x1b@` | Reset printer | `\\x1b@` |
| **Text Size** | `\\x1b!` | Set text attributes | `\\x1b!\\x30` (large) |
| **Bold On/Off** | `\\x1b\\x45` | Bold formatting | `\\x1b\\x45\\x01` (on) |
| **Align** | `\\x1ba` | Text alignment | `\\x1ba\\x01` (center) |
| **Feed Lines** | `\\x1b\\x64` | Line spacing | `\\x1b\\x64\\x02` (2 lines) |
| **Cut Paper** | `\\x1b\\x69` | Full cut | `\\x1b\\x69` |

### Text Attributes

```javascript
const ESC_POS = {
  // Initialize
  INIT: '\\x1b@',
  
  // Text Size
  NORMAL: '\\x1b!\\x00',
  DOUBLE_HEIGHT: '\\x1b!\\x10', 
  DOUBLE_WIDTH: '\\x1b!\\x20',
  LARGE: '\\x1b!\\x30',
  
  // Text Style
  BOLD_ON: '\\x1b\\x45\\x01',
  BOLD_OFF: '\\x1b\\x45\\x00',
  UNDERLINE_ON: '\\x1b\\x2d\\x01',
  UNDERLINE_OFF: '\\x1b\\x2d\\x00',
  
  // Alignment
  ALIGN_LEFT: '\\x1ba\\x00',
  ALIGN_CENTER: '\\x1ba\\x01',
  ALIGN_RIGHT: '\\x1ba\\x02',
  
  // Paper Control
  FEED_2_LINES: '\\x1b\\x64\\x02',
  FEED_5_LINES: '\\x1b\\x64\\x05',
  PARTIAL_CUT: '\\x1b\\x6d',
  FULL_CUT: '\\x1b\\x69'
};
```

### Receipt Generation Helper

```javascript
function generateReceipt(data) {
  let commands = '';
  
  // Header
  commands += '\\x1b@';  // Initialize
  commands += '\\x1ba\\x01';  // Center
  commands += `\\x1b!\\x30${data.storeName}\\n`;  // Large store name
  commands += '\\x1b!\\x00' + data.address + '\\n';
  commands += '\\x1b\\x64\\x02';  // Feed 2 lines
  
  // Items
  commands += '\\x1ba\\x00';  // Left align
  commands += '================================\\n';
  data.items.forEach(item => {
    commands += `${item.name}\\x09${item.qty}\\x09$${item.price}\\n`;
  });
  
  // Total
  commands += '--------------------------------\\n';
  commands += `\\x1b\\x45\\x01TOTAL: $${data.total}\\x1b\\x45\\x00\\n`;
  
  // Footer
  commands += '\\x1ba\\x01Thank you!\\n';
  commands += '\\x1b\\x64\\x05\\x1b\\x69';  // Feed and cut
  
  return commands;
}
```

### 80mm Paper Optimization

- **Line Width**: 48 characters maximum
- **Font**: Font A (12x24 dots) recommended  
- **Margins**: 2-4 characters on each side
- **Separator**: Use `=` or `-` for visual separation
- **Price Alignment**: Right-align using `\\x09` (tab) characters

## 🔧 Integration Examples

### JavaScript/Web Application

```javascript
class PrinterService {
  constructor(baseUrl = 'http://localhost:8081/api') {
    this.baseUrl = baseUrl;
  }

  async getPrinters() {
    const response = await fetch(`${this.baseUrl}/printers`);
    const data = await response.json();
    return data.printers;
  }

  async printText(printerName, content) {
    const response = await fetch(`${this.baseUrl}/print/text`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ printerName, content })
    });
    const data = await response.json();
    return data.jobId;
  }

  async printReceipt(printerName, receiptData) {
    const escPosCommands = this.generateEscPosReceipt(receiptData);
    const response = await fetch(`${this.baseUrl}/print/escpos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ printerName, rawData: escPosCommands })
    });
    const data = await response.json();
    return data.jobId;
  }

  async getJobStatus(jobId) {
    const response = await fetch(`${this.baseUrl}/jobs/${jobId}`);
    const data = await response.json();
    return data.job;
  }

  generateEscPosReceipt(data) {
    let commands = '';
    
    // Initialize and header
    commands += '\\x1b@\\x1ba\\x01';
    commands += `\\x1b!\\x30${data.businessName}\\n`;
    commands += `\\x1b!\\x00${data.address}\\n`;
    commands += '\\x1b\\x64\\x02';
    
    // Invoice details
    commands += '\\x1ba\\x00================================\\n';
    commands += `Invoice: ${data.invoiceCode}\\n`;
    commands += `Date: ${data.date}\\n`;
    commands += '================================\\n';
    
    // Items
    commands += '\\x1b\\x45\\x01ITEM\\x09\\x09QTY\\x09PRICE\\x1b\\x45\\x00\\n';
    commands += '--------------------------------\\n';
    data.items.forEach(item => {
      commands += `${item.name}\\x09${item.quantity}\\x09$${item.price}\\n`;
    });
    
    // Totals
    commands += '--------------------------------\\n';
    commands += `\\x1b\\x45\\x01TOTAL: $${data.total}\\x1b\\x45\\x00\\n`;
    commands += '================================\\n';
    
    // Footer
    commands += '\\x1ba\\x01Thank you for your business!\\n';
    commands += '\\x1b\\x64\\x05\\x1b\\x69';
    
    return commands;
  }
}

// Usage example
const printer = new PrinterService();

async function printSampleReceipt() {
  const printers = await printer.getPrinters();
  console.log('Available printers:', printers);
  
  const receiptData = {
    businessName: 'TECH STORE',
    address: '123 Main Street\\nSilicon Valley, CA',
    invoiceCode: 'INV-001',
    date: new Date().toLocaleString(),
    items: [
      { name: 'Laptop', quantity: 1, price: '999.99' },
      { name: 'Mouse', quantity: 2, price: '25.00' }
    ],
    total: '1049.99'
  };
  
  if (printers.length > 0) {
    const jobId = await printer.printReceipt(printers[0].name, receiptData);
    console.log('Print job queued:', jobId);
    
    // Monitor job status
    setTimeout(async () => {
      const status = await printer.getJobStatus(jobId);
      console.log('Job status:', status.status);
    }, 2000);
  }
}
```

### Laravel Integration

```php
<?php
// Laravel Controller example
namespace App\\Http\\Controllers;

use Illuminate\\Http\\Request;
use Illuminate\\Support\\Facades\\Http;

class PrinterController extends Controller
{
    private $printerServiceUrl = 'http://localhost:8081/api';
    
    public function getPrinters()
    {
        $response = Http::get("{$this->printerServiceUrl}/printers");
        return response()->json($response->json());
    }
    
    public function printReceipt(Request $request)
    {
        $validated = $request->validate([
            'printerName' => 'required|string',
            'businessName' => 'required|string',
            'items' => 'required|array',
            'items.*.name' => 'required|string',
            'items.*.quantity' => 'required|numeric|min:1',
            'items.*.price' => 'required|numeric|min:0',
        ]);
        
        // Generate ESC/POS commands
        $escPosData = $this->generateEscPosReceipt($validated);
        
        // Send to printer service
        $response = Http::post("{$this->printerServiceUrl}/print/escpos", [
            'printerName' => $validated['printerName'],
            'rawData' => $escPosData
        ]);
        
        return response()->json($response->json());
    }
    
    private function generateEscPosReceipt($data)
    {
        $commands = '';
        
        // Header
        $commands .= "\\x1b@\\x1ba\\x01";  // Init + center
        $commands .= "\\x1b!\\x30{$data['businessName']}\\n";  // Large business name
        $commands .= "\\x1b!\\x00{$data['address']}\\n";  // Address
        $commands .= "\\x1b\\x64\\x02";  // Feed 2 lines
        
        // Items
        $commands .= "\\x1ba\\x00================================\\n";
        $commands .= "\\x1b\\x45\\x01ITEM\\x09\\x09QTY\\x09PRICE\\x1b\\x45\\x00\\n";
        $commands .= "--------------------------------\\n";
        
        foreach ($data['items'] as $item) {
            $commands .= "{$item['name']}\\x09{$item['quantity']}\\x09\${$item['price']}\\n";
        }
        
        // Total
        $total = array_sum(array_map(fn($item) => $item['quantity'] * $item['price'], $data['items']));
        $commands .= "--------------------------------\\n";
        $commands .= "\\x1b\\x45\\x01TOTAL: \${$total}\\x1b\\x45\\x00\\n";
        
        // Footer
        $commands .= "\\x1ba\\x01Thank you!\\n";
        $commands .= "\\x1b\\x64\\x05\\x1b\\x69";  // Feed and cut
        
        return $commands;
    }
}
```

### Python Integration

```python
import requests
import json
import time

class PrinterService:
    def __init__(self, base_url="http://localhost:8081/api"):
        self.base_url = base_url
    
    def get_printers(self):
        response = requests.get(f"{self.base_url}/printers")
        return response.json()["printers"]
    
    def print_text(self, printer_name, content):
        data = {"printerName": printer_name, "content": content}
        response = requests.post(f"{self.base_url}/print/text", json=data)
        return response.json()["jobId"]
    
    def print_escpos(self, printer_name, raw_data):
        data = {"printerName": printer_name, "rawData": raw_data}
        response = requests.post(f"{self.base_url}/print/escpos", json=data)
        return response.json()["jobId"]
    
    def get_job_status(self, job_id):
        response = requests.get(f"{self.base_url}/jobs/{job_id}")
        return response.json()["job"]
    
    def wait_for_job(self, job_id, timeout=30):
        start_time = time.time()
        while time.time() - start_time < timeout:
            job = self.get_job_status(job_id)
            if job["status"] in ["completed", "failed"]:
                return job
            time.sleep(1)
        return None

# Usage example
printer = PrinterService()

# List available printers
printers = printer.get_printers()
print("Available printers:", [p["name"] for p in printers])

# Print test message
if printers:
    job_id = printer.print_text(printers[0]["name"], "Test print from Python")
    print(f"Job queued: {job_id}")
    
    # Wait for completion
    final_job = printer.wait_for_job(job_id)
    if final_job:
        print(f"Job {final_job['status']}: {final_job.get('error', 'Success')}")
```

### Vue.js Composable

```javascript
// composables/usePrinterService.js
import { ref } from 'vue'

export function usePrinterService() {
  const printers = ref([])
  const loading = ref(false)
  const error = ref('')

  const baseUrl = 'http://localhost:8081/api'

  const getPrinters = async () => {
    loading.value = true
    try {
      const response = await fetch(`${baseUrl}/printers`)
      const data = await response.json()
      printers.value = data.printers
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const printText = async (printerName, content) => {
    const response = await fetch(`${baseUrl}/print/text`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ printerName, content })
    })
    const data = await response.json()
    return data.jobId
  }

  return {
    printers,
    loading,
    error,
    getPrinters,
    printText
  }
}
```

## 🖥️ Platform Support

### Windows
- **Detection**: PowerShell + WMI (`Get-WmiObject -Class Win32_Printer`)
- **Printing**: 
  - Text: Notepad command-line printing
  - ESC/POS: Win32 Print Spooler API for raw printing
- **Requirements**: Windows 7+, PowerShell enabled, Print Spooler service
- **Printer Types**: All Windows-compatible printers

```powershell
# Verify Print Spooler service
net start spooler

# Check available printers
Get-WmiObject -Class Win32_Printer | Select-Object Name, Default, Status
```

### Linux
- **Detection**: CUPS (`lpstat -a`)
- **Printing**: 
  - Text: CUPS (`lp -d <printer>`)
  - ESC/POS: CUPS with raw option (`lp -d <printer> -o raw`)
- **Requirements**: CUPS service running, printer drivers installed
- **Printer Types**: CUPS-supported printers, USB thermal printers

```bash
# Check CUPS service
systemctl status cups
sudo systemctl start cups

# List available printers
lpstat -a

# Test printer
echo "Test" | lp -d PrinterName
```

### macOS
- **Detection**: Built-in CUPS (`lpstat -a`)
- **Printing**: 
  - Text: CUPS (`lp -d <printer>`)
  - ESC/POS: CUPS with raw option
- **Requirements**: macOS 10.9+, printer drivers installed
- **Printer Types**: AirPrint, USB, network printers

```bash
# List printers
lpstat -a

# System printer settings
open "System Preferences" -b com.apple.preference.printfax
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Service Connection Failed

```bash
# Check if service is running
curl http://localhost:8081/health
# Expected response: "Printer service is running"

# Verify port availability
netstat -an | grep 8081
lsof -i :8081  # Mac/Linux

# Check firewall settings
sudo ufw allow 8081  # Linux
netsh firewall set portopening TCP 8081 "Printer Service"  # Windows
```

#### 2. No Printers Found

**Windows:**
```powershell
# Check Print Spooler service
net start spooler
Get-Service -Name Spooler

# List printers via PowerShell
Get-WmiObject -Class Win32_Printer | Select-Object Name, Status
```

**Linux/macOS:**
```bash
# Check CUPS service
systemctl status cups
sudo systemctl start cups

# List available printers  
lpstat -a
lpstat -p  # Detailed printer status
```

#### 3. Print Jobs Fail

**General Debugging:**
- Verify printer is online and has paper/ink
- Test with system print dialog first
- Check printer drivers are properly installed
- Try different printer if available

**Check Logs:**
```bash
# Enable debug mode
DEBUG=true go run main.go

# Monitor logs
tail -f /var/log/cups/error_log  # Linux/macOS
# Windows Event Viewer > Windows Logs > System
```

#### 4. ESC/POS Commands Not Working

**Thermal Printer Issues:**
- Ensure direct USB connection (not network)
- Verify printer supports ESC/POS command set
- Check printer manual for supported commands
- Test with simple ESC/POS commands first:

```bash
# Test basic ESC/POS
curl -X POST http://localhost:8081/api/print/escpos \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "Your-Thermal-Printer",
    "rawData": "\\x1b@Test\\n\\x1b\\x69"
  }'
```

#### 5. Permission Issues

**Windows:**
- Run service as Administrator if needed
- Check UAC settings
- Verify user has printer access

**Linux/macOS:**
```bash
# Add user to printer groups
sudo usermod -a -G lpadmin $USER
sudo usermod -a -G lp $USER

# Check printer permissions
ls -la /dev/usb/lp*
```

### Performance Optimization

#### Job Queue Tuning
- Default queue size: 100 jobs
- Job processing delay: 500ms between jobs
- Modify in `main.go` if needed:

```go
queue: make(chan *PrintJobStatus, 200),  // Increase queue size
time.Sleep(200 * time.Millisecond)       // Reduce delay
```

#### Memory Management
- Large content (>1MB) may cause delays
- Consider splitting large jobs
- Monitor memory usage with system tools

### Debug Commands

```bash
# Service health
curl -v http://localhost:8081/health

# List all endpoints
curl -X OPTIONS http://localhost:8081/api/

# Detailed printer info
curl -s http://localhost:8081/api/printers | jq

# Monitor job queue
watch -n 2 'curl -s http://localhost:8081/api/jobs?limit=5 | jq'
```

## 📚 Advanced Usage

### Production Deployment

#### Build for Production

```bash
# Current platform
go build -ldflags="-s -w" -o printer-service

# Cross-platform builds
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o printer-service.exe
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o printer-service-linux  
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o printer-service-mac
```

#### Service Management

**Windows Service (Optional):**
```batch
# Install as Windows service using NSSM
nssm install PrinterService "C:\path\to\printer-service.exe"
nssm set PrinterService DisplayName "Printer Service"
nssm set PrinterService Description "Local printer service for web applications"
nssm start PrinterService
```

**Linux Systemd:**
```ini
# /etc/systemd/system/printer-service.service
[Unit]
Description=Printer Service
After=network.target

[Service]
Type=simple
User=printer
WorkingDirectory=/opt/printer-service
ExecStart=/opt/printer-service/printer-service
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl enable printer-service
sudo systemctl start printer-service
sudo systemctl status printer-service
```

#### Security Configuration

For production use, consider these security measures:

```go
// Example security middleware (add to main.go)
func securityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // CORS restrictions
        allowedOrigins := []string{"http://localhost:3000", "https://yourdomain.com"}
        origin := r.Header.Get("Origin")
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        // Rate limiting (implement with middleware)
        // Authentication (implement as needed)
        
        next.ServeHTTP(w, r)
    })
}
```

### API Integration Patterns

#### Polling Pattern for Job Status

```javascript
class JobMonitor {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
  }

  async waitForCompletion(jobId, timeout = 30000) {
    const startTime = Date.now();
    
    while (Date.now() - startTime < timeout) {
      const response = await fetch(`${this.baseUrl}/jobs/${jobId}`);
      const data = await response.json();
      
      if (data.success && data.job) {
        const { status, progress } = data.job;
        
        if (status === 'completed') {
          return { success: true, job: data.job };
        } else if (status === 'failed') {
          return { success: false, error: data.job.error };
        }
        
        // Report progress
        this.onProgress?.(progress);
      }
      
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
    
    return { success: false, error: 'Timeout waiting for job completion' };
  }
}

// Usage
const monitor = new JobMonitor('http://localhost:8081/api');
monitor.onProgress = (progress) => console.log(`Progress: ${progress}%`);

const result = await monitor.waitForCompletion(jobId);
if (result.success) {
  console.log('Print job completed successfully');
} else {
  console.error('Print job failed:', result.error);
}
```

#### Batch Printing

```javascript
class BatchPrinter {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
    this.maxConcurrent = 3;
  }

  async printBatch(printerName, items) {
    const jobs = [];
    const chunks = this.chunkArray(items, this.maxConcurrent);
    
    for (const chunk of chunks) {
      const chunkJobs = await Promise.all(
        chunk.map(item => this.printSingle(printerName, item))
      );
      jobs.push(...chunkJobs);
      
      // Small delay between batches
      await new Promise(resolve => setTimeout(resolve, 500));
    }
    
    return jobs;
  }

  async printSingle(printerName, content) {
    const response = await fetch(`${this.baseUrl}/print/text`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ printerName, content })
    });
    
    const data = await response.json();
    return data.jobId;
  }

  chunkArray(array, size) {
    const chunks = [];
    for (let i = 0; i < array.length; i += size) {
      chunks.push(array.slice(i, i + size));
    }
    return chunks;
  }
}
```

### Custom ESC/POS Utilities

#### Receipt Builder Class

```javascript
class ReceiptBuilder {
  constructor() {
    this.commands = [];
    this.width = 48;
  }

  init() {
    this.commands.push('\\x1b@');
    return this;
  }

  text(content, options = {}) {
    if (options.bold) this.commands.push('\\x1b\\x45\\x01');
    if (options.size === 'large') this.commands.push('\\x1b!\\x30');
    if (options.align === 'center') this.commands.push('\\x1ba\\x01');
    
    this.commands.push(content);
    
    if (options.bold) this.commands.push('\\x1b\\x45\\x00');
    if (options.size === 'large') this.commands.push('\\x1b!\\x00');
    if (options.align) this.commands.push('\\x1ba\\x00');
    
    return this;
  }

  line() {
    this.commands.push('\\n');
    return this;
  }

  separator() {
    this.commands.push('='.repeat(this.width) + '\\n');
    return this;
  }

  table(rows) {
    rows.forEach(row => {
      const line = this.formatTableRow(row);
      this.commands.push(line + '\\n');
    });
    return this;
  }

  formatTableRow(cells) {
    const cellWidths = [20, 8, 12]; // Adjust as needed
    return cells.map((cell, i) => {
      const width = cellWidths[i] || 10;
      return String(cell).padEnd(width).substring(0, width);
    }).join('');
  }

  cut() {
    this.commands.push('\\x1b\\x69');
    return this;
  }

  build() {
    return this.commands.join('');
  }
}

// Usage
const receipt = new ReceiptBuilder()
  .init()
  .text('TECH STORE', { bold: true, size: 'large', align: 'center' })
  .line()
  .text('123 Main Street', { align: 'center' })
  .line()
  .separator()
  .table([
    ['Item', 'Qty', 'Price'],
    ['Laptop', '1', '$999.99'],
    ['Mouse', '2', '$50.00']
  ])
  .separator()
  .text('TOTAL: $1049.99', { bold: true, align: 'center' })
  .line()
  .text('Thank you!', { align: 'center' })
  .cut()
  .build();
```

## 🔒 Security & Best Practices

### Development vs Production

**Development Mode (Default):**
- CORS allows all origins (`*`)
- No authentication required
- Debug logging enabled
- All IPs accepted

**Production Recommendations:**
- Restrict CORS to specific domains
- Implement API key authentication
- Use HTTPS with reverse proxy
- Enable rate limiting
- Monitor and log access

### Network Security

```bash
# Firewall configuration (Linux)
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 8081/tcp    # Printer service (local only)
sudo ufw --force enable

# For production, consider running behind reverse proxy
# nginx configuration example:
# location /printer/ {
#     proxy_pass http://127.0.0.1:8081/;
#     proxy_set_header Host $host;
#     proxy_set_header X-Real-IP $remote_addr;
# }
```

### Resource Limits

```go
// Example: Add to main.go for production
func configureServer() *http.Server {
    return &http.Server{
        Addr:         ":8081",
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  60 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }
}
```

## 📞 Support & Contributing

### Documentation & Examples
- 📖 **Full API Documentation**: `API_DOCUMENTATION.md`
- 🌐 **Test Interface**: `http://localhost:8081/`
- 💾 **Source Code**: [GitHub Repository](https://github.com/Premod1/printer-service)
- 🔗 **Laravel Integration**: [Laravel Print Gateway](https://github.com/Premod1/laravel-local-print-gateway)

### Contributing
1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### License & Support
- **License**: MIT License
- **Issues**: Create GitHub issue with details
- **Questions**: Check documentation and examples first
- **Feature Requests**: Submit detailed GitHub issue

---

**Version**: 2.0.0  
**Last Updated**: December 11, 2025  
**Compatibility**: Go 1.21+, Windows 7+, macOS 10.9+, Linux (CUPS)  
**API Versions**: REST API v1, WebSocket v1
