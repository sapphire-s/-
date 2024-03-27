## 统一认证开发中

### 配置文件

创建一个config文件夹，在里面添加config.json文件，格式如下：

```json
{
  "database": {
    "type": "mysql",
    "host": ":3306",
    "name": "",
    "user": "",
    "passwd": "",
    "charset": "utf8mb4",
    "parsetime": "True"
  },
  "server": {
    "addr": ":11111"
  },
  "jwt": {
    "key": ""
  },
  "feishu": {
    "app_id": "",
    "app_secret": "",
    "base_url": "",
    "redirect_uri":"",
    "enableTokenCache": true
  }
}
```

## TODO: Try casbin for perm management

