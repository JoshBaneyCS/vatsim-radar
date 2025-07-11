package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Pilot represents a VATSIM aircraft
type Pilot struct {
	Callsign     string  `json:"callsign"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Altitude     int     `json:"altitude"`
	Groundspeed  int     `json:"groundspeed"`
	Heading      int     `json:"heading"`
	PlannedDep   string  `json:"planned_depairport"`
	PlannedDest  string  `json:"planned_destairport"`
	AircraftType string  `json:"aircraft"`
}

// RealFlight represents a real-world aircraft
type RealFlight struct {
	Callsign  string  `json:"callsign"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) SetContext(ctx context.Context) {
	a.ctx = ctx
	go a.startFetcher()
}

func (a *App) FrontendAssets() http.FileSystem {
	return os.DirFS("./frontend/dist")
}

func (a *App) startFetcher() {
	for {
		vatsimData := a.fetchVatsim()
		if vatsimData != nil {
			runtime.EventsEmit(a.ctx, "vatsim-data", vatsimData)
		}

		realData := a.fetchRealFlights()
		if realData != nil {
			runtime.EventsEmit(a.ctx, "real-data", realData)
		}

		time.Sleep(15 * time.Second)
	}
}

func (a *App) fetchVatsim() []Pilot {
	resp, err := http.Get("https://data.vatsim.net/v3/vatsim-data.json")
	if err != nil {
		log.Println("Error fetching VATSIM:", err)
		return nil
	}
	defer resp.Body.Close()

	var result struct {
		Pilots []Pilot `json:"pilots"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding VATSIM JSON:", err)
		return nil
	}

	return result.Pilots
}

func (a *App) fetchRealFlights() []RealFlight {
	resp, err := http.Get("https://opensky-network.org/api/states/all")
	if err != nil {
		log.Println("Error fetching OpenSky:", err)
		return nil
	}
	defer resp.Body.Close()

	var result struct {
		States [][]interface{} `json:"states"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding OpenSky JSON:", err)
		return nil
	}

	var flights []RealFlight
	for _, s := range result.States {
		if len(s) < 11 {
			continue
		}
		flight := RealFlight{
			Callsign:  asString(s[1]),
			Longitude: asFloat(s[5]),
			Latitude:  asFloat(s[6]),
			Altitude:  asFloat(s[7]),
			Speed:     asFloat(s[9]),
			Heading:   asFloat(s[10]),
		}
		flights = append(flights, flight)
	}

	return flights
}

func asFloat(v interface{}) float64 {
	if val, ok := v.(float64); ok {
		return val
	}
	return 0
}

func asString(v interface{}) string {
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}
