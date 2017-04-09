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
	1-й XYZ-аргумент: цвет для преобразования
	2-й XYZ-аргумент: цвет нормализации (см. функции типа colorSpace.Normalize)

	Пример:
		var normalize colorSpace.Normalize
		XYZ2Lab(RGB2XYZ(color.RGBA{30,20,10,255}), normalize.D65())
	
##### DeltaECIE2000(Lab, Lab, float64, float64, float64) float64
	1-й, 2-й Lab-аргументы: сравниваемые цвета
	3-й, 4-й, 5-й float64-аргументы: коэффициенты Kl, Kc, Kh - задавать равным 1.0

##### Normalize.D65() XYZ 
	возвращает цветовую "нормализацию" D65

##### Hex2RGB(string) (color.RGBA, error)
	string-аргумент: hex-формат цвета

	Пример:
		#ff0032
		ff0032

##### RGB2Hex(color.RGBA) string
	возвращает цвет строкой hex-формате

