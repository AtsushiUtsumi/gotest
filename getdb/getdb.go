package getdb

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Run関数は，データベースとの接続を確立し、APIサーバーを起動する
func Run() {
	// sql.Open関数を使用して，指定されたデータベースサーバーへの接続を開く
	db, err := sql.Open("postgres", "取得したいDBサーバーのリンク")

	// エラーが発生した場合は，panic関数でプログラムを終了する
	if err != nil {
		panic(err)
	}
	// defer db.Close()により，関数の最後にデータベース接続を閉じるようにする
	defer db.Close()

	// Ginフレームワークを使用してルーターを作成する
	router := gin.Default()

	// router.GETメソッドを使用して，/getdbエンドポイントを定義する
	router.GET("getdb", func(c *gin.Context) {
		// db.Query関数を使用し，指定されたSQLクエリを実行し，結果を取得する
		rows, err := db.Query("SELECT * FROM \"テーブル名\";")
		// エラーが発生した場合は，ステータスコード500（Internal Server Error）とエラーメッセージをJSONレスポンスで返す
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// defer rows.Close()により，関数の最後に結果セットを閉じるようにする
		defer rows.Close()

		// data変数を定義して，取得したデータを格納するためのスライスを初期化する
		var data []map[string]interface{}
		// rows.Columns()関数を使用して，結果セットのカラム名を取得する
		columns, _ := rows.Columns()
		// rows.Next()メソッドを使用して，結果セットの各行をループ処理する
		for rows.Next() {
			// 各行の値を格納するためのvaluesスライスとvaluePtrsスライスを作成する
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			// rows.Scan関数を使用して，現在の行の値をスキャンし，valuesスライスに格納する
			err := rows.Scan(valuePtrs...)
			// エラーが発生した場合は，ステータスコード500（Internal Server Error）とエラーメッセージをJSONレスポンスで返す
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// item変数を定義して，カラム名と値のマップを作成する
			item := make(map[string]interface{})
			for i, col := range columns {
				item[col] = values[i]
			}
			// dataスライスにitemを追加する
			data = append(data, item)
		}

		// c.JSON関数を使用して，ステータスコード200（OK）と取得したデータをJSONレスポンスで返す
		c.JSON(http.StatusOK, data)
	})

	// router.Runメソッドを使用して，指定されたポート番号でAPIサーバーを起動する
	// ポート番号は8080
	router.Run(":" + "8080")
}
