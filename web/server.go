package web

type Configuration struct {
	Hostname  string
	UseHTTP   bool
	UseHTTPS  bool
	HTTPPort  int
	HTTPSPort int
	CertFile  string
	KeyFile   string
}

func startHTTPS()

func Run(httpHandlers http.Handler, httpsHandlers http.Handler, s Configuration) {
	if s.UseHTTP && s.UseHTTPS {
		go func() {
			startHTTPS(httpsHandlers, s)
		}()

		startHTTP(httpHandlers, s)
	} else if s.UseHTTP {
		startHTTP(httpHandlers, s)
	} else if s.UseHTTPS {
		startHTTPS(httpsHandlers, s)
	}
}
