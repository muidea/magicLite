package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"
	engine "github.com/muidea/magicEngine"

	_ "github.com/muidea/magicLite/config"
	"github.com/muidea/magicLite/core"
)

var listenPort = "8880"
var endpointName = "magicLite"

var logConfig = `
<seelog type="sync">
	<outputs formatid="main">
		<console/>
	</outputs>
	<formats>
		<format id="main" format="%Date %Time [%LEV] %Msg%n"/>
	</formats>
</seelog>`

func initPprofMonitor(listenPort string) error {
	var err error
	addr := ":1" + listenPort

	go func() {
		err = http.ListenAndServe(addr, nil)
		if err != nil {
			log.Critical("funcRetErr=http.ListenAndServe||err=%s", err.Error())
		}
	}()

	return err
}

func main() {
	flag.StringVar(&listenPort, "ListenPort", listenPort, "magicLite listen address")
	flag.StringVar(&endpointName, "EndpointName", endpointName, "magicLite endpoint name.")
	flag.Parse()

	initPprofMonitor(listenPort)

	logger, _ := log.LoggerFromConfigAsBytes([]byte(logConfig))
	log.ReplaceLogger(logger)

	log.Info("magicLite V1.0")

	router := engine.NewRouter()
	core, err := core.New(endpointName)

	if err == nil {
		core.Startup(router)

		svr := engine.NewHTTPServer(listenPort)
		svr.Bind(router)

		svr.Run()
	} else {
		log.Critical("start magicLite failed.")
	}

	core.Teardown()
}
