using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class CmdInputAction : Command
{
    short cmdType = MsgInsideDef.MsgError;
    WritePacket wp;

    public CmdInputAction(short t, long aid, byte action, byte not, long to = 0)
    {
        cmdType = t;
        switch(not)
        {
            case MsgBroadcastType.Peer:
                wp = PacketFunc.CreatePeerPacket(to);
                break;
            case MsgBroadcastType.Broadcast:
                wp = PacketFunc.CreateBroadcastPacket(true);
                break;
            case MsgBroadcastType.BroadcastNot:
                wp = PacketFunc.CreateBroadcastPacket(false);
                break;
            default:
                wp = PacketFunc.CreateBroadcastPacket(true);
                break;
        }
        wp.WriteInt16(MsgInsideDef.MsgInputAction);
        wp.WriteInt8(action);
        wp.WriteInt64(aid);
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
        byte action = rp.ReadInt8();
        long aid = rp.ReadInt64();
        agent.Action(aid, action, rp);
    }
}
