package lib

const (
	//---------------------------------------------------------------------------
	// Name - имя файла
	// Type - тип контекста(video,img)
	// Ext - расширение файла
	// Size - размер в байтах файла
	// Data - base64 если `img`, и обычный массив байтов для `video`
	// Time - дата добавление записи в базу данных
	// Request - ссылка на картинку/видео
	// Desc - описание картинки
	//---------------------------------------------------------------------------

	TABLE_CREATE_SQL = `
		CREATE TABLE IF NOT EXISTS files(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Name TEXT UNIQUE,
		Type TEXT,
		Ext TEXT,
		Size INTEGER,
		Data BLOB,
		Time INTEGER,
		Request TEXT,
		Desc TEXT
	);`
	TABLE_INSERT_SQL = `
	INSERT OR REPLACE INTO files(
		Name,
		Type,
		Ext,
		Size,
		Data,
		Time,
		Request,
		Desc
	) values(?, ?, ?, ?, ?, ?, ?, ?);
	`
	TABLE_SELECT_ALL  = `
		SELECT * FROM files
		`

	TABLE_SELECT_SINGLE_BY_ID = `
		SELECT * FROM files
		WHERE Id == ?`

	TABLE_SELECT_ALL_BY_TYPE = `
		SELECT * FROM files
		WHERE Type == "?"
	`

	TABLE_DELETE_BY_ID_LIST = `
		SELECT * FROM files
		WHERE Id IN (?);`
)