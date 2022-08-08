package lib

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type Island struct {
	CategoryId  int     `csv:"category_id"`
	IslandId    int     `csv:"island_id"`
	Atoll       int     `csv:"atoll"`
	EnglishName string  `csv:"english_name"`
	DhivehiName string  `csv:"dhivehi_name"`
	ArabicName  string  `csv:"arabic_name"`
	Latitude    float64 `csv:"latitude"`
	Longitude   float64 `csv:"longitude"`
	Status      int8    `csv:"status"`
}

type Atoll struct {
	CategoryId  int    `csv:"category_id"`
	Name        string `csv:"island_id"`
	ArabicName  string `csv:"arabic_name"`
	DhivehiName string `csv:"dhivehi_name"`
}

type PrayerTime struct {
	CategoryId int `csv:"category_id"`
	Date       int `csv:"date"`
	Fajr       int `csv:"fajr"`
	Sunrise    int `csv:"sunrise"`
	Duhr       int `csv:"duhr"`
	Asr        int `csv:"asr"`
	Maghrib    int `csv:"maghrib"`
	Isha       int `csv:"isha"`
}

func ParseCSV[V Island | Atoll | PrayerTime](name string) []*V {
	if file, err := os.Open(fmt.Sprintf("assets/%s.csv", name)); err != nil {
		panic(err)
	} else {
		defer file.Close()

		data := []*V{}

		if err := gocsv.UnmarshalFile(file, &data); err != nil {
			panic(err)
		}

		return data
	}
}

func (prayer PrayerTime) GetValue(p string) int {
	var pr int

	switch p {
	case "fajr":
		pr = prayer.Fajr
	case "sunrise":
		pr = prayer.Sunrise
	case "duhr":
		pr = prayer.Duhr
	case "asr":
		pr = prayer.Asr
	case "maghrib":
		pr = prayer.Maghrib
	case "isha":
		pr = prayer.Isha
	}

	return pr
}
