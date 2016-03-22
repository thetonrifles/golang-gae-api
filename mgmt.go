package api

import (
  "fmt"
  "errors"
  "net/http"
  "crypto/md5"
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
)

type App struct {
  Id string           `json:"id"`
  Name string         `json:"name"`
  Owner string        `json:"owner"`
  Android string      `json:"android_package"`
  IOS string          `json:"ios_bundle"`
}

type Device struct {
  Id string           `json:"id"`
  Model string        `json:"model"`
  Vendor string       `json:"vendor"`
  OS string           `json:"os"`
  ApiVersion string   `json:"api_version"`
  Keys []ApiKey       `json:"keys"`
}

type ApiKey struct {
  Key string          `json:"key"`
  AppId string        `json:"app_id"`
}

/**
 *  Create a new app in Google Datastore.
 */
func PostApp(r *http.Request, app App) (bool, error) {
  context := appengine.NewContext(r)
  key := hash(app.Id)
  appKey := datastore.NewKey(context, "application", key, 0, nil)
  err := datastore.Get(context, appKey, &app)
  if err != nil {
    _, err := datastore.Put(context, appKey, &app)
    if err != nil {
      return false, err
    } else {
      return true, nil
    }
  } else {
    return false, errors.New("app already exists")
  }
}

/**
 *  Get an app from Google Datastore.
 */
func GetApp(r *http.Request, appId string) *App {
  context := appengine.NewContext(r)
  key := hash(appId)
  appKey := datastore.NewKey(context, "application", key, 0, nil)
  var app App
  datastore.Get(context, appKey, &app)
	return &app
}

/**
 *  Get all apps from Google Datastore.
 */
func GetApps(r *http.Request) []*App {
  context := appengine.NewContext(r)
  q := datastore.NewQuery("application")
  var apps []*App
  _, err := q.GetAll(context, &apps)
  if err==nil && apps!=nil {
    return apps
  } else {
    return []*App{}
  }
}

func hash(s string) string {
  data := []byte(s)
  return fmt.Sprintf("%x", md5.Sum(data))
}
