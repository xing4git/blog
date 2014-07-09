### LZO安装和配置

我们的日志存储在hdfs上, 考虑到hdfs磁盘空间不足, 决定引入LZO压缩. 由于LZO的GPL开源协议和Hadoop采用的开源协议不兼容, 所以只能自己安装LZO. 我们的环境如下:

```
OS: CentOS release 6.5
Hadoop: 2.3.0-cdh5.0.1
```

* 在hadoop集群的机器上增加yum源(/etc/yum.repos.d/文件夹下任意一个文件):

```
[cloudera-gplextras5]
name=Cloudera's GPLExtras, Version 5
baseurl=http://archive.cloudera.com/gplextras5/redhat/5/x86_64/gplextras/5/
gpgkey = http://archive.cloudera.com/gplextras5/redhat/5/x86_64/gplextras/RPM-GPG-KEY-cloudera 
gpgcheck = 1

```

* 在hadoop集群上安装lzo:

```
yum install hadoop-lzo
```

* 修改core-site.xml:

```
<property>
    <name>io.compression.codecs</name>
	<value>com.hadoop.compression.lzo.LzoCodec,com.hadoop.compression.lzo.LzopCodec</value>
</property>
```
在value中增加`com.hadoop.compression.lzo.LzoCodec,com.hadoop.compression.lzo.LzopCodec`, 原有的codecs依旧保留.


* 重启hadoop集群

* git clone https://github.com/twitter/hadoop-lzo, 修改pom文件, 更改私服地址. 执行`mvn deploy`命令, 将jar包update到私服. 之所以要upload, 是因为外部可下载的jar包不包含native代码.


* 用户工程依赖私服上的lzo jar包. 代码如下:

```
Configuration conf = new Configuration();
conf.set("io.compression.codecs", "com.hadoop.compression.lzo.LzoCodec,com.hadoop.compression.lzo.LzopCodec");
CompressionCodecFactory ccf = new CompressionCodecFactory(conf);
CompressionCodec codec = ccf.getCodecByName(LzopCodec.class.getName());
// codec.createInputStream()
// codec.createOutputStream()
```
使用此方法可以避免native依赖找不到的问题, 因为LZO的native依赖已经在上一步打成的包里了. 外部可下载的hadoop-lzo都是不包含native依赖的, 需要在程序运行的机器上安装了LZO, 并指定hadoop lzo的native依赖路径. 这是完全无法接受的, 因为hdfs相关的程序会部署在很多台机器上, 而且迁移变更的几率比较高, 要求运行hdfs相关程序的机器都安装lzo是不可能的.