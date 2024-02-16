package poker

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

var ProcessedResources = make(map[string]bool)

func ReSnapshot(listenerName string, wtf *db.WTF) {
	baseURL := fmt.Sprintf("http://%s/poke", wtf.Config.GrpcService)

	params := url.Values{}
	params.Add("service", listenerName)

	wtf.Logger.Infof("new version added to snapshot: %s", listenerName)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("HTTP request failed: %s\n", err)
	}
	defer resp.Body.Close()
}

func DetectChangedResource(gType models.GTypes, resourceName string, wtf *db.WTF) {
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
	case models.Listener:
		if !ProcessedResources[resourceName] {
			ReSnapshot(resourceName, wtf)
			ProcessedResources[resourceName] = true
		}
	default:
		fmt.Println("sss")
	}
}

func ResetProcessedResources() {
	ProcessedResources = make(map[string]bool)
}
