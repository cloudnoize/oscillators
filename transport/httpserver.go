package transport

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudnoize/oscillators/oscillators"
	"github.com/gorilla/mux"
)

// NewHTTPHandler creates a new HTTP handler that serves the Sinkhole service
func NewHTTPHandler(contexts oscillators.ClientContexts, frq oscillators.Freq) http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = CreateDefaultHandler()
	r.HandleFunc("/register/{osc}", CreateRegisterationHandler(contexts))
	r.HandleFunc("/token/{token}/note/{note}", CreateChangeNoteHandler(contexts, frq))
	return r
}

func CreateDefaultHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Cache-Control", "max-age=7200") // 2hour

	}
}

func CreateRegisterationHandler(contexts oscillators.ClientContexts) http.HandlerFunc {
	rnd := rand.New(rand.NewSource(99))
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cl := rnd.Int() % 10000
		osc := vars["osc"]
		osc = strings.ToLower(osc)
		iosc := oscillators.Sin
		switch osc {
		case "sin":
			iosc = oscillators.Sin
		case "test":
			iosc = oscillators.Test
		default:
			fmt.Fprintf(w, "%s osc not supported", osc)
			return
		}
		println("Registering oscillator ", osc)
		contexts[uint(cl)] = &oscillators.ClientContext{Start: time.Now(), Misc: "Hi synth fan", Osc: oscillators.Oscillators(iosc)}
		w.Header().Set("registration", strconv.Itoa(cl)) //"registration"] = strconv.Itoa(rnd.Uint32())
		w.Header().Set("Cache-Control", "max-age=7200")  // 2hour

	}
}

func CreateChangeNoteHandler(contexts oscillators.ClientContexts, frq oscillators.Freq) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tk := vars["token"]
		note := vars["note"]
		println("Changin note ", note, " for client ", tk)
		icl, e := strconv.Atoi(tk)
		println("got ", icl)
		if e != nil {
			println(e.Error())
			return
		}
		cc, ok := contexts[uint(icl)]
		if cc == nil || !ok {
			println("got nil for cl ", icl)
			return
		}
		cc.SetFreq(frq.ToFreq(strings.ToLower(note)))
	}
}
