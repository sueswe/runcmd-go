package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
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
	home := os.Getenv("HOME")
	filename := home + "/runcmd_logging_rzomstp/runcmd.log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	CheckErr(err)
	defer f.Close()
	_, err2 := f.WriteString(message)
	CheckErr(err2)
}

func run_with_p(command string, p string) (string, int) {
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

	rtc := -1
	if err != nil {
		errorLog.Println("Program exited not as expected.")
		infoLog.Println("Runtime: " + diff.String())
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			//exitCode = ws.ExitStatus()
			errorLog.Println("ExitStatus: ", ws.ExitStatus())
			rtc = ws.ExitStatus()
			//rtc, _ = strconv.Atoi(rtc)
			//errorLog.Println(cmd.Stderr)
			//os.Exit(rtc)
		} else {
			errorLog.Fatalln("Could not get exit code for failed program!")
			os.Exit(120)
		}
	} else {
		infoLog.Println("Runtime: " + diff.String())
		infoLog.Println("Program exited OK.")
		rtc = 0
	}
	return diff.String(), rtc
}

func run(command string) (string, int) {
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

	rtc := -1
	if err != nil {
		errorLog.Println("Program exited not as expected.")
		// Endzeitpunkt ermitteln:

		infoLog.Println("Runtime: " + diff.String())
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			//exitCode = ws.ExitStatus()
			errorLog.Println("ExitStatus: ", ws.ExitStatus())
			rtc = ws.ExitStatus()
			//rtc, _ = strconv.Atoi(rtc)
			//errorLog.Println(cmd.Stderr)
			//os.Exit(rtc)
		} else {
			errorLog.Fatalln("Could not get exit code for failed program!")
			os.Exit(120)
		}
	} else {
		infoLog.Println("Runtime: " + diff.String())
		infoLog.Println("Program exited OK.")
		rtc = 0
	}
	return diff.String(), rtc
}

func main() {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	//warningLog := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	infoLog.Println("runcmd, Version ", version)

	runtime := ""
	returncode := -1
	command := os.Args[1]
	parameter := os.Args[2:]
	parameterlist := strings.Join(parameter, " ")
	if len(parameter) != 0 {
		runtime, returncode = run_with_p(command, parameterlist)
	} else {
		runtime, returncode = run(command)
	}
	r := strconv.Itoa(returncode)
	// 20250404; 2025-04-07_08:36:58; b32n59c ; lgkk_rest_caller.pl ; +w start +p lwbwzbz7 -d 0 ; t(s): 3 ; returncode: 0
	t := time.Now()
	// const YYYYMMDD = "2006-01-02"
	yyyymmdd := (t.Format("2006-01-02"))
	jobname := os.Getenv("SMA_JOBNAME")
	if len(jobname) == 0 {
		jobname = "n/a"
	}

	writeLog(yyyymmdd + "; " + t.Format(time.RFC3339) + "; " + jobname + "; " + command + "; " + parameterlist + "; " + "t(s): " + runtime + "; " + "returncode: " + r + "\n")

}
