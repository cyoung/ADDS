package ADDS

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"io/ioutil"
	"net/http"
	"strings"
)

type ADDSMETAR struct {
	Text           string   `xml:"raw_text"`
	StationID      string   `xml:"station_id"`
	Observation    ADDSTime `xml:"observation_time"`
	Latitude       float64  `xml:"latitude"`
	Longitude      float64  `xml:"longitude"`
	Temp           float64  `xml:"temp_c"`
	Dewpoint       float64  `xml:"dewpoint_c"`
	WindDirection  float64  `xml:"wind_dir_degrees"`
	WindSpeed      float64  `xml:"wind_speed_kt"`
	Visibility     float64  `xml:"visibility_statute_mi"`
	Altimeter      float64  `xml:"altim_in_hg"`
	FlightCategory string   `xml:"flight_category"`
}

func GetADDSMETARs(url string) ([]ADDSMETAR, error) {
	var ret ADDSResponse
	resp, err := http.Get(url)
	if err != nil || !strings.HasPrefix(resp.Status, "200") {
		return ret.Data.METARs, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret.Data.METARs, err
	}

	// Parse 'body'.
	err = xml.Unmarshal([]byte(body), &ret)

	// Empty response? Typically happens with invalid identifiers.
	if len(ret.Data.METARs) == 0 {
		return ret.Data.METARs, errors.New("No results.")
	}
	return ret.Data.METARs, nil
}

func GetADDSMETARsByIdent(ident string) ([]ADDSMETAR, error) {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=%s&hoursBeforeNow=1.5", ident)
	return GetADDSMETARs(url)
}

// Gets the most recent METARs that were obtained at least within the last 1.5 hours within a (lat, lng) defined rectangle.
func GetLatestADDSMETARsInRect(bottomLeft, topRight *geo.Point) ([]ADDSMETAR, error) {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&mostRecentForEachStation=constraint&hoursBeforeNow=1.25&minLat=%f&minLon=%f&maxLat=%f&maxLon=%f", bottomLeft.Lat(), bottomLeft.Lng(), topRight.Lat(), topRight.Lng())
	return GetADDSMETARs(url)
}

// Gets the most recent METARs that were obtained at least within the last 1.5 hours within "radius" of "pt".
func GetLatestADDSMETARsInRadiusOf(radius uint, pt *geo.Point) ([]ADDSMETAR, error) {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&mostRecentForEachStation=constraint&hoursBeforeNow=1.25&radialDistance=%d;%f,%f", radius, pt.Lng(), pt.Lat())
	return GetADDSMETARs(url)
}

// Gets the most recent METARs that were obtained at least within "route" staute miles of the defined route.
// "route" is in the format "lng1,lat1;lng2,lat2;..." or "ident1;ident2;..."
func GetLatestADDSMETARsAlongRoute(dist float64, route string) ([]ADDSMETAR, error) {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&mostRecentForEachStation=constraint&hoursBeforeNow=1.25&flightPath=%f;%s", dist, route)
	return GetADDSMETARs(url)
}

func GetLatestADDSMETARs(ident string) (ret ADDSMETAR, err error) {
	metars, errn := GetADDSMETARsByIdent(ident)
	if errn != nil {
		return ret, errn
	}

	// Get the latest observation time.
	for _, v := range metars {
		if v.Observation.Time.After(ret.Observation.Time) {
			// This observation is later than the current one.
			ret = v
		}
	}

	return
}
