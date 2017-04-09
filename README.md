# Пакет colorSpace
#### Конвертация некоторых цветовых пространств между собой


## Типы

* XYZ {
   X: float64,
   Y: float64,
   Z: float64
}

* Lab {
	L: float64,
	A: float64,
	B: float64
}

* Normalize bool

## Функции

##### RGB2XYZ(color.RGB) XYZ

##### XYZ2Lab(XYZ, XYZ) Lab
	прв. XYZ-параметр: цвет для приобразования
	втр. XYZ-параметр: цвет нормализации (см. функции типа colorSpace.Normalize)

	Пример:
		var normalize colorSpace.Normalize
		XYZ2Lab(RGB2XYZ(color.RGBA{30,20,10,255}), normalize.D65())
	
##### DeltaECIE2000(Lab, Lab, float64, float64, float64) float64
	прв., втр. Lab-параметры: сравниваемые цвета
	тр., чтв., пт. float64-параметры: коэффициенты (=1.0)

##### Normalize.D65() XYZ 
	возвращает цветовую "нормализацию" D65

##### Hex2RGB(string) (color.RGBA, error)
	string-параметр: hex-формат цвета

	Пример:
		#ff0032
		ff0032

##### RGB2Hex(color.RGBA) string

