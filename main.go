package main

import (
	"encoding/json"
	"log"
	"net/http"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	_ "github.com/lib/pq"
)

const (
	host     = "lallah.db.elephantsql.com"
	port     = 5432
	user     = "xprfqdkt"
	password = "mDyybaXNlvSnzq0gcuiM9dalbfbPmGrE"
	dbname   = "xprfqdkt"
)

type Meeting struct {
	Id           int      `json:"Id"`
	Title        string   `json:"title"`
	Participants []string `json:"participants"`
	StartTime    int64    `json:start`
	EndTime      int64    `json:end`
}
type Participant struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Email  string             `json:"email,omitempty" bson:"email,omitempty"`
	RSVP  string             `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
	Meeting  []Meeting             `json:"meeting,omitempty" bson:"meeting,omitempty"`

}

type Error struct {
	ErrorId int    `json:"ErrorId"`
	Desc    string `json:"Desc"`
}

var Meetings Meeting
var error Error

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func scheduleMeeting(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Just scheduled a meeting!")

	//business logic
	if Meetings.StartTime > Meetings.EndTime {
		error = Error{ErrorId: 101, Desc: "Meeting end time before start time"}
		json.NewEncoder(w).Encode(error)
		return
	}
	if Meetings.Title == "" {
		error = Error{ErrorId: 102, Desc: "Meeting has no title"}
		json.NewEncoder(w).Encode(error)
		return
	}
	if len(Meetings.Participants) == 0 {
		error = Error{ErrorId: 103, Desc: "Meeting has no participant"}
		json.NewEncoder(w).Encode(error)
		return
	}
	//DB LOGIC
	client, err := mongo.NewClient(options.Client().ApplyURI("<mongodb+srv://mohit:<mohit143kumath>@cluster0.9yilk.mongodb.net/<dbname>?retryWrites=true&w=majority>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}
func CheckRsvp(participant Participant) bool{
	if participant.RSVP == "yes"{
		return false
	}
	return true
}


func GetMeeting(response http.ResponseWriter, request *http.Request)  {
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Request-Method","GET")
	var meeting Meeting

	log.Print(meeting.Title)
	log.Print(request.URL.Query().Get("id"))
	idUser :=request.URL.Query().Get("id")
	collection := client.Database("mydb").Collection("meeting")

	ctx,_ := context.WithTimeout(context.Background(), 5*time.Second)

	err := collection.FindOne(ctx, bson.M{"_id": idUser}).Decode(&meeting)
	//log.Print(err)
	if err!=nil {
		log.Print(err)
		json.NewEncoder(response).Encode("no result found")
		return
	}

	json.NewEncoder(response).Encode(meeting)
	client, err := mongo.NewClient(options.Client().ApplyURI("<mongodb+srv://mohit:<mohit143kumath>@cluster0.9yilk.mongodb.net/<dbname>?retryWrites=true&w=majority>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}
func CreateParticipant(participants Meeting, Meetings Meeting) (error,string) {

	if participant.Name=="" || participant.Email=="" || participant.RSVP == ""{
		return errors.New("please fill all the details"), string(0)
	}
	var error Participant
	collection := client.Database("databases").Collection("participant")

	ctx,_ := context.WithTimeout(context.Background(), 5*time.Second)

	err := collection.FindOne(ctx, bson.M{"email": participant.Email}).Decode(&error)

	var neresult Participant
	err = collection.FindOne(ctx, bson.M{"email": participant.Email}).Decode(&neresult)
	if !CheckRsvp(neresult){
		log.Print("error")
		return errors.New("Meeting can not be made due to participant"), string(0)

	}
	if err != nil {

		participant := &Participant{
			ID: primitive.NewObjectID(),
			Name: participant.Name,
Email: participant.Email,
RSVP: "yes",
Meeting:participant.Meeting,
		}
		result,_ := collection.InsertOne(ctx, participant)
		log.Print(result,"made")
return nil,"success"
	}


	resultUpdate, err := collection.UpdateOne(
		ctx,
		bson.M{"email": participant.Email},
		bson.M{
			"$set": bson.M{
				"rsvp":"yes",

			},
			//"$push":bson.M{
			//
			//},
		},
	)
	log.Print(resultUpdate,"updated")
return nil, "success"
client, err := mongo.NewClient(options.Client().ApplyURI("<mongodb+srv://mohit:<mohit143kumath>@cluster0.9yilk.mongodb.net/<dbname>?retryWrites=true&w=majority>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}
databases := `
    INSERT INTO meeting (title, starttime, endtime, TimeofScheduling)
    VALUES ('abc', 921, 1100, 1500)`
    idRes, err = db.Exec(databases)
    if err != nil {
      Fatal(err)
    }
    Meetings.Id = idRes;

func handleRequests() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/meeting", scheduleMeeting)
	//meetingRouter.HandleFunc("/meeting/{title}/{participants=[]}/{start}/{end}", scheduleMeetingDetails)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {

	/*CREATE TABLE meeting (
	  id SERIAL PRIMARY KEY,
	  title VARCHAR,
	  starttime INT,
	  endtime INT,
	  TimeofScheduling INT
	);*/

	Meetings = Meeting{Id: 1, Title: "ABC", Participants: []string{"mohitkumath@gmail.com", "rishi@yahoo.com"}, StartTime: 920, EndTime: 1130}

	//Meeting with start time greater than end time
	//Meetings  = Meeting{Id: 1,Title: "ABC", Participants:[]string{"shashank@gmail.com","rishi@yahoo.com"},StartTime:1100,EndTime:930}

	//Meeting with no title
	//Meetings  = Meeting{Id: 1,Title: "", Participants:[]string{"shashank@gmail.com","rishi@yahoo.com"},StartTime:920,EndTime:1130}
	fmt.Println("Starting the application...")


	http.HandleFunc("/shedulemeeting",SheduleMeeting)
	http.HandleFunc("/getmeeting", GetMeeting)
	http.HandleFunc("/getparticipantmeeting", MeetingOfParticipant)
	handleRequests()
}
