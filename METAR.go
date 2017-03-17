package ADDS

import (
	"github.com/kellydunn/golang-geo"
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
	r, err := GetADDSData(url)
	return r.METARs, err
}

func GetADDSMETARsByIdent(ident string) ([]ADDSMETAR, error) {
	return GetADDSMETARs(urlADDSDataByIdent("metars", ident))
}

// Gets the most recent METARs that were obtained at least within the last 1.5 hours within a (lat, lng) defined rectangle.
func GetLatestADDSMETARsInRect(bottomLeft, topRight *geo.Point) ([]ADDSMETAR, error) {
	return GetADDSMETARs(urlLatestADDSDataInRect("metars", bottomLeft, topRight))
}

// Gets the most recent METARs that were obtained at least within the last 1.5 hours within "radius" of "pt".
func GetLatestADDSMETARsInRadiusOf(radius uint, pt *geo.Point) ([]ADDSMETAR, error) {
	return GetADDSMETARs(urlLatestADDSDataInRadiusOf("metars", radius, pt))
}

// Gets the most recent METARs that were obtained at least within "route" staute miles of the defined route.
// "route" is in the format "lng1,lat1;lng2,lat2;..." or "ident1;ident2;..."
func GetLatestADDSMETARsAlongRoute(dist float64, route string) ([]ADDSMETAR, error) {
	return GetADDSMETARs(urlLatestADDSDataAlongRoute("metars", dist, route))
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
