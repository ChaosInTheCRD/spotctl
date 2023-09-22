package internal

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	//"google.golang.org/api/iterator"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"math/rand"
)

func Generate(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func UpdateSecret(projectID, secretID, value string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err.Error())
	}
	defer client.Close()

	greq := &secretmanagerpb.GetSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID),
	}

	resp, err := client.GetSecretVersion(ctx, greq)
	if err != nil {
		return fmt.Errorf("failed to get latest secret version: %v", err.Error())
	}

	dreq := &secretmanagerpb.DestroySecretVersionRequest{
		Name: resp.Name,
	}

	_, err = client.DestroySecretVersion(ctx, dreq)
	if err != nil {
		err = fmt.Errorf("failed to delete secret version %s: %v", resp.Name, err.Error())
		fmt.Println(err)
	}

	// Build the request.
	areq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%s/secrets/%s", projectID, secretID),
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(value),
		},
	}

	// Call the API.
	_, err = client.AddSecretVersion(ctx, areq)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %v", err)
	}

	//fmt.Printf("added secret version: %s\n", result.Name)
	return nil

}

// GetLatestSecret returns the payload for the secret under the `/latest` tag.
func GetLatestSecret(projectID, secretID string) (string, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	fmt.Printf("retrieved payload for: %s\n", result.Name)
	return string(result.Payload.Data), nil
}
