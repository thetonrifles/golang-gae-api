package api

import (
  "fmt"
  "time"
  "errors"
  "net/http"
  "math/rand"
  "crypto/md5"
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
  "google.golang.org/appengine/log"
)

type App struct {
  Id string               `json:"id"`
  Name string             `json:"name"`
  Owner string            `json:"owner"`
  Android string          `json:"android_package"`
  IOS string              `json:"ios_bundle"`
  CreatedOn time.Time     `json:"created_on"`
  LastUpdate time.Time    `json:"last_update"`
}

type Device struct {
  Id string               `json:"device_id"`
  Model string            `json:"model"`
  Vendor string           `json:"vendor"`
  Platform string         `json:"platform"`
  PlatformVersion string  `json:"platform_version"`
  Keys []ApiKey           `json:"keys"`
  CreatedOn time.Time     `json:"created_on"`
  LastUpdate time.Time    `json:"last_update"`
}

type ApiKey struct {
  Key string              `json:"key"`
  AppId string            `json:"app_id"`
}

/**
 *  Create a new app in Google Datastore.
 */
func PostApp(r *http.Request, app *App) (bool, error) {
  context := appengine.NewContext(r)
  key := hash(app.Id)
  appKey := datastore.NewKey(context, "application", key, 0, nil)
  err := datastore.Get(context, appKey, app)
  if err != nil {
    (*app).CreatedOn = time.Now()
    (*app).LastUpdate = time.Now()
    _, err := datastore.Put(context, appKey, app)
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
func GetApp(r *http.Request, appId string) App {
  context := appengine.NewContext(r)
  key := hash(appId)
  appKey := datastore.NewKey(context, "application", key, 0, nil)
  var app App
  datastore.Get(context, appKey, &app)
	return app
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

/**
 *  Create a new device in Google Datastore.
 */
func PostDevice(r *http.Request, device *Device) (bool, error) {
  context := appengine.NewContext(r)
  key := device.Id
  deviceKey := datastore.NewKey(context, "device", key, 0, nil)
  var existingDevice Device
  err := datastore.Get(context, deviceKey, &existingDevice)
  if err != nil {
    // device do not exists... let's insert it
    (*device).Keys = []ApiKey{}
    (*device).CreatedOn = time.Now()
    (*device).LastUpdate = time.Now()
    _, err := datastore.Put(context, deviceKey, device)
    if err != nil {
      return false, err
    } else {
      return true, nil
    }
  } else {
    // device exists... let's update it
    if (*device).Keys == nil {
      if existingDevice.Keys == nil {
        (*device).Keys = []ApiKey{}
      } else {
        (*device).Keys = existingDevice.Keys
      }
    }
    (*device).LastUpdate = time.Now()
    _, err := datastore.Put(context, deviceKey, device)
    if err != nil {
      return false, err
    } else {
      return true, nil
    }
  }
}

func GetApiKey(r *http.Request, device *Device, appId string) ApiKey {
  context := appengine.NewContext(r)
  if (*device).Keys == nil {
    apiKey := ApiKey{Key:random(10),AppId:appId}
    (*device).Keys = []ApiKey{apiKey}
    log.Debugf(context, fmt.Sprintf("new key for device: %v", (*device)))
    PostDevice(r, device)
    return apiKey
  } else {
    apiKey := ApiKey{Key:"", AppId:appId}
    for _, key := range (*device).Keys {
      if key.AppId == appId {
        apiKey.Key = key.Key
      }
    }
    if len(apiKey.Key) == 0 {
      apiKey.Key = random(10)
      (*device).Keys = append((*device).Keys, apiKey)
      log.Debugf(context, fmt.Sprintf("new key for device: %v", (*device)))
      PostDevice(r, device)
    }
    return apiKey
  }
}

func hash(s string) string {
  data := []byte(s)
  return fmt.Sprintf("%x", md5.Sum(data))
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
