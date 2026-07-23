package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("parse date: %w", err)
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("missing days interval")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid days interval: %w", err)
		}
		if days < 1 || days > 400 {
			return "", fmt.Errorf("days interval out of range: %d", days)
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format: %q", repeat)
	}

	return date.Format(dateFormat), nil
}

func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	now := time.Now()
	if nowParam != "" {
		parsed, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "invalid now parameter", http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
