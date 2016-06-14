package main
import (
  "net/http"
)
type authHandler struct {
  next http.Handler
}
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
  if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
    // No Auth
    w.Header().Set("Location", "/login")
    w.WriteHeader(http.StatusTemporaryRedirect)
  }else if err != nil {
    // Something error is occured.
    panic(err.Error())
  }else{
    // Success. Call a handler that is wrapped.
    h.next.ServeHTTP(w, r)
  }
}
func MustAuth(handler http.Handler) http.Handler  {
  return &authHandler{next: handler}
}
