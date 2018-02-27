package bingo

import (
	"log"
	"os"
	"runtime"
	"fmt"
)

const(
	LEVEL_DEBUG=1
	LEVEL_INFO=2
	LEVEL_WARN=3
	LEVEL_ERROR=4
)
var LEVEL_NAMES=[]string{"debug","info","warn","error"}
type LogConfig struct {
	IsDebug bool
	FileName string
}
type LogFactory struct {
	loglevel int
	logfile *os.File
	l *log.Logger
}


func (this *LogFactory)SetConfig(config LogConfig) {
	logFile,err:= os.OpenFile(config.FileName,os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModeAppend)
	log.Println("open log file",config.FileName)
	if err!=nil {
		log.Fatal("openfile error!",err.Error())
	}else {
		this.logfile = logFile
		this.l = log.New(this.logfile, "", log.LstdFlags)
	}


}

func (this *LogFactory)GetLog(name string)Log {
	return &logImp{name,this}
}

func (this *LogFactory)Close(){
	if this .logfile!=nil {
		this.logfile.Close()
	}
}

func (this *LogFactory)write(level int,prefix string,fmt string,obj ...interface{}) {
	   if level >= this.loglevel { //判断loglevel是否大于指定的level，如果大于则输出，否则直接抛弃
		   content := this.formatHeader(prefix, level) + fmt + "\n"
		   if this.l == nil {
			   go log.Printf(content, obj...)
		   } else {
			   go this.l.Printf(content, obj...)
		   }
	   }
}

func (this *LogFactory) formatHeader(prefix string,level int) string {
	_, thefile, line, _ := runtime.Caller(3)
	short := thefile
	for i := len(thefile) - 1; i > 0; i-- {
		if thefile[i] == '/' {
			short = thefile[i+1:]
			break
		}
	}
	return fmt.Sprintf("[%s][%s][%s(%d)]- ",LEVEL_NAMES[level-1],prefix,short,line)
}

type Log interface {
	Info(msg string,obj ...interface{})
	Debug(msg string,obj ...interface{})
	Error(msg string,obj ...interface{})
	Warning(msg string,obj ...interface{})
}

type logImp struct {
	module string
	factory *LogFactory
}
func (this *logImp)Info(msg string,obj ...interface{}){
   this.factory.write(LEVEL_INFO,this.module,msg,obj...)
}

func (this *logImp)Debug(msg string,obj ...interface{}){
	this.factory.write(LEVEL_DEBUG,this.module,msg,obj...)
}
func (this *logImp)Error(msg string,obj ...interface{}){
	this.factory.write(LEVEL_ERROR,this.module,msg,obj...)
}
func (this *logImp)Warning(msg string,obj ...interface{}){
	this.factory.write(LEVEL_WARN,this.module,msg,obj...)
}