package tlv

import "github.com/Ovler-Young/MiraiGo/binary"

func T545(qimei []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x545)
		w.WriteBytesShort(qimei)
	})
}
