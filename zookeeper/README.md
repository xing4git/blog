Index
-----

####ZooKeeper  数据模型
*2013-05-17*
ZooKeeper的数据结构, 与普通的文件系统极为类似. 见下图:...[Read More](zookeeper/ZooKeeper--数据模型.md)

####Zookeeper  安装和配置
*2013-05-17*
Zookeeper的安装和配置十分简单, 既可以配置成单机模式, 也可以配置成集群模式.下面将分别进行介绍....[Read More](zookeeper/Zookeeper--安装和配置.md)

####ZooKeeper Java API
*2013-05-18*
ZooKeeper提供了Java和C的binding. 本文关注Java相关的API....[Read More](zookeeper/ZooKeeper Java API.md)

####ZooKeeper示例 分布式锁
*2013-05-18*
### 场景描述在分布式应用, 往往存在多个进程提供同一服务. 这些进程有可能在相同的机器上, 也有可能分布在不同的机器上. 如果这些进程共享了一些资源, 可能就需要分布式锁来锁定对这些资源的访问.  ...[Read More](zookeeper/ZooKeeper示例 分布式锁.md)

####ZooKeeper示例 实时更新server列表
*2013-05-18*
通过之前的3篇博文, 讲述了ZooKeeper的基础知识点. 可以看出, ZooKeeper提供的核心功能是非常简单, 且易于学习的. 可能会给人留下ZooKeeper并不强大的印象, 事实并非如此, 基于ZooKeeper的核心功能, 我们可以扩展出很多非常有意思的应用. 接下来的几篇博文, 将陆续介绍ZooKeeper的典型应用场景....[Read More](zookeeper/ZooKeeper示例 实时更新server列表.md)

