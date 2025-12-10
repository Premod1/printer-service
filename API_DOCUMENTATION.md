# Printer Service REST API Documentation

## Overview

The Printer Service provides a RESTful API for managing print jobs with queuing and status tracking capabilities. This service supports both text printing and ESC/POS thermal printing across Windows, macOS, and Linux platforms.

**Base URL**: `http://localhost:8081/api`

## Authentication

Currently, no authentication is required. CORS is enabled for cross-origin requests.

## Response Format

All responses follow this standard format:

```json
{
  "success": true|false,
  "data": {...},
  "error": "error message (if applicable)"
}
```

---

## Endpoints

### 1. Get Available Printers

Retrieve a list of all available printers on the system.

**Endpoint**: `GET /api/printers`

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

**cURL Example**:
```bash
curl -X GET http://localhost:8081/api/printers
```

**Status Codes**:
- `200 OK` - Success
- `500 Internal Server Error` - Failed to detect printers

---

### 2. Print Text Content

Queue a text printing job.

**Endpoint**: `POST /api/print/text`

**Request Body**:
```json
{
  "printerName": "string",
  "content": "string"
}
```

**Response**:
```json
{
  "success": true,
  "jobId": "job_1765360872587975241",
  "message": "Print job queued successfully"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8081/api/print/text \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "HP LaserJet Pro",
    "content": "Hello World!\nThis is a test print.\nTimestamp: 2025-12-10 15:30:00"
  }'
```

**Status Codes**:
- `200 OK` - Job queued successfully
- `400 Bad Request` - Missing required fields
- `500 Internal Server Error` - Server error

---

### 3. Print ESC/POS Content

Queue an ESC/POS printing job for thermal printers.

**Endpoint**: `POST /api/print/escpos`

**Request Body**:
```json
{
  "printerName": "string",
  "rawData": "string"
}
```

**Response**:
```json
{
  "success": true,
  "jobId": "job_1765360745910409307",
  "message": "ESC/POS print job queued successfully"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8081/api/print/escpos \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "Thermal Receipt Printer",
    "rawData": "\\x1b@\\x1ba\\x01\\x1b!\\x30RECEIPT\\n\\x1b!\\x00\\x1ba\\x00Item 1\\x09\\x09$10.00\\nItem 2\\x09\\x09$15.00\\n\\x1b\\x64\\x02\\x1b\\x45\\x01TOTAL: $25.00\\x1b\\x45\\x00\\n\\x1b\\x64\\x05\\x1b\\x69"
  }'
```

**Complete Receipt Example**:
```bash
curl -X POST http://localhost:8081/api/print/escpos \
  -H "Content-Type: application/json" \
  -d '{
    "printerName": "POS-58",
    "rawData": "\\x1b@\\x1ba\\x01\\x1b!\\x30MY STORE\\n\\x1b!\\x00\\x1ba\\x01123 Business St\\nCity, State 12345\\nTel: (555) 123-4567\\n\\x1b\\x64\\x02\\x1ba\\x00================================\\n\\x1bTIME: 2025-12-10 15:30:00\\nCASHIER: John Doe\\nRECEIPT #: 001234\\n================================\\n\\x1b\\x45\\x01ITEM\\x09\\x09QTY\\x09PRICE\\x1b\\x45\\x00\\n--------------------------------\\nCoffee\\x09\\x091\\x09$3.50\\nSandwich\\x09\\x091\\x09$8.99\\nChips\\x09\\x091\\x09$2.50\\n--------------------------------\\n\\x1b\\x45\\x01SUBTOTAL:\\x09\\x09$14.99\\nTAX:\\x09\\x09\\x09$1.20\\nTOTAL:\\x09\\x09\\x09$16.19\\x1b\\x45\\x00\\n================================\\n\\x1ba\\x01PAYMENT: CASH $20.00\\nCHANGE: $3.81\\n\\x1b\\x64\\x02\\x1ba\\x01Thank you for shopping!\\nPlease come again!\\n\\x1b\\x64\\x05\\x1b\\x69"
  }'
```

**ESC/POS Commands Reference**:

| Command | Description | Example |
|---------|-------------|---------|
| `\x1b@` | Initialize printer | `\x1b@` |
| `\x1b!\x00` | Normal text | `\x1b!\x00Normal Text` |
| `\x1b!\x10` | Double height | `\x1b!\x10BIG TEXT` |
| `\x1b!\x20` | Double width | `\x1b!\x20WIDE TEXT` |
| `\x1b!\x30` | Double height + width | `\x1b!\x30LARGE` |
| `\x1b!\x08` | Bold text | `\x1b!\x08Bold Text` |
| `\x1ba\x00` | Left align | `\x1ba\x00Left` |
| `\x1ba\x01` | Center align | `\x1ba\x01Center` |
| `\x1ba\x02` | Right align | `\x1ba\x02Right` |
| `\x1b\x45\x01` | Emphasize ON | `\x1b\x45\x01Bold` |
| `\x1b\x45\x00` | Emphasize OFF | `\x1b\x45\x00Normal` |
| `\x1b\x2d\x01` | Underline ON | `\x1b\x2d\x01Underline` |
| `\x1b\x2d\x00` | Underline OFF | `\x1b\x2d\x00Normal` |
| `\x1b\x64\x02` | Feed 2 lines | `\x1b\x64\x02` |
| `\x1b\x6d` | Partial cut | `\x1b\x6d` |
| `\x1b\x69` | Full cut | `\x1b\x69` |
| `\x1b\x61\x00` | Left justify | `\x1b\x61\x00` |
| `\x1b\x61\x01` | Center justify | `\x1b\x61\x01` |
| `\x1b\x61\x02` | Right justify | `\x1b\x61\x02` |

**Advanced ESC/POS Examples**:

1. **Receipt Header**:
```
\x1b@\x1ba\x01\x1b!\x30STORE NAME\n\x1b!\x00\x1ba\x01123 Main St\nCity, State 12345\n\x1b\x64\x02
```

2. **Item Line with Price Alignment**:
```
\x1ba\x00Item Name\x09\x09\x09$10.99\n
```

3. **Bold Total Line**:
```
\x1b\x45\x01\x1ba\x02TOTAL: $25.99\x1b\x45\x00\n
```

4. **Barcode (Code128)**:
```
\x1d\x6b\x49\x0c123456789012\x00
```

**Status Codes**:
- `200 OK` - Job queued successfully
- `400 Bad Request` - Missing required fields
- `500 Internal Server Error` - Server error

---

### 4. Get All Print Jobs

Retrieve all print jobs with optional filtering.

**Endpoint**: `GET /api/jobs`

**Query Parameters**:
- `status` (optional) - Filter by job status: `pending`, `printing`, `completed`, `failed`
- `limit` (optional) - Limit number of results (integer)

**Response**:
```json
{
  "success": true,
  "jobs": [
    {
      "jobId": "job_1765360872587975241",
      "printerName": "HP LaserJet Pro",
      "content": "Hello World Test Print",
      "status": "completed",
      "createdAt": "2025-12-10T15:31:12.587977806+05:30",
      "completedAt": "2025-12-10T15:31:12.665607362+05:30",
      "error": "",
      "jobType": "text",
      "progress": 100
    }
  ],
  "count": 1
}
```

**cURL Examples**:
```bash
# Get all jobs
curl -X GET http://localhost:8081/api/jobs

# Get only completed jobs
curl -X GET "http://localhost:8081/api/jobs?status=completed"

# Get last 5 jobs
curl -X GET "http://localhost:8081/api/jobs?limit=5"

# Get pending jobs with limit
curl -X GET "http://localhost:8081/api/jobs?status=pending&limit=10"
```

**Status Codes**:
- `200 OK` - Success
- `500 Internal Server Error` - Server error

---

### 5. Get Specific Job Status

Retrieve detailed information about a specific print job.

**Endpoint**: `GET /api/jobs/{jobId}`

**Path Parameters**:
- `jobId` - The unique job identifier

**Response**:
```json
{
  "success": true,
  "job": {
    "jobId": "job_1765360872587975241",
    "printerName": "HP LaserJet Pro",
    "content": "Hello World Test Print",
    "status": "completed",
    "createdAt": "2025-12-10T15:31:12.587977806+05:30",
    "completedAt": "2025-12-10T15:31:12.665607362+05:30",
    "error": "",
    "jobType": "text",
    "progress": 100
  }
}
```

**cURL Example**:
```bash
curl -X GET http://localhost:8081/api/jobs/job_1765360872587975241
```

**Status Codes**:
- `200 OK` - Job found
- `404 Not Found` - Job not found
- `500 Internal Server Error` - Server error

---

## Job Status States

| Status | Description |
|--------|-------------|
| `pending` | Job is queued and waiting to be processed |
| `printing` | Job is currently being sent to the printer |
| `completed` | Job finished successfully |
| `failed` | Job failed with an error |

## Progress Values

Progress is represented as an integer from 0 to 100:
- `0` - Job just created
- `50` - Job is being processed
- `100` - Job completed (success or failure)

---

---

## ESC/POS Command Guide

### Character Encoding

ESC/POS commands must be properly escaped in JSON:

| Raw Command | JSON Escaped | Description |
|-------------|--------------|-------------|
| `ESC @` | `\\x1b@` | Initialize printer |
| `ESC !` | `\\x1b!` | Text attributes |
| `ESC a` | `\\x1ba` | Alignment |
| `ESC d` | `\\x1b\\x64` | Line feed |
| `GS k` | `\\x1d\\x6b` | Barcode |
| `Tab` | `\\x09` or `\\t` | Horizontal tab |
| `Newline` | `\\n` | Line break |

### Text Formatting Commands

**Font Styles**:
```javascript
const styles = {
    normal: "\\x1b!\\x00",
    bold: "\\x1b!\\x08",
    doubleHeight: "\\x1b!\\x10", 
    doubleWidth: "\\x1b!\\x20",
    large: "\\x1b!\\x30",
    underline: "\\x1b\\x2d\\x01",
    noUnderline: "\\x1b\\x2d\\x00"
};
```

**Text Alignment**:
```javascript
const alignment = {
    left: "\\x1ba\\x00",
    center: "\\x1ba\\x01", 
    right: "\\x1ba\\x02"
};
```

### Paper Control

**Line Feeds and Cuts**:
```javascript
const paperControl = {
    feedLine: "\\n",
    feed2Lines: "\\x1b\\x64\\x02",
    feed5Lines: "\\x1b\\x64\\x05",
    partialCut: "\\x1b\\x6d",
    fullCut: "\\x1b\\x69"
};
```

### Advanced Features

**Barcodes** (Code128):
```javascript
function createBarcode(data) {
    return `\\x1d\\x6b\\x49${String.fromCharCode(data.length)}${data}`;
}
```

**QR Codes** (if supported):
```javascript
function createQRCode(data) {
    const size = Math.min(data.length + 3, 255);
    return `\\x1d(k\\x04\\x00\\x31\\x41\\x32\\x00` +  // QR Model
           `\\x1d(k\\x03\\x00\\x31\\x43\\x03` +        // Error correction
           `\\x1d(k${String.fromCharCode(size)}\\x00\\x31\\x50\\x30${data}` + // Data
           `\\x1d(k\\x03\\x00\\x31\\x51\\x30`;         // Print
}
```

**Tables and Columns**:
```javascript
function formatTableRow(col1, col2, col3, width = 32) {
    const col1Width = 12;
    const col2Width = 8;
    const col3Width = 12;
    
    return col1.padEnd(col1Width).substring(0, col1Width) + 
           col2.padStart(col2Width).substring(0, col2Width) + 
           col3.padStart(col3Width).substring(0, col3Width) + "\\n";
}

// Usage
let tableData = formatTableRow("Item", "Qty", "Price");
tableData += formatTableRow("Coffee", "2", "$7.00");
```

## Error Handling

### Common Error Responses

**400 Bad Request**:
```json
{
  "error": "printerName and content are required"
}
```

**404 Not Found**:
```json
{
  "error": "Job not found"
}
```

**500 Internal Server Error**:
```json
{
  "error": "Failed to detect printers"
}
```

### Print Job Errors

When a print job fails, the error will be captured in the job's `error` field:

```json
{
  "jobId": "job_123",
  "status": "failed",
  "error": "print failed: exit status 1, output: lp: Error - The printer or class does not exist.",
  "progress": 100
}
```

---

## Integration Examples

### JavaScript/Web Application

```javascript
// Get printers
async function getPrinters() {
    const response = await fetch('/api/printers');
    const data = await response.json();
    return data.printers;
}

// Print text
async function printText(printerName, content) {
    const response = await fetch('/api/print/text', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            printerName: printerName,
            content: content
        })
    });
    const data = await response.json();
    return data.jobId;
}

// Check job status
async function getJobStatus(jobId) {
    const response = await fetch(`/api/jobs/${jobId}`);
    const data = await response.json();
    return data.job;
}

// Print ESC/POS receipt
async function printEscPosReceipt(printerName, receiptData) {
    const escPosCommands = generateEscPosReceipt(receiptData);
    const response = await fetch('/api/print/escpos', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            printerName: printerName,
            rawData: escPosCommands
        })
    });
    const data = await response.json();
    return data.jobId;
}

// Generate ESC/POS commands for a receipt
function generateEscPosReceipt(data) {
    let commands = '';
    
    // Initialize and center store name
    commands += '\\x1b@';  // Initialize
    commands += '\\x1ba\\x01';  // Center align
    commands += `\\x1b!\\x30${data.storeName}\\n`;  // Large store name
    
    // Store info
    commands += '\\x1b!\\x00\\x1ba\\x01';  // Normal text, center
    commands += `${data.address}\\n`;
    commands += `${data.phone}\\n`;
    commands += '\\x1b\\x64\\x02';  // Feed 2 lines
    
    // Receipt header
    commands += '\\x1ba\\x00';  // Left align
    commands += '================================\\n';
    commands += `DATE: ${data.date}\\n`;
    commands += `CASHIER: ${data.cashier}\\n`;
    commands += `RECEIPT #: ${data.receiptNumber}\\n`;
    commands += '================================\\n';
    
    // Items header
    commands += '\\x1b\\x45\\x01ITEM\\x09QTY\\x09PRICE\\x1b\\x45\\x00\\n';
    commands += '--------------------------------\\n';
    
    // Items
    data.items.forEach(item => {
        commands += `${item.name}\\x09${item.qty}\\x09$${item.price}\\n`;
    });
    
    // Totals
    commands += '--------------------------------\\n';
    commands += `\\x1b\\x45\\x01SUBTOTAL:\\x09$${data.subtotal}\\n`;
    commands += `TAX:\\x09\\x09$${data.tax}\\n`;
    commands += `TOTAL:\\x09\\x09$${data.total}\\x1b\\x45\\x00\\n`;
    commands += '================================\\n';
    
    // Payment info
    commands += '\\x1ba\\x01';  // Center align
    commands += `PAYMENT: ${data.paymentMethod} $${data.payment}\\n`;
    if (data.change > 0) {
        commands += `CHANGE: $${data.change}\\n`;
    }
    
    // Footer
    commands += '\\x1b\\x64\\x02';  // Feed 2 lines
    commands += 'Thank you for shopping!\\n';
    commands += 'Please come again!\\n';
    commands += '\\x1b\\x64\\x05';  // Feed 5 lines
    commands += '\\x1b\\x69';  // Cut paper
    
    return commands;
}

// Example usage
const receiptData = {
    storeName: "TECH STORE",
    address: "123 Tech Street\\nSilicon Valley, CA 94000",
    phone: "Tel: (555) 123-4567",
    date: new Date().toLocaleString(),
    cashier: "Alice Johnson",
    receiptNumber: "R001234",
    items: [
        { name: "Laptop", qty: 1, price: "999.99" },
        { name: "Mouse", qty: 2, price: "25.00" },
        { name: "Keyboard", qty: 1, price: "75.00" }
    ],
    subtotal: "1099.99",
    tax: "87.99",
    total: "1187.98",
    paymentMethod: "CREDIT CARD",
    payment: "1187.98",
    change: 0
};

// Print the receipt
printEscPosReceipt("Thermal Printer", receiptData)
    .then(jobId => {
        console.log("Receipt job queued:", jobId);
        monitorJob(jobId);
    });

// QR Code generation (if printer supports it)
function generateQRCode(data) {
    // QR Code ESC/POS commands
    return `\\x1d(k\\x04\\x00\\x31\\x41\\x32\\x00` + // QR Code model
           `\\x1d(k\\x03\\x00\\x31\\x43\\x03` + // Error correction level
           `\\x1d(k${String.fromCharCode(data.length + 3)}\\x00\\x31\\x50\\x30${data}` + // Store QR data
           `\\x1d(k\\x03\\x00\\x31\\x51\\x30`; // Print QR code
}

// Barcode generation (Code128)
function generateBarcode(data) {
    return `\\x1d\\x6b\\x49${String.fromCharCode(data.length)}${data}\\x00`;
}
```

### Python Application

```python
import requests
import time
import json

class PrinterService:
    def __init__(self, base_url="http://localhost:8081/api"):
        self.base_url = base_url
    
    def get_printers(self):
        response = requests.get(f"{self.base_url}/printers")
        return response.json()
    
    def print_text(self, printer_name, content):
        data = {
            "printerName": printer_name,
            "content": content
        }
        response = requests.post(f"{self.base_url}/print/text", json=data)
        return response.json()
    
    def print_escpos(self, printer_name, raw_data):
        data = {
            "printerName": printer_name,
            "rawData": raw_data
        }
        response = requests.post(f"{self.base_url}/print/escpos", json=data)
        return response.json()
    
    def get_job_status(self, job_id):
        response = requests.get(f"{self.base_url}/jobs/{job_id}")
        return response.json()
    
    def wait_for_job(self, job_id, timeout=30):
        start_time = time.time()
        while time.time() - start_time < timeout:
            job_response = self.get_job_status(job_id)
            if job_response["success"]:
                job = job_response["job"]
                if job["status"] in ["completed", "failed"]:
                    return job
            time.sleep(1)
        return None

# Usage example
printer_service = PrinterService()

# List printers
printers = printer_service.get_printers()
print("Available printers:", printers)

# Submit print job
result = printer_service.print_text("HP LaserJet Pro", "Test print from Python")
if result["success"]:
    job_id = result["jobId"]
    print(f"Job queued: {job_id}")
    
    # Wait for completion
    final_job = printer_service.wait_for_job(job_id)
    if final_job:
        print(f"Job completed with status: {final_job['status']}")
        if final_job["status"] == "failed":
            print(f"Error: {final_job['error']}")
```

---

## Rate Limiting & Best Practices

### Recommendations

1. **Job Polling**: When monitoring job status, poll every 1-2 seconds rather than continuously
2. **Batch Operations**: For multiple print jobs, submit them sequentially to avoid overwhelming the queue
3. **Error Handling**: Always check the `success` field in responses
4. **Printer Validation**: Verify printer exists before submitting jobs
5. **Content Validation**: Ensure content is not empty and properly formatted

### Performance Considerations

- The job queue can handle up to 100 concurrent jobs
- Each job is processed sequentially to prevent printer conflicts
- Large content (>1MB) may cause slower processing
- ESC/POS jobs typically process faster than plain text

---

## Testing & Development

### Health Check

**Endpoint**: `GET /health`

```bash
curl -X GET http://localhost:8081/health
# Response: "Printer service is running"
```

### Test UI

Access the built-in test interface at: `http://localhost:8081/`

The test UI provides:
- Printer discovery and selection
- Interactive print job submission
- Real-time job monitoring
- Auto-refresh capabilities

---

## WebSocket Alternative

For real-time applications, the service also provides WebSocket support at `ws://localhost:8081/ws`. See WebSocket documentation for message formats and real-time job updates.

## Support

For issues or questions about the API, please refer to the project repository or submit an issue with:
- API endpoint used
- Request/response examples
- Error messages
- System information (OS, printer type)