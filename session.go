package session

import (
	"crypto/rand"
	"io"
	"encoding/base64"
	"time"
)

type Session struct {
	usid , ssid string
	uptime int64
}

type SessionMethods interface {
	VerifiedInfo (cssid string) bool
	GetTime () int64
	SetTime ()
	GetSsid () string
	SetSsid ()
	Destroy ()
}

var sessionMap = make(map[string]Session)

func (session *Session) VerifiedInfo (cssid string) bool {
	// cssid 是客户端( client )发送的 sessionId
	ssid := session.GetSsid( )
	if ssid == cssid {
		return true
	}
	return false
}

func (session *Session) GetTime () int64 {
	return session.uptime
}

func (session *Session) SetTime ()  {
	session.uptime=time.Now().Unix()
}

func (session *Session) GetSsid () string {
	return session.ssid
}

func (session *Session) SetSsid ()  {
	session.ssid = newSessionID()
	session.SetTime()
}

func (session *Session) Destroy () {
	delete( sessionMap , session.usid )
}

func GetSession(usid string )  *Session  {
	for k , v :=range sessionMap {
		if k == usid {
			return &v
		}
	}
	return nil
}

func InitSession( usid string , timeout int64 ) *Session {
	sess := Session {
		usid : usid  ,
		ssid : newSessionID(),
		uptime :time.Now().Unix(),
	}
	sessionMap[ usid ] = sess
	sessionTimeout(timeout , &sess ) // 设置 timeout 超时 , timeout=60 为60秒
	return &sess
}

func newSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func sessionTimeout(timeout int64, sess *Session)  {
	time.AfterFunc(time.Duration(timeout)  * time.Second , func () {
		//  获取 session 中的最后更新时间，并判断是否在当前时间的半小时前，是则删除 session ，否则重新计算超时时间
		if sess != nil {
			uptime := time.Unix(sess.GetTime(),0)
			half , _ := time.ParseDuration("-0.5h")
			beforeHalf := time.Now().Add(half)
			if uptime.Before( beforeHalf ) {
				sess.Destroy()
			}else{
				sessionTimeout(timeout , sess)
			}
		}
	})
}
