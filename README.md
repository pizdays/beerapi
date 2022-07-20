# beerapi

1. docker-compose up -d 
2. ทำการสร้าง database ชื่อว่า test ใน maria db
3. ทำการสร้าง database ชื่อว่า test ใน mongodb และทำการสร้าง user เพื่อทำการ access database
##############################################
   db.createUser(
   {
     user: "root",
     pwd: "password",  
     roles: [ "readWrite", "dbAdmin" ]
   }
  )
##############################################
4. ทำการ run คำสั่ง goreload main.go หรือ go run main.go
5. นำไฟล์ beerapi-insomnia.json ไป import ในโปรแกรม insomnia หรือโปรแกรมอื่นๆ
