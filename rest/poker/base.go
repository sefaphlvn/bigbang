package poker

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoFilters struct {
	Collection string
	Filter     bson.D
}

type Processed struct {
	Listeners []string
	Depends   []string
}

func ReSnapshot(listenerName string, context *db.AppContext) {
	baseURL := fmt.Sprintf("http://%s/poke", context.Config.BIGBANG_ADDRESS)

	params := url.Values{}
	params.Add("service", listenerName)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		context.Logger.Debugf("Creating request failed: %s\n", err)
		return
	}

	req.Header.Set("bigbang-controller", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		context.Logger.Debugf("HTTP request failed: %s\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		context.Logger.Debugf("Request successful: %s\n", resp.Status)
	} else {
		context.Logger.Debugf("HTTP request returned status code: %d\n", resp.StatusCode)
	}
}

func DetectChangedResource(gType models.GTypes, resourceName string, context *db.AppContext, processed *Processed) *Processed {
	pathWithGtype := gType.String() + "===" + resourceName
	if gType != models.Listener {
		processed.Depends = append(processed.Depends, pathWithGtype)
	}

	if handler, exists := handlers[gType]; exists {
		handler.Handle(context, resourceName, processed)
	} else if gType == models.Listener {
		if !helper.Contains(processed.Listeners, resourceName) {
			ReSnapshot(resourceName, context)
			processed.Listeners = append(processed.Listeners, resourceName)

			result := strings.Join(processed.Depends, " \n ")
			context.Logger.Infof("new version added to snapshot for (%s) processed resource paths: \n %s", resourceName, result)
		}
	} else {
		context.Logger.Infof("not covered gtype: %s", gType)
	}

	return processed
}
