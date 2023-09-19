package main

import (
    "fmt"
    "io"
    "net/http"
    "encoding/json"
    "time"
    "regexp"
    "strings"
    "strconv"
)

type speechRecord struct{
    //会議URL_発言番号
    SpeechID         string     `json:"speechID"` 
    //URLに組み込まれている会議ID
    IssueID          string     `json:"issueID"` 
    //投稿内容種類?
    ImageKind        string     `json:"imageKind"` 
    //これ、発言番号と同じでは?
    SearchObject     int        `json:"searchObject"` 
    //第XXX回国会
    Session          int        `json:"session"` 
    //参議院/
    NameOfHouse      string     `json:"nameOfHouse"` 
    ///XX委員会
    NameOfMeeting    string     `json:"nameOfMeeting"` 
    //第XX号
    Issue            string     `json:"issue"` 
    //20XX-0X-0X
    Date             string     `json:"date"` 
    //国会が閉会中かどうか?nullが開いてる?
    Closing          string     `json:"closing"` 
    //会議における発言番号
    SpeechOrder      int        `json:"speechOrder"` 
    //発言者。姓名はスペースで区切らない。
    Speaker          string     `json:"speaker"` 
    //ひらがなで
    SpeakerYomi      string     `json:"speakerYomi"` 
    //政党
    SpeakerGroup     string     `json:"speakerGroup"` 
    //?
    SpeakerPosition  string     `json:"speakerPosition"` 
    //?
    SpeakerRole      string     `json:"speakerRole"` 
    Speech           string     `json:"speech"` 
    //PDFと同じ。その番号。
    StartPage        int        `json:"startPage"` 
    //会議URL/発言番号
    SpeechURL        string     `json:"speechURL"` 
    //会議URL
    MeetingURL       string     `json:"meetingURL"` 
    //会議がPDF化されている。そのページ番号まで。
    PdfURL           string     `json:"pdfURL"` 
}


//type x struct{
type public struct{
    ////Articles        string  `-`
    NumberOfRecords         int             `json:"NumberOfRecords"`
    NumberOfReturn          int             `json:"NumberOfReturn"` 
    StartRecord             int             `json:"StartRecord"` 
    NextRecordPosition      int             `json:"NextRecordPosition"` 
    SpeechRecord            []speechRecord  `json:"SpeechRecord"` 
}


var KokkaiArray []string 
var OldKokkai string
var Url         = "https://kokkai.ndl.go.jp/api/speech"
func do(){
    //words := strings.Split( "AI 著作権" , "," )
    i := 0
    words := strings.Split( Config.Words , "," )
    for( i < len( words )) {
        if words[i] == "" {
            continue
        }
        CheckWord( words[i] , "" )
        time.Sleep(3 * time.Second)
        i++
    }
    i = 0
    words = strings.Split( Config.Human , "," )
    for( i < len( words )) {
        if words[i] == "" {
            continue
        }
        CheckWord( words[i] , words[i] )
        time.Sleep(3 * time.Second)
        i++
    }
}


func CheckWord( Word string , Human string ) {
	re := regexp.MustCompile("T.*")

    req, _ := http.NewRequest("GET", Url , nil)
    q := req.URL.Query()
    q.Add("recordPacking"       , "json"    )

    //ここを空にすることはできない。
    if Human == "" {
        q.Add( "any"          , Word      )
    } else {
        q.Add( "any"          , Word      )
        q.Add( "speaker"      , Human     )
    }
    q.Add( "searchRange"  , "本文"    )


    //now := time.Now().AddDate( 0, 0, -3).Format(time.RFC3339)
    //q.Add( "from"         , now       )
    //q.Add( "until"         , now       )
    from := time.Now().AddDate( 0, 0, Config.From).Format(time.RFC3339)
	from = re.ReplaceAllString( from , "" )

    until := time.Now().AddDate( 0, 0, Config.Until).Format(time.RFC3339)
    until = re.ReplaceAllString( until , "" )
    q.Add( "from"         , from       )
    q.Add( "until"         , until       )
    fmt.Print( "検索開始日" )
    fmt.Println( from )
    fmt.Print( "検索終了日" )
    fmt.Println( until )
    //q.Add( "from"         , "2023-05-16"    )
    //q.Add( "until"         , "2023-05-16"    )

    req.URL.RawQuery = q.Encode()
    var client *http.Client = &http.Client{}
    resp, _ := client.Do( req )
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    var t public
    json.Unmarshal( body , &t)

    //検索結果が無かった場合表示しない。
    if Config.SearchZero == true && len( t.SpeechRecord ) == 0 {
        return
    }
    _text := "検索対象 : " + Word + "\\n"
    _i := 0
    for( _i < len( t.SpeechRecord ) ){
        _text = _text + "\\_a[OnKokkaiUrl,"
        _text = _text + t.SpeechRecord[_i].SpeechURL
        _text = _text + "]"

        _text = _text + t.SpeechRecord[_i].Date + " "

        _text = _text + t.SpeechRecord[_i].NameOfHouse + " "
        _text = _text + t.SpeechRecord[_i].NameOfMeeting + " "
        _text = _text + t.SpeechRecord[_i].Issue + " "
        _text = _text + "#" + strconv.Itoa( t.SpeechRecord[_i].SpeechOrder ) + " "
        if Human == "" {
            _text = _text + t.SpeechRecord[_i].Speaker + ""
        }

        _text = _text + "\\_a\\n"
        _i++
    }
    KokkaiArray = append( KokkaiArray , _text )
}






