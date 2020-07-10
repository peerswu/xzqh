# 功能
* 县级已上行政区划转成SQL INSERT语句
* 添加拼音
* 添加声母
* 添加与行政片区的关系
* 添加行政区层级关系

# 安装

```bash
go get -u github.com/peerswu/xzqh
```

# 如何使用
## 下载最新的行政区划表

* 行政区划可以从 中华人民共和国民政部 下载，搜索关键字：中华人民共和国民政部 行政区划
* 下载后，编辑下载内容，只保留 行政区划部分，同时去除多余空白字符

如: [2020年2月中华人民共和国县以上行政区划代码](http://www.mca.gov.cn/article/sj/xzqh/2020/202003/20200300026218.shtml)，处理后保存为xzqh.txt，内容如下

```
110000 北京市 
110101 东城区 
....
820000 澳门特别行政区
```

## 执行命令，生成结果
```bash
$GOPATH/bin/xzqh xzqh.txt
```
```sql
INSERT INTO fog_addr_area(id, name) VALUES (1, '华北'),(2, '东北'),(3, '华东'),(4, '西北'),(5, '华中'),(6, '华南'),(7, '西南'),(8, '其他');
INSERT INTO fog_addr_region(id, name, parent_id, province_id, city_id, area_id, pinyin, title_pinyin, initials, first_letter, is_leaf) VALUES 
(110000,'北京市',0,110000,0,1,'BeiJingShi','beijingshi','bjsh','bjs',0),
(110101,'东城区',110000,110000,110101,1,'DongChengQu','dongchengqu','dchq','dcq',1),
...
(710000,'台湾省',0,710000,0,8,'TaiWanSheng','taiwansheng','tsh','tws',1),
(810000,'香港特别行政区',0,810000,0,8,'XiangGangTeBieHangZhengQu','xianggangtebiehangzhengqu','xgtbhzhq','xgtbhzq',1),
(820000,'澳门特别行政区',0,820000,0,8,'AoMenTeBieHangZhengQu','aomentebiehangzhengqu','mtbhzhq','Ãmtbhzq',1);
```
