package repo

import (
	"context"
	"ims-server/microservices/work/internal/dal/model"
	ioginx "ims-server/pkg/ginx"
	"time"
)

var WeatherSelect = []string{
	"id",
	"date_time",
	"location",
	"hour",
	"temperature",
	"wind_direction",
	"wind_force",
	"precipitation",
	"humidity",
}

type weatherRepo struct {
	ioginx.IRepo[model.Weather]
}

func NewWeatherRepo() *weatherRepo {
	return &weatherRepo{}
}

func (w *weatherRepo) GetByCreateBefore(ctx context.Context, createBefore time.Time) ([]model.Weather, error) {
	var weathers []model.Weather
	err := w.DB().WithContext(ctx).Where("created_at < ?", createBefore).Find(&weathers).Error
	if err != nil {
		return nil, err
	}
	return weathers, nil
}
