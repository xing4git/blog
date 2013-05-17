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

将zookeeper的目录拷贝2份:
```
|--zookeeper0
|--zookeeper1
|--zookeeper2
```

更改zookeeper0/conf/zoo.cfg文件为:
```
tickTime=2000  
initLimit=5  
syncLimit=2  
dataDir=/Users/apple/zookeeper0/data  
dataLogDir=/Users/apple/zookeeper0/logs  
clientPort=4180
server.0=127.0.0.1:8880:7770  
server.1=127.0.0.1:8881:7771  
server.2=127.0.0.1:8882:7772
```
新增了几个参数, 其含义如下:
+ initLimit: zookeeper集群中的包含多台server, 其中一台为leader, 集群中其余的server为follower. initLimit参数配置初始化连接时, follower和leader之间的最长心跳时间. 此时该参数设置为5, 说明时间限制为5倍tickTime, 即5*2000=10000ms=10s.
+ syncLimit: 该参数配置leader和follower之间发送消息, 请求和应答的最大时间长度. 此时该参数设置为2, 说明时间限制为2倍tickTime, 即4000ms.
+ server.X=A:B:C 其中X是一个数字, 表示这是第几号server. A是该server所在的IP地址. B配置该server和集群中的leader交换消息所使用的端口. C配置选举leader时所使用的端口. 由于配置的是伪集群模式, 所以各个server的B, C参数必须不同.












links
-----
+ [目录](../golang)
