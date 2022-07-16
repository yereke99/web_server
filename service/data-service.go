package service

import (
	"fmt"
	"log"
	"pro1/dto"
	"pro1/entity"
	"pro1/repository"
	"net/http"
	"io/ioutil"
	"github.com/mashingan/smapping"
)

type DataService interface {
	InsertData(d dto.DataCreateDTO) entity.Data
	UpdateData(d dto.DataUpdateDTO) entity.Data
	DeleteData(d entity.Data)
	All() []entity.Data
	FindDataByID(dataID uint64) entity.Data
	FindDataByString(data string) string
	IsAllowedToEdit(userID string, dataID uint64) bool
}

type dataService struct {
	dataRepository repository.DataRepository
}

func NewDataService(dataRepo repository.DataRepository) DataService {
	return &dataService{
		dataRepository: dataRepo,
	}
}

func (s *dataService) InsertData(d dto.DataCreateDTO) entity.Data {
	data := entity.Data{}
	err := smapping.FillStruct(&data, smapping.MapFields(&d))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := s.dataRepository.InsertData(data)
	return res
}

func (s *dataService) UpdateData(d dto.DataUpdateDTO) entity.Data {
	data := entity.Data{}
	err := smapping.FillStruct(&data, smapping.MapFields(&d))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := s.dataRepository.UpdateData(data)
	return res
}

func (s *dataService) DeleteData(d entity.Data) {
	s.dataRepository.DeleteData(d)
}

func (s *dataService) All() []entity.Data {
	return s.dataRepository.AllData()
}

func (s *dataService) FindDataByID(dataID uint64) entity.Data {
	return s.dataRepository.FindByDataID(dataID)
}

func (s *dataService) IsAllowedToEdit(userID string, dataID uint64) bool {
	d := s.dataRepository.FindByDataID(dataID)
	id := fmt.Sprintf("%v", d.UserID)

	return userID == id
}


func (s *dataService) FindDataByString(data string) string {
	url := fmt.Sprintf("https://serpapi.com/search.json?engine=google&q=%s&api_key=e5ce6ba62e732aa23b48be8cf9726852cdd885f9482061ca5661bc2f13c32042", data)
  req, _ := http.NewRequest("GET", url, nil)

  res, _ := http.DefaultClient.Do(req)
  defer res.Body.Close()

  body, _ := ioutil.ReadAll(res.Body)

  return string(body)
}
