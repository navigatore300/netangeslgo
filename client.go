package netangelsgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	apiUrl   = "https://api-ms.netangels.ru/api/v1"            //dns/records/
	panelUrl = "https://panel.netangels.ru/api/gateway/token/" //gateway/token/
)

type RecordType string

// DnsRecord представляет структуру DNS-записи
type DnsRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type TXTDnsRecord struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	TTL   int    `json:"ttl,omitempty"`
}

// CreateDnsResponse представляет структуру ответа при создании DNS-записи
type CreateDnsResponse struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type Zone_entities struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Zone представляет отдельную зону
type Zone struct {
	Comment         string          `json:"comment,omitempty"`
	CreatedAt       string          `json:"created_at,omitempty"`
	Editable        bool            `json:"editable,omitempty"`
	ID              int             `json:"id,omitempty"`
	IsInTransfer    bool            `json:"is_in_transfer,omitempty"`
	IsTechnicalZone bool            `json:"is_technical_zone,omitempty"`
	Name            string          `json:"name,omitempty"`
	RecordsCount    int             `json:"records_count,omitempty"`
	SecondaryDns    []Zone_entities `json:"secondary_dns,omitempty"`
	SOAEmail        string          `json:"soa_email,omitempty"`
	TTL             int             `json:"ttl,omitempty"`
	UpdatedAt       string          `json:"updated_at,omitempty"`
}

// SecondaryDNS представляет дополнительный DNS
type SecondaryDNS struct {
	Entities []Zone `json:"entities"`
}

// Zones представляет список зон
type Zones struct {
	Count    int    `json:"count"`
	Entities []Zone `json:"entities"`
}

// TokenResponse представляет структуру ответа от API
type TokenResponse struct {
	Token string `json:"token"`
}

var recordTypes = [...]RecordType{"A", "AAAA", "CNAME", "MX", "NS", "TXT", "SRV", "CAA"}

func validateRecordType(recordType RecordType) bool {
	for _, t := range recordTypes {
		if t == recordType {
			return true
		}
	}
	return false
}

func CreateNetangelsClient(accountName string, apiKey string) NetangelsClient {
	return NetangelsClient{
		Credentials: Credentials{
			AccountName: accountName,
			ApiKey:      apiKey,
		},
		ApiToken: "",
		Logger:   log.New(), // Логгер по умолчанию
	}
}

// NetangelsClient base type
type NetangelsClient struct {
	Credentials Credentials `json:"credentials"`
	ApiToken    string      `json:"apitoken"`
	Logger      *log.Logger
}

type Details struct {
	ID         int        `json:"id,omitempty"`
	Type       RecordType `json:"type,omitempty"`
	Name       string     `json:"name,omitempty"`
	Value      string     `json:"data,omitempty"`
	Priority   int        `json:"priority,omitempty"`
	TTL        int        `json:"ttl,omitempty"`
	IP         int        `json:"ip,omitempty"`
	Hostname   string     `json:"hostname,omitempty"`
	Port       string     `json:"port,omitempty"`
	Protocol   string     `json:"protocol,omitempty"`
	Service    string     `json:"service,omitempty"`
	Weight     int        `json:"weight,omitempty"`
	Domain     string     `json:"domain,omitempty"`
	Nameserver string     `json:"nameserver,omitempty"`
	Flag       string     `json:"flag,omitempty"`
	Tag        string     `json:"tag,omitempty"`
}

// RecordResponse api type
type RecordResponse struct {
	Records []struct {
		ID        int        `json:"id,omitempty"`
		Zone_id   int        `json:"zone_id,omitempty"`
		Name      string     `json:"name,omitempty"`
		TTL       int        `json:"ttl,omitempty"`
		Value     string     `json:"value,omitempty"`
		Type      RecordType `json:"type,omitempty"`
		CreatedAt string     `json:"created_at,omitempty"`
		UpdatedAt string     `json:"updated_at,omitempty"`
		Priority  int        `json:"priority,omitempty"`
		IP        string     `json:"ip,omitempty"`
		Details   Details    `json:"details,omitempty"`
	} `json:"records"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type GetRecordDetails struct {
	Value string `json:"value,omitempty"`
}

type GetRecord struct {
	CreatedAt string             `json:"created_at,omitempty"`
	Details   []GetRecordDetails `json:"details,omitempty"`
	ID        int                `json:"id,omitempty"`
	TTL       int                `json:"ttl,omitempty"`
	Name      string             `json:"name,omitempty"`
	Type      RecordType         `json:"type,omitempty"`
	UpdatedAt string             `json:"updated_at,omitempty"`
}
type GetRecordsList struct {
	Count    int         `json:"count"`
	Entities []GetRecord `json:"getrecords"`
}

// type RecordResponse struct {
// 	ID        int    `json:"id"`
// 	Name      string `json:"name"`
// 	Type      string `json:"type"`
// 	Value     string `json:"value"`
// 	TTL       int    `json:"ttl"`
// 	CreatedAt string `json:"created_at"`
// 	UpdatedAt string `json:"updated_at"`
// }

// CreateUpdateRecordBody api type

type CreateUpdateRecordBody Details

// CreateRecordResponse api type
type CreateRecordResponse struct {
	Record struct {
		Id int `json:"id,omitempty"`
	} `json:"record,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type Credentials struct {
	AccountName string `json:"account_name"`
	ApiKey      string `json:"apikey"`
}

// GetToken получает токен и сохраняет его в ApiToken
func (c *NetangelsClient) GetToken() error {
	requestData := fmt.Sprintf("api_key=%s", c.Credentials.ApiKey)
	resp, err := http.Post(panelUrl, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(requestData)))
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса токена: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка: статус %d, ответ: %s", resp.StatusCode, body)
	}

	var tokenResponse struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return fmt.Errorf("ошибка при декодировании токена в JSON: %v", err)
	}

	c.ApiToken = tokenResponse.Token // Сохраняем токен в ApiToken
	return nil
}

// Получение списка зон
func GetZoneID(zone_name string, c *NetangelsClient) (int, error) {
	// Создание нового HTTP-запроса

	req, err := http.NewRequest("GET", apiUrl+"/dns/zones/", nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Установка заголовка авторизации
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ошибка: статус %d, ответ: %s", resp.StatusCode, body)
	}

	// Декодирование JSON-ответа в структуру Zones
	var zones Zones
	if err := json.Unmarshal(body, &zones); err != nil {
		return 0, fmt.Errorf("ошибка декодирования JSON: %v", err)
	}

	// Ищем ID зоны с именем "k8s-services.ru"
	var zoneID int
	for _, zone := range zones.Entities {
		if zone.Name == zone_name {
			zoneID = zone.ID
			break // Выходим из цикла, если нашли нужную зону
		}
	}

	if zoneID != 0 {
		return zoneID, nil
	} else {
		return 0, fmt.Errorf("запись зоны %s отсутствует", zone_name)
	}
}

// AddRecord Add record to netangels
func (c *NetangelsClient) AddRecord(FQDNName string, Value string, recordType RecordType, ttl int) (int, error) {
	if !validateRecordType(recordType) {
		log.Errorln("invalid record type: ", recordType)
		return 0, errors.New("invalid record type")
	}

	if ttl <= 0 { // Проверяем, если ttl не задан или равен 0
		ttl = 300 // Значение по умолчанию
	}

	TXTRecordBody := TXTDnsRecord{
		Name:  FQDNName,
		Value: Value,
		TTL:   ttl,
	}
	postBody, err := json.Marshal(TXTRecordBody)
	if err != nil {
		return 0, fmt.Errorf("ошибка при маршалинге записи: %v", err)
	}
	req, err := http.NewRequest("POST", apiUrl+"/dns/records/", bytes.NewBuffer(postBody))
	if err != nil {
		return 0, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	//req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer response.Body.Close()

	// Обработка кодов ответов
	switch response.StatusCode {
	case http.StatusCreated: // Код 201
		var createResponse CreateDnsResponse
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return 0, fmt.Errorf("ошибка чтения ответа: %v", err)
		}
		if err := json.Unmarshal(body, &createResponse); err != nil {
			return 0, fmt.Errorf("ошибка декодирования JSON: %v", err)
		}
		return createResponse.ID, nil

	case http.StatusBadRequest: // Код 400
		return 0, fmt.Errorf("неверный формат данных или отсутствуют обязательные параметры")

	default:
		return 0, fmt.Errorf("неожиданный код ответа: %d", response.StatusCode)
	}
}

// RemoveRecord Remove record from netangels
// func (c *NetangelsClient) RemoveRecord(ID int, DnsName string) error {
func (c *NetangelsClient) RemoveRecord(ID int) error {
	//dnsName := cutTrailingDotIfExist(DnsName)
	req, err := http.NewRequest("DELETE", apiUrl+"/dns/records/"+strconv.Itoa(ID)+"/", nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		fmt.Println("Запись успешно удалена.")
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("ошибка: DNS-запись с ID %d не найдена", ID)
	default:
		return fmt.Errorf("ошибка: статус %d, ответ: %s", response.StatusCode, body)
	}
}

// GetRecord Fetch record by FQDNName, RecordData and RecordType,TTL and priority are ignored, returns id of first record found.
func (c *NetangelsClient) GetRecordID(FQDNName string, RecordData string, recordType RecordType) (int, error) {
	fqdnName := cutTrailingDotIfExist(FQDNName)
	responseData, err := getRecords(fqdnName, c)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения списка записей домена: %v", err)
	}

	// Ищем запись в списке записей полученых из DNS"
	var record_id int
	for _, item := range responseData.Entities {
		if item.Type == recordType {
			if item.Name == RecordData {
				record_id = item.ID
				break // Выходим из цикла, если нашли нужную запись
			}
		}
	}

	return record_id, nil

	// if err == nil {
	// 	for i := 0; i < len(records.Records); i++ {
	// 		if records.Records[i].Value == RecordData && records.Records[i].Type == string(recordType) && records.Records[i].Name == domainutil.Subdomain(fqdnName) {
	// 			return records.Records[i].ID, records.Records[i].Value, nil
	// 		}
	// 	}
	// 	if RecordData == "" {
	// 		for i := 0; i < len(records.Records); i++ {
	// 			if records.Records[i].Type == string(recordType) && records.Records[i].Name == domainutil.Subdomain(fqdnName) {
	// 				return records.Records[i].ID, records.Records[i].Value, nil
	// 			}
	// 		}
	// 	}
	// 	log.Errorln("Record not found")
	// 	return 0, "", errors.New("record not found")
	// } else {
	// 	log.Errorln("Failed to unmarshal body with error: ", err)
	// 	return 0, "", err
	// }
}

// GetRecords Fetch records by FQDNName in struct GetRecordsList
func getRecords(FQDNName string, c *NetangelsClient) (GetRecordsList, error) {
	fqdnName := cutTrailingDotIfExist(FQDNName)
	zone_id, err := GetZoneID(fqdnName, c)
	if err != nil {
		return GetRecordsList{}, err
	}
	req, err := http.NewRequest("GET", apiUrl+"/dns/zones/"+strconv.Itoa(zone_id)+"/records/", nil)
	if err != nil {
		return GetRecordsList{}, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	// Установка заголовка авторизации
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		log.Errorln("Failed request, response code: ", response.StatusCode)
		return GetRecordsList{}, err
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorln("Error on read: ", err)
		return GetRecordsList{}, err
	}

	// Декодирование JSON-ответа в структуру Zones
	var recoed_list GetRecordsList
	if err := json.Unmarshal(responseData, &recoed_list); err != nil {
		return GetRecordsList{}, fmt.Errorf("ошибка декодирования JSON: %v", err)
	}

	return recoed_list, nil
}

// UpdateRecord Update record by ID, FQDNName, Value and RecordType
// func (c *NetangelsClient) UpdateRecord(ID int, FQDNName string, Value string, recordType RecordType) (bool, error) {
// 	if !validateRecordType(recordType) {
// 		log.Errorln("Invalid record type: ", recordType)
// 		return false, errors.New("invalid record type")
// 	}
// 	// Trim one trailing dot
// 	fqdnName := cutTrailingDotIfExist(FQDNName)
// 	TXTRecordBody := CreateUpdateRecordBody{
// 		Type:     recordType,
// 		Name:     domainutil.Subdomain(fqdnName),
// 		Value:    Value,
// 		Priority: 1,
// 		TTL:      3600,
// 	}
// 	putBody, _ := json.Marshal(TXTRecordBody)
// 	req, err := http.NewRequest("PUT", apiUrl+"/my/products/"+domainutil.Domain(fqdnName)+"/dns/records/"+strconv.Itoa(ID), bytes.NewBuffer(putBody))
// 	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	req.SetBasicAuth(c.Credentials.AccountName, c.Credentials.ApiKey)
// 	client := &http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil || response.StatusCode != 200 {
// 		log.Errorln("Failed request, response code: ", response.StatusCode)
// 		return false, err
// 	}
// 	return true, nil
// }

// // DDNS update/create record if Ip omited Simply api will use client IP
// func (c *NetangelsClient) UpdateDDNS(FQDNName string, Ip string) (bool, error) {
// 	fqdnName := cutTrailingDotIfExist(FQDNName)
// 	domain := domainutil.Domain(fqdnName)
// 	path := ""
// 	if Ip == "" {
// 		path = "/ddns/?domain=" + domain + "&hostname=" + fqdnName
// 	} else {
// 		path = "/ddns/?domain=" + domain + "&hostname=" + fqdnName + "&myip=" + Ip
// 	}
// 	//body, _ := json.Marshal("")
// 	req, err := http.NewRequest("POST", apiUrl+path, nil)
// 	req.SetBasicAuth(c.Credentials.AccountName, c.Credentials.ApiKey)
// 	client := &http.Client{}
// 	response, err := client.Do(req)

// 	if err != nil || response.StatusCode != 200 {
// 		log.Errorln("Failed request, response code: ", response.StatusCode)
// 		return false, err
// 	}
// 	return true, nil
// }

func cutTrailingDotIfExist(FQDNName string) string {
	fqdnName := FQDNName
	if last := len(fqdnName) - 1; last >= 0 && fqdnName[last] == '.' {
		fqdnName = fqdnName[:last]
	}
	return fqdnName
}
