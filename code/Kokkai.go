package main

import (
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

func do() {
    //前日までの情報を拾えるように。
	re := regexp.MustCompile("T.*")
    now := time.Now().AddDate( 0, 0, -1).Format(time.RFC3339)
	now = re.ReplaceAllString( now, "" )
    //fmt.Println(now)


    url     := "https://kokkai.ndl.go.jp/api/speech"
    req, _ := http.NewRequest("GET", url , nil)

    q := req.URL.Query()
    q.Add("recordPacking"       , "json"    )
    //q.Add("speaker"       , "安倍晋三"    )
    //q.Add("any"       , "放送法"    )

    //ここを空にすることはできない。
    q.Add( "any"          , "案件"    )
    q.Add( "searchRange"  , "冒頭"    )
    q.Add( "from"         , now       )
    q.Add( "util"         , now       )
    //q.Add( "from"         , "2023-06-20"    )
    //q.Add( "util"         , "2023-06-20"    )

    req.URL.RawQuery = q.Encode()
    var client *http.Client = &http.Client{}
    resp, _ := client.Do( req )
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    //fmt.Println( string( body ) )

    var t public
    json.Unmarshal( body , &t)

    //最後にn/取得件数を表示しようかね。
    num := strconv.Itoa( t.NumberOfRecords )
    _i := 0
    var _array []string
    for( _i < len( t.SpeechRecord ) ){
        _text := "\\0\\b[2]\\_q"
        //fmt.Println( t.SpeechRecord[_i].SpeechID )
        //fmt.Println( t.SpeechRecord[_i].Date )
        //fmt.Println( t.SpeechRecord )
        //fmt.Print( "# " )

        _text = _text + "\\_a[OnKokkaiUrl,"
        _text = _text + t.SpeechRecord[_i].SpeechURL
        _text = _text + "]\\n"

        _text = _text + t.SpeechRecord[_i].NameOfHouse
        _text = _text + " "
        _text = _text + t.SpeechRecord[_i].NameOfMeeting
        _text = _text + " "
        _text = _text + t.SpeechRecord[_i].Issue 

        _text = _text + "\\_a\\n"

        str    := t.SpeechRecord[_i].Speech 
        str     = strings.Replace( str , "　" , " " , -1 )
        lines  := strings.Split( str , "\n" )
        //str = strings.Replace( str , "" , "" , -1 )
        _x := 0
        for( _x < len( lines )){
            //ここで特定のワードを含む内容化精査するぜ。
            check_seigan    := regexp.MustCompile("請願")
            check_tinjyou   := regexp.MustCompile("陳情書")
            check_hourituan := regexp.MustCompile("法律案")
            check_tyousyo   := regexp.MustCompile("調書")

            if check_seigan.MatchString( lines[ _x ] ) {
                _text = _text + lines[ _x ] + "\\n"
            }else if check_tinjyou.MatchString( lines[ _x ] ) {
                _text = _text + lines[ _x ] + "\\n"
            }else if check_hourituan.MatchString( lines[ _x ] ) {
                _text = _text + lines[ _x ] + "\\n"
            }else if check_tyousyo.MatchString( lines[ _x ] ) {
                _text = _text + lines[ _x ] + "\\n"
            }
            _x++
        }

        _I := strconv.Itoa( _i )
        _text = _text + _I + "/" + num + "\\_q"
        _array = append( _array , _text )
        _i++
    }
    KokkaiArray = _array
}




//まあ、普通に考えて、500を超える件数を見るのはしんどい。
//とりあえず、会議ごとに情報を回収するか。
//案件内容が泣ければ、表示しなくてもよいだろか?
//それはそれで、仕事しているのか監視したいな。
//会議の数自体は、50程度だったから、会議ごとに一分間隔で回すか。
//配列で回収してそれを廻そう。




