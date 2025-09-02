## PreRequest
- Cloud Server
- Installed Mysql And Version >=5.7 

## Config
- replace your-config in configs/config-template.yaml and change his name to config.yaml
## Run
```shell

mysql -u {your-mysql-name} -p < sql/main.sql
make install
```

