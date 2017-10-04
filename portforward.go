package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	rcvport				string
	dstaddr				string
	err						error
	logFile				string
	incommingFile string
	outgoingFile	string
	inoutFile			string
	fi, fo, fio		*os.File
	flog					*os.File
	dupinaddr, dupoutaddr, dupioaddr	string
	dupinconn, dupoutconn, dupioconn	net.Conn
	client, target net.Conn
)

func init() {
	flag.StringVar(&logFile, "flog", "", "file to save log")
	flag.StringVar(&incommingFile, "fi", "", "file to save incomming packet")
	flag.StringVar(&outgoingFile, "fo", "", "file to save outgoing packet")
	flag.StringVar(&inoutFile, "fio", "", "file to save incomming/outgoing packet")
	flag.StringVar(&rcvport, "rcvport", "localhost:52767", "listening port <:port>")
	flag.StringVar(&dstaddr, "dstaddr", "127.0.0.1:8081", "forward to address <host:port>")
	flag.StringVar(&dupinaddr, "dupinaddr", "", "address <host:port> for send duplicate incomming packet")
	flag.StringVar(&dupoutaddr, "dupoutaddr", "", "address <host:port> for send duplicate outgoing packet")
	flag.StringVar(&dupioaddr, "dupioaddr", "", "address <host:port> for send duplicate incomming/outgoing packet")
}

func show(mode string, len int, buf []byte, src net.Conn) {
	if flog==nil { return }
	t:=time.Now()
	fmt.Fprintf(flog, "<%s-%s> [", t.Format("15:04:05.000"), mode)
	for i:=0;i<len;i++ {
		if i>0 { fmt.Fprintf(flog, ":") }
		fmt.Fprintf(flog, "%02x",buf[i])
	}
	fmt.Fprintf(flog, "] (%d) @%s\n", len, src.RemoteAddr())
	flog.Sync()
}

const BUFFER_SZ = 1024 
func forward(getin bool, dst net.Conn, src net.Conn) {
	// io.Copy(dst,src)
	defer dst.Close()
	defer src.Close()

	if getin==true {
		if incommingFile!="" {
			fi,err = os.Create(incommingFile)
			if err!=nil { panic(err) }
			defer fi.Close()
		}

		if inoutFile!="" {
			fio,err = os.Create(inoutFile)
			if err!=nil { panic(err) }
			defer fio.Close()
		}
		
		if dupinaddr!="" {
			if dupinconn!=nil { dupinconn.Close() }
			dupinconn,err = net.Dial("tcp", dupinaddr)
			if err!=nil { log.Fatalf("cannot connect to %s: %v", dupinaddr, err) }
			fmt.Printf("connection to dup incomming server %v established!\n", dupinconn.RemoteAddr())
			defer dupinconn.Close()
		}
		
		if dupioaddr!="" {
			if dupioconn!=nil { dupioconn.Close() }
			dupioconn,err = net.Dial("tcp", dupioaddr)
			if err!=nil { log.Fatalf("cannot connect to %s: %v", dupioaddr, err) }
			fmt.Printf("connection to dup incomming/outgoing server %v established!\n", dupioconn.RemoteAddr())
			defer dupioconn.Close()
		}
	} else {
		if outgoingFile!="" {
			fo,err = os.Create(outgoingFile)
			if err!=nil { panic(err) }
			defer fo.Close()
		}
		if dupoutaddr!="" {
			if dupoutconn!=nil { dupoutconn.Close() }
			dupoutconn,err = net.Dial("tcp", dupoutaddr)
			if err!=nil { log.Fatalf("cannot connect to %s: %v", dupoutaddr, err) }
			fmt.Printf("connection to dup outgoing server %v established!\n", dupoutconn.RemoteAddr())
			defer dupoutconn.Close()
		}
	}
	
	buf:=make([]byte, BUFFER_SZ)
	r:=bufio.NewReader(src)
	for {
		n,err:= r.Read(buf)
		if err!=nil {
			if err==io.EOF { return }
			fmt.Println(err)
			return
		}
		// data:=string(buf[:n])
		if getin==true {
			// fmt.Printf("Toinside : [%s] (%d)\n", data,n)
			// fmt.Print(">")
			show("c", n, buf, src) // from client
			if fi!=nil {
				fi.Write(buf[:n])
				fi.Sync()
			}
			if fio!=nil {
				fio.Write(buf[:n])
				fio.Sync()
			}
			if dupinconn!=nil { dupinconn.Write(buf[:n]) }
			if dupioconn!=nil { dupioconn.Write(buf[:n]) }
		} else {
			// fmt.Printf("Tooutside: [%s] (%d)\n", data,n)
			// fmt.Print("<")
			show("s", n, buf, src)		// from server
			if fo!=nil {
				fo.Write(buf[:n])
				fo.Sync()
			}
			if fio!=nil {
				fio.Write(buf[:n])
				fio.Sync()
			}
			if dupoutconn!=nil { dupoutconn.Write(buf[:n]) }
			if dupioconn!=nil { dupioconn.Write(buf[:n]) }
		}
		
		dst.Write(buf[:n])
	}
}

func main() {
	flag.Parse()
	fmt.Printf("== portforward version:25600720 ==\n")
	fmt.Printf("  use rcvport %v forward to dstaddr %v\n", rcvport, dstaddr)

	done:= make(chan bool, 1)
	sigs:= make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig:= <-sigs
		fmt.Printf("received signal [%v], stopping...\n", sig)
		done <- true
		os.Exit(1)
	}()
	
	conn, err:= net.Listen("tcp", rcvport)
	if err!=nil { log.Fatalf("cannot listen port %v: %v", rcvport, err) }
	fmt.Printf("== %s ==\n  Server running on %s\n", time.Now(), rcvport)
	defer conn.Close()
	
	if logFile!="" {
		flog,err = os.Create(logFile)
		if err!=nil { panic(err) }
		defer flog.Close()
	}

	for {
		client, err= conn.Accept()
		if err!=nil { log.Fatal("cannot accept client connection", err) }
		fmt.Printf("\n\n== %s ==\n", time.Now())
		fmt.Printf("connected from client %v\n", client.RemoteAddr())
		
		target, err= net.Dial("tcp", dstaddr)
		if err!=nil { log.Fatalf("cannot connect to %s: %v", dstaddr, err) }
		fmt.Printf("connection to server %v established!\n", target.RemoteAddr())
		
		go forward(true,  target, client)
		go forward(false, client, target)
	}
	<-done
}