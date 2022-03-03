package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
)

type Schema struct {
	Name   string            `json:"name"`
	Schema map[string]string `json:"schema"`
}

const MAX_UPLOAD_SIZE = 10 * 1024 * 1024

func GetSchema(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		JSONMarshalBadRequest("File size should not be greater then 10 MB", 400).WriteToResponse(w)
		return
	}

	var value = make(map[string]string)

	file, handler, _ := r.FormFile("file")

	if file != nil {

		fileName := handler.Filename

		extensionFlag := getFileExtension(fileName)

		if extensionFlag {
			JSONMarshalBadRequest("Please provide CSV file", 500).WriteToResponse(w)
			return
		}

		csvLines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			JSONMarshalErr(err.Error(), 500).WriteToResponse(w)
			return
		}

		if len(csvLines) == 0 || len(csvLines) == 1 {
			JSONMarshalBadRequest("File does not contain any data", 400).WriteToResponse(w)
			return
		}

		headerRow := csvLines[0]

		dataRows := csvLines[1]

		for index := range headerRow {
			result := getDataType(dataRows[index])
			value[headerRow[index]] = result
		}

		csvSchema := &Schema{
			Name:   fileName,
			Schema: value,
		}

		json.NewEncoder(w).Encode(csvSchema)
		return
	}

	JSONMarshalNoContent("Please upload the file", 204).WriteToResponse(w)
}

func getDataType(value string) string {

	var result string

	if _, err := strconv.ParseInt(value, 10, 32); err == nil {
		result = "int"
	} else if _, err := strconv.ParseFloat(value, 32); err == nil {
		result = "float"
	} else if _, err := strconv.ParseFloat(value, 64); err == nil {
		result = "float"
	} else if _, err := dateparse.ParseLocal(value); err == nil {
		result = "datetime"
	} else {
		result = "string"
	}
	return result
}

func getFileExtension(fileName string) bool {

	index := strings.Index(fileName, ".")

	if index > 0 {
		extension := fileName[index+1:]

		if extension == "csv" || extension == "CSV" {
			return false
		}
	}

	return true
}
