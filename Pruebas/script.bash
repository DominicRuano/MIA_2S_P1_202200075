#Crea disco mia1
mkdisk    -size=1000   -unit=K -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia

#Crea disco mia2
mkdisk   -Size=2000   -unit=K -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco2.mia"

#Crea disco mia3
mkdisk -Size=3000 -unit=K -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"

#Elimina los 3 discos
#rmdisk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco2.mia" -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"

rmdisk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia
rmdisk -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco2.mia"
rmdisk -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"

#Genera el reporte del mbr del disco mia1
rep -name=mbr -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia

#Crea una particion en el disco mia1
fdisk -type=P -unit=K -fit=BF -name=Particion1 -size=100 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

#Monta la particion1 del disco mia1
mount -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia -name=Particion1


#
#
#


#> Crea el Disco1
mkdisk    -size=1000   -unit=K -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia

#> Elimina el Disco1
rmdisk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia

#> Elimina el Disco2 que no existe
rmdisk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco2.mia

#> Crea el Disco1 de nuevo
mkdisk    -size=1000   -unit=K -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia

#> Crea Disco3
mkdisk -Size=3000 -unit=K -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"

#> Crea todas las particiones del disco1
fdisk -type=P -unit=K -fit=BF -name=Particion1 -size=500 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=e -unit=k -fit=BF -name=Particion2 -size=10 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion3 -size=400 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion4 -size=50 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=l -unit=k -fit=BF -name=Particion5 -size=5 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=l -unit=k -fit=BF -name=Particion6 -size=1 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=l -unit=k -fit=BF -name=Particion7 -size=1 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=l -unit=k -fit=BF -name=Particion8 -size=2 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

#> Crea todas las particiones del disco3
fdisk -type=e -unit=k -fit=BF -name=Part2 -size=500 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=l -unit=k -fit=BF -name=Part8 -size=100 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=l -unit=k -fit=BF -name=Part6 -size=100 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=l -unit=k -fit=BF -name=Part7 -size=100 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=l -unit=k -fit=BF -name=Part5 -size=100 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=P -unit=K -fit=BF -name=Part3 -size=400 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=P -unit=K -fit=BF -name=Part4 -size=50 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  
fdisk -type=P -unit=K -fit=BF -name=Part1 -size=10 -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"  

#> Monta las particiones de Disco1
mount -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia -name=Particion1
mount -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia -name=Particion2
mount -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia -name=Particion3
mount -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia -name=Particion4

#> Monta las particiones de Disco3
mount -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia" -name=Part4
mount -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia" -name=Part2
mount -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia" -name=Part3
mount -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia" -name=Part1

#> Formatea una particion
mkfs -id=750A
mkfs -id=752A

#> Genera el reporte mbr del Disco1
rep -id=751A -name=mbr -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Reportes/repMBR1.pdf
rep -id=750B -name=mbr -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Reportes/repMBR3.pdf

#> Genera el reporte disk del Disco1 y Disco3
rep -id=751A -name=disk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Reportes/repDSK1.pdf
rep -id=750B -name=disk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Reportes/repDSK3.pdf
