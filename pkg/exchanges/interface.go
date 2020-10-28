package exchanges

import (
	"context"
	"time"

	"github.com/Akagi201/cryptotrader/pkg/exchanges/models"
	"github.com/Akagi201/cryptotrader/pkg/parameters"
)

// IExchange is a unified exchange interface
type IExchange interface {
	FetchTickers(ctx context.Context, symbols []string, params ...parameters.Params) (map[string]models.Ticker, error)
	FetchTicker(ctx context.Context, symbol string, params ...parameters.Params) (models.Ticker, error)
	FetchOHLCV(ctx context.Context, symbol, duration string, since time.Time, limit int, params ...parameters.Params) ([]*models.OHLCV, error)
	FetchOrderBook(ctx context.Context, symbol string, limit int, params ...parameters.Params) (*models.OrderBook, error)
	FetchL2OrderBook(ctx context.Context, symbol string, limit int, params ...parameters.Params) (*models.OrderBook, error)
	FetchTrades(ctx context.Context, symbol string, since time.Time, params ...parameters.Params) ([]*models.Trade, error)
	FetchOrder(ctx context.Context, id string, symbol string, params ...parameters.Params) (*models.Order, error)
	FetchOrders(ctx context.Context, symbol string, since time.Time, limit int, params ...parameters.Params) ([]*models.Order, error)
	FetchOpenOrders(ctx context.Context, symbol string, since time.Time, limit int, params ...parameters.Params) ([]*models.Order, error)
	FetchClosedOrders(ctx context.Context, symbol string, since time.Time, limit int, params ...parameters.Params) ([]*models.Order, error)
	FetchMyTrades(ctx context.Context, symbol string, since time.Time, limit int, params ...parameters.Params) ([]*models.Trade, error)
	FetchBalance(ctx context.Context, params ...parameters.Params) (*models.Balances, error)
	FetchCurrencies(ctx context.Context, params ...parameters.Params) (map[string]*models.Currency, error)
	FetchMarkets(ctx context.Context, params ...parameters.Params) ([]*models.Market, error)

	CreateOrder(ctx context.Context, symbol, typ, side string, amount float64, price float64, params ...parameters.Params) (*models.Order, error)
	CancelOrder(ctx context.Context, id string, symbol string, params ...parameters.Params) error
	CreateLimitBuyOrder(ctx context.Context, symbol string, amount float64, price float64, params ...parameters.Params) (*models.Order, error)
	CreateLimitSellOrder(ctx context.Context, symbol string, amount float64, price float64, params ...parameters.Params) (*models.Order, error)
	CreateMarketBuyOrder(ctx context.Context, symbol string, amount float64, params ...parameters.Params) (*models.Order, error)
	CreateMarketSellOrder(ctx context.Context, symbol string, amount float64, params ...parameters.Params) (*models.Order, error)
	GetInfo(ctx context.Context) *models.ExchangeInfo

	GetMarkets(ctx context.Context) map[string]*models.Market
	GetMarketsByID(ctx context.Context, id string) map[string]*models.Market
	GetCurrencies(ctx context.Context) map[string]*models.Currency
	GetCurrenciesByID(ctx context.Context, id string) map[string]*models.Currency
	GetMarket(symbol string) (*models.Market, error)
}
