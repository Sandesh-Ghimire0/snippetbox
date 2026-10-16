package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Sandesh-Ghimire0/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; model=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// method against app to get access of infolog
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// returns the internal server error response during panic
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// this function will execute when unwinding the stack during panic
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// redirects the unauthenticated user to /login
func (app *application) requiresAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			// return from the middleware chain so that no subsequent handlers in the chain are executed.
			return
		}

		// prevent browser from caching the content so that when an unauthenticated call is made
		// broswer donot reach out to cache to get the content where there is not authentication check
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// adds the isAutheticated key in the current request cycle
// which can be used to authenticate the user down the middlware chain and handler
// so that there is no need to make DB calls each time we want to authenticate the user
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r.Context(), "authenticatedUserID")

		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		userId := app.session.GetInt(r.Context(), "authenticatedUserID")
		user, err := app.users.Get(userId)
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r.Context(), "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.ServerError(w, err)
			return
		}

		// request is coming from the valid user
		// add the isAuthenticated key in the context with value true
		// pass the request to next handler with the context
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx)) // r.withContext() return the new copy of the request with context value

		// now use the context value down the middleware chain or handlers for this request cycle
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
