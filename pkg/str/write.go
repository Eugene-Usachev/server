package go_json

func (b *Builder) Write(buf byte) {
	b.buf = append(b.buf, buf)
}

func (b *Builder) WriteByte(c byte) {
	b.buf = append(b.buf, c)
}

func (b *Builder) WriteString(s string) {
	b.buf = append(b.buf, s...)
}
