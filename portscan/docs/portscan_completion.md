## portscan completion

Generate bash completion for your command

### Synopsis

To load your completion run
source <(portscan completion)

To load completions automatically on login, add this line to your .bashrc file:
$ ~/.bashrc
source <(portscan completion)

```
portscan completion [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string       config file (default is $HOME/.portscan.yaml)
  -f, --hosts-file string   portScan hosts file (default "portScan.hosts")
```

### SEE ALSO

* [portscan](portscan.md)	 - Fast TCP port scanner

###### Auto generated by spf13/cobra on 9-Jan-2023
