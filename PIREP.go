package ADDS

import (
	"github.com/kellydunn/golang-geo"
)

/*
	ADDS now calls these PIREPs "aircraftreports".
*/

type ADDSPIREP struct {
	ReceiptTime     ADDSTime `xml:"receipt_time"`
	ObservationTime ADDSTime `xml:"observation_time"`
	AircraftRef     string   `xml:"aircraft_ref"`
	Latitude        float64  `xml:"latitude"`
	Longitude       float64  `xml:"longitude"`
	Altitude        float64  `xml:"altitude_ft_msl"` // Feet, MSL.
	//TODO: "turbulence_condition", "wind_dir_degrees", "wind_dir_degrees", "icing_condition", etc.
	ReportType string `xml:"report_type"`
	Text       string `xml:"raw_text"`
}

func GetADDSPIREPs(url string) ([]ADDSPIREP, error) {
	r, err := GetADDSData(url)
	return r.PIREPs, err
}

// Gets the most recent PIREPs that were obtained at least within the last 1.5 hours within a (lat, lng) defined rectangle.
func GetLatestADDSPIREPsInRect(bottomLeft, topRight *geo.Point) ([]ADDSPIREP, error) {
	return GetADDSPIREPs(urlLatestADDSDataInRect("aircraftreports", bottomLeft, topRight))
}

// Gets the most recent PIREPs that were obtained at least within the last 1.5 hours within "radius" of "pt".
func GetLatestADDSPIREPsInRadiusOf(radius uint, pt *geo.Point) ([]ADDSPIREP, error) {
	return GetADDSPIREPs(urlLatestADDSDataInRadiusOf("aircraftreports", radius, pt))
}

// Gets the most recent PIREPs that were obtained at least within "route" staute miles of the defined route.
// "route" is in the format "lng1,lat1;lng2,lat2;..." or "ident1;ident2;..."
func GetLatestADDSPIREPsAlongRoute(dist float64, route string) ([]ADDSPIREP, error) {
	return GetADDSPIREPs(urlLatestADDSDataAlongRoute("aircraftreports", dist, route))
}
