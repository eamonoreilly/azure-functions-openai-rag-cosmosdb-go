// Build for linux if it is hosted in Azure Functions on Linux
// env GOOS=linux GOARCH=amd64 go build -o main .
// Build for Windows 32 bit if it is hosted in Azure Functions on Windows
// env GOOS=windows GOARCH=386 go build -o main .
// For local testing
// go build -o main .

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// InvokeRequest represents the structure of the request sent to the handle
type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

// InvokeResponse represents the structure of the response returned by the handler
type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type IngestRequest struct {
	Metadata map[string]interface{} `json:"Metadata"`
}

// promptHandler handles the /prompt endpoint
func promptHandler(w http.ResponseWriter, r *http.Request) {
	logs := []string{"Starting to process information for prompt call..."}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	logs = append(logs, "Received data (raw): "+string(body))

	var invokeReq InvokeRequest
	json.Unmarshal(body, &invokeReq)

	// Extract the "Response" value from SemanticSearchInput
	responseText := extractResponseFromSemanticSearchInput(invokeReq.Data)
	invokeResponse := InvokeResponse{
		Logs:        logs,
		ReturnValue: responseText,
	}

	responseJSON, _ := json.Marshal(invokeResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// supportHandler handles the /getsupportinformation MCP tool
func supportHandler(w http.ResponseWriter, r *http.Request) {
	logs := []string{"Starting to process information for getsupportinformation tool..."}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	logs = append(logs, "Received data (raw): "+string(body))

	var invokeReq InvokeRequest
	json.Unmarshal(body, &invokeReq)

	// Extract the "Response" value from SemanticSearchInput
	responseText := extractResponseFromSemanticSearchInput(invokeReq.Data)
	invokeResponse := InvokeResponse{
		Logs:        logs,
		ReturnValue: responseText,
	}

	responseJSON, _ := json.Marshal(invokeResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// extractResponseFromSemanticSearchInput extracts the "Response" field from the SemanticSearchInput JSON string.
func extractResponseFromSemanticSearchInput(data map[string]interface{}) string {
	semanticInputRaw, _ := data["SemanticSearchInput"].(string)

	var semanticInput string
	json.Unmarshal([]byte(semanticInputRaw), &semanticInput)

	var semanticObj map[string]interface{}
	json.Unmarshal([]byte(semanticInput), &semanticObj)

	resp, _ := semanticObj["Response"].(string)
	return resp
}

// ingestHandler handles the /ingest endpoint
func ingestHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize custom logs
	logs := []string{
		"Starting to process ingested content for Cosmos DB NoSQL...",
	}

	// Read the request body
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	logs = append(logs, "Received documents (raw): "+string(body))

	// Unmarshal the request body into a InjestRequest struct
	var injestReq IngestRequest
	json.Unmarshal(body, &injestReq)
	metadataJson, _ := json.Marshal(injestReq)
	logs = append(logs, "Extracted Metadata: "+string(metadataJson))

	returnValue := "The Cosmos DB documents have been processed successfully."

	// The title value is JSON encoded, so we need to decode it to get the actual string value
	var titleStr string
	if titleVal, ok := injestReq.Metadata["title"]; ok {
		// Convert titleVal to string
		titleStrRaw, ok := titleVal.(string)
		if ok {
			json.Unmarshal([]byte(titleStrRaw), &titleStr)
			if titleStr == "" {
				titleStr = titleStrRaw
			}
		}
	}

	logs = append(logs, "Title extracted: "+titleStr)
	// Send the value to the output binding name called "EmbeddingsStoreOutput" in function.json
	outputs := make(map[string]interface{})
	outputs["EmbeddingsStoreOutput"] = map[string]interface{}{"title": titleStr}

	invokeResponse := InvokeResponse{outputs, logs, returnValue}

	responseJson, _ := json.Marshal(invokeResponse)

	// Set the response headers and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

// main initializes the HTTP server and listens for incoming requests
func main() {
	// Determine the port to listen on
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	// Register the handler for the /ingest endpoint
	http.HandleFunc("/ingest", ingestHandler)

	// Register the handler for the /prompt endpoint
	http.HandleFunc("/prompt", promptHandler)

	// Register the handler for the /prompt endpoint
	http.HandleFunc("/getsupportinformation", supportHandler)

	// Start the HTTP server
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
