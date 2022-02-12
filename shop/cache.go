package shop

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (c Client) ClearCache(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodDelete, "/api/_action/cache", nil)

	if err != nil {
		return errors.Wrap(err, "ClearCache")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("ClearCache: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}
