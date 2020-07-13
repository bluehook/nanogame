using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public interface Command
{
    short Type();
    WritePacket Packet();
    void Send(Agent agent);
    void Handle(Agent agent, ReadPacket rp);
}
