package controller

import (
	"fmt"
	"os"

	"encoding/json"

	"strings"

	"time"

	"github.com/golang/glog"
	"github.com/knative-sample/weather-store/pkg/tablestore"
	"github.com/knative-sample/weather-store/pkg/weather"
)

// 存储北京海淀和郑州市区域天气
func StoreWeather() {
	cityCodes := strings.Split(city, "\n")
	client := tablestore.InitClient()
	for i, cityCode := range cityCodes {
		// 控制并发操作
		if i % 50 == 0 {
			time.Sleep(2 * time.Second)
		}
		go func(cityCode string) {
			glog.Infof("start query city: %s", cityCode)
			key := os.Getenv("WEATHER_API_KEY")
			queryApi := fmt.Sprintf("%s?key=%s&city=%s&extensions=all", weather.WEATHER_API, key, cityCode)
			res, err := weather.QueryWeather(queryApi, "")
			if err != nil {
				glog.Errorf("QueryWeather error: %s", err.Error())
				return
			}
			glog.Infof("weather info: %s", res)
			wr := weather.WeatherResponse{}
			err = json.Unmarshal(res, &wr)
			if err != nil {
				glog.Errorf("QueryWeather Unmarshal error: %s", err.Error())
				return
			}
			if wr.Status == "1" && len(wr.Forecasts) > 0 {
				glog.Infof("start store city: %s", cityCode)
				err := client.Store(wr.Forecasts[0])
				if err != nil {
					glog.Errorf("Weather Store error: %s", err.Error())
					return
				}
				glog.Infof("store city successfully: %s", cityCode)
			}
		}(cityCode)

	}
}

const city = `100190
450000`
