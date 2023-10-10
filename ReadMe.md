# ukagakaPlugin_NihonKokkai
日本の国会をチェックするためのSSPプラグイン。<br>
[国会会議録検索システム](https://kokkai.ndl.go.jp/#/result)<br>の検索APIを使用しています。
検索対象は本文中の内容や発言者です。<br>
検索にヒットした会議へのリンクをリスト状にしてバルーンに表示します。<br>


## 設定項目
プラグインフォルダ直下のConfig.jsonを編集します。
```
{
     "StartSec"     : 600 ,
     "IntervalSec"  : 40  ,
     "Words"        : "AI 著作権,原発,放送法"        ,
     "Human"        : "浜田聡,杉田水脈,高市早苗,赤松健,神谷宗幣",
     "SearchZero"   : true ,
     "From"         : -31,
     "Until"        : -31
}
```
- StartSec<br>
    SSP起動後に何秒後に検索内容を表示し始めるかを秒数で設定します。<br>
    検索自体は起動時に行われますが、早すぎる場合、検索結果の取得が間に合わず、指定した秒数より後になる場合があります。<br>
- IntervalSec<br>
    次の検索結果を表示するまでの間隔を秒数で指定します。<br>
- Words    <br>
    検索対象の文字列を指定します。<br>
    半角スペースで区切る事でAND検索になります。<br>
    カンマで複数指定することができます。その場合、IntervalSecが使用され複数回検索されます。<br>
- Human<br>
    発言者名で検索します。<br>
    姓名をスペースで区切らずに入力してください。<br>
    カンマで区切る事で複数回検索されます。<br>
- SearchZero<br>
    falseを指定すると、検索結果が無かった場合でも検索ワードの通知を行います。<br>
    trueを指定すると、検索結果がなかった回は通知を行わず、次の検索結果を表示します。<br>
- From<br>
    起動時を起点として-N日を指定します。<br>
    その日からの会議を検索します。<br>
- Until<br>
    起動時を起点として-N日を指定します。<br>
    その日までの会議を検索します。<br>
    FromとUtilを-30にした場合は30日前の会議を検索することになります。<br>


## 必要なもの
ネットワーク環境.


## 注意事項
会議が文字起こしされるのは少々遅れるようで、2週間からひと月ほど遅くなるように見えます。<br>
なので、From/Utilは-30ぐらいを起点とした方がよいでしょう。<br>

また、Golangを使用したDllなので、freelibrary等するとフリーズしてしまう問題があります。<br>
エクスプローラーからプラグインのOnOffなどをしない限り基本的に問題ないと思われます。<br>


## 実行環境
下記で動作確認をしています。<br>
- SSP 2.6.48以降
- Windows10



## License
MIT


## 他
このプラグインを使ったいかなる問題や損害に対して私は責任を負いません。


## Author
ambergon




