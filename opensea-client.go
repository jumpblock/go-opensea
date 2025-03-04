package opensea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"time"
)

var (
	mainnetAPI = "https://api.opensea.io"
	rinkebyAPI = "https://rinkeby-api.opensea.io"
)

type Opensea struct {
	API        string
	APIKey     string
	httpClient *http.Client
	proxy      string
}

type errorResponse struct {
	Success bool   `json:"success" bson:"success"`
	Msg     string `json:"msg" bson:"msg"`
}

func (e errorResponse) Error() string {
	return e.Msg
}

func NewOpensea(apiKey string) (*Opensea, error) {
	o := &Opensea{
		API:        mainnetAPI,
		APIKey:     apiKey,
		httpClient: defaultHttpClient(),
	}
	return o, nil
}
func NewOpenseaWithProxy(apiKey, proxy string) (*Opensea, error) {
	o := &Opensea{
		API:        mainnetAPI,
		APIKey:     apiKey,
		httpClient: defaultHttpClient(),
		proxy:      proxy,
	}
	return o, nil
}

func NewOpenseaRinkeby(apiKey string) (*Opensea, error) {
	o := &Opensea{
		API:        rinkebyAPI,
		APIKey:     apiKey,
		httpClient: defaultHttpClient(),
	}
	return o, nil
}

// TODO
//func (o Opensea) GetAssets(params GetAssetsParams) (*AssetResponse, error) {
//	ctx := context.TODO()
//	return o.GetAssetsWithContext(ctx, params)
//}

// TODO
//func (o Opensea) GetAssetsWithContext(ctx context.Context, params GetAssetsParams) (*AssetResponse, error) {
//	path := fmt.Sprintf("/api/v1/assets")
//	b, err := o.GetPath(ctx, path)
//	if err != nil {
//		return nil, err
//	}
//	ret := new(AssetResponse)
//	return ret, json.Unmarshal(b, ret)
//}

func (o Opensea) GetCollections(offset, limit int) ([]CollectionSingle, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/api/v1/collections?offset=%d&limit=%d", offset, limit)
	b, err := o.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	resp := new(collectionsResp)
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Collections, nil
}
func (o Opensea) GetSingleCollection(slug string) (CollectionSingle, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/api/v1/collection/%s", slug)
	b, err := o.GetPath(ctx, path)
	if err != nil {
		return CollectionSingle{}, err
	}
	resp := new(CollectionSingleResponse)
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return CollectionSingle{}, err
	}
	return resp.Collection, nil
}

func (o Opensea) GetSingleAsset(assetContractAddress string, tokenID *big.Int) (*Asset, error) {
	ctx := context.TODO()
	return o.GetSingleAssetWithContext(ctx, assetContractAddress, tokenID)
}

func (o Opensea) GetSingleAssetWithContext(ctx context.Context, assetContractAddress string, tokenID *big.Int) (*Asset, error) {
	path := fmt.Sprintf("/api/v1/asset/%s/%s", assetContractAddress, tokenID.String())
	b, err := o.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	ret := new(Asset)
	return ret, json.Unmarshal(b, ret)
}
func (o Opensea) GetAssetDetail(assetContractAddress string) (*Asset, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/api/v1/assets?asset_contract_address=%s&limit=1", assetContractAddress)
	b, err := o.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	var res assetsResp
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	if res.Assets != nil && len(res.Assets) > 0 {
		return &res.Assets[0], nil
	} else {
		return nil, fmt.Errorf("no asset return")
	}
}

func (o Opensea) GetPath(ctx context.Context, path string) ([]byte, error) {
	return o.getURL(ctx, o.API+path)
}
func (o Opensea) PostPath(ctx context.Context, path string, data []byte) ([]byte, error) {
	client := o.httpClient
	target := o.API + path
	if o.proxy != "" {
		target = o.proxy
	}
	req, err := http.NewRequestWithContext(ctx, "POST", target, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-KEY", o.APIKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if o.proxy != "" {
		req.Header.Add("__ddd__", o.API+path)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		e := new(errorResponse)
		err = json.Unmarshal(body, e)
		if err != nil {
			return nil, fmt.Errorf("Backend returns status %d msg: %s", resp.StatusCode, string(body))
		}
		if !e.Success {
			e.Msg = resp.Status
			return nil, e
		}

		return nil, fmt.Errorf("Backend returns status %d msg: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (o Opensea) getURL(ctx context.Context, url string) ([]byte, error) {
	client := o.httpClient
	target := url
	if o.proxy != "" {
		target = o.proxy
	}
	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	req.Header.Add("X-API-KEY", o.APIKey)
	req.Header.Add("Accept", "application/json")
	if o.proxy != "" {
		req.Header.Add("__ddd__", url)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		e := new(errorResponse)
		err = json.Unmarshal(body, e)
		if err != nil {
			return nil, fmt.Errorf("Backend returns status %d msg: %s", resp.StatusCode, string(body))
		}
		if !e.Success {
			e.Msg = resp.Status
			return nil, e
		}

		return nil, fmt.Errorf("Backend returns status %d msg: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (o Opensea) SetHttpClient(httpClient *http.Client) {
	o.httpClient = httpClient
}

func defaultHttpClient() *http.Client {
	client := new(http.Client)
	var transport http.RoundTripper = &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		DisableKeepAlives:  false,
		DisableCompression: false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 300 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client.Transport = transport
	return client
}
