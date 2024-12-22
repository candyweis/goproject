package session

import "github.com/gorilla/sessions"

// Store — хранилище сессий (CookieStore).
var Store = sessions.NewCookieStore([]byte("my-secret-key-123"))
