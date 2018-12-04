package ADDS

import (
	"github.com/kellydunn/golang-geo"
)

type ADDSTAF struct {
	Text          string   `xml:"raw_text"`
	StationID     string   `xml:"station_id"`
	IssueTime     ADDSTime `xml:"issue_time"`
	BulletinTime  ADDSTime `xml:"bulletin_time"`
	ValidTimeFrom ADDSTime `xml:"valid_time_from"`
	ValidTimeTo   ADDSTime `xml:"valid_time_to"`
	Remarks       string   `xml:"remarks"`
	//TODO: "Forecast" parsing., wind_dir_degrees, wind_speed_kt, etc.
}

func GetADDSTAFs(url string) ([]ADDSTAF, error) {
	r, err := GetADDSData(url)
	return r.TAFs, err
}

func GetADDSTAFsByIdent(ident string) ([]ADDSTAF, error) {
	return GetADDSTAFs(urlADDSDataByIdent("tafs", ident))
}

// Gets the most recent TAFs that were obtained at least within the last 1.5 hours within a (lat, lng) defined rectangle.
func GetLatestADDSTAFsInRect(bottomLeft, topRight *geo.Point) ([]ADDSTAF, error) {
	return GetADDSTAFs(urlLatestADDSDataInRect("tafs", bottomLeft, topRight))
}

// Gets the most recent TAFs that were obtained at least within the last 1.5 hours within "radius" statute miles of "pt".
func GetLatestADDSTAFsInRadiusOf(radius uint, pt *geo.Point) ([]ADDSTAF, error) {
	return GetADDSTAFs(urlLatestADDSDataInRadiusOf("tafs", radius, pt))
}

// Gets the most recent TAFs that were obtained at least within "route" staute miles of the defined route.
// "route" is in the format "lng1,lat1;lng2,lat2;..." or "ident1;ident2;..."
func GetLatestADDSTAFsAlongRoute(dist float64, route string) ([]ADDSTAF, error) {
	return GetADDSTAFs(urlLatestADDSDataAlongRoute("tafs", dist, route))
}
