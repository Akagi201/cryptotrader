# binance api

## Supported API list

| Func           	| API endpoint                      	| Meanings                                            	|
|----------------	|-----------------------------------	|-----------------------------------------------------	|
| GetTicker      	| GET /api/v1/ticker/24hr           	| 24hr ticker price change statistics                 	|
| Ping           	| GET /api/v1/ping                  	| Test connectivity to the Rest API                   	|
| GetTime        	| GET /api/v1/time                  	| Check server time                                   	|
| GetDepth       	| GET /api/v1/depth                 	| Order book                                          	|
| GetTickers     	| GET /api/v1/ticker/allPrices      	| Symbols price ticker                                	|
| GetTrades      	| GET /api/v1/aggTrades             	| Compressed/Aggregate trades list                    	|
| GetRecords     	| GET /api/v1/ticker/allPrices      	| Kline/candlesticks                                  	|
| GetBookTickers 	| GET /api/v1/ticker/allBookTickers 	| Symbols order book ticker                           	|
| GetAccount     	| GET /api/v3/account               	| Account information                                 	|
| Trade          	| POST /api/v3/order                	| Send in a new order                                 	|
| GetOrder       	| GET /api/v3/order                 	| Check an order's status                             	|
| CancelOrder    	| DELETE /api/v3/order              	| Cancel an active order                              	|
| GetOrders      	| GET /api/v3/openOrders            	| Get all open orders on a symbol                     	|
| GetAllOrders   	| GET /api/v3/allOrders             	| Get all account orders; active, canceled, or filled 	|
| GetMyTrades    	| GET /api/v3/myTrades              	| Get trades for a specific account and symbol        	|
