package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"https://github.com/infa-kmoholka/CDGC_Regression/config"
	"github.com/olivere/elastic/v7"
)

//GetReleaseData ...
func GetReleaseData(buildNum string, release string, env string, iteration string, service string) (map[string]*config.TimesResponse, map[string]*config.TimesResponse) {

	//Hostname := "irl62dqd07"
	conf := ReadConfig()

	client, err := elastic.NewClient(
		elastic.SetURL(conf.ElasticURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {

	}

	//ping to check connectivity

	info, code, err := client.Ping(conf.ElasticURL).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Get doc for the specific buildnumber
	filterByBuildQuery := elastic.NewTermQuery("BuildNumber", buildNum)
	filterByReleaseQuery := elastic.NewTermQuery("ReleaseNumber", release)
	//filterByLabelQuery:=elastic.NewRegexpQuery("Label","*REC*")
	searchQuery := elastic.NewTermQuery("Environment.keyword", env)
	filterByIterationQuery := elastic.NewTermQuery("Iteration", iteration)
	filterByServiceQuery := elastic.NewTermQuery("ServiceName.keyword", service)
	//fmt.Println(filterByBuildQuery,filterByReleaseQuery,searchQuery,filterByIterationQuery,filterByServiceQuery)
	filterQuery := elastic.NewBoolQuery().Must(filterByReleaseQuery).Must(filterByBuildQuery).Must(searchQuery).Must(filterByIterationQuery).Must(filterByServiceQuery)

	//for filter based on last build num use aggregation max with release

	SearchResult, err := client.Search().
		Index(conf.ElasticSearchReportIndex). // search in index mentioned in config file
		Query(filterQuery).
		From(0).Size(1000).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	if SearchResult.Hits.TotalHits.Value > 0 {
		fmt.Printf("Found a total of %d hits\n", SearchResult.Hits.TotalHits.Value)

		var t config.TimesResponse

		//declaring the map for storing the data returned from elasticsearch
		newTaskMap := make(map[string]*config.TimesResponse) //key is of type string and value os of type *config.TimesResponse
		newMap := make(map[string]*config.TimesResponse)

		for _, hit := range SearchResult.Hits.Hits {

			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
				fmt.Printf("%s", "Error during deserialization")
			}

			//fmt.Println("label: ",t.Label,"average: ",t.Average,"percentile: ",t.Percentile99)

			//key := strings.Split(string(t.ResourceName), "_")
			//key := t.Scenario + "__" + string(t.Iteration)
			key := t.Label
			//Printing only selected API
			if strings.Contains(key, "HAWK_TXN") {
				//fmt.Println(tasktimes)

				newMap[key] = &config.TimesResponse{Label: t.Label, Average: t.Average, Percentile99: t.Percentile99}
				//dummy map to be replaced in future if complete report is required.
				newTaskMap[key] = &config.TimesResponse{Label: t.Label, Average: t.Average, Percentile99: t.Percentile99}

			}

		}

		return newMap, newTaskMap

	}

	// No hits
	msg := fmt.Sprintf("Found no hits for buildNumber: %s", buildNum)

	fmt.Println(msg)

	return nil, nil

}
