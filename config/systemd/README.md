## Systemd

### Install 
```
cp wuproxy.service /etc/systemd/system
systemctl daemon-reload
systemctl enable --now wuproxy
```

### View logs
```
journalctl --follow --unit wuproxy.service
```

### Custom runtime values
To pass custom options to the server you need to add them to /etc/defaults/wuproxy, or /etc/sysconfig/wuproxy (depending on which your OS uses.)  

Example:
`/etc/default/wuproxy`
```
OPTIONS="--forword --publish mqtt://localhost:1883 --topic my/weather/topic"
```
