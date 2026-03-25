package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"rip_project/internal/app/ds"
	"rip_project/internal/app/serializer"
)

func (r *Repository) GetEngineerByID(id int) (ds.Engineers, error) {
	var engineer ds.Engineers
	if id <= 0 {
		return ds.Engineers{}, fmt.Errorf("неверный id: должен быть > 0")
	}
	err := r.db.Where("engineer_id = ?", id).First(&engineer).Error
	if err != nil {
		return ds.Engineers{}, err
	}
	return engineer, nil
}

func (r *Repository) GetEngineerByLogin(login string) (ds.Engineers, error) {
	var engineer ds.Engineers
	if login == "" {
		return ds.Engineers{}, errors.New("логин не может быть пустым")
	}
	err := r.db.Where("login = ?", login).First(&engineer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Engineers{}, fmt.Errorf("%w: пользователь с логином %s не найден", ErrNotFound, login)
		}
		return ds.Engineers{}, err
	}
	return engineer, nil
}

func (r *Repository) CreateEngineer(j serializer.EngineerJSON) (ds.Engineers, error) {
	engineer := serializer.EngineerFromJSON(j)
	if engineer.Login == "" {
		return ds.Engineers{}, errors.New("логин обязателен для заполнения")
	}
	if engineer.Password == "" {
		return ds.Engineers{}, errors.New("пароль обязателен для заполнения")
	}
	_, err := r.GetEngineerByLogin(engineer.Login)
	if err == nil {
		return ds.Engineers{}, fmt.Errorf("%w: пользователь с логином %s уже существует", ErrAlreadyExists, engineer.Login)
	}
	if !errors.Is(err, ErrNotFound) {
		return ds.Engineers{}, err
	}
	if err := r.db.Create(&engineer).Error; err != nil {
		return ds.Engineers{}, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return engineer, nil
}

func (r *Repository) SignIn(j serializer.EngineerJSON) (ds.Engineers, error) {
	if j.Login == "" {
		return ds.Engineers{}, errors.New("логин обязателен для заполнения")
	}
	if j.Password == "" {
		return ds.Engineers{}, errors.New("пароль обязателен для заполнения")
	}
	engineer, err := r.GetEngineerByLogin(j.Login)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ds.Engineers{}, errors.New("неверный логин или пароль")
		}
		return ds.Engineers{}, err
	}
	if engineer.Password != j.Password {
		return ds.Engineers{}, errors.New("неверный логин или пароль")
	}
	return engineer, nil
}
