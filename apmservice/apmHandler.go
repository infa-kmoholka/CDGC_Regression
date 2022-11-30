package apmservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/infa-kmoholka/CDGC_Regression/config"
	"github.com/infa-kmoholka/CDGC_Regression/utils"
)

// test method to validate api is up and running
func test(w http.ResponseWriter, r *http.Request) {

	body := config.Body{ResponseCode: 200, Message: "OK"}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

}

// handler function for release comparison
func compareRelease(w http.ResponseWriter, r *http.Request) {

	//storing the query params as variables
	oldBuildNum := r.URL.Query().Get("oldBuildNum") //427.4
	oldRelease := r.URL.Query().Get("oldRelease")   // "10.2.2"
	newBuildNum := r.URL.Query().Get("newBuildNum") //427.5
	newRelease := r.URL.Query().Get("newRelease")
	Environment := r.URL.Query().Get("environment")
	Iteration := r.URL.Query().Get("iteration")
	ServiceName := r.URL.Query().Get("serviceName")
	Scenario := r.URL.Query().Get("scenario")
	cc := r.URL.Query().Get("email")
	metric := r.URL.Query().Get("metric")

	//storing the map returned from GetReleaseData method containing the data returned from elasticsearch
	oldReleaseData, _ := utils.GetReleaseData(oldBuildNum, oldRelease, Environment, Iteration, ServiceName)
	newReleaseData, _ := utils.GetReleaseData(newBuildNum, newRelease, Environment, Iteration, ServiceName)

	//fmt.Println(oldReleaseData)
	//checking if the data returned is empty
	if (len(newReleaseData) == 0) || (len(oldReleaseData) == 0) {

		utils.RespondWithJSON("BuildNumber/Release not correct or not enough data ", w, r)

	} else {
		//declaring subject to be used in email report
		subject := fmt.Sprintf("Release Comparison Report for %s (%s) & %s (%s)", oldRelease, oldBuildNum, newRelease, newBuildNum)
		//declaring the header to be used in html report
		p := fmt.Sprintf("<body style='background:white'><h3 style='background:#0790bd;color:#fff;padding:5px;text-align:center;border-radius:5px;'> Release Comparison for %s (%s) & %s (%s) </h3> <br/> <br/>", oldRelease, oldBuildNum, newRelease, newBuildNum)
		p = p + fmt.Sprintf("<div style='background:yellow;text-align:center'><p><b>Scenario Summary : %s </p> </b></div>", Scenario)

		//updating the variable with the table body declaration
		p = p + fmt.Sprintf("<table style='backgound:#fff;border-collapse: collapse;' border = '1' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:#444;color:white;'><b>Scenario : %s | Iteration : %s </b></td></tr><tr><th>Transaction</th><th>Release: %s </th><th>Release: %s</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", Scenario, Iteration, oldRelease, newRelease)
		//sort old release data
		oldReleaseDataSorted := utils.SortingMap(oldReleaseData)
		//iterating over the map values of old data and storing the key/value pair as Label/_
		for _, Label := range oldReleaseDataSorted {

			//checking if the label key is same for both old and new release. If not same then comparison won't happen
			if _, ok := newReleaseData[Label]; ok {

				//checking if the newReleaseData label is not nil
				if newReleaseData[Label] != nil {

					//fmt.Println(Label)
					var timeOld int64
					var timeNew int64
					//timeNew := 0
					if strings.Contains(metric, "average") {
						timeOld = oldReleaseData[Label].Average
						timeNew = newReleaseData[Label].Average
					} else if strings.Contains(metric, "99") {
						timeOld = oldReleaseData[Label].Percentile99
						timeNew = newReleaseData[Label].Percentile99
					}
					diff := timeNew - timeOld
					//if difference is same or less then show it as green using color code :#80CA80
					if diff <= 0 {
						percDiff := utils.CalcPerc(diff, timeOld)
						p = p + "<tr style='background:#80CA80'><td>" + Label + "</td><td>" + strconv.FormatInt(timeOld, 10) + "</td><td>" + strconv.FormatInt(timeNew, 10) + "</td><td>" + strconv.FormatInt(diff, 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"

					} else {
						//else if difference is more then show it as red using color code :#ff9e82
						percDiff := utils.CalcPerc(diff, timeOld)
						p = p + "<tr style='background:#ff9e82'><td>" + Label + "</td><td>" + strconv.FormatInt(timeOld, 10) + "</td><td>" + strconv.FormatInt(timeNew, 10) + "</td><td>" + strconv.FormatInt(diff, 10) + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
					}

				}

			}

		}

		p = p + "</tbody></table></body><br/><br/>"
		//fmt.Println(p)
		//read the config from the ReadConfig method
		conf := utils.ReadConfig()
		p = p + "<b>Dashboard URL : </b>" + conf.DashboardURL
		//send the mail by passing the body variable-p and subject and cc
		utils.SendMail(p, subject, cc)

		//write to file
		fileName := conf.HtmlFolderPath + oldRelease + "_" + oldBuildNum + "vs" + newRelease + "_" + newBuildNum + "_" + Iteration + ".html"
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		_, err = f.WriteString(p)
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Println(p)
		utils.RespondWithJSON("Email Sent Successfully", w, r)
	}

}

func htmlReport(w http.ResponseWriter, r *http.Request) {

	p := "./" + r.URL.Path
	//fmt.Println(r.URL.Path)

	http.ServeFile(w, r, p)

}
