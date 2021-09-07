package captcha

import (
	"bytes"
	"image/color"
	"io"
	"sync"
	"time"
)

type Graffiti interface {
	// 添加字体
	SetFont(path ...string) error
	// 设置平面长宽
	SetSize(width int, height int)
	// 设置前景色
	SetFrontColor(colors ...color.Color)
	// 设置背景色
	SetBkgColor(colors ...color.Color)
	// 设置干扰等级
	SetDisturbance(level DisturLevel)
	// 设置字符长度
	SetCharNum(num int)
	// 绘图
	Draw(w io.Writer) (chars string, err error)
}

type Service struct {
	Gra        Graffiti
	poolBuffer *sync.Pool

	picOnLine  []*Picture
	picOffLine []*Picture
	picOffset  int
	imgLen     int
}

type Picture struct {
	chars string
	data  bytes.Buffer
}

func (p *Picture) GetChars() string {
	return p.chars
}

// 不要修改内容！
func (p *Picture) GetPNG() []byte {
	return p.data.Bytes()
}

func InitService(gra Graffiti) (svc *Service) {
	svc = &Service{
		Gra:       gra,
		picOffset: 0,
		imgLen:    16,
		poolBuffer: &sync.Pool{New: func() interface{} {
			img := &Picture{
				data: bytes.Buffer{},
			}
			return img
		}},
	}
	svc.fillImg()
	go svc.defaultStrategyIncOffset()
	go svc.resetImgOnline()
	return
}

// 初始化online图片数据
func (s *Service) fillImg() {
	if s.picOnLine == nil {
		var err error
		for i := 0; i < s.imgLen; i++ {
			img := s.poolBuffer.Get().(*Picture)
			img.chars, err = s.Gra.Draw(&img.data)
			if err != nil {
				panic(err)
			}
			s.picOnLine = append(s.picOnLine, img)
		}
	}
}

// 默认偏移策略
func (s *Service) defaultStrategyIncOffset() {
	tk := time.NewTicker(300 * time.Millisecond)
	for {
		<-tk.C
		offset := s.picOffset + 1
		if offset >= len(s.picOnLine) {
			offset = 0
		}
		s.picOffset = offset
	}
}

// 对图片数据进行定时重置
func (s *Service) resetImgOnline() {
	tk := time.NewTicker(1 * time.Second)
	var err error
	for {
		<-tk.C
		img := s.poolBuffer.Get().(*Picture)
		if img.data.Len() > 0 {
			// 1/3概率进行复用
			if s.picOffset%3 != 0 {
				img.data.Reset()
			}
		}
		if img.data.Len() == 0 {
			img.chars, err = s.Gra.Draw(&img.data)
			if err != nil {
				continue
			}
		}

		s.picOffLine = append(s.picOffLine, img)
		if len(s.picOffLine) >= s.imgLen {
			// 替换
			online := s.picOnLine
			s.picOnLine = s.picOffLine
			for i, _ := range online {
				s.poolBuffer.Put(online[i])
			}
			online = online[0:0]
			s.picOffLine = s.picOffLine[0:0]
		}
	}
}

func (s *Service) Draw() (pic *Picture, err error) {
	pic = s.picOnLine[s.picOffset]
	return pic, nil
}
