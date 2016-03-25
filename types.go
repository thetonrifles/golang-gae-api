package api

import (
  "encoding/json"
)

type ItemsWrapper struct {
  User User               `json:"user"`
  Type string             `json:"type"`
  Items []json.RawMessage `json:"items"`
}

/*
type RawItems struct {
  Items []RawItem         `json:"items"`
}

type NormalizedItems struct {
  Items []NormalizedItem  `json:"items"`
}

type ActivityItems struct {
  Items []ActivityItem    `json:"items"`
}

type RawItem struct {
  Timestamp int           `json:"start"`
}

type NormalizedItem struct {
  Timestamp int           `json:"timestamp"`
}

type ActivityItem struct {
  Timestamp int           `json:"timestamp"`
}
*/

type ItemSent struct {
  Count int               `json:"count"`
}
