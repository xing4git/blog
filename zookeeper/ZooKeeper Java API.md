ZooKeeper Java API
----

ZooKeeper提供了Java和C的binding. 本文关注Java相关的API.

### 准备工作
拷贝ZooKeeper安装目录下的zookeeper.x.x.x.jar文件到项目的classpath路径下.

### 创建连接和回调接口

首先需要创建ZooKeeper对象, 后续的一切操作都是基于该对象进行的.
<pre name="code" class="java">
ZooKeeper(String connectString, int sessionTimeout, Watcher watcher) throws IOException
</pre>
以下为各个参数的详细说明:
+ connectString. zookeeper server列表, 以逗号隔开. ZooKeeper对象初始化后, 将从server列表中选择一个server, 并尝试与其建立连接. 如果连接建立失败, 则会从列表的剩余项中选择一个server, 并再次尝试建立连接.
+ sessionTimeout. 指定连接的超时时间. 
+ watcher. 事件回调接口.

注意, 创建ZooKeeper对象时, 只要对象完成初始化便立刻返回. 建立连接是以异步的形式进行的, 当连接成功建立后, 会回调watcher的process方法. 如果想要同步建立与server的连接, 需要自己进一步封装.
<pre name="code" class="java">
public class ZKConnection {
	/**
	 * server列表, 以逗号分割
	 */
	protected String hosts = "localhost:4180,localhost:4181,localhost:4182";
	/**
	 * 连接的超时时间, 毫秒
	 */
	private static final int SESSION_TIMEOUT = 5000;
	private CountDownLatch connectedSignal = new CountDownLatch(1);
	protected ZooKeeper zk;

	/**
	 * 连接zookeeper server
	 */
	public void connect() throws Exception {
		zk = new ZooKeeper(hosts, SESSION_TIMEOUT, new ConnWatcher());
		// 等待连接完成
		connectedSignal.await();
	}

	public class ConnWatcher implements Watcher {
		public void process(WatchedEvent event) {
			// 连接建立, 回调process接口时, 其event.getState()为KeeperState.SyncConnected
			if (event.getState() == KeeperState.SyncConnected) {
				// 放开闸门, wait在connect方法上的线程将被唤醒
				connectedSignal.countDown();
			}
		}
	}
}
</pre>







links
-----
+ [目录](../zookeeper)
+ 上一节: [Zookeeper  安装和配置](Zookeeper--安装和配置.md)
