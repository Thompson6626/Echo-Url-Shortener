package main

import (
	"Url-Shortener/internal/auth"
	"Url-Shortener/internal/database"
	"Url-Shortener/internal/env"
	"Url-Shortener/internal/ratelimiter"
	"Url-Shortener/internal/store"
	"Url-Shortener/internal/store/cache"
	"context"
	"expvar"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "1.1.0"

func main() {
	cfg := config{
		addr:   env.GetString("PORT", ":8080"),
		apiURL: env.GetString("API_URL", "localhost:8080"),
		db: dbConfig{
			host:     env.GetString("DB_HOST", "localhost"),
			port:     env.GetString("DB_PORT", "27017"),
			dbName:   env.GetString("DB_NAME", "shortener_db"),
			username: env.GetString("DB_USERNAME", "admin"),
			password: env.GetString("DB_ROOT_PASSWORD", "adminpass"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gophersocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := database.New(
		cfg.db.host,
		cfg.db.port,
		cfg.db.dbName,
		cfg.db.username,
		cfg.db.password,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Info("database connection pool established")

	// Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	storage := store.NewStorage(db.Database(cfg.db.dbName))
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:        cfg,
		store:         storage,
		cacheStorage:  cacheStorage,
		logger:        logger,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("mongo_status", expvar.Func(func() any {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx, nil); err != nil {
			return "disconnected"
		}
		return "connected"
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
