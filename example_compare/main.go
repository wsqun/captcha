package main

import (
	"image/color"
	"image/png"
	"net/http"

	"github.com/wsqun/captcha"
)

var cap *captcha.Captcha

func main() {

	cap = captcha.New().(*captcha.Captcha)

	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	/*
	   //We can load font not only from localfile, but also from any []byte slice
	   	fontContenrs, err := ioutil.ReadFile("comic.ttf")
	   	if err != nil {
	   		panic(err.Error())
	   	}
	   	err = cap.AddFontFromBytes(fontContenrs)
	   	if err != nil {
	   		panic(err.Error())
	   	}
	*/

	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img, _,_ := cap.Create(4)
		png.Encode(w, img)
	})

	http.ListenAndServe(":8085", nil)

}