# VIRSH QUICK REFERENCE

- start vm

```bash
sudo virsh start --domain <vm_name>
```

- autostart vm

```bash
sudo virsh autostart --domain <vm_name>
```

- disable autostart vm

```bash
sudo virsh autostart --disable <vm_name>
```

- suspend vm

```bash
sudo virsh suspend --domain <vm_name>
```

- resume suspended vm

```bash
sudo virsh resume --domain <vm_name>
```
