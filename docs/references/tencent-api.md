# Tencent Finance API Reference

## Endpoint

```
GET http://qt.gtimg.cn/q=<sym1>,<sym2>,...
```

Multiple symbols comma-separated in a single request (batch).

## Symbol Format

| Market | Format | Example |
|--------|--------|---------|
| Shanghai A-share | `sh<code>` | `sh600519` |
| Shenzhen A-share | `sz<code>` | `sz000858` |
| Hong Kong | `hk<code>` | `hk00700` |
| US stocks | `us<ticker>` | `usAAPL` |

## Response Format

- Encoding: **GBK** (must be decoded to UTF-8)
- Structure: one line per symbol, tilde (`~`) delimited
- Example line:

```
v_sh600519="1~贵州茅台~600519~1474.98~1485.00~1480.00~..."
```

Field positions (0-indexed after splitting by `~`):

| Index | Meaning |
|-------|---------|
| 0 | Market prefix (e.g. `1` = SH) |
| 1 | Stock name |
| 2 | Stock code |
| 3 | Current price |
| 4 | Previous close |
| 5 | Open price |
| 31 | High price |
| 32 | Low price |
| 36 | Volume (lots) |
| 37 | Amount (yuan) |

## Encoding Note

Tencent's response is GBK-encoded. Use `golang.org/x/text/encoding/simplifiedchinese` to
transcode to UTF-8 before parsing. See `internal/api/tencent.go`.
