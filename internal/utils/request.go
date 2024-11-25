package utils

import (
	"context"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"os"
	errors2 "parser/internal/errors"
	"strings"
	"time"
)

func MakeJsonRequestVisitor(ctx context.Context, link string, proxy string, visitorId string) (http.Response, errors2.ServiceError) {
	// Заменяем центральное двоеточие на @
	parts1 := strings.Replace(proxy, ":", ",", -1)
	fmt.Println(proxy)
	parts := strings.Split(parts1, ",")
	if len(parts) != 4 {
		fmt.Println("Неверный формат строки прокси")
		os.Exit(1)
	}
	// 148.129.126.60:9375:a704Ddn:gY21pM7
	// 185.181.244.152:3000,tYcQJm91,Etw7MVvN

	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])

	Url, err := url.Parse(proxyUrl)

	if err != nil {

		return http.Response{}, errors2.WrongProxyUrl(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    0,
			MaxConnsPerHost: 0,
			Proxy:           http.ProxyURL(Url),
			IdleConnTimeout: 3 * time.Second},
	}
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	//req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("pragma", "no-cache")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	//req.Header.Add("sec-fetch-user", "?0")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", uarand.GetRandom())
	req.AddCookie(&http.Cookie{
		Name:  "visitor_id",
		Value: visitorId,
	})
	req.AddCookie(&http.Cookie{
		Name:  "siteversion",
		Value: "2",
	})

	//req.AddCookie(&http.Cookie{Name: "NSC_xxx.fnfy.sv", Value: uuid.NewString()})

	if err != nil {
		return http.Response{}, errors2.UnableToCreateReq(err)
	}
	resp, e := makeRequest(ctx, client, req)
	if e != nil {
		return http.Response{}, e
	}

	return resp, nil

}

func MakeJsonRequest(ctx context.Context, link string, proxy string) (http.Response, errors2.ServiceError) {
	// Заменяем центральное двоеточие на @
	parts := strings.Split(proxy, ":")
	if len(parts) != 4 {
		fmt.Println("Неверный формат строки прокси")
		os.Exit(1)
	}

	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])

	Url, err := url.Parse(proxyUrl)

	if err != nil {

		return http.Response{}, errors2.WrongProxyUrl(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    0,
			MaxConnsPerHost: 0,
			Proxy:           http.ProxyURL(Url),
			IdleConnTimeout: 3 * time.Second},
	}
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	//req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("pragma", "no-cache")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	//req.Header.Add("sec-fetch-user", "?0")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", uarand.GetRandom())
	req.AddCookie(&http.Cookie{
		Name:  "visitor_id",
		Value: uuid.New().String(),
	})
	req.AddCookie(&http.Cookie{
		Name:  "siteversion",
		Value: "2",
	})

	//req.AddCookie(&http.Cookie{Name: "NSC_xxx.fnfy.sv", Value: uuid.NewString()})

	if err != nil {
		return http.Response{}, errors2.UnableToCreateReq(err)
	}
	resp, e := makeRequest(ctx, client, req)
	if e != nil {
		return http.Response{}, e
	}

	return resp, nil

}

func MakeJsonRequestPost(ctx context.Context, link string, proxy string) (http.Response, errors2.ServiceError) {
	// Заменяем центральное двоеточие на @
	parts := strings.Split(proxy, ":")
	if len(parts) != 4 {
		fmt.Println("Неверный формат строки прокси")
		os.Exit(1)
	}

	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", parts[2], parts[3], parts[0], parts[1])

	Url, err := url.Parse(proxyUrl)

	if err != nil {

		return http.Response{}, errors2.WrongProxyUrl(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    0,
			MaxConnsPerHost: 0,
			Proxy:           http.ProxyURL(Url),
			IdleConnTimeout: 3 * time.Second},
	}
	req, err := http.NewRequest("POST", link, nil)
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	//req.Header.Add("cache-control", "no-cache")
	//req.Header.Add("pragma", "no-cache")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	//req.Header.Add("sec-fetch-user", "?0")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", uarand.GetRandom())
	//req.AddCookie(&http.Cookie{
	//	Name:  "visitor_id",
	//	Value: uuid.New().String(),
	//})
	//req.AddCookie(&http.Cookie{
	//	Name:  "siteversion",
	//	Value: "2",
	//})
	//
	//req.AddCookie(&http.Cookie{Name: "NSC_xxx.fnfy.sv", Value: uuid.NewString()})

	if err != nil {
		return http.Response{}, errors2.UnableToCreateReq(err)
	}
	resp, e := makeRequest(ctx, client, req)
	if e != nil {
		return http.Response{}, e
	}

	return resp, nil

}

func makeRequest(ctx context.Context, client http.Client, req *http.Request) (http.Response, errors2.ServiceError) {
	//resp, err := http.DefaultClient.Do(req)
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return http.Response{}, errors2.WrongProxyUrl(err)
	}
	return *resp, nil
}
