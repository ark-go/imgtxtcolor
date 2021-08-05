# imgtxtcolor

Все команды действуют с начала вставки и до конца текста, их можно переназначать,
вставляются в текст command:value и должны отделятся пробелами, все пробелы после команды будут убраны

### В начале строки

команды для строк, вставляются в любом месте но начинают действовать с начала текущей строки

- fontSize, size
- fontColor, color
- lineSpacing - межстрочное расстояние
- align
  - left
  - center
  - right

### Для следущего Image

Следущие команды начнут новый Image, они должны указываться подряд в любом порядке с начала строки,
между ними не должно быть текста, только пробелы, иначе будут создаваться новые страницы-image

- alignH, alignV
  - top
  - center
  - bottom
- padding
- paddingTop, top
- paddingLeft, left
- paddingRigh, right
- paddingBottom, bottom
- round - радиус углов
  - auto или не число, радиус углов будет в половину размера paddingTop
- bgcolor - цвет фона
- width
- height
- rect:tg

### Завершить Image

принудительно разорвать текст

- break:page

time:now - вставит время
