//
// Library to interact with DnsNode v3 API
//
// The DnsNode token and API URL is passed as a struct to New()
//

package dnsnode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const DEFAULT_DNSNODE_API string = "https://dnsnodeapi.netnod.se/apiv3"

var ALLOWED_ALGORITMS = map[string]int{
	"hmac-md5":    1,
	"hmac-sha1":   1,
	"hmac-sha256": 1,
	"hmac-sha512": 1,
}

var Loglevels = map[string]log.Level{
	"debug":   log.DebugLevel,
	"error":   log.ErrorLevel,
	"warning": log.WarnLevel,
	"info":    log.InfoLevel,
}

// ---------------------------------------------------------------------------
//   Types
// ---------------------------------------------------------------------------

type TSIGType struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Alg  string `json:"alg"`
}

type PrimaryType struct {
	IP   string `json:"ip"`
	TSIG string `json:"tsig"`
}

type ZoneType struct {
	Name        string        `json:"name"`
	Primaries   []PrimaryType `json:"masters"`
	Product     string        `json:"product"`
	Endcustomer string        `json:"endcustomer"`
}

type ZoneCreateRequestType struct {
	Name        string        `json:"name"`
	Primaries   []PrimaryType `json:"masters"`
	Product     string        `json:"product"`
	Endcustomer string        `json:"endcustomer"`
}

type DistPrimaryType struct {
	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`
	Serial      uint32 `json:"serial"`
	Timestamp   uint32 `json:"timestamp"`
}

type SiteType struct {
	Name      string `json:"name"`
	Serial    uint32 `json:"serial"`
	Timestamp uint32 `json:"timestamp"`
}

type ResponseType struct {
	ConfiguredSites         int               `json:"configured_sites"`
	SitesInMaintenance      int               `json:"sites_in_maintenance"`
	ConfiguredDistPrimaries int               `json:"configured_distmasters"`
	CurrentSerial           uint32            `json:"current_serial"`
	CurrentTimestamp        uint32            `json:"current_timestamp"`
	DistPrimaries           []DistPrimaryType `json:"distmasters"`
	Sites                   []SiteType        `json:"sites"`
}

// ---------------------------------------------------------------------------
//   Helpers to create a instance of DnsNodeClient and read configuration file
// ---------------------------------------------------------------------------

// Dnsnode client Parameters
type DnsNodeParam struct {
	Token string
	URL   string
	Debug bool
}

// DnsNodeClient instance
type DnsNodeClient struct {
	p DnsNodeParam
}

// Create a new DnsNode client
func New(param DnsNodeParam) *DnsNodeClient {
	if param.URL == "" {
		param.URL = DEFAULT_DNSNODE_API
	}

	client := new(DnsNodeClient)
	client.p = param
	return client
}

// ---------------------------------------------------------------------------
//   Internal Utils
// ---------------------------------------------------------------------------

// Verify that algoritm is supported
// An empty alg is allowed, it signals "leave unchanged" for partial updates
func verifyAlgKey(alg string, key string) error {
	if alg == "" {
		return nil
	}
	if _, ok := ALLOWED_ALGORITMS[alg]; !ok {
		return fmt.Errorf("algorithm %q not recognized", alg)
	}
	return nil
}

// Call DnsNode API
// If data is non-nil, encode as JSON and POST
// Return response as []byte. It is up to caller to decode the response
func (dnsnodec *DnsNodeClient) call(method string, endpoint string, data any) ([]byte, error) {
	if dnsnodec.p.Debug {
		log.Debug(dnsnodec.p.URL + endpoint)
	}
	var req *http.Request
	var err error
	if data != nil {
		var jsonData []byte
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
		log.Debugf("data=%s", jsonData)
		req, err = http.NewRequest(method, dnsnodec.p.URL+endpoint, bytes.NewReader(jsonData))
	} else {
		req, err = http.NewRequest(method, dnsnodec.p.URL+endpoint, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+dnsnodec.p.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("%s %s %s\n", method, endpoint, body)
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("dnsnode api: %s %s: %s: %s", method, endpoint, resp.Status, body)
	}
	return body, nil
}

// ---------------------------------------------------------------------------
//   Anomalies
//   Check for anomalies for one zone
// ---------------------------------------------------------------------------

func (dnsnodec *DnsNodeClient) Anomalies(zone_name string) (*ResponseType, error) {
	if zone_name == "" {
		return nil, errors.New("zone name must be specified")
	}
	body, err := dnsnodec.call("GET", "/anomalies/serial/"+zone_name, nil)
	if err != nil {
		return nil, err
	}
	data := new(ResponseType)
	err = json.Unmarshal(body, data)
	return data, err
}

// List all zones where the serial number differs between sites/distmasters
func (dnsnodec *DnsNodeClient) AnomaliesList() ([]string, error) {
	body, err := dnsnodec.call("GET", "/anomalies/serial/", nil)
	if err != nil {
		return nil, err
	}
	var data []string
	err = json.Unmarshal(body, &data)
	return data, err
}

// ---------------------------------------------------------------------------
//   Product
//   List available products
// ---------------------------------------------------------------------------

func (dnsnodec *DnsNodeClient) Product() ([]string, error) {
	body, err := dnsnodec.call("GET", "/product/", nil)
	if err != nil {
		return nil, err
	}
	var data []string
	err = json.Unmarshal(body, &data)
	return data, err
}

// ---------------------------------------------------------------------------
//   Sites
//   Return all sites, with information on location, maintenance status etc
// ---------------------------------------------------------------------------

type SiteResponse struct {
	Name        string  `json:"name"`
	City        string  `json:"city"`
	Countrycode string  `json:"countrycode"`
	Lat         float32 `json:"lat,string"`
	Long        float32 `json:"long,string"`
	Enabled     bool    `json:"enabled"`
	Maintenance bool    `json:"maintenance"`
}

func (dnsnodec *DnsNodeClient) Sites() (*[]SiteResponse, error) {
	body, err := dnsnodec.call("GET", "/sites/", nil)
	if err != nil {
		return nil, err
	}
	log.Debugf("Body: %s", body)
	data := new([]SiteResponse)
	err = json.Unmarshal(body, &data)
	return data, err
}

// Return all available sites, regardless of whether they are configured for this customer
func (dnsnodec *DnsNodeClient) SitesAll() (*[]SiteResponse, error) {
	body, err := dnsnodec.call("GET", "/sites/all/", nil)
	if err != nil {
		return nil, err
	}
	data := new([]SiteResponse)
	err = json.Unmarshal(body, &data)
	return data, err
}

// ---------------------------------------------------------------------------
//   TSIG
//   Create, read, update, and destroy TSIG keys for use in your zone transfers
// ---------------------------------------------------------------------------

// Get TSIG(s)
func (dnsnodec *DnsNodeClient) Tsig(name string) ([]TSIGType, error) {
	body, err := dnsnodec.call("GET", "/tsig/"+name, nil)
	if err != nil {
		return nil, err
	}
	var data []TSIGType
	if name != "" {
		var temp TSIGType
		err = json.Unmarshal(body, &temp)
		if err != nil {
			return nil, err
		}
		data = append(data, temp)
	} else {
		err = json.Unmarshal(body, &data)
	}
	return data, err
}

// Create one TSIG
func (dnsnodec *DnsNodeClient) TsigCreate(name string, alg string, key string) error {
	if !strings.HasPrefix(name, "netnod-") {
		return errors.New("TSIG name must start with netnod-")
	}
	err := verifyAlgKey(alg, key)
	if err != nil {
		return err
	}
	data := map[string]string{
		"name": name,
		"alg":  alg,
		"key":  key,
	}
	_, err = dnsnodec.call("POST", "/tsig/", data)
	return err
}

// Update one TSIG
// alg and key are optional, an empty value leaves that field unchanged
func (dnsnodec *DnsNodeClient) TsigUpdate(name string, alg string, key string) error {
	if name == "" {
		return errors.New("TSIG name must be specified")
	}
	err := verifyAlgKey(alg, key)
	if err != nil {
		return err
	}

	data := map[string]string{}
	if alg != "" {
		data["alg"] = alg
	}
	if key != "" {
		data["key"] = key
	}
	_, err = dnsnodec.call("PATCH", "/tsig/"+name, data)
	return err
}

// Delete a TSIG
func (dnsnodec *DnsNodeClient) TsigDelete(name string) error {
	if name == "" {
		return errors.New("TSIG name must be specified")
	}
	_, err := dnsnodec.call("DELETE", "/tsig/"+name, nil)
	return err
}

// ---------------------------------------------------------------------------
//   Status
// ---------------------------------------------------------------------------

// Get status for a Zone
func (dnsnodec *DnsNodeClient) Status(name string) (*ResponseType, error) {
	if name == "" {
		return nil, errors.New("zone name must be specified")
	}
	body, err := dnsnodec.call("GET", "/status/"+name, nil)
	if err != nil {
		return nil, err
	}
	data := new(ResponseType)
	err = json.Unmarshal(body, data)
	return data, err
}

// ---------------------------------------------------------------------------
//   Statistics
// ---------------------------------------------------------------------------

type StatisticsValue struct {
	Site string    `json:"site"`
	QPS  []float32 `json:"qps"`
}
type StatisticsData struct {
	Timestamps []int32           `json:"timestamps"`
	Values     []StatisticsValue `json:"values"`
}

// Get statistics for a Zone
func (dnsnodec *DnsNodeClient) Statistics(name string) (*StatisticsData, error) {
	body, err := dnsnodec.call("GET", "/statistics/graph/"+name, nil)
	if err != nil {
		return nil, err
	}
	data := new(StatisticsData)
	err = json.Unmarshal(body, data)
	return data, err
}

// ---------------------------------------------------------------------------
//   Zone
// ---------------------------------------------------------------------------

// List zone(s)
func (dnsnodec *DnsNodeClient) Zone(name string) ([]ZoneType, error) {
	var err error
	var data []ZoneType
	var body []byte
	if name != "" {
		body, err = dnsnodec.call("GET", "/zone/"+name, nil)
	} else {
		body, err = dnsnodec.call("GET", "/zone/", nil)
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

// Create one Zone
func (dnsnodec *DnsNodeClient) ZoneCreate(name string, primaries []PrimaryType, product string, endcustomer string) error {

	if name == "" {
		return errors.New("Zone name must be specified")
	}

	data := ZoneCreateRequestType{
		Name:        name,
		Primaries:   primaries,
		Product:     product,
		Endcustomer: endcustomer,
	}
	_, err := dnsnodec.call("POST", "/zone/", data)
	return err
}

// Update one zone
// primaries, product and endcustomer are optional, an empty/nil value leaves that field unchanged
func (dnsnodec *DnsNodeClient) ZoneUpdate(name string, primaries []PrimaryType, product string, endcustomer string) error {

	if name == "" {
		return errors.New("Zone name must be specified")
	}

	data := map[string]any{}
	if len(primaries) > 0 {
		data["masters"] = primaries
	}
	if product != "" {
		data["product"] = product
	}
	if endcustomer != "" {
		data["endcustomer"] = endcustomer
	}

	_, err := dnsnodec.call("PATCH", "/zone/"+name, data)
	return err
}

// Delete one zone
func (dnsnodec *DnsNodeClient) ZoneDelete(name string) error {
	if name == "" {
		return errors.New("Zone name must be specified")
	}
	_, err := dnsnodec.call("DELETE", "/zone/"+name, nil)
	return err
}

// ---------------------------------------------------------------------------
//   Logs
//   Get transfer/notify logs from distmasters for a zone
// ---------------------------------------------------------------------------

type LogEntryType struct {
	ID     string         `json:"_id"`
	Source map[string]any `json:"_source"`
}

type LogsResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []LogEntryType `json:"hits"`
	} `json:"hits"`
}

// Get transfer/notify logs for a zone
func (dnsnodec *DnsNodeClient) LogsXfer(zone_name string) (*LogsResponse, error) {
	if zone_name == "" {
		return nil, errors.New("zone name must be specified")
	}
	body, err := dnsnodec.call("GET", "/logs/xfer/"+zone_name, nil)
	if err != nil {
		return nil, err
	}
	data := new(LogsResponse)
	err = json.Unmarshal(body, data)
	return data, err
}
