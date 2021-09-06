package main

import (
	"image/color"
	"net/http"
	"os"

	"github.com/wsqun/captcha"
)


func main() {

	cap := captcha.New()

	dir, _ := os.Getwd()
	if err := cap.SetFont(dir + "/examples/comic.ttf"); err != nil {
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
	cap.SetCharNum(4)

	svc := captcha.InitService(cap)

	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		pic,err := svc.Draw()
		if err == nil {
			w.Write(pic.GetPNG())
		}
	})

	http.ListenAndServe(":8085", nil)

}
