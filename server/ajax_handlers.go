package server

import (
	//"fmt"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	"github.com/open-geotechnical/ogt-ags-go/ogtags"
)

// SendAjaxPayload writes payload to the response in requested format
func SendAjaxPayload(resp http.ResponseWriter, request *http.Request, payload interface{}) {

	// pretty returns indents data and readable (notably json) is ?pretty=1 in url
	pretty := true //request.URL.Query().Get("pretty") == "0"

	// Determine which encoding from the mux/router
	vars := mux.Vars(request)
	enc := vars["ext"]
	if enc == "" {
		enc = "json"
	}
	// TODO validate encoding and serialiser
	// eg yaml, json/js, html, xlsx, ags4,

	// TODO map[string] = encoding

	// Lets get ready to encode folks...
	var bites []byte
	var err error
	var content_type string = "text/plain"

	if enc == "yaml" {
		bites, err = yaml.Marshal(payload)
		content_type = "text/yaml"

	} else if enc == "json" || enc == "js" {
		if pretty {
			bites, err = json.MarshalIndent(payload, "", "  ")
		} else {
			bites, err = json.Marshal(payload)
		}
		content_type = "application/json"

	} else {
		bites = []byte("OOPs no `.ext` recognised ")
	}

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", content_type)
	resp.Write(bites)
}

type ErrorPayload struct {
	Success bool   ` json:"success" `
	Error   string ` json:"error" `
}

func SendAjaxError(resp http.ResponseWriter, request *http.Request, err error) {
	SendAjaxPayload(resp, request, ErrorPayload{Success: true, Error: err.Error()})
}




type UnitsPayload struct {
	Success bool        ` json:"success" `
	Units   []ogtags.Unit ` json:"units" `
}

// handles /ags/4/units.*
func AX_Units(resp http.ResponseWriter, req *http.Request) {

	payload := new(UnitsPayload)
	payload.Success = true
	payload.Units = ogtags.GetUnits()

	SendAjaxPayload(resp, req, payload)
}



var EndPoints = map[string]string{
	"/":            "Data and Sys information",
	"/ags/4/all":   "AGS4: All data",
	"/ags/4/units": "AGS4: Units",
}

func AX_Info(resp http.ResponseWriter, req *http.Request) {
	payload := map[string]interface{}{
		"repos":      "https://bitbucket.org/daf0dil/ags-def-json",
		"version":    "0.1-alpha",
		"server_utc": time.Now().UTC().Format("2006-01-02 15:04:05"),
		"endpoints":  EndPoints,
	}
	SendAjaxPayload(resp, req, payload)
}

type AbbrsPayload struct {
	Success bool           ` json:"success" `
	Abbrs []*ogtags.AbbrDD ` json:"abbrs" `
}

// handles /ags4/abbreviations
func AX_Abbrs(resp http.ResponseWriter, req *http.Request) {

	var e error
	payload := new(AbbrsPayload)
	payload.Success = true
	payload.Abbrs, e = ogtags.GetAbbrsDD()
	if e != nil {
		SendAjaxError(resp, req, e)
		return
	}
	SendAjaxPayload(resp, req, payload)
}


type AbbrPayload struct {
	Success bool         ` json:"success" `
	Found   bool         ` json:"found" `
	Abbr  *ogtags.AbbrDD   ` json:"abbr" `
}

// handles /ajax/ags4/abbr/<head_code>*
func AX_Abbr(resp http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	abbr, found, err := ogtags.GetAbbrDD(vars["head_code"])
	if err != nil {
		SendAjaxError(resp, req, err)
		return
	}
	payload := new(AbbrPayload)
	payload.Success = true
	payload.Found = found
	payload.Abbr = abbr
	SendAjaxPayload(resp, req, payload)
}

type GroupsPayload struct {
	Success bool          ` json:"success" `
	Groups  []*ogtags.GroupDD ` json:"groups" `
	GroupsCount  int      ` json:"groups_count" `
}

// handles /ags4/groups
func AX_Groups(resp http.ResponseWriter, req *http.Request) {

	var e error
	payload := new(GroupsPayload)
	payload.Success = true
	payload.Groups, e = ogtags.GetGroupsDD()
	if e != nil {
		SendAjaxError(resp, req, e)
		return
	}
	payload.GroupsCount = len(payload.Groups)

	SendAjaxPayload(resp, req, payload)
}


type GroupPayload struct {
	Success bool        ` json:"success" `
	Group   *ogtags.GroupDD ` json:"group" `
}

// handles /ajax/ags4/group
func AX_Group(resp http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	grp, err := ogtags.GetGroupDD(vars["group_code"])
	if err != nil {
		SendAjaxError(resp, req, err)
		return
	}
	payload := new(GroupPayload)
	payload.Success = true
	payload.Group = grp

	SendAjaxPayload(resp, req, payload)
}


type ExamplesPayload struct {
	Success bool        	` json:"success" `
	Examples   []ExampleRow ` json:"examples" `
}
type ExampleRow struct {
	FileName   string ` json:"file_name" `
}

// handles /ags/4/units.*
func AX_Examples(resp http.ResponseWriter, req *http.Request) {

	payload := new(ExamplesPayload)
	payload.Success = true

	recs, err := ogtags.GetExamples()
	if err != nil {
		SendAjaxError(resp, req, err)
		return
	}

	payload.Examples = make([]ExampleRow, 0)
	for _, ex := range recs{
		payload.Examples = append(payload.Examples, ExampleRow{FileName: ex})
	}

	SendAjaxPayload(resp, req, payload)
}

type DocumentPayload struct {
	Success 	bool        ` json:"success" `
	Document  	ogtags.Document ` json:"document" `
}

// handles /ags/4/units.*
func AX_Parse(resp http.ResponseWriter, req *http.Request) {

	example := req.URL.Query().Get("example")

	payload := new(DocumentPayload)
	payload.Success = true

	doc, err := ogtags.ParseExample(example)
	if err != nil {
		SendAjaxError(resp, req, err)
		return
	}
	payload.Document = *doc



	SendAjaxPayload(resp, req, payload)

}
