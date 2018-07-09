package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/newdag/config"
)

const (
	newdag_Version = "0.0.1"
	newdag_cfgfile = "newdag.cfg"

	jsonPeer_File = "peers.json"
	wallet_File   = "newdag_wallet.dat"
	pemKey_File   = "priv_key.pem"
)

type Config struct {
	WalletAddr       string
	Datadir          string
	NoClient         bool
	MaxPool          int
	CacheSize        int
	SyncLimit        int
	HeartbeatTimeout time.Duration
	TCPTimeout       time.Duration
	StoreType        string
	StorePath        string
	NodePort         int
	ProxyPort        int
	ServicePort      int
	ClientAddress    string
	Logger           *Logger
}

func NewConfig(nodePort int, proxyPort int, servicePort int, clientAddress string, logger *Logger) *Config {
	conf := new(Config)
	conf.Logger = logger
	cfg, err := config.ReadFile(newdag_cfgfile)
	if err != nil {
		if nodePort > 0 {
			fmt.Println("can't found %s in current dir!", newdag_cfgfile)
			return nil
		} else {
			cfg = new(config.Config)
		}
	}
	conf.WalletAddr = cfg.GetValue("", "address", "").(string)
	conf.Datadir = cfg.GetValue("server", "datadir", DefaultDataDir()).(string)
	conf.NoClient = cfg.GetValue("server", "noClient", false).(bool)
	conf.MaxPool = cfg.GetValue("server", "maxPool", 2).(int)
	conf.CacheSize = cfg.GetValue("server", "cacheSize", 40960).(int)
	conf.SyncLimit = cfg.GetValue("server", "syncLimit", 500).(int)
	conf.HeartbeatTimeout = time.Duration(cfg.GetValue("server", "heartBeat", 50).(int)) * time.Millisecond
	conf.TCPTimeout = time.Duration(cfg.GetValue("server", "tcpTimeout", 200).(int)) * time.Millisecond
	conf.StoreType = cfg.GetValue("server", "storeType", "inmem").(string)
	conf.StorePath = cfg.GetValue("server", "storePath", filepath.Join(conf.Datadir, "newdag_db")).(string)

	conf.NodePort = cfg.GetValue("node", "nodePort", nodePort).(int)
	conf.ProxyPort = cfg.GetValue("node", "proxyPortPort", proxyPort).(int)
	conf.ServicePort = cfg.GetValue("node", "servicePort", servicePort).(int)
	conf.ClientAddress = cfg.GetValue("node", "clientAddress", clientAddress).(string)

	//fmt.Println(*conf)

	return conf
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	//home := HomeDir()
	//if home != ""
	{
		if runtime.GOOS == "darwin" {
			return "./" //filepath.Join(home, ".newdag")
		} else if runtime.GOOS == "windows" {
			return "" //filepath.Join(home, "newdag")
		} else {
			return "./" //filepath.Join(home, ".newdag")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func ReadFile(sModule string) ([]byte, error) {
	sFile := ""
	switch sModule {
	case "wallet":
		sFile = wallet_File
	case "peer":
		sFile = jsonPeer_File
	case "peer_key":
		sFile = pemKey_File
	}

	if _, err := os.Stat(sFile); os.IsNotExist(err) {
		return nil, err
	}

	fileContent, err := ioutil.ReadFile(sFile)
	if err != nil {
		log.Panic(err)
	}

	return fileContent, err
}

// SaveToFile saves wallets to a file
func WriteFile(sModule string, content []byte) error {
	sFile := ""
	switch sModule {
	case "wallet":
		sFile = wallet_File
	case "peer":
		sFile = jsonPeer_File
	case "peer_key":
		sFile = pemKey_File
	}
	fmt.Println("WriteFile", sFile)

	err := ioutil.WriteFile(sFile, content, 0644)
	if err != nil {
		log.Panic(err)
	}
	return err
}

func GetVersion() string {
	return newdag_Version
}
