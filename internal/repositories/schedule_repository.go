package repositories

import (
	"ptm/internal/db"
	"ptm/internal/models"
)

type ScheduleRepository struct{}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{}
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	return db.DB.Create(schedule).Error
}

func (r *ScheduleRepository) FindByID(id uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := db.DB.First(&schedule, id).Error
	return &schedule, err
}

func (r *ScheduleRepository) Update(schedule *models.Schedule) error {
	return db.DB.Save(schedule).Error
}

func (r *ScheduleRepository) Delete(id uint) error {
	return db.DB.Delete(&models.Schedule{}, id).Error
}

func (r *ScheduleRepository) List() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := db.DB.Find(&schedules).Error
	return schedules, err
}
