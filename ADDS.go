package ADDS

import (
	"encoding/xml"
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
}

type ADDSResponse struct {
	RequestIndex int      `xml:"request_index"`
	Data         ADDSData `xml:"data"`
}
