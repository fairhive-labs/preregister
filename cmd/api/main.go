package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fairhive-labs/preregister/internal/crypto"
	"github.com/fairhive-labs/preregister/internal/crypto/cipher"
	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/fairhive-labs/preregister/internal/limiter"
	"github.com/fairhive-labs/preregister/internal/mailer"
)

type App struct {
	db     data.DB
	jwt    crypto.Token
	mailer mailer.Mailer
	wg     sync.WaitGroup
	rl     *limiter.RateLimiter
}

const (
	tableName = "Users"
)

var (
	jwts = map[string]crypto.Token{}
	ek   string
)

func init() {
	k, _ := cipher.GenerateKey(32)
	jwts["HS512"] = crypto.NewJWTHS512(k)
	k, _ = cipher.GenerateKey(16)
	jwts["HS256"] = crypto.NewJWTHS256(k)
	jwts["ES256"], _ = crypto.NewJWTES256()
	jwts["ES512"], _ = crypto.NewJWTES512()
	log.Println("🔐 JWT Services: OK")

	ek = os.Getenv("FAIRHIVE_ENCRYPTION_KEY")
	if ek == "" {
		panic("encryption key is missing")
	}
	log.Println("🔑 Encryption Key: OK")
}

func NewApp() *App {
	db, err := data.NewDynamoDB(tableName, ek)
	if err != nil {
		panic(err)
	}
	return &App{
		db:  db,
		jwt: jwts["ES256"],
		mailer: mailer.New(os.Getenv("FAIRHIVE_GSUITE_USER"),
			os.Getenv("FAIRHIVE_GSUITE_PASSWORD"),
			"smtp.gmail.com",
			587),
		rl: limiter.New(0.1, 1),
	}
}

var jwtregexp = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)

func (app App) register(c *gin.Context) {
	var u data.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := app.jwt.Create(&u, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash := app.jwt.Hash(token)
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.mailer.SendActivationEmail(u.Email, fmt.Sprintf("http://fairhive.io/activate/%s", token), hash)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"hash": hash,
	})
}

func (app App) activate(c *gin.Context) {
	t := c.Param("token")
	h := c.Param("hash")
	if !jwtregexp.MatchString(t) || app.jwt.Hash(t) != h {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	u, err := app.jwt.Extract(t) // verify + extract
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = app.db.Save(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		app.mailer.SendConfirmationEmail(u.Email)
	}()

	c.JSON(http.StatusCreated, gin.H{
		"token":     t,
		"activated": true,
	})
}

func (app App) limit(c *gin.Context) {
	ip := c.ClientIP()
	l := app.rl.GetAccess(ip)
	if !l.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "Too Many Requests",
			"ip":    ip,
		})
		return
	}
	c.Next()
}

func setupRouter(app App) *gin.Engine {
	r := gin.Default()
	r.Use(app.limit)
	r.POST("/", app.register)
	r.POST("/activate/:token/:hash", app.activate)
	return r
}

func main() {
	app := *NewApp()
	// gin.SetMode(gin.ReleaseMode)
	r := setupRouter(app)

	var addr string
	if p := os.Getenv("PORT"); p != "" {
		addr = "" + p
	} else {
		addr = ":8080" // default port
	}

	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-quit
		log.Printf("🚨 Shutdown signal \"%v\" received\n", s)

		log.Printf("🚦 Here we go for a graceful Shutdown...\n")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("⚠️ HTTP server Shutdown: %v", err)
		}

		log.Printf("⏳ Waiting the end of all go-routines...")
		app.wg.Wait() // wait for all go-routines
		log.Printf("👍 go-routines are over")
		close(idleConnsClosed)
	}()

	log.Printf("✅ Listening and serving HTTP on %s\n", addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("👹 HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Printf("😴 Server stopped")
}
