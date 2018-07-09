package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/newdag/wallet"
)

var m_proxy *SocketProxy

func http_Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sHTML = "finding index.html ..."
		fi, err := os.Open("index.html")

		if err == nil {
			fd, _ := ioutil.ReadAll(fi)
			sHTML = string(fd)
			fi.Close()
		}

		fmt.Fprintln(w, sHTML)
	}
}

func http_GetMsg(sType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch sType {
		case "alltrans":
			fmt.Fprintln(w, GetTransactions(-1))
			return
		case "balance":
			sAddr := r.FormValue("addr")
			fmt.Fprintf(w, "%d", client_GetBalance(sAddr))
		case "addrnew":
			wal := wallet.CreateWalletToFile("")
			if w != nil {
				ss := fmt.Sprintf("%s", wal.GetAddress())
				fmt.Fprintln(w, ss)
			} else {
				fmt.Fprintln(w, "CreateWalletToFile return is nil!")
			}
		case "addrlist":
			ws, err := wallet.GetWallets()
			if err != nil {
				fmt.Fprintln(w, fmt.Sprintf("Error in addrlist: %v\n", err))
				return
			}
			addresses := ws.GetAddresses()
			ss := strings.Join(addresses, "\n")
			fmt.Fprintln(w, ss)
		}
		//end switch
	}
}

func http_PostMsg(sType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case sType == "trans":
			sFrom := r.FormValue("from")
			sTo := r.FormValue("to")
			sAmount := r.FormValue("amount")
			sMsg := r.FormValue("data")
			tran, err := ConvertToTrans(sFrom, sTo, sAmount, sMsg)
			if err != nil {
				sError := fmt.Sprintf("Error in ConvertToTrans: %v\n", err)
				//fmt.Printf(sError)
				fmt.Fprintln(w, sError)
			} else {
				err := m_proxy.SubmitTx(*tran)
				if err != nil {
					sError := fmt.Sprintf("Error in SubmitTx: %v\n", err)
					fmt.Fprintln(w, sError)
					return
				}
				fmt.Fprintln(w, "Submit Tx successed!")
			}
		} //end switch
	}
}

//---------------------------------------------------------------------------------------------------------------------

func http_server_main(listenFlag string, proxy *SocketProxy) {
	m_proxy = proxy
	mux := http.NewServeMux()
	mux.HandleFunc("/", withAppHeaders(http_Home()))
	mux.HandleFunc("/post", withAppHeaders(http_PostMsg("trans")))
	mux.HandleFunc("/get", withAppHeaders(http_GetMsg("alltrans")))
	mux.HandleFunc("/getBalance", withAppHeaders(http_GetMsg("balance")))
	mux.HandleFunc("/keygen", withAppHeaders(http_GetMsg("addrnew")))
	mux.HandleFunc("/addrlist", withAppHeaders(http_GetMsg("addrlist")))
	listenFlag = ":" + listenFlag
	server := &http.Server{Addr: listenFlag, Handler: mux}
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[INFO] server is listening on %s\n", listenFlag)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERR] server exited with: %s", err)
		}
		close(serverCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh // Wait for interrupt

	log.Printf("[INFO] received interrupt, shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}

func withAppHeaders(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-App-Name", "newDAG")
		h(w, r)
	}
}
