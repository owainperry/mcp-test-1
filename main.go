// Command mcp-test-1 is a minimal MCP server used to exercise an MCP gateway.
//
// It exposes a single tool, "hello", which echoes back "Hello, <text>!" for
// whatever text it is given. The server speaks MCP over the streamable HTTP
// transport at /mcp, and also serves a plain /healthz probe.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// version is stamped at build time with -ldflags "-X main.version=...".
var version = "dev"

// HelloParams defines the parameters for the hello tool.
type HelloParams struct {
	Text string `json:"text" jsonschema:"Text to greet; anything you like"`
}

// hello implements the tool: it returns "hello" plus whatever text it is given.
func hello(ctx context.Context, req *mcp.CallToolRequest, params *HelloParams) (*mcp.CallToolResult, any, error) {
	text := strings.TrimSpace(params.Text)

	greeting := "Hello!"
	if text != "" {
		greeting = fmt.Sprintf("Hello, %s!", text)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: greeting},
		},
	}, nil, nil
}

func newServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-test-1",
		Version: version,
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "hello",
		Description: "Say hello. Returns \"Hello, <text>!\" for the text you pass in.",
	}, hello)

	return server
}

func main() {
	addr := os.Getenv("MCP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := newServer()

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)
	mux.Handle("/mcp/", handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "ok")
	})

	log.Printf("mcp-test-1 %s listening on %s (MCP endpoint: /mcp)", version, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
