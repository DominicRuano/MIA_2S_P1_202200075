package commands

import global "Backend/Global"

func Logout() string {
	if !global.Islogged {
		return "Error: No hay una sesión activa.\n"
	} else {
		global.Islogged = false
		global.User = ""
		global.Permiso = ""
		return "Sesión cerrada correctamente.\n"
	}
}
