ZooKeeper示例 分布式锁
----

### 场景描述
在分布式应用, 往往存在多个进程提供同一服务. 这些进程有可能在相同的机器上, 也有可能分布在不同的机器上. 如果这些进程共享了一些资源, 可能就需要分布式锁来锁定对这些资源的访问.  
本文将介绍如何利用zookeeper实现分布式锁.

### 思路
进程需要访问共享数据时, 就在"/locks"节点下创建一个sequence类型的子节点, 称为thisPath. 当thisPath在所有子节点中最小时, 说明该进程获得了锁. 进程获得锁之后, 就可以访问共享资源了. 访问完成后, 需要将thisPath删除. 锁由新的最小的子节点获得.  
有了清晰的思路之后, 还需要补充一些细节. 进程如何知道thisPath是所有子节点中最小的呢? 可以在创建的时候, 通过getChildren方法获取子节点列表, 然后在列表中找到排名比thisPath前1位的节点, 称为waitPath, 然后在waitPath上注册监听, 当waitPath被删除后, 进程获得通知, 此时说明该进程获得了锁.

### 实现
以一个DistributedClient对象模拟一个进程的形式, 演示zookeeper分布式锁的实现.

```java
public class DistributedClient {
	// 超时时间
	private static final int SESSION_TIMEOUT = 5000;
	// zookeeper server列表
	private String hosts = "localhost:4180,localhost:4181,localhost:4182";
	private String groupNode = "locks";
	private String subNode = "sub";

	private ZooKeeper zk;
	// 当前client创建的子节点
	private String thisPath;
	// 当前client等待的子节点
	private String waitPath;

	private CountDownLatch latch = new CountDownLatch(1);

	/**
	 * 连接zookeeper
	 */
	public void connectZookeeper() throws Exception {
		zk = new ZooKeeper(hosts, SESSION_TIMEOUT, new Watcher() {
			public void process(WatchedEvent event) {
				try {
					// 连接建立时, 打开latch, 唤醒wait在该latch上的线程
					if (event.getState() == KeeperState.SyncConnected) {
						latch.countDown();
					}

					// 发生了waitPath的删除事件
					if (event.getType() == EventType.NodeDeleted && event.getPath().equals(waitPath)) {
						doSomething();
					}
				} catch (Exception e) {
					e.printStackTrace();
				}
			}
		});

		// 等待连接建立
		latch.await();

		// 创建子节点
		thisPath = zk.create("/" + groupNode + "/" + subNode, null, Ids.OPEN_ACL_UNSAFE,
				CreateMode.EPHEMERAL_SEQUENTIAL);
		
		// wait一小会, 让结果更清晰一些
		Thread.sleep(10);

		// 注意, 没有必要监听"/locks"的子节点的变化情况
		List<String> childrenNodes = zk.getChildren("/" + groupNode, false);

		// 列表中只有一个子节点, 那肯定就是thisPath, 说明client获得锁
		if (childrenNodes.size() == 1) {
			doSomething();
		} else {
			String thisNode = thisPath.substring(("/" + groupNode + "/").length());
			// 排序
			Collections.sort(childrenNodes);
			int index = childrenNodes.indexOf(thisNode);
			if (index == -1) {
				// never happened
			} else if (index == 0) {
				// inddx == 0, 说明thisNode在列表中最小, 当前client获得锁
				doSomething();
			} else {
				// 获得排名比thisPath前1位的节点
				this.waitPath = "/" + groupNode + "/" + childrenNodes.get(index - 1);
				// 在waitPath上注册监听器, 当waitPath被删除时, zookeeper会回调监听器的process方法
				zk.getData(waitPath, true, new Stat());
			}
		}
	}

	private void doSomething() throws Exception {
		try {
			System.out.println("gain lock: " + thisPath);
			Thread.sleep(2000);
			// do something
		} finally {
			System.out.println("finished: " + thisPath);
			// 将thisPath删除, 监听thisPath的client将获得通知
			// 相当于释放锁
			zk.delete(this.thisPath, -1);
		}
	}

	public static void main(String[] args) throws Exception {
		for (int i = 0; i < 10; i++) {
			new Thread() {
				public void run() {
					try {
						DistributedClient dl = new DistributedClient();
						dl.connectZookeeper();
					} catch (Exception e) {
						e.printStackTrace();
					}
				}
			}.start();
		}
		
		Thread.sleep(Long.MAX_VALUE);
	}
}
```


links
-----
+ [目录](../zookeeper)
+ 上一节: [ZooKeeper Java API](ZooKeeper Java API.md)
+ 下一节: [ZooKeeper示例 实时更新server列表](ZooKeeper示例 实时更新server列表.md)
