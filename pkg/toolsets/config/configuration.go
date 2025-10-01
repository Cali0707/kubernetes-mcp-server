package config

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"k8s.io/utils/ptr"

	"github.com/containers/kubernetes-mcp-server/pkg/api"
	"github.com/containers/kubernetes-mcp-server/pkg/output"
)

func initConfiguration() []api.ServerTool {
	tools := []api.ServerTool{
		{Tool: api.Tool{
			Name:        "contexts_list",
			Description: "List all available contexts from the kubeconfig file. Shows context names for all available contexts",
			InputSchema: &jsonschema.Schema{
				Type: "object",
			},
			Annotations: api.ToolAnnotations{
				Title:           "Contexts: List",
				ReadOnlyHint:    ptr.To(true),
				DestructiveHint: ptr.To(false),
				IdempotentHint:  ptr.To(true),
				OpenWorldHint:   ptr.To(false),
			},
		}, Handler: contextsList},
		{Tool: api.Tool{
			Name:        "configuration_view",
			Description: "Get the current Kubernetes configuration content as a kubeconfig YAML",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"minified": {
						Type: "boolean",
						Description: "Return a minified version of the configuration. " +
							"If set to true, keeps only the current-context and the relevant pieces of the configuration for that context. " +
							"If set to false, all contexts, clusters, auth-infos, and users are returned in the configuration. " +
							"(Optional, default true)",
					},
				},
			},
			Annotations: api.ToolAnnotations{
				Title:           "Configuration: View",
				ReadOnlyHint:    ptr.To(true),
				DestructiveHint: ptr.To(false),
				IdempotentHint:  ptr.To(false),
				OpenWorldHint:   ptr.To(true),
			},
		}, Handler: configurationView},
	}
	return tools
}

func contextsList(params api.ToolHandlerParams) (*api.ToolCallResult, error) {
	contexts, err := params.GetTargets(params.Context)
	if err != nil {
		return api.NewToolCallResult("", fmt.Errorf("failed to list contexts: %v", err)), nil
	}

	if len(contexts) == 0 {
		return api.NewToolCallResult("No contexts found in kubeconfig", nil), nil
	}

	defaultContext := params.GetDefaultTarget()

	result := fmt.Sprintf("Available Kubernetes contexts (%d total, default: %s):\n\n", len(contexts), defaultContext)
	result += "Format: [*] CONTEXT_NAME\n"
	result += " (* indicates the default context used in tools if context is not set)\n\n"
	result += "Contexts:\n---------\n"
	for _, context := range contexts {
		marker := " "
		if context == defaultContext {
			marker = "*"
		}

		result += fmt.Sprintf("%s%s\n", marker, context)
	}
	result += "---------\n\n"

	result += "To use a specific context with any tool, set the 'context' parameter in the tool call arguments"

	return api.NewToolCallResult(result, nil), nil
}

func configurationView(params api.ToolHandlerParams) (*api.ToolCallResult, error) {
	minify := true
	minified := params.GetArguments()["minified"]
	if _, ok := minified.(bool); ok {
		minify = minified.(bool)
	}
	ret, err := params.ConfigurationView(minify)
	if err != nil {
		return api.NewToolCallResult("", fmt.Errorf("failed to get configuration: %v", err)), nil
	}
	configurationYaml, err := output.MarshalYaml(ret)
	if err != nil {
		err = fmt.Errorf("failed to get configuration: %v", err)
	}
	return api.NewToolCallResult(configurationYaml, err), nil
}
