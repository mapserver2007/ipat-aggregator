package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/file_gateway"
)

type NetKeibaCollector interface {
	Client() *colly.Collector
	Cookies(ctx context.Context) ([]*http.Cookie, error)
	Cache(c bool) bool
	Login(ctx context.Context) error
}

type netKeibaCollector struct {
	client        *colly.Collector
	pathOptimizer file_gateway.PathOptimizer
}

const (
	netKeibaBaseUrl     = "https://www.netkeiba.com"
	collectorConfigName = "netkeiba_collector_config.json"
	collyCacheDir       = "colly"
)

func NewNetKeibaCollector(
	pathOptimizer file_gateway.PathOptimizer,
) NetKeibaCollector {
	client := colly.NewCollector()
	client.AllowURLRevisit = true
	client.DetectCharset = true
	collector := &netKeibaCollector{
		client:        client,
		pathOptimizer: pathOptimizer,
	}
	collector.Cache(true)

	return collector
}

func (n *netKeibaCollector) Client() *colly.Collector {
	return n.client
}

func (n *netKeibaCollector) Cache(c bool) bool {
	if c {
		rootPath, err := n.pathOptimizer.GetProjectRoot()
		if err != nil {
			return false
		}
		cachePath, _ := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, collyCacheDir))
		n.client.CacheDir = cachePath
	} else {
		n.client.CacheDir = ""
	}

	return true
}

func (n *netKeibaCollector) Cookies(ctx context.Context) ([]*http.Cookie, error) {
	rootPath, err := n.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, collectorConfigName))
	if err != nil {
		return nil, err
	}

	configBytes, err := os.ReadFile(secretFilePath)
	if err != nil {
		return nil, err
	}

	var rawNetKeibaCollectorConfig raw_entity.NetKeibaCollectorConfigs
	if err = json.Unmarshal(configBytes, &rawNetKeibaCollectorConfig); err != nil {
		return nil, err
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

	return cookies, nil
}

func (n *netKeibaCollector) Login(ctx context.Context) error {
	cookies, err := n.Cookies(ctx)
	if err != nil {
		return err
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
