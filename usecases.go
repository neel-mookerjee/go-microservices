package main

import (
	"errors"
	"time"
)

type UsecaseInterface interface {
	GetNextPreview() (DownSamplingPreviewObject, error)
	GetNextQuery() (DownsamplingObject, error)
	DeployPreviewItem(id string) error
	GetDownsamplingItem(id string) (DownsamplingObject, error)
	GetMetricsStacks() (MetricsStacks, error)
}

type Usecase struct {
	db      DbInterface
	Configs *Config
}

func NewUsecase(config *Config) (*Usecase, error) {
	db, err := NewDynamodb(config)
	if err != nil {
		return nil, err
	}
	return &Usecase{Configs: config, db: db}, nil
}

func (u *Usecase) GetNextPreview() (DownSamplingPreviewObject, error) {
	var ds DownsamplingObject
	var q DownSamplingPreviewObject
	dsList, err := u.db.GetPendingDownsamplePreviewItems()
	if err != nil {
		return q, err
	}

	DownsampleObjects(dsList).SortByTime()

	if len(dsList) > 0 {
		ds = dsList[0]
	}

	if ds.QueryId != "" {
		q = ds.CreatePreviewObject()
	}

	return q, err
}

func (u *Usecase) GetNextQuery() (DownsamplingObject, error) {
	var ds DownsamplingObject
	dsList, err := u.db.GetPendingDownsamplePreviewItems()
	if err != nil {
		return ds, err
	}

	DownsampleObjects(dsList).SortByTime()

	if len(dsList) > 0 {
		ds = dsList[0]
	}

	return ds, err
}

func (u *Usecase) DeployPreviewItem(id string) error {
	var ds DownsamplingObject
	ds, err := u.db.GetDownsamplingItem(id)
	if err != nil {
		return err
	}

	if ds.QueryId == "" {
		return errors.New("object not found")
	} else if ds.QueryState != "PREVIEW_PENDING" {
		return errors.New("status changed")
	}

	ds.QueryState = "PREVIEW_DEPLOYED"
	ds.PreviewExpiresAt = time.Now().Add(45 * time.Minute).Format(time.RFC3339)

	_, err = u.db.UpdateDownsamplePreviewItem(ds)
	return err
}

func (u *Usecase) GetDownsamplingItem(id string) (DownsamplingObject, error) {
	return u.db.GetDownsamplingItem(id)
}

func (u *Usecase) GetMetricsStacks() (MetricsStacks, error) {
	return u.db.GetMetricsStacks()
}
