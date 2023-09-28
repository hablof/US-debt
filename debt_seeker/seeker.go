package debtseeker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// const (
// 	usTreasuryEndpoint = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v2/accounting/od/debt_to_penny"
// 	rawParams          = "fields=tot_pub_debt_out_amt,record_date&sort=-record_date&page[size]=1"
// )

type dataRecord struct {
	TotalDebt  string `json:"tot_pub_debt_out_amt"`
	RecordDate string `json:"record_date"`
}

func (dr *dataRecord) getDebt() (uint64, error) {
	if dr == nil {
		return 0, fmt.Errorf("empty data")
	}

	u, err := strconv.ParseUint(strings.Split(dr.TotalDebt, ".")[0], 10, 64)
	return uint64(u), err
}

func (dr *dataRecord) getDate() (time.Time, error) {
	if dr == nil {
		return time.Time{}, fmt.Errorf("empty data")
	}

	return time.Parse("2006-01-02", dr.RecordDate)
}

type endpointOutputScheme struct {
	Data []dataRecord `json:"data"`
	Err  string       `json:"error"`
}

// var (
// 	headers = map[string][]string{
// 		"fields":     {"tot_pub_debt_out_amt", "record_date"},
// 		"filter":     {"record_date:lte:2023-05-19"},
// 		"sort":       {"-record_date"},
// 		"page[size]": {"1"},
// 	}
// )

type DebtSeeker struct {
	c        *http.Client
	tempData *dataRecord

	apiEndpoint string `yaml:"api-endpoint"`
	apiParams   string `yaml:"api-params"`
}

func NewSeeker() (*DebtSeeker, error) {
	c := http.Client{Timeout: 15 * time.Second}

	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := struct {
		ApiEndpoint string `yaml:"api-endpoint"`
		ApiParams   string `yaml:"api-params"`
	}{}

	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &DebtSeeker{
		c:           &c,
		tempData:    nil,
		apiEndpoint: cfg.ApiEndpoint,
		apiParams:   cfg.ApiParams,
	}, nil
}

func (ds *DebtSeeker) FetchData() error {

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ds.apiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("не удалось создать http-request: %w", err)
	}

	req.URL.RawQuery = ds.apiParams
	// req.URL.RawQuery = url.Values(headers).Encode()

	resp, err := ds.c.Do(req)
	if err != nil {
		return fmt.Errorf("не удалось выполнить http-request: %w", err)
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	var dataSample endpointOutputScheme
	if err := json.NewDecoder(resp.Body).Decode(&dataSample); err != nil {
		return fmt.Errorf("не удалось прочитать тело ответа: %w", err)
	}

	if len(dataSample.Data) == 0 {
		return fmt.Errorf("no data in response")
	}

	if dataSample.Err != "" {
		return fmt.Errorf("some error in fetched data: %s", dataSample.Err)
	}

	ds.tempData = &dataSample.Data[0]

	return nil
}

func (ds *DebtSeeker) GetDebt() (uint64, error) {

	return ds.tempData.getDebt()
}

func (ds *DebtSeeker) GetDate() (time.Time, error) {

	return ds.tempData.getDate()
}
