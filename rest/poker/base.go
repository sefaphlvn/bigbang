package poker

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

type Processed struct {
	Listeners []string
	Depends   []string
}

var processed = Processed{Listeners: []string{}, Depends: []string{}}

func ReSnapshot(listenerName string, wtf *db.WTF) {
	baseURL := "http://localhost/poke"

	params := url.Values{}
	params.Add("service", listenerName)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Yeni bir HTTP isteği oluşturun
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		wtf.Logger.Debugf("Creating request failed: %s\n", err)
		return
	}

	// İstek nesnesine header ekleyin
	req.Header.Set("bigbang-controller", "1")

	// Bir HTTP client kullanarak isteği gönderin
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wtf.Logger.Debugf("HTTP request failed: %s\n", err)
		return
	}

	defer resp.Body.Close()

	// İsteğin yanıtını işleyin
	if resp.StatusCode == http.StatusOK {
		// İşlem başarılı, yanıtı okuyun veya loglayın
		wtf.Logger.Debugf("Request successful: %s\n", resp.Status)
	} else {
		// Yanıttaki hata durumlarını ele alın
		wtf.Logger.Debugf("HTTP request returned status code: %d\n", resp.StatusCode)
	}
}

func DetectChangedResource(gType models.GTypes, resourceName string, wtf *db.WTF) {
	if gType != models.Listener {
		pathWithGtype := gType.String() + "===" + resourceName
		processed.Depends = append(processed.Depends, pathWithGtype)
	}

	switch gType {
	case models.Endpoint:
		PokerEds(wtf, resourceName)
	case models.Cluster:
		PokerCds(wtf, resourceName)
	case models.Router:
		PokerRouter(wtf, resourceName)
	case models.Route:
		PokerRoute(wtf, resourceName)
	case models.HTTPConnectionManager:
		PokerHCM(wtf, resourceName)
	case models.TcpProxy:
		PokerTcpProxy(wtf, resourceName)
	case models.DownstreamTlsContext, models.UpstreamTlsContext, models.TlsCertificate, models.CertificateValidationContext:
		PokerTLS(wtf, resourceName, gType)
	case models.Listener:
		if !helper.Contains(processed.Listeners, resourceName) {
			ReSnapshot(resourceName, wtf)
			processed.Listeners = append(processed.Listeners, resourceName)
			result := strings.Join(processed.Depends, " \n ")
			wtf.Logger.Infof("new version added to snapshot for (%s) processed resource paths: \n %s", resourceName, result)
		}
	default:
		wtf.Logger.Infof("not covered gtype: %s", gType)
	}
}

func ResetProcessedResources() {
	processed = Processed{Listeners: []string{}, Depends: []string{}}
}
