package benchmark

import (
	"CellServer/net"
	"testing"
)

func BenchmarkPacket(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wb := net.CreateWritePacket()
		wb.WriteString("hello")
		wb.WriteFloat64(30.5)
		wb.WriteInt32(300)
		wb.WriteString("hello")
		wb.WriteFloat64(30.5)
		wb.WriteInt32(300)

		rb := net.CreateReadPacket(wb.BytesBuffer().Bytes())
		_ = rb.ReadString()
		_ = rb.ReadFloat64()
		_ = rb.ReadInt32()
		_ = rb.ReadString()
		_ = rb.ReadFloat64()
		_ = rb.ReadInt32()
	}
}
