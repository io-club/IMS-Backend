package model

import "gorm.io/gorm"

type Weather struct {
	gorm.Model
	DateTime string `gorm:"type:char(10);comment:日期"`
	Location string `gorm:"type:char(20);comment:地点"` // 如：天津 - 西青，天津 - 城区，贵州 - 安顺 - 城区（没有空格，插件导致）

	Hour          string `gorm:"type:char(2);comment:时间"`
	Temperature   int8   `gorm:"type:float;comment:温度"`
	WindDirection string `gorm:"type:char(3);comment:风向"`
	WindForce     int8   `gorm:"type:int;comment:风级"`
	Precipitation int8   `gorm:"type:int;comment:降水量"`
	Humidity      int8   `gorm:"type:int;comment:湿度"`
}

func (Weather) TableName() string {
	return "ims_work_weather"
}
