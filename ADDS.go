package ADDS

import (
	"encoding/xml"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ADDSTime struct {
	Time time.Time
}

func (t *ADDSTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	timeFormat := "2006-01-02T15:04:05Z"
	var inTime string
	d.DecodeElement(&inTime, &start)
	t2, err := time.Parse(timeFormat, inTime)
	if err != nil {
		return err
	}
	*t = ADDSTime{t2}
	return nil
}

type ADDSData struct {
	METARs []ADDSMETAR `xml:"METAR"`
	PIREPs []ADDSPIREP `xml:"AircraftReport"`
}

type ADDSResponse struct {
	RequestIndex int      `xml:"request_index"`
	Data         ADDSData `xml:"data"`
}

func GetADDSData(url string) (ADDSData, error) {
	fmt.Printf("URL: %s\n", url)
	var ret ADDSResponse
	resp, err := http.Get(url)
	if err != nil || !strings.HasPrefix(resp.Status, "200") {
		return ret.Data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret.Data, err
	}

	// Parse 'body'.
	err = xml.Unmarshal([]byte(body), &ret)

	//FIXME: Parse the error code and "warning" text that is sent.

	return ret.Data, nil
}

var reportFlags = map[string]string{
	"metars":          "&mostRecentForEachStation=constraint",
	"aircraftreports": "",
}

func urlADDSDataByIdent(dataSource string, ident string) string {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=%s&requestType=retrieve&format=xml&stationString=%s&hoursBeforeNow=1.5%s", dataSource, ident, reportFlags[dataSource])
	return url
}

func urlLatestADDSDataInRect(dataSource string, bottomLeft, topRight *geo.Point) string {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=%s&requestType=retrieve&format=xml&hoursBeforeNow=1.25&minLat=%f&minLon=%f&maxLat=%f&maxLon=%f%s", dataSource, bottomLeft.Lat(), bottomLeft.Lng(), topRight.Lat(), topRight.Lng(), reportFlags[dataSource])
	return url
}

func urlLatestADDSDataInRadiusOf(dataSource string, radius uint, pt *geo.Point) string {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=%s&requestType=retrieve&format=xml&hoursBeforeNow=1.25&radialDistance=%d;%f,%f%s", dataSource, radius, pt.Lng(), pt.Lat(), reportFlags[dataSource])
	return url
}

func urlLatestADDSDataAlongRoute(dataSource string, dist float64, route string) string {
	url := fmt.Sprintf("https://aviationweather.gov/adds/dataserver_current/httpparam?dataSource=%s&requestType=retrieve&format=xml&hoursBeforeNow=1.25&flightPath=%f;%s%s", dataSource, dist, route, reportFlags[dataSource])
	return url
}
