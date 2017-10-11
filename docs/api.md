# 通用交易所 API

## 协议
* HTTP Rest
* WebSocket
* Fix

## Request

```
{
  "access_key": "XXX",
  "secret_key": "XXX",
  "nonce": "1502345965295519602",
  "method": "ticker",
  "params": { // 请求参数, 不同 method 的请求参数不同.
    "symbol": "btc",
  }
}
```

## GetTicker: 用于获取当前行情数据
* method: ticker.
* request.params.symbol: 交易币种.
* return:

```
{
  "data": {
    "time": 1500793319499, // 毫秒时间戳, 整型
    "buy": 1000.0, // 以下浮点型
    "sell": 1001.0,
    "last": 1002.0,
    "high": 1100
  }
}
```

## GetRecords: 用于获取交易所提供的 K 线数据
* method: records.
* request.params.symbol: 交易币种.
* request.params.period: 分钟数, 表示的周期.
* return:

```
{
  "data": [
    [1500793319499, 1.1, 2.2, 3.3, 4.4, 5.5], // "Time" 为整型, 后面都是浮点 "Open", "High", "Low", "Close", "Volume"
    [1500793259499, 1.01, 2.02, 3.03, 4.04, 5.05],
    ...,
  ]
}
```

## GetDepth: 获取交易所的深度信息
* method: depth.
* request.params.symbol: 交易币种.
* return:

```
{
  "data": {
    "time" : 1500793319499,
    "asks" : [ [1000, 0.5], [1001, 0.23], [1004, 2.1], ... ],
    "bids" : [ [999, 0.25], [998, 0.8], [995, 1.4], ... ],
  }
}
```

## GetTrades: 获取交易所提供的 整个交易所一定时间内的交易记录
* method: trades.
* request.params.symbol: 交易币种.
* return:

```
{
  "data": [
    {
      "id": 12232153,
      "price": 1000,
      "amount": 0.5,
      "type": "buy"
    }, {
      "id": 12545664,
      "price": 1001,
      "amount": 1,
      "type": "sell"
    }, {
      ...
    }
  ]
}
```

## GetAccount: 获取账户资产信息
* method: accounts.
* return:

```
{
    "data": [
        {"currency": "btc", "free": 1.2, "frozen": 0.1},
        {"currency": "ltc", "free": 25, "frozen": 2.1},
        {"currency": "ltc", "free": 25, "frozen": 2.1},
        ...
    ]
}
```

## Trade: 发送委托单, 下单交易 (市价单, 限价单)
* method: trade.
* request.params.symbol: 交易币种
* request.params.type: "buy" or "sell".
* request.params.price: 价格
* request.params.amount: 数量
* return:

```
{
    "data": {
        "id": 125456,      // 下单后 返回的订单id
    }
}
```

## GetOrder: 获取指定订单号的订单信息
* method: order.
* request.params.symbol: 交易币种
* request.params.id: 订单 ID.
* return:

```
{
    "data": {
        "id": 2565244,
        "amount": 0.15,
        "price": 1002,
        "status": "open",    // "open": 挂起状态, "closed": 完成关闭状态, "cancelled": 已取消.
        "deal_amount": 0,
        "type": "buy",       // "buy", "sell"
    }
}
```

## GetOrders: 获取所有未完成的订单信息
* method: orders.
* request.params.symbol: 交易币种
* return:

```
{
    "data": [{
        "id": 542156,
        "amount": 0.25,
        "price": 1005,
        "deal_amount": 0,
        "type": "buy",   // "buy" 、 "sell"
        "status": "open", // "open"
    },{
        ...
    }]
}
```

## CancelOrder: 取消指定订单号的订单委托
* method: cancel.
* request.params.symbol: 交易币种
* request.params.id: 订单 ID.
* return:

```
{
    "data": true,        // true or false
}
```

## Refs
* [BotVS 托管者通用协议](https://www.botvs.com/bbs-topic/1052)
