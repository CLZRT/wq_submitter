## PreRequest
- Cloud Server
- Installed Mysql And Version >=5.7 
- Go env

## Config
- replace your-config in configs/config-template.yaml and change his name to config.yaml
## Run
```shell

mysql -u {your-mysql-name} -p < sql/main.sql
make install
```



# API文档

## Hello

### Ping

> 检查服务状态

```Yaml
Get
  {Your-Host-Name:Your-Port}/wq_submitter/hello
```



## Alpha

### Upload

> 上传AlphaList

```Yaml
Upload: Post- {Your-Host-Name:Your-Port}/wq_submitter/alpha/upload
```

```json
{
    "idea": {
        "ideaAlphaTemplate": "this is template",
        "ideaTitle": "hello",
        "ideaDesc": "hello world",
        "concurrencyNum": 2
    },
    "alphaList": [
        {
            "type": "REGULAR",
            "regular": "ts_scale(liabilities,3)/ts_scale(assets,3)",
            "settings": {
                "instrumentType": "EQUITY",
                "region": "USA",
                "universe": "TOP500",
                "delay": 1,
                "decay": 0,
                "neutralization": "MARKET",
                "truncation": 0.01,
                "pasteurization": "ON",
                "nanHandling": "OFF",
                "language": "FASTEXPR",
                "unitHandling": "VERIFY",
                "visualization": false
            }
        },
        {
            "type": "REGULAR",
            "regular": "-ts_scale(liabilities,3)/ts_scale(assets,3)",
            "settings": {
                "instrumentType": "EQUITY",
                "region": "USA",
                "universe": "TOP1000",
                "delay": 1,
                "decay": 0,
                "neutralization": "MARKET",
                "testPeriod": "P1Y2M",
                "truncation": 0.08,
                "pasteurization": "ON",
                "nanHandling": "ON",
                "unitHandling": "VERIFY",
                "language": "FASTEXPR",
                "visualization": false
            }
        }
    ]
}

```

### List

> 获取AlphaList For Idea

```yaml
Get
{Your-Host-Name:Your-Port}/wq_submitter/alpha/list?ideaId=3
```



## Idea

### GetALL

```Yaml
Get
{Your-Host-Name:Your-Port}/wq_submitter/idea/all
```

### GetRun

```Yaml
Get
{Your-Host-Name:Your-Port}/wq_submitter/idea/run
```

### GetUnFinish

```Yaml
Get
{Your-Host-Name:Your-Port}/wq_submitter/idea/unfinish
```

### Delete-Post

> 删除Idea及其对应的Alpha

```Yaml
Post
{Your-Host-Name:Your-Port}/wq_submitter/idea/delete
```

```json
{
    "id": 1
}
```

### Concurrency

> 改变Idea的并发度

```Yaml
Post
{Your-Host-Name:Your-Port}/wq_submitter/idea/concurrenct
```

```json
{
    "id": 1,
    "num": 3
}
```

