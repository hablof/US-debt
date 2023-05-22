package debtseeker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	usTreasuryEndpoint = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v2/accounting/od/debt_to_penny"
	rawParams          = "fields=tot_pub_debt_out_amt,record_date&filter=record_date:lte:2023-05-19&sort=-record_date&page[size]=1"
)

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
}

func NewSeeker() *DebtSeeker {
	c := http.Client{Timeout: 5 * time.Second}

	return &DebtSeeker{
		c:        &c,
		tempData: nil,
	}
}

func (ds *DebtSeeker) FetchData() error {

	ctx, cf := context.WithTimeout(context.Background(), 3*time.Second)
	defer cf()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, usTreasuryEndpoint, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = rawParams
	// req.URL.RawQuery = url.Values(headers).Encode()

	resp, err := ds.c.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var dataSample endpointOutputScheme
	if err := json.Unmarshal(body, &dataSample); err != nil {
		return err
	}

	if len(dataSample.Data) == 0 {
		return fmt.Errorf("no data in response")
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
