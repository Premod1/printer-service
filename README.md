# Printer Service

A cross-platform WebSocket service for printing from web applications to local system printers. Supports both plain text and ESC/POS thermal receipt printing with 80mm receipt optimization.

## 🚀 Features

- **Cross-platform support** (Windows, macOS, Linux)
- **WebSocket-based API** for real-time communication
- **Multiple print formats**: Plain text and ESC/POS thermal receipts
- **80mm thermal printer optimization** with perfect formatting
- **Frontend invoice generation** with customizable templates
- **Vue.js ready** with composables and examples
- **Auto printer detection** across all platforms
- **Real-time connection status** and error handling

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [API Reference](#api-reference)
- [Vue.js Integration](#vuejs-integration)
- [ESC/POS Printing](#escpos-printing)
- [Platform Support](#platform-support)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## ⚡ Quick Start

### 1. Start the Service

```bash
# Clone and run the service
git clone <repository-url>
cd printer-service
go mod download
go run main.go
```

Service will start on `ws://localhost:8081/ws`

### 2. Test the Connection

Open `test.html` in your browser or use the Vue.js examples below.

## 📦 Installation

### Prerequisites

- **Go 1.21+**
- **System printer drivers installed**
- **CUPS** (Linux/macOS) or **Windows Print Spooler**

### Download Dependencies

```bash
go mod download
```

### Run the Service

```bash
go run main.go
```

### Build for Production

```bash
# Current platform
go build -o printer-service

# Cross-platform builds
GOOS=windows GOARCH=amd64 go build -o printer-service.exe
GOOS=linux GOARCH=amd64 go build -o printer-service-linux
GOOS=darwin GOARCH=amd64 go build -o printer-service-mac
```

## 🔌 API Reference

### WebSocket Endpoint
```
ws://localhost:8081/ws
```

### Message Format
```json
{
  "type": "message_type",
  "payload": { /* message data */ }
}
```

### Available Message Types

#### 1. Get Printers
**Request:**
```json
{
  "type": "get_printers"
}
```

**Response:**
```json
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

#### 2. Print Text
**Request:**
```json
{
  "type": "print",
  "payload": {
    "printerName": "HP LaserJet Pro",
    "content": "Hello World!\nThis is a test print.",
    "jobId": "job_123456"
  }
}
```

**Response:**
```json
{
  "type": "print_success",
  "payload": {
    "jobId": "job_123456"
  }
}
```

#### 3. Print ESC/POS (Raw Commands)
**Request:**
```json
{
  "type": "print_raw_escpos",
  "payload": {
    "printerName": "Receipt Printer",
    "jobId": "escpos_001",
    "rawData": "\x1b@\x1ba\x01RECEIPT\n\x1dV\x41\x00"
  }
}
```

**Response:**
```json
{
  "type": "raw_escpos_print_success",
  "payload": {
    "jobId": "escpos_001"
  }
}
```

#### 4. Error Response
```json
{
  "type": "error",
  "payload": {
    "message": "Printer not found: Invalid Printer"
  }
}
```

## 🔧 Vue.js Integration

### 1. Vue Composable (Recommended)

Create `composables/usePrinterService.js`:

```javascript
import { ref, onUnmounted } from 'vue'

export function usePrinterService() {
  const ws = ref(null)
  const isConnected = ref(false)
  const printers = ref([])
  const lastError = ref('')

  const connect = () => {
    if (ws.value) {
      ws.value.close()
    }

    ws.value = new WebSocket('ws://localhost:8081/ws')

    ws.value.onopen = () => {
      isConnected.value = true
      console.log('Connected to printer service')
    }

    ws.value.onmessage = (event) => {
      const message = JSON.parse(event.data)
      handleMessage(message)
    }

    ws.value.onclose = () => {
      isConnected.value = false
      console.log('Disconnected from printer service')
      
      // Auto-reconnect after 3 seconds
      setTimeout(() => {
        if (!isConnected.value) {
          connect()
        }
      }, 3000)
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
      lastError.value = 'Connection failed'
    }
  }

  const handleMessage = (message) => {
    switch (message.type) {
      case 'printers_list':
        printers.value = message.payload
        break
      case 'print_success':
        console.log('Print job completed:', message.payload.jobId)
        break
      case 'raw_escpos_print_success':
        console.log('ESC/POS print completed:', message.payload.jobId)
        break
      case 'error':
        lastError.value = message.payload.message
        console.error('Printer error:', message.payload.message)
        break
    }
  }

  const sendMessage = (message) => {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify(message))
    } else {
      lastError.value = 'Not connected to printer service'
    }
  }

  const getPrinters = () => {
    sendMessage({ type: 'get_printers' })
  }

  const printText = (printerName, content) => {
    const jobId = 'text_' + Date.now()
    sendMessage({
      type: 'print',
      payload: {
        printerName,
        content,
        jobId
      }
    })
    return jobId
  }

  const printReceipt = (printerName, receiptData) => {
    const escPosCommands = generateReceiptESCPOS(receiptData)
    const jobId = 'receipt_' + Date.now()
    
    sendMessage({
      type: 'print_raw_escpos',
      payload: {
        printerName,
        jobId,
        rawData: escPosCommands
      }
    })
    return jobId
  }

  // ESC/POS receipt generation for 80mm thermal printers
  const generateReceiptESCPOS = (data) => {
    const ESC_POS = {
      INIT: '\x1b@',
      BOLD_ON: '\x1bE\x01',
      BOLD_OFF: '\x1bE\x00',
      ALIGN_LEFT: '\x1ba\x00',
      ALIGN_CENTER: '\x1ba\x01',
      SIZE_NORMAL: '\x1d!\x00',
      SIZE_DOUBLE_HEIGHT: '\x1d!\x01',
      SIZE_DOUBLE_WIDTH: '\x1d!\x10',
      FONT_A: '\x1bM\x00',
      LINE_FEED: '\n',
      CUT_PAPER: '\x1dV\x41\x00',
      PAPER_WIDTH: 48
    }

    let escpos = ''
    const width = ESC_POS.PAPER_WIDTH
    const separator = '-'.repeat(width)

    // Initialize
    escpos += ESC_POS.INIT
    escpos += ESC_POS.FONT_A

    // Header
    escpos += ESC_POS.ALIGN_CENTER
    escpos += ESC_POS.BOLD_ON
    escpos += ESC_POS.SIZE_DOUBLE_HEIGHT
    escpos += (data.businessName || 'RESTAURANT RECEIPT')
    escpos += ESC_POS.LINE_FEED
    escpos += ESC_POS.SIZE_NORMAL
    escpos += ESC_POS.BOLD_OFF
    escpos += ESC_POS.LINE_FEED

    // Business details
    if (data.address) {
      escpos += data.address + ESC_POS.LINE_FEED
    }
    if (data.phone) {
      escpos += 'Phone: ' + data.phone + ESC_POS.LINE_FEED
    }
    escpos += ESC_POS.LINE_FEED

    // Invoice details
    escpos += ESC_POS.ALIGN_LEFT
    escpos += separator + ESC_POS.LINE_FEED
    escpos += `Invoice: ${data.invoiceCode || 'N/A'}` + ESC_POS.LINE_FEED
    escpos += `Table: ${data.tableNumber || 'N/A'}` + ESC_POS.LINE_FEED
    escpos += `Date: ${data.date || new Date().toLocaleString()}` + ESC_POS.LINE_FEED
    escpos += separator + ESC_POS.LINE_FEED

    // Items header
    escpos += ESC_POS.BOLD_ON
    escpos += 'ITEM                        QTY  PRICE'
    escpos += ESC_POS.LINE_FEED
    escpos += ESC_POS.BOLD_OFF
    escpos += separator + ESC_POS.LINE_FEED

    // Items
    data.items?.forEach((item, index) => {
      let itemName = item.name
      if (itemName.length > 26) {
        itemName = itemName.substring(0, 23) + '...'
      }
      
      const line = `${(index + 1).toString().padStart(2, '0')}. ${itemName.padEnd(24, ' ')} ${item.quantity.toString().padStart(2, ' ')} ${item.price.toFixed(2).padStart(6, ' ')}`
      escpos += line + ESC_POS.LINE_FEED
    })

    // Totals
    escpos += separator + ESC_POS.LINE_FEED
    if (data.subtotal) {
      escpos += `Subtotal:${('$' + data.subtotal.toFixed(2)).padStart(width - 9, ' ')}` + ESC_POS.LINE_FEED
    }
    if (data.discount) {
      escpos += `Discount:${('$' + data.discount).padStart(width - 9, ' ')}` + ESC_POS.LINE_FEED
    }
    if (data.tax) {
      escpos += `Tax:${('$' + data.tax.toFixed(2)).padStart(width - 4, ' ')}` + ESC_POS.LINE_FEED
    }
    escpos += separator + ESC_POS.LINE_FEED

    // Grand total
    escpos += ESC_POS.BOLD_ON
    const totalText = `TOTAL: $${data.total?.toFixed(2) || '0.00'}`
    escpos += ESC_POS.ALIGN_CENTER
    escpos += ESC_POS.SIZE_DOUBLE_WIDTH
    escpos += totalText + ESC_POS.LINE_FEED
    escpos += ESC_POS.SIZE_NORMAL
    escpos += ESC_POS.BOLD_OFF
    escpos += ESC_POS.ALIGN_LEFT

    // Payment info
    escpos += separator + ESC_POS.LINE_FEED
    if (data.cashPaid) {
      escpos += `Cash Paid:${('$' + data.cashPaid).padStart(width - 10, ' ')}` + ESC_POS.LINE_FEED
    }
    if (data.change !== undefined) {
      escpos += `Change:${('$' + data.change.toFixed(2)).padStart(width - 7, ' ')}` + ESC_POS.LINE_FEED
    }
    escpos += ESC_POS.LINE_FEED

    // Footer
    escpos += ESC_POS.ALIGN_CENTER
    escpos += separator + ESC_POS.LINE_FEED
    escpos += ESC_POS.BOLD_ON
    escpos += '★ THANK YOU FOR VISITING! ★' + ESC_POS.LINE_FEED
    escpos += ESC_POS.BOLD_OFF
    escpos += 'Please visit us again soon!' + ESC_POS.LINE_FEED
    escpos += separator + ESC_POS.LINE_FEED
    escpos += ESC_POS.LINE_FEED
    escpos += 'Customer Copy' + ESC_POS.LINE_FEED
    escpos += ESC_POS.LINE_FEED

    // Cut paper
    escpos += ESC_POS.CUT_PAPER

    return escpos
  }

  // Auto-connect on initialization
  connect()

  // Cleanup on unmount
  onUnmounted(() => {
    if (ws.value) {
      ws.value.close()
    }
  })

  return {
    // State
    isConnected,
    printers,
    lastError,
    
    // Methods
    connect,
    getPrinters,
    printText,
    printReceipt,
    generateReceiptESCPOS
  }
}
```

### 2. Vue Component Example

Create `components/PrinterManager.vue`:

```vue
<template>
  <div class="printer-manager">
    <div class="connection-status">
      <span :class="isConnected ? 'connected' : 'disconnected'">
        {{ isConnected ? '🟢 Connected' : '🔴 Disconnected' }}
      </span>
      <button @click="getPrinters" :disabled="!isConnected">
        Refresh Printers
      </button>
    </div>

    <div v-if="lastError" class="error">
      ⚠️ {{ lastError }}
    </div>

    <div class="printers-list">
      <h3>Available Printers:</h3>
      <div v-if="printers.length === 0" class="no-printers">
        No printers found. Click "Refresh Printers" to search.
      </div>
      
      <div 
        v-for="printer in printers" 
        :key="printer.name"
        class="printer-item"
      >
        <div class="printer-info">
          <strong>{{ printer.name }}</strong>
          <span class="status">{{ printer.status }}</span>
          <span v-if="printer.default" class="default-badge">Default</span>
        </div>
        
        <div class="printer-actions">
          <button @click="testPrint(printer.name)">
            Test Print
          </button>
          <button @click="printSampleReceipt(printer.name)">
            Print Sample Receipt
          </button>
        </div>
      </div>
    </div>

    <!-- Receipt Form -->
    <div class="receipt-form">
      <h3>Print Custom Receipt</h3>
      <form @submit.prevent="printCustomReceipt">
        <select v-model="selectedPrinter" required>
          <option value="">Select Printer</option>
          <option v-for="printer in printers" :key="printer.name" :value="printer.name">
            {{ printer.name }}
          </option>
        </select>

        <div class="form-group">
          <label>Business Name:</label>
          <input v-model="receiptData.businessName" />
        </div>

        <div class="form-group">
          <label>Invoice Code:</label>
          <input v-model="receiptData.invoiceCode" />
        </div>

        <div class="form-group">
          <label>Table Number:</label>
          <input v-model="receiptData.tableNumber" />
        </div>

        <div class="items-section">
          <h4>Items:</h4>
          <div v-for="(item, index) in receiptData.items" :key="index" class="item-row">
            <input v-model="item.name" placeholder="Item name" />
            <input v-model.number="item.quantity" type="number" min="1" />
            <input v-model.number="item.price" type="number" step="0.01" min="0" />
            <button type="button" @click="removeItem(index)">Remove</button>
          </div>
          <button type="button" @click="addItem">Add Item</button>
        </div>

        <div class="totals-section">
          <div class="form-group">
            <label>Subtotal: ${{ receiptData.subtotal.toFixed(2) }}</label>
          </div>
          <div class="form-group">
            <label>Tax:</label>
            <input v-model.number="receiptData.tax" type="number" step="0.01" min="0" />
          </div>
          <div class="form-group">
            <strong>Total: ${{ receiptData.total.toFixed(2) }}</strong>
          </div>
        </div>

        <button type="submit" :disabled="!selectedPrinter">
          Print Receipt
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { usePrinterService } from '@/composables/usePrinterService'

const {
  isConnected,
  printers,
  lastError,
  getPrinters,
  printText,
  printReceipt
} = usePrinterService()

const selectedPrinter = ref('')

const receiptData = ref({
  businessName: 'My Restaurant',
  invoiceCode: 'INV-001',
  tableNumber: '5',
  items: [
    { name: 'Coffee', quantity: 2, price: 3.50 },
    { name: 'Sandwich', quantity: 1, price: 8.00 }
  ],
  tax: 1.15,
  cashPaid: '20.00',
  change: 0
})

const subtotal = computed(() => {
  return receiptData.value.items.reduce((sum, item) => {
    return sum + (item.quantity * item.price)
  }, 0)
})

const total = computed(() => {
  return subtotal.value + receiptData.value.tax
})

// Update reactive values
receiptData.value.subtotal = subtotal
receiptData.value.total = total

const testPrint = (printerName) => {
  const testContent = `
TEST PRINT
==========
Printer: ${printerName}
Time: ${new Date().toLocaleString()}
This is a test print from Vue.js application.
  `.trim()
  
  printText(printerName, testContent)
}

const printSampleReceipt = (printerName) => {
  const sampleData = {
    businessName: 'Sample Restaurant',
    invoiceCode: 'SAMPLE-001',
    tableNumber: '10',
    date: new Date().toLocaleString(),
    items: [
      { name: 'Margherita Pizza', quantity: 1, price: 12.99 },
      { name: 'Caesar Salad', quantity: 1, price: 8.50 },
      { name: 'Coca Cola', quantity: 2, price: 2.50 }
    ],
    subtotal: 26.49,
    tax: 2.12,
    total: 28.61,
    cashPaid: '30.00',
    change: 1.39
  }
  
  printReceipt(printerName, sampleData)
}

const printCustomReceipt = () => {
  const data = {
    ...receiptData.value,
    date: new Date().toLocaleString(),
    subtotal: subtotal.value,
    total: total.value,
    change: parseFloat(receiptData.value.cashPaid) - total.value
  }
  
  printReceipt(selectedPrinter.value, data)
}

const addItem = () => {
  receiptData.value.items.push({ name: '', quantity: 1, price: 0 })
}

const removeItem = (index) => {
  receiptData.value.items.splice(index, 1)
}

// Load printers on mount
getPrinters()
</script>

<style scoped>
.printer-manager {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.connection-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.connected { color: green; }
.disconnected { color: red; }

.error {
  background: #fee;
  color: #c00;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.printer-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 10px;
}

.printer-info {
  display: flex;
  gap: 10px;
  align-items: center;
}

.default-badge {
  background: #007bff;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.printer-actions {
  display: flex;
  gap: 10px;
}

.receipt-form {
  margin-top: 30px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input, select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.item-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: center;
}

.item-row input {
  flex: 1;
}

button {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background: #0056b3;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.no-printers {
  text-align: center;
  color: #666;
  padding: 20px;
}
</style>
```

### 3. Usage in Vue App

```javascript
// main.js or in your component
import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)
app.mount('#app')

// In any component
import PrinterManager from '@/components/PrinterManager.vue'
```

## 🧾 ESC/POS Printing

### 80mm Thermal Printer Optimization

The service includes optimized ESC/POS commands for 80mm thermal printers:

- **48 characters per line** for perfect formatting
- **Font A optimization** (12x24 dots)
- **Professional receipt layout**
- **Bold headers and totals**
- **Right-aligned pricing**
- **Automatic paper cutting**

### ESC/POS Commands Reference

```javascript
const ESC_POS = {
  INIT: '\x1b@',                    // Initialize printer
  BOLD_ON: '\x1bE\x01',            // Bold text on
  BOLD_OFF: '\x1bE\x00',           // Bold text off
  ALIGN_CENTER: '\x1ba\x01',       // Center alignment
  ALIGN_LEFT: '\x1ba\x00',         // Left alignment
  SIZE_DOUBLE_HEIGHT: '\x1d!\x01', // Double height text
  SIZE_DOUBLE_WIDTH: '\x1d!\x10',  // Double width text
  FONT_A: '\x1bM\x00',             // Font A (optimal for 80mm)
  CUT_PAPER: '\x1dV\x41\x00',      // Cut paper
  LINE_FEED: '\n'                  // New line
}
```

## 🖥️ Platform Support

### Windows
- **Detection**: PowerShell + WMI
- **Printing**: Notepad command / Copy command for ESC/POS
- **Requirements**: Windows 7+, PowerShell enabled

### Linux
- **Detection**: CUPS (`lpstat -a`)
- **Printing**: CUPS (`lp -d <printer>`)
- **Requirements**: CUPS service running

### macOS
- **Detection**: CUPS (`lpstat -a`)
- **Printing**: CUPS (`lp -d <printer>`)
- **Requirements**: CUPS (built-in)

## 📱 Examples

### JavaScript/Vanilla JS

```javascript
const ws = new WebSocket('ws://localhost:8081/ws')

ws.onopen = () => {
  // Get printers
  ws.send(JSON.stringify({ type: 'get_printers' }))
}

ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  
  if (message.type === 'printers_list') {
    console.log('Available printers:', message.payload)
  }
}

// Print simple text
function printText(printerName, content) {
  ws.send(JSON.stringify({
    type: 'print',
    payload: {
      printerName,
      content,
      jobId: 'job_' + Date.now()
    }
  }))
}
```

### React Hook

```javascript
import { useState, useEffect, useRef } from 'react'

export function usePrinterService() {
  const [isConnected, setIsConnected] = useState(false)
  const [printers, setPrinters] = useState([])
  const wsRef = useRef(null)

  useEffect(() => {
    const connect = () => {
      wsRef.current = new WebSocket('ws://localhost:8081/ws')
      
      wsRef.current.onopen = () => setIsConnected(true)
      wsRef.current.onclose = () => setIsConnected(false)
      wsRef.current.onmessage = (event) => {
        const message = JSON.parse(event.data)
        if (message.type === 'printers_list') {
          setPrinters(message.payload)
        }
      }
    }

    connect()
    return () => wsRef.current?.close()
  }, [])

  const sendMessage = (message) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    }
  }

  const getPrinters = () => sendMessage({ type: 'get_printers' })
  
  const printText = (printerName, content) => {
    sendMessage({
      type: 'print',
      payload: { printerName, content, jobId: Date.now() }
    })
  }

  return { isConnected, printers, getPrinters, printText }
}
```

## 🔧 Troubleshooting

### Common Issues

#### 1. Connection Failed
- **Check if service is running**: `http://localhost:8081/health`
- **Verify port availability**: `lsof -i :8081` (Mac/Linux)
- **Firewall**: Allow port 8081

#### 2. No Printers Found
- **Windows**: Check Print Spooler service
- **Linux/Mac**: Verify CUPS: `systemctl status cups`
- **Permissions**: Ensure user can access printers

#### 3. Print Jobs Fail
- **Printer status**: Check if printer is online
- **Driver issues**: Reinstall printer drivers
- **Permissions**: Run service with appropriate permissions

#### 4. ESC/POS Not Working
- **Raw printing**: Ensure printer supports ESC/POS
- **USB connection**: Direct USB connection works best
- **Network printers**: May not support raw commands

### Debug Mode

```bash
DEBUG=true go run main.go
```

### Logs

```bash
# Save logs to file
go run main.go > printer-service.log 2>&1

# Monitor logs
tail -f printer-service.log
```

## 🔒 Security Notes

- **Local use only**: Service designed for localhost
- **No authentication**: Add auth for production use
- **CORS enabled**: All origins allowed (development)
- **File permissions**: Service needs write access for temp files

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

MIT License - see LICENSE file for details.

## 📞 Support

- 📧 Create an issue in the repository
- 📖 Check troubleshooting section above
- 💡 Review examples and API documentation

---

**Last Updated**: November 24, 2025  
**Version**: 1.0.0  
**Vue.js Compatible**: ✅ Vue 3 + Composition API
