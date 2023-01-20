/*
Copyright 2022 EscherCloud.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package web

import (
	"fmt"
	"github.com/eschercloudai/k8s-e2e-tester/pkg/constants"
	workloads "github.com/eschercloudai/k8s-e2e-tester/pkg/workloads/coreworkloads"
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
        $servername = "%s.%s";
        $username = "%s";
        $password = "%s";
        $pgconn = '';

        try {
            $pgconn = new PDO("pgsql:host=$servername;port=5432;dbname=%s", $username, $password);
        } catch(PDOException $e) {
            if ($attempt > 5){
                error_log("couldn't connect to the database!", 0);
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
    header('Content-type: application/json');
    if(!$conn){
        $err = "couldn't connect to Postgres backend: " . $e->getMessage();
        error_log("$err", 0);
        return;
    }else{
        try {
            // set the PDO error mode to exception
            $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            $stmt = $conn->query("SELECT value FROM web");
            $result = $stmt->fetch();

            if (!$result) {
                $err = json_encode($arr);
                error_log("no results in database: $err", 0);
                header("HTTP/1.0 500 Internal Server Error");
                return;
            }else{
                $arr['success'] = true;
                $arr['data'] = $result["value"];
                header("HTTP/1.0 200 OK");
            }
        } catch(PDOException $e) { 
            $err = "Failed to get data: " . $e->getMessage();
            error_log("$err", 0);
            header("HTTP/1.0 500 Internal Server Error");
            return;
        }
    }
    $conn = null;
    echo json_encode($arr);
?>`
	healthz := `<?php
    include 'common.php';
    $conn = connectToDB(0);
    header('Content-type: application/json');
    if(!$conn){
        $arr['success'] = false;
        $arr['data'] = "Couldn't connect to Postgres backend: " . $e->getMessage();
    	echo json_encode($arr);
		return;
    }else{
        try {
            // set the PDO error mode to exception
            $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            $stmt = $conn->query("SELECT value FROM web");
            $result = $stmt->fetch();

            if (!$result) {
                error_log("no results in database: $err", 0);
                header("HTTP/1.0 500 Internal Server Error");
                return;
            }
            $conn = null;
        } catch(PDOException $e) {
            $err = "Failed to get data: " . $e->getMessage();
            error_log("$err", 0);
            header("HTTP/1.0 500 Internal Server Error");
            return;
        }
    }

    header("HTTP/1.0 200 OK");
    echo json_encode(array("success" => true, "data" => "ok"));
?>`
	cm.Resource.Data = map[string]string{
		"common":  fmt.Sprintf(cmCommon, constants.PGSqlName, namespace, constants.DBUser, constants.DBPassword, constants.DBName),
		"index":   index,
		"healthz": healthz,
	}

	return cm
}
