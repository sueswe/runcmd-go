package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var version string = "0.5.5"

func CheckErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func writeLog(message string) {
	home := os.Getenv("HOME")
	current_time := time.Now()
	ts := current_time.Format("2006-01-02")
	filename := home + "/runcmd_logging_rzomstp/runcmd_" + ts + ".csv"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	CheckErr(err)
	defer f.Close()
	_, err2 := f.WriteString(message)
	CheckErr(err2)
}

// returns runtime as string, returncode as int:
func run_with_p(command string, p string) (string, int) {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)
	if len(p) == 0 {
		p = " "
	}
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
		fmt.Println(" " + scanner.Text())
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		errorLog.Fatalln("Scanner.Error: ", scanner.Err())
	}
	err = cmd.Wait()
	// Endzeitpunkt ermitteln:
	t2 := time.Now()
	diff := t2.Sub(t1).Seconds()
	// diff2 := diff.Round(time.Second).String()
	diff2 := fmt.Sprintf("%.2f", diff)
	// infoLog.Println("DIFF2: " + diff2)
	rtc := -1
	if err != nil {
		errorLog.Println("Program exited not as expected.")
		infoLog.Println("Runtime: " + diff2)
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
		rtc = 0
		infoLog.Println("Runtime: " + diff2)
		r := strconv.Itoa(rtc)
		infoLog.Println("Exitstatus: " + r + ", Program exited OK.")
	}
	// return diff.String(), rtc
	return diff2, rtc
}

func main() {
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	//warningLog := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	infoLog.Println("runcmd, Version ", version)
	if len(os.Args) <= 1 {
		infoLog.Println("Nothing to do.")
		os.Exit(1)
	}
	// la := strconv.Itoa(len(os.Args))
	// fmt.Println("Länge: " + la)

	runtime := ""
	returncode := -1
	command := os.Args[1]

	parameter := os.Args[2:]
	parameterlist := strings.Join(parameter, " ")

	if _, err := os.Stat(command); errors.Is(err, os.ErrNotExist) {
		// does not exists
	} else {
		command = "./" + command
	}
	runtime, returncode = run_with_p(command, parameterlist)

	r := strconv.Itoa(returncode)

	t := time.Now()
	yyyymmdd := os.Getenv("SMA_SCHEDULE_DATE")
	if len(yyyymmdd) == 0 {
		yyyymmdd = "no SMA_SCHEDULE_DATE"
	}
	jobname := os.Getenv("SMA_USER_SPECIFIED_JOBNAME")
	if len(jobname) == 0 {
		jobname = "no SMA_JOBNAME"
	}

	writeLog(yyyymmdd + "; " + t.Format(time.RFC3339) + "; " + jobname + "; " + command + "; " + parameterlist + "; " + "t(s): " + runtime + "; " + "returncode: " + r + "\n")

	os.Exit(returncode)
}
