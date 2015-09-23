<?php
// MySQL Database Information
define("DB_HOST", "127.0.0.1");
define("DB_NAME", "test");
define("DB_USER", "root");
define("DB_PASS", "");

if ($_GET['cmd'] =="0")
	echo "1";
elseif ($_GET['cmd'] =="test") //test connection
	echo "ok";
elseif ($_GET['cmd'] =="ip")
	echo $_SERVER['REMOTE_ADDR']; //Get IP
else
	echo "Error";

?>