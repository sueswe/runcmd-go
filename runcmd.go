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

	"github.com/pelletier/go-toml"
)

var REV = "DEV"
var version string = "0.6.6"
var configFile string = os.Getenv("HOME") + "/.runcmd.toml"
var home = os.Getenv("HOME")
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var debugLog = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

func CheckErr(e error) {
	if e != nil {
		infoLog.Println(e)
	}
}

// writes a LogFile
func writeLog(target string, message string) {
	current_time := time.Now()
	ts := current_time.Format(time.DateOnly)
	filename := target + "/runcmd_" + ts + ".csv"

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	CheckErr(err)
	defer f.Close()
	_, err2 := f.WriteString(message)
	CheckErr(err2)
}

// check for a ConfigFile
func readConfig(filename string) int {
	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	_, err := os.ReadFile(filename)
	if err != nil {
		infoLog.Println("(no config " + filename + " found)")
		return 1
	} else {
		infoLog.Println("Reading config: " + filename)
		return 0
	}

}

// returns runtime as string, returncode as int:
func run_with_p(command string, p string) (string, int) {

	if len(p) == 0 {
		p = " "
	}
	fields := strings.Fields(p)

	// https://devdocs.io/go/os/exec/index#Command
	// If name contains no path separators, Command uses
	// LookPath to resolve name to a complete path if possible. Otherwise it uses name directly as Path.
	cmd := exec.Command(command, fields...)
	infoLog.Println("running: ", command+" "+p)
	// Startzeitpunkt ermitteln:
	t1 := time.Now()
	infoLog.Println("Starttime: " + t1.String())
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

	c := 1
	for scanner.Scan() {
		z := strconv.Itoa(c)
		fmt.Println(z + ": " + scanner.Text())
		c = c + 1
	}
	if scanner.Err() != nil {
		cmd.Process.Kill()
		cmd.Wait()
		errorLog.Fatalln("Scanner.Error: ", scanner.Err())
		//errorLog.Println("")

	}

	err = cmd.Wait()
	t2 := time.Now()
	diff := t2.Sub(t1).Seconds()
	// diff2 := diff.Round(time.Second).String()
	diff2 := fmt.Sprintf("%.2f", diff)
	// infoLog.Println("DIFF2: " + diff2)
	rtc := -1
	if err != nil {
		errorLog.Println("Program exited not as expected.")
		infoLog.Println("Runtime: " + diff2 + " sec")
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
		infoLog.Println("Runtime: " + diff2 + " sec")
		r := strconv.Itoa(rtc)
		infoLog.Println("Exitstatus: " + r + ", Program exited OK.")
	}
	// return diff.String(), rtc
	return diff2, rtc
}

func main() {

	whereToLogTo := "nowhere"

	if readConfig(configFile) == 0 {
		// Static: the location of the configFile:
		config, err := toml.LoadFile(home + "/.runcmd.toml")
		CheckErr(err)
		baseDir := config.Get("default.RUNCMD_BASE").(string)
		runcmdPath := config.Get("default.RUNCMD_PATH").(string)
		mderr := os.Mkdir(os.Getenv(baseDir)+"/"+runcmdPath, 0750)
		CheckErr(mderr)
		whereToLogTo = os.Getenv(baseDir) + "/" + runcmdPath
	} else {
		infoLog.Println("Trying to use a default-log location.")
		whereToLogTo = home + "/runcmd_logging_rzomstp"
		mkdirResult := os.Mkdir(whereToLogTo, 0750)
		CheckErr(mkdirResult)
	}

	infoLog.Println("runcmd, Version ", version+", "+REV)
	if len(os.Args) <= 1 {
		infoLog.Println("Nothing to do.")
		os.Exit(1)
	}

	user := os.Getenv("USER")
	host := os.Getenv("HOSTNAME")
	yyyymmdd := os.Getenv("SMA_SCHEDULE_DATE")
	if len(yyyymmdd) == 0 {
		yyyymmdd = "NO_SMA_SCHEDULE_DATE"
	}
	jobname := os.Getenv("SMA_USER_SPECIFIED_JOBNAME")
	if len(jobname) == 0 {
		jobname = "NO_SMA_JOBNAME"
	}

	infoLog.Println("System: " + user + "@" + host + ", Job: " + jobname)

	runtime := ""
	returncode := -1
	command := os.Args[1]

	c1 := strings.Split(command, " ")

	parameter := os.Args[2:]
	parameterlist := strings.Join(parameter, " ")

	if len(c1) > 1 {
		// fmt.Println("Command split: ", c1[0])
		// fmt.Println("Param0 split: ", c1[1:])
		parameterlist = strings.Join(c1[1:], " ")
		command = c1[0]
	}

	path, err := exec.LookPath(command)
	if errors.Is(err, exec.ErrDot) {
		// debugLog.Println("IF reached: " + path)
		command = path
		command = "./" + command
		err = nil
	}
	if err != nil {
		// log.Fatal(err)
		errorLog.Println(err)
		os.Exit(126)
	}

	// RUNCMD_DRY exported?
	if os.Getenv("RUNCMD_DRY") != "" {
		returncode = 0
		t := time.Now()
		r := "undef"
		fmt.Println(command, parameter)
		writeLog(whereToLogTo, yyyymmdd+"; "+t.Format(time.RFC3339)+"; "+jobname+"; "+command+"; "+parameterlist+"; "+"t(s): "+runtime+"; "+"returncode: "+r+"\n")
		infoLog.Println("(Just RUNCMD_DRY exported. Nothing happened.)")
		os.Exit(returncode)
	} else {
		runtime, returncode = run_with_p(command, parameterlist)
	}
	r := strconv.Itoa(returncode)
	t := time.Now()
	writeLog(whereToLogTo, yyyymmdd+"; "+t.Format(time.RFC3339)+"; "+jobname+"; "+command+"; "+parameterlist+"; "+"t(s): "+runtime+"; "+"returncode: "+r+"\n")

	os.Exit(returncode)
}
