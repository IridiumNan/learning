# OS SIGNAL

> [!NOTE]
> You can find all sh script on `script/`

## SIGINT

We use the ^C to send the signal interrupt to the process.

Then it will stop for default.

However, we can define how this process react with the sigint (signal interrupt).

If we add these code on a bash script then run it.

You will find that this process will print `sigint received - ignoring lol` instead of interrupt.

```sh
on_sigint() {
    echo "sigint received - ignoring lol"
}

trap on_sigint SIGINT
```

> [!NOTE]
> trap action signal  
> means exec specific command when process receive this signal

You can run the `sigint.sh` directly to test this.

---

## SIGTERM

When `kill` command is called, it send `SIGTERM` (signal terminate) to this process.

That will make this process stop.

So we can use the kill command to terminate the `sigint.sh` script process

We put another function on `sigint-and-sigterm.sh` to handle the sigterm from `kill` command

```sh
on_sigterm() {
    echo "sigterm received - ignore lol"
}

trap on_sigterm SIGTERM
```

And we will find that the `kill id` can not terminate it properly.

Now we can use `kill -9 id`.

> [!NOTE]
> The Docker and Kubernetes use SIGTERM to stop a container

---

## SIGTSTP

When we use ^Z to freeze a process, it send the `SIGTSTP`.

We can also add another `trap` command line to handle this signal.

So that the process will run a function or exec command we want as it receive the `SIGTSTP`

See `all-lol.sh` or you write it manually

---

## PS

> [!NOTE]
> When all normal interrupt signal has been handled by a process  
> You can always use the kill -9 pid to kill it
> for getting pid, use the `ps aux | grep process_name`

```bash
❯ sleep 30 &
[1] 2350741
❯ ps aux | grep sleep
cai      2350741  0.0  0.0 230560  2436 pts/81   SN   16:10   0:00 sleep 30
```

so the 2350741 is its pid. We can kill this process by `kill -9 2350741`

---

For golang signal handling, see `go/`
