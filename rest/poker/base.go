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
	baseURL := fmt.Sprintf("http://%s/poke", wtf.Config.GrpcService)

	params := url.Values{}
	params.Add("service", listenerName)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		wtf.Logger.Debugf("HTTP request failed: %s\n", err)
	}
	if resp != nil {
		defer resp.Body.Close()
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
