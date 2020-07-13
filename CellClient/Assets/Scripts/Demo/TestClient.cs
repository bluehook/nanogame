using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using BestHTTP;
using BestHTTP.WebSocket;
using UnityEngine.UI;

public class TestClient : MonoBehaviour
{
    Agent agent;

    public void Awake()
    {
        //添加代理
        agent = gameObject.AddComponent<Agent>();

        //注册命令
        agent.AddCommand(new CmdChatText(MsgInsideDef.MsgChatText, "", false));
        agent.AddCommand(new CmdInputAction(MsgInsideDef.MsgInputAction, 0, 0, MsgBroadcastType.Broadcast));
    }

    /// <summary>
    /// 各种测试功能按钮
    /// </summary>

    //创建房间按钮
    public void OnCreateRoomBtn()
    {
        HTTPRequest request = new HTTPRequest(new System.Uri("http://127.0.0.1:36000/create"));
        request.Send();
    }

    //打开连接按钮
    public void OnOpenBtn()
    {
        agent.Init();
        agent.Open();
    }

    public void OnCloseBtn()
    {
        agent.Close();
    }

    //发送聊天信息
    public void OnChatOk()
    {
        string sendStr = GameObject.Find("ChatSendText").GetComponent<Text>().text;
        var cmd = new CmdChatText(MsgInsideDef.MsgChatText, sendStr, true);
        cmd.Send(agent);
    }

    //jump按钮
    public void OnJumpBtn()
    {
        var cmd = new CmdInputAction(MsgInsideDef.MsgInputAction, agent.GetConnectID(), InputType.Jump, MsgBroadcastType.Broadcast);
        cmd.Send(agent);
    }

    //fire1按钮
    public void OnFire1Btn()
    {
        var cmd = new CmdInputAction(MsgInsideDef.MsgInputAction, agent.GetConnectID(), InputType.Fire1, MsgBroadcastType.Broadcast);
        cmd.Send(agent);
    }

    //fire2按钮
    public void OnFire2Btn()
    {
        var cmd = new CmdInputAction(MsgInsideDef.MsgInputAction, agent.GetConnectID(), InputType.Fire2, MsgBroadcastType.Broadcast);
        WritePacket wp = cmd.Packet();
        wp.WriteFloat32(5.8f);
        cmd.Send(agent);
    }

}

