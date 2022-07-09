# 预警推送
简单的预警推送，可以用来推送一些简单的预警信息。

## 使用方式
### 直接使用
```bash
cd cmd
go bulid
```

## TODO
- [X] 调整启动方式，可通过包引入

### 缓存
- [X] 内存缓存
- [ ] mysql缓存
- [ ] redis缓存

### 推送
- [X] [pushDeer](https://github.com/easychen/pushdeer)推送
- [ ] server酱推送
- [ ] 发送邮件

### 任务
- [ ] 天气预报
- [ ] 气候预警
  - [x] 按照省份推送
  - [x] 按照省市区过滤
  - [ ] 增加每个地区的推送key列表，以达到每个地区发送不同接受者 
- [ ] 温差过大预警
  - [ ] 按照省份推送
  - [ ] 按照省市区过滤

## 配置文件样例
```yaml
log:
  path: './log/alertService.log'
  maxAge: '720h'
  rotationAge: '24h'
  level: 'info'
jobs:
  - type: "WeatherWarning"
    # 使用了 github.com/robfig/cron/v3 的 cron 定时器
    # schedule规则详情，可查看：https://godoc.org/github.com/robfig/cron/v3#h3-Schedule_Expressions
    schedule: "@every 10m"
    title: "气象预警"
      # 此处不填写的话，默认为上海市
      # 此处必须填写完整的省市县名称，不能使用简称，如果使用简称，那么会导致查询不到数据
      # 例如：上海市，不能使用“上海”，必须填写“上海市”
      #      西城区，不能使用“西城”，必须填写“西城区”
      #      正定县，不能使用“正定”，必须填写“正定县”
    location:
      - province: '上海市'
      - province: '山西省'
        city: '太原市'
        county: '尖草坪区'
      - province: '江苏省'
        city: '连云港市'
        county: '灌云县'

store:
  # 若使用内存缓存，重启会重新拉取数据
  type: "memory"

alert:
  #  - type: "email"
  #    to: " - "
  #    from: " subject: "
  #    body: ""
  - type: "pushDeer"
    keys:
      - "PDU1430T9LLQmSCITOG5kQQLUTDhfK25ac0rwzzq"
```