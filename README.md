# Fyers WebSocket Connect Demo

A small Go demo for connecting to the Fyers market data WebSocket and streaming live quote updates for NSE symbols.

## What this repo does

- Uses the Fyers Go SDK and WebSocket client to subscribe to live market data
- Connects to Fyers data socket using an access token
- Subscribes to symbol price updates (`SymbolUpdate`) for a small set of securities
- Parses incoming messages and keeps the latest last traded price (LTP) in memory
- Periodically logs symbol prices to CSV files
- Demonstrates auto-reconnect, error handling, and price fetching fallback logic

## Files

- `main.go` - Main demo application, WebSocket connection logic, message handling, CSV logger
- `decoding-data.go` - Data model definitions for market, index, and depth messages
- `go.mod` - Go module file with dependencies
- `NSE_NIFTY50_INDEX.csv`, `NSE_SBIN_EQ.csv`, `XP03754.json` - Example data files included in repo

## Prerequisites

- Go 1.23 or newer
- Valid Fyers API credentials and access token
- Network access to Fyers WebSocket endpoint

## Dependencies

This project depends on:

- `github.com/FyersDev/fyers-go-sdk`
- `github.com/alphadose/haxmap`
- `github.com/sainipankaj15/All-In-One-Broker`

## Setup

1. Clone the repository.
2. Install dependencies using `go mod tidy`.
3. Provide a Fyers access token.

### Access token

The code expects an access token in the format `appid:accesstoken`.

By default `main.go` reads the token using:

```go
accessToken, err := fyers.ReadingAccessToken_Fyers(fyers.ADMIN_FYERS)
```

If you want to hardcode values for testing, uncomment and configure the `appId` and `token` values in `main.go`.

## Run

From the repository root:

```bash
go run main.go
```

## Behavior

- On connect, the demo subscribes to predefined symbols:
  - `NSE:NIFTY50-INDEX`
  - `NSE:SBIN-EQ`
  - `NSE:MON100-EQ`
  - `NSE:TCS-EQ`
  - `NSE:INFY-EQ`
- Incoming messages are parsed and stored in a concurrent map
- `GetPrice(symbol)` returns the most recent LTP; if the cached tick is older than 5 seconds, it falls back to `fyers.LTP_Fyers`
- `StartCSVLogger` writes price history for configured symbols to CSV files in the current directory

## Customization

- Change `symbols` and `datatype` in `main.go` to subscribe to different instruments or message types
- Enable `DepthUpdate` by setting `datatype := "DepthUpdate"`
- Adjust reconnect behavior and logging options in the `NewFyersDataSocket` call

## Notes

- This repo is intended as a demo and proof of concept, not production-ready trading software
- Keep your Fyers tokens secure and do not commit them to source control
- The example uses simple logging and in-memory state for demonstration

## License

This repository does not include a license file. Add one if you intend to reuse or share this project publicly.
