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

//go:embed frontend
var f embed.FS

func main() {
	addr := "51.89.161.169:5027"

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	index, err := f.ReadFile("frontend/templates/index.tmpl")
	if err != nil {
		log.Fatalf("failed to read index tmpl file: %v", err)
	}
	indexTemplate := template.Must(template.New("index").Parse(string(index)))

	jqueryJS, err := f.ReadFile("frontend/static/assets/js/jquery.min.js")
	if err != nil {
		log.Fatalf("failed to read jquery js file: %v", err)
	}
	jqueryUiJS, err := f.ReadFile("frontend/static/assets/js/jquery-ui.min.js")
	if err != nil {
		log.Fatalf("failed to read jquery-ui js file: %v", err)
	}
	jqueryUiCSS, err := f.ReadFile("frontend/static/assets/css/jquery-ui.min.css")
	if err != nil {
		log.Fatalf("failed to read jquery-ui css file: %v", err)
	}
	AppMainJS, err := f.ReadFile("frontend/static/assets/js/app/main.js")
	if err != nil {
		log.Fatalf("failed to read app main js file: %v", err)
	}
	AppMainCSS, err := f.ReadFile("frontend/static/assets/css/app/main.css")
	if err != nil {
		log.Fatalf("failed to read app main css file: %v", err)
	}
	AppChunkJS, err := f.ReadFile("frontend/static/assets/js/app/chunk.js")
	if err != nil {
		log.Fatalf("failed to read app chunk js file: %v", err)
	}

	opts := []standalone.HandlerOption{
		standalone.WithIndexTemplate(indexTemplate),
		standalone.ServeAsset("jquery.min.js", jqueryJS),
		standalone.ServeAsset("jquery-ui.min.js", jqueryUiJS),
		standalone.ServeAsset("jquery-ui.min.css", jqueryUiCSS),
		standalone.ServeAsset("app-main.js", AppMainJS),
		standalone.ServeAsset("app-main.css", AppMainCSS),
		standalone.ServeAsset("app-chunk.js", AppChunkJS),
	}

	h, err := standalone.HandlerViaReflection(context.Background(), conn, addr, opts...)
	if err != nil {
		log.Fatalf("failed to create handle: %v", err)
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
