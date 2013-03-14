# What is Par?
Have you ever want to run a bunch of related processes at once, and maybe even collect their output? If so, this little thing is for you. Par spins up a fleet of goroutines and runs a command template you specify against standard input. It's best suited for tasks where you would expect concurrency to help you out -- things where network timeouts are possible and can slow your processing time down to a crawl.

# Installation

## On OS X:

```shell
$ brew install go
$ echo PATH=$PATH:`brew --prefix go`/bin >> ~/.bash_profile && source ~/.bash_profile
$ go get github.com/jroes/par
```

# Examples

## Webscaling your family

```shell
$ cat family.txt
Jon
Mai
Kelly
Tommy
$ par "echo hello %line" < family.txt
Running with 8 workers
[echo hello Mai] starting
[echo hello Mai] hello Mai
[echo hello Mai] finished
[echo hello Kelly] starting
[echo hello Kelly] hello Kelly
[echo hello Kelly] finished
[echo hello Jon] starting
[echo hello Jon] hello Jon
[echo hello Jon] finished
[echo hello Tommy] starting
[echo hello Tommy] hello Tommy
[echo hello Tommy] finished
```

## Find the IP address of every dictionary word .com domain

```shell
$ par "dig +short %line.com" < /usr/share/dict/words
Running with 8 workers
[dig +short aardvark.com] starting
[dig +short aardvark.com] 67.210.123.5
[dig +short aardvark.com] finished
[dig +short Aaronic.com] starting
[dig +short Aaronic.com] 67.15.74.21
[dig +short Aaronic.com] finished
[dig +short Aaronite.com] starting
[dig +short Aaronite.com] 82.165.125.59
[dig +short Aaronite.com] finished
```

## Mirroring a backup to a bunch of geographically dispersed servers

```shell
$ par "scp backup.dump jroes@%line:backups/" < ~/.backup_servers
Running with 8 workers
[scp backup.dump jroes@jroes.net:backups/] starting
[scp backup.dump jroes@jroes.net:backups/] finished
[scp backup.dump jroes@backup.myco.com:backups/] starting
[scp backup.dump jroes@backup.myco.com:backups/] finished
```
