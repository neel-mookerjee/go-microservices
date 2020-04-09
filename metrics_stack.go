package main

type MetricsStack struct {
	Nickname   string `json:"nickname"`
	DbName     string `json:"db_name"`
	ClusterUrl string `json:"cluster_url"`
}

type MetricsStacks []MetricsStack
