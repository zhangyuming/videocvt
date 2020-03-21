package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"syscall"
)

var m = map[int]exec.Cmd{}

type CmdStatus struct {
	Status string
	Pid int
	Error error
}

func Run(stc  chan CmdStatus,command string , arg ...string ) {
	logdone := make(chan struct{},2)
	defer func() {
		logdone <- struct{}{}
	}()
	cmd := exec.Command(command, arg...)
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		logrus.Debug("cmd.Start: %v",err)
		stc <- CmdStatus{
			Status: "stop",
			Pid:    -1,
			Error:  err,
		}
		return
	}

	cmdPid := cmd.Process.Pid //查看命令pid
	m[cmdPid] = *cmd
	logrus.Debug("  PID: ", cmdPid ,"    ",command,  arg, )

	//go logReader(stdout,logdone)

	stc <- CmdStatus{
		Status: "running",
		Pid:    cmdPid,
		Error:  nil,
	}

	var res int
	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			fmt.Println("cmd exit status")
			res = ex.Sys().(syscall.WaitStatus).ExitStatus()
		}
	}

	logrus.Debug("PID: " ,cmdPid , " is done , result code is ", res)
	logdone <- struct{}{}
	delete(m,cmdPid)

	stc <- CmdStatus{
		Status: "stop",
		Pid:    cmdPid,
		Error:  nil,
	}

	return

}

//func logReader(r io.Reader, done <- chan struct{})  {
//
//	for{
//		select {
//		case <- done:
//			return
//		default:
//
//			bufReader := bufio.NewReader(r)
//			line,err := bufReader.ReadBytes('\n')
//			fmt.Println(string(line))
//			if err == io.EOF {
//				return
//			}
//			if err != nil{
//				fmt.Println(err)
//				return
//			}
//		}
//	}
//
//}



func Kill(pid int) error {
	if cmd,ok := m[pid]; ok {
		if err := cmd.Process.Kill(); err != nil {
			return err
		}
	}
	return nil
}
