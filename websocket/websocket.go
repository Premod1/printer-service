package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"printer-service/printer"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type InvoiceData struct {
	InvoiceCode     string        `json:"invoiceCode"`
	TableNumber     string        `json:"tableNumber"`
	CurrentDateTime string        `json:"currentDateTime"`
	Items           []InvoiceItem `json:"items"`
	Subtotal        float64       `json:"subtotal"`
	Discount        string        `json:"discount"`
	ServiceFee      string        `json:"serviceFee"`
	GrandTotal      float64       `json:"grandTotal"`
	CashPaid        string        `json:"cashPaid"`
	CustomerBalance float64       `json:"customerBalance"`
}

type InvoiceItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type PrintJobEscPos struct {
	PrinterName string      `json:"printerName"`
	JobID       string      `json:"jobId"`
	Data        InvoiceData `json:"data"`
	Format      string      `json:"format"` // "text", "escpos", "pdf"
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	go client.writePump()
	client.readPump()
}

func (c *Client) readPump() {
	defer close(c.send)

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "get_printers":
		c.sendPrinters()
	case "print":
		c.handlePrint(msg.Payload)
	case "print_invoice":
		c.handleInvoicePrint(msg.Payload)
	}
}

func (c *Client) sendPrinters() {
	printers, err := printer.DetectPrinters()
	if err != nil {
		c.sendError(err.Error())
		return
	}

	response := Message{
		Type:    "printers_list",
		Payload: mustMarshal(printers),
	}

	c.conn.WriteJSON(response)
}

func (c *Client) handlePrint(payload json.RawMessage) {
	var printJob printer.PrintJob
	if err := json.Unmarshal(payload, &printJob); err != nil {
		c.sendError("Invalid print job format")
		return
	}

	err := printer.PrintText(printJob.PrinterName, printJob.Content)
	if err != nil {
		c.sendError(fmt.Sprintf("Print failed: %v", err))
		return
	}

	response := Message{
		Type:    "print_success",
		Payload: mustMarshal(map[string]string{"jobId": printJob.JobID}),
	}

	c.conn.WriteJSON(response)
}

func (c *Client) handleInvoicePrint(payload json.RawMessage) {
	var printJob PrintJobEscPos
	if err := json.Unmarshal(payload, &printJob); err != nil {
		c.sendError("Invalid invoice format")
		return
	}

	// For now, convert invoice data to simple text format
	// TODO: Implement ESC/POS generation when escpos.go is working
	content := c.generateInvoiceText(printJob.Data)

	err := printer.PrintText(printJob.PrinterName, content)
	if err != nil {
		c.sendError(fmt.Sprintf("Print failed: %v", err))
		return
	}

	response := Message{
		Type:    "print_success",
		Payload: mustMarshal(map[string]string{"jobId": printJob.JobID}),
	}

	c.conn.WriteJSON(response)
}

// Simple text invoice generator (temporary until ESC/POS is working)
func (c *Client) generateInvoiceText(data InvoiceData) string {
	var content string
	content += "========================================\n"
	content += "           RECEIPT                      \n"
	content += "========================================\n"
	content += fmt.Sprintf("Invoice: %s\n", data.InvoiceCode)
	content += fmt.Sprintf("Date: %s\n", data.CurrentDateTime)
	if data.TableNumber != "" {
		content += fmt.Sprintf("Table: %s\n", data.TableNumber)
	}
	content += "----------------------------------------\n"
	content += "ITEMS:\n"
	content += "----------------------------------------\n"

	for _, item := range data.Items {
		total := float64(item.Quantity) * item.Price
		content += fmt.Sprintf("%-20s %2d x %6.2f = %7.2f\n",
			item.Name, item.Quantity, item.Price, total)
	}

	content += "----------------------------------------\n"
	content += fmt.Sprintf("Subtotal: %30.2f\n", data.Subtotal)
	content += fmt.Sprintf("Total: %33.2f\n", data.GrandTotal)
	content += "========================================\n"
	content += "         THANK YOU!                     \n"
	content += "========================================\n"

	return content
}

func (c *Client) sendError(message string) {
	response := Message{
		Type:    "error",
		Payload: mustMarshal(map[string]string{"message": message}),
	}

	c.conn.WriteJSON(response)
}

func mustMarshal(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
