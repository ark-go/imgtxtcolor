package imgtxtcolor

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type logtype int

const (
	LogFileAndConsole logtype = iota
	LogFile
	LogConsole
	LogOff
)

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("Time [%v]: %v\n", msg, time.Since(start))
}

func startLogFile(logt logtype) (*os.File, error) {
	var wr []io.Writer
	var f *os.File = nil
	var err error
	switch logt {
	case LogFileAndConsole, LogFile:
		f, err = os.OpenFile("canvas.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err) // Exit(1)
		}
		wr = append(wr, f)
		if logt == LogFile {
			break
		}
		wr = append(wr, os.Stdout)
	case LogConsole:
		wr = append(wr, os.Stdout)
	case LogOff:
		wr = []io.Writer{ioutil.Discard} //- отключить вывод
	}
	wrt := io.MultiWriter(wr...)
	log.SetOutput(wrt)
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//logger := log.New(f, "префикс: ", log.LstdFlags)
	//log.SetOutput(ioutil.Discard) //- отключить вывод
	return f, err

}

// Старт !!!
func CreateImageText(text string, opt *stStartOptions) ([]*ImgCanvas, error) {
	return CreateImageTextLog(text, opt, LogOff)
}
func CreateImageTextLog(text string, opt *stStartOptions, logt logtype) ([]*ImgCanvas, error) {
	// maximize CPU usage for maximum performance
	// runtime.GOMAXPROCS(runtime.NumCPU())

	if f, err := startLogFile(logt); err != nil {
		log.Panic("ошибка создания log-файла ")
	} else if f != nil {
		defer f.Close()
	}

	defer duration(track("Все")) // меряем время
	startTime := time.Now()
	param, err := initCanvas(opt) // инициализация параметров
	if err != nil {
		return nil, err
	}
	// Выделяем строки
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") //\r\n у windows
	// Перебираем строки
	for _, line := range lines { // строки

		if err := param.parseAndDrawLine(strings.TrimRight(line, " ")); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//param.textAlign() // последний Canvas
	param.isNewCanvas = true
	log.Printf("Time [%v]: %v\n", "Текст пройден", time.Since(startTime))
	param.formatAllCanvas()
	log.Printf("Time [%v]: %v\n", "Формат пройден", time.Since(startTime))
	if param.opt.GifFileName != "" {
		param.ToGif()
	}
	log.Printf("Time [%v]: %v\n", "Gif создан", time.Since(startTime))
	return param.allCanvas, err
}
