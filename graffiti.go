package captcha

import "github.com/golang/freetype/truetype"

type Graffiti interface {
	// 添加字体
	AddFont(fonts []*truetype.Font)
	// 设置平面长宽
	SetSize(long int, wide int)
	// 设置字体颜色
	SetFontColor()
	// 设置背景色
	SetBackgroundColor()
	// 设置干扰等级
	SetNoiseLevel()
	// 绘图
	CreateImg()
}
