# Eastmoney API Reference

## Endpoint

```
GET https://push2.eastmoney.com/api/qt/stock/get?secid=<SECID>&fields=f43,f57,f58,f169,f170,f46,f44,f45,f47,f48,f60
```

## secid Format

| Market | Format | Example |
|--------|--------|---------|
| Shanghai A-share | `1.<code>` | `1.600519` |
| Shenzhen A-share | `0.<code>` | `0.000858` |
| Hong Kong | `116.<code>` | `116.00700` (zero-padded to 5 digits) |
| NASDAQ | `105.<ticker>` | `105.AAPL` |
| NYSE | `106.<ticker>` | `106.BRK-A` |

## US Stock Retry Logic

US stocks default to NASDAQ prefix `105.`. If the response contains no data (empty `f57`),
the client automatically retries with NYSE prefix `106.`.

## Response Format

JSON. Key fields used:

| Field | Meaning |
|-------|---------|
| `f43` | Current price (×100, integer) |
| `f57` | Stock code |
| `f58` | Stock name |
| `f169` | Change amount (×100) |
| `f170` | Change percent (×100) |
| `f46` | Open price (×100) |
| `f44` | High price (×100) |
| `f45` | Low price (×100) |
| `f47` | Volume (lots) |
| `f48` | Amount (yuan) |
| `f60` | Previous close (×100) |

All prices are integers representing the actual value × 100 (i.e., divide by 100 to get
the display value).
