package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/danielspk/tatu-lang/pkg/builder"
	"github.com/danielspk/tatu-lang/pkg/interpreter"
	"github.com/danielspk/tatu-lang/pkg/runtime"

	embed "github.com/danielspk/blockly-tatu-transpiler"
	"github.com/danielspk/blockly-tatu-transpiler/pkg/blockly"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Starting Blockly-Tatu-Transpiler web interface on port %s...", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	webFS, err := fs.Sub(embed.WebFiles, "web")
	if err != nil {
		log.Fatalf("failed to get web files: %s", err)
	}

	http.Handle("/", http.FileServer(http.FS(webFS)))
	http.HandleFunc("/transpile", handleTranspile)
	http.HandleFunc("/execute", handleExecute)

	addr := fmt.Sprintf(":%s", port)
	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("listen web server: %s", err)
	}
}

// handleTranspile handles the POST /transpile endpoint.
func handleTranspile(w http.ResponseWriter, r *http.Request) {
	var response struct {
		JSON     string `json:"json"`
		TatuCode string `json:"tatu_code"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	tatuCode, err := blockly.Transpile(jsonData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transpilation error: %v", err), http.StatusBadRequest)
		return
	}

	response.JSON = string(jsonData)
	response.TatuCode = tatuCode

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleExecute handles the POST /execute endpoint.
func handleExecute(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Code   string         `json:"code"`
		Inputs map[string]any `json:"inputs"`
	}
	var result runtime.Value

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	progBuilder := builder.NewProgramBuilder(builder.NewDefaultScanner(), builder.NewDefaultParser())

	_, ast, err := progBuilder.BuildFromSource([]byte(request.Code), "<web-execution>")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing code: %v", err), http.StatusBadRequest)
		return
	}

	inter, err := interpreter.NewInterpreter()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error initializing interpreter: %v", err), http.StatusInternalServerError)
		return
	}

	env := inter.Environment()
	for k, v := range request.Inputs {
		env.Define(k, inputToValue(v))
	}

	result, err = inter.EvalProgram(ast, env)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if result != nil {
		fmt.Fprintf(w, "%v", result)
	}
}

func inputToValue(v any) runtime.Value {
	switch val := v.(type) {
	case float64:
		return runtime.NewNumber(val)
	case string:
		return runtime.NewString(val)
	case bool:
		return runtime.NewBool(val)
	case []any:
		items := make([]runtime.Value, len(val))
		for i, item := range val {
			items[i] = inputToValue(item)
		}
		return runtime.NewVector(items)
	case map[string]any:
		m := make(map[string]runtime.Value)
		for k, item := range val {
			m[k] = inputToValue(item)
		}
		return runtime.NewMap(m)
	case nil:
		return runtime.NewNil()
	default:
		return runtime.NewNil()
	}
}
