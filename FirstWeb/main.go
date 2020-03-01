package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)
//用户结构

type UserMessage struct{
	Name string
	Id   string
	Phone string
}
//房间
type Room struct{
	Rtype string
	Rnum string
	Rsta string
	Message string
}

//入住预定
type Checkin struct{
	Username string
	Rtype string
	Rnum string
	Phone string
	Usetime string
}

//退房
type Checkout struct{
	Username string
	Rnum string
	Usetime string
}

//房间价格
type Roomcost struct{
	Rtype string
	Message string
	Rcost string
}

//预定信息
type RegisterMessage struct{
	Rname string
	Email string
	Phone string
	Need string
}

//入住预定记录
type userroom struct{
   Username string
   Rtype string
   Rnum string
   Phone string
   Usetime string
}

//退房记录
type userroom2 struct{
	Username string
	Rtype string
	Rnum string
	Phone string
	Starttime string
	Endtime string
	Alltime int
	Fee int
}

//数据库登录常量
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pan123"
	dbname   = "hotel"
)





//------------------------------处理函数------------------------------//
//计算时间差
func GetTimeArr(start, end string)float64{
	a,_:=time.Parse("2006-01-02",start)
	b,_:=time.Parse("2006-01-02",end)
	date:=b.Sub(a)
	time:=date.Hours()/24
	return time
}

//查找错误
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//验证用户的登陆信息
func UserTest(db *sql.DB, testname string, Usertype string) UserMessage{
	var user UserMessage
	var usertype string
	fmt.Println(" ************登录信息***************")
	fmt.Println("开始验证：","姓名：[",testname,"]","人员类型：[",Usertype,"]")
	if Usertype=="用户"{
		usertype="users"
	}else if Usertype=="管理"{
		usertype="manager"
	}else if Usertype=="接待"{
		usertype="receptionist"
	}else{
		fmt.Println("无效查询！！")
	}
	row, err := db.Query("select * from "+usertype+" where name ="+"'"+testname+"'")
	checkErr(err)
	for row.Next() {
		var name string
		var id string
		var phone string
		err = row.Scan(&name, &id, &phone)
		checkErr(err)
		user.Name=name
		user.Phone=phone
		user.Id=id
	}
	db.Close()
	return user
}

//添加预定信息
func RegisterMessageAdd(db *sql.DB, Name string, Email string,Phone string,Need string){
	stmt, err := db.Prepare("insert into RegisterMessage (Rname, Email, Phone,Need) values($1, $2, $3,$4)")
	if err != nil {
		fmt.Printf("could not Insert Mesage, %v", err)
	}
	_, err = stmt.Exec(Name,Email,Phone,Need)
	if err != nil {
		fmt.Printf("could not insert mesage, %v", err)
	}
	fmt.Println("success！！")
	db.Close()
}

//添加用户信息
func UserAdd(db *sql.DB, user UserMessage){

	stmt, err := db.Prepare("insert into users (name, id, phone) values($1, $2, $3)")
	if err != nil {
		fmt.Printf("could not Insert Mesage, %v", err)
	}
	_, err = stmt.Exec(user.Name,user.Id,user.Phone)
	if err != nil {
		fmt.Printf("could not insert mesage, %v", err)
	}
	fmt.Println("success！！")
	db.Close()
}

//添加入住信息
func Checkin_Add(db *sql.DB, check_in Checkin){
    fmt.Println("************入住登记************")
	stmt, err := db.Prepare("insert into userroom (username, rtype, rnum,phone,usetime) values($1, $2, $3,$4,$5)")
	if err != nil {
		fmt.Printf("could not Insert Mesage, %v", err)
	}
	_, err = stmt.Exec(check_in.Username,check_in.Rtype,check_in.Rnum,check_in.Phone,check_in.Usetime)
	if err != nil {
		fmt.Printf("could not insert mesage, %v", err)
	}
	fmt.Println("入住信息添加成功！")
	db.Close()
}

//添加预定信息
func reserve_Add(db *sql.DB, check_in Checkin){
	fmt.Println("************预定登记************")
	stmt, err := db.Prepare("insert into userroom2 (username, rtype, rnum,phone,usetime) values($1, $2, $3,$4,$5)")
	if err != nil {
		fmt.Printf("could not Insert Mesage, %v", err)
	}
	_, err = stmt.Exec(check_in.Username,check_in.Rtype,check_in.Rnum,check_in.Phone,check_in.Usetime)
	if err != nil {
		fmt.Printf("could not insert mesage, %v", err)
	}
	fmt.Println("预定信息添加成功！")
	db.Close()
}

//添加退房信息
func Checkout_Add(db *sql.DB, userroom3 userroom2){
	fmt.Println("************预定登记************")
	stmt, err := db.Prepare("insert into userroom3 (username, rtype, rnum,phone,starttime,endtime,alltime,fee) values($1, $2, $3,$4,$5,$6,$7,$8)")
	if err != nil {
		fmt.Printf("could not Insert Mesage, %v", err)
	}
	_, err = stmt.Exec(userroom3.Username,userroom3.Rtype,userroom3.Rnum,userroom3.Phone,userroom3.Starttime,userroom3.Endtime,userroom3.Alltime,userroom3.Fee)
	if err != nil {
		fmt.Printf("could not insert mesage, %v", err)
	}
	fmt.Println("退房信息添加成功！")
	db.Close()
}

//房间查找
func Room_find(db *sql.DB,rtype string,status string)[]Room{
	fmt.Println("************查找房间信息************")
	fmt.Println("房间类型："+"["+rtype+"]"+"房间状态："+status)
	rows, err := db.Query("select * from room where rtype ="+"'"+rtype+"'"+" and rsta ="+"'"+status+"'")
	checkErr(err)
	room:=[]Room{}
	for rows.Next() {
		room1:=Room{}
		var rtype string
		var rnum string
		var rsta string
		var message string
		err = rows.Scan(&rtype, &rnum, &rsta,&message)
		room1.Rtype=rtype
		room1.Rnum=rnum
		room1.Rsta=rsta
		room1.Message=message
		room=append(room, room1)
		checkErr(err)
	}
	db.Close()
	fmt.Println("查找结果：查找成功！")
	return room
}

func Room_find2(db *sql.DB)[]Room{
	rows, err := db.Query("select * from room where rsta='空闲'")
	checkErr(err)
	room:=[]Room{}
	for rows.Next() {
		room1:=Room{}
		var rtype string
		var rnum string
		var rsta string
		var message string
		err = rows.Scan(&rtype, &rnum, &rsta,&message)
		room1.Rtype=rtype
		room1.Rnum=rnum
		room1.Rsta=rsta
		room1.Message=message
		room=append(room, room1)
		checkErr(err)
	}
	db.Close()
	return room
}

func Room_find3(db *sql.DB,rnum string)[]Room{
	rows, err := db.Query("select * from room where rnum="+"'"+rnum+"'")
	checkErr(err)
	room:=[]Room{}
	for rows.Next() {
		room1:=Room{}
		var rtype string
		var rnum string
		var rsta string
		var message string
		err = rows.Scan(&rtype, &rnum, &rsta,&message)
		room1.Rtype=rtype
		room1.Rnum=rnum
		room1.Rsta=rsta
		room1.Message=message
		room=append(room, room1)
		checkErr(err)
	}
	db.Close()
	return room
}
//网订
func Net_find(db *sql.DB)[]RegisterMessage{
	rows, err := db.Query("select * from registermessage ")
	checkErr(err)
	registermessage:=[]RegisterMessage{}
	for rows.Next() {
		piece:=RegisterMessage{}
		var rname string
		var email string
		var phone string
		var need string
		err = rows.Scan(&rname,&email,&phone,&need)
		piece.Rname=rname
		piece.Email=email
		piece.Phone=phone
		piece.Need=need
		registermessage=append(registermessage, piece)
		checkErr(err)
	}
	db.Close()
	fmt.Println("查询结果：查找成功！")
	return registermessage
}




//房间删除
func Room_delete(db *sql.DB,num string){
	fmt.Println(" ************删除房间信息***************")
	fmt.Println("房间号码："+"["+num+"]")
	stmt, err := db.Prepare("delete from room where rnum=$1")
	if err != nil {
		fmt.Printf("could not Update Mesage, %v", err)
	}
	_, err = stmt.Exec(num)
	if err != nil {
		fmt.Printf("could not Update mesage, %v", err)
	}
	fmt.Println("删除结果：删除成功！")
	db.Close()
}

//房间价格查找
func Roomcost_find(db *sql.DB)[]Roomcost{
	fmt.Println("************查询房间价格表************")
	rows, err := db.Query("select * from roomcost ")
	checkErr(err)
	roomcost:=[]Roomcost{}
	for rows.Next() {
		room1:=Roomcost{}
		var rtype string
		var message string
		var rcost string
		err = rows.Scan(&rtype,&message,&rcost)
		room1.Rtype=rtype
		room1.Message=message
		room1.Rcost=rcost
		roomcost=append(roomcost, room1)
		checkErr(err)
	}
	db.Close()
	fmt.Println("查询结果：查找成功！")
	return roomcost
}

func Roomcost_find2(db *sql.DB,rtype string,rmessage string)[]Roomcost{
	rows, err := db.Query("select * from roomcost where rtype ="+"'"+rtype+"'"+" and rmessage ="+"'"+rmessage+"'")
	checkErr(err)
	roomcost:=[]Roomcost{}
	for rows.Next() {
		room1:=Roomcost{}
		var rtype string
		var message string
		var rcost string
		err = rows.Scan(&rtype,&message,&rcost)
		room1.Rtype=rtype
		room1.Message=message
		room1.Rcost=rcost
		roomcost=append(roomcost, room1)
		checkErr(err)
	}
	db.Close()
	fmt.Println("查询结果：查找成功！")
	return roomcost
}
func Userroom_search(db *sql.DB)[]userroom{
	rows, err := db.Query("select * from userroom ")
	checkErr(err)
	userroom1:=[]userroom{}
	for rows.Next() {
		room1:=userroom{}
		var username string
		var rtype string
		var rnum string
		var phone string
		var usetime string
		err = rows.Scan(&username,&rtype,&rnum,&phone,&usetime)
		room1.Rtype=rtype
		room1.Username=username
		room1.Rnum=rnum
		room1.Phone=phone
		room1.Usetime=usetime
		userroom1=append(userroom1, room1)
		checkErr(err)
	}
	db.Close()
	return userroom1
}
func Userroom_search_2(db *sql.DB,num string,username string)[]userroom{
	rows, err := db.Query("select * from userroom where rnum="+"'"+num+"'"+" and username ="+"'"+username+"'")
	checkErr(err)
	userroom1:=[]userroom{}
	for rows.Next() {
		room1:=userroom{}
		var username string
		var rtype string
		var rnum string
		var phone string
		var usetime string
		err = rows.Scan(&username,&rtype,&rnum,&phone,&usetime)
		room1.Rtype=rtype
		room1.Username=username
		room1.Rnum=rnum
		room1.Phone=phone
		room1.Usetime=usetime
		userroom1=append(userroom1, room1)
		checkErr(err)
	}
	db.Close()
	return userroom1
}

//预定记录查看
func Userroom_search2(db *sql.DB)[]userroom{
	rows, err := db.Query("select * from userroom2 ")
	checkErr(err)
	userroom1:=[]userroom{}
	for rows.Next() {
		room1:=userroom{}
		var username string
		var rtype string
		var rnum string
		var phone string
		var usetime string
		err = rows.Scan(&username,&rtype,&rnum,&phone,&usetime)
		room1.Rtype=rtype
		room1.Username=username
		room1.Rnum=rnum
		room1.Phone=phone
		room1.Usetime=usetime
		userroom1=append(userroom1, room1)
		checkErr(err)
	}
	db.Close()
	return userroom1
}

//退房记录查看
func Userroom_search3(db *sql.DB)[]userroom2{
	rows, err := db.Query("select * from userroom3 ")
	checkErr(err)
	userroom:=[]userroom2{}
	for rows.Next() {
		room1:=userroom2{}
		var username string
		var rtype string
		var rnum string
		var phone string
		var starttime string
		var endtime string
		var alltime int
		var fee int
		err = rows.Scan(&username,&rtype,&rnum,&phone,&starttime,&endtime,&alltime,&fee)
		room1.Rtype=rtype
		room1.Username=username
		room1.Rnum=rnum
		room1.Phone=phone
		room1.Starttime=starttime
		room1.Alltime=alltime
		room1.Endtime=endtime
		room1.Fee=fee
		userroom=append(userroom, room1)
		checkErr(err)
	}
	db.Close()
	return userroom
}

//房间价格修改处理器
func Roomcost_change(c string,c1 string,c3 string,db *sql.DB){
	fmt.Println(" ************修改房间信息***************")
	_,err := db.Query("update roomcost set rprice=" +"'"+c3+"'"+"where rtype ="+"'"+c+"'"+" and rmessage ="+"'"+c1+"'")
	checkErr(err)
	fmt.Println("房间类型："+"["+c+"]"+"房间配置:"+"["+c1+"]"+"价格修改为:"+"["+c3+"]")
	fmt.Println("修改结果：修改成功！")
	db.Close()
}

//连接数据库
func Connect()(*sql.DB){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

//修改房间状态
func room_update(db *sql.DB,num string){

	stmt, err := db.Prepare("UPDATE room set rsta='入住' where rnum=$1")
	if err != nil {
		fmt.Printf("could not Update Mesage, %v", err)
	}
	_, err = stmt.Exec(num)
	if err != nil {
		fmt.Printf("could not Update mesage, %v", err)
	}
	fmt.Println("房间状态修改成功！")
	db.Close()
}
func room_update2(db *sql.DB,num string){

	stmt, err := db.Prepare("UPDATE room set rsta='预定' where rnum=$1")
	if err != nil {
		fmt.Printf("could not Update Mesage, %v", err)
	}
	_, err = stmt.Exec(num)
	if err != nil {
		fmt.Printf("could not Update mesage, %v", err)
	}
	fmt.Println("修改房间状态成功！")
	db.Close()
}
func room_update3(db *sql.DB,num string){

	stmt, err := db.Prepare("UPDATE room set rsta='空闲' where rnum=$1")
	if err != nil {
		fmt.Printf("could not Update Mesage, %v", err)
	}
	_, err = stmt.Exec(num)
	if err != nil {
		fmt.Printf("could not Update mesage, %v", err)
	}
	db.Close()
}
func room_update4(db *sql.DB,num string){
	stmt, err := db.Prepare("delete from userroom  where rnum=$1")
	if err != nil {
		fmt.Printf("could not Update Mesage, %v", err)
	}
	_, err = stmt.Exec(num)
	if err != nil {
		fmt.Printf("could not Update mesage, %v", err)
	}
	db.Close()
}

//添加房间
func Roomadd(db *sql.DB,rtype string,num string,message string){
	fmt.Println(" ************添加房间信息***************")
	stmt, err := db.Prepare("insert into room (rtype, rnum,rsta,rmessage) values($1, $2, $3,$4)")
	if err != nil {
		fmt.Printf("could not Insert Mesage, %v", err)
	}
	_, err = stmt.Exec(rtype,num,"空闲",message)
	if err != nil {
		fmt.Printf("could not insert mesage, %v", err)
	}
	fmt.Println("房间类型："+"["+rtype+"]"+"房间配置:"+"["+message+"]"+"房间号为:"+"["+num+"]")

	db.Close()

}




//--------------------------------路由处理器------------------------------//
//登录页面处理器
func login(w http.ResponseWriter,r *http.Request){
	if r.Method=="GET"{
		t,_:=template.ParseFiles("login.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		r.ParseForm()
		Usertype:=r.FormValue("type")//获取前端传入的类型
		testname:=r.FormValue("name")//获取前端传入的用户名
		db := Connect()
		result := UserTest(db,testname,Usertype)
		if result.Id!=r.FormValue("password"){
			http.Redirect(w,r,"/",302)
			fmt.Println("验证结果："+"姓名:","[",testname,"]","验证错误，登录失败！")
		}else{
			fmt.Println("验证结果："+"姓名:","[",testname,"]","验证成功！")
			if Usertype=="接待"{
				http.Redirect(w,r,"/login",302)
			}else if Usertype=="管理"{
				http.Redirect(w,r,"/login2",302)
			}else if Usertype=="用户"{
				http.Redirect(w,r,"/login3",302)
			}
		}
	}
}

//注册页面处理器
func register(w http.ResponseWriter,r *http.Request){
	if r.Method=="GET"{
		t,_:=template.ParseFiles("register.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		r.ParseForm()
		var user UserMessage
		user.Name=r.FormValue("name")
		user.Phone=r.FormValue("phone")
		user.Id=r.FormValue("password")
		db := Connect()
		http.Redirect(w,r,"/",302)
		UserAdd(db,user)
	}
}

//入住登记处理器
func check_in(w http.ResponseWriter,r *http.Request){
	if r.Method=="GET"{
		t,_:=template.ParseFiles("check_in.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		r.ParseForm()
		var check_in Checkin
		check_in.Username=r.FormValue("name")
		check_in.Rtype=r.FormValue("type")
		check_in.Rnum=r.FormValue("number")
		check_in.Phone=r.FormValue("phone")
		check_in.Usetime=r.FormValue("time")
		db := Connect()
		Checkin_Add(db,check_in)
		db2 := Connect()
        room_update(db2,check_in.Rnum)
		fmt.Println(check_in.Username+"入住成功")
	}
}

//退房登记处理器
func check_out(w http.ResponseWriter,r *http.Request){
	if r.Method=="GET"{
		t,_:=template.ParseFiles("check_out.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		fmt.Println("***********退房申请***************")
		r.ParseForm()
		var check_out Checkout
		check_out.Username=r.FormValue("name")
		check_out.Rnum=r.FormValue("number")
		check_out.Usetime=r.FormValue("time")
		db:= Connect()
		userroom:=Userroom_search_2(db,check_out.Rnum,check_out.Username)
		time:=GetTimeArr(userroom[0].Usetime,check_out.Usetime)
        fmt.Println(time)
		db6:=Connect()
		room1:=Room_find3(db6,check_out.Rnum)

		db2:=Connect()
		room:=Roomcost_find2(db2,userroom[0].Rtype,room1[0].Message)
		cost,err:=strconv.Atoi(room[0].Rcost)
		if err!=nil{
			fmt.Println("类型转换失败！")
		}else{
			fmt.Println("房间单价："+room[0].Rcost)
		}
		fee:=int(time)*cost
		fmt.Println(fee)
        //插入退房记录
		var userroom2 userroom2
		userroom2.Username=check_out.Username
		userroom2.Rnum=check_out.Rnum
		userroom2.Fee=fee
		userroom2.Starttime=userroom[0].Usetime
		userroom2.Endtime=check_out.Usetime
		userroom2.Alltime=int(time)
		userroom2.Rtype=room[0].Rtype
		userroom2.Phone=userroom[0].Phone
		fmt.Println(userroom2)
		db3:=Connect()
		Checkout_Add(db3,userroom2)
        db4:=Connect()
		room_update3(db4, check_out.Rnum)
        db5:=Connect()
        room_update4(db5, check_out.Rnum)
		fmt.Println( check_out.Username+"退房成功")
	}
}

//预定登记处理器
func reserve(w http.ResponseWriter,r *http.Request){
	if r.Method=="GET"{
		t,_:=template.ParseFiles("reserve.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		r.ParseForm()
		var check_in Checkin
		check_in.Username=r.FormValue("name")
		check_in.Rtype=r.FormValue("type")
		check_in.Rnum=r.FormValue("number")
		check_in.Phone=r.FormValue("phone")
		check_in.Usetime=r.FormValue("time")
		db := Connect()
		reserve_Add(db,check_in)
		db2 := Connect()
		room_update2(db2,check_in.Rnum)
		fmt.Println(check_in.Username+"预定成功")
	}
}

//房间查找路由处理器
func room_search(w http.ResponseWriter,r *http.Request ){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm();
	if err!=nil{
		fmt.Println("解析失败！")
	}
	result:=r.PostForm
	c:=result.Get("type")
	c2:=result.Get("status")
	db := Connect()
	room:=Room_find(db,c,c2)
	json.NewEncoder(w).Encode(room)
}

//房间处理器
func room(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("room.html")
	t.Execute(w,nil)
}

//添加房间处理器
func room_add(w http.ResponseWriter,r *http.Request) {
	if r.Method=="GET"{
		t,_:=template.ParseFiles("room_add.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		err:=r.ParseForm();
		if err!=nil{
			fmt.Println("解析失败！")
		}
		result:=r.PostForm
		c1:=result.Get("type")
		c2:=result.Get("num")
		c3:=result.Get("messages")
		db := Connect()
		Roomadd(db,c1,c2,c3)
	}
}

//主页面处理器
func view(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("view.html")
	t.Execute(w,nil)

}
func view2(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("view2.html")
	t.Execute(w,nil)

}
func view3(w http.ResponseWriter,r *http.Request) {
	if r.Method=="GET"{
		t,_:=template.ParseFiles("view3.html")
		t.Execute(w,nil)
	}else if r.Method=="POST" {
		r.ParseForm()
		name := r.FormValue("name") //获取前端传入的类型
		email := r.FormValue("email") //获取前端传入的用户名
		need := r.FormValue("message") //获取前端传入的用户名
		phone := r.FormValue("phone") //获取前端传入的用户名
		db := Connect()
		RegisterMessageAdd(db, name , email ,phone ,need)
	}

}

//房间价格处理器
func roomcost(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("roomcost.html")
	t.Execute(w,nil)
}

//房间价格查询处理器
func roomcost_search(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm()
	if err!=nil{
		fmt.Println("解析失败！")
	}
	db := Connect()
	roomcost:=Roomcost_find(db)
	json.NewEncoder(w).Encode(roomcost)

}

//入住记录处理器
func userroom_search(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm()
	if err!=nil{
		fmt.Println("解析失败！")
	}
	db := Connect()
	userroom:=Userroom_search(db)
	json.NewEncoder(w).Encode(userroom)

}

//预定记录处理器
func userroom_search2(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm()
	if err!=nil{
		fmt.Println("解析失败！")
	}
	db := Connect()
	userroom:=Userroom_search2(db)
	json.NewEncoder(w).Encode(userroom)
}

//退房记录处理器
func userroom_search3(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm()
	if err!=nil{
		fmt.Println("解析失败！")
	}
	db := Connect()
	userroom:=Userroom_search3(db)
	json.NewEncoder(w).Encode(userroom)
}

//房间查询处理器
func Room_search(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm()
	if err!=nil{
		fmt.Println("解析失败！")
	}
	db := Connect()
	room:=Room_find2(db)
	json.NewEncoder(w).Encode(room)

}

//修改房间价格处理器
func roomcost_change(w http.ResponseWriter,r *http.Request) {
	if r.Method=="GET"{
		t,_:=template.ParseFiles("roomcost_change.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		err:=r.ParseForm();
		if err!=nil{
			fmt.Println("解析失败！")
		}
		result:=r.PostForm
		c1:=result.Get("Rtype")
		c2:=result.Get("Rmessage")
		c3:=result.Get("Rprice")
		db := Connect()
		Roomcost_change(c1,c2,c3,db)
	}
}

//删除房间处理器
func room_delete(w http.ResponseWriter,r *http.Request) {
	if r.Method=="GET"{
		t,_:=template.ParseFiles("room_delete.html")
		t.Execute(w,nil)
	}else if r.Method=="POST"{
		err:=r.ParseForm();
		if err!=nil{
			fmt.Println("解析失败！")
		}
		result:=r.PostForm
		c1:=result.Get("Rnum")
		db := Connect()
		Room_delete(db,c1)
	}
}

//入住记录处理器
func check_in_record(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("check_in_record.html")
	t.Execute(w,nil)
}

//预定记录处理器
func reserve_record(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("reserve_record.html")
	t.Execute(w,nil)
}

//退房记录处理器
func check_out_record(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("check_out_record.html")
	t.Execute(w,nil)
}

//网上信息处理器
func net_record(w http.ResponseWriter,r *http.Request) {
	t,_:=template.ParseFiles("net_record.html")
	t.Execute(w,nil)
}


func net_search(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	err:=r.ParseForm()
	if err!=nil{
		fmt.Println("解析失败！")
	}
	db := Connect()
	roomcost:=Net_find(db)
	json.NewEncoder(w).Encode(roomcost)

}


//主函数
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	/***设置路由器****/
	//登录
	http.HandleFunc("/", login)
	http.HandleFunc("/login", view)
	http.HandleFunc("/login2", view2)
	http.HandleFunc("/login3", view3)
	//注册
	http.HandleFunc("/register", register)
	//房间
	http.HandleFunc("/room", room)
	http.HandleFunc("/room_search", room_search)
	http.HandleFunc("/room_add", room_add)
	http.HandleFunc("/room_delete", room_delete)
	http.HandleFunc("/roomcost", roomcost)
	http.HandleFunc("/Room_search", Room_search)
	http.HandleFunc("/roomcost_search", roomcost_search)
	http.HandleFunc("/roomcost_change", roomcost_change)
	//注册
	http.HandleFunc("/check_in", check_in)
	//记录
	http.HandleFunc("/check_in_record", check_in_record)
	http.HandleFunc("/reserve_record", reserve_record)
	http.HandleFunc("/check_out_record", check_out_record)
	http.HandleFunc("/userroom_search", userroom_search)
	http.HandleFunc("/userroom_search2", userroom_search2)
	http.HandleFunc("/userroom_search3", userroom_search3)
	//预定
	http.HandleFunc("/reserve", reserve)
	//退房
	http.HandleFunc("/check_out", check_out)
	//网上信息
	http.HandleFunc("/net_record", net_record)
	http.HandleFunc("/net_search", net_search)
	err := http.ListenAndServe(":9090", nil)
	//监听
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}