package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env ...
type Env struct {
	DB  *sqlx.DB
	Log *logrus.Entry
}

// Index ...
// GET /
func (e *Env) Index(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(`{"result": "OK"}`))
	return nil
}
