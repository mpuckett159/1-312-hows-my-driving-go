package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	soda "github.com/SebastiaanKlippert/go-soda"
	"github.com/gin-gonic/gin"

	privateTemplates "1-312-hows-my-driving-go/templates"
)

// License : License Structure
type License struct {
	License       string
	Make          string
	Model         string
	Dept          string
	Descrip       string
	EquipmentType string
}

// Officer : Officer structure
type Officer struct {
	Badge      string
	Name       string
	Rank       string
	UnitDesc   string
	JobTitle   string
	HourlyRate float64
	ProjSalary float64
}

// LicenseRender : Rendering data for License lookup view
var LicenseRender = gin.H{
	"title":             "Seattle Public Vehicle Lookup",
	"entity_name_long":  "license plate",
	"entity_name_short": "License #",
	"data_source":       "https://data.seattle.gov/resource/enxu-fgzb",
	"lookup_url":        "license",
}

// BadgeRender : Rendering data for Badge lookup view
var BadgeRender = gin.H{
	"title":             "Seattle Officer Badge Lookup",
	"entity_name_long":  "badge number",
	"entity_name_short": "Badge #",
	"data_source":       "https://data.seattle.gov/resource/2khk-5ukd",
	"lookup_url":        "badge",
}

// NameRender : Rendering data for Name lookup view
var NameRender = gin.H{
	"title":             "Seattle Officer Name Lookup",
	"entity_name_long":  "last name",
	"entity_name_short": "Last name",
	"data_source":       "https://data.seattle.gov/resource/2khk-5ukd",
	"lookup_url":        "name",
}

// Convert url.Values type HTTP GET queries to map[string]string
func convertValuesToMap(queryParams url.Values) (queryMap map[string]string) {
	queryMap = make(map[string]string)
	for key, value := range queryParams {
		queryMap[key] = value[0]
	}
	return
}

// SODA Interface
func sodaQuery(datasetURL string, queryFilterMap map[string]string) (results []map[string]interface{}, err error) {
	sodareq := soda.NewGetRequest(datasetURL, "")
	sodareq.Format = "json"
	sodareq.Filters = queryFilterMap
	sodareq.Query.Limit = 20

	resp, err := sodareq.Get()
	if err != nil {
		fmt.Println("Error getting data")
	}
	defer resp.Body.Close()

	results = make([]map[string]interface{}, 0)
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		fmt.Println("There was an error with the SODA API")
		return nil, fmt.Errorf("SODA API: %d", err)
	}
	return results, nil
}

// Format license HTML
func formatLicenseHTML(queryResults []map[string]interface{}) (templateHTML template.HTML, err error) {
	var resultSlices []string
	for _, row := range queryResults {
		Data := License{
			License:       row["license"].(string),
			Make:          row["make"].(string),
			Model:         row["model"].(string),
			Dept:          row["dept"].(string),
			Descrip:       row["descrip"].(string),
			EquipmentType: row["equipment_type"].(string),
		}
		templateBytes := new(bytes.Buffer)
		if err := privateTemplates.LicenseTemplate.Execute(templateBytes, Data); err != nil {
			fmt.Printf("Error processing template: %s", err.Error())
		}
		resultSlices = append(resultSlices, templateBytes.String())
	}
	return template.HTML(strings.Join(resultSlices, "\n<br/>\n")), nil
}

// Format badge HTML
func formatOfficerHTML(mergedResults []map[string]string) (templateHTML template.HTML, err error) {
	var resultSlices []string
	for _, row := range mergedResults {
		fullName := strings.Join([]string{row["FirstName"], row["MiddleInitMostly"], row["Surname"]}, " ")
		hourlyRate, _ := strconv.ParseFloat(row["hourly_rate"], 32)
		projSalary := hourlyRate * 2000
		Data := Officer{
			Badge:      row["Serial"],
			Name:       fullName,
			Rank:       row["RankRole"],
			UnitDesc:   row["UnitDesc"],
			JobTitle:   row["job_title"],
			HourlyRate: hourlyRate,
			ProjSalary: projSalary,
		}
		templateBytes := new(bytes.Buffer)
		if err := privateTemplates.LicenseTemplate.Execute(templateBytes, Data); err != nil {
			fmt.Printf("Error processing template: %s", err.Error())
		}
		resultSlices = append(resultSlices, templateBytes.String())
	}
	return template.HTML(strings.Join(resultSlices, "\n<br/>\n")), nil
}

// URL Handler Functions
func redirectToLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/license")
		c.Next()
	}
}

func renderLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		licenseRenderLocal := LicenseRender
		queryParams := c.Request.URL.Query()
		if len(queryParams) > 0 {
			queryMap := convertValuesToMap(queryParams)
			queryResults, err := sodaQuery(LicenseRender["data_source"].(string), queryMap)
			if err != nil {
				fmt.Printf("Error querying SODA API: %s", err)
			}
			resultHTML, err := formatLicenseHTML(queryResults)
			if err != nil {
				fmt.Println("Error rendering license result data into template")
			}
			licenseRenderLocal["entityHTML"] = resultHTML
		}
		c.HTML(http.StatusOK, "index.html", licenseRenderLocal)
	}
}

func renderBadge() gin.HandlerFunc {
	return func(c *gin.Context) {
		badgeRenderLocal := BadgeRender
		/*
			queryParams := c.Request.URL.Query()
			if len(queryParams) > 0 {
				queryMap := convertValuesToMap(queryParams)
				csvResults := data.SearchCSVMap(queryMap)
				queryResults, err := sodaQuery(LicenseRender["data_source"].(string), queryMap)
				if err != nil {
					fmt.Printf("Error querying SODA API: %s", err)
				}
				resultHTML, err := formatOfficerHTML(mergedResults)
				badgeRenderLocal["entityHTML"] = resultHTML
			}
			fmt.Println(queryParams)
		*/
		c.HTML(http.StatusOK, "index.html", badgeRenderLocal)
	}
}

func renderName() gin.HandlerFunc {
	return func(c *gin.Context) {
		//queryParams := c.Request.URL.Query()
		/*
			if len(queryParams["license"]) > 0 {
				queryResults, err := sodaQuery(LicenseRender["data_source"].(string), "last_name", strings.Join(queryParams["name"], ","))
				if err != nil {
					fmt.Println("Error")
				} else {
					if len(queryResults) > 1 {

					}
				}
				var rowStrings []string
				for _, row := range queryResults {
					rowStrings = append(rowStrings, strings.Join(row, ","))
				}
				fmt.Println(strings.Join(rowStrings, ","))
			}
		*/
		c.HTML(http.StatusOK, "index.html", NameRender)
	}
}

//Main Function
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/public", "./public")
	router.GET("/", redirectToLicense())
	router.GET("/license", renderLicense())
	router.GET("/badge", renderBadge())
	router.GET("/name", renderName())
	router.Run(":5000")
}
