package jazzavatar

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"strconv"
	"strings"

	m "bitbucket.org/neetsdkasu/mersenne_twister_go"
	colorful "github.com/lucasb-eyer/go-colorful"
)

// Jazzavatar Jazzavatar Model
type Jazzavatar struct {
	Name    string
	Size    string
	Radius  string
	BgColor string
	Shapes  []shapeData

	colors   []string
	mersenne *m.MersenneTwister
}

type shapeData struct {
	Color  string
	Tx     float64
	Ty     float64
	Rotate float64
	Center float64
}

// colorList Color List
var colorList = []string{
	"#01888C", // teal
	"#FC7500", // bright orange
	"#034F5D", // dark teal
	"#F73F01", // orangered
	"#FC1960", // magenta
	"#C7144C", // raspberry
	"#F3C100", // goldenrod
	"#1598F2", // lightning blue
	"#2465E1", // sail blue
	"#F19E02", // gold
}

const shapeCount = 4

//Init Init Jazzavatar
func (ja *Jazzavatar) Init(name string, size string, radius string) (*Jazzavatar, error) {
	ja.Name = name
	ja.Size = size
	ja.Radius = radius
	ja.colors = colorList

	seed, err := ja.nameToSeed(name)
	if err != nil {
		return nil, err
	}
	ja.mersenne = m.NewMersenneTwister().Init(uint32(seed))

	remainingColors, err := ja.hueShift()
	if err != nil {
		return nil, err
	}

	bgColor, remainingColors := ja.genColor(remainingColors)
	ja.BgColor = bgColor

	for i := 0; i < shapeCount-1; i++ {
		shape, newRemainingColors, err := ja.genShape(remainingColors, i)
		if err != nil {
			return nil, err
		}
		remainingColors = newRemainingColors

		ja.Shapes = append(ja.Shapes, shape)
	}

	return ja, err
}

func (ja *Jazzavatar) nameToSeed(name string) (int64, error) {
	seed, err := ja.hexToNumber(name)
	if seed == 0 {
		seed, err = ja.stringToNumber(name)
	}
	return seed, err
}

func (ja *Jazzavatar) hexToNumber(hex string) (int64, error) {
	hexString := strings.Replace(hex, "0x", "", 1)
	if len(hexString) > 8 {
		hexString = hexString[:8]
	}
	return strconv.ParseInt(hexString, 16, 64)
}

func (ja *Jazzavatar) stringToNumber(str string) (int64, error) {
	md5Obj := md5.New()
	md5Obj.Write([]byte(str))
	md5Str := hex.EncodeToString(md5Obj.Sum(nil))
	return ja.hexToNumber(md5Str)
}

func (ja *Jazzavatar) hueShift() ([]string, error) {
	amount := ja.mersenne.Real2()*30 - 15

	var colors []string

	for _, c := range ja.colors {
		color, err := colorful.Hex(c)
		if err != nil {
			return nil, err
		}

		h, s, v := color.Hsv()
		h = h + float64(amount)
		newColor := colorful.Hsv(h, s, v)
		newC := newColor.Hex()

		colors = append(colors, newC)
	}

	return colors, nil
}

func (ja *Jazzavatar) genColor(colors []string) (string, []string) {
	_ = ja.mersenne.Real2()
	index := int(math.Floor(float64(len(colors)) * ja.mersenne.Real2()))
	genColor := colors[index:][0]
	colors = append(colors[:index], colors[index+1:]...)
	return genColor, colors
}

func (ja *Jazzavatar) genShape(colors []string, index int) (shapeData, []string, error) {
	diameter, err := strconv.ParseFloat(ja.Size, 64)
	if err != nil {
		return shapeData{}, nil, err
	}
	center := diameter / 2.0

	total := shapeCount - 1

	firstRot := ja.mersenne.Real2()
	angle := math.Pi * 2.0 * firstRot
	velocity := diameter/float64(total)*ja.mersenne.Real2() + (float64(index) * diameter / float64(total))

	tx := math.Cos(angle) * velocity
	ty := math.Sin(angle) * velocity

	secondRot := ja.mersenne.Real2()
	rot := firstRot*360.0 + secondRot*180.0
	rotate := math.Round(rot*10.0) / 10.0

	color, remainingColors := ja.genColor(colors)

	shap := shapeData{
		Color:  color,
		Tx:     tx,
		Ty:     ty,
		Rotate: rotate,
		Center: center,
	}

	return shap, remainingColors, nil
}
