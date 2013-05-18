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

### 创建znode
ZooKeeper对象的create方法用于创建znode.
<pre name="code" class="java">
String create(String path, byte[] data, List acl, CreateMode createMode);
</pre>
以下为各个参数的详细说明:
+ path. znode的路径.
+ data. 与znode关联的数据.
+ acl. 指定权限信息, 如果不想指定权限, 可以传入Ids.OPEN_ACL_UNSAFE.
+ 指定znode类型. CreateMode是一个枚举类, 从中选择一个成员传入即可. 关于znode类型的详细说明, 可参考本人的上一篇博文.

```java
/**
 * 创建临时节点
 */
public void create(String nodePath, byte[] data) throws Exception {
	zk.create(nodePath, data, Ids.OPEN_ACL_UNSAFE, CreateMode.EPHEMERAL);
}
```

### 获取子node列表
ZooKeeper对象的getChildren方法用于获取子node列表.
<pre name="code" class="java">
List<String> getChildren(String path, boolean watch);
</pre>
watch参数用于指定是否监听path node的子node的增加和删除事件, 以及path node本身的删除事件.

### 判断znode是否存在
ZooKeeper对象的exists方法用于判断指定znode是否存在.
<pre name="code" class="java">
Stat exists(String path, boolean watch);
</pre>
watch参数用于指定是否监听path node的创建, 删除事件, 以及数据更新事件. 如果该node存在, 则返回该node的状态信息, 否则返回null.

### 获取node中关联的数据
ZooKeeper对象的getData方法用于获取node关联的数据.
<pre name="code" class="java">
byte[] getData(String path, boolean watch, Stat stat);
</pre>
watch参数用于指定是否监听path node的删除事件, 以及数据更新事件, 注意, 不监听path node的创建事件, 因为如果path node不存在, 该方法将抛出KeeperException.NoNodeException异常.  
stat参数是个传出参数, getData方法会将path node的状态信息设置到该参数中.  

### 更新node中关联的数据
ZooKeeper对象的setData方法用于更新node关联的数据.
<pre name="code" class="java">
Stat setData(final String path, byte data[], int version);
</pre>
data为待更新的数据.  
version参数指定要更新的数据的版本, 如果version和真实的版本不同, 更新操作将失败. 指定version为-1则忽略版本检查.  
返回path node的状态信息.  

### 删除znode
ZooKeeper对象的delete方法用于删除znode.
<pre name="code" class="java">
void delete(final String path, int version);
</pre>
version参数的作用同setData方法.

### 其余接口
请查看ZooKeeper对象的API文档.

### 需要注意的几个地方
+ znode中关联的数据不能超过1M. zookeeper的使命是分布式协作, 而不是数据存储.
+ getChildren, getData, exists方法可指定是否监听相应的事件. 而create, delete, setData方法则会触发相应的事件的发生.
+ 以上介绍的几个方法大多存在其异步的重载方法, 具体请查看API说明.



links
-----
+ [目录](../zookeeper)
+ 上一节: [Zookeeper  安装和配置](Zookeeper--安装和配置.md)
+ 下一节: [ZooKeeper示例 实时更新server列表](ZooKeeper示例 实时更新server列表.md)
