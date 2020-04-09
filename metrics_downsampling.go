package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DownsamplingObject struct {
	Nickname          string `json:"nickname"`
	QueryId           string `json:"queryId"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	Db                string `json:"db"`
	Rp                string `json:"rp"`
	Measurement       string `json:"measurement"`
	TargetRp          string `json:"targetRp"`
	TargetMeasurement string `json:"targetMeasurement"`
	PreviewExpiresAt  string `json:"previewExpiresAt"`
	QueryState        string `json:"queryState"`
	Fields            []struct {
		Alias    string `json:"alias"`
		Field    string `json:"field"`
		Function string `json:"func"`
	} `json:"fields"`
	Tags     []string `json:"tags"`
	Interval int      `json:"interval"`
}

type DownsampleObjects []DownsamplingObject

func (slice DownsampleObjects) Len() int {
	return len(slice)
}

func (slice DownsampleObjects) Less(i, j int) bool {
	return slice[i].CreatedAt < slice[j].CreatedAt
}

func (slice DownsampleObjects) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice DownsampleObjects) SortByTime() {
	sort.Sort(slice)
}

type DownSamplingPreviewObject struct {
	QueryId           string `json:"queryId"`
	Query             string `json:"query"`
	Database          string `json:"database"`
	JobId             string `json:"job"`
	BeginSample       int64  `json:"beginSample"`
	EndSample         int64  `json:"endSample"`
	StartedAt         int64  `json:"startedAt"`
	CompletedSample   int64  `json:"completedSample"`
	CreatedAt         int64  `json:"createdAt"`
	UpdatedAt         int64  `json:"updatedAt"`
	Note              string `json:"note"`
	State             string `json:"state"`
	TargetMeasurement string `json:"targetMeasurement"`
	Fields            []struct {
		Alias    string `json:"alias"`
		Field    string `json:"field"`
		Function string `json:"func"`
	} `json:"fields"`
	Tags     []string `json:"tags"`
	Interval int      `json:"interval"`
	ready    bool
}

func (ds DownsamplingObject) CreatePreviewObject() DownSamplingPreviewObject {
	var query DownSamplingPreviewObject
	query.QueryId = ds.QueryId
	query.Query = ds.formQuery()
	query.Database = ds.Db
	beginSample := time.Now().Unix() - 3*24*60*60
	// begin sample set to previous hour
	beginSample = beginSample - (beginSample % 3600)
	query.BeginSample = beginSample
	endSample := time.Now().Unix()
	// end sample set to previous hour
	endSample = endSample - (endSample % 3600)
	query.EndSample = endSample
	query.CompletedSample = 0
	query.CreatedAt = time.Now().Unix()
	query.UpdatedAt = time.Now().Unix()
	query.Note = "Continuous query created"
	query.JobId = "previewjob"
	query.StartedAt = 0
	query.State = "Ready"
	query.TargetMeasurement = ds.TargetMeasurement
	query.Fields = ds.Fields
	query.Tags = ds.Tags
	query.Interval = ds.Interval

	return query
}

func (ds DownsamplingObject) formQuery() string {
	template := "SELECT %s FROM %s WHERE time >= $ds_start_ts AND time <= $ds_end_ts GROUP BY time(%ss),%s fill(none)"
	source := fmt.Sprintf("\"%s\".\"%s\".\"%s\"", ds.Db, "autogen", ds.Measurement)

	var tagsArr []string
	for _, t := range ds.Tags {
		tagsArr = append(tagsArr, fmt.Sprintf("\"%s\"", t))
	}
	tags := strings.Join(tagsArr, ",")

	var fieldsArr []string
	for _, f := range ds.Fields {
		str := ""
		if strings.Contains(f.Function, ":") {
			arr := strings.Split(f.Function, ":")
			str = fmt.Sprintf("%s(\"%s\",%s) AS \"p%s_%s\"", arr[0], f.Field, arr[1], arr[1], f.Field)
		} else {
			str = fmt.Sprintf("%s(\"%s\") AS \"%s_%s\"", f.Function, f.Field, strings.ToLower(f.Function), f.Field)
		}
		fieldsArr = append(fieldsArr, str)
	}
	fields := strings.Join(fieldsArr, ",")

	queryStr := fmt.Sprintf(template, fields, source, strconv.Itoa(ds.Interval), tags)

	return queryStr
}
