package opensea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type AssetEventsResponse struct {
	Next        string   `json:"next" bson:"next"`
	Previous    string   `json:"previous" bson:"next"`
	AssetEvents []*Event `json:"asset_events" bson:"asset_events"`
}

type Event struct {
	ID                  uint64              `json:"id" bson:"id"`
	Transaction         *Transaction        `json:"transaction" bson:"transaction"`
	PaymentToken        *PaymentToken       `json:"payment_token" bson:"payment_token"`
	Asset               *Asset              `json:"asset" bson:"asset"`
	AssetBundle         *AssetBundle        `json:"asset_bundle" bson:"asset_bundle"`
	WinnerAccount       *Account            `json:"winner_account" bson:"winner_account"`
	FromAccount         *Account            `json:"from_account" bson:"from_account"`
	ToAccount           *Account            `json:"to_account" bson:"to_account"`
	OwnerAccount        *Account            `json:"owner_account" bson:"owner_account"`
	ApprovedAccount     *Account            `json:"approved_account" bson:"approved_account"`
	Seller              *Account            `json:"seller" bson:"seller"`
	DevFeePaymentEvent  *DevFeePaymentEvent `json:"dev_fee_payment_event" bson:"dev_fee_payment_event"`
	CollectionSlug      string              `json:"collection_slug" bson:"collection_slug"`
	CreatedDate         TimeNano            `json:"created_date" bson:"created_date"`
	ModifiedDate        TimeNano            `json:"modified_date" bson:"modified_date"`
	ContractAddress     Address             `json:"contract_address" bson:"contract_address"`
	LogIndex            interface{}         `json:"log_index" bson:"log_index"`
	EventType           EventType           `json:"event_type" bson:"event_type"`
	AuctionType         string              `json:"auction_type" bson:"auction_type"`
	StartingPrice       string              `json:"starting_price" bson:"starting_price"`
	EndingPrice         string              `json:"ending_price" bson:"ending_price"`
	Duration            string              `json:"duration" bson:"duration"`
	MinPrice            Number              `json:"min_price" bson:"min_price"`
	OfferedTo           Number              `json:"offered_to" bson:"offered_to"`
	BidAmount           Number              `json:"bid_amount" bson:"bid_amount"`
	TotalPrice          Number              `json:"total_price" bson:"total_price"`
	CustomEventName     interface{}         `json:"custom_event_name" bson:"custom_event_name"`
	Quantity            string              `json:"quantity" bson:"quantity"`
	PayoutAmount        interface{}         `json:"payout_amount" bson:"payout_amount"`
	EventTimestamp      TimeNano            `json:"event_timestamp" bson:"event_timestamp"`
	Relayer             string              `json:"relayer" bson:"relayer"`
	Collection          uint64              `json:"collection" bson:"collection"`
	PayoutAccount       interface{}         `json:"payout_account" bson:"payout_account"`
	PayoutAssetContract interface{}         `json:"payout_asset_contract" bson:"payout_asset_contract"`
	PayoutCollection    interface{}         `json:"payout_collection" bson:"payout_collection"`
	BuyOrder            uint64              `json:"buy_order" bson:"buy_order"`
	SellOrder           uint64              `json:"sell_order" bson:"sell_order"`
	ListingTime         string              `json:"listing_time" bson:"listing_time"`
	IsPrivate           bool                `json:"is_private" bson:"is_private"`
}

func (e Event) IsBundle() bool {
	return e.AssetBundle != nil
}

type PaymentToken struct {
	Symbol   string      `json:"symbol" bson:"symbol"`
	Address  Address     `json:"address" bson:"address"`
	ImageURL string      `json:"image_url" bson:"image_url"`
	Name     string      `json:"name" bson:"name"`
	Decimals int64       `json:"decimals" bson:"decimals"`
	EthPrice interface{} `json:"eth_price" bson:"eth_price"`
	UsdPrice interface{} `json:"usd_price" bson:"usd_price"`
}

type Transaction struct {
	ID               int64    `json:"id" bson:"id"`
	FromAccount      Account  `json:"from_account" bson:"from_account"`
	ToAccount        Account  `json:"to_account" bson:"to_account"`
	CreatedDate      TimeNano `json:"created_date" bson:"created_date"`
	ModifiedDate     TimeNano `json:"modified_date" bson:"modified_date"`
	TransactionHash  string   `json:"transaction_hash" bson:"transaction_hash"`
	TransactionIndex string   `json:"transaction_index" bson:"transaction_index"`
	BlockNumber      string   `json:"block_number" bson:"block_number"`
	BlockHash        string   `json:"block_hash" bson:"block_hash"`
	Timestamp        string   `json:"timestamp" bson:"timestamp"`
}

// DevFeePaymentEvent is fee transfer event from OpenSea to Dev, It appears to be running in bulk on a regular basis.
type DevFeePaymentEvent struct {
	EventType      string       `json:"event_type" bson:"event_type"`
	EventTimestamp string       `json:"event_timestamp" bson:"event_timestamp"`
	AuctionType    interface{}  `json:"auction_type" bson:"auction_type"`
	TotalPrice     interface{}  `json:"total_price" bson:"total_price"`
	Transaction    Transaction  `json:"transaction" bson:"transaction"`
	PaymentToken   PaymentToken `json:"payment_token" bson:"payment_token"`
}

type EventType string

const (
	EventTypeNone               EventType = ""
	EventTypeCreated            EventType = "created"
	EventTypeSuccessful         EventType = "successful"
	EventTypeCancelled          EventType = "cancelled"
	EventTypeBidEntered         EventType = "bid_entered"
	EventTypeBidWithdrawn       EventType = "bid_withdrawn"
	EventTypeTransfer           EventType = "transfer"
	EventTypeApprove            EventType = "approve"
	EventTypeCompositionCreated EventType = "composition_created"
)

type AuctionType string

const (
	AuctionTypeNone     AuctionType = ""
	AuctionTypeEnglish  AuctionType = "english"
	AuctionTypeDutch    AuctionType = "dutch"
	AuctionTypeMinPrice AuctionType = "min-price"
)

type EventParams struct {
	AssetContractAddress string
	TokenID              int32
	AccountAddress       string
	EventType            EventType
	OnlyOpensea          bool
	AuctionType          AuctionType
	CollectionSlug       string
	Cursor               string
	OccurredBefore       int64
	OccurredAfter        int64
	Limit                int
}

func NewRetrievingEventsParams() *EventParams {
	return &EventParams{
		Limit: 300,
	}
}
func (p EventParams) Encode() string {
	q := url.Values{}

	if p.AssetContractAddress != "" {
		q.Set("asset_contract_address", p.AssetContractAddress)
	}
	if p.TokenID != 0 {
		q.Set("token_id", fmt.Sprintf("%d", p.TokenID))
	}
	if p.CollectionSlug != "" {
		q.Set("collection_slug", p.CollectionSlug)
	}
	if p.AccountAddress != "" {
		q.Set("account_address", p.AccountAddress)
	}
	if p.EventType != EventTypeNone {
		q.Set("event_type", string(p.EventType))
	}
	if p.OnlyOpensea {
		q.Set("only_opensea", "true")
	} else {
		q.Set("only_opensea", "false")
	}
	if p.AuctionType != AuctionTypeNone {
		q.Set("auction_type", string(p.AuctionType))
	}
	//q.Set("limit", fmt.Sprintf("%d", p.Limit))
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	if p.OccurredAfter != 0 {
		q.Set("occurred_after", fmt.Sprintf("%d", p.OccurredAfter))
	}
	if p.OccurredBefore != 0 {
		q.Set("occurred_before", fmt.Sprintf("%d", p.OccurredBefore))
	}
	if p.Limit != 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	return q.Encode()
}

func (o Opensea) RetrievingEvents(params *EventParams) ([]*Event, error) {
	ctx := context.TODO()
	return o.RetrievingEventsWithContext(ctx, params)
}

func (o Opensea) RetrievingEventsWithContext(ctx context.Context, params *EventParams) (events []*Event, err error) {
	if params == nil {
		params = NewRetrievingEventsParams()
	}

	events = []*Event{}
	for true {
		path := "/api/v1/events?" + params.Encode()
		b, err := o.GetPath(ctx, path)
		if err != nil {
			return nil, err
		}
		//fmt.Println(";;;",string(b))
		var eventsResp AssetEventsResponse
		err = json.Unmarshal(b, &eventsResp)
		if err != nil {
			return nil, err
		}
		//fmt.Println("events:", len(eventsResp.AssetEvents), path)

		events = append(events, eventsResp.AssetEvents...)
		if eventsResp.Next == "" {
			break
		}
		params.Cursor = eventsResp.Next
	}

	return
}
