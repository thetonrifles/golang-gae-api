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
  r.HandleFunc("/api/device", auth(PostDeviceHandler)).Methods("POST")
  r.HandleFunc("/api/auth", auth(GetApiKeyHandler)).Methods("GET")
  http.Handle("/", r)
}

/**
 * Register a new app.
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
    success, err := PostApp(r, &app)
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
    key := ApiKey{Key:random(10),AppId:appId}
    //app.Keys = append(app.Keys, key)
    //PostApp(r, *app)
    responseHandler(w, key)
  } else {
    errorHandler(w, r, http.StatusUnauthorized, "package not allowed")
  }
}

/**
 * Register a new device or update an existing one.
 * {
 *   "device_id" : "abcdefghijklmnopqrstuvwxyz0123456789"
 *   "model" : "Nexus"
 *   "vendor" : "LG"
 *   "platform" : "Android"
 *   "platform_version" : "23"
 *  }
 *  Authorization: a.petreri@atooma.com
 */
func PostDeviceHandler(w http.ResponseWriter, r *http.Request) {
  sender := strings.Split(r.Header.Get("Authorization"), ":")
  //apiVersion := sender[0]
  //deviceId := sender[1]
  //lang := sender[2]
  appId := sender[3]
  //pack := sender[4]
  decoder := json.NewDecoder(r.Body)
  var device Device
  err := decoder.Decode(&device)
  if err == nil {
    success, err := PostDevice(r, &device)
    if success {
      key := ApiKey{Key:random(10),AppId:appId}
      responseHandler(w, key)
    } else {
      errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
    }
  } else {
    errorHandler(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", err))
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
