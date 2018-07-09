# newdag
new blockchain with newdag consensus

使用方法：

1.将代码下载到进入目录 $GOPATH/src/github.com/ ,执行命令 ：git clone https://github.com/newDAG/newdag.git

2.代码下载完后，将$GOPATH/src/github.com/newdag/package目录下的 .a 文件,
 copy 到 $GOPATH/pkg/linux_amd64/github.com/newdag，如果没有目录，可以手动创建。

3.进入$GOPATH/src/github.com/newdag/docker 目录

运行 make bin, make conf ，make images

最后运行 make run

4.打开浏览器，进入地址 http://127.0.0.1:3881 即可看到，其中3881映射的是client1 docker,3882docker里的client2 docker, 依次类推，你可以在浏览器开多个窗口，查看试用区块同步操作情况

 ![image](https://github.com/newDAG/newdag/blob/master/screenshots/demo1.jpg)

 








 
