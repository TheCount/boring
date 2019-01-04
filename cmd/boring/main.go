package main

import (
	"os"

	"github.com/TheCount/boring/app"
	"github.com/TheCount/boring/config"
	"github.com/TheCount/boring/wallet"
	"github.com/TheCount/boring/web"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmnode "github.com/tendermint/tendermint/node"
	tmp2p "github.com/tendermint/tendermint/p2p"
	tmpv "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err.Error())
	}
	pv := tmpv.LoadOrGenFilePV(cfg.TMConfig.PrivValidatorFile())
	nodeKey, err := tmp2p.LoadOrGenNodeKey(cfg.TMConfig.NodeKeyFile())
	if err != nil {
		panic(err.Error())
	}
	walletManager, err := wallet.NewManager(cfg.WalletConfig)
	if err != nil {
		panic(err.Error())
	}
	app, err := app.NewApp(cfg.AppConfig)
	if err != nil {
		panic(err.Error())
	}
	clientCreator := proxy.NewLocalClientCreator(app)
	genesisDocProvider := tmnode.DefaultGenesisDocProviderFunc(cfg.TMConfig)
	metricsProvider := tmnode.DefaultMetricsProvider(cfg.TMConfig.Instrumentation)
	tmlogger := tmlog.NewTMLogger(os.Stdout)
	n, err := tmnode.NewNode(
		cfg.TMConfig, pv, nodeKey,
		clientCreator, genesisDocProvider, tmnode.DefaultDBProvider,
		metricsProvider, tmlogger,
	)
	if err != nil {
		panic(err.Error())
	}
	client := tmclient.NewLocal(n)
	webServer := web.NewWeb(client, walletManager)
	n.Start()
	go webServer.Serve()
	<-n.Quit()
	if err := walletManager.CloseAllWallets(); err != nil {
		panic(err.Error())
	}
}
