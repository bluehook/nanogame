using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class CmdChatText : Command
{
    short cmdType = MsgInsideDef.MsgError;
    WritePacket wp;

    public CmdChatText(short t, string msg, bool not)
    {
        cmdType = t;
        wp = PacketFunc.CreateBroadcastPacket(not);  //广播是否包括自己
        wp.WriteInt16(MsgInsideDef.MsgChatText);
        wp.WriteString(msg);
    }

    public short Type()
    {
        return cmdType;
    }

    public WritePacket Packet()
    {
        return wp;
    }

    public void Send(Agent agent)
    {
        agent.Send(wp.ToArray());
    }

    public void Handle(Agent agent, ReadPacket rp)
    {
        string msg = rp.ReadString();
        Debug.Log(msg);
    }
}
