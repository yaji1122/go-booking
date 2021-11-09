package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/model"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(model.Reservation{})

	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 30 * time.Minute
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //localhost use http, in product will be true

	testApp.Session = session

	appConfig = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	var header http.Header
	return header
}

func (mw *myWriter) WriteHeader(i int) {

}

func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
