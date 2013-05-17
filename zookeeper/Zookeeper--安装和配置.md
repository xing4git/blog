Zookeeper  安装和配置
----

Zookeeper的安装和配置十分简单, 既可以配置成单机模式, 也可以配置成集群模式.
下面将分别进行介绍.

### 单机模式
[点击这里](http://zookeeper.apache.org/releases.html)下载zookeeper的安装包之后, 解压到合适目录.
进入zookeeper目录下的conf子目录, 创建zoo.cfg:
```
tickTime=2000  
dataDir=/Users/apple/zookeeper/data  
dataLogDir=/Users/apple/zookeeper/logs  
clientPort=4180  
```
参数说明:
+ tickTime: zookeeper中使用的基本时间单位, 毫秒值.  
+ dataDir: 数据目录. 可以是任意目录.  
+ dataLogDir: log目录, 同样可以是任意目录. 如果没有设置该参数, 将使用和dataDir相同的设置.  
+ clientPort: 监听client连接的端口号.  
至此, zookeeper的单机模式已经配置好了. 启动server只需运行脚本:
```bash
bin/zkServer.sh start
```
Server启动之后, 就可以启动client连接server了, 执行脚本:
```bash
bin/zkCli.sh -server localhost:4180
```

### 伪集群模式
所谓伪集群, 是指在单台机器中启动多个zookeeper进程, 并组成一个集群. 以启动3个zookeeper进程为例.











links
-----
+ [目录](../golang)
