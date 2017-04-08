package colorSpaceConverter

import (
	"image/color"
	"math"
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

func (lab1 *Lab) DeltaE(lab2 Lab) float64 {

	//Шаг1. Подготовка первичных переменных
	//Ci = sqrt(ai^2 + bi^2), где i=1,2
	C1 := math.Sqrt(math.Pow(lab1.A, 2.0) + math.Pow(lab1.B, 2.0)) 
	C2 := math.Sqrt(math.Pow(lab2.A, 2.0) + math.Pow(lab2.B, 2.0))
	
	//C = (C1 - C2)/2
	C := (C1 - C2)/2.0

	//G = 0.5*(1-sqrt(C'^7/(C'^7+25^7)))
	G	:= 0.5*(1-math.Sqrt(math.Pow(C,7.0)/(math.Pow(C,7.0)+math.Pow(25.0,7.0))))
	
	//ai' = (1+G)*ai, где i=1,2
	a1_ := (1+G)*lab1.A
	a2_ := (1+G)*lab2.A

	//Ci_ = sqrt(ai'^2 + b'^2), где i=1,2
	C1_ := math.Sqrt(math.Pow(a1_,2.0) + math.Pow(lab1.B,2.0))
	C2_ := math.Sqrt(math.Pow(a2_,2.0) + math.Pow(lab2.B,2.0))

	
	// hi_ = 0, при labi.B = 0 И ai' = 0
	// иначе hi_ = arctg(bi,ai')
	// где i=1,2
	var h1_, h2_ float64
	if lab1.B == 0.0 && a1_ == 0.0 {
		h1_ = 0.0
	} else {
		h1_ = math.Atan2(lab1.B,a1_)
	} 
	if lab2.B == 0.0 && a2_ == 0.0 {
		h2_ = 0.0
	} else {
		h2_ = math.Atan2(lab2.B,a2_)
	} 

	//Шаг2. Расчёт DL', DC', DH'
	
	deltaL_ := lab2.L - lab1.L
	deltaC_ := C2_ - C1_
	
	halfPi := math.Pi/2.0

	//dh' = 0,							при C1'*C2'=0
	//dh' = h2' - h1',				при C1'*C2'!=0 И |h2'-h1'|<=180град.
	//dh' = (h2' - h1')-360град.,	при C1'*C2'!=0 И |h2'-h1'|>180град.
	//dh' = (h2' - h1')+360град.,	при C1'*C2'!=0 И |h2'-h1'|<-180град.
	var dh_ float64
	if C1_*C2_ == 0.0 {
		dh_ = 0.0
	} else {
		if math.Abs(h2_ - h1_) <= halfPi {
			dh_ = h2_ - h1_
		}
		if math.Abs(h2_ - h1_) > halfPi {
			dh_ = (h2_ - h1_)-math.Pi
		}
		if math.Abs(h2_ - h1_) < -(halfPi) {
			dh_ = (h2_ - h1_)+math.Pi
		}
	}
	
	//DH' = 2*sqrt(C1'*C2')*sin(dh'/2)
	deltaH_ := 2.0*math.Sqrt(C1_*C2_)*math.Sin(dh_/2.0)

	//Шаг3. Расчёт DE00
	
	//L' = (L1 + L2)/2
	L_ := (lab1.L + lab2.L)/2.0
	
	//C' = (C1' + C2')/2
	C_ := (C1_ + C2_)/2.0
	
	//h' = (h1' + h2')/2, 				при |h1' - h2'|<=180град. И C1'*C2'!=0
	//h' = (h1' + h2'+360град.)/2, 	при 180град. > |h1' - h2'| И |h1' + h2'| < 360град. И C1'*C2'!=0
	//h' = (h1' + h2'-360град.)/2, 	при 180град. > |h1' - h2'| И |h1' + h2'| >= 360град. И C1'*C2'!=0
	//h' = (h1' + h2'), при C1'*C2'=0
	var h_ float64
	if C1_*C2_ == 0.0 {
		h_ = (h1_ + h2_)
	} else {
		if math.Abs(h1_ - h2_) <= halfPi {
			h_ = (h1_ + h2_)/2.0
		} else if halfPi > math.Abs(h1_ - h2_) && math.Abs(h1_ + h2_) < math.Pi {
			h_ = (h1_ + h2_ + math.Pi)/2.0
		} else if halfPi > math.Abs(h1_ - h2_) && math.Abs(h1_ + h2_) >= math.Pi {
			h_ = (h1_ + h2_ - math.Pi)/2.0
		}
	}
	
	//T = 1-0.1cos(h'-30град.)+0.24cos(2h') + 0.32cos(3h'+6град.)-0.20cos(4h'-63град.)
	T := 1.0 - 0.17 * math.Cos(h_-grad2rad(30.0))+0.24*math.Cos(2.0*h_)+0.32*math.Cos(3.0*h_+grad2rad(6.0))-0.20*math.Cos(4.0*h_-(grad2rad(63.0)))
	
	//DPhi = 30exp(-[((h'-275град.)/25)^2])
	deltaPhi := 30.0*math.Exp(-(math.Pow(((h_-grad2rad(275.0))/25.0),2.0) ))

	//Rc = 2sqrt(C'^7/(C'^7+25^7))
	RC :=	2.0*math.Sqrt(math.Pow(C_,7.0)/(math.Pow(C_,7.0)+math.Pow(25.0,7.0)))

	//Sl = 1+(0.015(L'-50)^2/sqrt(20+(L'-50)^2))
	SL := 1.0+(0.015*math.Pow(L_-50.0,2.0)/math.Sqrt(20.0+math.Pow(L_-50.0,2)))
	
	//Sc = 1+0.045C'
	SC := 1.0+0.045*C_	

	//Sh = 1+0.015C'T
	SH := 1.0+0.015*C_*T

	//Rt = -(sin(2DPhi)Rc
	RT := -(math.Sin(2*deltaPhi))*RC	

	// DE00 = sqrt((DL'/klSl)+(DC'/kcSc)+(DH'/hhSh)+Rt(DC'/kcSc)(DH'/khSh))
	return math.Sqrt( math.Pow((deltaL_/SL),2.0) + math.Pow((deltaC_/SC),2.0) + math.Pow((deltaH_/SH),2.0) + RT*(deltaC_/SC)*(deltaH_/SH) )	
}

func grad2rad(grad float64) float64 {
	return (math.Pi*grad)/180.0
}

func (this *XYZ) D65() XYZ {
	return XYZ{X:95.047, Y:100.000, Z:108.883}
}
