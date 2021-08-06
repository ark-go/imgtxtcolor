package imgtxtcolor

import "golang.org/x/image/math/fixed"

// Вычисляет левую точку для вставки текста
//	str - одна строка
func (p *stParam) getHorizontalPos(str string) fixed.Int26_6 {
	switch p.opt.AlignHorizontal {
	case AlignHorizontalLeft:
		return fixed.I(p.canvas.padding.left) // влево
	case AlignHorizontalCenter: // для центра
		max := fixed.I(p.canvas.img.Rect.Max.X) // всего
		max -= fixed.I(p.canvas.padding.right)  // отнимаем справа
		max -= fixed.I(p.canvas.padding.left)   // отнимаем слева
		max -= p.drw.MeasureString(str)         // получаем свободное место
		max /= 2                                //место пополам
		max += fixed.I(p.canvas.padding.left)   // отодвигаем слева
		return max
	case AlignHorizontalRight: // вправо
		return (fixed.I(p.canvas.img.Rect.Dx()) - p.drw.MeasureString(str) - fixed.I(p.canvas.padding.right))
	default:
		return fixed.I(p.canvas.padding.left) // влево
	}
}
