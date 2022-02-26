## WebSocket
[http://coolaf.com/tool/chattest](http://coolaf.com/tool/chattest)
> 到上面的網站，然後輸入下面這串，可以即時取得最新資料
```bash
ws://127.0.0.1:8080/api/v1/trade/sub 
```
## API DOC
> 本機部署完後，HOST自行更改
[http://127.0.0.1:8080/swagger/index.html](http://127.0.0.1:8080/swagger/index.html)
## 部署方式
```bash
docker-compose up -d
# restart service
docker-compose restart
```
## 資料持久化設定
```yaml
# 打開config.yaml
# SAVE_PG true = 要做持久化儲存
SAVE_PG: true
# SAVE_PG_MAX 100 = 每一次要儲存的資料量
# 為避免造成DB IO的Loading過重，所以改成批次儲存
SAVE_PG_MAX: 100
```