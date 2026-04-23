package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	return id
}

func ParseWeight(rawWeight string) float64 {
	weight, err := strconv.ParseFloat(rawWeight, 64)
	if err != nil {
		fmt.Println(err)
	}
	return weight
}

func ParseKoli(rawKoli string) int {
	koli, err := strconv.Atoi(rawKoli)
	if err != nil {
		fmt.Println(err)
	}
	return koli
}

func ParseDate(input, layout string) (time.Time, error) {
	re := regexp.MustCompile(`\d{2}/\d{2}/\d{4}`)
	match := re.FindString(input)

	// loc, _ := time.LoadLocation("Asia/Jakarta")

	if match == "" {
		return time.Time{}, fmt.Errorf("no date pattern found in: %s", input)
	}
	fmt.Printf("The match date is: %s\n", match)

	// Parse the 'match' string, not 'cleanedDate'
	// t, err := time.ParseInLocation("02/01/2006", match, loc)
	t, err := time.Parse("02/01/2006", match)
	// fmt.Printf("%s\n", t)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
