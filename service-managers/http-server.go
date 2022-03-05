package manager

import (
	"bridge/common"
	"fmt"
	"net/http"
)

func MkHttpServer(cnf common.HttpConf, engine http.Handler) *http.Server {

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		Handler: engine,
		//	TLSConfig *tls.Config
		//ReadTimeout:    readTimeout,
		//WriteTimeout:   writeTimeout,
		//MaxHeaderBytes: maxHeaderBytes,
	}
}
