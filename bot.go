package main

import (
	"os"
	"log"
	//"fmt"
	"time"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Info struct {
	Now int `json:"now"`
	Now_dt string `json:"now_dt"`
	Np np `json:"info"`
	Fact fact `json:"fact"`
	Forecast []forecast `json:"forecasts"` 
}

type np struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
	Tzifo tz `json:"tzifo"`
	Def_pr_mm int `json:"def_pressure_mm"`
	Def_pr_pa int `json:"def_pressure_pa"`
	Yurl string `json:"url"`
}

type tz struct {
	Offset int `json:"offset"`
	Name string `json:"name"`
	Abbr string `json:"abbr"`
	Dst bool `json:"dst"`
}

type fact struct {
	Temp int `json:"temp"`
	Feels int `json:"feels_like"`
	Icon string `json:"icon"`
	Condition string `json:"condition"`
	WindSpeed int `json:"wind_speed"`
	WindGust float32 `json:"wind_gust"`
	WindDir string `json:"wind_dir"`
	Pr_mm int `json:"pressure_mm"`
	Pr_pa int `json:"pressure_pa"`
	Humidity int `json:"humidity"`
	Daytime string `json:"daytime"`
	Polar bool `json:"polar"`
	Season string `json:"season"`
	Prec_type int `json:"prec_type"`
	Prec_strength float32 `json:"prec_strength"`
	Cloudness int `json:"cloudness"`
	Obs_time int `json:"obs_time"`
}

type forecast struct {
	Date string `json:"date"`
	Week int `json:"week"`
	Sunrise string `json:"sunrize"`
	Sunset string `json:"sunset"`
	Parts part `json:"parts"`
	Hours []hour `json:"hours"`
}

type part struct {
	Night day_time `json:"night"`
	Morning day_time `json:"morning"`
	Day day_time `json:"day"`
	Evening day_time `json:"evening"`
	Day_short short_part `json:"day_short"`
	Night_short short_part `json:"night_short"`
}

type day_time struct {
	Temp_min int `json:"temp_min"`
	Temp_max int `json:"temp_max"`
	Temp_avg int `json:"temp_avg"`
	Feels int `json:"feels_like"`
	Windspeed int `json:"wind_speed"`
	Windgust int `json:"wind_gust"`
	Cloudness float64 `json:"cloudness"`
}

type short_part struct {
	Temp int `json:"temp"`
	Feels int `json:"feels_like"`
	Icon string `json:"icon"`
	Condition string `json:"condition"`
	WindSpeed int `json:"wind_speed"`
	WindGust float32 `json:"wind_gust"`
	WindDir string `json:"wind_dir"`
	Pr_mm int `json:"pressure_mm"`
	Pr_pa int `json:"pressure_pa"`
	Humidity int `json:"humidity"`
	Cloudness float64 `json:"cloudness"`
}

type hour struct {
	Hour int `json:"hour"`
	Temp int `json:"temp"`
	Feels int `json:"feels_like"`
	Icon string `json:"icon"`
	Condition string `json:"condition"`
	WindSpeed int `json:"wind_speed"`
	WindGust float32 `json:"wind_gust"`
	WindDir string `json:"wind_dir"`
	Pr_mm int `json:"pressure_mm"`
	Pr_pa int `json:"pressure_pa"`
	Humidity int `json:"humidity"`
}

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
    resp.Write([]byte("Hi there! I'm BoGoBot!"))
}

func main() {
	bot_token := os.Getenv("BOT_TOKEN")
	b, err := tb.NewBot(tb.Settings{
		Token: bot_token, 
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return
	}

	weekday := int(time.Now().Weekday())
	weekdays := [7]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресение"}

	condition_emo := map[string]string{
		"clear": "☀️",
		"partly-cloudy": "🌤️",
		"cloudy": "⛅",
		"overcast": "🌦️",
		"partly-cloudy-and-light-rain": "🌦️",
		"partly-cloudy-and-rain": "🌦️",
		"overcast-and-rain": "🌦️",
		"overcast-thunderstorms-with-rain": "⛈️",
		"cloudy-and-light-rain": "🌦️",
		"overcast-and-light-rain": "🌦️",
		"cloudy-and-rain": "🌦️",
		"overcast-and-wet-snow": "💧❄️",
		"partly-cloudy-and-light-snow": "🌨️",
		"partly-cloudy-and-snow": "🌨️", 
		"overcast-and-snow": "🌨️",
		"cloudy-and-light-snow": "🌨️",
		"overcast-and-light-snow": "🌨️",
		"cloudy-and-snow": "🌨️",
	}

	log.Printf("Authorized on account PoGoBot")

	var info Info

	act_data := tb.InlineButton{
		Unique: "AD",
		Text:   "🗂️ Получить актуальные данные",
	}
	
	fact := tb.InlineButton{
		Unique: "F",
		Text:   "☁️ Сейчас",
	}

	nextd := tb.InlineButton{
		Unique: "ND",
		Text:   "☁️ Завтра",
	}

	fact_per_hour := tb.InlineButton{
		Unique: "FPH",
		Text:   "🕒 Сегодня по часам",
	}

	tomorrow_per_hour := tb.InlineButton{
		Unique: "TPH",
		Text:   "🕒 Завтра по часам",
	}

	det := tb.InlineButton{
		Unique: "DT",
		Text:   "📉 Детальный прогноз",
	}

	per_week := tb.InlineButton{
		Unique: "PW",
		Text:   "📆 Прогноз на неделю",
	}

	per_weekend := tb.InlineButton{
		Unique: "PEND",
		Text:   "📆 Прогноз на выходные",
	}

	BackToMain := tb.InlineButton{
		Unique: "BM",
		Text:   "⬅️ Назад",
	}
	

	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{act_data},
		[]tb.InlineButton{fact, nextd},
		[]tb.InlineButton{fact_per_hour, tomorrow_per_hour},
		[]tb.InlineButton{det},
	}

	detInline := [][]tb.InlineButton{
		[]tb.InlineButton{per_week},
		[]tb.InlineButton{per_weekend},
		[]tb.InlineButton{BackToMain},
	}

	http.HandleFunc("/", MainHandler)
    go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	b.Handle("/start", func(m *tb.Message) {
			b.Send(m.Sender, "🌤️ Привет!\nЯ помогу узнать о погоде", &tb.ReplyMarkup{
				InlineKeyboard: mainInline,
			})

			b.Handle(&act_data, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": act_data")

				client := &http.Client{
				}
				
				req, err := http.NewRequest("GET", "https://api.weather.yandex.ru/v1/forecast?lat=55.715723&lon=37.459478", nil)
				ykey := "X-Yandex-API-Key"
				yval := os.Getenv("YA_TOKEN")
				req.Header.Add(ykey, yval)
				if err != nil {
					log.Fatalln(err)
				}
				resp, err := client.Do(req)
			
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalln(err)
				}

				json.Unmarshal(body, &info)

				defer resp.Body.Close()

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&nextd, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": nextd")
				cld := info.Forecast[1].Parts.Day_short.Condition
				temp := "\nЗавтра температура воздуха составит: " + strconv.Itoa(info.Forecast[1].Parts.Day_short.Temp) + " ℃"
				feels := "\nОщущается как: " + strconv.Itoa(info.Forecast[1].Parts.Day_short.Feels) + " ℃"
				wind := "\nСкорость ветра: " + strconv.Itoa(info.Forecast[1].Parts.Day_short.WindSpeed) + "м/с"

				ureq := condition_emo[cld] + condition_emo[cld] + condition_emo[cld] + temp + feels + wind

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&fact, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": fact ---")
				cld := info.Fact.Condition
				temp := "\nТемпература воздуха составляет: " + strconv.Itoa(info.Fact.Temp) + " ℃"
				feels := "\nОщущается как: " + strconv.Itoa(info.Fact.Feels) + " ℃"
				wind := "\nСкорость ветра: " + strconv.Itoa(info.Fact.WindSpeed)+ " м/с"

				ureq := condition_emo[cld] + condition_emo[cld] + condition_emo[cld] + temp + feels + wind 

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&fact_per_hour, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_week")

				ureq := "Сегодня\n"

				for i := 0; i < len(info.Forecast[0].Hours); i++ {
					pw_date := "Час: " + strconv.Itoa(i)
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(info.Forecast[0].Hours[i].Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(info.Forecast[0].Hours[i].Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(info.Forecast[0].Hours[i].WindSpeed)+ " м/с\n\n"
					cld := info.Forecast[0].Hours[i].Condition
					count := weekday+i-1
					if count > 6 {
						count = count - 7
					}
					ureq += "🕒 " + pw_date + " " + condition_emo[cld] + temp + feels + wind
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&tomorrow_per_hour, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_week")

				ureq := "Завтра\n"

				for i := 0; i < len(info.Forecast[1].Hours); i++ {
					pw_date := "Час: " + strconv.Itoa(i)
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(info.Forecast[1].Hours[i].Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(info.Forecast[1].Hours[i].Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(info.Forecast[1].Hours[i].WindSpeed)+ " м/с\n\n"
					cld := info.Forecast[1].Hours[i].Condition
					count := weekday+i-1
					if count > 6 {
						count = count - 7
					}
					ureq += "🕒 " + pw_date + " " + condition_emo[cld] + temp + feels + wind
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&per_week, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_week")

				var ureq string

				for i := 0; i < len(info.Forecast); i++ {
					pw_date := "Дата: " + info.Forecast[i].Date
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(info.Forecast[i].Parts.Day_short.Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(info.Forecast[i].Parts.Day_short.Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(info.Forecast[i].Parts.Day_short.WindSpeed)+ " м/с\n\n"
					cld := info.Forecast[i].Parts.Day_short.Condition
					count := weekday+i-1
					if count > 6 {
						count = count - 7
					}
					
					ureq += condition_emo[cld] + " " + weekdays[count] + "\n" + pw_date + temp + feels + wind 
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: detInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&per_weekend, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_weekend")

				var ureq string

				for i := 0; i < len(info.Forecast); i++ {
					pw_date := "Дата: " + info.Forecast[i].Date
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(info.Forecast[i].Parts.Day_short.Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(info.Forecast[i].Parts.Day_short.Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(info.Forecast[i].Parts.Day_short.WindSpeed)+ " м/с\n\n"
					cld := info.Forecast[i].Parts.Day_short.Condition
					count := weekday+i-1
					if count > 6 {
						count = i-weekday-2
					}
					if count > 4 {
						ureq += condition_emo[cld] + " " + weekdays[count] + "\n" + pw_date + temp + feels + wind 
					}
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: detInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&det, func(c *tb.Callback) {
				b.Edit(c.Message, "Детальный прогноз", &tb.ReplyMarkup{
					InlineKeyboard: detInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&BackToMain, func(c *tb.Callback) {
				b.Edit(c.Message, "🌤️ Привет!\nЯ помогу узнать о погоде", &tb.ReplyMarkup{
					InlineKeyboard: mainInline})
				b.Respond(c, &tb.CallbackResponse{})
			})
	})

	b.Start()
}