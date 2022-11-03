package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/fullstorydev/grpcui/standalone"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"

	"github.com/s77rt/rdpcloud/client/go/license"
)

var (
	Version string = "dev"

	licenseInfo *license.License
)

func init() {
	// Read License
	var err error
	licenseInfo, err = license.Read()
	if err != nil {
		log.Fatalf("Failed to read license: %v", err)
	}
}

type fakeResponseWriter struct {
	http.ResponseWriter

	buf *bytes.Buffer
}

func (fw *fakeResponseWriter) Write(b []byte) (int, error) {
	if fw.buf == nil {
		return len(b), nil
	}

	return fw.buf.Write(b)
}

//go:embed cert
var c embed.FS

//go:embed frontend
var f embed.FS

func main() {
	var port int
	var webport int
	var certFile string
	var keyFile string
	flag.IntVar(&port, "port", 5027, "port on which the server is listening")
	flag.IntVar(&webport, "webport", 5028, "port on which the web server will listen")
	flag.StringVar(&certFile, "certFile", "", "certFile which the web server will use for incoming HTTPS connections")
	flag.StringVar(&keyFile, "keyFile", "", "keyFile which the web server will use for incoming HTTPS connections")
	flag.Parse()

	if certFile == "" || keyFile == "" {
		if certFile == "" {
			fmt.Fprintln(os.Stderr, "missing required -certFile argument")
		}
		if keyFile == "" {
			fmt.Fprintln(os.Stderr, "missing required -keyFile argument")
		}
		flag.Usage()
		os.Exit(2) // the same exit code flag.Parse uses
	}

	log.Printf("Running RDPCloud Client (Version: %s)", Version)

	if licenseInfo.ExpDate.IsZero() {
		log.Printf("Licensed to %s (%s | %s)", licenseInfo.ServerName, licenseInfo.ServerLocalIP, licenseInfo.ServerPublicIP)
	} else {
		log.Printf("Licensed to %s (%s | %s) [Exp. Date: %s]", licenseInfo.ServerName, licenseInfo.ServerLocalIP, licenseInfo.ServerPublicIP, licenseInfo.ExpDate)
	}

	addr := fmt.Sprintf("%s:%d", licenseInfo.ServerPublicIP, port)
	target := licenseInfo.ServerName

	serverCert, err := c.ReadFile("cert/server-cert.pem")
	if err != nil {
		log.Fatalf("Failed to read server-cert pem file: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(serverCert) {
		log.Fatalf("Failed to add server certificate")
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, "")),
	}

	conn, err := grpc.Dial(addr, grpcOpts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	index, err := f.ReadFile("frontend/templates/index.tmpl")
	if err != nil {
		log.Fatalf("Failed to read index tmpl file: %v", err)
	}
	indexTemplate := template.Must(template.New("index").Parse(string(index)))

	// Minify the executed index template on the first http request
	// and store the minified version in this variable for further http requests (caching)
	var executedIndexTemplateMinified []byte

	favicon, err := f.ReadFile("frontend/static/assets/images/favicon.png")
	if err != nil {
		log.Fatalf("Failed to read favicon png file: %v", err)
	}

	jqueryJS, err := f.ReadFile("frontend/static/assets/js/jquery.min.js")
	if err != nil {
		log.Fatalf("Failed to read jquery js file: %v", err)
	}
	jqueryJSMinified, err := m.Bytes("text/javascript", jqueryJS)
	if err != nil {
		log.Fatalf("Failed to minify jquery js file contents: %v", err)
	}

	jqueryUiJS, err := f.ReadFile("frontend/static/assets/js/jquery-ui.min.js")
	if err != nil {
		log.Fatalf("Failed to read jquery-ui js file: %v", err)
	}
	jqueryUiJSMinified, err := m.Bytes("text/javascript", jqueryUiJS)
	if err != nil {
		log.Fatalf("Failed to minify jquery-ui js file contents: %v", err)
	}

	jqueryUiCSS, err := f.ReadFile("frontend/static/assets/css/jquery-ui.min.css")
	if err != nil {
		log.Fatalf("Failed to read jquery-ui css file: %v", err)
	}
	jqueryUiCSSMinified, err := m.Bytes("text/css", jqueryUiCSS)
	if err != nil {
		log.Fatalf("Failed to minify jquery-ui css file contents: %v", err)
	}

	grpcWebFormJS, err := f.ReadFile("frontend/static/assets/js/grpc-web-form.js")
	if err != nil {
		log.Fatalf("Failed to read grpc-web-form js file: %v", err)
	}
	grpcWebFormJSMinified, err := m.Bytes("text/javascript", grpcWebFormJS)
	if err != nil {
		log.Fatalf("Failed to minify grpc-web-form js file contents: %v", err)
	}

	grpcWebFormCSS, err := f.ReadFile("frontend/static/assets/css/grpc-web-form.css")
	if err != nil {
		log.Fatalf("Failed to read grpc-web-form css file: %v", err)
	}
	grpcWebFormCSSMinified, err := m.Bytes("text/css", grpcWebFormCSS)
	if err != nil {
		log.Fatalf("Failed to minify grpc-web-form css file contents: %v", err)
	}

	AppMainJS, err := f.ReadFile("frontend/static/assets/js/app/main.js")
	if err != nil {
		log.Fatalf("Failed to read app main js file: %v", err)
	}
	AppMainJSMinified, err := m.Bytes("text/javascript", AppMainJS)
	if err != nil {
		log.Fatalf("Failed to minify app main js file contents: %v", err)
	}

	AppMainCSS, err := f.ReadFile("frontend/static/assets/css/app/main.css")
	if err != nil {
		log.Fatalf("Failed to read app main css file: %v", err)
	}
	AppMainCSSMinified, err := m.Bytes("text/css", AppMainCSS)
	if err != nil {
		log.Fatalf("Failed to minify app main css file contents: %v", err)
	}

	AppChunkJS, err := f.ReadFile("frontend/static/assets/js/app/chunk.js")
	if err != nil {
		log.Fatalf("Failed to read app chunk js file: %v", err)
	}
	AppChunkJSMinified, err := m.Bytes("text/javascript", AppChunkJS)
	if err != nil {
		log.Fatalf("Failed to minify app chunk js file contents: %v", err)
	}

	grpcUiOpts := []standalone.HandlerOption{
		standalone.WithIndexTemplate(indexTemplate),
		standalone.ServeAsset("favicon.png", favicon),
		standalone.ServeAsset("jquery.min.js", jqueryJSMinified),
		standalone.ServeAsset("jquery-ui.min.js", jqueryUiJSMinified),
		standalone.ServeAsset("jquery-ui.min.css", jqueryUiCSSMinified),
		standalone.ServeAsset("grpc-web-form.js", grpcWebFormJSMinified),
		standalone.ServeAsset("grpc-web-form.css", grpcWebFormCSSMinified),
		standalone.ServeAsset("app-main.js", AppMainJSMinified),
		standalone.ServeAsset("app-main.css", AppMainCSSMinified),
		standalone.ServeAsset("app-chunk.js", AppChunkJSMinified),
	}

	h, err := standalone.HandlerViaReflection(context.Background(), conn, target, grpcUiOpts...)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Xss-Protection", "1; mode=block")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")

		if r.URL.Path != "/" {
			h.ServeHTTP(w, r)
			return
		}

		fw := &fakeResponseWriter{ResponseWriter: w}

		if len(executedIndexTemplateMinified) == 0 {
			fw.buf = new(bytes.Buffer)
		}

		h.ServeHTTP(fw, r)

		if len(executedIndexTemplateMinified) == 0 {
			executedIndexTemplate, err := io.ReadAll(fw.buf)
			if err != nil {
				log.Printf("Failed to read executed index template: %v", err)
				io.Copy(w, fw.buf)
				return
			}

			executedIndexTemplateMinified, err = m.Bytes("text/html", executedIndexTemplate)
			if err != nil {
				log.Printf("Failed to minify executed index template: %v", err)
				io.Copy(w, fw.buf)
				return
			}
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(executedIndexTemplateMinified)))
		w.Write(executedIndexTemplateMinified)
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/", wrappedHandler)
	serveMux.Handle("/favicon.ico", http.RedirectHandler("/s/favicon.png", 301))
	serveMux.Handle("/favicon.png", http.RedirectHandler("/s/favicon.png", 301))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", webport))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening at %v", lis.Addr())

	if err := http.ServeTLS(lis, serveMux, certFile, keyFile); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
