package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var version string = "0.4.3"

func CheckErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func writeLog(message string) {
	filename := "/tmp/runcmd.log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	CheckErr(err)
	defer f.Close()
	_, err2 := f.WriteString(message)
	CheckErr(err2)
}

func run_with_p(command string, p string) string {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)
	fields := strings.Fields(p)
	cmd := exec.Command(command, fields...)
	infoLog.Println("running: ", command+" "+p)
	// Startzeitpunkt ermitteln:
	t1 := time.Now()
	infoLog.Println("Starttime: " + t1.String())
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
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
		fmt.Println("out:" + scanner.Text())
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		errorLog.Fatalln("Scanner.Error: ", scanner.Err())
	}
	err = cmd.Wait()
	// Endzeitpunkt ermitteln:
	t2 := time.Now()
	diff := t2.Sub(t1)

	if err != nil {
		errorLog.Println("Program exited not as expected.")
		infoLog.Println("Runtime: " + diff.String())
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
		infoLog.Println("Runtime: " + diff.String())
		infoLog.Println("Program exited OK.")
	}
	return diff.String()
}

func run(command string) string {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)

	cmd := exec.Command(command)
	infoLog.Println("running: ", command)
	// Startzeitpunkt ermitteln:
	t1 := time.Now()
	infoLog.Println("Starttime: " + t1.String())
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
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
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		errorLog.Fatalln("Scanner.Error: ", scanner.Err())
	}
	err = cmd.Wait()
	t2 := time.Now()
	diff := t2.Sub(t1)

	if err != nil {
		errorLog.Println("Program exited not as expected.")
		// Endzeitpunkt ermitteln:

		infoLog.Println("Runtime: " + diff.String())
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
		infoLog.Println("Runtime: " + diff.String())
		infoLog.Println("Program exited OK.")
	}
	return diff.String()
}

func main() {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	//warningLog := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	infoLog.Println("runcmd, Version ", version)

	runtime := ""
	command := os.Args[1]
	parameter := os.Args[2:]
	parameterlist := strings.Join(parameter, " ")
	if len(parameter) != 0 {
		runtime = run_with_p(command, parameterlist)
	} else {
		runtime = run(command)
	}
	// 0250404; 2025-04-07_08:36:58; b32n59c ; lgkk_rest_caller.pl ; +w start +p lwbwzbz7 -d 0 ; t(s): 3 ; returncode: 0
	writeLog("Laufzeit: " + runtime + "\n")
}
