package api

import (
  "fmt"
  "strings"
  "time"
  "math/rand"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
)

type HttpResponse struct {
  Status string   `json:"status"`
  Message string  `json:"message"`
}

func init() {
  r := mux.NewRouter()
  r.HandleFunc("/api/app", auth(PostAppHandler)).Methods("POST")
  r.HandleFunc("/api/apps", auth(GetAppsHandler)).Methods("GET")
  r.HandleFunc("/api/auth", auth(GetApiKeyHandler)).Methods("GET")
  http.Handle("/", r)
}

/**
 * Register a new app or update an existing one.
 *  {
 *    "name" : "resonance-srv",
 *    "android_package" : "com.atooma.resonance.sdk",
 *    "ios_bundle" : ""
 *  }
 *  Authorization: a.petreri@atooma.com
 */
func PostAppHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  decoder := json.NewDecoder(r.Body)
  var app App
  err := decoder.Decode(&app)
  if err == nil {
    app.Owner = owner
    app.Id = hash(app.Name)
    success, err := PostApp(r, app)
    if success {
      responseHandler(w, app)
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
  }
}

/**
 *  Authorization: a.petreri@atooma.com
 */
func GetAppsHandler(w http.ResponseWriter, r *http.Request) {
  apps := GetApps(r)
  responseHandler(w, apps)
}

/**
 *  Authorization: com.atooma.resonance.sdk:ed51f97e73b10b974f765ecfad1579dd
 */
func GetApiKeyHandler(w http.ResponseWriter, r *http.Request) {
  sender := strings.Split(r.Header.Get("Authorization"), ":")
  androidPackage, appId, _ := sender[0], sender[1], sender[2]
  app := GetApp(r, appId)
  if app.Android == androidPackage {
    key := ApiKey{Key:random(10)}
    //app.Keys = append(app.Keys, key)
    //PostApp(r, *app)
    responseHandler(w, key)
  } else {
    errorHandler(w, r, http.StatusUnauthorized, "package not allowed")
  }
}

func random(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    owner := r.Header.Get("Authorization")
    if owner == "" {
      errorHandler(w, r, http.StatusUnauthorized, "unauthorized")
    } else {
      fn(w, r)
    }
  }
}

func responseHandler(w http.ResponseWriter, v interface{}) {
  w.Header().Set("Content-Type", "application/json")
  encoder := json.NewEncoder(w)
  encoder.Encode(v)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  response := HttpResponse{Status:"failure",Message:message}
  encoder := json.NewEncoder(w)
  encoder.Encode(response)
}




/*
func GetCalendarsHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  calendars := GetCalendars(r, owner)
  for _, calendar := range calendars {
    if calendar.Events == nil {
      calendar.Events = []Event{}
    }
  }
  responseHandler(w, calendars)
}

func GetCalendarHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  vars := mux.Vars(r)
  calendar := GetCalendar(r, owner, vars["id"])
  if calendar != nil {
    if calendar.Events == nil {
      calendar.Events = []Event{}
    }
    responseHandler(w, calendar)
  } else {
    errorHandler(w, r, http.StatusNotFound, "not found")
  }
}

func PostCalendarHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  decoder := json.NewDecoder(r.Body)
  var calendar Calendar
  err := decoder.Decode(&calendar)
  if err == nil {
    calendar.Owner = owner
    if calendar.Events == nil {
      calendar.Events = []Event{}
    }
    success, err := PostCalendar(r, calendar)
    if success {
      responseHandler(w, calendar)
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
  }
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
  owner := r.Header.Get("Authorization")
  vars := mux.Vars(r)
  decoder := json.NewDecoder(r.Body)
  var event Event
  err := decoder.Decode(&event)
  if err == nil {
    success, err := PostEvent(r, vars["calendarId"], owner, event)
    if success {
      responseHandler(w, event)
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
  }
}
*/
