# 宙斯监控系统

本项目为运维监控服务端程序，后期会继续宙斯监控系统开发。

1. 从个人git迁移过来[http://git.53it.net/zuoxiupeng/zeus](http://git.53it.net/zuoxiupeng/zeus) 

2. 项目中zql查询语言解析 解析库：[https://github.com/shiguanghuxian/zql](https://github.com/shiguanghuxian/zql) 在线语句转换工具：[https://github.com/shiguanghuxian/zql-convert](https://github.com/shiguanghuxian/zql-convert)

3. collectd输出到zues的一个插件，通过collectd可以实现zues的一些监控 [http://git.53it.net/zuoxiupeng/write-zues](http://git.53it.net/zuoxiupeng/write-zues)



## 1.概述
>一个从采集->处理->展现，的完整处理流程，本程序可以实现收集主机性能数据和日后扩展的其他数据，并分析处理，最终缓存到内存中，在用户查看系统时以最快的速度响应，提高用户体验


## 备注
> 由于最近工作和对一些概念认识，打算重构本项目
> 
> 重构目的，对于数据落地存储插件使用插件机制，重新整理项目结构。