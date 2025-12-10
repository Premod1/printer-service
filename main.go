package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"printer-service/config"
	"printer-service/printer"
	"printer-service/websocket"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Job management structures
type PrintJobStatus struct {
	JobID       string     `json:"jobId"`
	PrinterName string     `json:"printerName"`
	Content     string     `json:"content"`
	Status      string     `json:"status"` // "pending", "printing", "completed", "failed"
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	Error       string     `json:"error,omitempty"`
	JobType     string     `json:"jobType"`  // "text", "escpos"
	Progress    int        `json:"progress"` // 0-100
}

type JobManager struct {
	jobs  map[string]*PrintJobStatus
	queue chan *PrintJobStatus
	mutex sync.RWMutex
}

var jobManager = &JobManager{
	jobs:  make(map[string]*PrintJobStatus),
	queue: make(chan *PrintJobStatus, 100),
}

func (jm *JobManager) addJob(job *PrintJobStatus) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()
	jm.jobs[job.JobID] = job
	jm.queue <- job
}

func (jm *JobManager) getJob(jobID string) *PrintJobStatus {
	jm.mutex.RLock()
	defer jm.mutex.RUnlock()
	return jm.jobs[jobID]
}

func (jm *JobManager) getAllJobs() []*PrintJobStatus {
	jm.mutex.RLock()
	defer jm.mutex.RUnlock()
	var jobs []*PrintJobStatus
	for _, job := range jm.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (jm *JobManager) updateJobStatus(jobID string, status string, progress int, errorMsg string) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()
	if job, exists := jm.jobs[jobID]; exists {
		job.Status = status
		job.Progress = progress
		if errorMsg != "" {
			job.Error = errorMsg
		}
		if status == "completed" || status == "failed" {
			now := time.Now()
			job.CompletedAt = &now
		}
	}
}

// Job processor
func processJobs() {
	for job := range jobManager.queue {
		log.Printf("Processing job %s", job.JobID)

		jobManager.updateJobStatus(job.JobID, "printing", 0, "")

		var err error
		if job.JobType == "text" {
			jobManager.updateJobStatus(job.JobID, "printing", 50, "")
			err = printer.PrintText(job.PrinterName, job.Content)
		} else if job.JobType == "escpos" {
			jobManager.updateJobStatus(job.JobID, "printing", 50, "")
			err = printer.PrintEscPos(job.PrinterName, job.Content)
		}

		if err != nil {
			jobManager.updateJobStatus(job.JobID, "failed", 100, err.Error())
			log.Printf("Job %s failed: %v", job.JobID, err)
		} else {
			jobManager.updateJobStatus(job.JobID, "completed", 100, "")
			log.Printf("Job %s completed successfully", job.JobID)
		}

		// Small delay between jobs
		time.Sleep(500 * time.Millisecond)
	}
}

// REST API handlers
func getPrintersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	printers, err := printer.DetectPrinters()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"printers": printers,
	})
}

func printTextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		PrinterName string `json:"printerName"`
		Content     string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON request"}`, http.StatusBadRequest)
		return
	}

	if req.PrinterName == "" || req.Content == "" {
		http.Error(w, `{"error": "printerName and content are required"}`, http.StatusBadRequest)
		return
	}

	// Generate job ID
	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())

	// Create job
	job := &PrintJobStatus{
		JobID:       jobID,
		PrinterName: req.PrinterName,
		Content:     req.Content,
		Status:      "pending",
		CreatedAt:   time.Now(),
		JobType:     "text",
		Progress:    0,
	}

	jobManager.addJob(job)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"jobId":   jobID,
		"message": "Print job queued successfully",
	})
}

func printEscPosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		PrinterName string `json:"printerName"`
		RawData     string `json:"rawData"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON request"}`, http.StatusBadRequest)
		return
	}

	if req.PrinterName == "" || req.RawData == "" {
		http.Error(w, `{"error": "printerName and rawData are required"}`, http.StatusBadRequest)
		return
	}

	// Generate job ID
	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())

	// Create job
	job := &PrintJobStatus{
		JobID:       jobID,
		PrinterName: req.PrinterName,
		Content:     req.RawData,
		Status:      "pending",
		CreatedAt:   time.Now(),
		JobType:     "escpos",
		Progress:    0,
	}

	jobManager.addJob(job)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"jobId":   jobID,
		"message": "ESC/POS print job queued successfully",
	})
}

func getJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	jobID := vars["jobId"]

	job := jobManager.getJob(jobID)
	if job == nil {
		http.Error(w, `{"error": "Job not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"job":     job,
	})
}

func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get query parameters for filtering
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")

	allJobs := jobManager.getAllJobs()
	var jobs []*PrintJobStatus

	for _, job := range allJobs {
		if status == "" || job.Status == status {
			jobs = append(jobs, job)
		}
	}

	// Apply limit if specified
	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit < len(jobs) {
			jobs = jobs[:limit]
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"jobs":    jobs,
		"count":   len(jobs),
	})
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func serveTestUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Printer Service Test UI</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        h1, h2 {
            color: #333;
            border-bottom: 2px solid #007cba;
            padding-bottom: 10px;
        }
        .form-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input, select, textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 14px;
            box-sizing: border-box;
        }
        textarea {
            height: 100px;
            resize: vertical;
        }
        button {
            background-color: #007cba;
            color: white;
            padding: 12px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            margin-right: 10px;
        }
        button:hover {
            background-color: #005a8a;
        }
        button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        .status {
            padding: 10px;
            border-radius: 5px;
            margin: 10px 0;
        }
        .status.success {
            background-color: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .status.error {
            background-color: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .status.info {
            background-color: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }
        .jobs-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        .jobs-table th,
        .jobs-table td {
            border: 1px solid #ddd;
            padding: 12px;
            text-align: left;
        }
        .jobs-table th {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        .jobs-table tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .status-badge {
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: bold;
            text-transform: uppercase;
        }
        .status-pending {
            background-color: #fff3cd;
            color: #856404;
        }
        .status-printing {
            background-color: #cce5ff;
            color: #004085;
        }
        .status-completed {
            background-color: #d4edda;
            color: #155724;
        }
        .status-failed {
            background-color: #f8d7da;
            color: #721c24;
        }
        .progress-bar {
            width: 100%;
            background-color: #e9ecef;
            border-radius: 10px;
            overflow: hidden;
            height: 20px;
        }
        .progress-bar-fill {
            height: 100%;
            background-color: #007cba;
            transition: width 0.3s ease;
        }
        .two-column {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
        }
        @media (max-width: 768px) {
            .two-column {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🖨️ Printer Service Test UI</h1>
        <p>Test the enhanced printer service with job queuing and status tracking.</p>
    </div>

    <div class="two-column">
        <div class="container">
            <h2>📋 Available Printers</h2>
            <button onclick="loadPrinters()">🔄 Refresh Printers</button>
            <div id="printers-status"></div>
            <div id="printers-list"></div>
        </div>

        <div class="container">
            <h2>📝 Print Text</h2>
            <div class="form-group">
                <label for="text-printer">Printer:</label>
                <select id="text-printer">
                    <option value="">Select a printer...</option>
                </select>
            </div>
            <div class="form-group">
                <label for="text-content">Content:</label>
                <textarea id="text-content" placeholder="Enter text to print...">Hello from Printer Service!
This is a test print.
Current time: ${new Date().toLocaleString()}</textarea>
            </div>
            <button onclick="printText()">🖨️ Print Text</button>
            <div id="text-status"></div>
        </div>
    </div>

    <div class="container">
        <h2>⚡ Print ESC/POS</h2>
        <div class="form-group">
            <label for="escpos-printer">Printer:</label>
            <select id="escpos-printer">
                <option value="">Select a printer...</option>
            </select>
        </div>
        <div class="form-group">
            <label for="escpos-content">ESC/POS Commands:</label>
            <textarea id="escpos-content" placeholder="Enter ESC/POS commands...">\x1b@\x1b!\x10RECEIPT TEST\n\x1b!\x00\nTest Item 1 ............ $10.00\nTest Item 2 ............ $15.00\n\x1b!\x10TOTAL: $25.00\n\x1b\x64\x05\x1bi</textarea>
        </div>
        <button onclick="printEscPos()">⚡ Print ESC/POS</button>
        <div id="escpos-status"></div>
    </div>

    <div class="container">
        <h2>📊 Print Jobs</h2>
        <button onclick="loadJobs()">🔄 Refresh Jobs</button>
        <button onclick="autoRefresh = !autoRefresh; updateAutoRefreshButton()">⏰ Toggle Auto-refresh</button>
        <button onclick="clearJobs()">🗑️ Clear Display</button>
        <div id="jobs-status"></div>
        <div id="jobs-list"></div>
    </div>

    <script>
        let autoRefresh = false;
        let refreshInterval;

        function showStatus(elementId, message, type = 'info') {
            const element = document.getElementById(elementId);
            element.innerHTML = '<div class="status ' + type + '">' + message + '</div>';
        }

        async function loadPrinters() {
            try {
                showStatus('printers-status', 'Loading printers...', 'info');
                const response = await fetch('/api/printers');
                const data = await response.json();
                
                if (data.success) {
                    const printersHtml = data.printers.map(printer => 
                        '<div style="padding: 10px; border: 1px solid #ddd; margin: 5px 0; border-radius: 5px;">' +
                        '<strong>' + printer.name + '</strong><br>' +
                        'Status: ' + printer.status + '<br>' +
                        'Default: ' + (printer.default ? 'Yes' : 'No') +
                        '</div>'
                    ).join('');
                    
                    document.getElementById('printers-list').innerHTML = printersHtml;
                    
                    // Populate printer dropdowns
                    const options = data.printers.map(p => 
                        '<option value="' + p.name + '">' + p.name + '</option>'
                    ).join('');
                    document.getElementById('text-printer').innerHTML = 
                        '<option value="">Select a printer...</option>' + options;
                    document.getElementById('escpos-printer').innerHTML = 
                        '<option value="">Select a printer...</option>' + options;
                    
                    showStatus('printers-status', 'Found ' + data.printers.length + ' printers', 'success');
                } else {
                    showStatus('printers-status', 'Failed to load printers', 'error');
                }
            } catch (error) {
                showStatus('printers-status', 'Error: ' + error.message, 'error');
            }
        }

        async function printText() {
            const printer = document.getElementById('text-printer').value;
            const content = document.getElementById('text-content').value;
            
            if (!printer || !content) {
                showStatus('text-status', 'Please select a printer and enter content', 'error');
                return;
            }
            
            try {
                showStatus('text-status', 'Submitting print job...', 'info');
                const response = await fetch('/api/print/text', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        printerName: printer,
                        content: content
                    })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    showStatus('text-status', 'Job queued successfully! Job ID: ' + data.jobId, 'success');
                    loadJobs(); // Refresh job list
                } else {
                    showStatus('text-status', 'Print failed: ' + (data.error || 'Unknown error'), 'error');
                }
            } catch (error) {
                showStatus('text-status', 'Error: ' + error.message, 'error');
            }
        }

        async function printEscPos() {
            const printer = document.getElementById('escpos-printer').value;
            const content = document.getElementById('escpos-content').value;
            
            if (!printer || !content) {
                showStatus('escpos-status', 'Please select a printer and enter ESC/POS content', 'error');
                return;
            }
            
            try {
                showStatus('escpos-status', 'Submitting ESC/POS job...', 'info');
                const response = await fetch('/api/print/escpos', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        printerName: printer,
                        rawData: content
                    })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    showStatus('escpos-status', 'ESC/POS job queued! Job ID: ' + data.jobId, 'success');
                    loadJobs(); // Refresh job list
                } else {
                    showStatus('escpos-status', 'Print failed: ' + (data.error || 'Unknown error'), 'error');
                }
            } catch (error) {
                showStatus('escpos-status', 'Error: ' + error.message, 'error');
            }
        }

        async function loadJobs() {
            try {
                const response = await fetch('/api/jobs');
                const data = await response.json();
                
                if (data.success) {
                    if (data.jobs.length === 0) {
                        document.getElementById('jobs-list').innerHTML = 
                            '<p style="text-align: center; color: #666; font-style: italic;">No jobs found</p>';
                        showStatus('jobs-status', 'No jobs in queue', 'info');
                        return;
                    }

                    const jobsHtml = '<table class="jobs-table">' +
                        '<thead><tr>' +
                        '<th>Job ID</th>' +
                        '<th>Printer</th>' +
                        '<th>Type</th>' +
                        '<th>Status</th>' +
                        '<th>Progress</th>' +
                        '<th>Created</th>' +
                        '<th>Error</th>' +
                        '</tr></thead>' +
                        '<tbody>' +
                        data.jobs.map(job => {
                            const statusClass = 'status-' + job.status;
                            const progress = job.progress || 0;
                            return '<tr>' +
                                '<td>' + job.jobId.substring(0, 16) + '...</td>' +
                                '<td>' + job.printerName + '</td>' +
                                '<td>' + job.jobType.toUpperCase() + '</td>' +
                                '<td><span class="status-badge ' + statusClass + '">' + job.status + '</span></td>' +
                                '<td><div class="progress-bar"><div class="progress-bar-fill" style="width: ' + progress + '%"></div></div> ' + progress + '%</td>' +
                                '<td>' + new Date(job.createdAt).toLocaleString() + '</td>' +
                                '<td>' + (job.error || '-') + '</td>' +
                                '</tr>';
                        }).join('') +
                        '</tbody></table>';
                    
                    document.getElementById('jobs-list').innerHTML = jobsHtml;
                    showStatus('jobs-status', 'Loaded ' + data.jobs.length + ' jobs', 'success');
                } else {
                    showStatus('jobs-status', 'Failed to load jobs', 'error');
                }
            } catch (error) {
                showStatus('jobs-status', 'Error: ' + error.message, 'error');
            }
        }

        function clearJobs() {
            document.getElementById('jobs-list').innerHTML = '';
            showStatus('jobs-status', 'Display cleared', 'info');
        }

        function updateAutoRefreshButton() {
            const button = event.target;
            if (autoRefresh) {
                button.textContent = '⏰ Auto-refresh ON';
                button.style.backgroundColor = '#28a745';
                refreshInterval = setInterval(loadJobs, 2000);
            } else {
                button.textContent = '⏰ Auto-refresh OFF';
                button.style.backgroundColor = '#007cba';
                if (refreshInterval) {
                    clearInterval(refreshInterval);
                }
            }
        }

        // Initialize the page
        window.onload = function() {
            loadPrinters();
            loadJobs();
        };
    </script>
</body>
</html>
	`
	w.Write([]byte(html))
}

func main() {
	cfg := config.Load()

	// Start job processor
	go processJobs()

	// Create router
	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enableCORS(w, r)
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// REST API routes
	api := r.PathPrefix("/api").Subrouter()

	// Printer endpoints
	api.HandleFunc("/printers", getPrintersHandler).Methods("GET")

	// Print endpoints
	api.HandleFunc("/print/text", printTextHandler).Methods("POST")
	api.HandleFunc("/print/escpos", printEscPosHandler).Methods("POST")

	// Job management endpoints
	api.HandleFunc("/jobs", getJobsHandler).Methods("GET")
	api.HandleFunc("/jobs/{jobId}", getJobStatusHandler).Methods("GET")

	// WebSocket endpoint (existing functionality)
	r.HandleFunc("/ws", websocket.HandleWebSocket)

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Printer service is running"))
	})

	// Serve test UI
	r.HandleFunc("/", serveTestUI)

	log.Printf("Starting printer service on %s", cfg.WebSocketPort)
	log.Printf("REST API: http://localhost%s/api", cfg.WebSocketPort)
	log.Printf("WebSocket: ws://localhost%s/ws", cfg.WebSocketPort)
	log.Printf("Test UI: http://localhost%s/", cfg.WebSocketPort)

	if err := http.ListenAndServe(cfg.WebSocketPort, r); err != nil {
		log.Fatal("Server error:", err)
	}
}
