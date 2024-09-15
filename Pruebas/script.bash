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



#
#
#

mkdisk    -size=1000   -unit=K -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia


fdisk -type=P -unit=K -fit=BF -name=Particion1 -size=100 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion1 -size=200 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

fdisk -type=P -unit=K -fit=BF -name=Particion2 -size=300 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion2 -size=400 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

fdisk -type=P -unit=K -fit=BF -name=Particion3 -size=500 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion3 -size=600 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

fdisk -type=P -unit=K -fit=BF -name=Particion4 -size=700 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion4 -size=800 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

fdisk -type=P -unit=K -fit=BF -name=Particion5 -size=900 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   
fdisk -type=P -unit=K -fit=BF -name=Particion5 -size=1000 -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia   

rep -name=mbr -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia


