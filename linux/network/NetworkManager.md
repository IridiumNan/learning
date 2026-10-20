# NetworkManager

这是主流Linux发行版采用的网络服务，特点是资源占用稍高，功能强大

命令行对应 -> nmcli
这里记录常用的命令

## 常用命令

- 查看现有的连接

```bash
nmcli 
```

- 列出周围的wifi

```bash
nmcli device wifi list
```

- 连接wifi

```bash
nmcli device wifi connect <wifi 名称> password <wifi 密码>
```

---

## 配置静态ip地址

```bash
nmcli con mod <连接名> ipv4.method manual ipv4.address <想要配置的ip地址>/<子网掩码 配置， 如 24, 28>


# 设置网关
nmcli con mod <连接名> ipv4.gateway <gateway的ip地址>

# 配置dns服务器
nmcli con mod <连接名> ipv4.dns <dns的ip地址>
```

- 实例

```bash
nmcli con mod "iPhone" ipv4.method manual ipv4.address 192.20.10.5/28

nmcli con mod "iPhone" ipv4.gateway 172.20.10.1

nmcli con mod "iPhone" ipv4.dns 223.5.5.5
```
