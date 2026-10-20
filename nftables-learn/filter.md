# Filter

## 前置条件

在这次实践中，我们将用到三台 `debian 13` 虚拟机

windows平台建议使用 [virtual box](https://www.virtualbox.org/) 或者 [vmware](https://www.vmware.com/) 创建虚拟机

linux 平台建议使用 [virt-manager](https://virt-manager.org/) 来创建虚拟机

debian 13 的iso镜像可以从 [debian 官网](https://www.debian.org/), [tuna 镜像站](https://mirrors.tuna.tsinghua.edu.cn/#), 或者 [中科大镜像站](https://mirrors.ustc.edu.cn/) 下载， 如果你没有魔法， 就使用镜像站， 如果你有魔法， 也可以使用官网

这里需要注意， 我们做网络实验不需要使用图形化界面，所以选择 **amd64 Live CD with Standard** 的版本, 避免网络安装以及桌面的资源占用

创建完三台虚拟机之后， 需要使得三台虚拟机处于同一个网络， 也就是可以让它们之间相互联系，收发流量包。

这里建议配置一下 ssh 以及 静态 ip地址， 可以参考这个文章 -> [waterman static ip](https://www.waterman.xin/posts/static_ip_addr/)

我们需要这样的三台机器

*Server* : 192.168.122.2/24
*Client1* : 192.168.122.3/24
*Client2* : 192.168.122.4/24

当然这里的ip地址不一定是这样， 因你的虚拟机配置而异， 但是要保证三个机器处于同一个网段

之后使用 `ssh` 连接虚拟机测试一下网络连通性

在 Server 上
请根据你实际设置的ip进行操作

查看当前机器的ip地址: `ip a` 命令

```bash
ping -c 2 192.168.122.3 && echo "success" || echo -e "\n\nfail to connect client1"

ping -c 2 192.168.122.4 && echo "success" || echo -e "\n\nfail to connect client2"
```

如果两个都是 `success` 的话就代表前置工作已经完成了

之后我们会涉及到特定端口的连通性测试
我们会使用 `nc` 这个工具

使用如下命令

```bash
nc -w1 -vz 192.168.122.3 22
```

w1 表示 wait one second 也就是等待一秒
v 表示verbose 也就是输出详细内容
z 表示只扫描但是不发送真实数据

这个命令表示 **试探一下 192.168.122.3 这个主机的 22 端口是否开放**
如果网段不同或者被防火墙拦截不会显示success
如果成功的话会出现像下面这样的输出

```bash
$ nc -w1 -vz 192.168.122.3 22
Connection to 192.168.122.3 22 port [tcp/ssh] succeeded!
```

---

## 开始实验

我们将会使用 原子化更新的方法来更新 nftables 的规则, 也就是 [basic](/basic.md#原子化重载) 当中提到的方法

接下来的更改我们都会在 `Server` 这个虚拟机上面完成

### 最小默认配置

- 创建`/tmp/nftables` 文件

```bash
# 清除现有的规则 
# 当然这个不是真的清除，文件尚未加载， 只是在操作成功的时候会清除原有内容
echo "flush ruleset" | tee /tmp/nftables

# 把现有规则写入配置文件
sudo nft -s list ruleset >> /tmp/nftables
```

- 打开并进行编辑

写入以下内容

```conf
flush ruleset

table ip filter {
    chain input {
        type filter hook input priority filter; policy accept;
    }
}
```

这是一个极简配置，一个 table 当中有一个 chain, 一个 chain 当中只有简单的定义而没有具体的规则

我们只需要看到 `type filter hook input priority filter; policy accept;` 这一行

其中 `hook` 就是流量包在运送的过程中我们可以做出干预的点，可以查看 [流程图](/README.md#netfilter-chart)

- `hook input` 代表 这里的 hook 是 input.
- `priority filter` 等价于 priority 0, 当其他链存在时， priority **更小的会先执行**
- `policy accept` 则表示如果进入这个 chain, 并且没有被下面的任何规则捕获(也就是都不满足 accept 或者 drop 的条件), 则使用默认策略对待这个流量包

所以现在这里的完整意思就是规定了一个 hook = input, priority = filter (0), policy = accept 的 chain. 面向所有端口的所有流量包都会被接收

- 加载配置

```bash
sudo nft -f /tmp/nftables
# 查看当前的规则
sudo nft list ruleset
# 你会发现出现的就是我们自己编辑的那个规则
```

加载之后你会发现无事发生(毕竟是接收所有的流量包)

我们在 `client1` 以及 `client2` 上测试一下连通性

```bash
nc -w1 -vz 192.168.122.2 22
```

会出现下面的结果

*Connection to 192.168.122.2 22 port [tcp/ssh] succeeded!*

### 添加拦截规则

我们现在打开刚才的那个配置文件， 修改为下面的内容

```conf
flush ruleset

table ip filter {
    chain input {
        type filter hook input priority filter; policy accept;
        tcp dport 22 counter accept
    }
}
```

我们添加了一行 `tcp dport 22 counter accept`

- `tcp dport 22` 代表 如果这个数据包是 tcp 协议并且 `dport` (目标端口) 为 22, 这是满足的条件
- `counter` 代表计数器， 我们在后面使用 `sudo nft list ruleset` 的时候会看到它可以统计符合这条规则的流量包
- `accept` 代表满足前面的规则所采取的措施， 这里就是接受这个流量包

- 重新加载配置文件

```bash
sudo nft -f /tmp/nftables
```

- 查看当前生效的配置

```bash
sudo nft list ruleset
```

会看到类似这样的输出
可以看到我们已经成功对 22 端口的流量包做了统计

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy accept;      
                tcp dport 22 counter packets 1780 bytes 156172 accept       
        }
}
```

- 更改配置

我们来到最刺激的环节，我们将修改规则丢弃除了 22 端口的所有流量包。于是这个虚拟机基本就断网了 **一定要在虚拟机操作**

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy accept;      
                tcp dport 22 counter packets 1780 bytes 156172 accept       
                counter drop
        }
}
```

我们加入了一个规则在末尾， 表示 计数并丢弃所有流量包(因为前面没有添加任何的条件)流量包

**如果流量包已经符合上面的规则并已经被 accept 或者 drop,则不会继续往下匹配**, 也就是说只有 目标端口是 22 端口的流量包可以幸免于难

- 重载配置

```bash
sudo nft -f /tmp/nftables
```

- 查看规则

```bash
sudo nft list ruleset
```

- 测试

我们在 `Server` 上进行测试

```bash
ping -c 2 8.8.8.8
```

我们会发现 100% package loss
就是我们自己把对方返回来的 流量包给拦截了

- 然后再次查看配置

```bash
sudo nft list ruleset
```

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy accept;      
                tcp dport 22 counter packets 421 bytes 36860 accept         
                counter packets 2 bytes 168 drop
        }
}
```

已经看到正好有两个包被丢弃了. **我们刚才 -c 2 就是代表发两个包**
于是对方会的两个包都被 input 这个 chain 的最后一条规则拦截

现在的配置写法等价于下面这种

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy drop;      
                tcp dport 22 counter accept         
                counter
        }
}
```

### 允许本地回环

需要注意的是， 刚才的规则把同一个机器内的不同端口之间的流量包也给拦截了
也就是说我们连 127.0.0.1 都ping不通了. 接下来我们要允许本地回环

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy drop;      
                iif lo counter accept
                tcp dport 22 counter accept         
                counter
        }
}
```

添加一行 `iif lo counter accept`

iif 是 `input interface`, 一般我们写 `ip a` 命令里面会出现的虚拟网卡的编号或者网卡名都行. "lo" 表示 本地回环.一般来说 lo 的编号是 1 , 所以写 1 也可以， nft会自动替换成真实的网卡名称

这条规则的作用是接收本地回环的流量包

之后我们重新加载

```bash
sudo nft -f /tmp/nftables
```

测试 127.0.0.1 连通性

```bash
ping -c 2 127.0.0.1
```

我们会发现成功 收到本地回环的流量包

### 基于连接状态筛选

我们希望已经建立的连接(比如说现在连的ssh) 以及相关的连接的流量包被接收，只需要修改配置为下面这样即可

- 修改配置

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy drop;      
                ct state established, related counter accept
                iif lo counter accept
                tcp dport 22 counter accept         
                counter
        }
}
```

`ct state` 代表 connection track state 就是根据连接的状态
established, related 则代表已经建立的和相关的连接
有了这个规则， 如果我们的机器向另外一个机器主动发包了， 我们就可以接收对方机器的包

- 重新加载

```bash
sudo nft -f /tmp/nftables
```

- 测试

我们的这个规则可以让我们主动跟外界建立连接
尝试 ping 8.8.8.8

```bash
ping -c 2 8.8.8.8
```

### 指定可以访问的saddr

为了安全性， 我们并不希望所有的外来流量包都可以访问 特殊端口比如说22
现在我们要设置一个规则，只允许 client1 以及宿主机访问 22 端口

```conf
table ip filter {
        chain input {
                type filter hook input priority filter; policy drop;      
                ct state established, related counter accept
                iif lo counter accept
                ip saddr { 192.168.122.1, 192.168.122.3 } tcp dport 22 counter accept         
                counter
        }
}
```

我们在前面添加了一个条件 `ip saddr { 192.168.122.1/24, 192.168.122.3/24 }` 前面已经讲过 saddr 表示源ip 地址， 然后这里的 大括号则表示集合的意思， 只要流量包的源ip地址是其中一个， 且访问的端口是 22, 我们就放行这个流量

- 接下来在 `client1` 和 `client2` 上面使用 `nc` 来测试对于 22 端口的连通性

client1

```bash
# 这里的 -s [saddr] 是指定 发送包时候使用的ip地址
# 因为我们可以还有尚未清理的 dhcp ip地址
cai@debian13-122-3:~$ nc -w1 -vz -s 192.168.122.3 192.168.122.2 22
Connection to 192.168.122.2 22 port [tcp/ssh] succeeded!
```

client2

```bash
cai@debian13-122-4:~$ nc -w1 -vz -s 192.168.122.4 192.168.122.2 22
nc: connect to 192.168.122.2 port 22 (tcp) timed out: Operation now in progress
```

可以看到 client1 顺利连接但是 client2 被拦截了， 这正是我们想要的效果

### 集合的使用方法

我们可以使用集合方便的管理
集合是定义在 table 当中的
接下来是一个案例

```conf
table ip filter {
        set allowed_ips {
            typeof ip saddr
            elements = { 192.168.122.1, 192.168.122.2, 192.168.122.3 }
        }
        chain input {
                type filter hook input priority filter; policy drop;      
                ct state established, related counter accept
                iif lo counter accept
                ip saddr @allowed_ips tcp dport 22 counter accept         
                counter
        }
}
```

注意使用引用集合的时候要加上 @

下面是一个特殊的集合案例
可以进行更加精细化的管理

```conf
table ip filter {
        set allowed_sockets {
            typeof ip saddr . tcp dport
            elements = { 192.168.122.1 . 22, 
                         192.168.122.2 . 22, 
                         192.168.122.3 . 22,
                         192.168.122.1 . 8080 }
        }
        chain input {
                type filter hook input priority filter; policy drop;      
                ct state established, related counter accept
                iif lo counter accept
                ip saddr . tcp dport @allowed_sockets tcp dport 22 counter accept         
                counter
        }
}
```
