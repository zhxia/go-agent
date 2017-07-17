# go-agent

* 基于golang在服务端部署server，客户端使用client，操作server进行一些基本的操作：
创建文件(目录)、删除文件(目录)、上传文件、远程运行脚本、列出目录

* 使用：  使用gb构建,gb build
 ./bin/server 
 <pre>
  Usage:./bin/server [options]
  	Options are:
  		-d dataroot Set data root
  		-s shellroot Set shell script root
  		-h hostport Set agent listen host and port
 </pre>
 <pre>
 ubuntu@myserver:/workspace/agent$ ./bin/server -s /data/script
 2017/07/17 08:52:57 server is running...
 2017/07/17 08:52:57 dataroot: /data ,shellroot: /data/script
 </pre>
 
 ./bin/client 
 <pre>
 Usage:./bin/client [options]
 Options are:
 	-c cmd Command type
 	-d data Command body data
 	-h host agent server host and port
 </pre>
 
 <pre>
 运行脚本：
 ./bin/client -c execute -d '{"shell":"/data/script/echo.sh","args":[]}' 
 
 创建目录：
  ./bin/client -c create -d '{"path":"/data/script","is_dir":true}'
 
 删除目录：
 ./bin/client -c delete -d '{"files":["/data/aa/bb/cc/dd"]}'
 
 上传文件：
 ./bin/client -c upload -d '{"file":"/data/logs/php_errors.log","override":true,"size":0}'
 
 列出目录文件：
 ./bin/client -c list -d {"path":"/data"}
 </pre>