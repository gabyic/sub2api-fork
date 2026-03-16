package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const openaiPlatformEmbeddingsURL = "https://api.openai.com/v1/embeddings"

// ForwardAsEmbeddings forwards an OpenAI-compatible embeddings request to the
// selected upstream account and writes the upstream response back to the client.
func (s *OpenAIGatewayService) ForwardAsEmbeddings(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
) (*OpenAIForwardResult, error) {
	originalModel := strings.TrimSpace(gjson.GetBytes(body, "model").String())
	if originalModel == "" {
		return nil, fmt.Errorf("model is required")
	}

	mappedModel := resolveOpenAIForwardModel(account, originalModel, "")
	if mappedModel != "" && mappedModel != originalModel {
		updatedBody, err := sjson.SetBytes(body, "model", mappedModel)
		if err == nil {
			body = updatedBody
		}
	}

	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	upstreamReq, err := s.buildUpstreamEmbeddingsRequest(ctx, c, account, body, token)
	if err != nil {
		return nil, fmt.Errorf("build upstream request: %w", err)
	}

	proxyURL := ""
	if account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		safeErr := sanitizeUpstreamErrorMessage(err.Error())
		setOpsUpstreamError(c, 0, safeErr, "")
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: 0,
			Kind:               "request_error",
			Message:            safeErr,
		})
		return nil, fmt.Errorf("upstream request failed: %s", safeErr)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(respBody))

		upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
		upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
		if s.shouldFailoverOpenAIUpstreamResponse(resp.StatusCode, upstreamMsg, respBody) {
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  resp.Header.Get("x-request-id"),
				Kind:               "failover",
				Message:            upstreamMsg,
			})
			if s.rateLimitService != nil {
				s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
			}
			return nil, &UpstreamFailoverError{
				StatusCode:             resp.StatusCode,
				ResponseBody:           respBody,
				RetryableOnSameAccount: account.IsPoolMode() && (isPoolModeRetryableStatus(resp.StatusCode) || isOpenAITransientProcessingError(resp.StatusCode, upstreamMsg, respBody)),
			}
		}
			err := s.handleErrorResponsePassthrough(ctx, resp, c, account, body)
			return nil, err
		}

	if s.rateLimitService != nil {
		s.rateLimitService.UpdateSessionWindow(ctx, account, resp.Header)
	}
	usage, err := s.handleNonStreamingResponsePassthrough(ctx, resp, c)
	if err != nil {
		return nil, err
	}

	result := &OpenAIForwardResult{
		Model:           originalModel,
		BillingModel:    mappedModel,
		Usage:           *usage,
		Duration:        0,
		Stream:          false,
		RequestID:       resp.Header.Get("x-request-id"),
		ResponseHeaders: resp.Header.Clone(),
	}
	if strings.TrimSpace(mappedModel) == "" {
		result.BillingModel = originalModel
	}
	return result, nil
}

func (s *OpenAIGatewayService) buildUpstreamEmbeddingsRequest(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	body []byte,
	token string,
) (*http.Request, error) {
	targetURL := openaiPlatformEmbeddingsURL
	baseURL := account.GetOpenAIBaseURL()
	if baseURL != "" {
		validatedURL, err := s.validateUpstreamBaseURL(baseURL)
		if err != nil {
			return nil, err
		}
		targetURL = buildOpenAIEmbeddingsURL(validatedURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", "Bearer "+token)
	if c != nil && c.Request != nil {
		for key, values := range c.Request.Header {
			lower := strings.ToLower(strings.TrimSpace(key))
			if !openaiAllowedHeaders[lower] && lower != "accept" {
				continue
			}
			for _, v := range values {
				req.Header.Add(key, v)
			}
		}
	}
	req.Header.Del("x-api-key")
	req.Header.Del("x-goog-api-key")
	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/json")
	}

	customUA := account.GetOpenAIUserAgent()
	if customUA != "" {
		req.Header.Set("user-agent", customUA)
	}
	return req, nil
}

func buildOpenAIEmbeddingsURL(base string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	if strings.HasSuffix(normalized, "/embeddings") {
		return normalized
	}
	if strings.HasSuffix(normalized, "/v1") {
		return normalized + "/embeddings"
	}
	return normalized + "/v1/embeddings"
}
