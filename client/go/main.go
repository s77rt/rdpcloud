package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fullstorydev/grpcui/standalone"
)

var (
	Version string = "dev"

	DisplayName string = "RDPCloud"
)

//go:embed frontend
var f embed.FS

func main() {
	log.Printf("Running RDPCloud Client (Version: %s)", Version)

	addr := "51.89.161.169:5027"

	// target is used for display purposes
	// it's meant to be for the addr, but we have no use in displaying the addr
	// instead we use this variable to display a custom app name
	target := DisplayName

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(addr, grpcOpts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	index, err := f.ReadFile("frontend/templates/index.tmpl")
	if err != nil {
		log.Fatalf("Failed to read index tmpl file: %v", err)
	}
	indexTemplate := template.Must(template.New("index").Parse(string(index)))

	jqueryJS, err := f.ReadFile("frontend/static/assets/js/jquery.min.js")
	if err != nil {
		log.Fatalf("Failed to read jquery js file: %v", err)
	}
	jqueryUiJS, err := f.ReadFile("frontend/static/assets/js/jquery-ui.min.js")
	if err != nil {
		log.Fatalf("Failed to read jquery-ui js file: %v", err)
	}
	grpcWebFormJS, err := f.ReadFile("frontend/static/assets/js/grpc-web-form.js")
	if err != nil {
		log.Fatalf("Failed to read grpc-web-form js file: %v", err)
	}
	jqueryUiCSS, err := f.ReadFile("frontend/static/assets/css/jquery-ui.min.css")
	if err != nil {
		log.Fatalf("Failed to read jquery-ui css file: %v", err)
	}
	grpcWebFormCSS, err := f.ReadFile("frontend/static/assets/css/grpc-web-form.css")
	if err != nil {
		log.Fatalf("Failed to read grpc-web-form css file: %v", err)
	}
	AppMainJS, err := f.ReadFile("frontend/static/assets/js/app/main.js")
	if err != nil {
		log.Fatalf("Failed to read app main js file: %v", err)
	}
	AppMainCSS, err := f.ReadFile("frontend/static/assets/css/app/main.css")
	if err != nil {
		log.Fatalf("Failed to read app main css file: %v", err)
	}
	AppChunkJS, err := f.ReadFile("frontend/static/assets/js/app/chunk.js")
	if err != nil {
		log.Fatalf("Failed to read app chunk js file: %v", err)
	}

	grpcUiOpts := []standalone.HandlerOption{
		standalone.WithIndexTemplate(indexTemplate),
		standalone.ServeAsset("jquery.min.js", jqueryJS),
		standalone.ServeAsset("jquery-ui.min.js", jqueryUiJS),
		standalone.ServeAsset("grpc-web-form.js", grpcWebFormJS),
		standalone.ServeAsset("jquery-ui.min.css", jqueryUiCSS),
		standalone.ServeAsset("grpc-web-form.css", grpcWebFormCSS),
		standalone.ServeAsset("app-main.js", AppMainJS),
		standalone.ServeAsset("app-main.css", AppMainCSS),
		standalone.ServeAsset("app-chunk.js", AppChunkJS),
	}

	h, err := standalone.HandlerViaReflection(context.Background(), conn, target, grpcUiOpts...)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/", h)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening at %v", lis.Addr())

	if err := http.Serve(lis, serveMux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
