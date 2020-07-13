using System.Collections;
using System.Collections.Generic;
using UnityEngine;

//服务端外层包头
public class MsgOutsideDef
{
    public const byte MsgNull = 101;
    public const byte MsgPeer = 102;
    public const byte MsgBroadcast = 103;
    public const byte MsgRoomJoin = 104;
    public const byte MsgRoomLeave = 105;
    public const byte MsgRoomIds = 106;
}

//客户端自定义包头
public class MsgInsideDef
{
    public const short MsgError = 0;
    public const short MsgGameStart = 1;
    public const short MsgGameEnd = 2;
    public const short MsgChatText = 3;
    public const short MsgInputAction = 4;
}

//广播类型
public class MsgBroadcastType
{
    public const byte Peer = 0;
    public const byte Broadcast = 1;
    public const byte BroadcastNot = 2;
}

//网络状态
public class NetStatusType
{
    public const byte Open = 0;
    public const byte Ids = 1;
    public const byte Close = 2;
    public const byte Error = 3;
    public const byte Join = 4;
    public const byte Leave = 5;
}
