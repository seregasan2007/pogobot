package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Struct of cities.json

type CitiesInfo struct {
	Cities []cities_info `json:"cities"`
}

type cities_info struct {
	Name string `json:"Город"`
	Lat float64 `json:"Широта"`
	Lon float64 `json:"Долгота"`
}

// Struct of json response from Yandex.Weather API

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

// Function for keep session on Heroku 

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
    resp.Write([]byte("Hi there! I'm BoGoBot!"))
}

// Function for GET request on Yandex.Weather

func GetWeather(ulat, ulon float64) Info {
	var info Info
	lat := fmt.Sprintf("%f", ulat)
	lon := fmt.Sprintf("%f", ulon)

	// Add header on GET request

	client := &http.Client{
	}
	
	req, err := http.NewRequest("GET", "https://api.weather.yandex.ru/v1/forecast?lat="+lat+"&lon="+lon, nil)
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
	return info
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

	log.Printf("Authorized on account PoGoBot")

	var cities CitiesInfo
	user_info := make(map[string]Info)
	user_city := make(map[string]string)
	user_lat := make(map[string]float64)
	user_lon := make(map[string]float64)

	weekday := int(time.Now().Weekday())
	weekdays := [7]string{"Воскресение", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}

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

	condition_desc := map[string]string{
		"clear": "Ясно",
		"partly-cloudy": "Малооблачно",
		"cloudy": "Облачно с прояснениями",
		"overcast": "Пасмурно",
		"partly-cloudy-and-light-rain": "Малооблачно, небольшой дождь",
		"partly-cloudy-and-rain": "Малооблачно, дождь",
		"overcast-and-rain": "Значительная облачность, сильный дождь",
		"overcast-thunderstorms-with-rain": "Сильный дождь с грозой",
		"cloudy-and-light-rain": "Облачно, небольшой дождь",
		"overcast-and-light-rain": "Значительная облачность, небольшой дождь",
		"cloudy-and-rain": "Облачно, дождь",
		"overcast-and-wet-snow": "Дождь со снегом",
		"partly-cloudy-and-light-snow": "Небольшой снег",
		"partly-cloudy-and-snow": "Малооблачно, снег", 
		"overcast-and-snow": "Снегопад",
		"cloudy-and-light-snow": "Облачно, небольшой снег",
		"overcast-and-light-snow": "Значительная облачность, небольшой снег",
		"cloudy-and-snow": "Облачно, снег",
	}
	
	// Open our cities.json 
	jsonCities, err := ioutil.ReadFile("cities.json")
	if err != nil {
		log.Println(err)
	}

	// The BOM identifies that the text is UTF-8 encoded, but it should be removed before decoding.
	jsonCities = bytes.TrimPrefix(jsonCities, []byte("\xef\xbb\xbf"))
	  
	err = json.Unmarshal(jsonCities, &cities)
    if err != nil {
        log.Println("error:", err)
    }

	// Set buttons on bot

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
	
	// Collect buttons on group

	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{fact, nextd},
		[]tb.InlineButton{fact_per_hour, tomorrow_per_hour},
		[]tb.InlineButton{det},
		[]tb.InlineButton{act_data},
	}

	detInline := [][]tb.InlineButton{
		[]tb.InlineButton{per_week},
		[]tb.InlineButton{per_weekend},
		[]tb.InlineButton{BackToMain},
	}

	http.HandleFunc("/", MainHandler)
    go http.ListenAndServe(":"+os.Getenv("PORT"), nil)     

	b.Handle("/start", func(m *tb.Message) {
			b.Send(m.Sender, "🌤️ Привет!\nЯ помогу узнать о погоде. Но сперва нужно выбрать город. Для выбора города отправь мне его название.")

			b.Handle(tb.OnText, func(m *tb.Message) {
				if_found := 0

				log.Println(m.Sender.Username, ":", m.Text)

				for  i := 0;  i < len(cities.Cities);  i++ {
					if m.Text == cities.Cities[i].Name {
						user_city[m.Sender.Username] = cities.Cities[i].Name
						user_lat[m.Sender.Username] = cities.Cities[i].Lat
						user_lon[m.Sender.Username] = cities.Cities[i].Lon
						if_found = 1
					}
				}

				if if_found == 0 {
					uresp := "Город не найден.\nНазвание города должно быть написано полностью и с большой буквы. Пример сообщения:\nМосква"
					b.Send(m.Sender, uresp)
				} else {
					user_info[m.Sender.Username] = GetWeather(user_lat[m.Sender.Username], user_lon[m.Sender.Username])
					uresp := "Выбран город: " + user_city[m.Sender.Username]
					b.Send(m.Sender, uresp, &tb.ReplyMarkup{
						InlineKeyboard: mainInline,
					})
				}
			})

			b.Handle(&act_data, func(c *tb.Callback) {

				log.Println(c.Sender.Username, ": act_data")

				user_info[m.Sender.Username] = GetWeather(user_lat[m.Sender.Username], user_lon[m.Sender.Username])

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&nextd, func(c *tb.Callback) {

				log.Println(c.Sender.Username, ": nextd")

				cld := user_info[c.Sender.Username].Forecast[1].Parts.Day_short.Condition
				city_header := "\nЗавтра в городе " + user_city[c.Sender.Username] + "\n"
				temp := "\nТемпература воздуха составит: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[1].Parts.Day_short.Temp) + " ℃"
				feels := "\nОщущается как: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[1].Parts.Day_short.Feels) + " ℃"
				wind := "\nСкорость ветра: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[1].Parts.Day_short.WindSpeed) + "м/с"

				ureq := condition_emo[cld] + condition_emo[cld] + condition_emo[cld] + city_header + condition_desc[cld] +  temp + feels + wind

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&fact, func(c *tb.Callback) {

				log.Println(c.Sender.Username, ": fact ---")

				cld := user_info[c.Sender.Username].Fact.Condition
				city_header := "\nПогода в городе " + user_city[c.Sender.Username] + "\n"
				temp := "\nТемпература воздуха составляет: " + strconv.Itoa(user_info[c.Sender.Username].Fact.Temp) + " ℃"
				feels := "\nОщущается как: " + strconv.Itoa(user_info[c.Sender.Username].Fact.Feels) + " ℃"
				wind := "\nСкорость ветра: " + strconv.Itoa(user_info[c.Sender.Username].Fact.WindSpeed)+ " м/с"

				ureq := condition_emo[cld] + condition_emo[cld] + condition_emo[cld] + city_header + condition_desc[cld] + temp + feels + wind 

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&fact_per_hour, func(c *tb.Callback) {

				log.Println(c.Sender.Username, ": per_week")

				ureq := "Сегодня в городе " + user_city[c.Sender.Username] +"\n"

				for i := 0; i < len(user_info[c.Sender.Username].Forecast[0].Hours); i++ {
					pw_date := "Час: " + strconv.Itoa(i)
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[0].Hours[i].Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[0].Hours[i].Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[0].Hours[i].WindSpeed)+ " м/с\n\n"
					cld := user_info[c.Sender.Username].Forecast[0].Hours[i].Condition

					ureq += "🕒 " + pw_date + " " + condition_emo[cld] + "\n" + condition_desc[cld] + temp + feels + wind
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&tomorrow_per_hour, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_week")

				ureq := "Завтра в городе " + user_city[c.Sender.Username] + "\n"

				for i := 0; i < len(user_info[c.Sender.Username].Forecast[1].Hours); i++ {
					pw_date := "Час: " + strconv.Itoa(i)
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[1].Hours[i].Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[1].Hours[i].Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[1].Hours[i].WindSpeed)+ " м/с\n\n"
					cld := user_info[c.Sender.Username].Forecast[1].Hours[i].Condition

					ureq += "🕒 " + pw_date + " " + condition_emo[cld] + "\n" + condition_desc[cld] + temp + feels + wind
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: mainInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&per_week, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_week")

				var ureq string

				for i := 0; i < len(user_info[c.Sender.Username].Forecast); i++ {
					pw_date := "Дата: " + user_info[c.Sender.Username].Forecast[i].Date
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[i].Parts.Day_short.Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[i].Parts.Day_short.Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[i].Parts.Day_short.WindSpeed)+ " м/с\n\n"
					cld := user_info[c.Sender.Username].Forecast[i].Parts.Day_short.Condition
					count := weekday+i
					if count > 6 {
						count = count - 7
					}
					
					ureq += condition_emo[cld] + " " + weekdays[count] + "\n" + pw_date + "\n" + condition_desc[cld] + temp + feels + wind 
				}

				b.Edit(c.Message, ureq, &tb.ReplyMarkup{
					InlineKeyboard: detInline,
				})

				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&per_weekend, func(c *tb.Callback) {

				log.Println(m.Sender.Username, ": per_weekend")

				var ureq string

				for i := 0; i < len(user_info[c.Sender.Username].Forecast); i++ {
					pw_date := "Дата: " + user_info[c.Sender.Username].Forecast[i].Date
					temp := "\nТемпература воздуха составит: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[i].Parts.Day_short.Temp) + " ℃"
					feels := "\nОщущается как: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[i].Parts.Day_short.Feels) + " ℃"
					wind := "\nСкорость ветра: " + strconv.Itoa(user_info[c.Sender.Username].Forecast[i].Parts.Day_short.WindSpeed)+ " м/с\n\n"
					cld := user_info[c.Sender.Username].Forecast[i].Parts.Day_short.Condition
					count := weekday+i
					if count > 6 {
						count += -7
					}
					if count == 6 || count == 0 {
						ureq += condition_emo[cld] + " " + weekdays[count] + "\n" + pw_date + "\n" + condition_desc[cld] + temp + feels + wind 
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
