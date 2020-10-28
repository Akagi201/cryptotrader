package models

import (
	"time"

	"go.uber.org/ratelimit"
)

// Ticker struct
type Ticker struct {
	Symbol      string    `json:"symbol"`
	Ask         float64   `json:"ask"`
	Bid         float64   `json:"bid"`
	Last        float64   `json:"last"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	BaseVolume  float64   `json:"base_volume"`
	QuoteVolume float64   `json:"quote_volume"`
	Change      float64   `json:"change"`
	Timestamp   time.Time `json:"timestamp"`
}

// OHLCV open, high, low, close, volume
type OHLCV struct {
	O         float64   `json:"o"`
	H         float64   `json:"h"`
	L         float64   `json:"l"`
	C         float64   `json:"c"`
	V         float64   `json:"v"`
	Timestamp time.Time `json:"timestamp"`
}

type OrderBook struct {
	Symbol    string
	Asks      []BookEntry `json:"asks"`
	Bids      []BookEntry `json:"bids"`
	Timestamp time.Time   `json:"timestamp"`
	Nonce     int64       `json:"nonce"`
}

type BookEntry struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type Trade struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Amount    float64   `json:"amount"`
	Price     float64   `json:"price"`
	Order     string    `json:"order"`
	Type      string    `json:"type"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}

type Order struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Side      string    `json:"side"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Cost      float64   `json:"cost"`
	Filled    float64   `json:"filled"`
	Remaining float64   `json:"remaining"`
	Fee       float64   `json:"fee"`
	Timestamp time.Time `json:"timestamp"`
}

type Balance struct {
	Free  float64 `json:"free"`
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

type Balances struct {
	Currencies map[string]Balance `json:"currencies"`
	Free       map[string]float64 `json:"free"`
	Used       map[string]float64 `json:"used"`
	Total      map[string]float64 `json:"total"`
}

// Currency struct
type Currency struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Precision int    `json:"precision"`
}

type Market struct {
	ID        string    `json:"id"`     // exchange specific
	Symbol    string    `json:"symbol"` // cryptotrader unified
	Base      string    `json:"base"`
	Quote     string    `json:"quote"`
	Precision Precision `json:"precision"`
}

type Precision struct {
	Amount int `json:"amount"`
	Base   int `json:"base"`
	Price  int `json:"price"`
	Cost   int `json:"cost"`
}

// HasDescription for exchange functionality
type HasDescription struct {
	CancelAllOrders      bool `json:"cancel_all_orders"`
	CancelOrder          bool `json:"cancel_order"`
	CancelOrders         bool `json:"cancel_orders"`
	CORS                 bool `json:"cors"`
	CreateDepositAddress bool `json:"create_deposit_address"`
	CreateLimitOrder     bool `json:"create_limit_order"`
	CreateMarketOrder    bool `json:"create_market_order"`
	CreateOrder          bool `json:"create_order"`
	Deposit              bool `json:"deposit"`
	EditOrder            bool `json:"edit_order"`
	FetchBalance         bool `json:"fetch_balance"`
	FetchBidsAsks        bool `json:"fetch_bids_asks"`
	FetchClosedOrders    bool `json:"fetch_closed_orders"`
	FetchCurrencies      bool `json:"fetch_currencies"`
	FetchDepositAddress  bool `json:"fetch_deposit_address"`
	FetchDeposits        bool `json:"fetch_deposits"`
	FetchFundingFees     bool `json:"fetch_funding_fees"`
	FetchL2OrderBook     bool `json:"fetch_l2_orderbook"`
	FetchLedger          bool `json:"fetch_ledger"`
	FetchMarkets         bool `json:"fetch_markets"`
	FetchMyTrades        bool `json:"fetch_mytrades"`
	FetchOHLCV           bool `json:"fetch_ohlcv"`
	FetchOpenOrders      bool `json:"fetch_open_orders"`
	FetchOrder           bool `json:"fetch_order"`
	FetchOrderBook       bool `json:"fetch_orderbook"`
	FetchOrderBooks      bool `json:"fetch_orderbooks"`
	FetchOrders          bool `json:"fetch_orders"`
	FetchTicker          bool `json:"fetch_ticker"`
	FetchTickers         bool `json:"fetch_tickers"`
	FetchTrades          bool `json:"fetch_trades"`
	FetchTradingFee      bool `json:"fetch_trading_fee"`
	FetchTradingFees     bool `json:"fetch_trading_fees"`
	FetchTradingLimits   bool `json:"fetch_trading_limits"`
	FetchTransactions    bool `json:"fetch_transactions"`
	FetchWithdrawals     bool `json:"fetch_withdrawals"`
	Withdraw             bool `json:"withdraw"`
}

type ExchangeInfo struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	RateLimit  ratelimit.Limiter `json:"rate_limit"`
	Has        HasDescription    `json:"has"`
	UserAgents map[string]string `json:"userAgents"`
	Market     Market            `json:"market"`
	BaseURL    string            `json:"base_url"`
	APIKey     string            `json:"api_key"`
	APISecret  string            `json:"api_secret"`
}
