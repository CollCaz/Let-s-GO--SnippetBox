package models

import "errors"

// Returned when Get method does not find any sql records matching the query
var ErrorNoRecord = errors.New("models: no matching record found")
