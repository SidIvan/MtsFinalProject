package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
	"io"
	"math"
	"math/rand"
	"net/http"
)

var (
	Main *zap.Logger
	mode string
)

func main() {
	flag.StringVar(&mode, "mode", "200", "mocked response status")
	flag.Parse()
	var err error
	Main, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/drivers", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if mode == "404" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if mode == "500" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	reqBodyData, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		Main.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var reqData QueryBody
	err = json.Unmarshal(reqBodyData, &reqData)
	if err != nil {
		Main.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var respData []Driver
	for i := 0; i <= rand.Intn(9)+1; i++ {
		id, err := gonanoid.New()
		if err != nil {
			Main.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		dist := rand.Float64() * reqData.Radius
		angle := rand.Float64() * 360
		latLngLiteral := LatLngLiteral{
			Lat: reqData.Lat + dist*math.Sin(angle),
			Lng: reqData.Lng + dist*math.Cos(angle),
		}
		for latLngLiteral.Lng > 180 {
			latLngLiteral.Lng -= 360
		}
		for latLngLiteral.Lng < -180 {
			latLngLiteral.Lng += 360
		}
		for latLngLiteral.Lat > 180 {
			latLngLiteral.Lat -= 360
		}
		if latLngLiteral.Lat > 90 {
			latLngLiteral.Lat = 180 - latLngLiteral.Lat
		}
		if latLngLiteral.Lat < -90 {
			latLngLiteral.Lat = -180 - latLngLiteral.Lat
		}
		for latLngLiteral.Lat < -180 {
			latLngLiteral.Lat += 360
		}
		respData = append(respData, Driver{
			Id:            id,
			LatLngLiteral: latLngLiteral,
		})
	}
	respByteData, err := json.Marshal(respData)
	numBytes, err := w.Write(respByteData)
	if err != nil {
		Main.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if numBytes != binary.Size(respByteData) {
		Main.Error(fmt.Sprintf("Not full response sent, need %d, sent %d", binary.Size(respByteData), numBytes))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type QueryBody struct {
	LatLngLiteral
	Radius float64 `json:"radius"`
}

type Driver struct {
	LatLngLiteral
	Id string `json:"id"`
}

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
