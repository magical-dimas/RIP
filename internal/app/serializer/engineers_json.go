package serializer

import "rip_project/internal/app/ds"

type EngineerJSON struct {
	ID           uint   `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	IsTechnician bool   `json:"is_technician"`
}

func EngineerToJSON(user ds.Engineers) EngineerJSON {
	return EngineerJSON{
		ID:           user.EngineerID,
		Login:        user.Login,
		Password:     user.Password,
		IsTechnician: user.IsTechnician,
	}
}

func EngineerFromJSON(j EngineerJSON) ds.Engineers {
	return ds.Engineers{
		Login:        j.Login,
		Password:     j.Password,
		IsTechnician: j.IsTechnician,
	}
}
