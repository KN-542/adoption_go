package repository

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DB, Redis, AWS, Google以外のIF
type IOuterIFRepository interface {
	// 日本の休日取得
	HolidaysJp(year int) ([]time.Time, error)
}

type OuterIFRepository struct{}

func NewOuterRepository() IOuterIFRepository {
	return &OuterIFRepository{}
}

// 日本の休日取得
func (o *OuterIFRepository) HolidaysJp(year int) ([]time.Time, error) {
	url := "https://holidays-jp.github.io/api/v1/date.json"

	result, err := http.Get(url)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var holidays map[string]string
	if err := json.Unmarshal(body, &holidays); err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var dates []time.Time

	for dateStr, _ := range holidays {
		if dateStr[:4] == strconv.Itoa(year) {
			parts := strings.Split(dateStr, "-")
			month, _ := strconv.Atoi(parts[1])
			day, _ := strconv.Atoi(parts[2])

			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			dates = append(dates, date)
		}
	}

	for d := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC); d.Year() == year; d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			if !isDateInSlice(d, dates) {
				dates = append(dates, d)
			}
		}
	}

	return dates, nil
}

func isDateInSlice(a time.Time, list []time.Time) bool {
	for _, b := range list {
		if a.Equal(b) {
			return true
		}
	}
	return false
}
