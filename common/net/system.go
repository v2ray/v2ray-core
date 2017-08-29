package net

import "net"

var DialTCP = net.DialTCP
var DialUDP = net.DialUDP

var Listen = net.Listen
var ListenTCP = net.ListenTCP
var ListenUDP = net.ListenUDP

var LookupIP = net.LookupIP

var SplitHostPort = net.SplitHostPort

type Addr = net.Addr
type Conn = net.Conn

type TCPAddr = net.TCPAddr
type TCPConn = net.TCPConn

type UDPAddr = net.UDPAddr
type UDPConn = net.UDPConn

type UnixConn = net.UnixConn

type IP = net.IP

type Error = net.Error
type AddrError = net.AddrError

type Dialer = net.Dialer
type Listener = net.Listener
type TCPListener = net.TCPListener
type UDPListener = net.UDPListener
