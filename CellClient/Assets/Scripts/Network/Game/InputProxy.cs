//动作代理
public class InputProxy
{
    public delegate void Callback(long id, byte t, ReadPacket rp);
    public Callback ActionCallBack;

    public bool IsAi = true;

    public void OnActionTriger(long id, byte t, ReadPacket rp)
    {
        if (ActionCallBack != null)
        {
            ActionCallBack(id, t, rp);
        }
    }
}
