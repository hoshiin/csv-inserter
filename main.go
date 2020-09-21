package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

const (
	driverName = "mysql"
	user       = "docker_user"
	password   = "docker_user_pwd"
	db         = "docker_db"
)

func main() {
	records := readCSV()
	err := insert(records)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readCSV() [][]string {
	currentDir, _ := os.Getwd()
	filepath := currentDir + "/sample_data.csv"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func getParamString(param string) string {
	return os.Getenv(param)
}

func xormConn() (*xorm.Engine, error) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	db := os.Getenv("MYSQL_MYSQL_DATABASE")
	engine, err := xorm.NewEngine(driverName,
		fmt.Sprintf("%s:%s@/%s", user, password, db))
	if err != nil {
		return nil, err
	}
	return engine, nil
}

type sequencerReport struct {
	GetTime    time.Time
	Speed      float64
	StopSignal int8
	CreatedAt  time.Time `xorm:"created TIMESTAMPZ"`
}

type sequencerReports []sequencerReport

// TableName .
func (rs sequencerReports) TableName() string {
	return "sequencer_reports"
}

func newSequencerReport(str []string) (sequencerReport, error) {
	getTime, err := time.Parse("2006/1/2 15:04:05.000", str[2])
	if err != nil {
		return sequencerReport{}, err
	}

	speed, err := strconv.ParseFloat(str[0], 64)
	if err != nil {
		return sequencerReport{}, err
	}

	stopSignal, err := strconv.Atoi(str[1])
	if err != nil {
		return sequencerReport{}, err
	}

	return sequencerReport{
		GetTime:    getTime,
		Speed:      speed,
		StopSignal: int8(stopSignal),
	}, nil
}

func newSequencerReports(arr [][]string) (sequencerReports, error) {
	rs := make(sequencerReports, len(arr))
	for i, line := range arr {
		r, err := newSequencerReport(line)
		if err != nil {
			return sequencerReports{}, err
		}
		rs[i] = r
	}
	return rs, nil
}

func insert(arr [][]string) error {
	db, err := xormConn()
	if err != nil {
		return err
	}
	reports, err := newSequencerReports(arr)
	if err != nil {
		return err
	}
	_, err = db.Table(&reports).Insert(&reports)
	if err != nil {
		return err
	}
	return nil
}
