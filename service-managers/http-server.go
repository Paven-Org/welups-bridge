package manager

import (
	"bridge/common"
	"fmt"
	"net/http"
)

func MkHttpServer(cnf common.HttpConf, engine http.Handler) *http.Server {
	//tlsConfig := &tls.Config{
	//	MinVersion:               tls.VersionTLS12,
	//	PreferServerCipherSuites: true,
	//}

	//tlsConfig.Certificates = make([]tls.Certificate, 1)
	//var err error
	//if tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(cnf.X509CertFile, cnf.X509KeyFile); err != nil {
	//	logger.Get().Err(err).Msg("Unable to load TLS certificates, aborting...")
	//	panic(err)
	//}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		Handler: engine,
		//TLSConfig: tlsConfig,
		//ReadTimeout:    readTimeout,
		//WriteTimeout:   writeTimeout,
		//MaxHeaderBytes: maxHeaderBytes,
	}
}
