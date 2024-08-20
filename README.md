# MIA_2S_P1_Carnet

Este repositorio contiene el Proyecto 1 de la clase Manejo e Implementación de Archivos del segundo semestre de 2024. El objetivo de este proyecto es desarrollar una aplicación web para la gestión de un sistema de archivos EXT2, permitiendo su administración desde cualquier sistema operativo a través de una interfaz web. La aplicación incluye tanto el desarrollo de un frontend interactivo como un backend robusto implementado en Go.

## Descripción

La aplicación desarrollada en este proyecto permite:

- Administrar archivos y escribir estructuras en el sistema de archivos EXT2.
- Realizar formateo rápido y completo de particiones.
- Crear y gestionar particiones y usuarios en el sistema de archivos.
- Generar reportes visuales utilizando Graphviz.

## Arquitectura

- **Frontend**: Desarrollado utilizando Angular (o cualquier otro framework moderno). Se encarga de la interfaz de usuario y la interacción con el backend mediante APIs RESTful.
- **Backend**: Implementado en Go, maneja las operaciones sobre el sistema de archivos y expone APIs para su interacción con el frontend.

## Características

- **Gestión de discos y particiones**: Comandos como `MKDISK`, `RMDISK`, `FDISK`, y `MOUNT` para la creación y manipulación de discos y particiones virtuales.
- **Administración de usuarios**: Comandos como `MKUSR`, `RMUSR`, `LOGIN`, y `LOGOUT` para la gestión de usuarios y sus permisos en el sistema.
- **Operaciones sobre archivos**: Comandos como `MKFILE`, `MKDIR`, `CAT` para la creación y manipulación de archivos y directorios.
- **Reportes**: Generación de reportes detallados del sistema de archivos, como el MBR, inodos, bloques, y más, utilizando Graphviz.

## Requisitos

- **Frontend**: Angular (o cualquier otro framework moderno para el frontend).
- **Backend**: Go (Golang).
- **Sistema Operativo**: Distribución GNU/Linux para la ejecución del proyecto.

## Instrucciones de Uso

1. **Clonar el repositorio**:
    ```bash
    git clone https://github.com/usuario/MIA_2S_P1_Carnet.git
    ```

2. **Configuración del entorno**:
   - Instalar las dependencias necesarias para el frontend (Angular, etc.).
   - Instalar Go y las dependencias del backend.

3. **Ejecución del proyecto**:
   - Iniciar el backend con Go.
   - Iniciar el frontend y acceder a la interfaz web para interactuar con el sistema de archivos.

4. **Comandos disponibles**:
   - **MKDISK**: Crear un disco virtual.
   - **RMDISK**: Eliminar un disco virtual.
   - **FDISK**: Administrar particiones dentro del disco virtual.
   - **MOUNT**: Montar particiones.
   - **MKUSR**: Crear un nuevo usuario en el sistema.
   - **MKFILE**: Crear un archivo en el sistema.
   - **REP**: Generar reportes del sistema de archivos.

## Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo `LICENSE` para más detalles.

