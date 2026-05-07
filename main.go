package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alphadose/haxmap"

	fyersws "github.com/FyersDev/fyers-go-sdk/websocket"
	fyers "github.com/sainipankaj15/All-In-One-Broker/Fyers"
)

type Tick struct {
	LTP       float64
	Timestamp time.Time
}

// Key = Symbol
var symbolMap = haxmap.New[string, Tick]()

func main() {
	// appId := "AAAAAAAAA-100"
	// token := "eyjb...."
	// accessToken := fmt.Sprintf("%s:%s", appId, token)
	accessToken, err := fyers.ReadingAccessToken_Fyers(fyers.ADMIN_FYERS)
	if err != nil {
		fmt.Printf("Error reading access token: %v\n", err)
		return
	}
	symbols := []string{"NSE:NIFTY50-INDEX", "NSE:SBIN-EQ", "NSE:MON100-EQ", "NSE:TCS-EQ", "NSE:INFY-EQ"}
	// symbols := []string{"NSE:SBIN-EQ"}
	datatype := "SymbolUpdate" // "SymbolUpdate", "DepthUpdate"
	// datatype := "DepthUpdate" // "SymbolUpdate", "DepthUpdate"

	var dataSocket *fyersws.FyersDataSocket
	onConnect := func() {
		dataSocket.Subscribe(symbols, datatype)
	}

	dataSocket = fyersws.NewFyersDataSocket(
		accessToken, // Access token in the format "appid:accesstoken"
		"",          // Log path - leave empty to auto-create logs in the current directory
		false,       // Lite mode disabled. Set to true if you want a lite response
		false,       // Save response in a log file instead of printing it
		true,        // Enable auto-reconnection to WebSocket on disconnection
		50,          // reconnectRetry: max reconnect attempts (same as Python default; cap 50)
		onConnect,   // Callback: subscribe on every connect (first + after reconnect)
		onClose,     // Callback function to handle WebSocket connection close events
		onError,     // Callback function to handle WebSocket errors
		onMessage,   // Callback function to handle incoming messages from the WebSocket
	)

	err = dataSocket.Connect()
	if err != nil {
		fmt.Printf("failed to connect to Data Socket: %v", err)
		return
	}
	time.Sleep(2 * time.Second)

	price, err := GetPrice("NSE:SBIN-EQ")
	if err != nil {
		fmt.Printf("Error fetching price for NSE:SBIN-EQ: %v\n", err)
	} else {
		fmt.Printf("Current price for NSE:SBIN-EQ: %.2f\n", price)
	}

	price, err = GetPrice("NSE:NIFTY50-INDEX")
	if err != nil {
		fmt.Printf("Error fetching price for NSE:NIFTY50-INDEX: %v\n", err)
	} else {
		fmt.Printf("Current price for NSE:NIFTY50-INDEX: %.2f\n", price)
	}

	// Start CSV loggers in separate goroutines
	go StartCSVLogger("NSE:SBIN-EQ")
	go StartCSVLogger("NSE:NIFTY50-INDEX")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nReceived interrupt signal, closing connection...")

	dataSocket.CloseConnection()
	fmt.Println("Data Socket connection closed")

}

func onMessage(message fyersws.DataResponse) {
	fmt.Printf("Response: %s\n", message)

	// Convert map → JSON bytes
	msgBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Marshal failed:", err)
		return
	}

	// fmt.Printf("Response after marshaling: %s\n", msgBytes)

	var base BaseMessage
	if err := json.Unmarshal(msgBytes, &base); err != nil {
		fmt.Println("Failed to unmarshal base:", err)
		return
	}

	switch base.Type {

	case "sf":
		var data MarketData
		if err := json.Unmarshal(msgBytes, &data); err != nil {
			fmt.Println("Unmarshal failed (sf):", err)
			return
		}

		fmt.Printf("[SF] %s | LTP: %.2f\n", data.Symbol, data.LTP)
		symbolMap.Set(data.Symbol, Tick{
			LTP:       data.LTP,
			Timestamp: time.Now(),
		})

	case "if":
		var data IndexData
		if err := json.Unmarshal(msgBytes, &data); err != nil {
			fmt.Println("Unmarshal failed (if):", err)
			return
		}

		fmt.Printf("[IF] %s | LTP: %.2f\n", data.Symbol, data.LTP)
		symbolMap.Set(data.Symbol, Tick{
			LTP:       data.LTP,
			Timestamp: time.Now(),
		})

	case "dp":
		var data DepthData
		if err := json.Unmarshal(msgBytes, &data); err != nil {
			fmt.Println("Unmarshal failed (dp):", err)
			return
		}

		// Best Bid/Ask (Level 1)
		// fmt.Printf("[DP] %s | Bid: %.2f (%d) | Ask: %.2f (%d)\n",
		// 	data.Symbol,
		// 	data.BidPrice1, data.BidSize1,
		// 	data.AskPrice1, data.AskSize1)
	default:
		// ignore but already logged
		return
	}

}

func onError(message fyersws.DataError) {
	fmt.Printf("Error: %s\n", message)
}

func onClose(message fyersws.DataClose) {
	fmt.Printf("Connection closed: %s\n", message)
}

func GetPrice(symbol string) (float64, error) {

	tick, ok := symbolMap.Get(symbol)

	if !ok {
		return 0, errors.New("symbol not found")
	}

	if time.Since(tick.Timestamp) > 5*time.Second {
		// If the tick data is older than 5 seconds, consider it outdated
		// and do your things which you want to do when price is outdated
		// I will take price from other source and update the symbolMap with new price and timestamp

		fmt.Printf("Outdated price for %s\n", symbol)
		fmt.Printf("Fetching from other sources")

		ltp, err := fyers.LTP_Fyers(symbol, "XP03754")
		if err != nil {
			return 0, errors.New("failed to fetch updated price")
		}
		symbolMap.Set(symbol, Tick{
			LTP:       ltp,
			Timestamp: time.Now(),
		})
		return ltp, nil
	}

	return tick.LTP, nil
}

// Write price every second into CSV
func StartCSVLogger(symbol string) {

	// Replace special chars for safe filename
	fileName := strings.ReplaceAll(symbol, ":", "_")
	fileName = strings.ReplaceAll(fileName, "-", "_")
	fileName += ".csv"

	file, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)

	if err != nil {
		fmt.Println("file open error:", err)
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Optional CSV Header
	fileInfo, _ := file.Stat()

	if fileInfo.Size() == 0 {
		writer.Write([]string{"timestamp", "price"})
		writer.Flush()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		price, err := GetPrice(symbol)

		if err != nil {
			fmt.Println("price fetch error:", err)
			continue
		}

		record := []string{
			time.Now().Format(time.RFC3339),
			fmt.Sprintf("%.2f", price),
		}

		err = writer.Write(record)

		if err != nil {
			fmt.Println("csv write error:", err)
			continue
		}

		writer.Flush()

		fmt.Println("Written:", record)
	}
}
