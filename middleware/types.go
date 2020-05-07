package middleware

import "knocker"

type Skipper func(knocker.Context) bool


func DefaultSkipper(knocker.Context) bool {
	return false
}