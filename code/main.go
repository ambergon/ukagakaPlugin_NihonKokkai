package main

/*
   #include <windows.h>
   #include <stdlib.h>
   #include <string.h>
*/
import "C"

import (
    "fmt"
    "unsafe"
    "strings"
    "regexp"
)



func main() {
    fmt.Println( "test" )
}

var Sec         int
var Interval    int

var Directory    string
var References []string 
var CheckID         = regexp.MustCompile("^ID: ")
var CheckReference  = regexp.MustCompile("^Reference.+?: ")

type ResponseStruct struct {
    Shiori  string
    Sender  string
    Charset string
    Marker  string
    Value   string
}
func GetResponse( r *ResponseStruct ) string {
    V := ""
    if r.Value  != "" { V = "Value: "  + r.Value     + "\r\n" }
    res :=  r.Shiori    + "\r\n" + 
            r.Sender    + "\r\n" + 
            r.Charset   + "\r\n" + 
            V + "\r\n\r\n"
    return res
}


//export load
func load(h C.HGLOBAL, length C.long ) C.BOOL {
    fmt.Println( "load NihonKokkai" )
    Directory = C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( length ))
    fmt.Println( Directory  )

    //設定読み込み
    LoadJson()

    //初期化。
    Sec      = 0
    Interval = Config.IntervalSec

    //国会回収
    go do()

	C.GlobalFree( h )
	return C.TRUE
}


//export unload
func unload() bool {
    fmt.Println( "unload NihonKokkai" )
	return true
}


//export request
func request( h C.HGLOBAL, length *C.long ) C.HGLOBAL {
	RequestText := C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( *length ))
	C.GlobalFree( h )


    Value           := ""
    ID              := ""
    References      = []string{}
    //var NOTIFY bool = false

    Response := new( ResponseStruct )
    Response.Sender  = "Sender: NihonKokkai"
    Response.Charset = "Charset: UTF-8"

    //IDとReference
    //必要な情報を分解する。
    RequestLines := strings.Split( RequestText , "\r\n" )
    for _ , line := range RequestLines {
        if( line == "NOTIFY PLUGIN/2.0" ){
            //"GET PLUGIN/2.0";
            //NOTIFY = true

        } else if CheckID.MatchString( line )  {
            //fmt.Println( line )
            ID = CheckID.ReplaceAllString( line , "" )

        } else if CheckReference.MatchString( line )  {
            //fmt.Println( line )
            ref := CheckReference.ReplaceAllString( line , "" )
            References = append( References , ref )

        } else {
            //fmt.Println( line )
        }
    }

    //実行関数
    if ID == "OnOtherGhostTalk" {
    } else if ID == "OnSecondChange" && len( KokkaiArray ) != 0 {
        //3時間までカウント。
        if 10800 > Sec {
            Sec++
        }
        //実行開始
        if Sec >= Config.StartSec {
            if Interval >= Config.IntervalSec {
                Interval = 0
                Value         = KokkaiArray[0]
                KokkaiArray   = KokkaiArray[1:]
            }
            Interval++
        }

    } else if ID == "OnMenuExec"  {
    } else if ID == "OnKokkaiUrl"  {
        Value = "\\j[" + References[0] + "]" 
    } else {
        //fmt.Println( "no touch :" + ID )
        //fmt.Print( "NOTIFY : " )
        //fmt.Println( NOTIFY )
        //fmt.Print( "References : " )
        //fmt.Println( References )
        //fmt.Println( "" )
    }


    if Value == "" {
        Response.Shiori  = "PLUGIN/2.0 204 No Content"
    } else {
        Response.Shiori = "PLUGIN/2.0 200 OK"
        Response.Value  = Value
    }

    res_buf := C.CString( GetResponse( Response ))
    defer C.free( unsafe.Pointer( res_buf ))

	res_size := C.strlen( res_buf )
	ret      := C.GlobalAlloc( C.GPTR , ( C.SIZE_T )( res_size ))
	C.memcpy(( unsafe.Pointer )( ret ) , ( unsafe.Pointer )( res_buf ) , res_size )
	*length = ( C.long )( res_size )
	return ret
}


















