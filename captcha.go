package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image/png"
	"io"

	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

type Captcha struct {
	frontColors []color.Color
	bkgColors   []color.Color
	disturlvl   DisturLevel
	charNum     int
	fonts       []*truetype.Font
	size        image.Point

	rand *rand.Rand
}

var _ Graffiti = &Captcha{}

type DisturLevel int

const (
	NORMAL DisturLevel = 4
	MEDIUM DisturLevel = 8
	HIGH   DisturLevel = 16
)

func New() Graffiti {
	c := &Captcha{
		disturlvl: NORMAL,
		size:      image.Point{X:82, Y:32},
		charNum:   4,

		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	c.frontColors = []color.Color{color.Black}
	c.bkgColors = []color.Color{color.White}
	return c
}


// AddFont 添加一个字体
func (c *Captcha) AddFont(path string) (err error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontData)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

//AddFontFromBytes allows to load font from slice of bytes, for example, load the font packed by https://github.com/jteeuwen/go-bindata
func (c *Captcha) AddFontFromBytes(contents []byte) error {
	font, err := freetype.ParseFont(contents)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

// SetFont 设置字体 可以设置多个
func (c *Captcha) SetFont(paths ...string) error {
	for _, v := range paths {
		if erro := c.AddFont(v); erro != nil {
			return erro
		}
	}
	return nil
}

func (c *Captcha) SetDisturbance(d DisturLevel) {
	if d > 0 {
		c.disturlvl = d
	}
}

func (c *Captcha) SetCharNum(charNum int) {
	if charNum > 0 {
		c.charNum = charNum
	}
}

func (c *Captcha) SetFrontColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.frontColors = c.frontColors[:0]
		for _, v := range colors {
			c.frontColors = append(c.frontColors, v)
		}
	}
}

func (c *Captcha) SetBkgColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.bkgColors = c.bkgColors[:0]
		for _, v := range colors {
			c.bkgColors = append(c.bkgColors, v)
		}
	}
}

func (c *Captcha) SetSize(w, h int) {
	if w < 48 {
		w = 48
	}
	if h < 20 {
		h = 20
	}
	c.size = image.Point{w, h}
}

func (c *Captcha) randFont() *truetype.Font {
	return c.fonts[rand.Intn(len(c.fonts))]
}

// 绘制背景
func (c *Captcha) drawBkg(img *Image) {
	//填充主背景色
	bgcolorindex := c.rand.Intn(len(c.bkgColors))
	bkg := image.NewUniform(c.bkgColors[bgcolorindex])
	img.FillBkg(bkg)
}

// 绘制噪点
func (c *Captcha) drawNoises(img *Image) {
	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	dlen := int(c.disturlvl)
	// 绘制干扰斑点
	for i := 0; i < dlen; i++ {
		x := c.rand.Intn(size.X)
		y := c.rand.Intn(size.Y)
		r := c.rand.Intn(size.Y/20) + 1
		colorIndex := c.rand.Intn(len(c.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, c.frontColors[colorIndex])
	}

	// 绘制干扰线
	for i := 0; i < dlen; i++ {
		x := c.rand.Intn(size.X)
		y := c.rand.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := c.rand.Intn(size.Y) * o
		h := c.rand.Intn(size.Y/10) * o
		colorIndex := c.rand.Intn(len(c.frontColors))
		img.DrawLine(x, y, x+w, y+h, c.frontColors[colorIndex])
		colorIndex++
	}

}

// 绘制文字
func (c *Captcha) drawString(img *Image, str string) (err error) {

	if c.fonts == nil {
		panic("没有设置任何字体")
	}
	tmp := NewImage(c.size.X, c.size.Y)

	// 文字大小为图片高度的 0.6
	fsize := int(float64(c.size.Y) * 0.6)
	// 用于生成随机角度

	// 文字之间的距离
	// 左右各留文字的1/4大小为内部边距
	padding := fsize / 4
	gap := (c.size.X - padding*2) / (len(str))

	// 逐个绘制文字到图片上
	for i, char := range str {
		// 创建单个文字图片
		// 以文字为尺寸创建正方形的图形
		str := NewImage(fsize, fsize)
		// str.FillBkg(image.NewUniform(color.Black))
		// 随机取一个前景色
		colorindex := c.rand.Intn(len(c.frontColors))

		//随机取一个字体
		font := c.randFont()
		err = str.DrawString(font, c.frontColors[colorindex], string(char), float64(fsize))
		if err != nil {
			return err
		}

		// 转换角度后的文字图形
		rs := str.Rotate(float64(c.rand.Intn(40) - 20))
		// 计算文字位置
		s := rs.Bounds().Size()
		left := i*gap + padding
		top := (c.size.Y - s.Y) / 2
		// 绘制到图片上
		draw.Draw(tmp, image.Rect(left, top, left+s.X, top+s.Y), rs, image.ZP, draw.Over)
	}
	if c.size.Y >= 48 {
		// 高度大于48添加波纹 小于48波纹影响用户识别
		tmp.distortTo(float64(fsize)/10, 200.0)
	}

	draw.Draw(img, tmp.Bounds(), tmp, image.ZP, draw.Over)
	return nil
}

// 绘图
func (c *Captcha) Draw(w io.Writer) ( chars string, err error) {
	dst, chars, err := c.Create(c.charNum)
	if err != nil {
		return
	}

	if err = png.Encode(w, dst); err != nil {
		return
	}
	return chars, nil
}

// Create 生成一个验证码图片
func (c *Captcha) Create(num int) (dst *Image, str string, err error) {
	if num <= 0 {
		num = 4
	}
	dst = NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(dst)

	str = c.randStr(num)
	err = c.drawString(dst, str)

	return
}

func (c *Captcha) CreateCustom(str string) (dst *Image, err error) {
	if len(str) == 0 {
		str = "unkown"
	}
	dst = NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(dst)
	err = c.drawString(dst, str)
	return
}

var letters = []string{"3", "4", "5", "7", "8", "a", "c", "d", "e", "f", "g", "h", "j", "k", "m", "n", "p", "q", "s", "t", "w", "x", "y", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "M", "N", "P", "Q", "R", "S", "V", "W", "X", "Y"}

// 生成随机字符串
// size 个数 kind 模式
func (c *Captcha) randStr(size int) (res string) {

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		res = res + letters[rand.Intn(len(letters))]
	}
	return res
}
