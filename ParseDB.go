package ADDS

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"github.com/kellydunn/golang-geo"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"sort"
)

var CSVRequiredFields = []string{"ident", "type", "name", "latitude_deg", "longitude_deg", "continent", "iso_country"}

// Imports database from http://ourairports.com/data/ file.
func ImportCSVToNewSQLite(csvFile, dbFile string) error {

	// Open SQLite database.
	os.Remove(dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create new table.
	sqlStmt := `
	create table airports (id integer not null primary key, ident text, type text, name text, latitude_deg real, longitude_deg real, continent text, iso_country text);
	delete from airports;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	// Prepare insert statement.
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into airports(ident, type, name, latitude_deg, longitude_deg, continent, iso_country) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Open CSV file and start parsing.
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer f.Close()

	colLabels := make(map[string]int, 0)
	reqLen := 0

	// Start reading in the CSV.
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(colLabels) == 0 {
			// Load in the labels.
			for i := 0; i < len(record); i++ {
				colLabels[record[i]] = i
			}
			// Make sure we have the required fields.
			for _, v := range CSVRequiredFields {
				if fieldPos, ok := colLabels[v]; !ok {
					return errors.New("missing required fields.")
				} else {
					if fieldPos > reqLen {
						reqLen = fieldPos // Track the maximum field position in the required fields.
					}
				}
			}
			reqLen++
			continue
		}

		// Extract ident, type, name, latitude_deg, longitude_deg, continent, iso_country.
		if len(record) < reqLen {
			continue
		}

		vals := make([]interface{}, 0)
		for _, field := range CSVRequiredFields {
			vals = append(vals, record[colLabels[field]])
		}

		_, err = stmt.Exec(vals...)
		if err != nil {
			return err
		}

	}

	tx.Commit()

	return nil
}

type AirportDB struct {
	db    *sql.DB
	Cache []Airport // Loaded into memory when initialized.
}

type Airport struct {
	ID        int
	Ident     string
	Type      string
	Name      string
	Lat       float64
	Lng       float64
	Continent string
	Country   string
}

func NewAirportDB(dbFile string) (*AirportDB, error) {
	ret := new(AirportDB)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return ret, err
	}
	ret.db = db

	// Load db into memory.
	rows, err := ret.db.Query("select id, ident, type, name, latitude_deg, longitude_deg, continent, iso_country from airports")
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		var thisAirport Airport
		err = rows.Scan(&thisAirport.ID, &thisAirport.Ident, &thisAirport.Type, &thisAirport.Name, &thisAirport.Lat, &thisAirport.Lng, &thisAirport.Continent, &thisAirport.Country)
		if err != nil {
			return ret, err
		}
		ret.Cache = append(ret.Cache, thisAirport)
	}
	err = rows.Err()
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// Airport sort-by-distance.
type AirportDistance struct {
	ThisAirport Airport
	Distance    float64
}

type AirportDistances []AirportDistance

func (a AirportDistances) Len() int {
	return len(a)
}
func (a AirportDistances) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AirportDistances) Less(i, j int) bool {
	return a[i].Distance < a[j].Distance
}

func (a *AirportDB) FindClosestAirports(lat, lng float64) []AirportDistance {
	p := geo.NewPoint(lat, lng)

	s := make([]AirportDistance, 0)
	// Calculate distances to all airports in 'Cache'.
	for _, airport := range a.Cache {
		p2 := geo.NewPoint(airport.Lat, airport.Lng)
		dist := p.GreatCircleDistance(p2) // Distance in km.
		distInNM := dist * 0.539957       // Distance in nm.
		s = append(s, AirportDistance{ThisAirport: airport, Distance: distInNM})
	}

	sort.Sort(AirportDistances(s))

	return s
}

