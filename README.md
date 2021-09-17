# DozenPlans Server

> 学习养成计划主要是针对考研学生用户的任务管理系统， 用户可以添加~~学习~~任务，
> 每个任务都有自己的主题、 完成时间和优先级， 根据设定的完成时间和优先级对
> 用户进行不同频率的提醒， 直到其完成该任务， 任务可以设定提醒频率和时间期限，
> 避免导致任务越堆积越多， 除了用户自己可以设定任务以外， 本软件后台也会针对不
> 同科目来进行推送任务， 软件还提供任务分析功能， 针对用户每天的任务添加数和完
> 成率进行统计， 来让使用者更好的知道自己的复习进度和完成率。 通过可视化的图表
> 展示数据， 使用者能够有完成任务的成就感。综合实训项目

## 项目规划

### 需要实现的功能

1. 用户注册

2. 用户登录

3. 用户信息修改

4. 邮件进行提醒 

5. 邮件每天发个日报汇报进度

6. 找回密码 *

7. 添加任务

   - 设置主题
   - 设置优先级 (在前端，优先级会改变后面表单的预设项)
   - 设置ddl
   - 可以开启按间隔提醒  
   - 可以开启定时提醒 二选一

8. 删除任务

9. 方便演示 要提供直接触发发送日报的方法、发送提醒的方法

10. 进度可视化 前端使用echats 

    ![image-20210621171050542](https://img.aoyouer.com/images/2021/06/21/20210621172335.png)

    做过的显示成绿色，完成越多颜色越深。没做的显示成红色，鼠标移上去显示具体内容

    可以考虑加一个时间轴来显示做完的任务

11. tauri 打包成桌面应用 *

### 数据库表

1. 用户表
   - 用户名
   - 密钥
   - 邮箱
   - 身份 （用户 管理员
2. 任务表
   - 主题(类似标题)
   - 关联用户id
   - 细分阶段 待定
   - 优先级
   - 创建时间
   - 完成时间
   - 提醒频率
   - 当天还需要提醒的次数 或者 提醒间隔  或者是 提醒的时间   （从添加开始就每天提醒
   - 标签/类别 同一个类别的可以聚合提醒

### 技术细节

- 认证部分使用jwt jwt-go 中间件

- 密码使用bcrypt https://blog.csdn.net/m0_37609579/article/details/100785947  **bcrypt处理后的串拼接了盐，但是知道盐也没那么容易破解，因为bcrypt多次加密后，计算时间会边长，大大增加了开销**   **加盐慢哈希——函数非常慢**

- 前端使用vue + element

- 后端 golang

- 使用第三方的邮件服务器 smtp发信

- vscode中的插件rest client https://zhuanlan.zhihu.com/p/54266685

- 使用jwt + 中间件的形式进行认证 https://juejin.cn/post/6844903905424310279  https://zhuanlan.zhihu.com/p/70275218 （项目中使用的是HS256）

- post请求的格式说明 https://imququ.com/post/four-ways-to-post-data-in-http.html

- golang的时间解析 https://www.lagou.com/lgeduarticle/115634.html

- 跨域请求的处理 https://segmentfault.com/a/1190000022781975

- cron 定时任务

  

  > 这是正确的，任何人都可以使用哈希函数，然后输入Header和Payload来生成结果。但HS256签名不止这样，我们拿到Header、Payload外，还要加上一个密码，将这三个输入值一起哈希。输出结果是一个SHA-256 HMAC或者基于哈希的MAC。如果需要重复生成，则需要同时拥有Header、Payload和密码才可以。这也意味着，哈希函数的输出结果是一个数字签名，因为输出结果就表示了Payload是由拥有密码的角色生成并加签了的，没有其它方式可以生成这样的输出值了。
  >
  > *即使篡改了身份信息，并重新生成了签名，但是服务器验证的时候还会把 header + payload + **key*** 一起进行hash。 **验证方需要有这个密码**
  >
  > **HS256** 要求生成jwt方和验证jwt方都要用同样的密码，另外弱密码可能会被暴力破解
  >
  > **RS256则是使用了非对称加密，加密解密用不同的密钥**
  >
  > RS256使用一种特殊的密钥，叫RSA密钥。RSA是一种加解密算法，使用一个密钥进行加密，然后用另外一个密钥解密。值得注意的是，RSA不是哈希函数，从定义上来说，这种方式加密是可逆的，也就是我们可以从加密后的内容得到原始内容。
  >
  > **之所以不用RSA直接加密payload，是因为加密速度较慢，尤其是payload较大的时候，所以可以用sha256先做hash，获得摘要后，对摘要加密**
  >
  > 接收者将:
  >
  > 1. 取出Header和Payload，然后使用SHA-256进行哈希。
  > 2. 使用公钥解密数字签名，得到签名的哈希值。
  > 3. 接收者将解密签名得到的哈希值和刚使用Header和Payload参与计算的哈希值进行比较。如果两个哈希值相等，则证明JWT确实是由认证服务器创建的。
  >
  > 使用RS256，黑客可以轻松实现创建签名的第一步，即根据盗来的JWT Header和Payload生成SHA-256哈希值，之后他还要暴力破解RSA才能继续生成签名 **安全性更高**

- 也许会用swagger生成文档..(不过有点麻烦，也许就不做了)

### 使用技术

1. gin https://juejin.cn/post/6844903889272045575
2. gorm https://gorm.io/zh_CN/docs/models.html  **基于约定 **https://learnku.com/articles/40600



## 环境准备

创建数据库

```mysql
MariaDB [(none)]> CREATE USER 'dozenplans'@'%' IDENTIFIED BY 'dozenplans';
MariaDB [(none)]> CREATE DATABASE plansdb;
MariaDB [(none)]> GRANT ALL ON *.* TO 'dozenplans'@'%';
```

使用 *github.com/go-sql-driver/mysql* 连接数据库



golang的命名返回值 https://www.jianshu.com/p/ce58bc8885e2

## 接口设计

### User

Authorization 部分请登陆后复制出来

1. POST /api/login 登录

   ```
   Content-Type: application/x-www-form-urlencoded
   
   username=user
   &email=user@qq.com
   &password=password
   ```

2. POST /api/users/ 注册

   ```
   Content-Type: application/x-www-form-urlencoded
   
   username=user
   &email=user@qq.com
   &password=password
   ```

3. GET /api/users/:uid 获取用户信息 当前没有权限要求 之后加上管理员权限 权限级=2

4. GET /api/users/ 获取所有用户的信息 当前没有权限要求

5. PUT /api/users/ 更新用户信息 （uid会使用token里面的id，不需要声明）

   ```
   Content-Type: application/x-www-form-urlencoded
   Authorization: Bearer xxxxx
   
   username=user
   &email=user@qq.com
   &password=password
   
   ```

### Task

添加task的时候如果开启提醒， 有两种提醒模式 注意，传入的时间格式必须为RFC3339

- 定时提醒 NotifyMode = timing

  传入的NotifyTime就是提醒时间 （应该为当前时间+间隔 **前端来处理**）

- 间隔提醒 NotifyMode = interval

  传入的NotifyTime是首次提醒的时间,NotifyInterval是提醒的间隔  **单位为分钟**

1. *POST* /api/tasks 添加任务

   ```
   Content-Type: application/json
   Authorization: Bearer xxxxxx
   
   {
     "TaskName": "一个定时提醒任务",
     "Priority": 1,
     "DeadlineAt": "2021-06-24T14:12:12+08:00",
     "Status": "undone",
     "NotifyMode": "interval",   // 还有timing模式 或者留空
     "NotifyTime": "2021-06-25T15:57:00+08:00", 前端需要计算第一次提醒时间
     "NotifyInterval":120, 
     "Tags": "吃饭 睡觉 学习",
     "Category": "分类2"
   }
   ```

2. GET /api/tasks/1 根据本人id列出任务

3. *PUT* http://localhost:8080/api/tasks/:tid  更新一个任务

   ```
   Content-Type: application/json
   
   Authorization: Bearer xxxxx
   
   {
   //传入要更新的字段
   }
   ```

4. *DELETE* http://localhost:8080/api/tasks/:tid 删除任务 （）

   ```
   Authorization: Bearer xxxxx
   ```



### tag 和 category

```
### 获取一个用户的所有的tag
GET http://localhost:8080/api/tags
Content-Type: application/json
Authorization: Bearer xxxx

### 获取一个用户的一个tag下的所有task
GET http://localhost:8080/api/tags/10
Content-Type: application/json
Authorization: Bearer xxxxx
```



```
### 获取一个用户所有的分类
GET http://localhost:8080/api/categories
Authorization: Bearer xxxxxx
### 获取一个分类下所有的task id
GET http://localhost:8080/api/categories/1
Authorization: Bearer xxxx
```





