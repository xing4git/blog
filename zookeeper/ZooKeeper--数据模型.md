ZooKeeper  数据模型
----

ZooKeeper的数据结构, 与普通的文件系统极为类似. 见下图:

![ZooKeeper数据结构](./model.jpg)

图片引用自[developerworks](http://www.ibm.com/developerworks/cn/opensource/os-cn-zookeeper/)

图中的每个节点称为一个znode. 每个znode由3部分组成:
- stat. 此为状态信息, 描述该znode的版本, 权限等信息.
- data. 与该znode关联的数据.
- children. 该znode下的子节点.

### ZooKeeper命令
在深入znode的各个部分之前, 首先需要熟悉一些常用的ZooKeeper命令.
- zkCli.sh 连接server: `bin/zkCli.sh -server 10.1.39.43:4180`
- ls: 列出指定node的子node

```
[zk: 10.1.39.43:4180(CONNECTED) 9] ls /
[hello, filesync, zookeeper, xing, server, group, log]
[zk: 10.1.39.43:4180(CONNECTED) 10] ls /hello
[]
```
- create: 创建znode节点, 并关联数据. `create /hello world`, 创建节点/hello, 并将字符串"world"关联到该节点中.
- get: 获取znode的数据和状态信息.

```
[zk: 10.1.39.43:4180(CONNECTED) 7] get /hello
world
cZxid = 0x10000042c
ctime = Fri May 17 17:57:33 CST 2013
mZxid = 0x10000042c
mtime = Fri May 17 17:57:33 CST 2013
pZxid = 0x10000042c
cversion = 0
dataVersion = 0
aclVersion = 0
ephemeralOwner = 0x0
dataLength = 5
numChildren = 0
```















links
-----
+ [目录](../zookeeper)
+ 下一节: [Zookeeper  安装和配置](Zookeeper--安装和配置.md)
