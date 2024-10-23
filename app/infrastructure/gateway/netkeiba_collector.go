package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"net/http"
	"os"
	"path/filepath"
)

type NetKeibaCollector interface {
	Client() *colly.Collector
	Login(ctx context.Context) error
}

type netKeibaCollector struct {
	client *colly.Collector
}

const (
	netKeibaBaseUrl     = "https://www.netkeiba.com"
	collectorConfigName = "netkeiba_collector_config.json"
	collyCacheDir       = "colly"
)

func NewNetKeibaCollector() NetKeibaCollector {
	client := colly.NewCollector()
	client.AllowURLRevisit = true
	client.DetectCharset = true

	rootPath, _ := os.Getwd()
	cachePath, _ := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, collyCacheDir))
	client.CacheDir = cachePath

	return &netKeibaCollector{
		client: client,
	}
}

func (n *netKeibaCollector) Client() *colly.Collector {
	return n.client
}

func (n *netKeibaCollector) Login(ctx context.Context) error {
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, collectorConfigName))
	if err != nil {
		return err
	}

	configBytes, err := os.ReadFile(secretFilePath)
	if err != nil {
		return err
	}

	var rawNetKeibaCollectorConfig raw_entity.NetKeibaCollectorConfigs
	if err = json.Unmarshal(configBytes, &rawNetKeibaCollectorConfig); err != nil {
		return err
	}

	var cookies []*http.Cookie
	for _, data := range rawNetKeibaCollectorConfig.Cookies {
		cookies = append(cookies, &http.Cookie{
			Name:     data.Name,
			Value:    data.Value,
			Path:     data.Path,
			Domain:   data.Domain,
			Secure:   data.Secure,
			HttpOnly: data.HttpOnly,
		})
	}
	err = n.client.SetCookies(netKeibaBaseUrl, cookies)
	if err != nil {
		return err
	}

	ok, err := n.isLoggedIn()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("netkeiba login failed")
	}

	return nil
}

func (n *netKeibaCollector) isLoggedIn() (bool, error) {
	loggedIn := false
	n.client.OnHTML("span.header_nickname", func(e *colly.HTMLElement) {
		loggedIn = true
	})
	n.client.OnError(func(r *colly.Response, err error) {
		loggedIn = false
	})

	err := n.client.Visit(netKeibaBaseUrl)
	if err != nil {
		return false, err
	}

	return loggedIn, nil
}