using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class Demo : MonoBehaviour
{
    const int MinPlayerCount = 2;
    const int MaxPlayerCount = 2;

    bool isStart = false;
    bool isMain = false;

    Agent agent;

    [SerializeField]
    GameObject blue;

    [SerializeField]
    GameObject red;

    [SerializeField]
    Button startGame;

    // Start is called before the first frame update
    void Start()
    {
        agent = GameObject.FindObjectOfType<Agent>();

        Debug.Log(agent);

        //监听网络状态
        agent.SetStatusListener(OnStatus);

        //监听动作
        agent.SetActionListener(OnAction);
    }

    // Update is called once per frame
    void Update()
    {
        if (isStart && Input.GetMouseButtonDown(0))
        {
            Debug.Log("鼠标左键点击");
            var cmd = new CmdInputAction(MsgInsideDef.MsgInputAction, agent.GetConnectID(), InputType.ButtonC, MsgBroadcastType.Broadcast);
            WritePacket wp = cmd.Packet();
            wp.WriteInt32(isMain?1:0);
            cmd.Send(agent);
        }

        if (Input.GetMouseButtonDown(1))
        {
            blue.GetComponent<Rigidbody>().AddForce(0, 500, 0);
            Debug.Log("右键点击");
        }
    }

    public void OnStartGame()
    {
        if (!isStart && agent.GetIds().Count >= MinPlayerCount)
        {
            var cmd = new CmdInputAction(MsgInsideDef.MsgInputAction, agent.GetConnectID(), InputType.ButtonA, MsgBroadcastType.BroadcastNot);
            WritePacket wp = cmd.Packet();
            cmd.Send(agent);

            Debug.Log("Click startGame");
        }
        else
        {
            Debug.Log("Player Count: " + agent.GetIds().Count);
        }
    }

    public void OnStatus(byte t, Agent agent, long id = 0)
    {
        switch (t)
        {
            case NetStatusType.Open:
                Debug.Log("Status open");
                break;
            case NetStatusType.Ids:
                Debug.Log("Status ids");
                break;
            case NetStatusType.Close:
                Debug.Log("Status close");
                break;
            case NetStatusType.Error:
                Debug.Log("Status error");
                break;
            case NetStatusType.Join:
                Debug.Log("Status join: " + id);
                break;
            case NetStatusType.Leave:
                Debug.Log("Status leave: " + id);
                break;
        }
    }

    /// <summary>
    /// 动作回调
    /// </summary>

    public void OnAction(long id, byte action, ReadPacket rp)
    {
        switch(action)
        {
            case InputType.ButtonA:
                var cmd = new CmdInputAction(MsgInsideDef.MsgInputAction, agent.GetConnectID(), InputType.ButtonD, MsgBroadcastType.Peer, id);
                WritePacket wp = cmd.Packet();
                cmd.Send(agent);
                isStart = true;
                startGame.gameObject.SetActive(false);
                break;
            case InputType.ButtonC:
                Debug.Log("ButtonC");
                int isM = rp.ReadInt32();
                if (isM == 1)
                {
                    blue.GetComponent<Rigidbody>().AddForce(0, 500, 0);
                }
                else
                {
                    red.GetComponent<Rigidbody>().AddForce(0, 500, 0);
                }
                break;
            case InputType.ButtonD:
                isStart = true;
                isMain = true;
                startGame.gameObject.SetActive(false);
                break;
            case InputType.Fire2:
                Debug.Log("收到动作命令, ID: " + id + " 动作: " + action + " 附加数据: " + rp.ReadFloat32());
                break;
            default:
                Debug.Log("收到动作命令, ID: " + id + " 动作: " + action);
                break;
        }
    }
}
