package duino

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Conn struct {
	net.Conn
}

type JobResult struct {
	// If the job result was successful or not
	Success bool
	// The reason for the job result not being successful
	Reason string
}

func NewConn(pool string) (Conn, error) {
	tcpconn, err := net.Dial("tcp", pool)
	if err != nil {
		return Conn{}, err
	}

	return WrapConn(tcpconn)
}

func WrapConn(tcpconn net.Conn) (Conn, error) {
	conn := Conn{tcpconn}
	return conn, conn.ReadVersion()
}

func (conn *Conn) ReadVersion() error {
	_, err := conn.Read(make([]byte, 3)) // read 3 bytes
	return err
}

func (conn *Conn) GetJob(username, difficulty string) (Job, error) {
	_, err := conn.Write([]byte("JOB," + username + "," + difficulty))
	if err != nil {
		return Job{}, nil
	}

	buf := make([]byte, 128)

	n, err := conn.Read(buf)
	if err != nil {
		return Job{}, err
	}

	jobInfo := strings.Split(string(buf[:n-1]), ",")
	if len(jobInfo) != 3 {
		return Job{}, fmt.Errorf("duinocoin returned %s", buf[:n])
	}

	diffInt, err := strconv.Atoi(jobInfo[2])
	return Job{
		Base:       jobInfo[0],
		Expected:   jobInfo[1],
		Difficulty: diffInt,
	}, err

}

func (conn *Conn) SubmitJob(res int, hashrate float32, minerName, rigIdentifier string) (JobResult, error) {
	var data string
	if hashrate != -1 {
		data = fmt.Sprintf("%d,%f,%s,%s", res, hashrate, minerName, rigIdentifier)
	} else {
		data = fmt.Sprintf("%d,,%s,%s", res, minerName, rigIdentifier)
	}

	_, err := conn.Write([]byte(data))
	if err != nil {
		return JobResult{}, err
	}

	buf := make([]byte, 64)

	n, err := conn.Read(buf)
	if err != nil {
		return JobResult{}, err
	}

	split := strings.Split(string(buf[:n-1]), ",")
	feedback := split[0]
	if feedback != "GOOD" && feedback != "BAD" && feedback != "BLOCK" {
		return JobResult{}, fmt.Errorf("duinocoin returned %s", buf[:n])
	}

	jobresult := JobResult{
		Success: feedback != "BAD",
	}

	if len(split) > 1 {
		jobresult.Reason = split[1]
	}

	return jobresult, nil
}
