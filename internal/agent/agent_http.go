package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
)

type Http struct {
	cfg *Config
}

func NewHttp(cfg *Config) *Http {
	return &Http{cfg: cfg}
}

func (h Http) Close() error {
	return nil
}

func (h Http) SendMetricsByOne(m *memory.Metrics) error {
	postURL := url.URL{
		Scheme: "http",
		Host:   h.cfg.Address,
		Path:   "/update/",
	}

	for _, value := range m.Collection {
		value.Hash = value.HashValue(h.cfg.HashKey)

		processedMetrics, err := json.Marshal(value)
		if err != nil {
			return errors.New("couldn't marshal metric")
		}
		processedMetrics, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, h.cfg.CryptoKey, processedMetrics, nil)
		if err != nil {
			return fmt.Errorf("couldn't encrypt metric: %s", err)
		}

		gzipMetrics := bytes.NewBuffer([]byte{})
		w := gzip.NewWriter(gzipMetrics)
		_, err = w.Write(processedMetrics)
		if err != nil {
			return fmt.Errorf("couldn't gzip metric with error: %s", err)
		}

		err = w.Close()
		if err != nil {
			return fmt.Errorf("couldn't close gzip writer with error: %s", err)
		}

		request, err := http.NewRequest(http.MethodPost, postURL.String(), gzipMetrics)
		if err != nil {
			return fmt.Errorf("couldn't create request with error: %s", err)
		}

		request.Header.Set("Content-Encoding", "gzip")

		ip, err := getOutboundIP()
		if err != nil {
			return fmt.Errorf("couldn't get outbound ip: %s", err.Error())
		}

		request.Header.Set("X-Real-IP", ip.String())

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return fmt.Errorf("couldn't post url %s", postURL.String())
		}

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("coudln't send metrics, status code: %d, response: %s", resp.StatusCode, body)
		}

		if err := resp.Body.Close(); err != nil {
			return errors.New("couldn't close response body")
		}
	}

	return nil
}

func (h Http) SendMetricsAllTogether(m *memory.Metrics) error {
	postURL := url.URL{
		Scheme: "http",
		Host:   h.cfg.Address,
		Path:   "/updates/",
	}

	metrics := make([]models.Metric, 0, len(m.Collection))
	for _, value := range m.Collection {
		value.Hash = value.HashValue(h.cfg.HashKey)
		metrics = append(metrics, value)
	}

	processedMetrics, err := json.Marshal(metrics)
	if err != nil {
		return errors.New("couldn't marshal metrics")
	}

	gzipMetrics := bytes.NewBuffer([]byte{})
	w := gzip.NewWriter(gzipMetrics)
	_, err = w.Write(processedMetrics)
	if err != nil {
		return fmt.Errorf("couldn't gzip metrics with error: %s", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("couldn't close gzip writer with error: %s", err)
	}

	request, err := http.NewRequest(http.MethodPost, postURL.String(), gzipMetrics)
	if err != nil {
		return fmt.Errorf("couldn't create request with error: %s", err)
	}
	request.Header.Set("Content-Encoding", "gzip")

	ip, err := getOutboundIP()
	if err != nil {
		return fmt.Errorf("couldn't get outbound ip: %s", err.Error())
	}

	request.Header.Set("X-Real-IP", ip.String())

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("couldn't post url %s", postURL.String())
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("coudln't send metrics, status code: %d, response: %s", resp.StatusCode, body)
	}

	if err := resp.Body.Close(); err != nil {
		return errors.New("couldn't close response body")
	}

	return err
}
