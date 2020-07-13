using System;
using System.Text;
using UnityEngine;
using BestHTTP;
using BestHTTP.WebSocket;

public class WsClient : Singleton<WsClient>
{
    public delegate void callback();
    public callback OpenCallBack;
    public callback ErrorCallBack;
    public callback CloseCallBack;

    public delegate void bincallback(WebSocket ws, byte[] message);
    public bincallback BinCallBack;

    public string Addr
    {
        get { return addr; }
        set { addr = value; }
    }
    private string addr;

    private WebSocket socket = null;

    public void Start()
    {
    }

    public void Connect()
    {
        socket = new WebSocket(new Uri(addr));
        socket.OnOpen += OnOpen;
        socket.OnBinary += OnBinaryMessage;
        socket.OnError += OnError;
        socket.OnClosed += OnClosed;
        socket.Open();
    }

    public void Send(string str)
    {
        socket.Send(str);
    }

    public void Send(byte[] buf)
    {
        socket.Send(buf);
    }

    public void Close()
    {
        socket.Close();
    }

    public bool IsOpen()
    {
        return (socket != null && socket.IsOpen);
    }

    public void OnOpen(WebSocket ws)
    {
        OpenCallBack();
    }

    public void OnBinaryMessage(WebSocket ws, byte[] message)
    {
        BinCallBack(ws, message);
    }

    public void OnClosed(WebSocket ws, UInt16 code, string message)
    {
        CloseCallBack();
    }

    public void OnError(WebSocket ws, Exception ex)
    {
        ErrorCallBack();
        Debug.Log("An error occured: " + ex.Message);
    }

    public void OnDestroy()
    {
        if (socket != null && socket.IsOpen)
        {
            socket.Close();
        }
    }

    public byte[] GetBytes(string message)
    {
        byte[] buffer = Encoding.Default.GetBytes(message);
        return buffer;
    }
}