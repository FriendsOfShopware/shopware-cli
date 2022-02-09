package shop

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/oauth2"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"shopware-cli/tui"
)

type Client struct {
	url        string
	httpClient *http.Client
	TUI        *tui.TUI
}

type Credentials interface {
	getTokenSource(ctx context.Context, tokenURL string) (oauth2.TokenSource, error)
}

func NewShopClient(ctx context.Context, url string, creds Credentials, httpClient *http.Client) (*Client, error) {
	shopClient := &Client{url, httpClient, nil}
	if err := shopClient.authorize(ctx, creds); err != nil {
		return nil, err
	}

	return shopClient, nil
}

func (sc *Client) authorize(ctx context.Context, creds Credentials) error {
	if sc.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, sc.httpClient)
	}
	tokenSrc, err := creds.getTokenSource(ctx, sc.url+"/api/oauth/token")
	if err != nil {
		return err
	}
	sc.httpClient = oauth2.NewClient(ctx, tokenSrc)
	return nil
}

func (sc *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, sc.url+path, body)
}

func (sc *Client) UploadExtension(ctx context.Context, extensionZip io.Reader) error {
	var buf bytes.Buffer
	parts := multipart.NewWriter(&buf)
	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("Content-Disposition", `form-data; name="file"; filename="extension.zip"`)
	mimeHeader.Set("Content-Type", "application/zip")
	part, err := parts.CreatePart(mimeHeader)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, extensionZip); err != nil {
		return err
	}
	if err := parts.Close(); err != nil {
		return nil
	}
	var body io.Reader = &buf

	if sc.TUI != nil {
		bar := sc.TUI.NewUploadBar(buf.Len(), "app")
		progressReader := progressbar.NewReader(&buf, bar)
		body = &progressReader
	}

	req, err := sc.newRequest(ctx, http.MethodPost, "/api/_action/extension/upload", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", parts.FormDataContentType())
	if resp, err := sc.httpClient.Do(req); err != nil {
		return err
	} else {
		if resp.StatusCode != http.StatusNoContent {
			return errors.New("could not upload extension")
		}
	}
	return nil
}

func (sc *Client) InstallApp(ctx context.Context, appName string) error {
	req, err := sc.newRequest(ctx, http.MethodPost, "/api/_action/extension/install/app/"+appName, nil)
	if err != nil {
		return err
	}
	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("could not install app")
	}
	return nil
}

func (sc *Client) ActivateApp(ctx context.Context, appName string) error {
	req, err := sc.newRequest(ctx, http.MethodPut, "/api/_action/extension/activate/app/"+appName, nil)
	if err != nil {
		return err
	}
	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("could not activate app")
	}
	return nil
}
