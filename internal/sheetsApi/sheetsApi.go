package sheetsApi

import (
	"context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
	"log"
	"time"
)

const timeOut = 100 //константа задержки между отправками в мсек

type SheetsApi struct {
	sheet *spreadsheet.Sheet
}

func New() SheetsApi {

	data, err := ioutil.ReadFile("client_secret.json") //файл авторизации пользователя
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope) //файл конфигурации
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background()) //создаем клиента по конфигурации

	service := spreadsheet.NewServiceWithClient(client) //создаем сервис по клиенту

	spreadsheetID := "1qHZNYdV-auYL7AL_oISEHdWg6K95ofcL9xFNv4pKgVA" //ID таблицы
	spreadsheets, err := service.FetchSpreadsheet(spreadsheetID)    //Пытаемся достать таблицу
	if err != nil {
		log.Fatal(err)
	}

	sheet, err := spreadsheets.SheetByIndex(0) //выбираем именно этот лист
	if err != nil {
		log.Fatal(err)
	}
	return SheetsApi{
		sheet: sheet,
	}
}

func (s *SheetsApi) Put(row int, column int, data string) error {
	s.sheet.Update(row, column, data) //заносим данные

	err := s.sheet.Synchronize()           //синзронизируем
	time.Sleep(time.Millisecond * timeOut) //тайм-аут чтобы гугл апи не ругался на DDOS
	return err
}
