package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"quests/internal"
	"quests/repository"
	"testing"
)

type mockQuestRepo struct {
	repository.Quest
	expectedQuestID    int
	expectedErr        error
	expectedCheckQuest int
}

func (m *mockQuestRepo) CreateQuest(questDB internal.NewQuestDB) (int, error) {
	return m.expectedQuestID, m.expectedErr
}

/*
	func (m *mockQuestRepo) GetQuestWithoutStep() ([]internal.Quests, error) {
		return nil, m.expectedErr
	}

	func (m *mockQuestRepo) GetSteps(q internal.Quests) ([]internal.Steps, error) {
		return nil, m.expectedErr
	}
*/
func (m *mockQuestRepo) CheckQuest(questDB internal.NewQuestDB) int {
	return m.expectedCheckQuest
}

func (m *mockQuestRepo) CreateQuestStep(newQuest internal.NewQuestStep) (int, error) {
	return 0, m.expectedErr
}

/*
func (m *mockQuestRepo) CreateQuestSteps(newQuestSteps internal.NewQuestSteps) error {
	return m.expectedErr
}
func (m *mockQuestRepo) UpdateQuestSteps(updateQuestStep internal.UpdateQuestStep) error {
	return m.expectedErr
}*/

func TestCreateQuest(t *testing.T) {
	type APITestCase struct {
		name         string
		questRepo    *mockQuestRepo
		inputData    internal.NewQuest
		expectedData *internal.Quests
		expectedErr  error
	}

	tests := []APITestCase{
		{
			name:         "Создание задания без шагов",
			questRepo:    &mockQuestRepo{expectedQuestID: 1},
			inputData:    internal.NewQuest{Name: "Футбол"},
			expectedData: &internal.Quests{Id: "1", QuestName: "Футбол", Steps: []internal.Steps{}},
			expectedErr:  nil,
		},
		{
			name:        "Создание задания без шагов. Без имени",
			questRepo:   &mockQuestRepo{expectedQuestID: 1},
			inputData:   internal.NewQuest{Name: ""},
			expectedErr: errors.New("Имя задания должно содержать от 1 до 200 символов"),
		},
		{
			name:        "Создание задания без шагов. Задание уже существует",
			questRepo:   &mockQuestRepo{expectedCheckQuest: 1},
			inputData:   internal.NewQuest{Name: "Футбол"},
			expectedErr: QuestExsistErr{1},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := NewQuestService(tc.questRepo)

			quest, err := srv.CreateQuest(tc.inputData)

			assert.Equal(t, tc.expectedData, quest)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

/*func TestUpdateQuestSteps(t *testing.T) {
	type APITestCase struct {
		name         string
		questRepo    *mockQuestRepo
		inputData    internal.UpdateQuestSteps
		expectedData internal.UpdateQuestSteps
		expectedErr  error
	}

	tests := []APITestCase{
		{
			name:      "Создание задания без шагов",
			questRepo: &mockQuestRepo{expectedQuestID: 1},
			inputData: internal.UpdateQuestSteps{QuestSteps: []internal.UpdateQuestStep{
				{
					Id:      2,
					Bonus:   6000,
					IsMulti: false,
				},
				{
					Id:      3,
					Bonus:   60,
					IsMulti: false,
				},
			},
			},
			expectedData: internal.UpdateQuestSteps{QuestSteps: []internal.UpdateQuestStep{
				{
					Id:      2,
					Bonus:   6000,
					IsMulti: false,
				},
				{
					Id:      3,
					Bonus:   60,
					IsMulti: false,
				},
			},
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := NewQuestService(tc.questRepo)

			quest, err := srv.CreateQuest(tc.inputData)

			assert.Equal(t, tc.expectedData, quest)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}*/
