package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffcontext"
	"github.com/thomaspoignant/go-feature-flag/retriever/fileretriever"
)

func main() {
	// Initialize the feature flag client
	err := ffclient.Init(ffclient.Config{
		PollingInterval: 3 * time.Second,
		Retriever: &fileretriever.Retriever{
			Path: "./demo-flags.goff.yaml",
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize feature flag client: %v", err)
	}
	defer ffclient.Close()

	http.HandleFunc("/feature-flag/default", func(w http.ResponseWriter, r *http.Request) {
		flagKey := r.URL.Query().Get("flagKey")
		if flagKey == "" {
			http.Error(w, "flagKey is required", http.StatusBadRequest)
			return
		}

		// Context should be by company or user??
		evaluationContext := ffcontext.NewEvaluationContext("user-unique-key")

		// Evaluate the feature flag
		// TODO: Add error handling and evaluate multiple flags or diffent
		flagValue, err := ffclient.StringVariation(flagKey, evaluationContext, "default")
		if err != nil {
			http.Error(w, "Failed to get feature flag value", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"flagKey":   flagKey,
			"flagValue": flagValue,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/feature-flag", func(w http.ResponseWriter, r *http.Request) {
		// Context should be by company or user??
		evaluationContext := ffcontext.NewEvaluationContext("user-unique-key")

		flagsState := ffclient.AllFlagsState(evaluationContext)

		response := map[string]interface{}{
			"flagsState": flagsState,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
