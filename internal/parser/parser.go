package parser

import (
	"HFLabs-test/internal/sheetsApi"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const (
	url             = "https://confluence.hflabs.ru/pages/viewpage.action?pageId=1181220999"
	startColumnsCnt = 10
)

type Parser struct {
	sheetsApi sheetsApi.SheetsApi
}

func New(api sheetsApi.SheetsApi) Parser {
	return Parser{
		sheetsApi: api,
	}
}

func (p *Parser) Run() error {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//экземпляры структур для записи в них объектов html
	temp := make([]string, startColumnsCnt)
	items := make([][]string, 0)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	//ищем названия колонок
	s := doc.Find("th")
	var columnCnt int

	for i := range s.Nodes {
		columnCnt++
		single := s.Eq(i)
		temp[i] = single.Text()
	}
	temp = temp[0:columnCnt]
	items = append(items, temp)

	//обходим ячейки
	s = doc.Find("td")
	temp = make([]string, columnCnt)
	for i := range s.Nodes {
		single := s.Eq(i)
		indx := i % columnCnt
		sel := single.Find("li")
		if sel.Is("li") { //проверяем на список
			sed := single.Find("p")
			if sed.Is("p") { //проверяем на заголовок списка
				temp[indx] = sed.Text() + "\n"
			}
			sed = single.Find("span") //идем по списку
			for k := range sed.Nodes {
				sec := sed.Eq(k)
				temp[indx] += "· " + sec.Text() + "\n"
			}
		} else {
			temp[indx] = single.Text()
		}
		if indx == columnCnt-1 {
			tmp := []string{temp[0], temp[1]}
			items = append(items, tmp)
		}
	}

	err = p.post(items) //отправляем данные
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) post(items [][]string) error {
	for i := 0; i < len(items); i++ {
		for k := 0; k < len(items[0]); k++ {
			err := p.sheetsApi.Put(i, k, items[i][k])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
