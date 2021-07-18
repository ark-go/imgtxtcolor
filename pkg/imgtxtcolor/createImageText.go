package imgtxtcolor

import (
	"image"
	"log"
	"strings"
	"time"
)

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("Time [%v]: %v\n", msg, time.Since(start))
}

func CreateImageText(text string, opt *stStartOptions) ([]*image.RGBA, error) {
	if opt.FontSizeInt < 3 {
		opt.FontSizeInt = 3
	}
	// maximize CPU usage for maximum performance
	// runtime.GOMAXPROCS(runtime.NumCPU())
	defer duration(track("Все"))  // меряем время
	param, err := initCanvas(opt) // инициализация первого Canvas
	if err != nil {
		return nil, err
	}
	//tmpCanvas := []*image.RGBA{}                            // массив для сбора всех Canvas
	//param.allImages = append(param.allImages, param.canvas) // сразу положим туда первую Canvas

	// Выделяем строки
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") //\r\n у windows
	// Перебираем строки

	for _, line := range lines { // строки

		if err := parseAndDrawLine(param, strings.TrimRight(line, " ")); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	return param.allImages, err
}
