package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

type Exchange struct {
	anx         ANX
	btcchina    BTCChina
	bitstamp    Bitstamp
	bitfinex    Bitfinex
	btce        BTCE
	btcmarkets  BTCMarkets
	coinbase    Coinbase
	cryptsy     Cryptsy
	okcoinChina OKCoin
	okcoinIntl  OKCoin
	itbit       ItBit
	lakebtc     LakeBTC
	huobi       HUOBI
	kraken      Kraken
}

type Bot struct {
	config   Config
	exchange Exchange
	shutdown chan bool
}

var bot Bot

func main() {
	HandleInterrupt()
	log.Println("Loading config file config.json..")

	err := errors.New("")
	bot.config, err = ReadConfig()

	if err != nil {
		log.Println("Fatal error opening config.json file. Error:", err)
		return
	}

	err = CheckConfigValues()

	if err != nil {
		log.Println("Fatal error checking config values. Error:", err)
		return
	}

	log.Println("Config file loaded.")
	log.Printf("Bot '%s' started.\n", bot.config.Name)

	enabledExchanges := 0
	for _, exch := range bot.config.Exchanges {
		if exch.Enabled {
			enabledExchanges++
		}
	}

	if enabledExchanges == 0 {
		log.Println("Bot started with no exchanges supported. Exiting.")
		return
	}

	AdjustGoMaxProcs()

	smsSupport := false
	smsContacts := 0

	for _, sms := range bot.config.SMSContacts {
		if sms.Enabled {
			smsSupport = true
			smsContacts++
		}
	}

	if smsSupport {
		log.Printf("SMS support enabled. Number of SMS contacts %d.\n", smsContacts)
	} else {
		log.Println("SMS support disabled.")
	}

	log.Printf("Available Exchanges: %d. Enabled Exchanges: %d.\n", len(bot.config.Exchanges), enabledExchanges)
	log.Println("Bot Exchange support:")

	bot.exchange.anx.SetDefaults()
	bot.exchange.kraken.SetDefaults()
	bot.exchange.btcchina.SetDefaults()
	bot.exchange.bitstamp.SetDefaults()
	bot.exchange.bitfinex.SetDefaults()
	bot.exchange.btce.SetDefaults()
	bot.exchange.btcmarkets.SetDefaults()
	bot.exchange.coinbase.SetDefaults()
	bot.exchange.cryptsy.SetDefaults()
	bot.exchange.okcoinChina.SetURL(OKCOIN_API_URL_CHINA)
	bot.exchange.okcoinChina.SetDefaults()
	bot.exchange.okcoinIntl.SetURL(OKCOIN_API_URL)
	bot.exchange.okcoinIntl.SetDefaults()
	bot.exchange.itbit.SetDefaults()
	bot.exchange.lakebtc.SetDefaults()
	bot.exchange.huobi.SetDefaults()

	err = RetrieveConfigCurrencyPairs(bot.config)

	if err != nil {
		log.Println("Fatal error retrieving config currency AvailablePairs. Error: ", err)
	}

	for _, exch := range bot.config.Exchanges {
		if exch.Enabled {
			log.Printf("%s: Exchange support: %s (Authenticated API support: %s - Verbose mode: %s).\n", exch.Name, IsEnabled(exch.Enabled), IsEnabled(exch.AuthenticatedAPISupport), IsEnabled(exch.Verbose))
		} else {
			log.Printf("%s: Exchange support: %s\n", exch.Name, IsEnabled(exch.Enabled))
		}

		if bot.exchange.anx.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.anx.SetEnabled(false)
			} else {
				bot.exchange.anx.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.anx.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.anx.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.anx.Verbose = exch.Verbose
				bot.exchange.anx.Websocket = exch.Websocket
				bot.exchange.anx.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.anx.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.anx.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.anx.Run()
			}
		} else if bot.exchange.btcchina.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.btcchina.SetEnabled(false)
			} else {
				bot.exchange.btcchina.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.btcchina.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.btcchina.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.btcchina.Verbose = exch.Verbose
				bot.exchange.btcchina.Websocket = exch.Websocket
				bot.exchange.btcchina.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.btcchina.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.btcchina.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.btcchina.Run()
			}
		} else if bot.exchange.bitstamp.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.bitstamp.SetEnabled(false)
			} else {
				bot.exchange.bitstamp.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.bitstamp.SetAPIKeys(exch.ClientID, exch.APIKey, exch.APISecret)
				bot.exchange.bitstamp.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.bitstamp.Verbose = exch.Verbose
				bot.exchange.bitstamp.Websocket = exch.Websocket
				bot.exchange.bitstamp.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.bitstamp.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.bitstamp.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.bitstamp.Run()
			}
		} else if bot.exchange.bitfinex.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.bitfinex.SetEnabled(false)
			} else {
				bot.exchange.bitfinex.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.bitfinex.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.bitfinex.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.bitfinex.Verbose = exch.Verbose
				bot.exchange.bitfinex.Websocket = exch.Websocket
				bot.exchange.bitfinex.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.bitfinex.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.bitfinex.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.bitfinex.Run()
			}
		} else if bot.exchange.btce.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.btce.SetEnabled(false)
			} else {
				bot.exchange.btce.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.btce.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.btce.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.btce.Verbose = exch.Verbose
				bot.exchange.btce.Websocket = exch.Websocket
				bot.exchange.btce.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.btce.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.btce.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.btce.Run()
			}
		} else if bot.exchange.btcmarkets.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.btcmarkets.SetEnabled(false)
			} else {
				bot.exchange.btcmarkets.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.btcmarkets.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.btcmarkets.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.btcmarkets.Verbose = exch.Verbose
				bot.exchange.btcmarkets.Websocket = exch.Websocket
				bot.exchange.btcmarkets.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.btcmarkets.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.btcmarkets.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.btcmarkets.Run()
			}
		} else if bot.exchange.coinbase.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.coinbase.SetEnabled(false)
			} else {
				bot.exchange.coinbase.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.coinbase.SetAPIKeys(exch.ClientID, exch.APIKey, exch.APISecret)
				bot.exchange.coinbase.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.coinbase.Verbose = exch.Verbose
				bot.exchange.coinbase.Websocket = exch.Websocket
				bot.exchange.coinbase.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.coinbase.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.coinbase.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.coinbase.Run()
			}
		} else if bot.exchange.cryptsy.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.cryptsy.SetEnabled(false)
			} else {
				bot.exchange.cryptsy.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.cryptsy.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.cryptsy.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.cryptsy.Verbose = exch.Verbose
				bot.exchange.cryptsy.Websocket = exch.Websocket
				bot.exchange.cryptsy.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.cryptsy.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.cryptsy.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.cryptsy.Run()
			}
		} else if bot.exchange.okcoinChina.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.okcoinChina.SetEnabled(false)
			} else {
				bot.exchange.okcoinChina.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.okcoinChina.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.okcoinChina.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.okcoinChina.Verbose = exch.Verbose
				bot.exchange.okcoinChina.Websocket = exch.Websocket
				bot.exchange.okcoinChina.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.okcoinChina.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.okcoinChina.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.okcoinChina.Run()
			}
		} else if bot.exchange.okcoinIntl.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.okcoinIntl.SetEnabled(false)
			} else {
				bot.exchange.okcoinIntl.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.okcoinIntl.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.okcoinIntl.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.okcoinIntl.Verbose = exch.Verbose
				bot.exchange.okcoinIntl.Websocket = exch.Websocket
				bot.exchange.okcoinIntl.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.okcoinIntl.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.okcoinIntl.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.okcoinIntl.Run()
			}
		} else if bot.exchange.itbit.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.itbit.SetEnabled(false)
			} else {
				bot.exchange.itbit.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.itbit.SetAPIKeys(exch.APIKey, exch.APISecret, exch.ClientID)
				bot.exchange.itbit.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.itbit.Verbose = exch.Verbose
				bot.exchange.itbit.Websocket = exch.Websocket
				bot.exchange.itbit.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.itbit.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.itbit.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.itbit.Run()
			}
		} else if bot.exchange.kraken.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.kraken.SetEnabled(false)
			} else {
				bot.exchange.kraken.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.kraken.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.kraken.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.kraken.Verbose = exch.Verbose
				bot.exchange.kraken.Websocket = exch.Websocket
				bot.exchange.kraken.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.kraken.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.kraken.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.kraken.Run()
			}
		} else if bot.exchange.lakebtc.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.lakebtc.SetEnabled(false)
			} else {
				bot.exchange.lakebtc.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.lakebtc.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.lakebtc.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.lakebtc.Verbose = exch.Verbose
				bot.exchange.lakebtc.Websocket = exch.Websocket
				bot.exchange.lakebtc.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.lakebtc.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.lakebtc.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.lakebtc.Run()
			}
		} else if bot.exchange.huobi.GetName() == exch.Name {
			if !exch.Enabled {
				bot.exchange.huobi.SetEnabled(false)
			} else {
				bot.exchange.huobi.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
				bot.exchange.huobi.SetAPIKeys(exch.APIKey, exch.APISecret)
				bot.exchange.huobi.RESTPollingDelay = exch.RESTPollingDelay
				bot.exchange.huobi.Verbose = exch.Verbose
				bot.exchange.huobi.Websocket = exch.Websocket
				bot.exchange.huobi.BaseCurrencies = SplitStrings(exch.BaseCurrencies, ",")
				bot.exchange.huobi.AvailablePairs = SplitStrings(exch.AvailablePairs, ",")
				bot.exchange.huobi.EnabledPairs = SplitStrings(exch.EnabledPairs, ",")
				go bot.exchange.huobi.Run()
			}
		}
	}
	<-bot.shutdown
	Shutdown()
}

func AdjustGoMaxProcs() {
	log.Println("Adjusting bot runtime performance..")
	maxProcsEnv := os.Getenv("GOMAXPROCS")
	maxProcs := runtime.NumCPU()
	log.Println("Number of CPU's detected:", maxProcs)

	if maxProcsEnv != "" {
		log.Println("GOMAXPROCS env =", maxProcsEnv)
		env, err := strconv.Atoi(maxProcsEnv)

		if err != nil {
			log.Println("Unable to convert GOMAXPROCS to int, using", maxProcs)
		} else {
			maxProcs = env
		}
	}
	log.Println("Set GOMAXPROCS to:", maxProcs)
	runtime.GOMAXPROCS(maxProcs)
}

func HandleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Printf("Captured %v.", sig)
		Shutdown()
	}()
}

func Shutdown() {
	log.Println("Bot shutting down..")
	err := SaveConfig()

	if err != nil {
		log.Println("Unable to save config.")
	} else {
		log.Println("Config file saved successfully.")
	}

	log.Println("Exiting.")
	os.Exit(1)
}
