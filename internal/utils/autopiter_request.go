package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	errors2 "parser/internal/errors"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/google/uuid"
)

func MakeAutopiterGetRequest(ctx context.Context, link string, proxy string) (http.Response, errors2.ServiceError) {
	client, svcErr := autopiterHTTPClient(proxy)
	if svcErr != nil {
		return http.Response{}, svcErr
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return http.Response{}, errors2.UnableToCreateReq(err)
	}

	setAutopiterHeaders(req)
	return makeRequest(ctx, client, req)
}

func MakeAutopiterSearchRequest(ctx context.Context, link string, proxy string, body any) (http.Response, errors2.ServiceError) {
	payload, err := json.Marshal(body)
	if err != nil {
		return http.Response{}, errors2.UnableToUnmarshall(err)
	}

	client, svcErr := autopiterHTTPClient(proxy)
	if svcErr != nil {
		return http.Response{}, svcErr
	}

	req, err := http.NewRequest(http.MethodPost, link, bytes.NewReader(payload))
	if err != nil {
		return http.Response{}, errors2.UnableToCreateReq(err)
	}

	setAutopiterHeaders(req)
	req.Header.Set("content-type", "application/json")
	return makeRequest(ctx, client, req)
}

func setAutopiterHeaders(req *http.Request) {
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("origin", "https://autopiter.ru")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://autopiter.ru/")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", uarand.GetRandom())
	req.Header.Set("x-ap-request-id", uuid.NewString())
}

func autopiterHTTPClient(proxy string) (http.Client, errors2.ServiceError) {
	if proxy == "" {
		return http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    0,
				MaxConnsPerHost: 0,
				IdleConnTimeout: 3 * time.Second,
			},
		}, nil
	}

	proxy = strings.Replace(proxy, ",", ":", -1)
	parts := strings.Split(proxy, ":")
	if len(parts) != 4 {
		fmt.Println("Неверный формат строки прокси")
		os.Exit(1)
	}

	proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return http.Client{}, errors2.WrongProxyUrl(err)
	}

	return http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    0,
			MaxConnsPerHost: 0,
			Proxy:           http.ProxyURL(parsedURL),
			IdleConnTimeout: 3 * time.Second,
		},
	}, nil
}
