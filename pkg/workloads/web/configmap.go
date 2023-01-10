package web

import (
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	workloads "github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
)

// GenerateNginxConfigMap returns a configmap that will be consumed by the nginx service
func GenerateNginxConfigMap(namespace string) *workloads.ConfigMap {
	cm := &workloads.ConfigMap{}
	cm.Generate(map[string]string{"namespace": namespace, "name": constants.NginxConfName, "label": constants.NginxName})

	defaultConf := `log_format custom_format '$remote_addr - $remote_user [$time_local]'
'"$request" $status $body_bytes_sent'
'"$http_referer" "$http_user_agent"'
'$upstream_response_time';

server {
    listen       80 default_server;
    listen  [::]:80 default_server;
    server_name  _;
    error_log  /var/log/nginx/error.log;
    access_log /var/log/nginx/access.log custom_format;
    root /usr/share/nginx/html;
    index    index.html index.htm index.php;

    location ~ \.php$ {
        try_files $uri =404;
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_pass localhost:9000;
        fastcgi_index index.php;
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_param PATH_INFO $fastcgi_path_info;
    }
}`
	metrics := `server {
    listen       8080;
    listen  [::]:8080;
    server_name _;

    location /stub_status {
        stub_status on;
    }
}`

	cm.Resource.Data = map[string]string{
		"default": defaultConf,
		"metrics": metrics,
	}

	return cm
}

// GenerateWebpageConfigMap returns a configmap with the webpages configured
func GenerateWebpageConfigMap(namespace string) *workloads.ConfigMap {
	cm := &workloads.ConfigMap{}
	cm.Generate(map[string]string{"namespace": namespace, "name": constants.NginxPagesName, "label": constants.NginxName})

	cmCommon := `<?php
    $arr = array("success" => false, "data" => "");
    $retries = 5;

    function connectToDB($attempt){
        $servername = "database-e2e.%s";
        $username = "%s";
        $password = %s;
        $pgconn = '';

        try {
            $pgconn = new PDO("pgsql:host=$servername;port=5432;dbname=%s", $username, $password);
        } catch(PDOException $e) {
            if ($attempt > 5){
                return false;
        }
        sleep(5);
        connectToDB($attempt+1);
    }
    return $pgconn;
}
?>`

	index := `<?php
    include 'common.php';
    $conn = connectToDB(0);
    if(!$conn){
        $arr['success'] = false;
        $arr['data'] = "Couldn't connect to Postgres backend: " . $e->getMessage();
    }else{
        try {
            // set the PDO error mode to exception
            $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            $stmt = $conn->query("SELECT value FROM web");
            $result = $stmt->fetch();

            if (!$result) {
                $arr['success'] = false;
                $arr['data'] = "not ok";
                header("HTTP/1.0 500 Internal Server Error");
            }else{
                $arr['success'] = true;
                $arr['data'] = $result["value"];
                header("HTTP/1.0 200 OK");
            }
        } catch(PDOException $e) {
            $arr['success'] = false;
            $arr['data'] = "Failed to get data: " . $e->getMessage();
            header("HTTP/1.0 200 Internal Server Error");
        }
    }
    $conn = null;
    header('Content-type: application/json');
    echo json_encode($arr);
?>`
	healthz := `<?php
    include 'common.php';
    $conn = connectToDB(0);
    if(!$conn){
        $arr['success'] = false;
        $arr['data'] = "Couldn't connect to Postgres backend: " . $e->getMessage();
    }else{
        try {
            // set the PDO error mode to exception
            $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            $stmt = $conn->query("SELECT value FROM web");
            $result = $stmt->fetch();

            if (!$result) {
                header("HTTP/1.0 500 Internal Server Error");
                return;
            }
            $conn = null;
        } catch(PDOException $e) {
            header("HTTP/1.0 500 Internal Server Error");
            return;
        }
    }

    header("HTTP/1.0 200 OK");
    header('Content-type: application/json');
    echo json_encode(array("success" => true, "data" => "ok"));
?>`
	cm.Resource.Data = map[string]string{
		"common":  fmt.Sprintf(cmCommon, namespace, constants.DBUser, constants.DBPassword, constants.DBName),
		"index":   index,
		"healthz": healthz,
	}

	return cm
}
