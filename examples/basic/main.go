package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KhavrTrading/flowex/binance"
	"github.com/KhavrTrading/flowex/models"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)

	mgr := binance.NewManager()

	symbol := "BTCUSDT"
	fmt.Printf("Subscribing to %s on Binance...\n", symbol)

	// Optional: set hooks on the worker before subscribing
	worker := mgr.GetOrCreateWorker(symbol)
	worker.SetOnCandleUpdate(func(candles []models.CandleHLCV) {
		if len(candles) > 0 {
			last := candles[len(candles)-1]
			fmt.Printf("  Candle update: C=%.2f V=%.2f (%d candles)\n", last.Close, last.Volume, len(candles))
		}
	})

	err := mgr.SubscribeAll(symbol, nil, nil, nil)
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	// Print snapshots every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			snap := mgr.GetSnapshot(symbol)
			if snap == nil {
				fmt.Println("No snapshot yet...")
				continue
			}
			fmt.Printf(
				"[%s] Candles: %d | Trades: %d | Depth metrics: %d\n",
				snap.Timestamp.Format("15:04:05"),
				len(snap.Candles),
				len(snap.Trades),
				snap.DepthStore.Size(),
			)
			if len(snap.Candles) > 0 {
				last := snap.Candles[len(snap.Candles)-1]
				fmt.Printf("  Last candle: O=%.2f H=%.2f L=%.2f C=%.2f V=%.2f\n",
					last.Open, last.High, last.Low, last.Close, last.Volume)
			}
			metrics := worker.GetMetrics()
			fmt.Printf("  Processed: %d | Dropped: c=%d d=%d t=%d\n",
				metrics["processed"],
				metrics["candle_dropped"],
				metrics["depth_dropped"],
				metrics["trade_dropped"],
			)
		}
	}()

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nShutting down...")
	mgr.Shutdown()
}
