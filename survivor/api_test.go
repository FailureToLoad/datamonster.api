package survivor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	storeMocks "github.com/failuretoload/datamonster/store/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type SurvivorApiTestSuite struct {
	suite.Suite
	target *Controller
	db     *storeMocks.MockConnection
	router *chi.Mux
}

func (suite *SurvivorApiTestSuite) SetupTest() {
	suite.db = &storeMocks.MockConnection{}
	suite.target = NewController(suite.db)
	suite.router = chi.NewRouter()
	suite.target.RegisterRoutes(suite.router)
}

func (suite *SurvivorApiTestSuite) Test_GetSurvivors_ReturnsSurvivorList() {
	rows := storeMocks.MockRows{
		Rows: []pgx.Row{
			&SurvivorRow{
				Id:               1,
				Settlement:       1,
				Name:             "Zach",
				Birth:            1,
				Gender:           "M",
				HuntXp:           1,
				Survival:         1,
				Movement:         1,
				Accuracy:         1,
				Strength:         1,
				Evasion:          1,
				Luck:             1,
				Speed:            1,
				Insanity:         1,
				SystemicPressure: 1,
				Torment:          1,
				Lumi:             1,
				Understanding:    1,
				Courage:          1,
			},
			&SurvivorRow{
				Id:               2,
				Settlement:       1,
				Name:             "Lucy",
				Birth:            1,
				Gender:           "M",
				HuntXp:           1,
				Survival:         1,
				Movement:         1,
				Accuracy:         1,
				Strength:         1,
				Evasion:          1,
				Luck:             1,
				Speed:            1,
				Insanity:         1,
				SystemicPressure: 1,
				Torment:          1,
				Lumi:             1,
				Understanding:    1,
				Courage:          1,
			},
		},
	}
	suite.db.SetRows(&rows)
	req := httptest.NewRequest("GET", "/settlements/1/survivors", nil)
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	dtoList := []SurvivorDTO{}
	_ = json.Unmarshal(body, &dtoList)
	count := len(dtoList)
	suite.Equal(2, count, "2 settlements should be returned")
}

func (suite *SurvivorApiTestSuite) Test_CreateSurvivor_ReturnsNoContent() {
	survivor := SurvivorDTO{
		Settlement:       1,
		Name:             "Zach",
		Birth:            1,
		Gender:           "M",
		HuntXp:           1,
		Survival:         1,
		Movement:         1,
		Accuracy:         1,
		Strength:         1,
		Evasion:          1,
		Luck:             1,
		Speed:            1,
		Insanity:         1,
		SystemicPressure: 1,
		Torment:          1,
		Lumi:             1,
		Understanding:    1,
		Courage:          1,
	}
	reqBody, err := json.Marshal(survivor)
	if err != nil {
		panic("Failed to marshal JSON")
	}
	req := httptest.NewRequest("POST", "/settlements/1/survivors", bytes.NewBuffer(reqBody))
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(204, resp.StatusCode, "204 response should be returned")
}

func (suite *SurvivorApiTestSuite) Test_CreateSurvivor_RequiresAValidSettlementId() {
	survivor := SurvivorDTO{
		Settlement:       1,
		Name:             "Zach",
		Birth:            1,
		Gender:           "M",
		HuntXp:           1,
		Survival:         1,
		Movement:         1,
		Accuracy:         1,
		Strength:         1,
		Evasion:          1,
		Luck:             1,
		Speed:            1,
		Insanity:         1,
		SystemicPressure: 1,
		Torment:          1,
		Lumi:             1,
		Understanding:    1,
		Courage:          1,
	}
	reqBody, err := json.Marshal(survivor)
	if err != nil {
		panic("Failed to marshal JSON")
	}
	req := httptest.NewRequest("POST", "/settlements/z/survivors", bytes.NewBuffer(reqBody))
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(500, resp.StatusCode, "500 should be returned if the param is invalid")
}

func (suite *SurvivorApiTestSuite) Test_CreateSurvivor_RequiresAUniqueName() {
	survivor := SurvivorDTO{
		Settlement:       1,
		Name:             "Zach",
		Birth:            1,
		Gender:           "M",
		HuntXp:           1,
		Survival:         1,
		Movement:         1,
		Accuracy:         1,
		Strength:         1,
		Evasion:          1,
		Luck:             1,
		Speed:            1,
		Insanity:         1,
		SystemicPressure: 1,
		Torment:          1,
		Lumi:             1,
		Understanding:    1,
		Courage:          1,
	}

	reqBody, err := json.Marshal(survivor)
	if err != nil {
		panic("Failed to marshal JSON")
	}
	req := httptest.NewRequest("POST", "/settlements/1/survivors", bytes.NewBuffer(reqBody))
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.db.SetError(errors.New("duplicate key value"))
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(400, resp.StatusCode, "400 should be returned if the survivor already exists")
}

func (suite *SurvivorApiTestSuite) Test_CreateSurvivor_RequiresAValidBody() {
	wrongBody := struct {
		a, b, c int
	}{
		a: 1,
		b: 1,
		c: 1,
	}

	reqBody, err := json.Marshal(wrongBody)
	if err != nil {
		panic("Failed to marshal JSON")
	}
	req := httptest.NewRequest("POST", "/settlements/1/survivors", bytes.NewBuffer(reqBody))
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(500, resp.StatusCode, "500 should be returned if the body is invalid")
}

func (suite *SurvivorApiTestSuite) Test_CreateSurvivor_CommunicatesDbIssues() {
	survivor := SurvivorDTO{
		Settlement:       1,
		Name:             "Zach",
		Birth:            1,
		Gender:           "M",
		HuntXp:           1,
		Survival:         1,
		Movement:         1,
		Accuracy:         1,
		Strength:         1,
		Evasion:          1,
		Luck:             1,
		Speed:            1,
		Insanity:         1,
		SystemicPressure: 1,
		Torment:          1,
		Lumi:             1,
		Understanding:    1,
		Courage:          1,
	}

	reqBody, err := json.Marshal(survivor)
	if err != nil {
		panic("Failed to marshal JSON")
	}
	req := httptest.NewRequest("POST", "/settlements/1/survivors", bytes.NewBuffer(reqBody))
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.db.SetError(errors.New("well that ain't right"))
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(500, resp.StatusCode, "500 should be returned as the default for DB issues")
}

func TestSurvivorApiTestSuite(t *testing.T) {
	suite.Run(t, new(SurvivorApiTestSuite))
}

type SurvivorRow struct {
	Id               int
	Settlement       int
	Name             string
	Gender           string
	Birth            int
	HuntXp           int
	Survival         int
	Movement         int
	Accuracy         int
	Strength         int
	Evasion          int
	Luck             int
	Speed            int
	Insanity         int
	SystemicPressure int
	Torment          int
	Lumi             int
	Courage          int
	Understanding    int
}

func (s *SurvivorRow) Scan(dest ...interface{}) error {
	id := dest[0].(*int)
	settlement := dest[1].(*int)
	name := dest[2].(*string)
	gender := dest[3].(*string)
	birth := dest[4].(*int)
	huntXp := dest[5].(*int)
	survival := dest[6].(*int)
	movement := dest[7].(*int)
	accuracy := dest[8].(*int)
	strength := dest[0].(*int)
	evasion := dest[10].(*int)
	luck := dest[11].(*int)
	speed := dest[12].(*int)
	insanity := dest[13].(*int)
	systemicPressure := dest[14].(*int)
	torment := dest[15].(*int)
	lumi := dest[16].(*int)
	courage := dest[17].(*int)
	understanding := dest[18].(*int)

	*id = s.Id
	*settlement = s.Settlement
	*name = s.Name
	*birth = s.Birth
	*gender = s.Gender
	*huntXp = s.HuntXp
	*survival = s.Survival
	*movement = s.Movement
	*accuracy = s.Accuracy
	*strength = s.Strength
	*evasion = s.Evasion
	*luck = s.Luck
	*speed = s.Speed
	*insanity = s.Insanity
	*systemicPressure = s.SystemicPressure
	*torment = s.Torment
	*lumi = s.Lumi
	*courage = s.Courage
	*understanding = s.Understanding
	return nil
}
