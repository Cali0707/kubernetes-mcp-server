package mcp

import (
	"context"
	"fmt"

	"github.com/containers/kubernetes-mcp-server/pkg/api"
	"github.com/containers/kubernetes-mcp-server/pkg/mcplog"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type sessionElicitor struct{}

var _ api.Elicitor = &sessionElicitor{}

func (s *sessionElicitor) Elicit(ctx context.Context, message string, requestedSchema *jsonschema.Schema) (*api.ElicitResult, error) {
	session, ok := ctx.Value(mcplog.MCPSessionContextKey).(*mcp.ServerSession)
	if !ok || session == nil {
		return nil, fmt.Errorf("no MCP session found in context")
	}

	result, err := session.Elicit(ctx, &mcp.ElicitParams{Message: message, RequestedSchema: requestedSchema})
	// TODO: check if the error is because the client does not support elicitation, possibly return default value
	if err != nil {
		return nil, err
	}

	return &api.ElicitResult{Action: result.Action, Content: result.Content}, nil
}
