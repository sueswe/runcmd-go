package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var version string = "0.3"

func main() {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	//warningLog := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	infoLog.Println("runcmd, Version ", version)

	command := os.Args[1]
	parameter := os.Args[2:]
	parameterlist := strings.Join(parameter, " ")
	cmd := exec.Command(command, parameterlist)
	infoLog.Println("exec: ", command)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	if err != nil {
		errorLog.Println(err)
		os.Exit(3)
	}
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		errorLog.Println(err)
		os.Exit(3)
	}
	for scanner.Scan() {
		// outout:
		infoLog.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		errorLog.Fatalln("Scanner.Error: ", scanner.Err())
	}
	err = cmd.Wait()
	if err != nil {
		errorLog.Println("Program exited not as expected.")
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			//exitCode = ws.ExitStatus()
			errorLog.Println("ExitStatus: ", ws.ExitStatus())
			rtc := ws.ExitStatus()
			//rtc, _ = strconv.Atoi(rtc)
			//errorLog.Println(cmd.Stderr)
			os.Exit(rtc)
		} else {
			errorLog.Fatalln("Could not get exit code for failed program!")
			os.Exit(120)
		}
	} else {
		infoLog.Println("Program exited OK.")
	}
}
