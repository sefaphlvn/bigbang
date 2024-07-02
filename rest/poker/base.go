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

func ReSnapshot(listenerName string, context *db.AppContext) {
	baseURL := fmt.Sprintf("http://%s/poke", context.Config.BIGBANG_ADDRESS)

	params := url.Values{}
	params.Add("service", listenerName)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Yeni bir HTTP isteği oluşturun
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		context.Logger.Debugf("Creating request failed: %s\n", err)
		return
	}

	// İstek nesnesine header ekleyin
	req.Header.Set("bigbang-controller", "1")

	// Bir HTTP client kullanarak isteği gönderin
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		context.Logger.Debugf("HTTP request failed: %s\n", err)
		return
	}

	defer resp.Body.Close()

	// İsteğin yanıtını işleyin
	if resp.StatusCode == http.StatusOK {
		// İşlem başarılı, yanıtı okuyun veya loglayın
		context.Logger.Debugf("Request successful: %s\n", resp.Status)
	} else {
		// Yanıttaki hata durumlarını ele alın
		context.Logger.Debugf("HTTP request returned status code: %d\n", resp.StatusCode)
	}
}

func DetectChangedResource(gType models.GTypes, resourceName string, context *db.AppContext) {
	if gType != models.Listener {
		pathWithGtype := gType.String() + "===" + resourceName
		processed.Depends = append(processed.Depends, pathWithGtype)
	}

	switch gType {
	case models.Endpoint:
		PokerEds(context, resourceName)
	case models.Cluster:
		PokerCds(context, resourceName)
	case models.Router:
		PokerRouter(context, resourceName)
	case models.Route:
		PokerRoute(context, resourceName)
	case models.HTTPConnectionManager:
		PokerHCM(context, resourceName)
	case models.TcpProxy:
		PokerTcpProxy(context, resourceName)
	case models.DownstreamTlsContext, models.UpstreamTlsContext, models.TlsCertificate, models.CertificateValidationContext:
		PokerTLS(context, resourceName, gType)
	case models.Listener:
		if !helper.Contains(processed.Listeners, resourceName) {
			ReSnapshot(resourceName, context)
			processed.Listeners = append(processed.Listeners, resourceName)
			result := strings.Join(processed.Depends, " \n ")
			context.Logger.Infof("new version added to snapshot for (%s) processed resource paths: \n %s", resourceName, result)
		}
	default:
		context.Logger.Infof("not covered gtype: %s", gType)
	}
}

func ResetProcessedResources() {
	processed = Processed{Listeners: []string{}, Depends: []string{}}
}
