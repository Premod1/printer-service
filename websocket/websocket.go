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

type PrintJobEscPos struct {
	PrinterName string      `json:"printerName"`
	JobID       string      `json:"jobId"`
	Data        interface{} `json:"data"`   // Keep as interface{} for backward compatibility
	Format      string      `json:"format"` // "text", "escpos", "pdf"
}

type RawEscPosJob struct {
	PrinterName string `json:"printerName"`
	JobID       string `json:"jobId"`
	RawData     string `json:"rawData"` // Raw ESC/POS commands
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
	fmt.Printf("Received message: %s\n", msg.Type)
	switch msg.Type {
	case "get_printers":
		c.sendPrinters()
	case "print":
		c.handlePrint(msg.Payload)
	case "print_escpos":
		c.handlePrintEscPos(msg.Payload)
	case "print_raw_escpos":
		c.handlePrintRawEscPos(msg.Payload)
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

func (c *Client) handlePrintEscPos(payload json.RawMessage) {
	// This endpoint is deprecated - use print_raw_escpos instead
	c.sendError("This endpoint is deprecated. Use 'print_raw_escpos' with raw ESC/POS commands generated from frontend.")
}

func (c *Client) handlePrintRawEscPos(payload json.RawMessage) {
	var printJob RawEscPosJob
	if err := json.Unmarshal(payload, &printJob); err != nil {
		c.sendError("Invalid raw ESC/POS print job format")
		return
	}

	err := printer.PrintEscPos(printJob.PrinterName, printJob.RawData)
	if err != nil {
		c.sendError(fmt.Sprintf("Raw ESC/POS print failed: %v", err))
		return
	}

	response := Message{
		Type:    "raw_escpos_print_success",
		Payload: mustMarshal(map[string]string{"jobId": printJob.JobID}),
	}

	c.conn.WriteJSON(response)
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
