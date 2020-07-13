package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

func CreateRoomEntity() *RoomEntity {
	return &RoomEntity{
		ID:       RoomID(),
		Sessions: sync.Map{},
	}
}

func CreateRoomPacketEntity(t int8, rid string, not bool, from int64, to int64, data []byte) *RoomPacketEntity {
	return &RoomPacketEntity{
		Type:    t,
		RoomID:  rid,
		NotSelf: not,
		From:    from,
		To:      to,
		Data:    data,
	}
}

//房间实体
type RoomEntity struct {
	ID        string
	EmptyTime int64
	Sessions  sync.Map
}

//房间数据包实体
type RoomPacketEntity struct {
	Type    int8
	RoomID  string
	NotSelf bool
	From    int64
	To      int64
	Data    []byte
}

var roomId int64

func RoomID() string {
	shortUrl, err := TransformShortUrl(strconv.FormatInt(atomic.AddInt64(&roomId, 1), 10))
	if err == nil {
		return shortUrl[0][4:] + shortUrl[1][4:] + shortUrl[2][4:] + shortUrl[3][4:]
	} else {
		data := []byte(strconv.FormatInt(atomic.AddInt64(&roomId, 1), 10))
		has := md5.Sum(data)
		md5str := fmt.Sprintf("%x", has)
		return md5str[2:4] + md5str[10:12] + md5str[18:20] + md5str[26:28]
	}
}

func TransformShortUrl(longURL string) ([4]string, error) {
	alphabet := []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	md5Str := func(str string) string {
		m := md5.New()
		m.Write([]byte(str))
		c := m.Sum(nil)
		return hex.EncodeToString(c)
	}(longURL)
	//var hexVal int64
	var tempVal int64
	var result [4]string
	var tempUri []byte
	for i := 0; i < 4; i++ {
		tempSubStr := md5Str[i*8 : (i+1)*8]
		hexVal, err := strconv.ParseInt(tempSubStr, 16, 64)
		if err != nil {
			return result, nil
		}
		tempVal = int64(0x3FFFFFFF) & hexVal
		var index int64
		tempUri = []byte{}
		for i := 0; i < 6; i++ {
			index = 0x0000003D & tempVal
			tempUri = append(tempUri, alphabet[index])
			tempVal = tempVal >> 5
		}
		result[i] = string(tempUri)
	}
	return result, nil
}
