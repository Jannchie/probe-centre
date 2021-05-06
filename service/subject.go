package service

import (
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
)

func CreateSubject(subject model.Subject) error {
	if res := db.DB.Create(&subject); res.Error != nil {
		return res.Error
	}
	return nil
}
func UpdateSubject(subject model.Subject) error {
	if res := db.DB.Updates(&subject); res.Error != nil {
		return res.Error
	}
	return nil
}
func DeleteSubject(subject model.Subject) error {
	if res := db.DB.Delete(&subject); res.Error != nil {
		return res.Error
	}
	return nil
}

func ListSubjects(p int, ps int) ([]model.Subject, error) {
	var res []model.Subject
	result := db.DB.Offset((p - 1) * ps).Limit(ps).Find(&res)
	if result.Error != nil {
		return res, result.Error
	}
	return res, nil
}
