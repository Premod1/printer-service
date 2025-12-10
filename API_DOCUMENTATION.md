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
    "rawData": "\\x1b@\\x1b!\\x10RECEIPT TEST\\n\\x1b!\\x00\\nItem 1 ............ $10.00\\nItem 2 ............ $15.00\\n\\x1b!\\x10TOTAL: $25.00\\n\\x1b\\x64\\x05\\x1bi"
  }'
```

**Common ESC/POS Commands**:
- `\x1b@` - Initialize printer
- `\x1b!\x10` - Double height text
- `\x1b!\x00` - Normal text
- `\x1b\x64\x05` - Feed 5 lines
- `\x1bi` - Full cut

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

// Monitor job progress
async function monitorJob(jobId) {
    const checkStatus = async () => {
        const job = await getJobStatus(jobId);
        console.log(`Job ${jobId}: ${job.status} (${job.progress}%)`);
        
        if (job.status === 'completed') {
            console.log('Print job completed successfully!');
        } else if (job.status === 'failed') {
            console.error('Print job failed:', job.error);
        } else {
            // Still processing, check again in 1 second
            setTimeout(checkStatus, 1000);
        }
    };
    checkStatus();
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