package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
)

var currentPattern string
var lastRefresh time.Time

var patterns = []string{
	`s("bd*4")`,
	`s("bd cp")`,
	`s("hh*8")`,
	`note("c3 e3 g3")`,
	`s("bd*2, hh*4")`,
	`stack(s("bd*4"), s("hh*8"))`,
	`s("[bd, sd, hh]*4")`,
	`n("0 1 2 3").s("bd")`,
	`note("c3 e3 g3 c4").s("sawtooth")`,
	`s("bd sd").room(0.5)`,
}

func generatePattern(useAI bool, refresh int) string {
	if !useAI {
		return patterns[rand.Intn(len(patterns))]
	}

	if refresh > 0 && time.Since(lastRefresh) > time.Duration(refresh)*time.Second {
		lastRefresh = time.Now()
		return getAIPattern()
	}

	if currentPattern != "" {
		return currentPattern
	}

	return getAIPattern()
}

func getAIPattern() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	gatewayURL := os.Getenv("GATEWAY_URL")
	
	if apiKey == "" {
		log.Warn("OPENAI_API_KEY not set, using random pattern")
		return patterns[rand.Intn(len(patterns))]
	}

	if gatewayURL == "" {
		gatewayURL = "https://gateway.v.ai.vercel.ai/v1/chat/completions"
	}

	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type Request struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	reqBody := Request{
		Model: "gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Generate a short Strudel live coding pattern. Strudel is a live coding music language. Return ONLY the code, no explanation. Examples: s(\"bd*4\"), note(\"c3 e3 g3\"), stack(s(\"bd*4\"), s(\"hh*8\")). Make it interesting but short."},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", gatewayURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("AI gateway error", "error", err)
		return patterns[rand.Intn(len(patterns))]
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error("Parse error", "error", err)
		return patterns[rand.Intn(len(patterns))]
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return patterns[rand.Intn(len(patterns))]
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return patterns[rand.Intn(len(patterns))]
	}

	content, ok := message["content"].(string)
	if !ok {
		return patterns[rand.Intn(len(patterns))]
	}

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "```") {
			currentPattern = strings.Trim(line, "`")
			return currentPattern
		}
	}

	return patterns[rand.Intn(len(patterns))]
}

func main() {
	var (
		host    = flag.String("host", "0.0.0.0", "Host to listen on")
		port    = flag.Int("port", 23234, "Port to listen on")
		seed    = flag.Int64("seed", time.Now().UnixNano(), "Random seed")
		refresh = flag.Int("refresh", 0, "Refresh AI pattern every N seconds")
	)
	flag.Parse()

	useAI := os.Getenv("OPENAI_API_KEY") != ""
	rand.Seed(*seed)

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", *host, *port)),
		wish.WithHostKeyPath(".ssh/strudel_wish"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			logging.Middleware(),
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					pattern := generatePattern(useAI, *refresh)

					_, _, isPTY := s.Pty()
					
					if isPTY {
						lines := []string{
							"🌀  Strudel Wish - AI Music Stream",
							strings.Repeat("=", 35),
							"",
							"Now playing:",
							pattern,
							"",
							"Run locally:",
							"  ssh strudel@<host> | strudel",
							"",
							"(Ctrl+C to stop)",
							"",
						}
						output := strings.Join(lines, "\n")
						fmt.Fprint(s, output)
						fmt.Fprint(s, "\r\n")
						fmt.Fprint(s, pattern)
					} else {
						fmt.Fprint(s, pattern)
					}
				}
			},
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
		os.Exit(1)
	}

	log.Info("Strudel Wish server starting", "host", *host, "port", *port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info("Shutting down server...")
		s.Close()
	}()

	if err := s.ListenAndServe(); err != nil {
		log.Error("Server error", "error", err)
	}
}
