package parser

import (
	"HFLabs-test/internal/sheetsApi"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

const url = "https://confluence.hflabs.ru/pages/viewpage.action?pageId=1181220999"
const columnCnt = 2

type Parser struct {
	sheetsApi sheetsApi.SheetsApi
}

//Codes структура для хранения строчки таблицы
type codes struct {
	code        string
	description string
}

func New(api sheetsApi.SheetsApi) Parser {
	return Parser{
		sheetsApi: api,
	}
}

func (p *Parser) Run() error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//экземпляры структур для записи в них объектов html
	temp := codes{}
	items := make([]codes, 0)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	//ищем названия колонок
	s := doc.Find("th")
	for i := range s.Nodes {
		single := s.Eq(i)

		switch i % columnCnt { //проверка по номеру
		case 0:
			sel := single.Find("li")
			if sel.Is("li") { //проверка на список
				sed := single.Find("p") //заголовок списка
				temp.description = sed.Text() + "\n"

				sed = single.Find("span") //заголовок списка
				for k := range sed.Nodes {
					sec := sed.Eq(k)
					temp.description += "· " + sec.Text() + "\n"
				}
			} else {
				temp.code = single.Text()
			}
		case 1:
			sel := single.Find("li")
			if sel.Is("li") {
				sed := single.Find("p")
				temp.description = sed.Text() + "\n"
				sed = single.Find("span")
				for k := range sed.Nodes {
					sec := sed.Eq(k)
					temp.description += "· " + sec.Text() + "\n"
				}
			} else {
				temp.description = single.Text()
			}
			items = append(items, temp)
		default:
			log.Fatal("columnsCnt нее соответствует ширине таблицы")

		}
	}
	s = doc.Find("td")
	for i := range s.Nodes {
		single := s.Eq(i)
		switch i % columnCnt {
		case 0:
			sel := single.Find("li")
			if sel.Is("li") {
				sed := single.Find("p")
				temp.description = sed.Text() + "\n"
				sed = single.Find("span")
				for k := range sed.Nodes {
					sec := sed.Eq(k)
					temp.description += "· " + sec.Text() + "\n"
				}
			} else {
				temp.code = single.Text()
			}
		case 1:
			sel := single.Find("li")
			if sel.Is("li") {
				sed := single.Find("p")
				temp.description = sed.Text() + "\n"
				sed = single.Find("span")
				for k := range sed.Nodes {
					sec := sed.Eq(k)
					temp.description += "· " + sec.Text() + "\n"
				}
			} else {
				temp.description = single.Text()
			}
			items = append(items, temp)
		default:
			log.Fatal("Columns не соответствует ширине таблицы")
		}
	}

	err = p.post(items) //отправляем данные
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) post(items []codes) error {
	for i := 0; i < len(items); i++ {
		var err error
		if i == 0 {
			err = p.sheetsApi.Put(i, 0, "id") //для красоты
		} else {
			err = p.sheetsApi.Put(i, 0, strconv.Itoa(i)) //отправляем id
		}
		if err != nil {
			return err
		}

		err = p.sheetsApi.Put(i, 1, items[i].code) //отправляем код
		if err != nil {
			return err
		}

		err = p.sheetsApi.Put(i, 2, items[i].description) //отправляем описание
		if err != nil {
			return err
		}
	}
	return nil
}
