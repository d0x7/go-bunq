package bunq

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

// CreateContext registers a new API key and device, to create a session.
// The hereby created API context is then saved to the specified contextFile, and may be loaded at a later time using LoadContext.
// PermittedIps may either use bunq.WildcardIP or bunq.CurrentIP, or a custom list of specific IP addresses or ranges.
func CreateContext(ctx context.Context, url, apiKey, deviceDescription string, permittedIps []string, contextFile string) (*Client, error) {
	key, err := CreateNewKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "creating new key pair")
	}

	client := NewClient(ctx, url, key, apiKey, deviceDescription, permittedIps)

	if err := client.Init(); err != nil {
		return nil, errors.Wrap(err, "initializing bunq client")
	}

	err = SaveContext(client, contextFile)
	if err != nil {
		return nil, errors.Wrap(err, "saving client context")
	}

	return client, nil
}

// SaveContext saves the client's API context to the specified file.
func SaveContext(client *Client, contextFile string) error {
	clientContext, err := client.ExportClientContext()
	if err != nil {
		return errors.Wrap(err, "exporting 	client context")
	}

	indent, err := json.MarshalIndent(clientContext, "", "    ")
	if err != nil {
		return errors.Wrap(err, "marshaling client context")
	}

	if os.WriteFile(contextFile, indent, 0600) != nil {
		return errors.Wrap(err, "writing client context")
	}

	return nil
}

// LoadContext loads a previously created API context from the specified file and initializes a new client from it.
func LoadContext(ctx context.Context, file string) (*Client, error) {
	var clientContext ClientContext

	readFile, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "reading client context")
	}

	err = json.Unmarshal(readFile, &clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling client context")
	}

	client, err := NewClientFromContext(ctx, &clientContext)
	if err != nil {
		return nil, errors.Wrap(err, "creating client from context")
	}

	if client.Init() != nil {
		return nil, errors.Wrap(err, "initializing bunq client")
	}

	return client, nil
}
