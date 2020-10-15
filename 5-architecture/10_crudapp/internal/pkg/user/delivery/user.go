package delivery

import (
	"crudapp/internal/pkg/models"
	"crudapp/internal/pkg/user"
	"errors"
	"html/template"
	"net/http"

	"crudapp/internal/pkg/session"

	"go.uber.org/zap"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo user.Repository
	Sessions *session.SessionsManager
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := models.SessionFromContext(r.Context())
	if err == nil {
		http.Redirect(w, r, "/items", 302)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Authorize(r.FormValue("login"), r.FormValue("password"))
	if err == models.ErrNoUser {
		http.Error(w, `no user`, http.StatusBadRequest)
		return
	}
	if err == models.ErrBadPass {
		http.Error(w, `bad pass`, http.StatusBadRequest)
		return
	}

	sess, _ := h.Sessions.Create(w, u.ID)
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", 302)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/", 302)
}

func (h *UserHandler) GetUserByID(login string) (*models.User, error) {
	if login == "" {
		return nil, errors.New("login is empty")
	}

	return h.UserRepo.GetByLogin(login)
}
