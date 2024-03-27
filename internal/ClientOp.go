package internal

import "scutbot.cn/uniauth/model"

func AddClient(client *model.Client) error {
	return db.Create(client).Error
}

func GetAllClient() ([]model.Client, error) {
	var clients []model.Client
	return clients, db.Find(&clients).Error
}

func GetClientById(client_id string) (*model.Client, int64, error) {
	var client model.Client
	result := db.Where("client_id = ?", client_id).Find(&client)
	return &client, result.RowsAffected, result.Error
}

func UpdateClient(client *model.Client) error {
	return db.Save(client).Error
}
