package reader

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	utm "github.com/im7mortal/UTM"
)

var (
	ErrMissingFields     = errors.New("reader: missing fields in CSV")
	ErrInvalidTimeString = errors.New("reader: invalid format for time")
	ErrInvalidLocation   = errors.New("reader: invalid format for location")
)

var (
	utmStrRe = regexp.MustCompile(`UTM (\d+)(\w) (\d+) (\d+)`)
	gpsStrRe = regexp.MustCompile(`GPS (\w+)°(\w+)'(\w+)"(\w) (\w+)°(\w+)'(\w+)"(\w)`) //                                               1     2     3     4    5     6     7     8
)

type LatLon struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Alert struct {
	Type     string    `json:"type"`
	Time     time.Time `json:"time"`
	Location LatLon    `json:"location"`
	Message  string    `json:"message"`
}

func NewAlert(values []string) (*Alert, error) {
	if len(values) != 6 {
		return nil, ErrMissingFields
	}

	alert := &Alert{}
	alert.Type = strings.TrimSpace(values[0])
	// serialNo := values[1]
	timeStr := strings.TrimSpace(values[2])
	dateStr := strings.TrimSpace(values[3])
	finalTime, err := time.Parse("02/01/2006 1504 MST", dateStr+" "+timeStr)
	if err != nil {
		return nil, err
	}
	alert.Time = finalTime

	coordinates := strings.TrimSpace(values[4])
	switch {
	case strings.HasPrefix(coordinates, "UTM"):
		{
			matches := utmStrRe.FindStringSubmatch(coordinates)
			if matches == nil {
				return nil, ErrInvalidLocation
			}
			zoneNum, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, err
			}
			zoneLetter := matches[2]
			e, _ := strconv.ParseFloat(matches[3], 64)
			n, _ := strconv.ParseFloat(matches[4], 64)
			utmCoor := utm.Coordinate{
				ZoneLetter: zoneLetter,
				ZoneNumber: zoneNum,
				Easting:    e,
				Northing:   n,
			}
			latLon, err := utmCoor.ToLatLon()
			if err != nil {
				return nil, err
			}

			alert.Location = LatLon(latLon)
		}
	case strings.HasPrefix(coordinates, "GPS"):
		{
			matches := gpsStrRe.FindStringSubmatch(coordinates)
			if matches == nil {
				return nil, ErrInvalidLocation
			}

			latDir := matches[4]
			lat, err := dmsToDd(matches[1], matches[2], matches[3])
			if err != nil {
				return nil, err
			}
			switch latDir {
			case "N": // sanity check
			case "S":
				lat = -lat
			default:
				return nil, ErrInvalidLocation
			}

			longDir := matches[8]
			long, err := dmsToDd(matches[1], matches[2], matches[3])
			if err != nil {
				return nil, err
			}

			switch longDir {
			case "E": // again, sanity check
			case "W":
				long = -long
			}

			alert.Location = LatLon{lat, long}
		}
	default:
		return nil, ErrInvalidLocation
	}

	alert.Message = strings.TrimPrefix(values[5], "LABELLED AS")
	return alert, nil
}
