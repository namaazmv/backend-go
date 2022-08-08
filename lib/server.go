package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router

	Prayer Prayer
}

type PrayerJson struct {
	CategoryId int `json:"category_id"`
	Date       int `json:"date"`
	Fajr       int `json:"fajr"`
	Sunrise    int `json:"sunrise"`
	Duhr       int `json:"dhuhr"`
	Asr        int `json:"asr"`
	Maghrib    int `json:"maghrib"`
	Isha       int `json:"isha"`
}

type IslandJson struct {
	CategoryId  int     `json:"category_id"`
	IslandId    int     `json:"island_id"`
	Atoll       int     `json:"atoll"`
	EnglishName string  `json:"english_name"`
	DhivehiName string  `json:"dhivehi_name"`
	ArabicName  string  `json:"arabic_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Status      int8    `json:"status"`
}

type TodayJson struct {
	Island      IslandJson `json:"island"`
	PrayerTimes PrayerJson `json:"prayer_times"`
}

type NextJson struct {
	Call            string `json:"call"`
	Timestamp       int    `json:"timestamp"`
	TimestampString string `json:"timestamp_string"`
}

func NewServer() *Server {
	atolls := ParseCSV[Atoll]("atolls")
	islands := ParseCSV[Island]("islands")
	prayerTimes := ParseCSV[PrayerTime]("prayertimes")

	s := &Server{Router: mux.NewRouter(), Prayer: Prayer{
		Atolls:      atolls,
		Islands:     islands,
		PrayerTimes: prayerTimes,
		Timings:     []string{"fajr", "sunrise", "duhr", "asr", "maghrib", "isha"},
	}}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.HandleFunc("/", s.HelloWorld()).Methods("GET")
	s.HandleFunc("/today", s.GetToday()).Methods("GET")
	s.HandleFunc("/next", s.GetNext()).Methods("GET")
}

func (s *Server) HelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}
}

func (s *Server) GetToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if query, err := strconv.Atoi(r.URL.Query().Get("island")); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to parse the integer."))
			return
		} else {
			island := s.Prayer.GetIsland(query)
			if island == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid island provided."))
				return
			}

			today := s.Prayer.GetToday(island)
			data, err := json.Marshal(TodayJson{
				Island: IslandJson{
					CategoryId:  island.CategoryId,
					IslandId:    island.IslandId,
					Atoll:       island.Atoll,
					EnglishName: island.EnglishName,
					DhivehiName: island.DhivehiName,
					ArabicName:  island.ArabicName,
					Latitude:    island.Longitude,
					Longitude:   island.Longitude,
					Status:      island.Status,
				},
				PrayerTimes: PrayerJson{
					CategoryId: today.CategoryId,
					Date:       today.Date,
					Fajr:       today.Fajr,
					Sunrise:    today.Sunrise,
					Duhr:       today.Duhr,
					Asr:        today.Asr,
					Maghrib:    today.Maghrib,
					Isha:       today.Isha,
				},
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to marshal the data."))
				return
			}

			w.Write(data)
		}
	}
}

func (s *Server) GetNext() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if query, err := strconv.Atoi(r.URL.Query().Get("island")); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to parse the integer."))
			return
		} else {
			island := s.Prayer.GetIsland(query)
			if island == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid island provided."))
				return
			}

			prayerToday := s.Prayer.GetToday(island)

			now := time.Now().Local()
			var call *string = nil

			for i := range s.Prayer.Timings {
				timing := s.Prayer.Timings[i]

				if ConvertTimestampToDate(prayerToday.GetValue(timing)).After(now) {
					call = &timing
					break
				}

				fmt.Println(ConvertTimestampToDate(prayerToday.GetValue(timing)))
			}

			if call == nil {
				new_prayer := "fajr"
				call = &new_prayer
				prayerToday = s.Prayer.GetEntryFromDay((DaysIntoYear(now) + 1), island)
			}

			def_call := *call
			timestamp := prayerToday.GetValue(def_call)

			data, err := json.Marshal(NextJson{
				Call:            def_call,
				Timestamp:       timestamp,
				TimestampString: ConvertTimestampToString(timestamp),
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to marshal the data."))
				return
			}

			w.Write(data)
		}
	}
}

// Todos

// func logRequestHandler(h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {

// 		h.ServeHTTP(w, r)

// 		uri := r.URL.String()
// 		method := r.Method

// 		fmt.Println(uri, method)
// 	}

// 	return http.HandlerFunc(fn)
// }
