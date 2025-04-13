package services

import (
	"PrytkovaBot/config"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sort"
)

const BaseUrl = "https://securepay.tinkoff.ru/v2/"

func InitTransaction(amount float64) (string, string, error) {
	// Параметры запроса
	baseParams := map[string]string{
		"TerminalKey": config.TerminalData.TerminalKey,
		"Amount":      fmt.Sprint(amount * 100),
		"OrderId":     uuid.New().String(),
	}

	// Добавляем пароль для токена
	tokenParams := make(map[string]string)
	for k, v := range baseParams {
		tokenParams[k] = v
	}
	tokenParams["Password"] = config.TerminalData.TerminalPassword
	baseParams["Token"] = createToken(tokenParams)

	// Преобразуем параметры в JSON
	jsonData, err := json.Marshal(baseParams)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(BaseUrl+"Init", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Декодирование JSON
	var result map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	if url, ok := result["PaymentURL"].(string); !ok {
		return "", "", fmt.Errorf("нет PaymentURL")
	} else {
		return url, result["PaymentId"].(string), nil
	}
}

func GetState(paymentId string) (string, error) {
	baseParams := map[string]string{
		"TerminalKey": config.TerminalData.TerminalKey,
		"PaymentId":   paymentId,
	}
	tokenParams := make(map[string]string)
	for k, v := range baseParams {
		tokenParams[k] = v
	}
	tokenParams["Password"] = config.TerminalData.TerminalPassword
	baseParams["Token"] = createToken(tokenParams)
	jsonData, err := json.Marshal(baseParams)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(BaseUrl+"GetState", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	if _, ok := result["Status"].(string); !ok {
		return "", fmt.Errorf("нет Status")
	} else {
		return result["Status"].(string), nil
	}

}

func createToken(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	concatenated := ""
	for _, k := range keys {
		concatenated += params[k]
	}
	hash := sha256.Sum256([]byte(concatenated))
	return hex.EncodeToString(hash[:])
}
