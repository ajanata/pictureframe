package req

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetJson(ctx context.Context, url string, body any) error {
	reader, err := GetRaw(ctx, url)

	decode := json.NewDecoder(reader)
	err = decode.Decode(body)
	_ = reader.Close()
	if err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

func GetRaw(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	return resp.Body, nil
}
