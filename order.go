package opensea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Order struct {
	ID    int64 `json:"id" bson:"id"`
	Asset Asset `json:"asset" bson:"asset"`
	// AssetBundle          interface{}          `json:"asset_bundle" bson:"asset_bundle"`
	CreatedDate       *TimeNano `json:"created_date" bson:"created_date"`
	ClosingDate       *TimeNano `json:"closing_date" bson:"closing_date"`
	ClosingExtendable bool      `json:"closing_extendable" bson:"closing_extendable"`
	ExpirationTime    int64     `json:"expiration_time" bson:"expiration_time"`
	ListingTime       int64     `json:"listing_time" bson:"listing_time"`
	OrderHash         string    `json:"order_hash" bson:"order_hash"`
	Metadata          Metadata  `json:"metadata" bson:"metadata"`
	Exchange          Address   `json:"exchange" bson:"exchange"`
	Maker             Account2  `json:"maker" bson:"maker"`
	Taker             Account2  `json:"taker" bson:"taker"`
	CurrentPrice      Number    `json:"current_price" bson:"current_price"`
	// CurrentBounty        string               `json:"current_bounty" bson:"current_bounty"`
	BountyMultiple     string    `json:"bounty_multiple" bson:"bounty_multiple"`
	MakerRelayerFee    Number    `json:"maker_relayer_fee" bson:"maker_relayer_fee"`
	TakerRelayerFee    Number    `json:"taker_relayer_fee" bson:"taker_relayer_fee"`
	MakerProtocolFee   Number    `json:"maker_protocol_fee" bson:"maker_protocol_fee"`
	TakerProtocolFee   Number    `json:"taker_protocol_fee" bson:"taker_protocol_fee"`
	MakerReferrerFee   Number    `json:"maker_referrer_fee" bson:"maker_referrer_fee"`
	FeeRecipient       Account2  `json:"fee_recipient" bson:"fee_recipient"`
	FeeMethod          FeeMethod `json:"fee_method" bson:"fee_method"`
	Side               Side      `json:"side" bson:"side"` // 0 for buy orders and 1 for sell orders.
	SaleKind           SaleKind  `json:"sale_kind" bson:"sale_kind"`
	Target             Address   `json:"target" bson:"target"`
	HowToCall          HowToCall `json:"how_to_call" bson:"how_to_call"`
	Calldata           string    `json:"calldata" bson:"calldata"`
	ReplacementPattern string    `json:"replacement_pattern" bson:"replacement_pattern"`
	StaticTarget       Address   `json:"static_target" bson:"static_target"`
	StaticExtradata    string    `json:"static_extradata" bson:"static_extradata"`
	PaymentToken       Address   `json:"payment_token" bson:"payment_token"`
	// PaymentTokenContract PaymentTokenContract `json:"payment_token_contract" bson:"payment_token_contract"`
	BasePrice       Number `json:"base_price" bson:"base_price"`
	Extra           Number `json:"extra" bson:"extra"`
	Quantity        string `json:"quantity" bson:"quantity"`
	Salt            Number `json:"salt" bson:"salt"`
	V               *uint8 `json:"v" bson:"v"`
	R               string `json:"r" bson:"r"`
	S               string `json:"s" bson:"s"`
	ApprovedOnChain bool   `json:"approved_on_chain" bson:"approved_on_chain"`
	Cancelled       bool   `json:"cancelled" bson:"cancelled"`
	Finalized       bool   `json:"finalized" bson:"finalized"`
	MarkedInvalid   bool   `json:"marked_invalid" bson:"marked_invalid"`
	PrefixedHash    string `json:"prefixed_hash" bson:"prefixed_hash"`
}

func (o Order) IsPrivate() bool {
	if o.Taker.Address != NullAddress {
		return true
	}
	return false
}

type OrderV2 struct {
	OrderHash        string        `json:"order_hash" bson:"order_hash"`
	CreatedDate      *TimeNano     `json:"created_date" bson:"created_date"`
	ClosingDate      *TimeNano     `json:"closing_date" bson:"closing_date"`
	ExpirationTime   int64         `json:"expiration_time" bson:"expiration_time"`
	ListingTime      int64         `json:"listing_time" bson:"listing_time"`
	Maker            *Account2     `json:"maker" bson:"maker"`
	Taker            *Account2     `json:"taker" bson:"taker"`
	CurrentPrice     Number        `json:"current_price" bson:"current_price"`
	MakerFees        []AccountFee  `json:"maker_fees" bson:"maker_fees"`
	TakerFees        []AccountFee  `json:"taker_fees" bson:"taker_fees"`
	Side             string        `json:"side" bson:"side"`
	OrderType        string        `json:"order_type" bson:"order_type"`
	Cancelled        bool          `json:"cancelled" bson:"cancelled"`
	Finalized        bool          `json:"finalized" bson:"finalized"`
	MarkedInvalid    bool          `json:"marked_invalid" bson:"marked_invalid"`
	RelayId          string        `json:"relay_id" bson:"relay_id"`
	ClientSignature  string        `json:"client_signature" bson:"client_signature"`
	ProtocolAddress  string        `json:"protocol_address" bson:"protocol_address"`
	ProtocolData     *ProtocolData `json:"protocol_data"`
	MakerAssetBundle *AssetBundle  `json:"maker_asset_bundle"`
	TakerAssetBundle *AssetBundle  `json:"taker_asset_bundle"`
}

func (p OrderV2) String() string {
	by, _ := json.Marshal(p)
	return string(by)
}

type Side uint8

const (
	Buy Side = iota
	Sell
)

type SaleKind uint8

const (
	FixedOrMinBit SaleKind = iota // 0 for fixed-price sales or min-bid auctions
	DutchAuctions                 // 1 for declining-price Dutch Auctions
)

type HowToCall uint8

const (
	Call HowToCall = iota
	DelegateCall
)

type FeeMethod uint8

const (
	ProtocolFee FeeMethod = iota
	SplitFee
)

type Metadata struct {
	Asset  MetadataAsset `json:"asset"`
	Schema string        `json:"schema"`
}

type MetadataAsset struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}
type OrderParams struct {
	AssetContractAddress string   `json:"asset_contract_address"`
	PaymentTokenAddress  string   `json:"payment_token_address"`
	Maker                string   `json:"maker"`
	Taker                string   `json:"taker"`
	Owner                string   `json:"owner"`
	IsEnglish            string   `json:"is_english"`
	Bundled              string   `json:"bundled"`
	IncludeBundled       string   `json:"include_bundled"`
	ListedAfter          string   `json:"listed_after"`
	ListedBefore         string   `json:"listed_before"`
	TokenIds             []string `json:"token_ids"`
	Side                 string   `json:"side"`
	SaleKind             string   `json:"sale_kind"`
	Limit                string   `json:"limit"`
	Offset               string   `json:"offset"`
	OrderBy              string   `json:"order_by"`        //created_date,eth_price
	OrderDirection       string   `json:"order_direction"` //asc,desc
	Delay                uint     //in ms
}

type orderResp struct {
	Count  int64    `json:"count"`
	Orders []*Order `json:"orders"`
}
type listingsResp struct {
	Listings        []*Order   `json:"listings"`
	SeaportListings []*OrderV2 `json:"seaport_listings"`
}
type listingsRespV2 struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Orders   []*OrderV2 `json:"orders"`
}

func (o Opensea) GetOrders2(assetContractAddress string, listedAfter int64) ([]*Order, error) {
	ctx := context.TODO()
	return o.GetOrdersWithContext(ctx, assetContractAddress, listedAfter)
}
func (o Opensea) GetActiveListingsV2(assetAddress string, tokenIds []string) ([]*OrderV2, error) {
	path := fmt.Sprintf("/v2/orders/ethereum/seaport/listings?limit=50&asset_contract_address=%s", assetAddress)
	for _, v := range tokenIds {
		path += "&token_ids=" + v
	}
	ctx := context.TODO()
	by, err := o.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	var res listingsRespV2
	err = json.Unmarshal(by, &res)
	if err != nil {
		return nil, err
	}
	return res.Orders, nil
}
func (o Opensea) GetActiveListings(assetAddress string, tokenIds []string, interval time.Duration) ([]*Order, error) {
	var list []*Order
	for _, tokenId := range tokenIds {
		res, err := o.getActiveListings(assetAddress, tokenId)
		if err != nil {
			return nil, err
		}
		list = append(list, res...)
		if interval > 0 {
			time.Sleep(interval)
		}
	}
	return list, nil
}
func (o Opensea) getActiveListings(assetAddress string, tokenId string) ([]*Order, error) {
	path := fmt.Sprintf("/api/v1/asset/%s/%s/listings?limit=%d", assetAddress, tokenId, 50)
	b, err := o.GetPath(context.Background(), path)
	if err != nil {
		return nil, err
	}
	var res listingsResp
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	//for _, v := range res.Listings {
	//	//reset some info
	//	v.Asset = Asset{TokenID: v.Metadata.Asset.ID, AssetContract: &AssetContract{Address: Address(v.Metadata.Asset.Address)}}
	//}
	return res.Listings, nil
}
func (o Opensea) GetOrders(params OrderParams, findAll bool) ([]*Order, error) {
	if !findAll {
		return o.getOrders(params)
	}
	offset := 0
	limit := 50
	var orders []*Order
	for {
		params.Offset = fmt.Sprintf("%d", offset)
		params.Limit = fmt.Sprintf("%d", limit)
		ords, err := o.getOrders(params)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ords...)
		if len(ords) < limit {
			break
		}
		if params.Delay > 0 {
			time.Sleep(time.Duration(params.Delay) * time.Millisecond)
		}
		offset += limit
	}
	return orders, nil
}
func (o Opensea) getOrders(params OrderParams) ([]*Order, error) {
	q := url.Values{}
	if params.Offset == "" {
		q.Set("offset", "0")
	} else {
		q.Set("offset", params.Offset)
	}
	if params.Limit == "" {
		q.Set("limit", "50")
	} else {
		q.Set("limit", params.Limit)
	}
	if params.OrderBy == "" {
		q.Set("order_by", "created_date")
	} else {
		q.Set("order_by", params.OrderBy)
	}
	if params.OrderDirection == "" {
		q.Set("order_direction", "desc")
	} else {
		q.Set("order_by", params.OrderDirection)
	}
	if params.AssetContractAddress != "" {
		q.Set("asset_contract_address", params.AssetContractAddress)
	}
	if params.PaymentTokenAddress != "" {
		q.Set("payment_token_address", params.PaymentTokenAddress)
	}
	if params.Maker != "" {
		q.Set("maker", params.Maker)
	}
	if params.Taker != "" {
		q.Set("taker", params.Taker)
	}
	if params.Owner != "" {
		q.Set("owner", params.Owner)
	}
	if params.IsEnglish != "" {
		q.Set("is_english", params.IsEnglish)
	}
	if params.Bundled != "" {
		q.Set("bundled", params.Bundled)
	}
	if params.IncludeBundled != "" {
		q.Set("include_bundled", params.IncludeBundled)
	}
	if params.ListedAfter != "" {
		q.Set("listed_after", params.ListedAfter)
	}
	if params.ListedBefore != "" {
		q.Set("listed_before", params.ListedBefore)
	}
	if params.Side != "" {
		q.Set("side", params.Side)
	}
	if params.SaleKind != "" {
		q.Set("sale_kind", params.SaleKind)
	}
	path := "/wyvern/v1/orders?" + q.Encode()
	if params.TokenIds != nil {
		for _, v := range params.TokenIds {
			//q.Set("token_ids", v)
			path += fmt.Sprintf("&token_ids=%s", v)
		}
	}
	b, err := o.GetPath(context.Background(), path)
	if err != nil {
		return nil, err
	}
	out := &orderResp{}
	err = json.Unmarshal(b, out)
	if err != nil {
		return nil, err
	}
	return out.Orders, nil
}
func (o Opensea) GetOrdersWithContext(ctx context.Context, assetContractAddress string, listedAfter int64) (orders []*Order, err error) {
	offset := 0
	limit := 100

	q := url.Values{}
	q.Set("asset_contract_address", assetContractAddress)
	q.Set("listed_after", fmt.Sprintf("%d", listedAfter))
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("order_by", "created_date")
	q.Set("order_direction", "asc")

	orders = []*Order{}

	for true {
		q.Set("offset", fmt.Sprintf("%d", offset))
		path := "/wyvern/v1/orders?" + q.Encode()
		b, err := o.GetPath(ctx, path)
		if err != nil {
			return nil, err
		}

		out := &struct {
			Count  int64    `json:"count"`
			Orders []*Order `json:"orders"`
		}{}

		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}
		orders = append(orders, out.Orders...)

		if len(out.Orders) < limit {
			break
		}
		offset += limit
	}

	return
}
