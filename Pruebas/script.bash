#Crea disco mia1
mkdisk    -size=1000   -unit=K -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia

# Crea disco mia2
mkdisk   -Size=2000   -unit=K -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco2.mia"

#Crea disco mia3
mkdisk -Size=3000 -unit=K -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"

#Elimina los 3 discos
rmdisk -path=/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco1.mia
rmdisk -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/Disco2.mia"
rmdisk -path="/home/drop/Documentos/U/Lab_MIA/MIA_2S_P1_202200075/Pruebas/cosas mias/Disco3.mia"
