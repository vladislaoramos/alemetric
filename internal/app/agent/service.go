package agent

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/alemetric/internal/entity"
)

// WebAPIClient implements the client web-application for Agent.
type WebAPIClient struct {
	client    *resty.Client
	Key       string
	publicKey string
}

func NewWebAPI(client *resty.Client, key string, cryptoKey string) *WebAPIClient {
	return &WebAPIClient{
		client:    client,
		Key:       key,
		publicKey: cryptoKey,
	}
}

// SendMetrics sends a client request for a metrics update to the server.
func (wc *WebAPIClient) SendMetrics(
	metricsName,
	metricsType string,
	delta *entity.Counter,
	value *entity.Gauge,
) error {
	body := entity.Metrics{
		ID:    metricsName,
		MType: metricsType,
		Delta: delta,
		Value: value,
	}

	body.SignData("agent", wc.Key)

	// respBody, _ := json.Marshal(body)
	// log.Printf("req body: %s", string(respBody))

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	publicKey, err := loadPublicKeyFromFile(wc.publicKey)
	if err != nil {
		return err
	}

	encryptedBody := tryEncrypt(b, publicKey)

	resp, err := wc.client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(encryptedBody).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("cannot send metrics from agent: %w", err)
	}

	status := resp.StatusCode()
	if status != http.StatusOK {
		return fmt.Errorf("sending metrics from agent with not successful status code: %d", status)
	}

	return nil
}

// SendSeveralMetrics sends a client request for several metrics update to the server.
func (wc *WebAPIClient) SendSeveralMetrics(items []entity.Metrics) error {
	resp, err := wc.client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(items).
		Post("/updates/")
	if err != nil {
		return fmt.Errorf("cannot send several metrics from agent: %w", err)
	}

	status := resp.StatusCode()
	if status != http.StatusOK {
		return fmt.Errorf("sending several metrics from agent with not successful status code: %d", status)
	}
	return nil
}

func tryEncrypt(msg []byte, key *rsa.PublicKey) []byte {
	if key == nil {
		return msg
	}

	hash := sha512.New()

	result, err := rsa.EncryptOAEP(hash, rand.Reader, key, msg, nil)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func loadPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected public key type")
	}

	return rsaPubKey, nil
}
