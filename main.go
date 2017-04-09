package colorSpaceConverter

import (
	"image/color"
	"math"
   "strings"
   "strconv"
   "errors"
   "regexp"
	"log"
)

type XYZ struct {
	X float64
	Y float64
	Z float64
}

type Lab struct {
	L float64
	A float64
	B float64
}

type Normalize XYZ

func RGB2XYZ(c color.RGBA) XYZ {
	
	var r = float64(c.R)/255.0
	var g = float64(c.G)/255.0
	var b = float64(c.B)/255.0
	
	if r > 0.04045 {
		r = math.Pow(((r+0.055)/1.055),2.4)
	} else {
		r = r/12.92
	}
	if g > 0.04045 {
		g = math.Pow(((g+0.055)/1.055),2.4)
	} else {
		g = g/12.92
	}
	if b > 0.04045 {
		b = math.Pow(((b+0.055)/1.055),2.4)
	} else {
		r = b/12.92
	}

	r = r * 100
	g = g * 100
	b = b * 100


	return XYZ{
				X:r*0.4124+g*0.3576+b*0.1805,
					Y:r*0.2126+g*0.7152+b*0.0722,
					Z:r*0.0193+g*0.1192+b*0.9505}
}

func XYZ2Lab(c XYZ, n XYZ) Lab {

	var x = c.X/n.X
	var y = c.Y/n.Y
	var z = c.Z/n.Z
	var t = math.Pow(float64(6.0/29.0),3.0)

	if x > t	{
		x = greatestThenT(x)
	} else {
		x = lessThenT(x)
	}
	if y > t	{
		y = greatestThenT(y)
	} else {
		y = lessThenT(y)
	}
	if z > t	{
		z = greatestThenT(z)
	} else {
		z = lessThenT(z)
	}

	return Lab {
		L: (116*y)-16,
		A: 500*(x-y),
		B: 200*(y-z)}
}

func lessThenT(i float64) float64 {
	return float64(1.0/3.0)*math.Pow(float64(29.0/6.0),2.0)*i+float64(4.0/29.0)
} 

func greatestThenT(i float64) float64 {
	return math.Pow(i,float64(1.0/3.0))
}

func DeltaECIE2000(lab1,lab2 Lab, KL, KC, KH float64) float64 {

	C1 := math.Sqrt(math.Pow(lab1.A, 2.0) + math.Pow(lab1.B, 2.0)) 
	C2 := math.Sqrt(math.Pow(lab2.A, 2.0) + math.Pow(lab2.B, 2.0))

	C := (C1 + C2)/2.0

	G	:= 0.5*(1 - math.Sqrt(math.Pow(C,7.0)/(math.Pow(C,7.0)+math.Pow(25.0,7.0))) )

	a1_ := (1+G)*lab1.A
	a2_ := (1+G)*lab2.A

	C1_ := math.Sqrt(math.Pow(a1_,2.0) + math.Pow(lab1.B,2.0))
	C2_ := math.Sqrt(math.Pow(a2_,2.0) + math.Pow(lab2.B,2.0))

	
	var h1_, h2_ float64
	if lab1.B == 0.0 && a1_ == 0.0 {
		h1_ = 0.0
	} else {
		h1_ = 360-math.Abs(rad2deg(math.Atan2(lab1.B,a1_)))
	} 
	if lab2.B == 0.0 && a2_ == 0.0 {
		h2_ = 0.0
	} else {
		h2_ = 360-math.Abs(rad2deg(math.Atan2(lab2.B,a2_)))
	} 

	
	deltaL_ := lab2.L - lab1.L
	deltaC_ := C2_ - C1_
	

	var dh_ float64
	if C1_*C2_ == 0.0 {
		dh_ = 0.0
	} else {
		if math.Abs(h2_ - h1_) <= 180 {
			dh_ = h2_ - h1_
		}
		if math.Abs(h2_ - h1_) > 180 {
			dh_ = (h2_ - h1_)-360
		}
		if math.Abs(h2_ - h1_) < -(180) {
			dh_ = (h2_ - h1_)+360
		}
	}
	
	deltaH_ := 2.0*math.Sqrt(C1_*C2_)*math.Sin(deg2rad(dh_)/2.0)

	L_ := (lab1.L + lab2.L)/2.0
	
	C_ := (C1_ + C2_)/2.0
	
	var h_ float64
	if C1_*C2_ == 0.0 {
		h_ = (h1_ + h2_)
	} else {
		if math.Abs(h1_ - h2_) <= 180 {
			h_ = (h1_ + h2_)/2.0
		} else if 180 > math.Abs(h1_ - h2_) && math.Abs(h1_ + h2_) < 360 {
			h_ = (h1_ + h2_ + math.Pi)/2.0
		} else if 180 > math.Abs(h1_ - h2_) && math.Abs(h1_ + h2_) >= 360 {
			h_ = (h1_ + h2_ - math.Pi)/2.0
		}
	}

	T := 1.0 - 0.17 * math.Cos(deg2rad(h_-30.0))+0.24*math.Cos(deg2rad(2.0*h_))+0.32*math.Cos(deg2rad(3.0*h_+6.0))-0.20*math.Cos(deg2rad(4.0*h_-63.0))
	log.Printf("T: %f",T)

	deltaTeta := deg2rad(30.0*math.Exp(-(math.Pow(((h_-275.0)/25.0),2.0) )))

	RC :=	2.0*math.Sqrt(math.Pow(C_,7.0)/(math.Pow(C_,7.0)+math.Pow(25.0,7.0)))

	SL := 1.0+(0.015*math.Pow(L_-50.0,2.0)/math.Sqrt(20.0+math.Pow(L_-50.0,2)))
	
	SC := 1.0+0.045*C_	

	SH := 1.0+0.015*C_*T

	RT := -(math.Sin(2*deltaTeta))*RC	

	return math.Sqrt( math.Pow((deltaL_/SL*KL),2.0) + math.Pow((deltaC_/SC*KC),2.0) + math.Pow((deltaH_/SH*KH),2.0) + RT*(deltaC_/SC*KC)*(deltaH_/SH*KH) )	
}


func deg2rad(deg float64) float64 {
	return (math.Pi*deg)/180.0
}

func rad2deg(rad float64) float64 {
	return ((180.0*rad)/math.Pi)
}

func (this *Normalize) D65() XYZ {
	return XYZ{X:95.047, Y:100.0, Z:108.883}
}

func Hex2RGB(s string) (color.RGBA, error) {
   re, err := regexp.Compile("(?i)"+"#?[a-z0-9]{2}[a-z0-9]{2}[a-z0-9]{2}")
   if err != nil {
      return color.RGBA{}, err
   }

   if len(s) > 7 || !re.MatchString(s) {
      return color.RGBA{}, errors.New("Isn't hex-color.")
   }

   ss := strings.Split(s, "")
   if ss[0] == "#" {
      ss = ss[1:]
   }

   r, err := strconv.ParseInt(ss[0]+ss[1], 16, 64)
   if err != nil {
      return color.RGBA{}, err
   }
   g, err := strconv.ParseInt(ss[2]+ss[3], 16, 64)
   if err != nil {
      return color.RGBA{}, err
   }
   b, err := strconv.ParseInt(ss[4]+ss[5], 16, 64)
   if err != nil {
      return color.RGBA{}, err
   }

   return color.RGBA{R:uint8(r),G:uint8(g),B:uint8(b),A:uint8(255)}, nil
}

func RGB2Hex(c color.RGBA) string {
	var hex string
	
	hex += dec2hex(c.R)
	hex += dec2hex(c.G)
	hex += dec2hex(c.B)

	return hex
}

func dec2hex(dig uint8) string {
	return fmt.Sprintf("%02x", dig)
}
