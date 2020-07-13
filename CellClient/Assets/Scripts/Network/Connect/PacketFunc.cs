using System;
using System.IO;
using System.Text;
using UnityEngine;

//创建数据包帮助函数
public class PacketFunc
{
    //创建接收数据包
    public static ReadPacket CreateReadPacket(byte[] msg)
    {
        return new ReadPacket(msg);
    }

    //创建空数据包
    public static WritePacket CreateNullPacket()
    {
        WritePacket wp = new WritePacket();
        wp.WriteInt8(MsgOutsideDef.MsgNull);
        return wp;
    }

    //创建单发数据包
    public static WritePacket CreatePeerPacket(Int64 to)
    {
        WritePacket wp = new WritePacket();
        wp.WriteInt8(MsgOutsideDef.MsgPeer);
        wp.WriteInt64(to);
        return wp;
    }

    //创建群发数据包
    public static WritePacket CreateBroadcastPacket(bool notself)
    {
        WritePacket wp = new WritePacket();
        wp.WriteInt8(MsgOutsideDef.MsgBroadcast);
        wp.WriteInt64(0);
        wp.WriteInt8(notself ? (byte)1 : (byte)0);
        return wp;
    }
}

public class ReadPacket
{
    private BinaryReader reader;

    public ReadPacket(byte[] msg)
    {
        MemoryStream memory = new MemoryStream(msg);
        reader = new BinaryReader(memory);
    }

    public static string Base64Decode(string base64EncodedData)
    {
        var base64EncodedBytes = Convert.FromBase64String(base64EncodedData);
        return Encoding.UTF8.GetString(base64EncodedBytes);
    }

    public byte[] ReadBytes()
    {
        short size = reader.ReadInt16();
        return reader.ReadBytes(size);
    }

    public string ReadString()
    {
        short size = ReadInt16();
        return Base64Decode(Encoding.UTF8.GetString(reader.ReadBytes(size)));
    }

    public byte ReadInt8()
    {
        return reader.ReadByte();
    }

    public short ReadInt16()
    {
        byte[] bData = BitConverter.GetBytes(reader.ReadInt16());
        Array.Reverse(bData);
        return BitConverter.ToInt16(bData, 0);
    }

    public Int32 ReadInt32()
    {
        byte[] bData = BitConverter.GetBytes(reader.ReadInt32());
        Array.Reverse(bData);
        return BitConverter.ToInt32(bData, 0);
    }

    public Int64 ReadInt64()
    {
        byte[] bData = BitConverter.GetBytes(reader.ReadInt64());
        Array.Reverse(bData);
        return BitConverter.ToInt64(bData, 0);
    }

    public float ReadFloat32()
    {
        byte[] bData = BitConverter.GetBytes(reader.ReadSingle());
        Array.Reverse(bData);
        return BitConverter.ToSingle(bData, 0);
    }

    public double ReadFloat64()
    {
        byte[] bData = BitConverter.GetBytes(reader.ReadDouble());
        Array.Reverse(bData);
        return BitConverter.ToDouble(bData, 0);
    }

    public void Close()
    {
        reader.Close();
    }
}

public class WritePacket
{
    private MemoryStream memory;
    private BinaryWriter writer;

    public WritePacket(int size = 1024)
    {
        memory = new MemoryStream(size);
        writer = new BinaryWriter(memory);
    }

    public static string Base64Encode(string plainText)
    {
        byte[] plainTextBytes = Encoding.UTF8.GetBytes(plainText);
        return Convert.ToBase64String(plainTextBytes);
    }

    public byte[] ToArray()
    {
        return memory.ToArray();
    }

    public void WriteBytes(byte[] b)
    {
        writer.Write((short)b.Length);
        writer.Write(b);
    }

    public void WriteString(string str)
    {
        var strb = Base64Encode(str);
        WriteInt16((short)strb.Length);
        writer.Write(Encoding.UTF8.GetBytes(strb));
    }

    public void WriteInt8(byte i)
    {
        writer.Write((byte)i);
    }

    public void WriteInt16(short i)
    {
        byte[] bData = BitConverter.GetBytes(i);
        Array.Reverse(bData);
        writer.Write(bData);
    }

    public void WriteInt32(Int32 i)
    {
        byte[] bData = BitConverter.GetBytes(i);
        Array.Reverse(bData);
        writer.Write(bData);
    }

    public void WriteInt64(Int64 i)
    {
        byte[] bData = BitConverter.GetBytes(i);
        Array.Reverse(bData);
        writer.Write(bData);
    }

    public void WriteFloat32(float f)
    {
        byte[] bData = BitConverter.GetBytes(f);
        Array.Reverse(bData);
        writer.Write(bData);
    }

    public void WriteFloat64(double f)
    {
        byte[] bData = BitConverter.GetBytes(f);
        Array.Reverse(bData);
        writer.Write(bData);
    }

    public void Close()
    {
        writer.Close();
    }
}

public class ByteArray
{
    private MemoryStream _memoryStream;

    private const byte BooleanFalse = 2;
    private const byte BooleanTrue = 3;

    public ByteArray()
    {
        _memoryStream = new MemoryStream();
    }

    public ByteArray(MemoryStream ms)
    {
        _memoryStream = ms;
    }

    public ByteArray(byte[] buffer)
    {
        _memoryStream = new MemoryStream();
        _memoryStream.Write(buffer, 0, buffer.Length);
        _memoryStream.Position = 0;
    }
    public void dispose()
    {
        if (_memoryStream != null)
        {
            _memoryStream.Close();
            _memoryStream.Dispose();
        }
        _memoryStream = null;
    }

    public uint Length
    {
        get
        {
            return (uint)_memoryStream.Length;
        }
    }

    public uint Position
    {
        get { return (uint)_memoryStream.Position; }
        set { _memoryStream.Position = value; }
    }

    public uint BytesAvailable
    {
        get { return Length - Position; }
    }


    public byte[] GetBuffer()
    {
        return _memoryStream.GetBuffer();
    }

    public byte[] ToArray()
    {
        return _memoryStream.ToArray();
    }

    public MemoryStream MemoryStream
    {
        get
        {
            return _memoryStream;
        }
    }

    // Read

    public bool ReadBoolean()
    {
        return _memoryStream.ReadByte() == BooleanTrue;
    }

    public byte ReadByte()
    {
        return (byte)_memoryStream.ReadByte();
    }

    public void ReadBytes(byte[] bytes, uint offset, uint length)
    {
        _memoryStream.Read(bytes, (int)offset, (int)length);
    }

    public void ReadBytes(ByteArray bytes, uint offset, uint length)
    {
        uint tmp = bytes.Position;
        int count = (int)(length != 0 ? length : BytesAvailable);
        for (int i = 0; i < count; i++)
        {
            bytes._memoryStream.Position = i + offset;
            bytes._memoryStream.WriteByte(ReadByte());
        }
        bytes.Position = tmp;
    }

    private byte[] priReadBytes(uint c)
    {
        byte[] a = new byte[c];
        for (uint i = 0; i < c; i++)
        {
            a[i] = (byte)_memoryStream.ReadByte();
        }
        return a;
    }

    public double ReadDouble()
    {
        byte[] bytes = priReadBytes(8);
        byte[] reverse = new byte[8];
        //Grab the bytes in reverse order 
        for (int i = 7, j = 0; i >= 0; i--, j++)
        {
            reverse[j] = bytes[i];
        }
        double value = System.BitConverter.ToDouble(reverse, 0);
        return value;
    }

    public float ReadFloat()
    {
        byte[] bytes = priReadBytes(4);
        byte[] invertedBytes = new byte[4];
        //Grab the bytes in reverse order from the backwards index
        for (int i = 3, j = 0; i >= 0; i--, j++)
        {
            invertedBytes[j] = bytes[i];
        }
        float value = System.BitConverter.ToSingle(invertedBytes, 0);
        return value;
    }

    public int ReadInt()
    {
        byte[] bytes = priReadBytes(4);
        int value = (int)((bytes[0] << 24) | (bytes[1] << 16) | (bytes[2] << 8) | bytes[3]);
        return value;
    }

    public short ReadShort()
    {
        byte[] bytes = priReadBytes(2);
        return (short)((bytes[0] << 8) | bytes[1]);
    }

    public string ReadUTF()
    {
        byte[] bytes = priReadBytes(2);
        uint length = (ushort)(((bytes[0] & 0xff) << 8) | (bytes[1] & 0xff));
        return ReadUTFBytes(length);
    }

    public string ReadUTFBytes(uint length)
    {
        if (length == 0)
            return string.Empty;
        UTF8Encoding utf8 = new UTF8Encoding(false, true);
        byte[] encodedBytes = priReadBytes((uint)length);
        string decodedString = utf8.GetString(encodedBytes, 0, encodedBytes.Length);
        return decodedString;
    }

    // Write

    public void WriteBoolean(bool value)
    {
        WriteByte((byte)(value ? BooleanTrue : BooleanFalse));
    }

    public void WriteByte(byte value)
    {
        _memoryStream.WriteByte(value);
    }

    public void WriteBytes(byte[] bytes, int offset, int length)
    {
        for (int i = offset; i < offset + length; i++)
            _memoryStream.WriteByte(bytes[i]);
    }

    public void WriteDouble(double value)
    {
        byte[] bytes = System.BitConverter.GetBytes(value);
        WriteBigEndian(bytes);
    }

    private void WriteBigEndian(byte[] bytes)
    {
        if (bytes == null)
            return;
        for (int i = bytes.Length - 1; i >= 0; i--)
        {
            WriteByte(bytes[i]);
        }
    }

    public void WriteFloat(float value)
    {
        byte[] bytes = System.BitConverter.GetBytes(value);
        WriteBigEndian(bytes);
    }

    public void WriteInt(int value)
    {
        WriteInt32(value);
    }

    private void WriteInt32(int value)
    {
        byte[] bytes = System.BitConverter.GetBytes(value);
        WriteBigEndian(bytes);
    }


    public void WriteShort(short value)
    {
        byte[] bytes = System.BitConverter.GetBytes((ushort)value);
        WriteBigEndian(bytes);
    }

    public void WriteUTF(string value)
    {
        UTF8Encoding utf8Encoding = new UTF8Encoding();
        int byteCount = utf8Encoding.GetByteCount(value);
        byte[] buffer = utf8Encoding.GetBytes(value);
        this.WriteShort((short)byteCount);
        if (buffer != null && buffer.Length > 0)
        {
            this.WriteBytes(buffer, 0, buffer.Length);
        }
    }
}
