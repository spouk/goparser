package lib

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"log"
	"io"
	b64 "encoding/base64"
)

const (
	LOG_DBS_PREFIX = "[sqlitedbs-log]"
	LOG_DBS_FLAGS  = log.Ldate | log.Ltime | log.Lshortfile

	DB_LOG_MSG_SUCCESS_CREATE    = "база данных успешно создана/открыта для манипуляций"
	DB_LOG_MSG_SUCCESS_INSERT    = "запись успешно вставлена в базу данных"
	DB_LOG_MSG_SUCCESS_DELETE    = "запись(и) успешно удалена(ы) из базы данных"
	DB_LOG_MSG_SUCCESS_SELECT    = "запись(и) успешно выбрана(ы) из базы данных"
	DB_LOG_MSG_ERROR_TYPE_RECORD = "указан ошибочный тип данных для выборки"
	DB_LOG_MSG_ERROR_OPEN_DB     = "база данных не создана/не открыта"

	RECORD_TYPE_VIDEO = "video"
	RECORD_TYPE_IMG   = "img"
)

type SqliteDBS struct {
	db     *sqlite3.Conn
	dbname string
	logger *log.Logger
}
type FileObj struct {
	ID      int64
	Name    string
	Type    string
	Ext     string
	Size    int64
	Data    []byte
	Time    int64
	Request string
	Desc    string
}

func NewSqliteDBS(dbname string, logout io.Writer) *SqliteDBS {
	//создаю новый инстанс
	db := &SqliteDBS{
		dbname: dbname,
		logger: log.New(logout, LOG_DBS_PREFIX, LOG_DBS_FLAGS),
	}
	//создаю/открываю базу для манипуляций
	c, err := sqlite3.Open(db.dbname)
	if err != nil {
		db.logger.Printf(err.Error())
		panic(err)
	}
	db.db = c
	//создаю таблицу /если не создана
	err = db.db.Exec(TABLE_CREATE_SQL)
	if err != nil {
		panic(err)
	}
	//уведомление о результате операции
	db.logger.Printf(DB_LOG_MSG_SUCCESS_CREATE)
	//возвращаю результат
	return db
}
func (d *SqliteDBS) SelectAll(typeRecord string) []FileObj {
	//подготавливаю запрос
	sh, err := d.db.Prepare(TABLE_SELECT_ALL_BY_TYPE)
	if err != nil {
		d.logger.Printf(err.Error())
		return nil
	}
	defer sh.Close()

	switch typeRecord {
	case RECORD_TYPE_VIDEO:
		//делаю запрос к базе данных
		err = sh.Exec(RECORD_TYPE_VIDEO)
		if err != nil {
			d.logger.Printf(err.Error())
			return nil
		}
	case RECORD_TYPE_IMG:
		//делаю запрос к базе данных
		err = sh.Exec(RECORD_TYPE_IMG)
		if err != nil {
			d.logger.Printf(err.Error())
			return nil
		}

	default:
		d.logger.Printf(DB_LOG_MSG_ERROR_TYPE_RECORD)
		panic(DB_LOG_MSG_SUCCESS_CREATE)
	}

	//конвертирую результат
	result := []FileObj{}
	for {
		e := sh.Next()
		if e != nil {
			d.logger.Printf(err.Error())
			return result
		} else {
			fo := &FileObj{}
			err = sh.Scan(fo.ID, fo.Name, fo.Type, fo.Size, fo.Ext, fo.Time)
			if err != nil {
				d.logger.Printf(err.Error())
				return nil
			} else {
				result = append(result, *fo)
			}
		}
	}
	return result

}
func (d *SqliteDBS) SelectSinglebyID(id int64) *FileObj {
	//подготавливаю запрос
	sh, err := d.db.Prepare(TABLE_SELECT_SINGLE_BY_ID)
	if err != nil {
		d.logger.Printf(err.Error())
		return nil
	}
	defer sh.Close()
	//делаю запрос к базе данных
	err = sh.Exec(id)
	if err != nil {
		d.logger.Printf(err.Error())
		return nil
	}
	fo := &FileObj{}
	//конвертирую результат
	err = sh.Scan(fo.ID, fo.Name, fo.Type, fo.Size, fo.Ext, fo.Time)
	if err != nil {
		d.logger.Printf(err.Error())
		return nil
	}
	//возвращаю результат
	return fo
}
func (d *SqliteDBS) InsertRecord(fo FileObj) {
	if d.db == nil {
		d.logger.Printf(DB_LOG_MSG_ERROR_OPEN_DB)
		panic(DB_LOG_MSG_ERROR_OPEN_DB)
		return
	}
	sh, err := d.db.Prepare(TABLE_INSERT_SQL)
	if err != nil {
		panic(err)
	}
	defer sh.Close()
	//конвертирую в base64 если тип данных `img`
	blobtoWrite := []byte{}
	if fo.Type == RECORD_TYPE_IMG {
		blobtoWrite = []byte(b64.StdEncoding.EncodeToString(fo.Data))
	} else {
		blobtoWrite = fo.Data
	}
	//	Name,
	//	Type,
	//	Ext,
	//	Size,
	//	Data,
	//	Time,
	//	Request,
	//	Desc

	//запись в базу данных
	err = sh.Exec(fo.Name, fo.Type, fo.Ext, fo.Size, blobtoWrite, fo.Time, fo.Request, fo.Desc)
	if err != nil {
		d.logger.Printf(err.Error())
		//panic(err)
	}
	return
}
