package source

import (
	"sync"
	"net/http"
	"strconv"
	"io/ioutil"
	"fmt"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
)

type Http struct {
	Base
}

func (s *Http) Listen(wg sync.WaitGroup) {

	s.AppLog.Info("Start listening (http:%d)", s.Port)

	m := martini.New()
	route := martini.NewRouter()

	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		// Use indentations. &pretty=1
		pretty, _ := strconv.ParseBool(r.FormValue("pretty"))
		// Use null instead of empty object for json &null=1
		null, _ := strconv.ParseBool(r.FormValue("null"))
		// Some content negotiation
		switch r.Header.Get("Accept") {
		case "application/xml":
			c.MapTo(encoder.XmlEncoder{PrettyPrint: pretty}, (*encoder.Encoder)(nil))
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		default:
			c.MapTo(encoder.JsonEncoder{PrettyPrint: pretty, PrintNull: null}, (*encoder.Encoder)(nil))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
	})

	route.Post("/event", func(enc encoder.Encoder, w http.ResponseWriter, r *http.Request) (int, []byte) {

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.AppLog.Debug("%s", err)
			return http.StatusInternalServerError, []byte{}
		}

		result := s.processRaw(string(body))
		return http.StatusOK, encoder.Must(enc.Encode(result))
	})

	m.Action(route.Handle)

	port := fmt.Sprintf(":%d", s.Port)
	if err := http.ListenAndServe(port, m); err != nil {
		s.AppLog.Fatal(err)
	}

	wg.Done()
}
