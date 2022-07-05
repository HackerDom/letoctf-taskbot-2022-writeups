package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx"
	"github.com/minio/minio-go"
	"github.com/rs/zerolog"

	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/api"
	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/storage"
)

//go:embed files/flag.enc
var flagFile []byte

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	minioEndpoint := "localhost:9000"
	dbEndpoint := "localhost"
	if os.Getenv("INSIDE_DOCKER") == "1" {
		minioEndpoint = "minio:9000"
		dbEndpoint = "postgres"
	}

	logger.Info().Str("endpoint", minioEndpoint).Msg("Connecting to minio...")
	client, err := minio.New(
		minioEndpoint, "minioadmin", "minioadmin", false,
	)
	if err != nil {
		logger.Fatal().Err(err).Str("endpoint", minioEndpoint).Msg("connect to minio failed")
	}
	filestore, err := storage.NewFileStorage(client)
	if err != nil {
		logger.Fatal().Err(err).Msg("create file storage failed")
	}

	conf := pgx.ConnConfig{
		Host:     dbEndpoint,
		Port:     5432,
		Database: "filestore",
		User:     "filestore",
		Password: "filestore",
	}
	connectionUrl := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database,
	)
	logger.Info().Str("url", connectionUrl).Msg("Connecting to postgres...")
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: 20,
	})
	if err != nil {
		logger.Fatal().Err(err).
			Str("connection config", fmt.Sprint(conf)).
			Msg("connect to database failed: %v")
	}
	dbStore, err := storage.NewDbStorage(connPool)
	if err != nil {
		logger.Fatal().Err(err).Msg("create database storage failed")
	}

	sessionKey := []byte(getEnv(&logger, "SESSION_KEY"))
	sessionStore := sessions.NewCookieStore(sessionKey)

	serverKey := []byte(getEnv(&logger, "SERVER_KEY"))
	wasmPath := getEnv(&logger, "WASM_PATH")
	webuiPath := getEnv(&logger, "WEBUI_PATH")

	addFlag(&logger, serverKey, dbStore, filestore)

	opts := &api.ServerOptions{
		Filestore:    filestore,
		DbStore:      dbStore,
		SessionStore: sessionStore,
		HmacKey:      serverKey,
	}
	s := api.NewServer(&logger, http.Dir(wasmPath), http.Dir(webuiPath), opts)

	addr := ":8080"
	logger.Info().Str("addr", addr).Msg("Listening...")
	if err := http.ListenAndServe(addr, route(s)); err != nil {
		logger.Err(err).Str("address", addr).Msg("listen and serve failed")
	}
}

func route(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/api") {
			http.StripPrefix("/api", h).ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.RequestURI, "/static") {
			h.ServeHTTP(w, r)
			return
		}

		// если запрос не имеет префикса api, необходимо вернуть SPA приложение
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "/"
		r2.URL.RawPath = "/"
		h.ServeHTTP(w, r2)
	})
}

func addFlag(
	logger *zerolog.Logger,
	hmacKey []byte,
	dbStore storage.DbStorage,
	filestore storage.FileStorage,
) {
	adminId := uuid.MustParse("9c97249a-26d8-49f2-8f1b-92641a8901f8")
	adminUsername := "admin"
	adminPass := []byte(getEnv(logger, "ADMIN_PASS"))

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(adminPass)
	passHash := mac.Sum(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// проверка существования пользователя
	if _, err := dbStore.GetUserId(ctx, adminUsername, passHash); err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := dbStore.CreateUser(ctx, "admin", passHash, adminId); err != nil {
			logger.Fatal().Err(err).Msg("create admin user failed")
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		reader := bytes.NewReader(flagFile)
		opts := &storage.FileOptions{
			Name:      "flag.enc",
			Size:      int64(len(flagFile)),
			OwnerId:   adminId,
			Encrypted: true,
		}
		if err := filestore.Update(ctx, reader, opts); err != nil {
			logger.Fatal().Err(err).Msg("upload file with flag failed")
		}
	}
}

func getEnv(logger *zerolog.Logger, name string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		logger.Fatal().Str("name", name).Msg("variable is not set")
	}

	return v
}
