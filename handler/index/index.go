package index

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

const tpl = `<!DOCTYPE html>
<html><head>
  <meta charset="UTF-8">
</head><body>
{{if not .IsAdmin}}
  <a href="{{.Login}}">Sign in or register</a><br />
{{else}}
  Welcome, {{.User}}! (<a href="{{.Logout}}">sign out</a>)<br />
  <br />
  <a href="rate">Help Organize Photos!!!</a><br />
  <a href="best">See the best photos!!!</a><br />
{{end}}
</body></html>`

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	t, _ := template.New("webpage").Parse(tpl)
	login, _ := user.LoginURL(ctx, "/")
	logout, _ := user.LogoutURL(ctx, "/")

	data := struct {
		IsAdmin bool
		User    *user.User
		Login   string
		Logout  string
	}{
		IsAdmin: user.IsAdmin(ctx),
		User:    user.Current(ctx),
		Login:   login,
		Logout:  logout,
	}

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	t.Execute(w, data)
}
