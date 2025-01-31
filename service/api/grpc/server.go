package grpc

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"time"

	golang "github.com/fatemehkarimi/link-shortener/api-scheme/proto/src/golang"
	"github.com/fatemehkarimi/link-shortener/service/db"
	"github.com/fatemehkarimi/link-shortener/service/internal"
	"google.golang.org/grpc"
)

type Server struct {
	DB *sql.DB
	golang.UnimplementedLinkShortenerServer
}

var now = time.Now

func (s Server) CreateLink(context context.Context, request *golang.RequestCreateLink) (*golang.ResponseCreateLink, error) {
	linkHash := internal.GenerateLinkHash(now, request.URL)
	create_at := now()
	expires_at := create_at.AddDate(0, 0, 7)

	err := db.InsertIntoDatabase(s.DB, request.URL, linkHash, create_at, expires_at)
	if err != nil {
		return nil, err
	}

	return &golang.ResponseCreateLink{Hash: linkHash, CreateAt: create_at.UnixMilli(), ExpiresAt: expires_at.UnixMilli()}, nil
}

func (s Server) GetLinkByHash(context context.Context, request *golang.RequestGetLinkByHash) (*golang.ResponseGetLinkByHash, error) {
	originalURL, err := db.GetLinkByHashFromDb(s.DB, request.Hash)
	if err != nil {
		return nil, err
	}

	if originalURL != nil {
		return &golang.ResponseGetLinkByHash{URL: *originalURL}, nil
	}

	return nil, errors.New("url not found")
}

func (server Server) Start() {
	dbCredentials := db.GetDBCredentials()
	db, err := sql.Open("postgres", dbCredentials)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	server.DB = db

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	}

	serverRegistary := grpc.NewServer()

	golang.RegisterLinkShortenerServer(serverRegistary, server)
	err = serverRegistary.Serve(listener)

	if err != nil {
		log.Fatalf("impossible to serve %s", err)
	}
}

func NewServer() Server {
	return Server{}
}
