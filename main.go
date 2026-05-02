package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	fyersws "github.com/FyersDev/fyers-go-sdk/websocket"
	fyers "github.com/sainipankaj15/All-In-One-Broker/Fyers"
)

func main() {
	// appId := "AAAAAAAAA-100"
	// token := "eyjb...."
	// accessToken := fmt.Sprintf("%s:%s", appId, token)
	accessToken, err := fyers.ReadingAccessToken_Fyers(fyers.ADMIN_FYERS)
	if err != nil {
		fmt.Printf("Error reading access token: %v\n", err)
		return
	}
	// symbols := []string{"NSE:NIFTY50-INDEX"}
	symbols := []string{"NSE:SBIN-EQ"}
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

	symbols = []string{"NSE:MON100-EQ"}
	dataSocket.Subscribe(symbols, datatype)

	symbols = []string{"NSE:TCS-EQ"}
	dataSocket.Subscribe(symbols, datatype)

	symbols = []string{"NSE:INFY-EQ"}
	dataSocket.Subscribe(symbols, datatype)

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

	fmt.Printf("Response after marshaling: %s\n", msgBytes)

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

	case "if":
		var data IndexData
		if err := json.Unmarshal(msgBytes, &data); err != nil {
			fmt.Println("Unmarshal failed (if):", err)
			return
		}

		fmt.Printf("[IF] %s | LTP: %.2f\n", data.Symbol, data.LTP)

	case "dp":
		var data DepthData
		if err := json.Unmarshal(msgBytes, &data); err != nil {
			fmt.Println("Unmarshal failed (dp):", err)
			return
		}

		// 🔥 Best Bid/Ask (Level 1)
		fmt.Printf("[DP] %s | Bid: %.2f (%d) | Ask: %.2f (%d)\n",
			data.Symbol,
			data.BidPrice1, data.BidSize1,
			data.AskPrice1, data.AskSize1)

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
