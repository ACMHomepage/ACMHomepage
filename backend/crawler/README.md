# Crawler for CP
支持爬取各个 OJ 平台的数据到本系统。

## 为以下两个功能提供数据库支持
1. 前端展示
    + 用户主页：handle ac数量 rating
    + 提交记录：时间 题目链接 结果
2. 题目推荐(检索器)
    + 在固定的 problemset 中检索: rating 区间 | tag 列表 
    + 在自己/团队的 submission 中检索: submit_time 区间 | user 列表

## 数据库表设计
| 表名        	| 字段 	|              	|              	|              	|             	|                	|      	|
|-------------	|------	|--------------	|--------------	|--------------	|-------------	|----------------	|------	|
| oj_user     	| id   	| oj_id        	| handle       	| accept_count 	| max_rating  	| current_rating 	| link 	|
| problem     	| id   	| oj_id        	| problem_name 	| rating       	| link        	|                	|      	|
| submission  	| id   	| oj_user_id   	| problem_id   	| verdict_id   	| submit_time 	| link           	|      	|
| problem_tag 	| id   	| problem_id   	| tag_id       	|              	|             	|                	|      	|
| oj          	| id   	| oj_name      	|              	|              	|             	|                	|      	|
| tag         	| id   	| tag_name     	|              	|              	|             	|                	|      	|
| verdict     	| id   	| verdict_name 	|              	|              	|             	|                	|      	|

### 字段名解释
+ xxx_id: 指向 xxx 表的外键
+ link: 指向原 oj 平台的 url
+ verdict: 评测结果(如 AC, WA, TLE 等)
+ handle: 用户名

## Resource
### cftracker:
<https://github.com/mbashem/cftracker>

### cf_api:
<https://codeforces.com/apiHelp>
+ user.status: submission
+ user.info: rating
+ problemset.problems