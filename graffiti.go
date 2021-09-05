package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image/color"
	"io/ioutil"
)

type Graffiti interface {
	// 添加字体
	AddFonts(fonts []*truetype.Font)
	// 设置平面长宽
	SetSize(width int, height int)
	// 设置前景色
	SetFrontColor(colors []color.Color)
	// 设置背景色
	SetBkgColor(colors []color.Color)
	// 设置干扰等级
	SetDisturbance(level int)
	// 设置字符长度
	SetCharNum(num int)
	// 绘图
	Draw() (img []byte, chars string, err error)
}


// AddFont 添加一个字体
func GetFont(path string) (font *truetype.Font, err error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	font, err = freetype.ParseFont(fontData)
	if err != nil {
		return
	}
	return
}