using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using BestHTTP;
using BestHTTP.WebSocket;
using UnityEngine.UI;

public class Agent : MonoBehaviour
{
    [Header("房间地址")]
    public string Addr = "ws://127.0.0.1:36000/majq2eRv";

    WsClient client;
    bool keepPing = false;

    //进入房间后自己的连接ID
    long myID = -1;

    //房间ID列表(包括自己)
    List<long> ids;
    InputProxy inputProxy = new InputProxy();

    //注册命令列表
    List<Command> responseCommands = new List<Command>();

    public delegate void StatusCallback(byte t, Agent agent, long id = 0);
    public StatusCallback statusCallBack;

    public void Init()
    {
        if (WsClient.Instance.IsOpen())
        {
            WsClient.Instance.Close();
        }

        client = WsClient.Instance;
        client.Addr = Addr;

        client.OpenCallBack -= onConnected;
        client.OpenCallBack += onConnected;

        client.ErrorCallBack -= onError;
        client.ErrorCallBack += onError;

        client.BinCallBack -= onMessage;
        client.BinCallBack += onMessage;

        client.CloseCallBack -= onClose;
        client.CloseCallBack += onClose;
    }

    public void Send(byte[] buf)
    {
        client.Send(buf);
    }

    public void Open()
    {
        client.Connect();
    }

    public void Close()
    {
        client.Close();
    }

    public void AddCommand(Command cmd)
    {
        responseCommands.Add(cmd);
    }

    public bool IsConnected()
    {
        return client.IsOpen();
    }

    public long GetConnectID()
    {
        return myID;
    }

    public List<long> GetIds()
    {
        return ids;
    }

    //设置网络状态监听
    public void SetStatusListener(StatusCallback cb)
    {
        statusCallBack -= cb;
        statusCallBack += cb;
    }

    //设置输入动作监听
    public void SetActionListener(InputProxy.Callback cb)
    {
        inputProxy.ActionCallBack -= cb;
        inputProxy.ActionCallBack += cb;
    }

    public void Action(long id, byte action, ReadPacket rp)
    {
        inputProxy.OnActionTriger(id, action, rp);
    }

    //处理服务端消息
    void onMessage(WebSocket ws, byte[] message)
    {
        ReadPacket rp = PacketFunc.CreateReadPacket(message);
        //包头
        byte t = rp.ReadInt8();
        if (t == MsgOutsideDef.MsgNull) return;

        //发送者ID,广播为0
        long clientID = rp.ReadInt64();
        switch (t)
        {
            case MsgOutsideDef.MsgPeer:  //单播消息
                short gameType = rp.ReadInt16();
                foreach (Command cmd in responseCommands)
                {
                    if (cmd.Type() == gameType)
                    {
                        cmd.Handle(this, rp);
                        break;
                    }
                }
                break;
            case MsgOutsideDef.MsgBroadcast:  //广播消息
                byte _ = rp.ReadInt8();
                gameType = rp.ReadInt16();
                foreach (Command cmd in responseCommands)
                {
                    if (cmd.Type() == gameType)
                    {
                        cmd.Handle(this, rp);
                        break;
                    }
                }
                break;
            case MsgOutsideDef.MsgRoomJoin:  //加入房间
                long joinID = rp.ReadInt64();
                bool canAdd = true;
                foreach(var id in ids)
                {
                    if (id == joinID)
                    {
                        canAdd = false;
                        break;
                    }
                }
                if (canAdd)
                {
                    ids.Add(joinID);
                    if (statusCallBack != null)
                    {
                        statusCallBack(NetStatusType.Join, this, joinID);
                    }
                }
                Debug.Log("Room join id: " + joinID.ToString());
                break;
            case MsgOutsideDef.MsgRoomLeave:  //离开房间
                long leaveID = rp.ReadInt64();
                ids.Remove(leaveID);
                if (statusCallBack != null)
                {
                    statusCallBack(NetStatusType.Leave, this, leaveID);
                }
                Debug.Log("Room leave id: " + leaveID.ToString());
                break;
            case MsgOutsideDef.MsgRoomIds:  //获得房间ID列表
                myID = rp.ReadInt64();
                Debug.Log("myID: " + myID.ToString());
                ids = new List<long>();
                ids.Add(myID);
                long otherID;
                while ((otherID = rp.ReadInt64()) != 0)
                {
                    ids.Add(otherID);
                }
                if (statusCallBack != null)
                {
                    statusCallBack(NetStatusType.Ids, this);
                }
                Debug.Log("Room id count: " + ids.ToArray().Length.ToString());
                break;
        }
    }

    //连接成功
    void onConnected()
    {
        //开始心跳
        keepPing = true;
        StartCoroutine("ping");
        if (statusCallBack != null)
        {
            statusCallBack(NetStatusType.Open, this);
        }
    }

    //服务端断开连接
    void onError()
    {
        if (statusCallBack != null)
        {
            statusCallBack(NetStatusType.Error, this);
        }
    }

    void onClose()
    {
        keepPing = false;
        if (statusCallBack != null)
        {
            statusCallBack(NetStatusType.Close, this);
        }
    }

    IEnumerator ping()
    {
        while (keepPing && true)
        {
            WritePacket wp = PacketFunc.CreateNullPacket();
            client.Send(wp.ToArray());
            yield return new WaitForSeconds(8);
        }
    }
}
