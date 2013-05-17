ZooKeeper  数据模型
----

ZooKeeper的数据结构, 与普通的文件系统极为类似. 见下图:

![ZooKeeper数据结构](./model.jpg)

图片引用自[developerworks](http://www.ibm.com/developerworks/cn/opensource/os-cn-zookeeper/)

图中的每个节点称为一个znode. 每个znode由3部分组成:
- stat. 此为状态信息, 描述该znode的版本, 权限等信息.
- data. 与该znode关联的数据.
- children. 该znode下的子节点.


links
-----
+ [目录](../zookeeper)
+ 下一节: [Zookeeper  安装和配置](Zookeeper--安装和配置.md)
