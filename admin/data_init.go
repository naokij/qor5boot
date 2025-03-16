package admin

import (
	"fmt"

	"github.com/theplant/gofixtures"
	"gorm.io/gorm"
)

func EmptyDB(db *gorm.DB, tables []string) {
	for _, name := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", name)).Error; err != nil {
			panic(err)
		}
	}
}

// InitDB initializes the database with some initial data.
func InitDB(db *gorm.DB, tables []string) {
	// 重置序列号
	for _, name := range tables {
		if err := db.Exec(fmt.Sprintf("SELECT setval('%s_id_seq', (SELECT max(id) FROM %s));", name, name)).Error; err != nil {
			panic(err)
		}
	}
}

// composeS3Path to generate file path as https://cdn.qor5.com/system/media_libraries/236/file.jpeg.
func composeS3Path(filePath string) string {
	endPoint := s3Endpoint
	if endPoint == "" {
		endPoint = "https://cdn.qor5.com"
	}
	return fmt.Sprintf("%s/system/media_libraries%s", endPoint, filePath)
}

// GetNonIgnoredTableNames returns all table names except the ignored ones.
func GetNonIgnoredTableNames(db *gorm.DB) []string {
	ignoredTableNames := map[string]struct{}{
		"users":            {},
		"roles":            {},
		"user_role_join":   {},
		"login_sessions":   {},
		"qor_seo_settings": {},
	}

	var rawTableNames []string
	if err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema='public';").Scan(&rawTableNames).
		Error; err != nil {
		panic(err)
	}

	var tableNames []string
	for _, n := range rawTableNames {
		if _, ok := ignoredTableNames[n]; !ok {
			tableNames = append(tableNames, n)
		}
	}

	return tableNames
}

// 保留用于测试的示例数据
var OrdersExampleData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.orders (id, user_id, amount, created_at, updated_at) VALUES (1, 1, 100, '2023-01-01 00:00:00', '2023-01-01 00:00:00');
INSERT INTO public.orders (id, user_id, amount, created_at, updated_at) VALUES (2, 1, 200, '2023-01-02 00:00:00', '2023-01-02 00:00:00');
INSERT INTO public.orders (id, user_id, amount, created_at, updated_at) VALUES (3, 2, 300, '2023-01-03 00:00:00', '2023-01-03 00:00:00');
`, []string{"orders"}))

var PostsExampleData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.posts (id, title, body, hero_image, created_at, updated_at) VALUES (1, 'Demo', '<p>test edit</p>', '{"ID":1,"Url":"//qor5-test.s3.ap-northeast-1.amazonaws.com/system/media_libraries/1/file.jpeg","VideoLink":"","FileName":"demo image.jpeg","Description":"","FileSizes":{"@qor_preview":8917,"default":326350,"main":94913,"og":123973,"original":326350,"thumb":21199,"twitter-large":117784,"twitter-small":77615},"Width":750,"Height":1000}', '2023-01-05 00:00:00', '2023-01-05 00:00:00');
`, []string{"posts"}))
