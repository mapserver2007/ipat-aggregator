# 馬券集計

1. PATのページから購入結果CSVファイルを取得して`csv`に置く
2. Google SpreadSheetを新規作成する
3. SpreadSheetに権限を与えたときに取得できるjsonファイルを`secret/secret.json`として保存する
4. `secret/spreadsheet.json`を新規作成し、以下の値を設定する
```
{
  "spreadsheet_id": "xxxxxxxxxxxxxxxxxxxxxxxx", // SpreadSheetのURLに書いてあるID
  "spreadsheet_sheet_name": "シート名" // 書き出したいシート名
}
```
5. go mod tidy
6. go run cmd/main.go