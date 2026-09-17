package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHello(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"with text", "world", "Hello, world!"},
		{"empty text", "", "Hello!"},
		{"whitespace only", "   ", "Hello!"},
		{"trims surrounding space", "  gateway  ", "Hello, gateway!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _, err := hello(context.Background(), nil, &HelloParams{Text: tt.text})
			if err != nil {
				t.Fatalf("hello() error = %v", err)
			}
			got := textOf(t, res)
			if got != tt.want {
				t.Errorf("hello(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

// TestHelloOverMCP exercises the tool end to end through a real MCP session,
// which is what the gateway will do.
func TestHelloOverMCP(t *testing.T) {
	ctx := context.Background()

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverSession, err := newServer().Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer serverSession.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "hello" {
		t.Fatalf("unexpected tool list: %+v", tools.Tools)
	}

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "hello",
		Arguments: map[string]any{"text": "gateway"},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}
	if got, want := textOf(t, res), "Hello, gateway!"; got != want {
		t.Errorf("hello over MCP = %q, want %q", got, want)
	}
}

func textOf(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if res == nil || len(res.Content) != 1 {
		t.Fatalf("expected exactly one content item, got %+v", res)
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	return tc.Text
}
