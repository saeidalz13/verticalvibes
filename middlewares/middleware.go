package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/saeidalz13/verticalvibes/token"
)

type Middleware struct {
	logger       *log.Logger
	tokenManager *token.PasetoTokenManager
}

func NewMiddlewareHanlder(logger *log.Logger, tokenManager *token.PasetoTokenManager) *Middleware {
	mwLogger := log.New(log.Writer(), "MIDDLEWARE: ", logger.Flags())
	return &Middleware{
		mwLogger,
		tokenManager,
	}
}

func (m *Middleware) authorizeUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := r.Cookie("auth-token")
		if err != nil {
			m.logger.Println(err)
			http.Error(w, "Forbidden access", http.StatusForbidden)
			return
		}

		// pasetoPayload
		pp, err := m.tokenManager.ValidateToken(authToken.Value)
		if err != nil {
			m.logger.Println(err)
			http.Error(w, "Forbidden access", http.StatusForbidden)
			return
		}

        m.logger.Printf("%+v\n", pp)

		handler.ServeHTTP(w, r)
	})
}

func (m *Middleware) logRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.logger.Printf("new request -> Path: %s - RemoteAddr: %s - Host: %s", r.URL.Path, r.RemoteAddr, r.Host)
		handler.ServeHTTP(w, r)
		m.logger.Printf("Processed in %d ms", time.Since(start).Milliseconds())
	})
}

func (m *Middleware) Chain(handler http.Handler) http.Handler {
	return m.logRequests(m.authorizeUser(handler))
}
