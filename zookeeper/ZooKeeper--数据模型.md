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

连接server  
```
bin/zkCli.sh -server 10.1.39.43:4180
```

列出指定node的子node  
```
[zk: 10.1.39.43:4180(CONNECTED) 9] ls /
[hello, filesync, zookeeper, xing, server, group, log]
[zk: 10.1.39.43:4180(CONNECTED) 10] ls /hello
[]
```

创建znode节点, 并指定关联数据  
```
create /hello world
```
创建节点/hello, 并将字符串"world"关联到该节点中.

获取znode的数据和状态信息  
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

### znode节点的状态信息
使用get命令获取指定节点的数据时, 同时也将返回该节点的状态信息, 称为Stat. 其包含如下字段:
+ czxid. 节点创建时的zxid.
+ mzxid. 节点最新一次更新发生时的zxid.
+ ctime. 节点创建时的时间戳.
+ mtime. 节点最新一次更新发生时的时间戳.
+ dataVersion. 节点数据的更新次数.
+ cversion. 其子节点的更新次数.
+ aclVersion. 节点ACL(授权信息)的更新次数.
+ ephemeralOwner. 如果该节点为ephemeral节点, ephemeralOwner值表示与该节点绑定的session id. 如果该节点不是ephemeral节点, ephemeralOwner值为0. 至于什么是ephemeral节点, 请看后面的讲述.
+ dataLength. 节点数据的字节数.
+ numChildren. 子节点个数.

### zxid
znode节点的状态信息中包含czxid和mzxid, 那么什么是zxid呢?  
ZooKeeper状态的每一次改变, 都对应着一个递增的`Transaction id`, 该id称为zxid. 由于zxid的递增性质, 如果zxid1小于zxid2, 那么zxid1肯定先于zxid2发生. 创建任意节点, 或者更新任意节点的数据, 或者删除任意节点, 都会导致Zookeeper状态发生改变, 从而导致zxid的值增加.

### 节点类型
讲述节点状态的ephemeralOwner字段时, 提到过有的节点是ephemeral节点, 而有的并不是. 那么节点都具有哪些类型呢? 每种类型的节点又具有哪些特点呢?
+ 















links
-----
+ [目录](../zookeeper)
+ 下一节: [Zookeeper  安装和配置](Zookeeper--安装和配置.md)
